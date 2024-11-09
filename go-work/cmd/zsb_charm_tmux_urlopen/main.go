package main

import (
	"errors"
	"fmt"
	"os"
	"regexp"

	"example.com/workspace/lib/argsLib"
	"example.com/workspace/lib/open"
)

func main() {
	args, err := argsLib.Parse()
	if err != nil || args.Len == 0 {
		fmt.Println(err)
		os.Exit(1)
	}

	url, err := extractFirstURL(args.Get(0))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if err := open.Url(url); err != nil {
		fmt.Printf("Error opening URL %s: %v\n", url, err)
	}
}
func extractFirstURL(text string) (string, error) {
	re := regexp.MustCompile(`https?://[^\s"']+`)
	urls := re.FindStringSubmatch(text)
	if len(urls) == 0 {
		return "", errors.New("Error: No valid URL found in the input.")
	}
	return urls[0], nil
}
