package main

import (
	"example.com/workspace/lib/argsLib"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	args, err := argsLib.Parse()
	if err != nil || args.Len < 2 {
		fmt.Println(err)
		os.Exit(1)

	}

	n, err := strconv.Atoi(args.Get(0))
	if err != nil {
		fmt.Println("First argument must be a number")
		os.Exit(1)
	}

	word := args.Get(1)
	var builder strings.Builder
	for i := 0; i < n; i++ {
		builder.WriteString(word)
	}
	fmt.Println(builder.String())
}
