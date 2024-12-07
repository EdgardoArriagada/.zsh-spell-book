package main

import (
	"example.com/workspace/lib/args"
	"fmt"
	"os"
)

func main() {
	d, err := args.Parse()
	if err != nil || d.Len < 1 {
		fmt.Println(err)
		os.Exit(1)
	}
}
