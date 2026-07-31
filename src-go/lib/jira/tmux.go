package jira

import (
	"regexp"
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
