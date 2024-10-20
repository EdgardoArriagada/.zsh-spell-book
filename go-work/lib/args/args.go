package args

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func ParseArguments() ([]string, error) {
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
