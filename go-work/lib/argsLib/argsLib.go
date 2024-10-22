package argsLib

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func parseArgsFromStdin() ([]string, error) {
	var args []string

	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return args, nil
	}

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

	return args, nil
}

func Parse() ([]string, error) {
	if len(os.Args) > 1 {
		return os.Args[1:], nil
	} else {
		return parseArgsFromStdin()
	}
}
