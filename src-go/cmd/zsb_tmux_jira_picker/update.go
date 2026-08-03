package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"example.com/workspace/lib/tui"

	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.windowHeight = ws.Height
		m.availRows = tui.AvailableRows(m.windowHeight, tui.Title("Jira Tickets"), m.statusSection(), "\n"+m.footerSection()+"\n")
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availRows)
		return m, nil
	}
	if nc, ok := msg.(notifCountsMsg); ok {
		m.notifCounts = nc
		return m, tickNotifCountsCmd()
	}
	if m.mode == tui.SearchMode {
		return m.updateSearch(msg)
	}
	return m.updateList(msg)
}

func (m model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	if _, ok := msg.(tui.EditorDoneMsg); ok {
		m = m.reloadTickets()
		if m.err != nil {
			return m, nil
		}
		m.filtered = m.filterTickets(m.searchInput.Value())
		if m.cursor >= len(m.filtered) {
			m.cursor = max(0, len(m.filtered)-1)
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availRows)
		return m, nil
	}

	km, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	switch km.String() {
	case "q", "ctrl+c", "esc":
		return m, tea.Quit
	case "/":
		m.mode = tui.SearchMode
		return m, m.searchInput.Focus()
	case "j", "down", "tab":
		if len(m.filtered) > 0 {
			m.cursor = (m.cursor + 1) % len(m.filtered)
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availRows)
	case "k", "up", "shift+tab":
		if len(m.filtered) > 0 {
			m.cursor = (m.cursor - 1 + len(m.filtered)) % len(m.filtered)
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availRows)
	case "g":
		m.cursor = 0
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availRows)
	case "G":
		if len(m.filtered) > 0 {
			m.cursor = len(m.filtered) - 1
		}
		m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availRows)
	case "enter":
		if len(m.filtered) > 0 {
			t := m.filtered[m.cursor]
			m.selected = &t
		}
		return m, tea.Quit
	case "ctrl+g":
		if len(m.filtered) == 0 {
			return m, nil
		}
		home := os.Getenv("HOME")
		editor := os.Getenv("EDITOR")
		cmd := editTicketsCmd(editor, home+"/temp/tickets", m.filtered[m.cursor].Line)
		return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
			return tui.EditorDoneMsg{Err: err}
		})
	}
	return m, nil
}

func editTicketsCmd(editor, filename string, line int) *exec.Cmd {
	parts := strings.Fields(editor)
	if len(parts) == 0 {
		parts = []string{"vi"}
	}
	args := parts[1:]
	if filepath.Base(parts[0]) == "nvim" {
		args = append(args, "+"+strconv.Itoa(line))
	}
	return exec.Command(parts[0], append(args, filename)...)
}

func (m model) updateSearch(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyPressMsg); ok {
		switch km.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.searchInput.Blur()
			m.searchInput.SetValue("")
			m.filtered = m.tickets
			m.cursor = 0
			m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availRows)
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
	m.filtered = m.filterTickets(m.searchInput.Value())
	m.cursor = 0
	m.vp = m.vp.Clamp(m.cursor, len(m.filtered), m.availRows)
	return m, cmd
}
