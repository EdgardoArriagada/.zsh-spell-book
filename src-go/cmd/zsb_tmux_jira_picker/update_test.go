package main

import (
	"slices"
	"strings"
	"testing"

	"example.com/workspace/lib/jira"
	"example.com/workspace/lib/tui"
)

func TestEditTicketsCmdStartsNvimAtTicketLine(t *testing.T) {
	cmd := editTicketsCmd("nvim --clean", "/tmp/tickets", 4)
	if got, want := cmd.Args, []string{"nvim", "--clean", "+4", "/tmp/tickets"}; !slices.Equal(got, want) {
		t.Fatalf("args = %q, want %q", got, want)
	}
}

func TestRenderSectionTitlesOnlyWithoutSearch(t *testing.T) {
	tickets := []jira.Ticket{
		{Title: "Personal", Current: "JIRA-1", Label: "First"},
		{Title: "Work", Current: "JIRA-2", Label: "Second"},
	}
	m := model{
		tickets: tickets, filtered: tickets, current: -1, width: 80, availRows: 10,
		searchInput: tui.NewSearchInput(),
	}

	if got := m.render(); !strings.Contains(got, "Personal") || !strings.Contains(got, "Work") || strings.Contains(got, "Jira Tickets") {
		t.Fatalf("unexpected unfiltered view: %q", got)
	}
	m.searchInput.SetValue("First")
	m.filtered = tickets[:1]
	if got := m.render(); strings.Contains(got, "Personal") {
		t.Fatalf("search view contains section title: %q", got)
	}
}

func TestSectionTitlesCountAsViewportRows(t *testing.T) {
	tickets := []jira.Ticket{{Title: "Personal"}, {Title: "Work"}}
	if end, rows := visibleTicketEnd(tickets, 0, 4, true); end != 1 || rows != 4 {
		t.Fatalf("end, rows = %d, %d; want 1, 4", end, rows)
	}
	if got := clampTicketViewport(tui.Viewport{}, 1, tickets, 3, true).Offset; got != 1 {
		t.Fatalf("offset = %d, want 1", got)
	}
}

func TestEditTicketsCmdKeepsOtherEditors(t *testing.T) {
	cmd := editTicketsCmd("vim", "/tmp/tickets", 4)
	if got, want := cmd.Args, []string{"vim", "/tmp/tickets"}; !slices.Equal(got, want) {
		t.Fatalf("args = %q, want %q", got, want)
	}
}
