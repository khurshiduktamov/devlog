// Package clipboard copies text to the system clipboard using the
// platform's standard command-line tool, so devlog needs no extra
// dependencies.
package clipboard

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

// tool is a clipboard command that reads the text to copy from stdin.
type tool struct {
	name string
	args []string
}

// candidates returns the clipboard tools to try, in order, for goos.
// wayland reports whether a Wayland session is running (Linux only).
func candidates(goos string, wayland bool) []tool {
	switch goos {
	case "darwin":
		return []tool{{"pbcopy", nil}}
	case "windows":
		// PowerShell's Set-Clipboard handles Unicode correctly; clip.exe
		// can garble non-ASCII characters such as "•".
		ps := `[Console]::InputEncoding = [Text.Encoding]::UTF8; Set-Clipboard -Value ([Console]::In.ReadToEnd())`
		return []tool{
			{"powershell.exe", []string{"-NoProfile", "-NonInteractive", "-Command", ps}},
			{"clip.exe", nil},
		}
	default: // linux, *bsd
		x11 := []tool{
			{"xclip", []string{"-selection", "clipboard"}},
			{"xsel", []string{"--clipboard", "--input"}},
		}
		if wayland {
			return append([]tool{{"wl-copy", nil}}, x11...)
		}
		return x11
	}
}

// Copy puts text on the system clipboard.
func Copy(text string) error {
	tools := candidates(runtime.GOOS, os.Getenv("WAYLAND_DISPLAY") != "")

	for _, t := range tools {
		path, err := exec.LookPath(t.name)
		if err != nil {
			continue
		}
		cmd := exec.Command(path, t.args...)
		cmd.Stdin = strings.NewReader(text)
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%s failed: %v %s", t.name, err, strings.TrimSpace(string(out)))
		}
		return nil
	}

	names := make([]string, len(tools))
	for i, t := range tools {
		names[i] = t.name
	}
	return errors.New("no clipboard tool found (install one of: " + strings.Join(names, ", ") + ")")
}
