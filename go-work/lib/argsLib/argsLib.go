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
	var result []string

	if isStdinEmpty() {
		return result, nil
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		arg := scanner.Text()
		if arg != "" {
			result = append(result, arg)
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, errors.New("Error: Reading input")
	}

	return result, nil
}

func parseWholeStdin() (string, error) {
	var result strings.Builder

	if isStdinEmpty() {
		return "", nil
	}

	reader := bufio.NewReader(os.Stdin)
	for {
		line, err := reader.ReadString('\n')
		if err != nil && err != io.EOF {
			return "", errors.New("Error: Reading input")
		}
		result.WriteString(line)
		if err == io.EOF {
			break
		}
	}

	return result.String(), nil
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

func ReadArg(args *[]string, index int) string {
	if len(*args) > index {
		return (*args)[index]
	}

	return ""
}
