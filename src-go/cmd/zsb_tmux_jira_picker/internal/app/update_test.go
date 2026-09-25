package app

import (
	"slices"
	"strings"
	"testing"

	"example.com/workspace/lib/jira"
	"example.com/workspace/lib/tui"

	tea "charm.land/bubbletea/v2"
)

func TestEditTicketsCmdStartsNvimAtTicketLine(t *testing.T) {
	cmd := editTicketsCmd("nvim --clean", "/tmp/tickets", 4)
	if got, want := cmd.Args, []string{"nvim", "--clean", "+let g:zsb_prevent_renametab = 1", "+4", "/tmp/tickets"}; !slices.Equal(got, want) {
		t.Fatalf("args = %q, want %q", got, want)
	}
}

func TestPageKeysMoveTicketCursor(t *testing.T) {
	tickets := make([]jira.Ticket, 10)
	m := model{tickets: tickets, filtered: tickets, availRows: 4, current: -1, width: 80, searchInput: tui.NewSearchInput()}

	updated, _ := m.Update(tea.KeyPressMsg{Code: 'd', Mod: tea.ModCtrl})
	m = updated.(model)
	if m.cursor != 4 {
		t.Fatalf("ctrl+d cursor = %d, want 4", m.cursor)
	}
	updated, _ = m.Update(tea.KeyPressMsg{Code: 'u', Mod: tea.ModCtrl})
	if got := updated.(model).cursor; got != 0 {
		t.Fatalf("ctrl+u cursor = %d, want 0", got)
	}
}

func TestNumberFocusesOnlyVisibleTicketRowsUntilEnter(t *testing.T) {
	tickets := make([]jira.Ticket, 6)
	for i := range tickets {
		tickets[i] = jira.Ticket{Current: "JIRA-" + string(rune('0'+i)), Title: "X"}
	}
	tickets[3].Title = "Y"
	m := model{tickets: tickets, filtered: tickets, availRows: 6, current: -1, searchInput: tui.NewSearchInput()}
	m.vp.Offset = 2
	updated, cmd := m.Update(tea.KeyPressMsg{Text: "3"})
	m = updated.(model)
	if cmd != nil || m.selected != nil {
		t.Fatalf("3 selected hidden ticket: cmd=%v selected=%v", cmd, m.selected)
	}
	updated, cmd = m.Update(tea.KeyPressMsg{Text: "2"})
	m = updated.(model)
	if cmd != nil || m.cursor != 3 || m.selected != nil {
		t.Fatalf("2: cmd=%v cursor=%d selected=%v", cmd, m.cursor, m.selected)
	}
	updated, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(model)
	if cmd == nil || m.selected == nil || m.selected.Current != tickets[3].Current {
		t.Fatalf("enter selected %v, want %q", m.selected, tickets[3].Current)
	}
}

func TestYankSessionCmdGeneratesIDForTicketUnderCursor(t *testing.T) {
	m := model{
		filtered: []jira.Ticket{
			{Current: "JIRA-1", Label: "Other ticket"},
			{Current: "JIRA-2", Label: "Selected ticket"},
		},
		cursor: 1,
	}
	var copied string
	msg := m.yankSessionCmd(func(value string) error {
		copied = value
		return nil
	})()
	if got := msg.(yankSessionMsg).err; got != nil {
		t.Fatal(got)
	}
	if copied != "JIRA-2-selected-ticket" {
		t.Fatalf("copied %q", copied)
	}
}

func TestQuitKeysClearTicketSearchAndPreserveSelection(t *testing.T) {
	for _, key := range []tea.KeyPressMsg{{Text: "q"}, {Code: 'c', Mod: tea.ModCtrl}, {Code: tea.KeyEscape}} {
		tickets := []jira.Ticket{{Current: "JIRA-1"}, {Current: "JIRA-2"}, {Current: "JIRA-3"}}
		m := model{tickets: tickets, filtered: tickets[1:], cursor: 1, availRows: 4, searchInput: tui.NewSearchInput()}
		m.searchInput.SetValue("JIRA")

		updated, cmd := m.Update(key)
		m = updated.(model)
		if cmd != nil || m.searchInput.Value() != "" || len(m.filtered) != len(tickets) || m.cursor != 2 {
			t.Fatalf("%s did not clear search with selection preserved: cmd=%v term=%q len=%d cursor=%d", key.String(), cmd, m.searchInput.Value(), len(m.filtered), m.cursor)
		}
		if _, cmd = m.Update(key); cmd == nil {
			t.Fatalf("second %s should quit", key.String())
		}
	}
}

func TestSearchKeepsSelectionUntilEnter(t *testing.T) {
	tickets := []jira.Ticket{
		{Current: "JIRA-1", Label: "main"},
		{Current: "JIRA-2", Label: "feature"},
		{Current: "JIRA-3", Label: "fix"},
	}
	filterKeys := []string{"main", "feature", "fix"}
	press := func(m model, key string) model {
		updated, _ := m.Update(tea.KeyPressMsg{Text: key})
		return updated.(model)
	}
	newModel := func(cursor int) model {
		return model{
			tickets: tickets, filterKeys: filterKeys, filtered: tickets, cursor: cursor,
			current: -1, width: 80, availRows: 10, searchInput: tui.NewSearchInput(),
		}
	}

	m := press(press(newModel(0), "/"), "f")
	if len(m.filtered) != 2 || strings.Contains(m.render(), "▸") {
		t.Fatalf("search selected a match: len=%d cursorVisible=%t", len(m.filtered), strings.Contains(m.render(), "▸"))
	}
	updated, cmd := m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(model)
	if cmd != nil || m.mode != tui.ListMode || m.cursor != 0 || m.selected != nil {
		t.Fatalf("first enter selected instead of focusing first match: cmd=%v mode=%v cursor=%d selected=%v", cmd, m.mode, m.cursor, m.selected)
	}
	first := m.filtered[0].Current
	updated, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(model)
	if cmd == nil || m.selected == nil || m.selected.Current != first {
		t.Fatalf("second enter: cmd=%v selected=%v, want %q", cmd, m.selected, first)
	}

	m = press(press(newModel(2), "/"), "f")
	if m.filtered[m.cursor].Current != tickets[2].Current || !strings.Contains(m.render(), "▸") {
		t.Fatal("search did not preserve a matching selection")
	}

	m = press(press(newModel(0), "/"), "x")
	updated, cmd = m.Update(tea.KeyPressMsg{Code: tea.KeyEnter})
	m = updated.(model)
	if cmd == nil || m.selected == nil || m.selected.Current != tickets[2].Current {
		t.Fatalf("single match enter: cmd=%v selected=%v, want %q", cmd, m.selected, tickets[2].Current)
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
	if end, rows := visibleTicketEnd(tickets, 0, 3, true); end != 1 || rows != 3 {
		t.Fatalf("end, rows = %d, %d; want 1, 3", end, rows)
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
	m := model{tickets: tickets, filtered: tickets, current: -1, width: 80, availRows: 3, blinkOn: true, searchInput: tui.NewSearchInput()}
	if got := m.render(); strings.Count(got, "JIRA:") != 2 || !strings.Contains(got, "") || !strings.HasSuffix(got, m.footerSection()) {
		t.Fatalf("first page = %q", got)
	}
	m.vp.Offset = 1
	if got := m.render(); strings.Count(got, "JIRA:") != 1 || !strings.Contains(got, "") || !strings.Contains(got, "") {
		t.Fatalf("middle page = %q", got)
	}
}

func TestFooterShowsTotalAgentCountsOutsideSearch(t *testing.T) {
	m := model{notifCounts: map[string]jira.NotifCounts{
		"one": {Working: 2, Finished: 1},
		"two": {Working: 3, Finished: 4},
	}}
	if got := m.footerSection(); !strings.Contains(got, "󰔟") || !strings.Contains(got, "󰂚") || strings.Count(got, " 5") != 2 {
		t.Fatalf("footer = %q", got)
	}
	m.mode = tui.SearchMode
	if got := m.footerSection(); strings.Contains(got, "󰔟") || strings.Contains(got, "󰂚") {
		t.Fatalf("search footer = %q", got)
	}
}

func TestNotificationRefreshSchedulesNextPoll(t *testing.T) {
	updated, cmd := (model{}).Update(notifCountsMsg{"session": {Finished: 1}})
	if cmd == nil {
		t.Fatal("notification update did not schedule next poll")
	}
	if got := updated.(model).notifCounts["session"].Finished; got != 1 {
		t.Fatalf("finished notifications = %d, want 1", got)
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
