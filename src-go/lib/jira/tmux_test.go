package jira_test

import (
	"testing"

	"example.com/workspace/lib/jira"
)

func TestSanitizeTmuxID(t *testing.T) {
	cases := []struct {
		in   string
		want string
	}{
		{"abc123", "abc123"},          // purely alphanumeric — unchanged
		{"abc-def", "abc-def"},        // dashes allowed — unchanged
		{"abc.", "abc-"},              // single special char → dash
		{"abc...", "abc-"},            // consecutive special chars → single dash
		{".abc", "-abc"},              // leading special char → leading dash
		{"abc.", "abc-"},              // trailing special char → trailing dash
		{"abc.def-ghi", "abc-def-ghi"}, // mixed valid + invalid
	}
	for _, c := range cases {
		got := jira.SanitizeTmuxID(c.in)
		if got != c.want {
			t.Errorf("SanitizeTmuxID(%q) = %q, want %q", c.in, got, c.want)
		}
	}
}

func TestTmuxSessionID(t *testing.T) {
	cases := []struct {
		name    string
		current string
		label   string
		want    string
	}{
		{
			name:    "short normal",
			current: "PROJ-123",
			label:   "my ticket",
			want:    "PROJ-123-my-ticket",
		},
		{
			name:    "result exactly 38 chars",
			current: "PROJ-1234",
			label:   "abcdefghijklmnopqrstuvwxyzab",
			want:    "PROJ-1234-abcdefghijklmnopqrstuvwxyzab",
		},
		{
			name:    "label uppercase lowercased",
			current: "PROJ-1",
			label:   "MY LABEL",
			want:    "PROJ-1-my-label",
		},
		{
			name:    "label bracketed text removed",
			current: "PROJ-123",
			label:   "[team] Fix login [urgent]",
			want:    "PROJ-123-fix-login",
		},
		{
			name:    "truncated no trailing dash",
			current: "PROJ-123",
			label:   "abcdefghijklmnopqrstuvwxyzabcde",
			want:    "PROJ-123-abcdefghijklmnopqrstuvwxyzabc",
		},
		{
			// result[:38] ends with '-' → trimmed
			name:    "truncated trailing dash trimmed",
			current: "PROJ-123",
			label:   "abcdefghijklmnopqrstuvwxyzab more.",
			want:    "PROJ-123-abcdefghijklmnopqrstuvwxyzab",
		},
		{
			// full result ends with '-', but result[:38] does not → no trim
			name:    "truncated full ends dash but pos37 clean",
			current: "PROJ-123",
			label:   "abcdefghijklmnopqrstuvwxyzabc more.",
			want:    "PROJ-123-abcdefghijklmnopqrstuvwxyzabc",
		},
		{
			name:    "label has special chars",
			current: "A",
			label:   "fix: bug (urgent)",
			want:    "A-fix-bug-urgent",
		},
		{
			name:    "label consecutive spaces collapsed",
			current: "A",
			label:   "hello  world",
			want:    "A-hello-world",
		},
		{
			name:    "current with special chars sanitized",
			current: "PROJ.123",
			label:   "my ticket",
			want:    "PROJ-123-my-ticket",
		},
		{
			name:    "empty label produces trailing dash",
			current: "PROJ-123",
			label:   "",
			want:    "PROJ-123",
		},
		{
			name:    "label only special chars collapses to single dash",
			current: "A",
			label:   "...",
			want:    "A",
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			ticket := jira.Ticket{Current: c.current, Label: c.label}
			got := ticket.TmuxSessionID()
			if got != c.want {
				t.Errorf("got %q, want %q", got, c.want)
			}
		})
	}
}

func TestMatchesTmuxSession(t *testing.T) {
	cases := []struct {
		name     string
		session  string
		prefixes []string
		want     bool
	}{
		{
			name:     "matching session",
			session:  "PROJ-123-my-ticket",
			prefixes: jira.TmuxSessionPrefixes(map[string]bool{"PROJ-123": true}),
			want:     true,
		},
		{
			name:     "session without trailing dash suffix",
			session:  "PROJ-123",
			prefixes: jira.TmuxSessionPrefixes(map[string]bool{"PROJ-123": true}),
			want:     false,
		},
		{
			name:     "different prefix",
			session:  "OTHER-456-something",
			prefixes: jira.TmuxSessionPrefixes(map[string]bool{"PROJ-123": true}),
			want:     false,
		},
		{
			name:     "empty session",
			session:  "",
			prefixes: jira.TmuxSessionPrefixes(map[string]bool{"PROJ-123": true}),
			want:     false,
		},
		{
			name:     "current with special chars sanitized before compare",
			session:  "PROJ-123-my-ticket",
			prefixes: jira.TmuxSessionPrefixes(map[string]bool{"PROJ.123": true}),
			want:     true,
		},
		{
			name:     "longer current does not partially match",
			session:  "PROJ-12-something",
			prefixes: jira.TmuxSessionPrefixes(map[string]bool{"PROJ-1": true}),
			want:     false,
		},
		{
			name:     "matches one of multiple IDs",
			session:  "PROJ-456-foo",
			prefixes: jira.TmuxSessionPrefixes(map[string]bool{"PROJ-123": true, "PROJ-456": true}),
			want:     true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := jira.MatchesTmuxSession(c.session, c.prefixes)
			if got != c.want {
				t.Errorf("MatchesTmuxSession(%q, %v) = %v, want %v", c.session, c.prefixes, got, c.want)
			}
		})
	}
}
