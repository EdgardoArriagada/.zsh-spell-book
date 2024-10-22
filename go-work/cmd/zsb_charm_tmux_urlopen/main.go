package main

import (
	"example.com/workspace/lib/argsLib"
	"example.com/workspace/lib/open"
	"fmt"
	"os"
)

func main() {
	args, err := argsLib.Parse()
	if err != nil || len(args) == 0 {
		fmt.Println(err)
		os.Exit(1)
	}

	url := args[0]

	if err := open.Url(url); err != nil {
		fmt.Printf("Error opening URL %s: %v\n", url, err)
	}
}
