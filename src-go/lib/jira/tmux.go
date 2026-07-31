package jira

import (
	"regexp"
	"strings"
)

var (
	invalidChars      = regexp.MustCompile(`[^a-zA-Z0-9-]`)
	consecutiveDashes = regexp.MustCompile(`-{2,}`)
)

func (t Ticket) TmuxSessionID() string {
	kebab := strings.ToLower(strings.ReplaceAll(t.Label, " ", "-"))
	result := sanitizeTmuxID(t.Current + "-" + kebab)
	if len(result) > 25 {
		result = result[:25]
	}
	return strings.TrimRight(result, "-")
}

func TmuxSessionPrefixes(currentIDs map[string]bool) []string {
	prefixes := make([]string, 0, len(currentIDs))
	for id := range currentIDs {
		prefixes = append(prefixes, sanitizeTmuxID(id)+"-")
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

func sanitizeTmuxID(s string) string {
	s = invalidChars.ReplaceAllString(s, "-")
	s = consecutiveDashes.ReplaceAllString(s, "-")
	return s
}
