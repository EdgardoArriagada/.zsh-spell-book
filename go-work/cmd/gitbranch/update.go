package main

import (
	gitlib "example.com/workspace/lib/git"
	"example.com/workspace/lib/tui"

	tea "github.com/charmbracelet/bubbletea"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.vp.Height = ws.Height
		m.input.SetWidth(ws.Width)
		m.vp = m.vp.Clamp(m.cursor, len(m.branches))
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
		if len(m.branches) > 0 {
			next := (m.cursor + 1) % len(m.branches)
			for next != m.cursor && m.branches[next].IsWorktree {
				next = (next + 1) % len(m.branches)
			}
			if !m.branches[next].IsWorktree {
				m.cursor = next
			}
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.branches))
	case "k", "up", "shift+tab":
		m.statusMsg = ""
		if len(m.branches) > 0 {
			next := (m.cursor - 1 + len(m.branches)) % len(m.branches)
			for next != m.cursor && m.branches[next].IsWorktree {
				next = (next - 1 + len(m.branches)) % len(m.branches)
			}
			if !m.branches[next].IsWorktree {
				m.cursor = next
			}
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.branches))
	case "enter":
		if len(m.branches) > 0 && m.branches[m.cursor].IsWorktree {
			m.statusMsg = "worktree branches are not selectable"
			return m, nil
		}
		if len(m.branches) > 0 {
			m.selected = m.branches[m.cursor].Name
		}
		return m, tea.Quit
	case "a":
		m.mode = tui.AddMode
		m.input.SetValue("")
		m.err = nil
		m.statusMsg = ""
		return m, m.input.Focus()
	case "d":
		if len(m.branches) == 0 {
			break
		}
		if m.cursor == m.current && len(m.branches) == 1 {
			m.statusMsg = "cannot delete the only branch"
			break
		}
		if isDefaultBranch(m.branches[m.cursor].Name) {
			m.statusMsg = "cannot delete default branch"
			break
		}
		m.statusMsg = ""
		m.mode = tui.DeleteConfirmMode
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
			if err := createBranch(branch); err != nil {
				m.err = err
				return m, nil
			}
			branches, err := listBranches()
			if err != nil {
				m.err = err
				return m, nil
			}
			m.branches = branches
			m.firstWorktreeIdx = findFirstWorktreeIdx(branches)
			for i, b := range branches {
				if b.IsCurrent {
					m.current = i
					break
				}
			}
			m.mode = tui.ListMode
			m.input.Blur()
			m.err = nil
			for i, b := range branches {
				if b.Name == branch {
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

// postDeleteRefresh refreshes the branch list after a successful deletion,
// updates m.current, repositions the cursor, and returns to ListMode.
func postDeleteRefresh(m model, deletedIdx int, wasCurrentBranch bool) model {
	branches, err := listBranches()
	if err != nil {
		m.err = err
		m.mode = tui.ListMode
		return m
	}
	m.branches = branches
	m.firstWorktreeIdx = findFirstWorktreeIdx(branches)
	for i, b := range branches {
		if b.IsCurrent {
			m.current = i
			break
		}
	}
	if wasCurrentBranch || deletedIdx == 0 {
		m.cursor = 0
	} else {
		m.cursor = deletedIdx - 1
	}
	m.mode = tui.ListMode
	m.err = nil
	return m
}

func (m model) updateDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "y", "Y":
			deletedIdx := m.cursor
			wasCurrentBranch := m.cursor == m.current
			if wasCurrentBranch {
				switchTo := ""
				for _, b := range m.branches {
					if b.Name != m.branches[m.cursor].Name {
						switchTo = b.Name
						break
					}
				}
				if err := checkoutBranch(switchTo); err != nil {
					m.err = err
					m.mode = tui.ListMode
					return m, nil
				}
				for i, b := range m.branches {
					if b.Name == switchTo {
						m.current = i
						break
					}
				}
			}
			if err := deleteBranch(m.branches[m.cursor].Name); err != nil {
				if isUnmergedBranchError(err) {
					m.mode = tui.ForceDeleteConfirmMode
					m.err = nil
					return m, nil
				}
				m.err = err
				m.mode = tui.ListMode
				return m, nil
			}
			m = postDeleteRefresh(m, deletedIdx, wasCurrentBranch)
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
			deletedIdx := m.cursor
			wasCurrentBranch := m.cursor == m.current
			if err := forceDeleteBranch(m.branches[m.cursor].Name); err != nil {
				m.err = err
				m.mode = tui.ListMode
				return m, nil
			}
			m = postDeleteRefresh(m, deletedIdx, wasCurrentBranch)
		default:
			m.mode = tui.ListMode
		}
	}
	return m, nil
}
