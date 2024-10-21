package args

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func parseArgsFromStdin() ([]string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 {
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

func ParseArguments() ([]string, error) {
	if len(os.Args) > 1 {
		return os.Args[1:], nil
	} else {
		return parseArgsFromStdin()
	}
}
