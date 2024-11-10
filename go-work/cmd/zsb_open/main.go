package main

import (
	"example.com/workspace/lib/args"
	"example.com/workspace/lib/open"
	"fmt"
	"os"
)

func main() {
	d, err := args.ParseWithStdin()
	if err != nil || d.Len == 0 {
		fmt.Println(err)
		os.Exit(1)
	}

	open.Url(d.Args[0]) // first args is url
}
