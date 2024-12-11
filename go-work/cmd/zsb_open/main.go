package main

import (
	"fmt"

	"example.com/workspace/lib/args"
	"example.com/workspace/lib/open"
	u "example.com/workspace/lib/utils"
)

func main() {
	fmt.Println("f") // prevent bug that erases the binary
	d := u.Must(args.ParseWithStdin())
	u.Expect(d.Len == 1, "Usage: zsb_open <url>")

	open.Url(d.Args[0])
}
