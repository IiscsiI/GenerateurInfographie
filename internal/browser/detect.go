package browser

import (
	"context"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"
)

// Browser represents a detected Chromium-based browser.
type Browser struct {
	Name string
	Path string
}

// Detect scans the system for available Chromium-based browsers.
// Priority: Edge > Chrome > Chromium (Edge is always present on modern Windows).
//
// Chaque candidat est validé par un lancement réel (--version) : sur Ubuntu,
// /usr/bin/chromium-browser est souvent un stub snap qui existe sur disque
// mais échoue à l'exécution. Sans cette sonde, tous les exports échoueraient
// silencieusement alors qu'un autre navigateur fonctionnel est disponible.
func Detect() Browser {
	candidates := getCandidates()
	for _, c := range candidates {
		if _, err := os.Stat(c.Path); err == nil && probeBrowser(c.Path) {
			return c
		}
	}
	// Fallback: try PATH lookup
	for _, name := range []string{"msedge", "google-chrome", "google-chrome-stable", "chromium-browser", "chromium"} {
		if p, err := exec.LookPath(name); err == nil && probeBrowser(p) {
			return Browser{Name: name, Path: p}
		}
	}
	return Browser{}
}

// FromPath construit un Browser depuis un chemin fourni par l'utilisateur
// (flag -browser). Le chemin est sondé ; s'il ne s'exécute pas, on le garde
// quand même (l'utilisateur a explicitement choisi) mais on le signale.
func FromPath(path string) Browser {
	name := filepath.Base(path)
	if !probeBrowser(path) {
		log.Printf("ATTENTION: le navigateur indique ne repond pas a --version: %s", path)
	}
	return Browser{Name: name, Path: path}
}

// probeBrowser vérifie que le binaire s'exécute réellement.
//
// UNIQUEMENT hors Windows : sur Linux, /usr/bin/chromium-browser est parfois
// un stub snap qui existe mais ne s'exécute pas — la sonde --version l'écarte.
// Sous Windows en revanche, msedge.exe/chrome.exe sont des applications GUI :
// "--version" n'affiche rien et OUVRE LE NAVIGATEUR (fenêtre d'accueil
// parasite au démarrage de l'application). L'existence du fichier suffit,
// les chemins d'installation Windows étant fiables.
func probeBrowser(path string) bool {
	if runtime.GOOS == "windows" {
		return true // le os.Stat de l'appelant a déjà validé l'existence
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	return exec.CommandContext(ctx, path, "--version").Run() == nil
}

func getCandidates() []Browser {
	switch runtime.GOOS {
	case "windows":
		return []Browser{
			{Name: "Microsoft Edge", Path: `C:\Program Files (x86)\Microsoft\Edge\Application\msedge.exe`},
			{Name: "Microsoft Edge", Path: `C:\Program Files\Microsoft\Edge\Application\msedge.exe`},
			{Name: "Google Chrome", Path: `C:\Program Files\Google\Chrome\Application\chrome.exe`},
			{Name: "Google Chrome", Path: `C:\Program Files (x86)\Google\Chrome\Application\chrome.exe`},
			{Name: "Chromium", Path: `C:\Program Files\Chromium\Application\chromium.exe`},
		}
	case "linux":
		return []Browser{
			{Name: "Chromium", Path: "/usr/bin/chromium-browser"},
			{Name: "Chromium", Path: "/usr/bin/chromium"},
			{Name: "Google Chrome", Path: "/usr/bin/google-chrome"},
			{Name: "Google Chrome", Path: "/usr/bin/google-chrome-stable"},
			{Name: "Microsoft Edge", Path: "/usr/bin/microsoft-edge"},
			{Name: "Microsoft Edge", Path: "/usr/bin/microsoft-edge-stable"},
		}
	case "darwin":
		return []Browser{
			{Name: "Google Chrome", Path: "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"},
			{Name: "Microsoft Edge", Path: "/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge"},
			{Name: "Chromium", Path: "/Applications/Chromium.app/Contents/MacOS/Chromium"},
		}
	}
	return nil
}
