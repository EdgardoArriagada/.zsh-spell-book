package main

import (
	"example.com/workspace/lib/args"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	d, err := args.Parse()
	if err != nil || d.Len < 2 {
		fmt.Println(err)
		os.Exit(1)
	}

	n, err := strconv.Atoi(d.Args[0])
	if err != nil {
		fmt.Println("First argument must be a number")
		os.Exit(1)
	}

	word := d.Args[1]
	var builder strings.Builder
	for i := 0; i < n; i++ {
		builder.WriteString(word)
	}
	fmt.Println(builder.String())
}
