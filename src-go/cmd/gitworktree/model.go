package main

import (
	"context"
	"path/filepath"

	"example.com/workspace/lib/tui"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/textarea"
	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	worktrees    []Worktree
	filtered     []Worktree
	cursor       int
	mode         tui.Mode
	input        textarea.Model
	searchInput  textinput.Model
	searchTarget string
	spinner      spinner.Model
	width        int
	windowHeight int
	vp           tui.Viewport
	selected     string
	fallbackPath string // set when current worktree is deleted, so quit still cd's somewhere valid
	err          error
	statusMsg    string // transient info message shown in list mode
	current      int    // index of current worktree in worktrees, -1 if none
	deleteBranch bool   // when true, D key was used — also delete the branch after worktree removal
	prStatus     map[string]string
	prPending    int
	prCancel     context.CancelFunc
}

func applyWorktreeFilter(wts []Worktree, term string) []Worktree {
	names := make([]string, len(wts))
	branches := make([]string, len(wts))
	for i, wt := range wts {
		names[i] = filepath.Base(wt.Path)
		branches[i] = wt.Branch
	}
	filtered := make([]Worktree, 0, len(wts))
	seen := make([]bool, len(wts))
	for _, i := range append(tui.FuzzyFilter(term, names), tui.FuzzyFilter(term, branches)...) {
		if !seen[i] {
			filtered = append(filtered, wts[i])
			seen[i] = true
		}
	}
	return filtered
}

func initialModel() model {
	wts, err := listWorktrees()
	ti := tui.NewInput("branch-name")
	si := tui.NewSearchInput()
	sp := spinner.New(spinner.WithSpinner(spinner.Dot), spinner.WithStyle(tui.WarnStyle))

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
		spinner:     sp,
		current:     cur,
		width:       tui.DefaultWidth,
	}
}

func (m model) Init() tea.Cmd { return nil }
