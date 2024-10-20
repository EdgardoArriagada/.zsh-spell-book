package main

import (
	"example.com/workspace/lib/args"
	"example.com/workspace/lib/open"
	"fmt"
	"os"
)

func main() {
	parsedArgs, err := args.ParseArguments()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	url := parsedArgs[0]

	if err := open.Url(url); err != nil {
		fmt.Printf("Error opening URL %s: %v\n", url, err)
	}
}
