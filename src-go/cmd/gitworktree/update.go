package main

import (
	"os"
	"os/exec"
	"strings"

	gitlib "example.com/workspace/lib/git"
	"example.com/workspace/lib/tui"

	tea "github.com/charmbracelet/bubbletea"
)


func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.windowHeight = ws.Height
		m.input.SetWidth(ws.Width)
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availableRows())
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
			m.cursor = (m.cursor + 1) % len(m.filtered)
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availableRows())
	case "k", "up", "shift+tab":
		m.statusMsg = ""
		if len(m.filtered) > 0 {
			m.cursor = (m.cursor - 1 + len(m.filtered)) % len(m.filtered)
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availableRows())
	case "g":
		m.statusMsg = ""
		m.cursor = 0
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availableRows())
	case "G":
		m.statusMsg = ""
		if len(m.filtered) > 0 {
			m.cursor = len(m.filtered) - 1
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availableRows())
	case "enter":
		if len(m.filtered) > 0 {
			m.selected = m.filtered[m.cursor].Path
		}
		return m, tea.Quit
	case "a":
		m.mode = tui.AddMode
		m.input.SetValue("")
		m.err = nil
		m.statusMsg = ""
		return m, m.input.Focus()
	case "d", "D":
		if len(m.filtered) == 0 {
			break
		}
		if len(m.worktrees) > 0 && m.filtered[m.cursor].Path == m.worktrees[0].Path {
			m.statusMsg = "cannot delete main worktree"
		} else {
			m.statusMsg = ""
			m.deleteBranch = km.String() == "D"
			m.mode = tui.DeleteConfirmMode
		}
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
			m.filtered = m.worktrees
			m.cursor = 0
			m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availableRows())
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
	m.filtered = applyWorktreeFilter(m.worktrees, m.searchInput.Value())
	m.cursor = 0
	m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availableRows())
	return m, cmd
}

func (m model) updateAdd(msg tea.Msg) (tea.Model, tea.Cmd) {
	if edMsg, ok := msg.(tui.EditorDoneMsg); ok {
		if content, _ := os.ReadFile(edMsg.TmpFile); len(content) > 0 {
			if branch := strings.TrimSpace(string(content)); branch != "" {
				m.input.SetValue(branch)
			}
		}
		os.Remove(edMsg.TmpFile)
		return m, m.input.Focus()
	}
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.mode = tui.ListMode
			m.input.Blur()
			m.err = nil
			return m, nil
		case "ctrl+g":
			f, err := os.CreateTemp("", "zsb-branch-*")
			if err != nil {
				m.err = err
				return m, nil
			}
			if v := tui.ParseInputValue(m.input); v != "" {
				f.WriteString(v)
			}
			f.Close()
			editor := os.Getenv("EDITOR")
			if editor == "" {
				editor = "vi"
			}
			parts := strings.Fields(editor)
			cmd := exec.Command(parts[0], append(parts[1:], f.Name())...)
			return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
				return tui.EditorDoneMsg{TmpFile: f.Name(), Err: err}
			})
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
			m.filtered = applyWorktreeFilter(wts, m.searchInput.Value())
			m.mode = tui.ListMode
			m.input.Blur()
			m.err = nil
			for i, wt := range m.filtered {
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
			if len(m.filtered) == 0 {
				m.mode = tui.ListMode
				return m, nil
			}
			target := m.filtered[m.cursor]
			branch := target.Branch
			deletingCurrent := m.current >= 0 && target.Path == m.worktrees[m.current].Path
			if err := deleteWorktree(target.Path); err != nil {
				if isWorktreeDirtyError(err) {
					m.mode = tui.ForceDeleteConfirmMode
					m.err = nil
					return m, nil
				}
				m.err = err
				m.mode = tui.ListMode
				m.deleteBranch = false
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
				m.deleteBranch = false
				return m, nil
			}
			m.worktrees = wts
			m.current = currentWorktreeIndex(wts)
			m.filtered = applyWorktreeFilter(wts, m.searchInput.Value())
			if m.cursor >= len(m.filtered) {
				m.cursor = max(0, len(m.filtered)-1)
			}
			m.mode = tui.ListMode
			m.err = nil
			if m.deleteBranch && branch != "" {
				m.err = deleteWorktreeBranch(branch)
			}
			m.deleteBranch = false
		default:
			m.deleteBranch = false
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
			target := m.filtered[m.cursor]
			branch := target.Branch
			deletingCurrent := m.current >= 0 && target.Path == m.worktrees[m.current].Path
			if err := deleteWorktreeForce(target.Path); err != nil {
				m.err = err
				m.mode = tui.ListMode
				m.deleteBranch = false
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
				m.deleteBranch = false
				return m, nil
			}
			m.worktrees = wts
			m.current = currentWorktreeIndex(wts)
			m.filtered = applyWorktreeFilter(wts, m.searchInput.Value())
			if m.cursor >= len(m.filtered) {
				m.cursor = max(0, len(m.filtered)-1)
			}
			m.mode = tui.ListMode
			m.err = nil
			if m.deleteBranch && branch != "" {
				m.err = deleteWorktreeBranch(branch)
			}
			m.deleteBranch = false
		default:
			m.deleteBranch = false
			m.mode = tui.ListMode
		}
	}
	return m, nil
}

