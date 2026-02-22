package main

import (
	"example.com/workspace/lib/tui"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// Branch represents a single local git branch.
type Branch struct {
	Name       string
	IsCurrent  bool
	IsWorktree bool
}

type model struct {
	branches   []Branch
	cursor     int
	mode       tui.Mode
	input      textarea.Model
	width      int
	vp         tui.Viewport
	selected   string
	err        error
	statusMsg  string // transient info message shown in list mode
	current    int    // index of current branch, -1 if none
	inWorktree bool
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

	inWT := checkIsLinkedWorktree()
	vp := tui.Viewport{}
	if inWT {
		vp.ExtraOverhead = 2
	}

	return model{
		branches:   branches,
		cursor:     cursor,
		err:        err,
		input:      ti,
		current:    cur,
		width:      tui.DefaultWidth,
		inWorktree: inWT,
		vp:         vp,
	}
}

func (m model) Init() tea.Cmd { return nil }
