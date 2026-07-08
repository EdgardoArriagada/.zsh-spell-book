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

	out, err := exec.Command("tmux", "ls", "-F", "#{session_name}").Output()
	if err != nil {
		return
	}

	for s := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		if s == "" {
			continue
		}
		if !isJiraSession(s, currentIDs) {
			fmt.Println(s)
		}
	}
}

func isJiraSession(session string, currentIDs map[string]bool) bool {
	for id := range currentIDs {
		if strings.HasPrefix(session, id+"-") {
			return true
		}
	}
	return false
}
