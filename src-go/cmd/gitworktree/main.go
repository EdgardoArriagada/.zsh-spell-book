package main

import (
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func main() {
	p := tea.NewProgram(initialModel(), tea.WithOutput(os.Stderr))
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	mdl := m.(model)
	if mdl.selected != "" {
		fmt.Print(mdl.selected)
		os.Exit(0)
	}
	if mdl.fallbackPath != "" {
		fmt.Print(mdl.fallbackPath)
		os.Exit(0)
	}
	os.Exit(1)
}
