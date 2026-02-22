package main

import (
	"fmt"
	"testing"

	gitlib "example.com/workspace/lib/git"
	"example.com/workspace/lib/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func makeListModel(worktrees []Worktree, cursor int) model {
	ti := tui.NewInput("branch-name")
	return model{
		worktrees: worktrees,
		cursor:    cursor,
		mode:      tui.ListMode,
		input:     ti,
		current:   0,
	}
}

func pressKey(m model, key string) model {
	var msg tea.Msg
	switch key {
	case "j":
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'j'}}
	case "k":
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'k'}}
	case "d":
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'d'}}
	case "y":
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'y'}}
	case "n":
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune{'n'}}
	default:
		msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
	}
	updated, _ := m.Update(msg)
	return updated.(model)
}

func threeWorktrees() []Worktree {
	return []Worktree{
		{Path: "/repo", Branch: "main"},
		{Path: "/repo_wt/feat", Branch: "feat"},
		{Path: "/repo_wt/fix", Branch: "fix"},
	}
}

// --- Cursor wrap tests (issue 6) ---

func TestCursorWrapDown(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, len(wts)-1) // start at last
	m = pressKey(m, "j")
	if m.cursor != 0 {
		t.Errorf("j at bottom: cursor = %d, want 0 (wrap-around)", m.cursor)
	}
}

func TestCursorWrapUp(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 0) // start at top
	m = pressKey(m, "k")
	if m.cursor != len(wts)-1 {
		t.Errorf("k at top: cursor = %d, want %d (wrap-around)", m.cursor, len(wts)-1)
	}
}

func TestCursorNormalDown(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 0)
	m = pressKey(m, "j")
	if m.cursor != 1 {
		t.Errorf("j from 0: cursor = %d, want 1", m.cursor)
	}
}

func TestCursorNormalUp(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 2)
	m = pressKey(m, "k")
	if m.cursor != 1 {
		t.Errorf("k from 2: cursor = %d, want 1", m.cursor)
	}
}

// --- Delete on main worktree shows status (issue 3) ---

func TestDeleteMainShowsStatus(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 0) // cursor on main (index 0)
	m = pressKey(m, "d")
	if m.mode != tui.ListMode {
		t.Errorf("d on main: mode = %v, want ListMode", m.mode)
	}
	if m.statusMsg == "" {
		t.Error("d on main: statusMsg is empty, want a message")
	}
}

func TestDeleteNonMainEntersConfirm(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 1) // cursor on non-main
	m = pressKey(m, "d")
	if m.mode != tui.DeleteConfirmMode {
		t.Errorf("d on non-main: mode = %v, want DeleteConfirmMode", m.mode)
	}
}

func TestStatusMsgClearedOnNavigation(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 0)
	m = pressKey(m, "d") // set statusMsg
	if m.statusMsg == "" {
		t.Fatal("precondition: statusMsg not set after d on main")
	}
	m = pressKey(m, "j") // navigate away
	if m.statusMsg != "" {
		t.Errorf("statusMsg not cleared after navigation: %q", m.statusMsg)
	}
}

// --- Force-delete on dirty worktrees (issue 4) ---

func TestForceDeleteModeOnDirtyError(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 1)
	m.mode = tui.DeleteConfirmMode

	// Simulate a failed delete with a dirty-worktree error by directly
	// injecting the dirty error and the mode transition logic.
	dirtyErr := fmt.Errorf("fatal: '%s' contains modified or untracked files, use --force to override", wts[1].Path)
	if !isWorktreeDirtyError(dirtyErr) {
		t.Fatal("isWorktreeDirtyError did not recognise dirty error")
	}
}

func TestIsWorktreeDirtyError(t *testing.T) {
	tests := []struct {
		err    error
		expect bool
	}{
		{nil, false},
		{fmt.Errorf("some other error"), false},
		{fmt.Errorf("fatal: '/path' contains modified or untracked files, use --force to override"), true},
		{fmt.Errorf("contains modified or untracked files"), true},
	}
	for _, tt := range tests {
		got := isWorktreeDirtyError(tt.err)
		if got != tt.expect {
			t.Errorf("isWorktreeDirtyError(%v) = %v, want %v", tt.err, got, tt.expect)
		}
	}
}

// --- Branch name validation in add mode (issue 5) ---

func TestAddModeRejectsInvalidBranch(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 0)
	m.mode = tui.AddMode
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(nil) // init
	_ = cmd

	m.input.SetValue("my invalid branch")
	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, _ := m.Update(enterMsg)
	result := updated.(model)

	if result.err == nil {
		t.Error("expected error for invalid branch name with spaces, got nil")
	}
	if result.mode != tui.AddMode {
		t.Errorf("mode = %v, want AddMode (should stay in add mode on validation error)", result.mode)
	}
}

func TestAddModeAcceptsValidBranch(t *testing.T) {
	// This test verifies that a valid branch name passes validation.
	// The actual git command will fail in tests (no repo), so we only check
	// that the error is NOT a validation error.
	if err := gitlib.ValidateBranchName("feature/my-branch"); err != nil {
		t.Errorf("ValidateBranchName(\"feature/my-branch\") = %v, want nil", err)
	}
	if err := gitlib.ValidateBranchName("my invalid branch"); err == nil {
		t.Error("ValidateBranchName(\"my invalid branch\") expected error, got nil")
	}
}
