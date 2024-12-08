package main

import (
	"fmt"
	"strconv"
	"strings"

	"example.com/workspace/lib/args"
	u "example.com/workspace/lib/utils"
)

func main() {
	d := u.Must(args.Parse())
	u.Expect(d.Len == 2, "Usage: repeatstr <number> <word>")

	n, err := strconv.Atoi(d.Args[0])
	u.Expect(err == nil, "First argument must be a number")

	word := d.Args[1]
	var builder strings.Builder
	for i := 0; i < n; i++ {
		builder.WriteString(word)
	}
	fmt.Println(builder.String())
}
