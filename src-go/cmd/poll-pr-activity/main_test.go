package main

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
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
	if got := newActivity(seen, baseline, nil); len(got) != 1 {
		t.Fatalf("baseline has %d trackable items, want 1", len(got))
	}
	items := []activity{{ID: 1, Kind: "comment"}, {ID: 2, Kind: "review", State: "APPROVED"}, {ID: 3, Kind: "review comment"}}
	if got := newActivity(seen, items, nil); len(got) != 2 || label(got[0]) != "PR approved" || label(got[1]) != "New review comment" {
		t.Fatalf("new activity = %+v, want approval and review comment", got)
	}
	if got := newActivity(seen, items, nil); len(got) != 0 {
		t.Fatalf("repeat poll reported %d items, want 0", len(got))
	}
}

func TestIgnoredUsers(t *testing.T) {
	path := filepath.Join(t.TempDir(), "poll-pr-activity.conf")
	if err := os.WriteFile(path, []byte("# my account\n\n  MyUser  \nother-user\n"), 0600); err != nil {
		t.Fatal(err)
	}
	ignored, err := loadIgnoredUsers(path)
	if err != nil {
		t.Fatal(err)
	}
	if !ignored["myuser"] || !ignored["other-user"] || len(ignored) != 2 {
		t.Fatalf("ignored users = %v", ignored)
	}
	items := []activity{{ID: 1, Kind: "comment"}, {ID: 2, Kind: "review"}, {ID: 3, Kind: "review comment"}}
	items[0].User.Login = "MYUSER"
	items[1].User.Login = "Other-User"
	items[2].User.Login = "someone-else"
	seen := make(map[string]bool)
	if got := newActivity(seen, items, ignored); len(got) != 1 || got[0].ID != 3 {
		t.Fatalf("new activity = %+v, want only non-ignored user", got)
	}
	if len(seen) != 1 {
		t.Fatalf("tracked activity = %v, want only non-ignored user", seen)
	}
}

func TestMissingIgnoreConfig(t *testing.T) {
	ignored, err := loadIgnoredUsers(filepath.Join(t.TempDir(), "poll-pr-activity.conf"))
	if err != nil || len(ignored) != 0 {
		t.Fatalf("missing config: ignored = %v, err = %v", ignored, err)
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

func TestRunReportsMergeReady(t *testing.T) {
	dir := t.TempDir()
	gh := filepath.Join(dir, "gh")
	script := "#!/bin/sh\nif [ \"$1\" = api ]; then printf '[[]]\\n'; elif [ \"$3\" = --json ]; then printf '{\"url\":\"https://github.com/owner/repo/pull/42\"}\\n'; else printf '{\"state\":\"OPEN\",\"mergeStateStatus\":\"CLEAN\"}\\n'; fi\n"
	if err := os.WriteFile(gh, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var out bytes.Buffer
	if err := run(ctx, []string{"https://github.com/owner/repo/pull/42"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "PR ready to merge: https://github.com/owner/repo/pull/42") {
		t.Errorf("watcher output = %q, want merge readiness", out.String())
	}
}
