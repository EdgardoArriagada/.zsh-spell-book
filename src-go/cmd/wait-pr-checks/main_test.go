package main

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTmuxFlag(t *testing.T) {
	dir := t.TempDir()
	gh := filepath.Join(dir, "gh")
	if err := os.WriteFile(gh, []byte("#!/bin/sh\nprintf '[{\"bucket\":\"pass\"}]\\n'\n"), 0700); err != nil {
		t.Fatal(err)
	}
	notify := filepath.Join(dir, "zsb_tmux_agent_notification")
	if err := os.WriteFile(notify, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" >> \"$NOTIFICATION_ARGS\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	argsFile := filepath.Join(dir, "args")
	t.Setenv("PATH", dir)
	t.Setenv("NOTIFICATION_ARGS", argsFile)
	t.Setenv("TMUX_PANE", "")
	if err := run([]string{"42"}); err != nil {
		t.Fatalf("without tmux: %v", err)
	}
	if _, err := os.Stat(argsFile); !os.IsNotExist(err) {
		t.Fatalf("notification without -t: %v", err)
	}
	for _, flag := range []string{"-t", "--tmux"} {
		if err := run([]string{"42", flag}); err == nil || !strings.Contains(err.Error(), "inside a tmux pane") {
			t.Errorf("%s outside tmux: %v", flag, err)
		}
	}
	t.Setenv("TMUX_PANE", "%42")
	for _, args := range [][]string{{"-t", "42"}, {"42", "--tmux"}} {
		if err := run(args); err != nil {
			t.Fatalf("run(%q): %v", args, err)
		}
	}
	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatal(err)
	}
	want := "--working\n_\n%42\n--force-finished\n_\n%42\n"
	if string(got) != want+want {
		t.Errorf("notification args = %q, want %q", got, want+want)
	}
}

func TestStatus(t *testing.T) {
	cases := []struct {
		name         string
		checks       []check
		done, failed bool
	}{
		{"not started", nil, false, false},
		{"pending", []check{{Bucket: "pass"}, {Bucket: "pending"}}, false, false},
		{"success", []check{{Bucket: "pass"}, {Bucket: "skipping"}}, true, false},
		{"failure", []check{{Bucket: "pass"}, {Bucket: "cancel"}}, true, true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			done, failed := status(tc.checks)
			if done != tc.done || failed != tc.failed {
				t.Fatalf("status() = (%t, %t), want (%t, %t)", done, failed, tc.done, tc.failed)
			}
		})
	}
}

func TestFetchChecksFailedExit(t *testing.T) {
	dir := t.TempDir()
	gh := filepath.Join(dir, "gh")
	if err := os.WriteFile(gh, []byte("#!/bin/sh\nprintf '[{\"bucket\":\"fail\"}]\\n'\nexit 1\n"), 0700); err != nil {
		t.Fatal(err)
	}
	checks, err := fetchChecks(context.Background(), gh, []string{"https://github.com/owner/repo/pull/42"})
	if err != nil || len(checks) != 1 || checks[0].Bucket != "fail" {
		t.Fatalf("failed checks = %v, %v", checks, err)
	}
}
