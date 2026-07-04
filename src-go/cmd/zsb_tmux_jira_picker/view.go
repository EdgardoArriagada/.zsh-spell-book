package main

import (
	"fmt"
	"strings"

	"example.com/workspace/lib/tui"
)

func (m model) View() string {
	if m.err != nil && len(m.tickets) == 0 {
		return tui.ErrStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
	}

	sep := tui.Sep()
	var s strings.Builder
	s.WriteString(tui.Title("Jira Tickets"))

	maxVis := m.vp.MaxVisible(len(m.filtered))
	end := m.vp.Offset + maxVis
	if end > len(m.filtered) {
		end = len(m.filtered)
	}
	for i, t := range m.filtered[m.vp.Offset:end] {
		idx := i + m.vp.Offset
		cursor := "   "
		if idx == m.cursor {
			cursor = " " + tui.CursorStyle.Render("▸ ")
		}
		line := t.Parent + " | " + t.Current + " | " + t.Label
		if idx == m.cursor {
			s.WriteString(cursor + tui.ActiveStyle.Render(line) + "\n")
		} else {
			s.WriteString(cursor + tui.DimStyle.Render(line) + "\n")
		}
	}

	if m.mode == tui.SearchMode {
		s.WriteString(tui.RenderSearchInput(m.searchInput))
	} else {
		s.WriteString(tui.RenderActiveFilterHint(m.searchInput))
	}

	var footer string
	if m.mode == tui.SearchMode {
		footer = tui.SearchFooter()
	} else {
		footer = "  " + tui.Hint("↑/↓", "navigate") + sep +
			tui.Hint("enter", "select") + sep +
			tui.Hint("/", "search") + sep +
			tui.Hint("esc/q", "quit")
	}
	s.WriteString("\n" + footer + "\n")
	return s.String()
}
