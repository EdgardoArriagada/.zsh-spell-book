package argsLib

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

func isStdinEmpty() bool {
	stat, _ := os.Stdin.Stat()
	return (stat.Mode() & os.ModeCharDevice) != 0
}

func parseStdin() ([]string, error) {
	var args []string

	if isStdinEmpty() {
		return args, nil
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		arg := scanner.Text()
		if arg != "" {
			args = append(args, arg)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.New("Error: Reading input")
	}

	return args, nil
}

func parseWholeStdin() (string, error) {
	var input strings.Builder

	if isStdinEmpty() {
		return input.String(), nil
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", errors.New("Error: Reading input")
		}
		input.WriteString(line)
		if err == io.EOF {
			break
		}
	}

	return input.String(), nil
}

func Parse() ([]string, error) {
	if len(os.Args) > 1 {
		return os.Args[1:], nil
	} else {
		return parseStdin()
	}
}

func ParseWhole() (string, error) {
	if len(os.Args) > 1 {
		return strings.Join(os.Args[1:], "\n"), nil
	} else {
		return parseWholeStdin()
	}
}
