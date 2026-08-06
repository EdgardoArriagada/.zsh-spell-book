package app

import (
	"os"
	"time"

	"example.com/workspace/lib/jira"
	"example.com/workspace/lib/tui"

	"charm.land/bubbles/v2/textinput"
	tea "charm.land/bubbletea/v2"
)

type model struct {
	tickets      []jira.Ticket
	filterKeys   []string // search projections, one per ticket, computed at load
	filtered     []jira.Ticket
	cursor       int
	mode         tui.Mode
	searchInput  textinput.Model
	width        int
	windowHeight int
	availRows    int // cached tui.AvailableRows result; updated on WindowSizeMsg only
	vp           tui.Viewport
	selected     *jira.Ticket
	err          error
	current      int // index of current ticket in tickets, -1 if none
	notifCounts  map[string]jira.NotifCounts
	statusMsg    string
	blinkOn      bool
}

const defaultAvailableRows = 20

// notifCountsMsg delivers tmux notification counts loaded off the startup
// critical path, so the list paints before the tmux subprocess returns.
type notifCountsMsg map[string]jira.NotifCounts

func loadNotifCountsCmd() tea.Msg {
	return notifCountsMsg(jira.LoadNotificationCounts())
}

func tickNotifCountsCmd() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		return notifCountsMsg(jira.LoadNotificationCounts())
	})
}

type blinkMsg struct{}

func blinkCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(time.Time) tea.Msg { return blinkMsg{} })
}

func (m model) filterTickets(term string) []jira.Ticket {
	return tui.ApplyFilterKeys(m.tickets, m.filterKeys, term)
}

func (m model) showTitles() bool {
	return m.mode != tui.SearchMode && m.searchInput.Value() == ""
}

func (m model) clampViewport() model {
	m.vp = clampTicketViewport(m.vp, m.cursor, m.filtered, m.availRows, m.showTitles())
	return m
}

func (m model) reloadTickets() model {
	tickets, err := jira.LoadTickets()
	m.tickets = tickets
	m.err = err
	m.filterKeys = make([]string, len(tickets))
	m.current = -1
	currentTicket := os.Getenv("ZSB_CURRENT_TICKET")
	for i, t := range tickets {
		m.filterKeys[i] = t.Parent + " " + t.Current + " " + t.Label
		if m.current == -1 && t.Current == currentTicket {
			m.current = i
		}
	}
	return m
}

func initialModel() model {
	m := model{
		searchInput: tui.NewSearchInput(),
		width:       tui.DefaultWidth,
		availRows:   defaultAvailableRows,
		current:     -1,
	}
	m = m.reloadTickets()
	m.cursor = max(0, m.current)
	m.filtered = m.tickets
	return m
}

func (m model) Init() tea.Cmd { return tea.Batch(loadNotifCountsCmd, blinkCmd()) }
