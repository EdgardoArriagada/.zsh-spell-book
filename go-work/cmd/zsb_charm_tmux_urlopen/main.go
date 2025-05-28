package main

import (
	"errors"
	"regexp"

	"example.com/workspace/lib/args"
	"example.com/workspace/lib/open"
	u "example.com/workspace/lib/utils"
)

func main() {
	d := u.Must1(args.ParseWithStdin())
	u.Expect(d.Len == 1, "Usage: zsb_charm_tmux_urlopen <url>")

	url := u.Must1(ExtractFirstURL(d.Args[0]))
	open.Url(url)
}

func ExtractFirstURL(text string) (string, error) {
	re := regexp.MustCompile(`https?://[\w\-_\.%?/:+=&#%]+`)
	urls := re.FindStringSubmatch(text)
	if len(urls) == 0 {
		return "", errors.New("Error: No valid URL found in the input.")
	}
	return urls[0], nil
}
