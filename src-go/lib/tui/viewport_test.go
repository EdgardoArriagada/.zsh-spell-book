package tui_test

import (
	"testing"

	"example.com/workspace/lib/tui"
)

func TestPageCursor(t *testing.T) {
	if got := tui.PageCursor(2, 10, 4, 1); got != 6 {
		t.Fatalf("down = %d, want 6", got)
	}
	if got := tui.PageCursor(6, 10, 4, -1); got != 2 {
		t.Fatalf("up = %d, want 2", got)
	}
	if got := tui.PageCursor(8, 10, 4, 1); got != 9 {
		t.Fatalf("bottom = %d, want 9", got)
	}
}
