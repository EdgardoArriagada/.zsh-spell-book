package main

import (
	"fmt"
	"strings"
	"testing"

	"example.com/workspace/lib/tui"
)

func TestViewOnlyRendersVisibleBranches(t *testing.T) {
	brs := manyBranches(10)
	ti := tui.NewInput("branch-name")
	m := model{
		branches:         brs,
		filtered:         brs,
		cursor:           3,
		current:          0,
		windowHeight:     40,
		vp:               tui.Viewport{Offset: 2},
		input:            ti,
		firstWorktreeIdx: findFirstWorktreeIdx(brs),
	}

	view := m.View()
	maxVis := m.vp.MaxVisible(len(brs), m.availableRows())
	start := m.vp.Offset
	end := start + maxVis

	for i := start; i < end && i < len(brs); i++ {
		name := fmt.Sprintf("branch-%02d", i)
		if !strings.Contains(view, name) {
			t.Errorf("branch %q (index %d) should be visible but not found in view", name, i)
		}
	}

	for i := 0; i < len(brs); i++ {
		if i >= start && i < end {
			continue
		}
		name := fmt.Sprintf("branch-%02d", i)
		if strings.Contains(view, name) {
			t.Errorf("branch %q (index %d) should NOT be visible but found in view", name, i)
		}
	}
}

// --- Worktree section divider and footer hint ---

func TestViewShowsDividerBeforeWorktrees(t *testing.T) {
	brs := []Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feat"},
		{Name: "feat-wt", IsWorktree: true},
	}
	ti := tui.NewInput("branch-name")
	m := model{branches: brs, filtered: brs, cursor: 0, current: 0, input: ti, firstWorktreeIdx: findFirstWorktreeIdx(brs)}

	view := m.View()

	if !strings.Contains(view, "Worktree Branches") {
		t.Error("view should contain a 'Worktree Branches' section header when worktree branches exist")
	}
}

func TestViewNoDividerWithoutWorktrees(t *testing.T) {
	brs := []Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feat"},
	}
	ti := tui.NewInput("branch-name")
	m := model{branches: brs, filtered: brs, cursor: 0, current: 0, input: ti, firstWorktreeIdx: findFirstWorktreeIdx(brs)}

	view := m.View()

	if strings.Contains(view, "Worktree Branches") {
		t.Error("view should not contain 'Worktree Branches' section header when no worktree branches exist")
	}
}

func TestViewFooterShowsWorktreeHintWhenPresent(t *testing.T) {
	brs := []Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feat-wt", IsWorktree: true},
	}
	ti := tui.NewInput("branch-name")
	m := model{branches: brs, filtered: brs, cursor: 0, current: 0, input: ti, firstWorktreeIdx: findFirstWorktreeIdx(brs)}

	view := m.View()

	if !strings.Contains(view, "checked out elsewhere") {
		t.Error("footer should show 'checked out elsewhere' hint when worktree branches exist")
	}
}

func TestViewFooterNoWorktreeHintWhenAbsent(t *testing.T) {
	brs := []Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feat"},
	}
	ti := tui.NewInput("branch-name")
	m := model{branches: brs, filtered: brs, cursor: 0, current: 0, input: ti, firstWorktreeIdx: findFirstWorktreeIdx(brs)}

	view := m.View()

	if strings.Contains(view, "checked out elsewhere") {
		t.Error("footer should not show worktree hint when no worktree branches exist")
	}
}

func TestViewWorktreeWarningBanner(t *testing.T) {
	brs := []Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feat"},
	}
	ti := tui.NewInput("branch-name")
	m := model{branches: brs, filtered: brs, cursor: 0, current: 0, input: ti, inWorktree: true, firstWorktreeIdx: findFirstWorktreeIdx(brs)}

	view := m.View()

	if !strings.Contains(view, "worktree") {
		t.Error("view should contain worktree warning when inWorktree=true")
	}
}

func TestViewNoWorktreeWarningWhenNotInWorktree(t *testing.T) {
	brs := []Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feat"},
	}
	ti := tui.NewInput("branch-name")
	m := model{branches: brs, filtered: brs, cursor: 0, current: 0, input: ti, inWorktree: false, firstWorktreeIdx: findFirstWorktreeIdx(brs)}

	view := m.View()

	if strings.Contains(view, "⚠") {
		t.Error("view should not contain worktree warning when inWorktree=false")
	}
}

func TestViewNoScrollWhenHeightZero(t *testing.T) {
	// windowHeight=0 means "no terminal size received yet" → all branches must render
	brs := manyBranches(5)
	ti := tui.NewInput("branch-name")
	m := model{
		branches:         brs,
		filtered:         brs,
		cursor:           0,
		current:          0,
		input:            ti,
		firstWorktreeIdx: findFirstWorktreeIdx(brs),
	}

	view := m.View()

	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("branch-%02d", i)
		if !strings.Contains(view, name) {
			t.Errorf("branch %q should be visible when windowHeight=0 (no viewport) but not found in view", name)
		}
	}
}
