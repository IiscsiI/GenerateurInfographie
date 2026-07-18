//go:build windows

package browser

import (
	"os/exec"
	"syscall"
)

// hideWindow évite l'apparition d'une console parasite lors du lancement
// du navigateur d'export (l'application est compilée en -H windowsgui).
func hideWindow(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
