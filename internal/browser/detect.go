package browser

import (
	"os"
	"os/exec"
	"runtime"
)

// Browser represents a detected Chromium-based browser.
type Browser struct {
	Name string
	Path string
}

// Detect scans the system for available Chromium-based browsers.
// Priority: Edge > Chrome > Chromium (Edge is always present on modern Windows).
func Detect() Browser {
	candidates := getCandidates()
	for _, c := range candidates {
		if _, err := os.Stat(c.Path); err == nil {
			return c
		}
	}
	// Fallback: try PATH lookup
	for _, name := range []string{"msedge", "google-chrome", "google-chrome-stable", "chromium-browser", "chromium"} {
		if p, err := exec.LookPath(name); err == nil {
			return Browser{Name: name, Path: p}
		}
	}
	return Browser{}
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
