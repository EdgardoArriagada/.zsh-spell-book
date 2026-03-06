package main

import (
	"example.com/workspace/lib/tui"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// Branch represents a single local git branch.
type Branch struct {
	Name       string
	IsCurrent  bool
	IsWorktree bool
}

type model struct {
	branches         []Branch
	filtered         []Branch
	cursor           int
	mode             tui.Mode
	input            textarea.Model
	searchInput      textinput.Model
	width            int
	vp               tui.Viewport
	selected         string
	err              error
	statusMsg        string // transient info message shown in list mode
	current          int    // index of current branch in branches, -1 if none
	inWorktree       bool
	firstWorktreeIdx int // cached index of first worktree branch in branches, -1 if none
}

func findFirstWorktreeIdx(branches []Branch) int {
	for i, br := range branches {
		if br.IsWorktree {
			return i
		}
	}
	return -1
}

func branchNames(branches []Branch) []string {
	names := make([]string, len(branches))
	for i, br := range branches {
		names[i] = br.Name
	}
	return names
}

func applyBranchFilter(branches []Branch, term string) []Branch {
	indices := tui.FuzzyFilter(term, branchNames(branches))
	result := make([]Branch, len(indices))
	for i, idx := range indices {
		result[i] = branches[idx]
	}
	return result
}

func initialModel() model {
	branches, inWT, err := initBranchData()
	ti := tui.NewInput("branch-name")
	si := tui.NewSearchInput()

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

	vp := tui.Viewport{}
	if inWT {
		vp.ExtraOverhead = 2
	}

	return model{
		branches:         branches,
		filtered:         branches,
		cursor:           cursor,
		err:              err,
		input:            ti,
		searchInput:      si,
		current:          cur,
		width:            tui.DefaultWidth,
		inWorktree:       inWT,
		vp:               vp,
		firstWorktreeIdx: findFirstWorktreeIdx(branches),
	}
}

func (m model) Init() tea.Cmd { return nil }
