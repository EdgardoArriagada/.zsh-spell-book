package tui

import (
	"strings"
	"unicode/utf8"

	"charm.land/bubbles/v2/textinput"
	"charm.land/lipgloss/v2"
	"github.com/sahilm/fuzzy"
)

// ActiveSearchTerm returns the query only while search mode is active.
func ActiveSearchTerm(mode Mode, si textinput.Model) string {
	if mode == SearchMode {
		return si.Value()
	}
	return ""
}

// HighlightMatches renders fuzzy-matched runes with MatchStyle over base.
func HighlightMatches(text, term string, base lipgloss.Style) string {
	matches := fuzzy.FindNoSort(term, []string{text})
	if len(matches) == 0 {
		return base.Render(text)
	}

	var rendered strings.Builder
	start := 0
	highlight := MatchStyle.Inherit(base)
	for _, index := range matches[0].MatchedIndexes {
		if index > start {
			rendered.WriteString(base.Render(text[start:index]))
		}
		_, size := utf8.DecodeRuneInString(text[index:])
		rendered.WriteString(highlight.Render(text[index : index+size]))
		start = index + size
	}
	rendered.WriteString(base.Render(text[start:]))
	return rendered.String()
}

// RenderSearchInput renders the search input widget for SearchMode.
func RenderSearchInput(si textinput.Model) string {
	return "\n  " + si.View() + "\n"
}

// RenderActiveFilterHint renders a dim hint showing the active filter term.
// Always returns exactly 2 rows to keep the status area height fixed.
func RenderActiveFilterHint(si textinput.Model) string {
	if term := si.Value(); term != "" {
		return "\n" + DimStyle.Render("  / "+term) + "\n"
	}
	return "\n\n"
}

// SearchFooter returns the key hint footer for SearchMode.
func SearchFooter() string {
	return "  " + Hint("enter", "confirm") + Sep() + Hint("esc", "clear")
}
