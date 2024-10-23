package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"example.com/workspace/lib/argsLib"
)

func clipcopy(args *string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("pbcopy")
	default: // Assume Unix-like
		cmd = exec.Command("xclip", "-selection", "clipboard")
	}

	cmd.Stdin = strings.NewReader(*args)
	return cmd.Run()
}

func main() {
	args, err := argsLib.ParseWhole()
	if err != nil || len(args) == 0 {
		fmt.Println(err)
		os.Exit(1)
	}

	clipcopy(&args)
}
