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

const viewportOverhead = 6 // title(2) + status/confirm(2) + footer(2)

type model struct {
	branches  []Branch
	cursor    int
	mode      tui.Mode
	input     textarea.Model
	width     int
	height    int
	offset    int // first visible branch index
	selected  string
	err       error
	statusMsg string // transient info message shown in list mode
	current   int    // index of current branch, -1 if none
}

// maxVisible returns the number of branches that fit in the terminal.
// Returns len(branches) when height is 0 (terminal size not yet received).
func (m model) maxVisible() int {
	if m.height == 0 {
		return len(m.branches)
	}
	v := m.height - viewportOverhead
	if v < 1 {
		return 1
	}
	return v
}

// clampViewport adjusts m.offset so that m.cursor is always within the visible window.
func (m model) clampViewport() model {
	if m.height == 0 {
		return m
	}
	maxVis := m.maxVisible()
	if m.cursor < m.offset {
		m.offset = m.cursor
	} else if m.cursor >= m.offset+maxVis {
		m.offset = m.cursor - maxVis + 1
	}
	// Clamp offset so we don't scroll past the end of the list
	if maxOffset := len(m.branches) - maxVis; maxOffset >= 0 && m.offset > maxOffset {
		m.offset = maxOffset
	}
	if m.offset < 0 {
		m.offset = 0
	}
	return m
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
