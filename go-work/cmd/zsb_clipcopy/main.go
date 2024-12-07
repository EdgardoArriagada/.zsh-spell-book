package main

import (
	"os/exec"
	"runtime"
	"strings"

	"example.com/workspace/lib/args"
	u "example.com/workspace/lib/utils"
)

func copyToClipboard(arg *args.ParsedWhole) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	default: // Assume Unix-like
		cmd = exec.Command("xclip", "-selection", "clipboard")
	}

	cmd.Stdin = strings.NewReader(arg.Content)
	return cmd.Run()
}

func main() {
	d := u.Must(args.ParseWholeWithStdin())

	copyToClipboard(d)
}
