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

func worktreeNames(wts []Worktree) []string {
	names := make([]string, len(wts))
	for i, wt := range wts {
		names[i] = filepath.Base(wt.Path)
	}
	return names
}

func applyWorktreeFilter(wts []Worktree, term string) []Worktree {
	indices := tui.FuzzyFilter(term, worktreeNames(wts))
	result := make([]Worktree, len(indices))
	for i, idx := range indices {
		result[i] = wts[idx]
	}
	return result
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
