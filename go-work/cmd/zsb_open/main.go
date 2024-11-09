package main

import (
	"example.com/workspace/lib/argsLib"
	"example.com/workspace/lib/open"
	"fmt"
	"os"
)

func main() {
	args, err := argsLib.Parse(true)
	if err != nil || args.Len == 0 {
		fmt.Println(err)
		os.Exit(1)
	}

	open.Url(args.Args[0]) // first args is url
}
