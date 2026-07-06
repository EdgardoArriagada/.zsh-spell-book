package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"example.com/workspace/lib/tui"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

type Ticket struct {
	Parent  string
	Current string
	Label   string
}

type model struct {
	tickets     []Ticket
	filtered    []Ticket
	cursor      int
	mode        tui.Mode
	searchInput textinput.Model
	width       int
	vp          tui.Viewport
	selected    *Ticket
	err         error
	current     int // index of current ticket in tickets, -1 if none
}

func loadTickets() ([]Ticket, error) {
	home := os.Getenv("HOME")
	f, err := os.Open(home + "/temp/tickets")
	if err != nil {
		return nil, fmt.Errorf("no tickets file at %s/temp/tickets", home)
	}
	defer f.Close()

	var tickets []Ticket
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "|", 3)
		if len(parts) < 3 {
			continue
		}
		tickets = append(tickets, Ticket{
			Parent:  strings.TrimSpace(parts[0]),
			Current: strings.TrimSpace(parts[1]),
			Label:   strings.TrimSpace(parts[2]),
		})
	}
	if len(tickets) == 0 {
		return nil, fmt.Errorf("tickets file is empty")
	}
	return tickets, nil
}

func loadCurrentTicket() string {
	home := os.Getenv("HOME")
	data, err := os.ReadFile(home + "/temp/current-ticket.zsh")
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.Contains(line, "ZSB_CURRENT_TICKET=") {
			continue
		}
		start := strings.Index(line, "'")
		end := strings.LastIndex(line, "'")
		if start >= 0 && end > start {
			return line[start+1 : end]
		}
	}
	return ""
}

func applyTicketFilter(tickets []Ticket, term string) []Ticket {
	return tui.ApplyFilter(tickets, term, func(t Ticket) string {
		return t.Parent + " " + t.Current + " " + t.Label
	})
}

func initialModel() model {
	tickets, err := loadTickets()
	si := tui.NewSearchInput()

	cur := -1
	currentTicket := loadCurrentTicket()
	if currentTicket != "" {
		for i, t := range tickets {
			if t.Current == currentTicket {
				cur = i
				break
			}
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
