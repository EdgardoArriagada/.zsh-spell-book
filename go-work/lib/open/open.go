package open

import (
	"os/exec"
	"runtime"
)

func Url(url string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default: // Assume Unix-like
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Run()
}
