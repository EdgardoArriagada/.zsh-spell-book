package main

import (
	"fmt"
	"os"

	"gitbranch/internal/app"
)

func main() {
	selected, err := app.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if selected != "" {
		fmt.Print(selected)
		os.Exit(0)
	}
	os.Exit(1)
}
