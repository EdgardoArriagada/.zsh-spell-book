package app

import (
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"example.com/workspace/lib/open"
	"example.com/workspace/lib/tui"

	tea "charm.land/bubbletea/v2"
)

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if ws, ok := msg.(tea.WindowSizeMsg); ok {
		m.width = ws.Width
		m.windowHeight = ws.Height
		m.availRows = tui.AvailableRows(m.windowHeight, m.statusSection(), m.footerSection())
		m = m.clampViewport()
		return m, nil
	}
	if nc, ok := msg.(notifCountsMsg); ok {
		m.notifCounts = nc
		return m, tickNotifCountsCmd()
	}
	if _, ok := msg.(blinkMsg); ok {
		m.blinkOn = !m.blinkOn
		return m, blinkCmd()
	}
	if ob, ok := msg.(openBrowserMsg); ok {
		if ob.err != nil {
			m.statusMsg = "open browser: " + ob.err.Error()
		}
		return m, nil
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
		m = m.clampViewport()
		return m, nil
	}

	km, ok := msg.(tea.KeyPressMsg)
	if !ok {
		return m, nil
	}
	m.statusMsg = ""
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
		return m.clampViewport(), nil
	case "k", "up", "shift+tab":
		if len(m.filtered) > 0 {
			m.cursor = (m.cursor - 1 + len(m.filtered)) % len(m.filtered)
		}
		return m.clampViewport(), nil
	case "ctrl+d":
		m.cursor = tui.PageCursor(m.cursor, len(m.filtered), m.availRows, 1)
		return m.clampViewport(), nil
	case "ctrl+u":
		m.cursor = tui.PageCursor(m.cursor, len(m.filtered), m.availRows, -1)
		return m.clampViewport(), nil
	case "g":
		m.cursor = 0
		return m.clampViewport(), nil
	case "G":
		if len(m.filtered) > 0 {
			m.cursor = len(m.filtered) - 1
		}
		return m.clampViewport(), nil
	case "H":
		if len(m.filtered) > 0 {
			m.cursor = m.vp.Offset
		}
		return m.clampViewport(), nil
	case "M":
		if len(m.filtered) > 0 {
			end, _, _, _ := ticketViewport(m.filtered, m.vp.Offset, m.availRows, m.showTitles())
			m.cursor = m.vp.Offset + (end-1-m.vp.Offset)/2
		}
		return m.clampViewport(), nil
	case "L":
		if len(m.filtered) > 0 {
			end, _, _, _ := ticketViewport(m.filtered, m.vp.Offset, m.availRows, m.showTitles())
			m.cursor = end - 1
		}
		return m.clampViewport(), nil
	case "enter":
		if len(m.filtered) > 0 {
			t := m.filtered[m.cursor]
			m.selected = &t
		}
		return m, tea.Quit
	case "O":
		return m.handleOpenBrowser()
	case "alt+j":
		return m.handleOpenNotes()
	case "ctrl+g":
		return m.handleEditTickets()
	}
	return m, nil
}

func (m model) handleOpenBrowser() (tea.Model, tea.Cmd) {
	if len(m.filtered) == 0 {
		return m, nil
	}
	baseURL := strings.TrimRight(os.Getenv("ZSB_JIRA_BASEURL"), "/")
	if baseURL == "" {
		m.statusMsg = "ZSB_JIRA_BASEURL not set"
		return m, nil
	}
	return m, openBrowserCmd(baseURL + "/browse/" + m.filtered[m.cursor].Current)
}

func (m model) handleOpenNotes() (tea.Model, tea.Cmd) {
	if len(m.filtered) == 0 {
		return m, nil
	}
	ticketsDir := os.Getenv("ZSB_TICKETS_DIR")
	if ticketsDir == "" {
		m.statusMsg = "ZSB_TICKETS_DIR not set"
		return m, nil
	}
	t := m.filtered[m.cursor]
	notesPath := filepath.Join(ticketsDir, t.Parent, t.Current, "NOTES.md")
	if _, err := os.Stat(notesPath); os.IsNotExist(err) {
		m.statusMsg = "Notes not found"
		return m, nil
	}
	cmd := openNotesCmd(os.Getenv("EDITOR"), notesPath)
	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		return tui.EditorDoneMsg{Err: err}
	})
}

func (m model) handleEditTickets() (tea.Model, tea.Cmd) {
	if len(m.filtered) == 0 {
		return m, nil
	}
	cmd := editTicketsCmd(os.Getenv("EDITOR"), os.Getenv("HOME")+"/temp/tickets", m.filtered[m.cursor].Line)
	return m, tea.ExecProcess(cmd, func(err error) tea.Msg {
		return tui.EditorDoneMsg{Err: err}
	})
}

type openBrowserMsg struct{ err error }

func openBrowserCmd(url string) tea.Cmd {
	return func() tea.Msg { return openBrowserMsg{err: open.Url(url)} }
}

func openNotesCmd(editor, notesPath string) *exec.Cmd {
	parts := strings.Fields(editor)
	if len(parts) == 0 {
		parts = []string{"vi"}
	}
	args := parts[1:]
	if filepath.Base(parts[0]) == "nvim" {
		args = append(args, "+let g:zsb_prevent_renametab = 1", "-c", "cd "+filepath.Dir(notesPath))
	}
	return exec.Command(parts[0], append(args, notesPath)...)
}

func editTicketsCmd(editor, filename string, line int) *exec.Cmd {
	parts := strings.Fields(editor)
	if len(parts) == 0 {
		parts = []string{"vi"}
	}
	args := parts[1:]
	if filepath.Base(parts[0]) == "nvim" {
		args = append(args, "+let g:zsb_prevent_renametab = 1", "+"+strconv.Itoa(line))
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
			m = m.clampViewport()
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
	m = m.clampViewport()
	return m, cmd
}
