package main

import (
	"example.com/workspace/lib/argsLib"
	"example.com/workspace/lib/open"
	"fmt"
	"os"
	"regexp"
)

func main() {
	args, err := argsLib.Parse()
	if err != nil || len(args) == 0 {
		fmt.Println(err)
		os.Exit(1)
	}

	url := extractFirstURL(args[0])
	if url == "" {
		fmt.Println("No valid URL found in the input.")
		os.Exit(1)
	}

	if err := open.Url(url); err != nil {
		fmt.Printf("Error opening URL %s: %v\n", url, err)
	}
}
func extractFirstURL(text string) string {
	// Regular expression to find URLs
	re := regexp.MustCompile(`https?://[^\s]+`)
	urls := re.FindStringSubmatch(text)
	if len(urls) > 0 {
		return urls[0]
	}
	return ""
}
