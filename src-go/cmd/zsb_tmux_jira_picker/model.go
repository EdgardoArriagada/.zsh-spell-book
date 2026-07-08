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
	width       int
	vp          tui.Viewport
	selected    *jira.Ticket
	err         error
	current     int // index of current ticket in tickets, -1 if none
}

func applyTicketFilter(tickets []jira.Ticket, term string) []jira.Ticket {
	return tui.ApplyFilter(tickets, term, func(t jira.Ticket) string {
		return t.Parent + " " + t.Current + " " + t.Label
	})
}

func initialModel() model {
	tickets, err := jira.LoadTickets()
	si := tui.NewSearchInput()

	cur := -1
	currentTicket := os.Getenv("ZSB_CURRENT_TICKET")
	for i, t := range tickets {
		if t.Current == currentTicket {
			cur = i
			break
		}
	}

	cursor := 0
	if cur > 0 {
		cursor = cur
	}

	return model{
		tickets:     tickets,
		filtered:    tickets,
		cursor:      cursor,
		searchInput: si,
		current:     cur,
		width:       tui.DefaultWidth,
		err:         err,
	}
}

func (m model) Init() tea.Cmd { return nil }
