package main

import (
	"example.com/workspace/lib/tui"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

// Worktree represents a single git worktree.
type Worktree struct {
	Path   string
	Branch string
	IsBare bool
}

type model struct {
	worktrees    []Worktree
	cursor       int
	mode         tui.Mode
	input        textarea.Model
	width        int
	selected     string
	fallbackPath string // set when current worktree is deleted, so quit still cd's somewhere valid
	err          error
	statusMsg    string // transient info message shown in list mode
	current      int    // index of current worktree, -1 if none
}

func initialModel() model {
	wts, err := listWorktrees()
	ti := tui.NewInput("branch-name")

	cur := -1
	if err == nil {
		cur = currentWorktreeIndex(wts)
	}

	cursor := 0
	if cur > 0 {
		cursor = cur
	}

	return model{
		worktrees: wts,
		cursor:    cursor,
		err:       err,
		input:     ti,
		current:   cur,
		width:     80,
	}
}

func (m model) Init() tea.Cmd { return nil }
