package main

import (
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case addMode:
		return m.updateAdd(msg)
	case deleteConfirmMode:
		return m.updateDelete(msg)
	default:
		return m.updateList(msg)
	}
}

func (m model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch km.String() {
	case "q", "ctrl+c", "esc":
		return m, tea.Quit
	case "j", "down", "tab":
		if m.cursor < len(m.worktrees)-1 {
			m.cursor++
		}
	case "k", "up", "shift+tab":
		if m.cursor > 0 {
			m.cursor--
		}
	case "enter":
		if len(m.worktrees) > 0 {
			m.selected = m.worktrees[m.cursor].Path
		}
		return m, tea.Quit
	case "a":
		m.mode = addMode
		m.input.SetValue("")
		m.err = nil
		return m, m.input.Focus()
	case "d":
		if len(m.worktrees) > 0 && m.cursor != 0 {
			m.mode = deleteConfirmMode
		}
	}
	return m, nil
}

func (m model) updateAdd(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.mode = listMode
			m.input.Blur()
			m.err = nil
			return m, nil
		case "enter":
			branch := strings.TrimSpace(m.input.Value())
			if branch == "" {
				return m, nil
			}
			mainPath := m.worktrees[0].Path
			if err := createWorktree(mainPath, branch); err != nil {
				m.err = err
				return m, nil
			}
			wts, err := listWorktrees()
			if err != nil {
				m.err = err
				return m, nil
			}
			m.worktrees = wts
			m.current = currentWorktreeIndex(wts)
			m.mode = listMode
			m.input.Blur()
			m.err = nil
			for i, wt := range wts {
				if wt.Branch == branch {
					m.cursor = i
					break
				}
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) updateDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "y", "Y":
			deletingCurrent := m.cursor == m.current
			if err := deleteWorktree(m.worktrees[m.cursor].Path); err != nil {
				m.err = err
				m.mode = listMode
				return m, nil
			}
			// If we deleted the current worktree, chdir to main
			// so the process has a valid cwd.
			if deletingCurrent && len(m.worktrees) > 0 {
				os.Chdir(m.worktrees[0].Path)
			}
			wts, err := listWorktrees()
			if err != nil {
				m.err = err
				m.mode = listMode
				return m, nil
			}
			m.worktrees = wts
			m.current = currentWorktreeIndex(wts)
			if m.cursor >= len(m.worktrees) {
				m.cursor = len(m.worktrees) - 1
			}
			m.mode = listMode
			m.err = nil
		default:
			m.mode = listMode
		}
	}
	return m, nil
}
