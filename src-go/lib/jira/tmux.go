package jira

import (
	"os/exec"
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

// NotifCounts is the per-session tally of agent notification states:
// panes finished (value 1, red badge) and still working (value 2, blue badge).
type NotifCounts struct {
	Working  int
	Finished int
}

func LoadNotificationCounts() map[string]NotifCounts {
	out, err := exec.Command("tmux", "list-panes", "-a", "-F",
		"#{session_name}\t#{@zsb_agent_notif}").Output()
	if err != nil {
		return map[string]NotifCounts{}
	}
	return parseNotifCounts(string(out))
}

// parseNotifCounts tallies panes per session by state: value "1" → Finished,
// "2" → Working. Everything else (unset/empty/"0") is ignored.
func parseNotifCounts(out string) map[string]NotifCounts {
	counts := map[string]NotifCounts{}
	for line := range strings.SplitSeq(strings.TrimSpace(out), "\n") {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) < 2 {
			continue
		}
		c := counts[parts[0]]
		switch strings.TrimSpace(parts[1]) {
		case "1":
			c.Finished++
		case "2":
			c.Working++
		default:
			continue
		}
		counts[parts[0]] = c
	}
	return counts
}
