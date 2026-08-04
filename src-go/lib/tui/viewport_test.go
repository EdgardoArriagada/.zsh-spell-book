package tui

import "testing"

func TestPageCursor(t *testing.T) {
	if got := PageCursor(2, 10, 4, 1); got != 6 {
		t.Fatalf("down = %d, want 6", got)
	}
	if got := PageCursor(6, 10, 4, -1); got != 2 {
		t.Fatalf("up = %d, want 2", got)
	}
	if got := PageCursor(8, 10, 4, 1); got != 9 {
		t.Fatalf("bottom = %d, want 9", got)
	}
}
