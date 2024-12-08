package main

import (
	"example.com/workspace/lib/args"
	"example.com/workspace/lib/open"
	u "example.com/workspace/lib/utils"
)

func main() {
	d := u.Must(args.ParseWithStdin())
	u.Expect(d.Len == 1, "Usage: zsb_open <url>")

	open.Url(d.Args[0]) // first args is url
}
