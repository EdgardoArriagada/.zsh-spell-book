package main

import (
	"fmt"
	"path/filepath"
	"strings"

	"example.com/workspace/lib/tui"
)

func (m model) statusSection() string {
	var target Worktree
	if len(m.filtered) > 0 && m.cursor < len(m.filtered) {
		target = m.filtered[m.cursor]
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
		confirmLabel := "Delete"
		if m.deleteBranch {
			confirmLabel = "Delete + branch"
		}
		return "\n" + tui.WarnStyle.Render(
			fmt.Sprintf("  %s %s [%s]? (y/n)", confirmLabel, filepath.Base(target.Path), target.Branch),
		) + "\n"
	case tui.ForceDeleteConfirmMode:
		forceSuffix := "force delete"
		if m.deleteBranch {
			forceSuffix = "force delete + branch"
		}
		return "\n" + tui.ErrStyle.Render(
			fmt.Sprintf("  %s has uncommitted changes — %s? (y/n)", filepath.Base(target.Path), forceSuffix),
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
	switch m.mode {
	case tui.SearchMode:
		return tui.SearchFooter()
	case tui.AddMode:
		return "  " + tui.Hint("enter", "confirm") + sep +
			tui.Hint("ctrl+g", "open $EDITOR") + sep +
			tui.Hint("esc", "cancel")
	default:
		return "  " + tui.Hint("↑/↓", "navigate") + sep +
			tui.Hint("enter", "select") + sep +
			tui.Hint("a", "add") + sep +
			tui.Hint("d", "delete") + sep +
			tui.Hint("D", "delete+branch") + sep +
			tui.Hint("/", "search") + sep +
			tui.Hint("esc/q", "quit")
	}
}

func (m model) availableRows() int {
	return tui.AvailableRows(m.windowHeight, tui.Title("Git Worktrees"), m.statusSection(), "\n"+m.footerSection()+"\n")
}

func (m model) View() string {
	if m.err != nil && len(m.worktrees) == 0 {
		return tui.ErrStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
	}

	var currentPath string
	if m.current >= 0 && m.current < len(m.worktrees) {
		currentPath = m.worktrees[m.current].Path
	}

	var s strings.Builder
	s.WriteString(tui.Title("Git Worktrees"))

	maxVis := m.vp.MaxVisible(len(m.filtered), m.availableRows())
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

	actualDisplayed := end - m.vp.Offset
	if padding := maxVis - actualDisplayed; padding > 0 {
		s.WriteString(strings.Repeat("\n", padding))
	}

	s.WriteString(m.statusSection())
	s.WriteString("\n" + m.footerSection() + "\n")
	return s.String()
}
