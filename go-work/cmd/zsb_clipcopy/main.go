package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"example.com/workspace/lib/argsLib"
)

func clipcopy(arg *argsLib.ParsedWhole) error {
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
	arg, err := argsLib.ParseWhole(true)
	if err != nil || arg.Content == "" {
		fmt.Println(err)
		os.Exit(1)
	}

	clipcopy(arg)
}
