package tui

import "github.com/charmbracelet/bubbles/textinput"

// RenderSearchInput renders the search input widget for SearchMode.
func RenderSearchInput(si textinput.Model) string {
	return "\n  " + si.View() + "\n"
}

// RenderActiveFilterHint renders a dim hint showing the active filter term,
// or returns "" if no filter is active.
func RenderActiveFilterHint(si textinput.Model) string {
	if term := si.Value(); term != "" {
		return "\n" + DimStyle.Render("  / "+term) + "\n"
	}
	return ""
}

// SearchFooter returns the key hint footer for SearchMode.
func SearchFooter() string {
	return "  " + Hint("enter", "confirm") + Sep() + Hint("esc", "clear")
}
