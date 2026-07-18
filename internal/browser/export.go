package browser

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	_ "image/jpeg"
	"log"
	"math"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	cdbrowser "github.com/chromedp/cdproto/browser"
	"github.com/chromedp/cdproto/page"
	"github.com/chromedp/chromedp"
)

// ExportConfig holds export parameters.
type ExportConfig struct {
	URL         string
	Format      string  // "A3", "A4", "A2"
	Orientation string  // "portrait", "landscape"
	DPI         float64 // 72, 150, 300
	BrowserPath string
	// ProfileDir : répertoire où créer les profils navigateur dédiés aux
	// exports (persistants : évite les échecs liés au répertoire temporaire
	// système et accélère les lancements suivants). Vide = temp chromedp.
	ProfileDir string
}

// paperSizes in inches (width, height) for portrait orientation.
var paperSizes = map[string][2]float64{
	"A4": {8.27, 11.69},
	"A3": {11.69, 16.54},
	"A2": {16.54, 23.39},
}

// paperViewport returns CSS pixel dimensions matching the exact paper ratio.
// baseWidth is the reference width; height is computed from the ratio.
func paperViewport(format, orientation string, baseWidth int) (int, int) {
	size, ok := paperSizes[format]
	if !ok {
		size = paperSizes["A3"]
	}
	w, h := size[0], size[1]
	if orientation == "landscape" {
		w, h = h, w
	}
	ratio := h / w // height/width ratio
	vpH := int(math.Round(float64(baseWidth) * ratio))
	return baseWidth, vpH
}

// buildRenderURL appends format and orientation query params.
// ==========================================================================
// LANCEMENT DU NAVIGATEUR HEADLESS : STRATÉGIES MULTIPLES
// --------------------------------------------------------------------------
// Sur certains postes (Edge sous Windows notamment), le lancement headless
// peut échouer silencieusement selon la version du navigateur, la présence
// d'une instance déjà ouverte, les politiques d'entreprise ou l'antivirus.
// Plutôt qu'un échec sec, on tente successivement :
//   1. headless "new"       (mode moderne, Chrome/Edge >= 112)
//   2. headless classique    (flag booléen historique)
//   3. fenêtre hors écran    (pas de headless du tout : si le kiosk
//                             fonctionne, cette stratégie fonctionne)
// Chaque échec est journalisé et TOUTES les causes sont remontées dans
// l'erreur finale — plus jamais de "chrome failed to start:" muet.
// ==========================================================================

type launchStrategy struct {
	name string
	args []string
}

func launchStrategies() []launchStrategy {
	return []launchStrategy{
		{name: "headless nouveau", args: []string{"--headless=new"}},
		{name: "headless classique", args: []string{"--headless"}},
		// Dernier recours : fenêtre réelle hors écran. Si le kiosk s'ouvre,
		// ceci s'ouvre. PAS de minimisation (fenêtre minimisée = parfois
		// jamais rendue, même avec les flags anti-throttling).
		{name: "fenetre hors ecran", args: []string{"--window-position=-32000,-32000"}},
	}
}

// baseArgs : arguments communs de lancement.
func baseArgs(profileDir string, debugPort, vpWidth int) []string {
	return []string{
		"--user-data-dir=" + profileDir,
		// Port EXPLICITE, jamais 0 : certaines versions d'Edge n'annoncent
		// pas le port choisi sur stderr avec --remote-debugging-port=0,
		// ce qui rend le lancement indétectable pour un superviseur.
		fmt.Sprintf("--remote-debugging-port=%d", debugPort),
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-gpu",
		"--no-sandbox",
		"--disable-dev-shm-usage",
		"--disable-extensions",
		"--disable-background-networking",
		"--disable-component-update",
		"--disable-sync",
		"--disable-breakpad",
		"--mute-audio",
		// Anti-throttling : indispensable pour la stratégie fenêtre (une
		// fenêtre occultée voit son renderer suspendu par Chromium).
		"--disable-backgrounding-occluded-windows",
		"--disable-renderer-backgrounding",
		"--disable-background-timer-throttling",
		// Écrans bloquants sur profil vierge (choix du moteur de recherche
		// Chrome >= 127, accueil Edge, bulles de restauration...)
		"--disable-search-engine-choice-screen",
		"--disable-session-crashed-bubble",
		"--hide-crash-restore-bubble",
		"--disable-features=msFirstRunExperience,msEdgeWelcomePage,Translate",
		fmt.Sprintf("--window-size=%d,900", vpWidth),
		"about:blank",
	}
}

// runWithStrategies exécute les actions chromedp en essayant chaque
// stratégie de lancement jusqu'à la première qui fonctionne.
func runWithStrategies(parent context.Context, cfg ExportConfig, vpWidth int, label string, actions ...chromedp.Action) error {
	var attempts []string

	for i, st := range launchStrategies() {
		// Profil dédié par stratégie : un verrou laissé par un crash
		// n'empoisonne pas les stratégies suivantes.
		profileDir := filepath.Join(os.TempDir(), fmt.Sprintf("ig-export-%d", i))
		if cfg.ProfileDir != "" {
			profileDir = filepath.Join(cfg.ProfileDir, fmt.Sprintf("export-%d", i))
		}
		if err := os.MkdirAll(profileDir, 0755); err != nil {
			attempts = append(attempts, fmt.Sprintf("[%s] creation du profil impossible: %v", st.name, err))
			continue
		}

		err := runOnce(parent, cfg.BrowserPath, profileDir, st.args, vpWidth, actions)
		if err == nil {
			if i > 0 {
				log.Printf("Export %s: reussi via la strategie %q", label, st.name)
			}
			return nil
		}

		msg := strings.TrimSpace(err.Error())
		attempts = append(attempts, fmt.Sprintf("[%s] %s", st.name, msg))
		log.Printf("Export %s: strategie %q en echec: %s", label, st.name, msg)

		if parent.Err() != nil {
			break // requête annulée : inutile d'insister
		}
	}

	browserName := cfg.BrowserPath
	if browserName == "" {
		browserName = "navigateur systeme"
	}
	return fmt.Errorf(
		"impossible de lancer le navigateur pour l'export (%s).\nTentatives:\n  %s\nPistes: fermer toutes les fenetres du navigateur puis reessayer ; verifier antivirus / AppLocker ; ou designer un autre navigateur au lancement: infographic-generator.exe -browser \"C:\\Program Files\\Google\\Chrome\\Application\\chrome.exe\"",
		browserName, strings.Join(attempts, "\n  "))
}

// runOnce : UN lancement complet, piloté par nos soins.
//
// Pourquoi ne pas laisser chromedp lancer le navigateur ? Son ExecAllocator
// dépend de la ligne "DevTools listening on ws://..." que le navigateur écrit
// sur stderr. Or Edge sous Windows omet parfois cette ligne : le navigateur
// démarre (fenêtre about:blank visible) mais chromedp conclut à un
// "chrome failed to start" muet. Ici : port explicite + interrogation HTTP
// du endpoint DevTools — peu importe ce que le navigateur écrit ou n'écrit pas.
func runOnce(parent context.Context, browserPath, profileDir string, strategyArgs []string, vpWidth int, actions []chromedp.Action) error {
	debugPort, err := freeTCPPort()
	if err != nil {
		return fmt.Errorf("aucun port libre pour le debogage: %w", err)
	}

	args := append(baseArgs(profileDir, debugPort, vpWidth), strategyArgs...)

	launchCtx, cancelLaunch := context.WithCancel(parent)
	defer cancelLaunch()

	cmd := exec.CommandContext(launchCtx, browserPath, args...)
	var output bytes.Buffer
	cmd.Stdout = &output
	cmd.Stderr = &output
	hideWindow(cmd)

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("lancement du processus: %w", err)
	}
	// Toujours récolter le processus (pas de zombie) et le tuer en sortie.
	procDone := make(chan error, 1)
	go func() { procDone <- cmd.Wait() }()
	defer func() {
		cancelLaunch() // tue notre processus via CommandContext (si encore là)
		select {
		case <-procDone:
		case <-time.After(5 * time.Second):
			if cmd.Process != nil {
				cmd.Process.Kill()
			}
		}
	}()

	// Attendre que le endpoint DevTools réponde (ou que le processus meure).
	wsURL, err := waitForDevTools(parent, debugPort, profileDir, procDone, &output)
	if err != nil {
		return err
	}

	allocCtx, allocCancel := chromedp.NewRemoteAllocator(parent, wsURL, chromedp.NoModifyURL)
	defer allocCancel()
	taskCtx, taskCancel := chromedp.NewContext(allocCtx)
	defer taskCancel()

	// Fermeture PROPRE du navigateur en fin d'export, via DevTools.
	// Indispensable avec Edge : quand le processus lancé délègue à un frère
	// (Startup Boost), tuer notre cmd ne ferme pas le vrai navigateur — sans
	// cette commande, des processus Edge invisibles s'accumuleraient.
	defer func() {
		_ = chromedp.Run(taskCtx, chromedp.ActionFunc(func(ctx context.Context) error {
			c, cancel := context.WithTimeout(ctx, 2*time.Second)
			defer cancel()
			return cdbrowser.Close().Do(c)
		}))
	}()

	runCtx, timeoutCancel := context.WithTimeout(taskCtx, 45*time.Second)
	defer timeoutCancel()

	return chromedp.Run(runCtx, actions...)
}

// waitForDevTools attend l'URL websocket DevTools du navigateur (15 s max).
//
// Deux sources sont surveillées :
//  1. le port explicite passé en argument (--remote-debugging-port)
//  2. le fichier DevToolsActivePort que Chromium écrit dans le profil
//     (utile si le navigateur a retenu un autre port)
//
// CAS EDGE (Startup Boost) : le msedge.exe que nous lançons peut DÉLÉGUER
// à un processus frère puis se terminer avec le code 0. Cette sortie
// immédiate n'est PAS un échec : le navigateur délégué, qui a reçu nos
// arguments, va ouvrir le port. On ne s'arrête donc sur une fin de
// processus que si le code de sortie est non nul.
func waitForDevTools(ctx context.Context, port int, profileDir string, procDone <-chan error, output *bytes.Buffer) (string, error) {
	client := &http.Client{Timeout: 2 * time.Second}
	deadline := time.After(15 * time.Second)
	tick := time.NewTicker(250 * time.Millisecond)
	defer tick.Stop()

	delegated := false

	browserOut := func() string {
		s := strings.TrimSpace(output.String())
		if s == "" {
			return "(aucune sortie du navigateur)"
		}
		if len(s) > 600 {
			s = s[:600] + "..."
		}
		return s
	}

	tryPort := func(p int) string {
		resp, err := client.Get(fmt.Sprintf("http://127.0.0.1:%d/json/version", p))
		if err != nil {
			return ""
		}
		defer resp.Body.Close()
		var v struct {
			WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
		}
		if json.NewDecoder(resp.Body).Decode(&v) == nil {
			return v.WebSocketDebuggerURL
		}
		return ""
	}

	for {
		select {
		case <-ctx.Done():
			return "", ctx.Err()

		case err := <-procDone:
			if err != nil {
				// Vrai échec : le navigateur a refusé de démarrer.
				return "", fmt.Errorf("le navigateur s'est arrete avant d'ouvrir son port de debogage (%v). Sortie: %s", err, browserOut())
			}
			// Code 0 : très probablement une délégation (Edge Startup Boost).
			// Le navigateur réel va ouvrir le port : on continue d'attendre.
			delegated = true
			procDone = nil // ne plus lire ce canal

		case <-deadline:
			hint := ""
			if delegated {
				hint = " (le processus lance a delegue a une instance existante sans transmettre le debogage)"
			}
			return "", fmt.Errorf("le port de debogage %d n'a jamais repondu en 15s%s. Sortie: %s", port, hint, browserOut())

		case <-tick.C:
			if ws := tryPort(port); ws != "" {
				return ws, nil
			}
			// Source alternative : le fichier écrit par Chromium dans le profil
			if data, err := os.ReadFile(filepath.Join(profileDir, "DevToolsActivePort")); err == nil {
				lines := strings.SplitN(strings.TrimSpace(string(data)), "\n", 2)
				if p, convErr := strconv.Atoi(strings.TrimSpace(lines[0])); convErr == nil && p != port {
					if ws := tryPort(p); ws != "" {
						return ws, nil
					}
				}
			}
		}
	}
}

// freeTCPPort réserve puis libère un port TCP local.
func freeTCPPort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}

func buildRenderURL(baseURL, format, orientation string) string {
	sep := "?"
	if strings.Contains(baseURL, "?") {
		sep = "&"
	}
	return fmt.Sprintf("%s%sformat=%s&orientation=%s", baseURL, sep, format, orientation)
}

// ExportPNG captures a full-page screenshot with layout matching the paper ratio.
func ExportPNG(ctx context.Context, cfg ExportConfig) ([]byte, error) {
	dpi := cfg.DPI
	if dpi == 0 {
		dpi = 150
	}
	scale := dpi / 96.0
	if scale < 1 {
		scale = 1
	}

	orientation := cfg.Orientation
	if orientation == "" {
		orientation = "portrait"
	}

	baseWidth := 1200
	if orientation == "landscape" {
		baseWidth = 1600
	}
	vpW, vpH := paperViewport(cfg.Format, orientation, baseWidth)
	// fit=1 : la page render met le contenu à l'échelle du viewport après
	// chargement des images. Garantie : rien n'est tronqué, quel que soit
	// le volume de contenu par rapport au ratio papier.
	renderURL := buildRenderURL(cfg.URL, cfg.Format, orientation) + "&fit=1"

	var buf []byte

	log.Printf("Export PNG: %s (%s %s %gDPI, viewport %dx%d)", renderURL, cfg.Format, orientation, dpi, vpW, vpH)
	err := runWithStrategies(ctx, cfg, vpW, "PNG",
		// Le viewport doit être émulé AVANT le rendu : le calcul de fit
		// se base sur window.innerWidth/innerHeight.
		chromedp.EmulateViewport(int64(vpW), int64(vpH), chromedp.EmulateScale(scale)),
		chromedp.Navigate(renderURL),
		// Le marqueur ready n'apparaît qu'après : rendu + images chargées + fit.
		chromedp.WaitVisible("#infographic-ready", chromedp.ByID),
		chromedp.Sleep(200*time.Millisecond),
		// Capture exactly the viewport (not beyond) to match paper ratio
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			buf, err = page.CaptureScreenshot().
				WithFormat(page.CaptureScreenshotFormatPng).
				WithCaptureBeyondViewport(false).
				WithFromSurface(true).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("export PNG: %w", err)
	}

	log.Printf("Export PNG: OK (%d bytes)", len(buf))
	return buf, nil
}

// ExportPDF produit un PDF VECTORIEL via Page.printToPDF (texte réel,
// sélectionnable, net à toute échelle, fichiers légers). En cas d'échec,
// bascule automatique sur l'ancien pipeline raster (capture JPEG encapsulée).
func ExportPDF(ctx context.Context, cfg ExportConfig) ([]byte, error) {
	data, err := exportPDFVector(ctx, cfg)
	if err == nil {
		return data, nil
	}
	log.Printf("Export PDF vectoriel en echec (%v), bascule sur le mode raster", err)
	return exportPDFRaster(ctx, cfg)
}

// exportPDFVector : rendu print natif de Chromium (Page.printToPDF).
//
// Stratégie : reproduire EXACTEMENT le pipeline PNG (même viewport au ratio
// papier, même mise à l'échelle fit=1 côté page), puis imprimer avec
// scale = largeurPapier / largeurViewport. Chrome compose alors la page
// imprimée à la même largeur CSS que le viewport de capture : le PDF est
// pixel-identique au PNG, mais vectoriel, et tient sur une page exactement.
func exportPDFVector(parent context.Context, cfg ExportConfig) ([]byte, error) {
	size, ok := paperSizes[cfg.Format]
	if !ok {
		size = paperSizes["A3"]
	}
	paperWIn, paperHIn := size[0], size[1]
	orientation := cfg.Orientation
	if orientation == "" {
		orientation = "portrait"
	}
	if orientation == "landscape" {
		paperWIn, paperHIn = paperHIn, paperWIn
	}

	baseWidth := 1200
	if orientation == "landscape" {
		baseWidth = 1600
	}
	vpW, vpH := paperViewport(cfg.Format, orientation, baseWidth)
	renderURL := buildRenderURL(cfg.URL, cfg.Format, orientation) + "&fit=1"

	// Facteur d'impression : largeur papier (96 px CSS / pouce) / largeur viewport.
	// printToPDF n'accepte que [0.1 ; 2] — tous nos formats A2/A3/A4 y tiennent.
	scale := paperWIn * 96.0 / float64(vpW)
	if scale < 0.1 || scale > 2 {
		return nil, fmt.Errorf("scale d'impression hors bornes: %.3f", scale)
	}

	log.Printf("Export PDF (vectoriel): %s (%s %s, viewport %dx%d, scale %.3f)",
		renderURL, cfg.Format, orientation, vpW, vpH, scale)

	var pdf []byte
	if err := runWithStrategies(parent, cfg, vpW, "PDF",
		chromedp.EmulateViewport(int64(vpW), int64(vpH)),
		chromedp.Navigate(renderURL),
		// Le marqueur ready n'apparaît qu'après : rendu + images + fit.
		chromedp.WaitVisible("#infographic-ready", chromedp.ByID),
		chromedp.Sleep(200*time.Millisecond),
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			pdf, _, err = page.PrintToPDF().
				WithPaperWidth(paperWIn).
				WithPaperHeight(paperHIn).
				WithMarginTop(0).
				WithMarginBottom(0).
				WithMarginLeft(0).
				WithMarginRight(0).
				WithPrintBackground(true).
				WithScale(scale).
				WithPreferCSSPageSize(false).
				Do(ctx)
			return err
		}),
	); err != nil {
		return nil, fmt.Errorf("printToPDF: %w", err)
	}

	log.Printf("Export PDF (vectoriel): OK (%d octets)", len(pdf))
	return pdf, nil
}

// exportPDFRaster : ancien pipeline (capture JPEG encapsulée dans un PDF
// mono-page). Conservé comme solution de secours.
func exportPDFRaster(ctx context.Context, cfg ExportConfig) ([]byte, error) {
	size, ok := paperSizes[cfg.Format]
	if !ok {
		size = paperSizes["A3"]
	}

	paperWidthIn := size[0]
	paperHeightIn := size[1]
	orientation := cfg.Orientation
	if orientation == "" {
		orientation = "portrait"
	}
	if orientation == "landscape" {
		paperWidthIn, paperHeightIn = paperHeightIn, paperWidthIn
	}

	pageW := paperWidthIn * 72.0 // PDF points
	pageH := paperHeightIn * 72.0

	baseWidth := 1200
	if orientation == "landscape" {
		baseWidth = 1600
	}
	vpW, vpH := paperViewport(cfg.Format, orientation, baseWidth)
	renderURL := buildRenderURL(cfg.URL, cfg.Format, orientation) + "&fit=1"
	scale := 2.0

	var jpegBuf []byte

	log.Printf("Export PDF (raster): %s (%s %s, viewport %dx%d)", renderURL, cfg.Format, orientation, vpW, vpH)
	err := runWithStrategies(ctx, cfg, vpW, "PDF raster",
		chromedp.EmulateViewport(int64(vpW), int64(vpH), chromedp.EmulateScale(scale)),
		chromedp.Navigate(renderURL),
		chromedp.WaitVisible("#infographic-ready", chromedp.ByID),
		chromedp.Sleep(200*time.Millisecond),
		// Capture exactly the viewport — matches paper ratio
		chromedp.ActionFunc(func(ctx context.Context) error {
			var err error
			jpegBuf, err = page.CaptureScreenshot().
				WithFormat(page.CaptureScreenshotFormatJpeg).
				WithQuality(95).
				WithCaptureBeyondViewport(false).
				WithFromSurface(true).
				Do(ctx)
			return err
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("export PDF screenshot: %w", err)
	}

	imgCfg, _, err := image.DecodeConfig(bytes.NewReader(jpegBuf))
	if err != nil {
		return nil, fmt.Errorf("export PDF decode: %w", err)
	}

	log.Printf("Export PDF: screenshot %dx%d, page %.0fx%.0f pt", imgCfg.Width, imgCfg.Height, pageW, pageH)

	// Image fills the entire page — no aspect ratio adjustment needed
	// because viewport was set to exact paper ratio
	pdfBytes := buildPDFWithJPEG(jpegBuf, imgCfg.Width, imgCfg.Height, pageW, pageH)

	log.Printf("Export PDF: OK (%d bytes)", len(pdfBytes))
	return pdfBytes, nil
}

// buildPDFWithJPEG creates a minimal single-page PDF with the JPEG stretched to fill the page.
func buildPDFWithJPEG(jpegData []byte, imgW, imgH int, pageW, pageH float64) []byte {
	var buf bytes.Buffer
	offsets := make([]int, 7)

	// Image fills entire page (viewport already matches paper ratio)
	contentStream := fmt.Sprintf("q\n%.4f 0 0 %.4f 0 0 cm\n/Img1 Do\nQ\n",
		pageW, pageH)

	buf.WriteString("%PDF-1.4\n%\xe2\xe3\xcf\xd3\n")

	offsets[1] = buf.Len()
	buf.WriteString("1 0 obj\n<< /Type /Catalog /Pages 2 0 R >>\nendobj\n")

	offsets[2] = buf.Len()
	buf.WriteString("2 0 obj\n<< /Type /Pages /Kids [3 0 R] /Count 1 >>\nendobj\n")

	offsets[3] = buf.Len()
	buf.WriteString(fmt.Sprintf("3 0 obj\n<< /Type /Page /Parent 2 0 R /MediaBox [0 0 %.4f %.4f] /Contents 4 0 R /Resources << /XObject << /Img1 5 0 R >> >> >>\nendobj\n",
		pageW, pageH))

	offsets[4] = buf.Len()
	buf.WriteString(fmt.Sprintf("4 0 obj\n<< /Length %d >>\nstream\n", len(contentStream)))
	buf.WriteString(contentStream)
	buf.WriteString("endstream\nendobj\n")

	offsets[5] = buf.Len()
	buf.WriteString(fmt.Sprintf("5 0 obj\n<< /Type /XObject /Subtype /Image /Width %d /Height %d /ColorSpace /DeviceRGB /BitsPerComponent 8 /Filter /DCTDecode /Length %d >>\nstream\n",
		imgW, imgH, len(jpegData)))
	buf.Write(jpegData)
	buf.WriteString("\nendstream\nendobj\n")

	offsets[6] = buf.Len()
	buf.WriteString("6 0 obj\n<< /Producer (Generateur Infographie v2.0) /Creator (Generateur Infographie) >>\nendobj\n")

	xrefOffset := buf.Len()
	buf.WriteString("xref\n")
	buf.WriteString(fmt.Sprintf("0 %d\n", len(offsets)))
	buf.WriteString("0000000000 65535 f \n")
	for i := 1; i < len(offsets); i++ {
		buf.WriteString(fmt.Sprintf("%010d 00000 n \n", offsets[i]))
	}
	buf.WriteString(fmt.Sprintf("trailer\n<< /Size %d /Root 1 0 R /Info 6 0 R >>\n", len(offsets)))
	buf.WriteString(fmt.Sprintf("startxref\n%d\n%%%%EOF\n", xrefOffset))

	return buf.Bytes()
}
