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
	filterActive := m.searchInput.Value() != ""

	var currentName string
	if m.current >= 0 && m.current < len(m.branches) {
		currentName = m.branches[m.current].Name
	}

	var s strings.Builder
	s.WriteString(tui.Title("Git Branches"))
	if m.inWorktree {
		s.WriteString(tui.WarnStyle.Render("  ⚠  You are inside a worktree — branch management here is not recommended") + "\n\n")
	}

	maxVis := m.vp.MaxVisible(len(m.filtered))
	end := m.vp.Offset + maxVis
	if end > len(m.filtered) {
		end = len(m.filtered)
	}
	for i, br := range m.filtered[m.vp.Offset:end] {
		idx := i + m.vp.Offset
		// Only show divider when no filter active
		if !filterActive && idx == m.firstWorktreeIdx {
			s.WriteString(tui.Title("Worktree Branches"))
		}
		cursor := "   "
		if idx == m.cursor {
			cursor = " " + tui.CursorStyle.Render("▸ ")
		}

		isCurrent := currentName != "" && br.Name == currentName

		var line string
		switch {
		case br.IsWorktree:
			line = tui.WorktreeStyle.Render(br.Name) + tui.WorktreeStyle.Render()
		case isCurrent && idx == m.cursor:
			line = tui.CurrentMark.Render(br.Name) + tui.CurrentMark.Render("  ●")
		case isCurrent:
			line = tui.DimStyle.Render(br.Name) + tui.CurrentMark.Render("  ●")
		case idx == m.cursor:
			line = tui.ActiveStyle.Render(br.Name)
		default:
			line = tui.DimStyle.Render(br.Name)
		}

		s.WriteString(cursor + line + "\n")
	}

	var target Branch
	if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
		target = m.filtered[m.cursor]
	}

	switch m.mode {
	case tui.AddMode:
		s.WriteString("\n" + tui.PromptStyle.Render("  New branch") + "\n  ❱")
		s.WriteString(m.input.View() + "\n")
		if m.err != nil {
			s.WriteString(tui.ErrStyle.Render(fmt.Sprintf("  %v", m.err)) + "\n")
		}
	case tui.DeleteConfirmMode:
		msg := fmt.Sprintf("  Delete branch %s? (y/n)", target.Name)
		if currentName != "" && target.Name == currentName {
			msg = fmt.Sprintf("  Switch away and delete branch %s? (y/n)", target.Name)
		}
		s.WriteString("\n" + tui.WarnStyle.Render(msg) + "\n")
	case tui.ForceDeleteConfirmMode:
		s.WriteString("\n" + tui.ErrStyle.Render(
			fmt.Sprintf("  %s is not fully merged — force delete? (y/n)", target.Name),
		) + "\n")
	case tui.SearchMode:
		s.WriteString("\n  " + m.searchInput.View() + "\n")
	default:
		if m.statusMsg != "" {
			s.WriteString("\n" + tui.StatusStyle.Render("  "+m.statusMsg) + "\n")
		} else if m.err != nil {
			s.WriteString("\n" + tui.ErrStyle.Render(fmt.Sprintf("  %v", m.err)) + "\n")
		} else if term := m.searchInput.Value(); term != "" {
			s.WriteString("\n" + tui.DimStyle.Render("  / "+term) + "\n")
		}
	}

	var footer string
	if m.mode == tui.SearchMode {
		footer = "  " + tui.Hint("enter", "confirm") + sep + tui.Hint("esc", "clear")
	} else {
		footer = "  " + tui.Hint("↑/↓", "navigate") + sep +
			tui.Hint("enter", "select") + sep +
			tui.Hint("a", "add") + sep +
			tui.Hint("d", "delete") + sep +
			tui.Hint("/", "search") + sep +
			tui.Hint("esc/q", "quit")
		if !filterActive && m.firstWorktreeIdx >= 0 {
			footer += sep + tui.Hint("wt", "checked out elsewhere")
		}
	}
	s.WriteString("\n" + footer + "\n")
	return s.String()
}
