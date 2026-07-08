package main

import (
	"fmt"
	"strings"

	"example.com/workspace/lib/tui"
)

func (m model) statusSection() string {
	if m.mode == tui.SearchMode {
		return tui.RenderSearchInput(m.searchInput)
	}
	return tui.RenderActiveFilterHint(m.searchInput)
}

func (m model) footerSection() string {
	sep := tui.Sep()
	if m.mode == tui.SearchMode {
		return tui.SearchFooter()
	}
	return "  " + tui.Hint("↑/↓", "navigate") + sep +
		tui.Hint("enter", "select") + sep +
		tui.Hint("/", "search") + sep +
		tui.Hint("ctrl+g", "edit tickets") + sep +
		tui.Hint("esc/q", "quit")
}

func (m model) availableRows() int {
	return tui.AvailableRows(m.windowHeight, tui.Title("Jira Tickets"), m.statusSection(), "\n"+m.footerSection()+"\n")
}

func (m model) View() string {
	if m.err != nil && len(m.tickets) == 0 {
		return tui.ErrStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
	}

	var s strings.Builder
	s.WriteString(tui.Title("Jira Tickets"))

	var currentTicket string
	if m.current >= 0 && m.current < len(m.tickets) {
		currentTicket = m.tickets[m.current].Current
	}

	maxVis := m.vp.MaxVisible(len(m.filtered), m.availableRows())
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
		isCurrent := currentTicket != "" && t.Current == currentTicket

		var renderedLine string
		switch {
		case isCurrent && idx == m.cursor:
			renderedLine = tui.CurrentMark.Render(line) + tui.CurrentMark.Render("  ●")
		case isCurrent:
			renderedLine = tui.DimStyle.Render(line) + tui.CurrentMark.Render("  ●")
		case idx == m.cursor:
			renderedLine = tui.ActiveStyle.Render(line)
		default:
			renderedLine = tui.DimStyle.Render(line)
		}
		s.WriteString(cursor + renderedLine + "\n")
	}

	actualDisplayed := end - m.vp.Offset
	if padding := maxVis - actualDisplayed; padding > 0 {
		s.WriteString(strings.Repeat("\n", padding))
	}

	s.WriteString(m.statusSection())
	s.WriteString("\n" + m.footerSection() + "\n")
	return s.String()
}
