// Package clipboard writes text to the system clipboard using native tools.
// No CGO required — works on macOS, Linux (X11/Wayland), and Windows.
package clipboard

import (
	"bytes"
	"os/exec"
	"runtime"
)

// Write copies text to the system clipboard.
// Returns nil on success, or an error if no clipboard tool is available.
func Write(text string) error {
	cmd := clipboardCmd()
	if cmd == nil {
		return nil // no clipboard available — silent no-op
	}
	cmd.Stdin = bytes.NewBufferString(text)
	return cmd.Run()
}

func clipboardCmd() *exec.Cmd {
	switch runtime.GOOS {
	case "darwin":
		return exec.Command("pbcopy")
	case "windows":
		return exec.Command("clip")
	default: // linux and others
		// Prefer Wayland, fall back to X11
		if path, err := exec.LookPath("wl-copy"); err == nil {
			return exec.Command(path)
		}
		if path, err := exec.LookPath("xclip"); err == nil {
			return exec.Command(path, "-selection", "clipboard")
		}
		if path, err := exec.LookPath("xsel"); err == nil {
			return exec.Command(path, "--clipboard", "--input")
		}
		return nil
	}
}
