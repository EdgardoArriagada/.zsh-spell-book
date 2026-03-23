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
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered))
		return m, nil
	}
	switch m.mode {
	case tui.AddMode:
		return m.updateAdd(msg)
	case tui.DeleteConfirmMode:
		return m.updateDelete(msg)
	case tui.ForceDeleteConfirmMode:
		return m.updateForceDelete(msg)
	case tui.SearchMode:
		return m.updateSearch(msg)
	default:
		return m.updateList(msg)
	}
}

func (m model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	filterActive := m.searchInput.Value() != ""
	switch km.String() {
	case "q", "ctrl+c", "esc":
		return m, tea.Quit
	case "/":
		m.mode = tui.SearchMode
		m.statusMsg = ""
		return m, m.searchInput.Focus()
	case "j", "down", "tab":
		m.statusMsg = ""
		if len(m.filtered) > 0 {
			if filterActive {
				m.cursor = (m.cursor + 1) % len(m.filtered)
			} else {
				next := (m.cursor + 1) % len(m.filtered)
				for next != m.cursor && m.filtered[next].IsWorktree {
					next = (next + 1) % len(m.filtered)
				}
				if !m.filtered[next].IsWorktree {
					m.cursor = next
				}
			}
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered))
	case "k", "up", "shift+tab":
		m.statusMsg = ""
		if len(m.filtered) > 0 {
			if filterActive {
				m.cursor = (m.cursor - 1 + len(m.filtered)) % len(m.filtered)
			} else {
				next := (m.cursor - 1 + len(m.filtered)) % len(m.filtered)
				for next != m.cursor && m.filtered[next].IsWorktree {
					next = (next - 1 + len(m.filtered)) % len(m.filtered)
				}
				if !m.filtered[next].IsWorktree {
					m.cursor = next
				}
			}
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered))
	case "g":
		m.statusMsg = ""
		m.cursor = 0
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered))
	case "G":
		m.statusMsg = ""
		if len(m.filtered) > 0 {
			if filterActive {
				m.cursor = len(m.filtered) - 1
			} else {
				last := len(m.filtered) - 1
				// find last index before worktree section
				lastWorktree := -1
				for i, br := range m.filtered {
					if br.IsWorktree {
						if lastWorktree < 0 {
							lastWorktree = i
						}
					}
				}
				if lastWorktree >= 0 {
					last = lastWorktree - 1
				}
				if last >= 0 {
					m.cursor = last
				}
			}
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered))
	case "enter":
		if len(m.filtered) > 0 && m.filtered[m.cursor].IsWorktree {
			m.statusMsg = "worktree branches are not selectable"
			return m, nil
		}
		if len(m.filtered) > 0 {
			m.selected = m.filtered[m.cursor].Name
		}
		return m, tea.Quit
	case "a":
		m.mode = tui.AddMode
		m.input.SetValue("")
		m.err = nil
		m.statusMsg = ""
		return m, m.input.Focus()
	case "d":
		if len(m.filtered) == 0 {
			break
		}
		br := m.filtered[m.cursor]
		currentName := ""
		if m.current >= 0 && m.current < len(m.branches) {
			currentName = m.branches[m.current].Name
		}
		if br.Name == currentName && len(m.branches) == 1 {
			m.statusMsg = "cannot delete the only branch"
			break
		}
		if isDefaultBranch(br.Name) {
			m.statusMsg = "cannot delete default branch"
			break
		}
		m.statusMsg = ""
		m.mode = tui.DeleteConfirmMode
	}
	return m, nil
}

func (m model) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.searchInput.Blur()
			m.searchInput.SetValue("")
			m.filtered = m.branches
			m.cursor = 0
			m.vp = m.vp.Clamp(m.cursor, len(m.filtered))
			m.mode = tui.ListMode
			return m, nil
		case "enter":
			m.searchInput.Blur()
			m.mode = tui.ListMode
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.searchInput, cmd = m.searchInput.Update(msg)
	m.filtered = applyBranchFilter(m.branches, m.searchInput.Value())
	m.cursor = 0
	m.vp = m.vp.Clamp(m.cursor, len(m.filtered))
	return m, cmd
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
			m.filtered = applyBranchFilter(branches, m.searchInput.Value())
			m.mode = tui.ListMode
			m.input.Blur()
			m.err = nil
			for i, b := range m.filtered {
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
	m.filtered = applyBranchFilter(branches, m.searchInput.Value())
	if wasCurrentBranch || deletedIdx == 0 {
		m.cursor = 0
	} else {
		m.cursor = deletedIdx - 1
	}
	if m.cursor >= len(m.filtered) {
		m.cursor = max(0, len(m.filtered)-1)
	}
	m.mode = tui.ListMode
	m.err = nil
	return m
}

func (m model) updateDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "y", "Y":
			if len(m.filtered) == 0 {
				m.mode = tui.ListMode
				return m, nil
			}
			br := m.filtered[m.cursor]
			currentName := ""
			if m.current >= 0 && m.current < len(m.branches) {
				currentName = m.branches[m.current].Name
			}
			wasCurrentBranch := br.Name == currentName
			if wasCurrentBranch {
				switchTo := ""
				for _, b := range m.branches {
					if b.Name != br.Name {
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
			if err := deleteBranch(br.Name); err != nil {
				if isUnmergedBranchError(err) {
					m.mode = tui.ForceDeleteConfirmMode
					m.err = nil
					return m, nil
				}
				m.err = err
				m.mode = tui.ListMode
				return m, nil
			}
			m = postDeleteRefresh(m, m.cursor, wasCurrentBranch)
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
			if len(m.filtered) == 0 {
				m.mode = tui.ListMode
				return m, nil
			}
			br := m.filtered[m.cursor]
			currentName := ""
			if m.current >= 0 && m.current < len(m.branches) {
				currentName = m.branches[m.current].Name
			}
			wasCurrentBranch := br.Name == currentName
			if err := forceDeleteBranch(br.Name); err != nil {
				m.err = err
				m.mode = tui.ListMode
				return m, nil
			}
			m = postDeleteRefresh(m, m.cursor, wasCurrentBranch)
		default:
			m.mode = tui.ListMode
		}
	}
	return m, nil
}

