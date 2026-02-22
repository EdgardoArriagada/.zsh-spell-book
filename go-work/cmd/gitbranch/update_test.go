package main

import (
	"fmt"
	"testing"

	"example.com/workspace/lib/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func makeListModel(branches []Branch, cursor, current int) model {
	ti := tui.NewInput("branch-name")
	return model{
		branches: branches,
		cursor:   cursor,
		mode:     tui.ListMode,
		input:    ti,
		current:  current,
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

func threeBranches() []Branch {
	return []Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feat", IsCurrent: false},
		{Name: "fix", IsCurrent: false},
	}
}

// --- Cursor wrap tests ---

func TestCursorWrapDown(t *testing.T) {
	brs := threeBranches()
	m := makeListModel(brs, len(brs)-1, 0) // start at last
	m = pressKey(m, "j")
	if m.cursor != 0 {
		t.Errorf("j at bottom: cursor = %d, want 0 (wrap-around)", m.cursor)
	}
}

func TestCursorWrapUp(t *testing.T) {
	brs := threeBranches()
	m := makeListModel(brs, 0, 0) // start at top
	m = pressKey(m, "k")
	if m.cursor != len(brs)-1 {
		t.Errorf("k at top: cursor = %d, want %d (wrap-around)", m.cursor, len(brs)-1)
	}
}

func TestCursorNormalDown(t *testing.T) {
	brs := threeBranches()
	m := makeListModel(brs, 0, 0)
	m = pressKey(m, "j")
	if m.cursor != 1 {
		t.Errorf("j from 0: cursor = %d, want 1", m.cursor)
	}
}

func TestCursorNormalUp(t *testing.T) {
	brs := threeBranches()
	m := makeListModel(brs, 2, 0)
	m = pressKey(m, "k")
	if m.cursor != 1 {
		t.Errorf("k from 2: cursor = %d, want 1", m.cursor)
	}
}

// --- Delete on current branch shows status ---

func TestDeleteCurrentShowsStatus(t *testing.T) {
	brs := threeBranches()
	m := makeListModel(brs, 0, 0) // cursor on current branch
	m = pressKey(m, "d")
	if m.mode != tui.ListMode {
		t.Errorf("d on current: mode = %v, want ListMode", m.mode)
	}
	if m.statusMsg == "" {
		t.Error("d on current: statusMsg is empty, want a message")
	}
}

func TestDeleteNonCurrentEntersConfirm(t *testing.T) {
	brs := threeBranches()
	m := makeListModel(brs, 1, 0) // cursor on non-current
	m = pressKey(m, "d")
	if m.mode != tui.DeleteConfirmMode {
		t.Errorf("d on non-current: mode = %v, want DeleteConfirmMode", m.mode)
	}
}

func TestStatusMsgClearedOnNavigation(t *testing.T) {
	brs := threeBranches()
	m := makeListModel(brs, 0, 0)
	m = pressKey(m, "d") // set statusMsg
	if m.statusMsg == "" {
		t.Fatal("precondition: statusMsg not set after d on current")
	}
	m = pressKey(m, "j") // navigate away
	if m.statusMsg != "" {
		t.Errorf("statusMsg not cleared after navigation: %q", m.statusMsg)
	}
}

// --- Force-delete on unmerged branches ---

func TestForceDeleteModeOnUnmergedError(t *testing.T) {
	brs := threeBranches()
	m := makeListModel(brs, 1, 0)
	m.mode = tui.DeleteConfirmMode

	unmergedErr := fmt.Errorf("error: The branch 'feat' is not fully merged")
	if !isUnmergedBranchError(unmergedErr) {
		t.Fatal("isUnmergedBranchError did not recognise unmerged error")
	}
}

// --- Branch name validation in add mode ---

func TestAddModeRejectsInvalidBranch(t *testing.T) {
	brs := threeBranches()
	m := makeListModel(brs, 0, 0)
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
