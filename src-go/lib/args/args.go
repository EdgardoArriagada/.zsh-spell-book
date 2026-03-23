package args

import (
	"os"
	"strings"
)

func parse(stdinFallback bool) (*ParsedArgs, error) {
	var args []string
	if len(os.Args) > 1 {
		args = os.Args[1:]
	} else if stdinFallback {
		var err error
		args, err = parseStdin()
		if err != nil {
			return nil, err
		}
	}
	return &ParsedArgs{Args: args, Len: len(args)}, nil
}

func parseWhole(stdinFallback bool) (*ParsedWhole, error) {
	var content string
	if len(os.Args) > 1 {
		content = strings.Join(os.Args[1:], "\n")
	} else if stdinFallback {
		var err error
		content, err = parseWholeStdin()
		if err != nil {
			return nil, err
		}
	}
	return &ParsedWhole{Content: content}, nil
}

func Parse() (*ParsedArgs, error) {
	return parse(false)
}

func ParseWithStdin() (*ParsedArgs, error) {
	return parse(true)
}

func ParseWhole() (*ParsedWhole, error) {
	return parseWhole(false)
}

func ParseWholeWithStdin() (*ParsedWhole, error) {
	return parseWhole(true)
}
