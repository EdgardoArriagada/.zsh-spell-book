package main

import (
	"fmt"
	"strings"
	"testing"

	"example.com/workspace/lib/tui"
)

func TestViewOnlyRendersVisibleBranches(t *testing.T) {
	// height=10, viewportOverhead=6 → maxVisible=4
	// offset=2 → visible indices 2..5
	brs := manyBranches(10)
	ti := tui.NewInput("branch-name")
	m := model{
		branches: brs,
		cursor:   3,
		current:  0,
		height:   10,
		offset:   2,
		input:    ti,
	}

	view := m.View()

	// branches at indices 2..5 must appear
	for i := 2; i <= 5; i++ {
		name := fmt.Sprintf("branch-%02d", i)
		if !strings.Contains(view, name) {
			t.Errorf("branch %q (index %d) should be visible (offset=2, maxVis=4) but not found in view", name, i)
		}
	}

	// branches outside the window must NOT appear
	for _, i := range []int{0, 1, 6, 7, 8, 9} {
		name := fmt.Sprintf("branch-%02d", i)
		if strings.Contains(view, name) {
			t.Errorf("branch %q (index %d) should NOT be visible (offset=2, maxVis=4) but found in view", name, i)
		}
	}
}

func TestViewNoScrollWhenHeightZero(t *testing.T) {
	// height=0 means "no terminal size received yet" → all branches must render
	brs := manyBranches(5)
	ti := tui.NewInput("branch-name")
	m := model{
		branches: brs,
		cursor:   0,
		current:  0,
		height:   0,
		offset:   0,
		input:    ti,
	}

	view := m.View()

	for i := 0; i < 5; i++ {
		name := fmt.Sprintf("branch-%02d", i)
		if !strings.Contains(view, name) {
			t.Errorf("branch %q should be visible when height=0 (no viewport) but not found in view", name)
		}
	}
}
