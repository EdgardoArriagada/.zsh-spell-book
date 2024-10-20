package main

import (
	"example.com/workspace/lib/args"
	"fmt"
	"os"
	"os/exec"
	"runtime"
)

func openURL(url string) error {
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

func main() {
	parsedArgs, err := args.ParseArguments()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	url := parsedArgs[0]

	if err := openURL(url); err != nil {
		fmt.Printf("Error opening URL %s: %v\n", url, err)
	}
}
