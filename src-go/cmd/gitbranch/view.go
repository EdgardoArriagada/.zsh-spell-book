package main

import (
	"fmt"
	"strings"

	"example.com/workspace/lib/tui"
)

func (m model) headerSection() string {
	s := tui.Title("Git Branches")
	if m.inWorktree {
		s += tui.WarnStyle.Render("  ⚠  You are inside a worktree — branch management here is not recommended") + "\n\n"
	}
	return s
}

func (m model) statusSection() string {
	var target Branch
	if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
		target = m.filtered[m.cursor]
	}
	var currentName string
	if m.current >= 0 && m.current < len(m.branches) {
		currentName = m.branches[m.current].Name
	}
	switch m.mode {
	case tui.AddMode:
		s := "\n" + tui.PromptStyle.Render("  New branch") + "\n  ❱"
		s += m.input.View() + "\n"
		if m.err != nil {
			s += tui.ErrStyle.Render(fmt.Sprintf("  %v", m.err)) + "\n"
		}
		return s
	case tui.DeleteConfirmMode:
		msg := fmt.Sprintf("  Delete branch %s? (y/n)", target.Name)
		if currentName != "" && target.Name == currentName {
			msg = fmt.Sprintf("  Switch away and delete branch %s? (y/n)", target.Name)
		}
		return "\n" + tui.WarnStyle.Render(msg) + "\n"
	case tui.ForceDeleteConfirmMode:
		return "\n" + tui.ErrStyle.Render(
			fmt.Sprintf("  %s is not fully merged — force delete? (y/n)", target.Name),
		) + "\n"
	case tui.SearchMode:
		return tui.RenderSearchInput(m.searchInput)
	default:
		if m.statusMsg != "" {
			return "\n" + tui.StatusStyle.Render("  "+m.statusMsg) + "\n"
		} else if m.err != nil {
			return "\n" + tui.ErrStyle.Render(fmt.Sprintf("  %v", m.err)) + "\n"
		}
		return tui.RenderActiveFilterHint(m.searchInput)
	}
}

func (m model) footerSection() string {
	sep := tui.Sep()
	filterActive := m.searchInput.Value() != ""
	switch m.mode {
	case tui.SearchMode:
		return tui.SearchFooter()
	case tui.AddMode:
		return "  " + tui.Hint("enter", "confirm") + sep +
			tui.Hint("ctrl+g", "open $EDITOR") + sep +
			tui.Hint("esc", "cancel")
	default:
		footer := "  " + tui.Hint("↑/↓", "navigate") + sep +
			tui.Hint("enter", "select") + sep +
			tui.Hint("a", "add") + sep +
			tui.Hint("d", "delete") + sep +
			tui.Hint("/", "search") + sep +
			tui.Hint("esc/q", "quit")
		if !filterActive && m.firstWorktreeIdx >= 0 {
			footer += sep + tui.Hint("wt", "checked out elsewhere")
		}
		return footer
	}
}

func (m model) availableRows() int {
	return tui.AvailableRows(m.windowHeight, m.headerSection(), m.statusSection(), "\n"+m.footerSection()+"\n")
}

func (m model) View() string {
	if m.err != nil && len(m.branches) == 0 {
		return tui.ErrStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
	}

	filterActive := m.searchInput.Value() != ""

	var currentName string
	if m.current >= 0 && m.current < len(m.branches) {
		currentName = m.branches[m.current].Name
	}

	var s strings.Builder
	s.WriteString(m.headerSection())

	maxVis := m.vp.MaxVisible(len(m.filtered), m.availableRows())
	end := m.vp.Offset + maxVis
	if end > len(m.filtered) {
		end = len(m.filtered)
	}
	hasDivider := false
	for i, br := range m.filtered[m.vp.Offset:end] {
		idx := i + m.vp.Offset
		if !filterActive && idx == m.firstWorktreeIdx {
			s.WriteString(tui.Title("Worktree Branches"))
			hasDivider = true
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

	actualDisplayed := end - m.vp.Offset
	dividerLines := 0
	if hasDivider {
		dividerLines = 3
	}
	if padding := maxVis - actualDisplayed - dividerLines; padding > 0 {
		s.WriteString(strings.Repeat("\n", padding))
	}

	s.WriteString(m.statusSection())
	s.WriteString("\n" + m.footerSection() + "\n")
	return s.String()
}
