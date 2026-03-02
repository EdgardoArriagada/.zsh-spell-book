package main

import (
	"os"

	gitlib "example.com/workspace/lib/git"
	"example.com/workspace/lib/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.vp.Height = ws.Height
		m.input.SetWidth(ws.Width)
		m.vp = m.vp.Clamp(m.cursor, len(m.worktrees))
		return m, nil
	}
	switch m.mode {
	case tui.AddMode:
		return m.updateAdd(msg)
	case tui.DeleteConfirmMode:
		return m.updateDelete(msg)
	case tui.ForceDeleteConfirmMode:
		return m.updateForceDelete(msg)
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
		m.statusMsg = ""
		if len(m.worktrees) > 0 {
			m.cursor = (m.cursor + 1) % len(m.worktrees)
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.worktrees))
	case "k", "up", "shift+tab":
		m.statusMsg = ""
		if len(m.worktrees) > 0 {
			m.cursor = (m.cursor - 1 + len(m.worktrees)) % len(m.worktrees)
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.worktrees))
	case "g":
		m.statusMsg = ""
		m.cursor = 0
		m.vp = m.vp.Clamp(m.cursor, len(m.worktrees))
	case "G":
		m.statusMsg = ""
		if len(m.worktrees) > 0 {
			m.cursor = len(m.worktrees) - 1
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.worktrees))
	case "enter":
		if len(m.worktrees) > 0 {
			m.selected = m.worktrees[m.cursor].Path
		}
		return m, tea.Quit
	case "a":
		m.mode = tui.AddMode
		m.input.SetValue("")
		m.err = nil
		m.statusMsg = ""
		return m, m.input.Focus()
	case "d":
		if len(m.worktrees) == 0 {
			break
		}
		if m.cursor == 0 {
			m.statusMsg = "cannot delete main worktree"
		} else {
			m.statusMsg = ""
			m.mode = tui.DeleteConfirmMode
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
			m.mode = tui.ListMode
			m.input.Blur()
			m.err = nil
			return m, nil
		case "enter":
			branch := tui.ParseInputValue(m.input)
			if branch == "" {
				return m, nil
			}
			if err := gitlib.ValidateBranchName(branch); err != nil {
				m.err = err
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
			m.mode = tui.ListMode
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
				if isWorktreeDirtyError(err) {
					m.mode = tui.ForceDeleteConfirmMode
					m.err = nil
					return m, nil
				}
				m.err = err
				m.mode = tui.ListMode
				return m, nil
			}
			// If we deleted the current worktree, chdir to main
			// so the process has a valid cwd.
			if deletingCurrent && len(m.worktrees) > 0 {
				m.fallbackPath = m.worktrees[0].Path
				os.Chdir(m.fallbackPath)
			}
			wts, err := listWorktrees()
			if err != nil {
				m.err = err
				m.mode = tui.ListMode
				return m, nil
			}
			m.worktrees = wts
			m.current = currentWorktreeIndex(wts)
			if m.cursor >= len(m.worktrees) {
				m.cursor = len(m.worktrees) - 1
			}
			m.mode = tui.ListMode
			m.err = nil
		default:
			m.mode = tui.ListMode
		}
	}
	return m, nil
}

func (m model) updateForceDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "y", "Y":
			deletingCurrent := m.cursor == m.current
			if err := deleteWorktreeForce(m.worktrees[m.cursor].Path); err != nil {
				m.err = err
				m.mode = tui.ListMode
				return m, nil
			}
			if deletingCurrent && len(m.worktrees) > 0 {
				m.fallbackPath = m.worktrees[0].Path
				os.Chdir(m.fallbackPath)
			}
			wts, err := listWorktrees()
			if err != nil {
				m.err = err
				m.mode = tui.ListMode
				return m, nil
			}
			m.worktrees = wts
			m.current = currentWorktreeIndex(wts)
			if m.cursor >= len(m.worktrees) {
				m.cursor = len(m.worktrees) - 1
			}
			m.mode = tui.ListMode
			m.err = nil
		default:
			m.mode = tui.ListMode
		}
	}
	return m, nil
}
