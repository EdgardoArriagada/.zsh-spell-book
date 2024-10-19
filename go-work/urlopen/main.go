package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
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
	if len(os.Args) > 1 {
		// If command-line arguments are provided, use them
		for _, url := range os.Args[1:] {
			if err := openURL(url); err != nil {
				fmt.Printf("Error opening URL %s: %v\n", url, err)
			}
		}
	} else {
		// If no arguments, check if there's piped input
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Data is being piped to stdin
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				url := strings.TrimSpace(scanner.Text())
				if url != "" {
					if err := openURL(url); err != nil {
						fmt.Printf("Error opening URL %s: %v\n", url, err)
					}
				}
			}
			if err := scanner.Err(); err != nil {
				fmt.Printf("Error reading input: %v\n", err)
			}
		} else {
			// No arguments and no piped input
			fmt.Println("Usage: urlopen <url> [<url2> ...]")
			fmt.Println("   or: echo <url> | urlopen")
			os.Exit(1)
		}
	}
}
