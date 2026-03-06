package main

import (
	"path/filepath"

	"example.com/workspace/lib/tui"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	worktrees    []Worktree
	filtered     []Worktree
	cursor       int
	mode         tui.Mode
	input        textarea.Model
	searchInput  textinput.Model
	width        int
	vp           tui.Viewport
	selected     string
	fallbackPath string // set when current worktree is deleted, so quit still cd's somewhere valid
	err          error
	statusMsg    string // transient info message shown in list mode
	current      int    // index of current worktree in worktrees, -1 if none
}

func applyWorktreeFilter(wts []Worktree, term string) []Worktree {
	return tui.ApplyFilter(wts, term, func(wt Worktree) string { return filepath.Base(wt.Path) })
}

func initialModel() model {
	wts, err := listWorktrees()
	ti := tui.NewInput("branch-name")
	si := tui.NewSearchInput()

	cur := -1
	if err == nil {
		cur = currentWorktreeIndex(wts)
	}

	cursor := 0
	if cur > 0 {
		cursor = cur
	}

	return model{
		worktrees:   wts,
		filtered:    wts,
		cursor:      cursor,
		err:         err,
		input:       ti,
		searchInput: si,
		current:     cur,
		width:       tui.DefaultWidth,
	}
}

func (m model) Init() tea.Cmd { return nil }
