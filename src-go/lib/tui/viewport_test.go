package tui_test

import (
	"testing"

	"example.com/workspace/lib/tui"

	"charm.land/lipgloss/v2"
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

func TestHighlightMatchesStylesEveryFuzzyRune(t *testing.T) {
	base := lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9")).Underline(true)
	highlight := tui.MatchStyle.Inherit(base)
	want := highlight.Render("g") + base.Render("it") + highlight.Render("b") + base.Render("ranch")
	if got := tui.HighlightMatches("gitbranch", "gb", base); got != want {
		t.Fatalf("highlighted text = %q, want %q", got, want)
	}
}

func TestActiveSearchTermOnlyReturnsValueInSearchMode(t *testing.T) {
	si := tui.NewSearchInput()
	si.SetValue("feat")
	if got := tui.ActiveSearchTerm(tui.SearchMode, si); got != "feat" {
		t.Fatalf("search term = %q, want feat", got)
	}
	if got := tui.ActiveSearchTerm(tui.ListMode, si); got != "" {
		t.Fatalf("confirmed search term = %q, want empty", got)
	}
}
