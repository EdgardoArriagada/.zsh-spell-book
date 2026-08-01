package jira

import (
	"os/exec"
	"strconv"
	"strings"
)

const maxTmuxIDLen = 38

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

// SanitizeTmuxID replaces every char outside [a-zA-Z0-9-] with '-' and
// collapses runs of dashes into one, in a single pass (no regex).
func SanitizeTmuxID(s string) string {
	var b strings.Builder
	b.Grow(len(s))
	prevDash := false
	for _, r := range s {
		ok := r == '-' ||
			(r >= 'a' && r <= 'z') ||
			(r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9')
		if !ok {
			r = '-'
		}
		if r == '-' {
			if prevDash {
				continue
			}
			prevDash = true
		} else {
			prevDash = false
		}
		b.WriteByte(byte(r)) // ponytail: all valid output chars are ASCII [a-zA-Z0-9-]
	}
	return b.String()
}

func LoadNotificationCounts() map[string]int {
	counts := map[string]int{}
	out, err := exec.Command("tmux", "list-panes", "-a", "-F",
		"#{session_name}\t#{@zsb_agent_notif}").Output()
	if err != nil {
		return counts
	}
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) < 2 {
			continue
		}
		n, _ := strconv.Atoi(strings.TrimSpace(parts[1])) // empty/unset → 0
		counts[parts[0]] += n
	}
	return counts
}
