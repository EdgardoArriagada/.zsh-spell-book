package main

import (
	"os/exec"
	"runtime"
	"strings"

	"example.com/workspace/lib/args"
	u "example.com/workspace/lib/utils"
)

func main() {
	d := u.Must1(args.ParseWholeWithStdin())

	copyToClipboard(d)
}

func copyToClipboard(d *args.ParsedWhole) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	default: // Assume Unix-like
		cmd = exec.Command("xclip", "-selection", "clipboard")
	}

	cmd.Stdin = strings.NewReader(d.Content)
	return cmd.Run()
}
