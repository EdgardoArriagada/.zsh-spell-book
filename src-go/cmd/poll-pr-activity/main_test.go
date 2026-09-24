package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
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
		if got := out.String(); !strings.Contains(got, "Usage: poll-pr-activity") || !strings.Contains(got, "-h, --help") || !strings.Contains(got, "-t, --tmux") {
			t.Errorf("run(%q) output = %q, want usage and flags", flag, got)
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

func TestTmuxFlag(t *testing.T) {
	t.Setenv("TMUX_PANE", "")
	for _, args := range [][]string{{"-t"}, {"--tmux"}, {"42", "-t"}} {
		if err := run(context.Background(), args, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "inside a tmux pane") {
			t.Errorf("run(%q) = %v, want tmux pane error", args, err)
		}
	}
}

func TestTmuxNotification(t *testing.T) {
	command := filepath.Join(t.TempDir(), "notify")
	if err := os.WriteFile(command, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$NOTIFICATION_ARGS\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	argsFile := filepath.Join(t.TempDir(), "args")
	t.Setenv("NOTIFICATION_ARGS", argsFile)
	if err := tmuxNotification(context.Background(), command, "%42"); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "--force-finished\n_\n%42\n" {
		t.Errorf("notification args = %q", got)
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
