package browser

import (
	"embed"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

//go:embed all:profile
var profileTemplate embed.FS

// OpenKioskAndWait launches the browser in app mode and BLOCKS until the window is closed.
// Returns true if this was a fresh profile (first launch) — in that case cmd.Wait()
// may return prematurely because Edge/Chromium restarts itself during profile init.
// The caller should NOT trigger server shutdown when freshProfile is true.
func OpenKioskAndWait(browserPath, url, dataDir string) bool {
	profileDir := filepath.Join(dataDir, ".browser-profile")

	// First launch: extract the embedded profile template
	freshProfile := ensureProfile(profileDir)

	args := []string{
		"--app=" + url,
		"--window-size=1500,900",
		"--user-data-dir=" + profileDir,
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-default-apps",
		"--disable-extensions",
		"--disable-sync",
		"--disable-translate",
		// msFirstRunExperience / WelcomePage : Edge ouvre sinon une fenêtre
		// d'accueil séparée au premier lancement d'un profil, malgré
		// --no-first-run.
		"--disable-features=TranslateUI,EdgeCollections,msEdgeSidebarV2,msFirstRunExperience,msEdgeWelcomePage,msSeamlessWebToBrowserSignIn",
		"--disable-client-side-phishing-detection",
		"--disable-component-update",
		"--disable-infobars",
	}

	cmd := exec.Command(browserPath, args...)
	if err := cmd.Start(); err != nil {
		log.Printf("Impossible d'ouvrir en mode kiosk: %v", err)
		OpenURL(url)
		return freshProfile
	}

	if err := cmd.Wait(); err != nil {
		log.Printf("Navigateur ferme (code: %v)", err)
	} else {
		log.Println("Navigateur ferme normalement.")
	}

	return freshProfile
}

// ensureProfile extracts the embedded profile template if the profile doesn't exist yet.
// Returns true if the profile was just created (fresh), false if it already existed.
func ensureProfile(profileDir string) bool {
	marker := filepath.Join(profileDir, "First Run")
	if _, err := os.Stat(marker); err == nil {
		return false // Profile already exists
	}

	log.Println("Initialisation du profil navigateur...")
	os.MkdirAll(profileDir, 0755)

	err := fs.WalkDir(profileTemplate, "profile", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel("profile", path)
		dest := filepath.Join(profileDir, rel)

		if d.IsDir() {
			return os.MkdirAll(dest, 0755)
		}

		data, err := profileTemplate.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(dest, data, 0644)
	})

	if err != nil {
		log.Printf("Erreur extraction profil: %v", err)
	} else {
		log.Println("Profil navigateur pret.")
	}

	return true
}

// OpenURL opens a URL in the default system browser (non-blocking).
func OpenURL(url string) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "windows":
		cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
	case "darwin":
		cmd = exec.Command("open", url)
	default:
		cmd = exec.Command("xdg-open", url)
	}

	if err := cmd.Start(); err != nil {
		log.Printf("Impossible d'ouvrir le navigateur: %v", err)
		log.Printf("Ouvrez manuellement: %s", url)
	}
}
