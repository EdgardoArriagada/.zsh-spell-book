package main

import (
	"fmt"
	"os/exec"
	"strings"

	"example.com/workspace/lib/jira"
)

func main() {
	tickets, _ := jira.LoadTickets()

	currentIDs := make(map[string]bool, len(tickets))
	for _, t := range tickets {
		currentIDs[t.Current] = true
	}

	prefixes := jira.TmuxSessionPrefixes(currentIDs)

	out, err := exec.Command("tmux", "ls", "-F", "#{session_name}").Output()
	if err != nil {
		return
	}

	for s := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		if s == "" {
			continue
		}
		if !jira.MatchesTmuxSession(s, prefixes) {
			fmt.Println(s)
		}
	}
}
