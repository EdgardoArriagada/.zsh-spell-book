package jira

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Ticket struct {
	Parent  string
	Current string
	Label   string
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
		t := Ticket{
			Parent:  strings.TrimSpace(parts[0]),
			Current: strings.TrimSpace(parts[1]),
			Label:   strings.TrimSpace(parts[2]),
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
