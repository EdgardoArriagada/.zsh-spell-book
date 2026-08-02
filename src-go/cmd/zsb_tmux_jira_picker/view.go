package main

import (
	"fmt"
	"strconv"
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

func (m model) View() string {
	if m.err != nil && len(m.tickets) == 0 {
		return tui.ErrStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
	}

	title := tui.Title("Jira Tickets")
	statusSec := m.statusSection()
	footerSec := "\n" + m.footerSection() + "\n"

	var s strings.Builder
	s.Grow(len(m.filtered) * 80)
	s.WriteString(title)

	var currentTicket string
	if m.current >= 0 && m.current < len(m.tickets) {
		currentTicket = m.tickets[m.current].Current
	}

	maxVis := m.vp.MaxVisible(len(m.filtered), m.availRows)
	end := min(m.vp.Offset+maxVis, len(m.filtered))
	for i, t := range m.filtered[m.vp.Offset:end] {
		idx := i + m.vp.Offset
		cursor := "   "
		if idx == m.cursor {
			cursor = " " + tui.CursorStyle.Render("▸ ")
		}
		const cursorWidth = 3
		const badgeReserve = 11 // "  ●" (3) + " 999" working + " 999" finished (~4 each)
		fixedWidth := len(t.Current) + 2 // ponytail: JIRA IDs are ASCII-only, byte len == rune len
		label := truncateLabel(t.Label, m.width-cursorWidth-fixedWidth-badgeReserve)
		line := t.Current + ": " + label
		isCurrent := currentTicket != "" && t.Current == currentTicket

		var renderedLine string
		switch {
		case isCurrent && idx == m.cursor:
			renderedLine = tui.CurrentMark.Render(line) + tui.CurrentMark.Render(" ●")
		case isCurrent:
			renderedLine = tui.DimStyle.Render(line) + tui.CurrentMark.Render(" ●")
		case idx == m.cursor:
			renderedLine = tui.ActiveStyle.Render(line)
		default:
			renderedLine = tui.DimStyle.Render(line)
		}
		c := m.notifCounts[t.SessionID]
		if c.Working > 0 {
			renderedLine += tui.NotifWorking.Render(" " + strconv.Itoa(c.Working))
		}
		if c.Finished > 0 {
			renderedLine += tui.NotifBadge.Render(" " + strconv.Itoa(c.Finished))
		}
		s.WriteString(cursor)
		s.WriteString(renderedLine)
		s.WriteByte('\n')
	}

	actualDisplayed := end - m.vp.Offset
	if padding := maxVis - actualDisplayed; padding > 0 {
		s.WriteString(strings.Repeat("\n", padding))
	}

	s.WriteString(statusSec)
	s.WriteString(footerSec)
	return s.String()
}

func truncateLabel(s string, maxWidth int) string {
	if maxWidth < 1 {
		maxWidth = 1
	}
	runes := []rune(s)
	if len(runes) <= maxWidth {
		return s
	}
	return string(runes[:maxWidth-1]) + "…"
}
