package argsLib

import (
	"bufio"
	"errors"
	"io"
	"os"
	"strings"
)

type ParsedArgs struct {
	Args []string
	Len  int
}

func (p *ParsedArgs) Get(index int) string {
	if index < p.Len {
		return p.Args[index]
	}
	return ""
}

type ParsedWhole struct {
	Content string
}

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

func Parse() (*ParsedArgs, error) {
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	} else {
		var err error
		args, err = parseStdin()
		if err != nil {
			return nil, err
		}
	}
	return &ParsedArgs{Args: args, Len: len(args)}, nil
}

func ParseWhole() (*ParsedWhole, error) {
	var content string
	if len(os.Args) > 1 {
		content = strings.Join(os.Args[1:], "\n")
	} else {
		var err error
		content, err = parseWholeStdin()
		if err != nil {
			return nil, err
		}
	}
	return &ParsedWhole{Content: content}, nil
}
