package main

import (
	"bufio"
	"errors"
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

func parseArguments() ([]string, error) {
	if len(os.Args) > 1 {
		// If command-line arguments are provided, use them
		return os.Args[1:], nil
	} else {
		// If no arguments, check if there's piped input
		stat, _ := os.Stdin.Stat()
		if (stat.Mode() & os.ModeCharDevice) == 0 {
			// Data is being piped to stdin
			var args []string
			scanner := bufio.NewScanner(os.Stdin)
			for scanner.Scan() {
				arg := strings.TrimSpace(scanner.Text())
				if arg != "" {
					args = append(args, arg)
				}
			}
			if err := scanner.Err(); err != nil {
				return nil, errors.New("Error: Reading input")
			}
			if len(args) == 0 {
				return nil, errors.New("Error: Empty arguments")
			}
			return args, nil
		} else {
			return nil, errors.New("Error: Empty arguments")
		}
	}
}

func main() {
	args, err := parseArguments()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	url := args[0]

	if err := openURL(url); err != nil {
		fmt.Printf("Error opening URL %s: %v\n", url, err)
	}
}
