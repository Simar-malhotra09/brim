package utils 
import (
	"os/exec"
	"runtime"
	tea "charm.land/bubbletea/v2"
)

func OpenURL(url string) tea.Cmd {
	return func() tea.Msg {
		var cmd *exec.Cmd
		switch runtime.GOOS {
		case "darwin":
			cmd = exec.Command("open", url)
		case "linux":
			cmd = exec.Command("xdg-open", url)
		case "windows":
			cmd = exec.Command("rundll32", "url.dll,FileProtocolHandler", url)
		}
		_ = cmd.Start()  // fire and forget; don't Wait
		return nil       // no Msg back
	}
}
