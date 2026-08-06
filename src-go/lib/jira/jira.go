package jira

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// AppendTicketRow appends parentKey|issueKey|summary to ticketsPath.
// Returns (existingLine, nil) if already present, (0, nil) on success, (0, err) on write failure.
func AppendTicketRow(ticketsPath, parentKey, issueKey, summary string) (int, error) {
	if lnum := findTicketLine(ticketsPath, issueKey); lnum > 0 {
		return lnum, nil
	}
	f, err := os.OpenFile(ticketsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("open tickets: %w", err)
	}
	defer f.Close()
	_, err = fmt.Fprintf(f, "%s|%s|%s\n", parentKey, issueKey, summary)
	return 0, err
}

func findTicketLine(path, issueKey string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	needle := "|" + issueKey + "|"
	sc := bufio.NewScanner(f)
	for lnum := 1; sc.Scan(); lnum++ {
		if strings.Contains(sc.Text(), needle) {
			return lnum
		}
	}
	_ = sc.Err() // ponytail: find-only; 0 = not found, error indistinguishable from miss
	return 0
}

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
