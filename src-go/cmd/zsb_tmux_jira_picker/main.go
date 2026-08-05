package main

import (
	"fmt"
	"os"

	"zsb_tmux_jira_picker/internal/app"
)

func main() {
	selected, err := app.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if !selected {
		os.Exit(1)
	}
}
