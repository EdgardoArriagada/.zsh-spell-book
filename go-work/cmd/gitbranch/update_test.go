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

func makeListModelWithHeight(branches []Branch, cursor, current, height int) model {
	ti := tui.NewInput("branch-name")
	return model{
		branches: branches,
		cursor:   cursor,
		mode:     tui.ListMode,
		input:    ti,
		current:  current,
		vp:       tui.Viewport{Height: height},
	}
}

func manyBranches(n int) []Branch {
	branches := make([]Branch, n)
	for i := range branches {
		branches[i] = Branch{Name: fmt.Sprintf("branch-%02d", i)}
	}
	branches[0].IsCurrent = true
	return branches
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

// --- Viewport / scrolling tests ---
// viewportOverhead=6, so height=10 → maxVisible=4, height=20 → maxVisible=14

func TestViewportOffsetIncreasesWhenCursorMovesBelow(t *testing.T) {
	// height=10, maxVisible=4; 10 branches starting at cursor=0, offset=0
	brs := manyBranches(10)
	m := makeListModelWithHeight(brs, 0, 0, 10)

	// Move down 4 times: cursor ends at 4, past the initial [0,4) window
	for i := 0; i < 4; i++ {
		m = pressKey(m, "j")
	}

	if m.cursor != 4 {
		t.Fatalf("cursor = %d, want 4", m.cursor)
	}
	maxVis := m.vp.MaxVisible(len(m.branches))
	if m.cursor < m.vp.Offset || m.cursor >= m.vp.Offset+maxVis {
		t.Errorf("cursor %d not visible in viewport [%d, %d)", m.cursor, m.vp.Offset, m.vp.Offset+maxVis)
	}
}

func TestViewportOffsetDecreasesWhenCursorMovesAbove(t *testing.T) {
	// cursor at top of a scrolled-down viewport; k should scroll up
	brs := manyBranches(10)
	m := makeListModelWithHeight(brs, 6, 0, 10) // cursor=6, maxVisible=4
	m.vp.Offset = 6                              // cursor at the very top of viewport

	m = pressKey(m, "k") // cursor → 5, viewport must scroll up

	if m.cursor != 5 {
		t.Fatalf("cursor = %d, want 5", m.cursor)
	}
	if m.vp.Offset != 5 {
		t.Errorf("offset = %d, want 5 (cursor now at top of viewport)", m.vp.Offset)
	}
}

func TestViewportNoScrollWhenAllBranchesFit(t *testing.T) {
	// 3 branches, height=20 → maxVisible=14; all fit, offset must stay 0
	brs := threeBranches()
	m := makeListModelWithHeight(brs, 0, 0, 20)

	m = pressKey(m, "j")
	m = pressKey(m, "j")

	if m.vp.Offset != 0 {
		t.Errorf("offset = %d, want 0 (3 branches fit in 14 visible rows)", m.vp.Offset)
	}
}

func TestViewportWrapAroundResetsOffset(t *testing.T) {
	// cursor at last item with scrolled viewport; j wraps to 0, offset must go to 0
	brs := manyBranches(10)
	m := makeListModelWithHeight(brs, 9, 0, 10) // cursor=9 (last), maxVisible=4
	m.vp.Offset = 6                              // viewport scrolled down

	m = pressKey(m, "j") // wraps to cursor=0

	if m.cursor != 0 {
		t.Fatalf("cursor = %d, want 0 (wrap around)", m.cursor)
	}
	if m.vp.Offset != 0 {
		t.Errorf("offset = %d, want 0 after wrap to top", m.vp.Offset)
	}
}

func TestViewportClampsOnWindowResize(t *testing.T) {
	// cursor=8 in a large terminal; resize to small → offset must clamp so cursor stays visible
	brs := manyBranches(10)
	m := makeListModelWithHeight(brs, 8, 0, 20) // big terminal, cursor=8
	m.vp.Offset = 0

	updated, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 10}) // maxVisible=4
	m = updated.(model)

	maxVis := m.vp.MaxVisible(len(m.branches))
	if m.cursor < m.vp.Offset || m.cursor >= m.vp.Offset+maxVis {
		t.Errorf("cursor %d not visible in viewport [%d, %d) after resize", m.cursor, m.vp.Offset, m.vp.Offset+maxVis)
	}
}

// --- Stale m.current fix: checkout updates m.current before entering ForceDeleteConfirmMode ---

// TestCurrentUpdatedAfterCheckoutInDeleteFlow verifies that when deleting the current branch
// fails with "not fully merged", m.current is updated to reflect the newly checked-out branch
// before entering ForceDeleteConfirmMode. Without the fix, m.current stays stale (== cursor),
// causing updateForceDelete to mis-identify wasCurrentBranch and land the cursor at 0 instead
// of deletedIdx-1 after the force delete.
func TestCurrentUpdatedAfterCheckoutInDeleteFlow(t *testing.T) {
	origCheckout := checkoutBranch
	checkoutBranch = func(_ string) error { return nil }
	defer func() { checkoutBranch = origCheckout }()

	origDelete := deleteBranch
	deleteBranch = func(_ string) error {
		return fmt.Errorf("error: The branch 'current' is not fully merged")
	}
	defer func() { deleteBranch = origDelete }()

	origForce := forceDeleteBranch
	forceDeleteBranch = func(_ string) error { return nil }
	defer func() { forceDeleteBranch = origForce }()

	origList := listBranches
	listBranches = func() ([]Branch, error) {
		return []Branch{
			{Name: "main", IsCurrent: true},
			{Name: "feat"},
		}, nil
	}
	defer func() { listBranches = origList }()

	// branches[2] ("current") is the current branch; cursor is also on it
	brs := []Branch{
		{Name: "main"},
		{Name: "feat"},
		{Name: "current", IsCurrent: true},
	}
	m := makeListModel(brs, 2, 2)
	m.mode = tui.DeleteConfirmMode

	// Confirm delete: checkout succeeds, deleteBranch fails → ForceDeleteConfirmMode
	m = pressKey(m, "y")

	if m.mode != tui.ForceDeleteConfirmMode {
		t.Fatalf("mode after unmerged delete = %v, want ForceDeleteConfirmMode", m.mode)
	}
	// m.current must be updated to main (index 0) after checkout, not left stale at 2
	if m.current != 0 {
		t.Errorf("m.current after checkout = %d, want 0 (main); stale m.current not fixed", m.current)
	}

	// Confirm force delete
	m = pressKey(m, "y")

	if m.mode != tui.ListMode {
		t.Fatalf("mode after force delete = %v, want ListMode", m.mode)
	}
	// cursor should be at deletedIdx-1 = 1, not 0
	// (with stale m.current, wasCurrentBranch was wrongly true → cursor went to 0)
	if m.cursor != 1 {
		t.Errorf("cursor after force delete = %d, want 1 (deletedIdx-1); cursor misplaced", m.cursor)
	}
}

// --- postDeleteRefresh helper ---

func TestPostDeleteRefreshOnListError(t *testing.T) {
	origList := listBranches
	listBranches = func() ([]Branch, error) {
		return nil, fmt.Errorf("git error")
	}
	defer func() { listBranches = origList }()

	brs := threeBranches()
	m := makeListModel(brs, 1, 0)
	m.mode = tui.DeleteConfirmMode

	result := postDeleteRefresh(m, 1, false)

	if result.err == nil {
		t.Error("expected err to be set on listBranches failure")
	}
	if result.mode != tui.ListMode {
		t.Errorf("mode = %v, want ListMode", result.mode)
	}
}

func TestPostDeleteRefreshCursorWasCurrentBranch(t *testing.T) {
	origList := listBranches
	listBranches = func() ([]Branch, error) {
		return []Branch{
			{Name: "main", IsCurrent: true},
			{Name: "feat"},
		}, nil
	}
	defer func() { listBranches = origList }()

	brs := threeBranches()
	m := makeListModel(brs, 2, 2)

	result := postDeleteRefresh(m, 2, true)

	if result.cursor != 0 {
		t.Errorf("cursor = %d, want 0 when wasCurrentBranch=true", result.cursor)
	}
	if result.mode != tui.ListMode {
		t.Errorf("mode = %v, want ListMode", result.mode)
	}
	if result.err != nil {
		t.Errorf("err = %v, want nil", result.err)
	}
}

func TestPostDeleteRefreshCursorDeletedAtZero(t *testing.T) {
	origList := listBranches
	listBranches = func() ([]Branch, error) {
		return []Branch{
			{Name: "feat", IsCurrent: true},
			{Name: "fix"},
		}, nil
	}
	defer func() { listBranches = origList }()

	brs := threeBranches()
	m := makeListModel(brs, 0, 1)

	result := postDeleteRefresh(m, 0, false)

	if result.cursor != 0 {
		t.Errorf("cursor = %d, want 0 when deletedIdx=0", result.cursor)
	}
}

func TestPostDeleteRefreshCursorNormal(t *testing.T) {
	origList := listBranches
	listBranches = func() ([]Branch, error) {
		return []Branch{
			{Name: "main", IsCurrent: true},
			{Name: "feat"},
		}, nil
	}
	defer func() { listBranches = origList }()

	brs := threeBranches()
	m := makeListModel(brs, 2, 0)

	result := postDeleteRefresh(m, 2, false)

	if result.cursor != 1 {
		t.Errorf("cursor = %d, want 1 (deletedIdx-1)", result.cursor)
	}
}

func TestPostDeleteRefreshUpdatesBranchesAndCurrent(t *testing.T) {
	newBranches := []Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feat"},
	}
	origList := listBranches
	listBranches = func() ([]Branch, error) { return newBranches, nil }
	defer func() { listBranches = origList }()

	brs := threeBranches()
	m := makeListModel(brs, 1, 0)

	result := postDeleteRefresh(m, 1, false)

	if len(result.branches) != 2 {
		t.Errorf("branches len = %d, want 2", len(result.branches))
	}
	if result.current != 0 {
		t.Errorf("current = %d, want 0 (main is IsCurrent)", result.current)
	}
}

// --- Enter on worktree item ---

// TestEnterOnWorktreeShowsStatusMsg verifies that pressing enter while the cursor
// is on a worktree branch shows a statusMsg and does not quit the program.
func TestEnterOnWorktreeShowsStatusMsg(t *testing.T) {
	brs := []Branch{
		{Name: "main", IsCurrent: true},
		{Name: "feat-wt", IsWorktree: true},
	}
	// Force cursor onto the worktree entry (bypassing normal navigation guards)
	m := makeListModel(brs, 1, 0)

	enterMsg := tea.KeyMsg{Type: tea.KeyEnter}
	updated, cmd := m.Update(enterMsg)
	result := updated.(model)

	if cmd != nil {
		t.Error("expected nil cmd (no quit) when enter on worktree branch")
	}
	if result.statusMsg == "" {
		t.Error("expected statusMsg to be set when enter on worktree branch")
	}
	if result.selected != "" {
		t.Errorf("selected = %q, want empty when enter on worktree branch", result.selected)
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
