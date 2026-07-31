package jira

import "testing"

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
		got := sanitizeTmuxID(c.in)
		if got != c.want {
			t.Errorf("sanitizeTmuxID(%q) = %q, want %q", c.in, got, c.want)
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
			name:    "result exactly 25 chars",
			current: "PROJ-1234",
			label:   "abcdefghijklmno",
			want:    "PROJ-1234-abcdefghijklmno",
		},
		{
			name:    "label uppercase lowercased",
			current: "PROJ-1",
			label:   "MY LABEL",
			want:    "PROJ-1-my-label",
		},
		{
			name:    "truncated no trailing dash",
			current: "PROJ-123",
			label:   "abcdefghijklmnopqrstuvwxyz",
			want:    "PROJ-123-abcdefghijklmnop",
		},
		{
			// full result ends with '-', char at index 24 IS '-' → trimmed
			name:    "truncated trailing dash trimmed",
			current: "PROJ-123",
			label:   "abcdefghijklmno more.",
			want:    "PROJ-123-abcdefghijklmno",
		},
		{
			// full result ends with '-', char at index 24 is NOT '-' → no-op trim
			name:    "truncated full ends dash but pos24 clean",
			current: "PROJ-123",
			label:   "abcdefghijklmnop more.",
			want:    "PROJ-123-abcdefghijklmnop",
		},
		{
			// full result does NOT end with '-', but position 24 IS '-' → trailing dash kept
			name:    "truncated pos24 is dash but full result clean",
			current: "PROJ-123",
			label:   "abcdefghijklmno more words",
			want:    "PROJ-123-abcdefghijklmno",
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
			ticket := Ticket{Current: c.current, Label: c.label}
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
			prefixes: TmuxSessionPrefixes(map[string]bool{"PROJ-123": true}),
			want:     true,
		},
		{
			name:     "session without trailing dash suffix",
			session:  "PROJ-123",
			prefixes: TmuxSessionPrefixes(map[string]bool{"PROJ-123": true}),
			want:     false,
		},
		{
			name:     "different prefix",
			session:  "OTHER-456-something",
			prefixes: TmuxSessionPrefixes(map[string]bool{"PROJ-123": true}),
			want:     false,
		},
		{
			name:     "empty session",
			session:  "",
			prefixes: TmuxSessionPrefixes(map[string]bool{"PROJ-123": true}),
			want:     false,
		},
		{
			name:     "current with special chars sanitized before compare",
			session:  "PROJ-123-my-ticket",
			prefixes: TmuxSessionPrefixes(map[string]bool{"PROJ.123": true}),
			want:     true,
		},
		{
			name:     "longer current does not partially match",
			session:  "PROJ-12-something",
			prefixes: TmuxSessionPrefixes(map[string]bool{"PROJ-1": true}),
			want:     false,
		},
		{
			name:     "matches one of multiple IDs",
			session:  "PROJ-456-foo",
			prefixes: TmuxSessionPrefixes(map[string]bool{"PROJ-123": true, "PROJ-456": true}),
			want:     true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got := MatchesTmuxSession(c.session, c.prefixes)
			if got != c.want {
				t.Errorf("MatchesTmuxSession(%q, %v) = %v, want %v", c.session, c.prefixes, got, c.want)
			}
		})
	}
}
