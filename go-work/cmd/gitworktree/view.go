package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

func (m model) View() string {
	if m.err != nil && len(m.worktrees) == 0 {
		return errStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
	}

	var s strings.Builder
	s.WriteString(titleStyle.Render("Git Worktrees") + "\n\n")

	for i, wt := range m.worktrees {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("> ")
		}

		name := filepath.Base(wt.Path)
		branch := ""
		if wt.Branch != "" {
			branch = " " + branchStyle.Render("["+wt.Branch+"]")
		}
		if wt.IsBare {
			branch += dimStyle.Render(" (bare)")
		}

		var line string
		switch {
		case i == m.current:
			line = currentMark.Render(name) + branch + currentMark.Render(" ●")
		case i == m.cursor:
			line = activeStyle.Render(name) + branch
		default:
			line = dimStyle.Render(name) + branch
		}

		s.WriteString(cursor + line + "\n")
	}

	switch m.mode {
	case addMode:
		s.WriteString("\n" + promptStyle.Render("New branch: ") + m.input.View() + "\n")
		if m.err != nil {
			s.WriteString(errStyle.Render(fmt.Sprintf("  %v", m.err)) + "\n")
		}
	case deleteConfirmMode:
		wt := m.worktrees[m.cursor]
		s.WriteString("\n" + errStyle.Render(
			fmt.Sprintf("Delete worktree %s [%s]? (y/n)", filepath.Base(wt.Path), wt.Branch),
		) + "\n")
	default:
		if m.err != nil {
			s.WriteString("\n" + errStyle.Render(fmt.Sprintf("%v", m.err)) + "\n")
		}
	}

	s.WriteString("\n" + dimStyle.Render("j/k navigate • enter select • a add • d delete • esc/q quit") + "\n")
	return s.String()
}
