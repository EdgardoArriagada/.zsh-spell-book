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
	default: // Assume Unix-like
		cmd = exec.Command("xdg-open", url)
	}

	return cmd.Run()
}

func File(filename string) error {
	return Url(filename)
}
