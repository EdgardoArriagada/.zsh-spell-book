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

func TestClampTicketViewportLargeJumpKeepsLastPageFull(t *testing.T) {
	tickets := make([]jira.Ticket, 100)
	if got := clampTicketViewport(tui.Viewport{}, len(tickets)-1, tickets, 20, false).Offset; got != 81 {
		t.Fatalf("offset = %d, want 81", got)
	}
}

func TestRenderUsesBoundedRowsBeforeWindowSize(t *testing.T) {
	tickets := make([]jira.Ticket, defaultAvailableRows+1)
	for i := range tickets {
		tickets[i] = jira.Ticket{Current: "JIRA", Label: "ticket"}
	}
	m := model{tickets: tickets, filtered: tickets, current: -1, width: 80, searchInput: tui.NewSearchInput()}
	if got := strings.Count(m.render(), "JIRA:"); got != defaultAvailableRows {
		t.Fatalf("rendered %d tickets, want %d", got, defaultAvailableRows)
	}
}

func TestRenderShowsScrollIndicatorsAndUsesViewport(t *testing.T) {
	tickets := make([]jira.Ticket, 5)
	for i := range tickets {
		tickets[i] = jira.Ticket{Current: "JIRA", Label: "ticket"}
	}
	m := model{tickets: tickets, filtered: tickets, current: -1, width: 80, availRows: 3, searchInput: tui.NewSearchInput()}
	if got := m.render(); strings.Count(got, "JIRA:") != 2 || !strings.Contains(got, "") || !strings.HasSuffix(got, m.footerSection()) {
		t.Fatalf("first page = %q", got)
	}
	m.vp.Offset = 1
	if got := m.render(); strings.Count(got, "JIRA:") != 1 || !strings.Contains(got, "") || !strings.Contains(got, "") {
		t.Fatalf("middle page = %q", got)
	}
}

func TestTruncateLabelPreservesShortUnicode(t *testing.T) {
	if got := truncateLabel("café", 4); got != "café" {
		t.Fatalf("label = %q", got)
	}
}

func TestEditTicketsCmdKeepsOtherEditors(t *testing.T) {
	cmd := editTicketsCmd("vim", "/tmp/tickets", 4)
	if got, want := cmd.Args, []string{"vim", "/tmp/tickets"}; !slices.Equal(got, want) {
		t.Fatalf("args = %q, want %q", got, want)
	}
}
