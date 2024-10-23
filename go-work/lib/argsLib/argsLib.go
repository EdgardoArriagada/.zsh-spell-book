package argsLib

import (
	"bufio"
	"errors"
	"io"
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
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) != 0 {
		return "", nil
	}

	reader := bufio.NewReader(os.Stdin)
	var input strings.Builder
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
		return parseArgsFromStdin()
	}
}

func ParseWhole() (string, error) {
	if len(os.Args) > 1 {
		return strings.Join(os.Args[1:], "\n"), nil
	} else {
		return parseWholeStdin()
	}
}
