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

func TestShortcutIndexAndPrefix(t *testing.T) {
	for _, tc := range []struct {
		key  string
		want int
		ok   bool
	}{
		{"1", 4, true}, {"9", 12, true}, {"0", 13, true},
		{"2", 5, true}, {"x", 0, false}, {"11", 0, false},
	} {
		got, ok := tui.ShortcutIndex(tc.key, 4, 14)
		if ok != tc.ok || ok && got != tc.want {
			t.Fatalf("key %q: index=%d ok=%t, want %d %t", tc.key, got, ok, tc.want, tc.ok)
		}
	}
	if _, ok := tui.ShortcutIndex("0", 4, 13); ok {
		t.Fatal("0 should not select an entry outside the visible page")
	}
	gutter := tui.GutterStyle.Render(" │ ")
	if got := tui.ShortcutPrefix(0, false); got != " "+tui.GutterStyle.Render("1")+gutter {
		t.Fatalf("first prefix = %q", got)
	}
	if got := tui.ShortcutPrefix(9, false); got != " "+tui.GutterStyle.Render("0")+gutter {
		t.Fatalf("tenth prefix = %q", got)
	}
	if got := tui.ShortcutPrefix(10, false); got != "  "+gutter {
		t.Fatalf("eleventh prefix = %q", got)
	}
	if got := tui.ShortcutPrefix(0, true); got != " "+tui.CursorStyle.Render("▸")+gutter {
		t.Fatalf("cursor prefix = %q", got)
	}
	for _, position := range []int{0, 9, 10} {
		if width := lipgloss.Width(tui.ShortcutPrefix(position, false)); width != 5 {
			t.Fatalf("prefix %d width = %d, want 5", position, width)
		}
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
