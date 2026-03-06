package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"example.com/workspace/lib/tui"
)

func (m model) View() string {
	if m.err != nil && len(m.worktrees) == 0 {
		return tui.ErrStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
	}

	sep := tui.Sep()
	var currentPath string
	if m.current >= 0 && m.current < len(m.worktrees) {
		currentPath = m.worktrees[m.current].Path
	}

	var s strings.Builder
	s.WriteString(tui.Title("Git Worktrees"))

	maxVis := m.vp.MaxVisible(len(m.filtered))
	end := m.vp.Offset + maxVis
	if end > len(m.filtered) {
		end = len(m.filtered)
	}
	for i, wt := range m.filtered[m.vp.Offset:end] {
		idx := i + m.vp.Offset
		cursor := "   "
		if idx == m.cursor {
			cursor = " " + tui.CursorStyle.Render("▸ ")
		}

		name := filepath.Base(wt.Path)
		branch := ""
		if wt.Branch != "" {
			branch = " " + tui.BranchStyle.Render(wt.Branch)
		}
		if wt.IsBare {
			branch += tui.DimStyle.Render(" bare")
		}

		isCurrent := currentPath != "" && wt.Path == currentPath

		var line string
		switch {
		case isCurrent && idx == m.cursor:
			line = tui.CurrentMark.Render(name) + branch + tui.CurrentMark.Render("  ●")
		case isCurrent:
			line = tui.DimStyle.Render(name) + tui.DimStyle.Render(branch) + tui.CurrentMark.Render("  ●")
		case idx == m.cursor:
			line = tui.ActiveStyle.Render(name) + branch
		default:
			line = tui.DimStyle.Render(name) + tui.DimStyle.Render(branch)
		}

		s.WriteString(cursor + line + "\n")
	}

	var target Worktree
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
		s.WriteString("\n" + tui.WarnStyle.Render(
			fmt.Sprintf("  Delete %s [%s]? (y/n)", filepath.Base(target.Path), target.Branch),
		) + "\n")
	case tui.ForceDeleteConfirmMode:
		s.WriteString("\n" + tui.ErrStyle.Render(
			fmt.Sprintf("  %s has uncommitted changes — force delete? (y/n)", filepath.Base(target.Path)),
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
	}
	s.WriteString("\n" + footer + "\n")
	return s.String()
}
