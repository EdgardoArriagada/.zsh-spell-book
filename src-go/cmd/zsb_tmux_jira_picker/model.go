package main

import (
	"os"

	"example.com/workspace/lib/jira"
	"example.com/workspace/lib/tui"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	tickets     []jira.Ticket
	filtered    []jira.Ticket
	cursor      int
	mode        tui.Mode
	searchInput textinput.Model
	width        int
	windowHeight int
	vp           tui.Viewport
	selected    *jira.Ticket
	err         error
	current     int // index of current ticket in tickets, -1 if none
}

func applyTicketFilter(tickets []jira.Ticket, term string) []jira.Ticket {
	return tui.ApplyFilter(tickets, term, func(t jira.Ticket) string {
		return t.Parent + " " + t.Current + " " + t.Label
	})
}

func (m model) reloadTickets() model {
	tickets, err := jira.LoadTickets()
	m.tickets = tickets
	m.err = err
	currentTicket := os.Getenv("ZSB_CURRENT_TICKET")
	m.current = -1
	for i, t := range tickets {
		if t.Current == currentTicket {
			m.current = i
			break
		}
	}
	return m
}

func initialModel() model {
	m := model{
		searchInput: tui.NewSearchInput(),
		width:       tui.DefaultWidth,
		current:     -1,
	}
	m = m.reloadTickets()
	m.cursor = max(0, m.current)
	m.filtered = m.tickets
	return m
}

func (m model) Init() tea.Cmd { return nil }
