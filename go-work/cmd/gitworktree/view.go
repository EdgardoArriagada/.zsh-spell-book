package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func hint(key, desc string) string {
	return keyStyle.Render(key) + dimStyle.Render(" "+desc)
}

func (m model) View() string {
	if m.err != nil && len(m.worktrees) == 0 {
		return errStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
	}

	sep := dimStyle.Render("  ·  ")

	var s strings.Builder
	s.WriteString(titleStyle.Render("  Git Worktrees") + "\n\n")

	for i, wt := range m.worktrees {
		cursor := "   "
		if i == m.cursor {
			cursor = " " + cursorStyle.Render("▸ ")
		}

		name := filepath.Base(wt.Path)
		branch := ""
		if wt.Branch != "" {
			branch = " " + branchStyle.Render(wt.Branch)
		}
		if wt.IsBare {
			branch += dimStyle.Render(" bare")
		}

		var line string
		switch {
		case i == m.current && i == m.cursor:
			line = currentMark.Render(name) + branch + currentMark.Render("  ●")
		case i == m.current:
			line = dimStyle.Render(name) + dimStyle.Render(branch) + currentMark.Render("  ●")
		case i == m.cursor:
			line = activeStyle.Render(name) + branch
		default:
			line = dimStyle.Render(name) + dimStyle.Render(branch)
		}

		s.WriteString(cursor + line + "\n")
	}

	switch m.mode {
	case addMode:
		s.WriteString("\n" + promptStyle.Render("  New branch") + "\n  ❱")
		s.WriteString(m.input.View() + "\n")
		if m.err != nil {
			s.WriteString(errStyle.Render(fmt.Sprintf("  %v", m.err)) + "\n")
		}
	case deleteConfirmMode:
		wt := m.worktrees[m.cursor]
		s.WriteString("\n" + warnStyle.Render(
			fmt.Sprintf("  Delete %s [%s]? (y/n)", filepath.Base(wt.Path), wt.Branch),
		) + "\n")
	case forceDeleteConfirmMode:
		wt := m.worktrees[m.cursor]
		s.WriteString("\n" + errStyle.Render(
			fmt.Sprintf("  %s has uncommitted changes — force delete? (y/n)", filepath.Base(wt.Path)),
		) + "\n")
	default:
		if m.statusMsg != "" {
			s.WriteString("\n" + statusStyle.Render("  "+m.statusMsg) + "\n")
		} else if m.err != nil {
			s.WriteString("\n" + errStyle.Render(fmt.Sprintf("  %v", m.err)) + "\n")
		}
	}

	footer := "  " + hint("j/k", "navigate") + sep +
		hint("enter", "select") + sep +
		hint("a", "add") + sep +
		hint("d", "delete") + sep +
		hint("esc/q", "quit")
	s.WriteString("\n" + footer + "\n")
	return s.String()
}
