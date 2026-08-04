package jira

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Ticket struct {
	Title   string
	Parent  string
	Current string
	Label   string
	Line    int
	// SessionID is the cached tmux session id, computed once at load so the
	// render path (per row, per frame) does a plain field read instead of
	// re-running SanitizeTmuxID.
	SessionID string
}

func LoadTickets() ([]Ticket, error) {
	home := os.Getenv("HOME")
	f, err := os.Open(home + "/temp/tickets")
	if err != nil {
		return nil, fmt.Errorf("no tickets file at %s/temp/tickets", home)
	}
	defer f.Close()

	var tickets []Ticket
	var title string
	sc := bufio.NewScanner(f)
	for lineNumber := 1; sc.Scan(); lineNumber++ {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parent, rest, ok := strings.Cut(line, "|")
		if !ok {
			title = line
			continue
		}
		current, label, ok := strings.Cut(rest, "|")
		if !ok {
			continue
		}
		t := Ticket{
			Title:   title,
			Parent:  strings.TrimSpace(parent),
			Current: strings.TrimSpace(current),
			Label:   strings.TrimSpace(label),
			Line:    lineNumber,
		}
		t.SessionID = t.TmuxSessionID()
		tickets = append(tickets, t)
	}
	if err := sc.Err(); err != nil {
		return nil, err
	}
	if len(tickets) == 0 {
		return nil, fmt.Errorf("tickets file is empty")
	}
	return tickets, nil
}
