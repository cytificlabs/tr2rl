package clipboard

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

// ReadAll returns the current text content of the system clipboard.
// It uses native OS commands to avoid CGO dependencies.
func ReadAll() (string, error) {
	switch runtime.GOOS {
	case "windows":
		// PowerShell Get-Clipboard is reliable on modern Windows, but we must ensure encoding is UTF8
		// otherwise Go executables might receive CP1252 or UTF-16.
		cmd := exec.Command("powershell", "-NoProfile", "-Command", "$OutputEncoding = [Console]::OutputEncoding = [System.Text.Encoding]::UTF8; Get-Clipboard")
		out, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("clipboard read failed: %w", err)
		}
		// Windows allows formatting; we want raw text usually, but simple text is fine.
		// Remove BOM if present (EF BB BF)
		text := string(out)
		text = strings.TrimPrefix(text, "\uFEFF")
		return strings.ReplaceAll(text, "\r\n", "\n"), nil

	case "darwin":
		// pbpaste is standard on macOS
		cmd := exec.Command("pbpaste")
		out, err := cmd.Output()
		if err != nil {
			return "", fmt.Errorf("clipboard read failed: %w", err)
		}
		return string(out), nil

	case "linux":
		// Try xclip first, then xsel
		if _, err := exec.LookPath("xclip"); err == nil {
			cmd := exec.Command("xclip", "-selection", "clipboard", "-o")
			out, err := cmd.Output()
			return string(out), err
		}
		if _, err := exec.LookPath("xsel"); err == nil {
			cmd := exec.Command("xsel", "--clipboard", "--output")
			out, err := cmd.Output()
			return string(out), err
		}
		return "", fmt.Errorf("no clipboard tool found (install xclip or xsel)")

	default:
		return "", fmt.Errorf("clipboard not supported on %s", runtime.GOOS)
	}
}
