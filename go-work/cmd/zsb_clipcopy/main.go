package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"example.com/workspace/lib/args"
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
	d, err := args.ParseWholeWithStdin()
	if err != nil || d.Content == "" {
		fmt.Println(err)
		os.Exit(1)
	}

	copyToClipboard(d)
}
