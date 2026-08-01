package jira

import (
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const maxTmuxIDLen = 38

var (
	invalidChars      = regexp.MustCompile(`[^a-zA-Z0-9-]`)
	consecutiveDashes = regexp.MustCompile(`-{2,}`)
)

func (t Ticket) TmuxSessionID() string {
	kebab := strings.ToLower(strings.ReplaceAll(t.Label, " ", "-"))
	result := SanitizeTmuxID(t.Current + "-" + kebab)
	if len(result) > maxTmuxIDLen {
		result = result[:maxTmuxIDLen]
	}
	return strings.TrimRight(result, "-")
}

func TmuxSessionPrefixes(currentIDs map[string]bool) []string {
	prefixes := make([]string, 0, len(currentIDs))
	for id := range currentIDs {
		prefixes = append(prefixes, SanitizeTmuxID(id)+"-")
	}
	return prefixes
}

func MatchesTmuxSession(session string, prefixes []string) bool {
	for _, p := range prefixes {
		if strings.HasPrefix(session, p) {
			return true
		}
	}
	return false
}

func SanitizeTmuxID(s string) string {
	s = invalidChars.ReplaceAllString(s, "-")
	s = consecutiveDashes.ReplaceAllString(s, "-")
	return s
}

func LoadNotificationCounts() map[string]int {
	counts := map[string]int{}
	out, err := exec.Command("tmux", "list-panes", "-a", "-F",
		"#{session_name}\t#{@zsb_agent_notif}").Output()
	if err != nil {
		return counts
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) < 2 {
			continue
		}
		n, _ := strconv.Atoi(strings.TrimSpace(parts[1])) // empty/unset → 0
		counts[parts[0]] += n
	}
	return counts
}
