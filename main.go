package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"infographic-generator/internal/browser"
	"infographic-generator/internal/project"
	"infographic-generator/internal/server"
)

// version est injectée au build via -ldflags "-X main.version=..." ;
// la valeur ci-dessous sert de repli pour un build sans script (go build / go run).
var version = "2.3.0"

type lockInfo struct {
	PID  int    `json:"pid"`
	Port int    `json:"port"`
	URL  string `json:"url"`
}

func main() {
	port := flag.Int("port", 0, "Port HTTP (0 = auto)")
	dataDir := flag.String("data", "", "Repertoire de donnees (defaut: ./data)")
	noKiosk := flag.Bool("no-kiosk", false, "Ne pas ouvrir le navigateur en mode application")
	noBrowser := flag.Bool("no-browser", false, "Ne pas ouvrir le navigateur du tout")
	showVersion := flag.Bool("version", false, "Afficher la version")
	browserPath := flag.String("browser", "", "Chemin explicite du navigateur Chromium (Chrome/Edge/Brave...) ; sinon detection automatique")
	flag.Parse()

	if *showVersion {
		fmt.Printf("Generateur d'Infographie v%s\n", version)
		os.Exit(0)
	}

	// Resolve data directory
	dd := *dataDir
	if dd == "" {
		exe, err := os.Executable()
		if err != nil {
			dd = "./data"
		} else {
			dd = filepath.Join(filepath.Dir(exe), "data")
		}
	}
	os.MkdirAll(dd, 0755)
	lockPath := filepath.Join(dd, ".lock")

	// ============================================
	// SINGLE INSTANCE CHECK
	// ============================================
	if existingURL := checkExistingInstance(lockPath); existingURL != "" {
		log.Println("Instance deja en cours d'execution.")
		log.Printf("Ouverture de %s", existingURL)
		br := browser.Detect()
		if !*noBrowser {
			if *noKiosk || br.Path == "" {
				browser.OpenURL(existingURL)
			} else {
				browser.OpenURL(existingURL) // Don't kiosk, just open a tab
			}
		}
		os.Exit(0)
	}

	// ============================================
	// NORMAL STARTUP
	// ============================================
	store, err := project.NewStorage(dd)
	if err != nil {
		log.Fatalf("Erreur initialisation stockage: %v", err)
	}

	br := browser.Detect()
	if *browserPath != "" {
		// Choix explicite de l'utilisateur : prioritaire sur la detection.
		br = browser.FromPath(*browserPath)
	}
	if br.Path != "" {
		log.Printf("Navigateur: %s (%s)", br.Name, br.Path)
	} else {
		log.Println("ATTENTION: Aucun navigateur Chromium detecte.")
	}

	actualPort := *port
	if actualPort == 0 {
		actualPort, err = findFreePort()
		if err != nil {
			log.Fatalf("Impossible de trouver un port libre: %v", err)
		}
	}

	addr := fmt.Sprintf("127.0.0.1:%d", actualPort)
	baseURL := fmt.Sprintf("http://%s", addr)

	// Write lock
	writeLock(lockPath, actualPort, baseURL)

	// Shutdown channel: any goroutine can trigger shutdown
	shutdownCh := make(chan string, 1)

	// Build server with shutdown hook
	srv := server.New(store, &br, baseURL)
	srv.OnShutdown = func() {
		select {
		case shutdownCh <- "shutdown-api":
		default:
		}
	}

	httpServer := &http.Server{
		Handler:      srv.Router(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 120 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Bind the port FIRST, then serve. This guarantees the port is open
	// before we launch the browser. No more race condition.
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("Impossible d'ecouter sur %s: %v", addr, err)
	}
	log.Printf("Serveur pret: %s (PID %d)", baseURL, os.Getpid())

	go func() {
		if err := httpServer.Serve(ln); err != http.ErrServerClosed {
			log.Fatalf("Erreur serveur: %v", err)
		}
	}()

	// ============================================
	// OPEN BROWSER (port is guaranteed to be bound)
	// ============================================
	useKiosk := !*noBrowser && !*noKiosk && br.Path != ""

	if useKiosk {
		// Le kiosk est lancé mais la fin de son processus n'est PAS un
		// signal de fermeture fiable : Edge (Startup Boost) délègue parfois
		// à un processus frère et notre cmd se termine immédiatement alors
		// que la fenêtre est bien ouverte. Le cycle de vie est donc piloté
		// par le battement de coeur de la page (voir plus bas).
		go func() {
			browser.OpenKioskAndWait(br.Path, baseURL, dd)
			log.Println("Processus kiosk termine (informatif ; la fermeture est pilotee par le heartbeat).")
		}()
	} else if !*noBrowser {
		browser.OpenURL(baseURL)
	}

	// ============================================
	// VEILLE D'INACTIVITE (fermeture de la fenetre)
	// ============================================
	// L'éditeur envoie un battement toutes les 3 s. Après un premier
	// battement, 15 s de silence = plus aucune fenêtre ouverte -> extinction.
	// Sans navigateur (-no-browser), le serveur vit indéfiniment.
	if !*noBrowser {
		go func() {
			const silence = 15 * time.Second
			seen := false
			t := time.NewTicker(2 * time.Second)
			defer t.Stop()
			for range t.C {
				last := srv.LastHeartbeat()
				if last.IsZero() {
					continue // fenêtre jamais ouverte : on n'éteint pas
				}
				seen = true
				if seen && time.Since(last) > silence {
					log.Println("Aucun battement de l'interface depuis 15s: fenetre fermee, extinction.")
					select {
					case shutdownCh <- "fenetre-fermee":
					default:
					}
					return
				}
			}
		}()
	}

	// ============================================
	// WAIT FOR SHUTDOWN TRIGGER
	// ============================================
	// Three possible triggers:
	// 1. Ctrl+C / SIGTERM (console or service manager)
	// 2. Browser window closed (kiosk mode)
	// 3. /api/shutdown called from the web UI
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	var reason string
	select {
	case sig := <-sigCh:
		reason = fmt.Sprintf("signal %v", sig)
	case reason = <-shutdownCh:
	}

	log.Printf("Arret: %s", reason)

	// Cleanup
	removeLock(lockPath)
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	httpServer.Shutdown(ctx)
	log.Println("Serveur arrete.")
}

// ============================================
// LOCK FILE MANAGEMENT
// ============================================

func checkExistingInstance(lockPath string) string {
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return ""
	}

	var info lockInfo
	if err := json.Unmarshal(data, &info); err != nil {
		os.Remove(lockPath)
		return ""
	}

	if !isProcessAlive(info.PID) {
		log.Printf("Lock perrime (PID %d mort). Nettoyage.", info.PID)
		os.Remove(lockPath)
		return ""
	}

	url := info.URL
	if url == "" {
		url = fmt.Sprintf("http://127.0.0.1:%d", info.Port)
	}

	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(url + "/api/browser/status")
	if err != nil {
		log.Printf("Lock perrime (serveur muet). Nettoyage.")
		os.Remove(lockPath)
		return ""
	}
	defer resp.Body.Close()

	if resp.StatusCode == 200 {
		return url
	}
	os.Remove(lockPath)
	return ""
}

func writeLock(lockPath string, port int, url string) {
	info := lockInfo{PID: os.Getpid(), Port: port, URL: url}
	data, _ := json.MarshalIndent(info, "", "  ")
	os.WriteFile(lockPath, data, 0644)
}

func removeLock(lockPath string) {
	os.Remove(lockPath)
}

func isProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))
	if err == nil {
		return true
	}
	errMsg := err.Error()
	return strings.Contains(errMsg, "Access is denied") ||
		strings.Contains(errMsg, "not supported")
}

func findFreePort() (int, error) {
	l, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return 0, err
	}
	defer l.Close()
	return l.Addr().(*net.TCPAddr).Port, nil
}
