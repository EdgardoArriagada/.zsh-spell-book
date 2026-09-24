package main

import (
	"bytes"
	"context"
	"strings"
	"testing"
)

func TestHelp(t *testing.T) {
	t.Setenv("PATH", "")
	for _, flag := range []string{"-h", "--help"} {
		var out bytes.Buffer
		if err := run(context.Background(), []string{flag}, &out); err != nil {
			t.Fatalf("run(%q) = %v", flag, err)
		}
		if got := out.String(); !strings.Contains(got, "Usage: poll-pr-activity") || !strings.Contains(got, "-h, --help") {
			t.Errorf("run(%q) output = %q, want usage and help flags", flag, got)
		}
	}
}

func TestNewActivity(t *testing.T) {
	seen := make(map[string]bool)
	baseline := []activity{{ID: 1, Kind: "comment"}, {ID: 2, Kind: "review", State: "PENDING"}}
	if got := newActivity(seen, baseline); len(got) != 1 {
		t.Fatalf("baseline has %d trackable items, want 1", len(got))
	}
	items := []activity{{ID: 1, Kind: "comment"}, {ID: 2, Kind: "review", State: "APPROVED"}, {ID: 3, Kind: "review comment"}}
	if got := newActivity(seen, items); len(got) != 2 || label(got[0]) != "PR approved" || label(got[1]) != "New review comment" {
		t.Fatalf("new activity = %+v, want approval and review comment", got)
	}
	if got := newActivity(seen, items); len(got) != 0 {
		t.Fatalf("repeat poll reported %d items, want 0", len(got))
	}
}

func TestValidPR(t *testing.T) {
	for _, tc := range []struct {
		value string
		want  bool
	}{
		{"42", true},
		{"https://github.com/owner/repo/pull/42", true},
		{"--repo=other/repo", false},
		{"https://github.com.evil.test/owner/repo/pull/42", false},
		{"https://github.com/owner/repo/pull/42?token=secret", false},
	} {
		if got := validPR(tc.value); got != tc.want {
			t.Errorf("validPR(%q) = %t, want %t", tc.value, got, tc.want)
		}
	}
}
