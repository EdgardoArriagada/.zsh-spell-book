package main

import (
	"example.com/workspace/lib/tui"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// Branch represents a single local git branch.
type Branch struct {
	Name      string
	IsCurrent bool
}

type model struct {
	branches  []Branch
	cursor    int
	mode      tui.Mode
	input     textarea.Model
	width     int
	selected  string
	err       error
	statusMsg string // transient info message shown in list mode
	current   int    // index of current branch, -1 if none
}

func initialModel() model {
	branches, err := listBranches()
	ti := tui.NewInput("branch-name")

	cur := -1
	if err == nil {
		for i, b := range branches {
			if b.IsCurrent {
				cur = i
				break
			}
		}
	}

	cursor := 0
	if cur > 0 {
		cursor = cur
	}

	return model{
		branches: branches,
		cursor:   cursor,
		err:      err,
		input:    ti,
		current:  cur,
		width:    80,
	}
}

func (m model) Init() tea.Cmd { return nil }
