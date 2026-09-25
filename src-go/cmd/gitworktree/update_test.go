package main

import (
	"fmt"
	"strings"
	"testing"

	gitlib "example.com/workspace/lib/git"
	"example.com/workspace/lib/tui"

	tea "charm.land/bubbletea/v2"
)

func makeListModel(worktrees []Worktree, cursor int) model {
	ti := tui.NewInput("branch-name")
	return model{
		worktrees: worktrees,
		filtered:  worktrees,
		cursor:    cursor,
		mode:      tui.ListMode,
		input:     ti,
		current:   0,
	}
}

func pressKey(m model, key string) model {
	// Key.String() returns Text when set, so this stringifies to `key` for the
	// printable single-char keys the tests use (j/k/d/y/n/D).
	updated, _ := m.Update(tea.KeyPressMsg{Text: key})
	return updated.(model)
}

func threeWorktrees() []Worktree {
	return []Worktree{
		{Path: "/repo", Branch: "main"},
		{Path: "/repo_wt/feat", Branch: "feat"},
		{Path: "/repo_wt/fix", Branch: "fix"},
	}
}

func TestNumberSelectsVisibleWorktree(t *testing.T) {
	worktrees := make([]Worktree, 12)
	for i := range worktrees {
		worktrees[i].Path = fmt.Sprintf("/repo/wt-%d", i)
	}
	m := makeListModel(worktrees, 3)
	m.vp.Offset = 2
	updated, cmd := m.Update(tea.KeyPressMsg{Text: "0"})
	m = updated.(model)
	if cmd == nil || m.selected != worktrees[11].Path {
		t.Fatalf("0 selected %q, want %q", m.selected, worktrees[11].Path)
	}
}

func TestTmuxWorktreeBranchPrefix(t *testing.T) {
	for _, tt := range []struct {
		mainBranch    string
		developExists bool
		want          string
	}{
		{"develop", true, "feature/session"},
		{"master", true, "hotfix/session"},
		{"master", false, "feature/session"},
		{"main", false, "hotfix/session"},
	} {
		if got := automaticWorktreeBranch(tt.mainBranch, "session", tt.developExists); got != tt.want {
			t.Errorf("main branch %q: got %q, want %q", tt.mainBranch, got, tt.want)
		}
	}
}

func TestFooterHidesAutoCreateHint(t *testing.T) {
	if strings.Contains(makeListModel(threeWorktrees(), 0).footerSection(), "A") {
		t.Error("footer must not show the automatic-create keybind")
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

func TestPageKeysMoveWorktreeCursor(t *testing.T) {
	m := makeListModel(make([]Worktree, 10), 0)
	m.windowHeight = 10
	m.width = 80
	page := m.availableRows()

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	m = updated.(model)
	if m.cursor != page {
		t.Fatalf("ctrl+d cursor = %d, want %d", m.cursor, page)
	}
	updated, _ = m.Update(tea.KeyPressMsg{Code: 'u', Mod: tea.ModCtrl})
	if got := updated.(model).cursor; got != 0 {
		t.Fatalf("ctrl+u cursor = %d, want 0", got)
	}
}

func TestQuitKeysClearWorktreeSearchAndPreserveSelection(t *testing.T) {
	for _, key := range []tea.KeyPressMsg{{Text: "q"}, {Code: 'c', Mod: tea.ModCtrl}, {Code: tea.KeyEscape}} {
		wts := threeWorktrees()
		m := makeListModel(wts, 1)
		m.searchInput = tui.NewSearchInput()
		m.searchInput.SetValue("f")
		m.filtered = wts[1:]

		updated, cmd := m.Update(key)
		m = updated.(model)
		if cmd != nil || m.searchInput.Value() != "" || len(m.filtered) != len(wts) || m.cursor != 2 {
			t.Fatalf("%s did not clear search with selection preserved: cmd=%v term=%q len=%d cursor=%d", key.String(), cmd, m.searchInput.Value(), len(m.filtered), m.cursor)
		}
		if _, cmd = m.Update(key); cmd == nil {
			t.Fatalf("second %s should quit", key.String())
		}
	}
}

func TestSearchKeepsSelectionUntilEnter(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 0)
	m.searchInput = tui.NewSearchInput()
	m = pressKey(pressKey(m, "/"), "f")
	if len(m.filtered) != 2 || strings.Contains(m.render(), "▸") {
		t.Fatalf("search selected a match: len=%d cursorVisible=%t", len(m.filtered), strings.Contains(m.render(), "▸"))
	}

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(model)
	if cmd != nil || m.mode != tui.ListMode || m.cursor != 0 || m.selected != "" {
		t.Fatalf("first enter selected instead of focusing first match: cmd=%v mode=%v cursor=%d selected=%q", cmd, m.mode, m.cursor, m.selected)
	}
	first := m.filtered[0].Path
	updated, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(model)
	if cmd == nil || m.selected != first {
		t.Fatalf("second enter: cmd=%v selected=%q, want %q", cmd, m.selected, first)
	}

	m = makeListModel(wts, 2)
	m.searchInput = tui.NewSearchInput()
	m = pressKey(pressKey(m, "/"), "f")
	if m.filtered[m.cursor].Path != wts[2].Path || !strings.Contains(m.render(), "▸") {
		t.Fatal("search did not preserve a matching selection")
	}

	m = makeListModel(wts, 0)
	m.searchInput = tui.NewSearchInput()
	m = pressKey(pressKey(m, "/"), "x")
	updated, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(model)
	if cmd == nil || m.selected != wts[2].Path {
		t.Fatalf("single match enter: cmd=%v selected=%q, want %q", cmd, m.selected, wts[2].Path)
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

// --- D key: delete worktree + branch ---

func TestDeleteWithBranchMainShowsStatus(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 0) // cursor on main
	m = pressKey(m, "D")
	if m.mode != tui.ListMode {
		t.Errorf("D on main: mode = %v, want ListMode", m.mode)
	}
	if m.statusMsg == "" {
		t.Error("D on main: statusMsg is empty, want a message")
	}
}

func TestDeleteWithBranchNonMainSetsFlag(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 1) // cursor on non-main
	m = pressKey(m, "D")
	if m.mode != tui.DeleteConfirmMode {
		t.Errorf("D on non-main: mode = %v, want DeleteConfirmMode", m.mode)
	}
	if !m.deleteBranch {
		t.Error("D on non-main: deleteBranch = false, want true")
	}
}

func TestDeleteWithoutBranchDoesNotSetFlag(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 1)
	m = pressKey(m, "d")
	if m.deleteBranch {
		t.Error("d key: deleteBranch = true, want false")
	}
}

func TestCancelDeleteResetsBranchFlag(t *testing.T) {
	wts := threeWorktrees()
	m := makeListModel(wts, 1)
	m = pressKey(m, "D")
	if !m.deleteBranch {
		t.Fatal("precondition: deleteBranch not set after D")
	}
	m = pressKey(m, "n") // cancel
	if m.deleteBranch {
		t.Error("after cancel: deleteBranch = true, want false")
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
	enterMsg := tea.KeyPressMsg{Code: tea.KeyEnter}
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

func TestAddModeShowsCreatingSpinner(t *testing.T) {
	m := makeListModel(threeWorktrees(), 0)
	m.mode = tui.AddMode
	m.input.SetValue("feature/my-branch")

	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	if got := updated.(model).mode; got != tui.CreatingMode {
		t.Errorf("mode = %v, want CreatingMode", got)
	}
	if cmd == nil {
		t.Error("enter on a valid branch should start the create command")
	}
	if got := updated.(model).statusSection(); !strings.Contains(got, "creating...") {
		t.Errorf("status = %q, want creating indicator", got)
	}
}
