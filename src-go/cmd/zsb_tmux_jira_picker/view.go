package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode/utf8"

	"example.com/workspace/lib/jira"
	"example.com/workspace/lib/tui"

	tea "charm.land/bubbletea/v2"
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

func (m model) View() tea.View {
	v := tea.NewView(m.render())
	v.AltScreen = true
	return v
}

func (m model) render() string {
	if m.err != nil && len(m.tickets) == 0 {
		return tui.ErrStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
	}

	statusSec := m.statusSection()
	footerSec := "\n" + m.footerSection() + "\n"

	var currentTicket string
	if m.current >= 0 && m.current < len(m.tickets) {
		currentTicket = m.tickets[m.current].Current
	}

	end, usedRows := visibleTicketEnd(m.filtered, m.vp.Offset, m.availRows, m.showTitles())
	var s strings.Builder
	s.Grow((end - m.vp.Offset) * 80)
	for i, t := range m.filtered[m.vp.Offset:end] {
		idx := i + m.vp.Offset
		if m.showTitles() && sectionStarts(m.filtered, idx) && !(idx == m.vp.Offset && m.availRows == 1) {
			title := truncateLabel(t.Title, m.width-2)
			s.WriteByte('\n')
			s.WriteString(tui.TitleStyle.Render("  " + title))
			s.WriteByte('\n')
			s.WriteByte('\n')
		}
		cursor := "   "
		if idx == m.cursor {
			cursor = " " + tui.CursorStyle.Render("▸ ")
		}
		const cursorWidth = 3
		const badgeReserve = 11          // "  ●" (3) + " 999" working + " 999" finished (~4 each)
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

	if padding := m.availRows - usedRows; m.availRows > 0 && padding > 0 {
		s.WriteString(strings.Repeat("\n", padding))
	}

	s.WriteString(statusSec)
	s.WriteString(footerSec)
	return s.String()
}

func sectionStarts(tickets []jira.Ticket, index int) bool {
	return tickets[index].Title != "" && (index == 0 || tickets[index-1].Title != tickets[index].Title)
}

func visibleTicketEnd(tickets []jira.Ticket, start, availableRows int, showTitles bool) (end, usedRows int) {
	if availableRows <= 0 {
		availableRows = defaultAvailableRows
	}
	for i := start; i < len(tickets); i++ {
		rows := 1
		if showTitles && sectionStarts(tickets, i) && !(i == start && availableRows == 1) {
			rows += 3
		}
		if usedRows+rows > availableRows {
			break
		}
		usedRows += rows
		end = i + 1
	}
	return end, usedRows
}

func clampTicketViewport(vp tui.Viewport, cursor int, tickets []jira.Ticket, availableRows int, showTitles bool) tui.Viewport {
	if availableRows <= 0 || len(tickets) == 0 {
		return vp
	}
	if cursor < vp.Offset {
		vp.Offset = cursor
	}
	if end, _ := visibleTicketEnd(tickets, vp.Offset, availableRows, showTitles); cursor == len(tickets)-1 && cursor-end >= availableRows {
		vp.Offset = cursor
	}
	for end, _ := visibleTicketEnd(tickets, vp.Offset, availableRows, showTitles); cursor >= end; end, _ = visibleTicketEnd(tickets, vp.Offset, availableRows, showTitles) {
		vp.Offset++
	}
	for vp.Offset > 0 {
		end, _ := visibleTicketEnd(tickets, vp.Offset-1, availableRows, showTitles)
		if end < len(tickets) {
			break
		}
		vp.Offset--
	}
	return vp
}

func truncateLabel(s string, maxWidth int) string {
	if maxWidth < 1 {
		maxWidth = 1
	}
	if utf8.RuneCountInString(s) <= maxWidth {
		return s
	}
	runes := []rune(s)
	return string(runes[:maxWidth-1]) + "…"
}
