package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: urlopen <url>")
		os.Exit(1)
	}

	url := os.Args[1]
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin":
		cmd = exec.Command("open", url)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", url)
	default: // Assume Unix-like
		cmd = exec.Command("xdg-open", url)
	}

	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error opening URL: %v\n", err)
		os.Exit(1)
	}
}
