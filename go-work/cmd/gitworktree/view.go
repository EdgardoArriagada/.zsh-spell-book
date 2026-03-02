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

	var s strings.Builder
	s.WriteString(tui.Title("Git Worktrees"))

	maxVis := m.vp.MaxVisible(len(m.worktrees))
	end := m.vp.Offset + maxVis
	if end > len(m.worktrees) {
		end = len(m.worktrees)
	}
	for i, wt := range m.worktrees[m.vp.Offset:end] {
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

		var line string
		switch {
		case idx == m.current && idx == m.cursor:
			line = tui.CurrentMark.Render(name) + branch + tui.CurrentMark.Render("  ●")
		case idx == m.current:
			line = tui.DimStyle.Render(name) + tui.DimStyle.Render(branch) + tui.CurrentMark.Render("  ●")
		case idx == m.cursor:
			line = tui.ActiveStyle.Render(name) + branch
		default:
			line = tui.DimStyle.Render(name) + tui.DimStyle.Render(branch)
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
		wt := m.worktrees[m.cursor]
		s.WriteString("\n" + tui.WarnStyle.Render(
			fmt.Sprintf("  Delete %s [%s]? (y/n)", filepath.Base(wt.Path), wt.Branch),
		) + "\n")
	case tui.ForceDeleteConfirmMode:
		wt := m.worktrees[m.cursor]
		s.WriteString("\n" + tui.ErrStyle.Render(
			fmt.Sprintf("  %s has uncommitted changes — force delete? (y/n)", filepath.Base(wt.Path)),
		) + "\n")
	default:
		if m.statusMsg != "" {
			s.WriteString("\n" + tui.StatusStyle.Render("  "+m.statusMsg) + "\n")
		} else if m.err != nil {
			s.WriteString("\n" + tui.ErrStyle.Render(fmt.Sprintf("  %v", m.err)) + "\n")
		}
	}

	footer := "  " + tui.Hint("j/k", "navigate") + sep +
		tui.Hint("g/G", "top/bottom") + sep +
		tui.Hint("enter", "select") + sep +
		tui.Hint("a", "add") + sep +
		tui.Hint("d", "delete") + sep +
		tui.Hint("esc/q", "quit")
	s.WriteString("\n" + footer + "\n")
	return s.String()
}
