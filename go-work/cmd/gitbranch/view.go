package main

import (
	"fmt"
	"strings"

	"example.com/workspace/lib/tui"
)

func (m model) View() string {
	if m.err != nil && len(m.branches) == 0 {
		return tui.ErrStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
	}

	sep := tui.Sep()

	var s strings.Builder
	s.WriteString(tui.TitleStyle.Render("  Git Branches") + "\n\n")

	for i, br := range m.branches {
		cursor := "   "
		if i == m.cursor {
			cursor = " " + tui.CursorStyle.Render("▸ ")
		}

		var line string
		switch {
		case i == m.current && i == m.cursor:
			line = tui.CurrentMark.Render(br.Name) + tui.CurrentMark.Render("  ●")
		case i == m.current:
			line = tui.DimStyle.Render(br.Name) + tui.CurrentMark.Render("  ●")
		case i == m.cursor:
			line = tui.ActiveStyle.Render(br.Name)
		default:
			line = tui.DimStyle.Render(br.Name)
		}

		s.WriteString(cursor + line + "\n")
	}

	switch m.mode {
	case tui.AddMode:
		s.WriteString("\n" + tui.PromptStyle.Render("  New branch") + "\n  ❱")
		s.WriteString(m.input.View() + "\n")
		if m.err != nil {
			s.WriteString(tui.ErrStyle.Render(fmt.Sprintf("  %v", m.err)) + "\n")
		}
	case tui.DeleteConfirmMode:
		br := m.branches[m.cursor]
		msg := fmt.Sprintf("  Delete branch %s? (y/n)", br.Name)
		if m.cursor == m.current {
			msg = fmt.Sprintf("  Switch away and delete branch %s? (y/n)", br.Name)
		}
		s.WriteString("\n" + tui.WarnStyle.Render(msg) + "\n")
	case tui.ForceDeleteConfirmMode:
		br := m.branches[m.cursor]
		s.WriteString("\n" + tui.ErrStyle.Render(
			fmt.Sprintf("  %s is not fully merged — force delete? (y/n)", br.Name),
		) + "\n")
	default:
		if m.statusMsg != "" {
			s.WriteString("\n" + tui.StatusStyle.Render("  "+m.statusMsg) + "\n")
		} else if m.err != nil {
			s.WriteString("\n" + tui.ErrStyle.Render(fmt.Sprintf("  %v", m.err)) + "\n")
		}
	}

	footer := "  " + tui.Hint("j/k", "navigate") + sep +
		tui.Hint("enter", "select") + sep +
		tui.Hint("a", "add") + sep +
		tui.Hint("d", "delete") + sep +
		tui.Hint("esc/q", "quit")
	s.WriteString("\n" + footer + "\n")
	return s.String()
}
