package main

import (
	"bytes"
	"context"
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

const ghReadyScript = `#!/bin/sh
if [ "$1" = api ]; then
  printf '[[]]\n'
elif [ "$2" = checks ]; then
  printf '[{"bucket":"%s"}]\n' "$CHECK_BUCKET"
  [ "$CHECK_BUCKET" != fail ]
elif [ "$4" = url ]; then
  printf '{"url":"https://github.com/owner/repo/pull/42"}\n'
elif [ "$5" = headRefOid ]; then
  printf '{"headRefOid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}\n'
else
  printf '{"state":"OPEN","mergeStateStatus":"CLEAN","headRefOid":"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}\n'
fi
`

func TestHelp(t *testing.T) {
	t.Setenv("PATH", "")
	for _, flag := range []string{"-h", "--help"} {
		var out bytes.Buffer
		if err := run(context.Background(), []string{flag}, &out); err != nil {
			t.Fatalf("run(%q) = %v", flag, err)
		}
		if got := out.String(); !strings.Contains(got, "Usage: watch-pr-events [-t|--tmux]") || !strings.Contains(got, "-t, --tmux") || !strings.Contains(got, "--codex-uuid UUID") || !strings.Contains(got, "-h, --help") {
			t.Errorf("run(%q) output = %q, want usage and supported options", flag, got)
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
	path := filepath.Join(t.TempDir(), "watch-pr-events.conf")
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
	ignored, err := loadIgnoredUsers(filepath.Join(t.TempDir(), "watch-pr-events.conf"))
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

func TestCodexUUIDFlag(t *testing.T) {
	for _, args := range [][]string{{"--codex-uuid"}, {"--codex-uuid", "bad"}, {"--codex-uuid", "12345678-1234-1234-1234-123456789abc", "--codex-uuid", "12345678-1234-1234-1234-123456789abc"}} {
		if err := run(context.Background(), args, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "--codex-uuid requires a UUID") {
			t.Errorf("run(%q) = %v, want UUID flag error", args, err)
		}
	}
}

func TestQueueCodexEvent(t *testing.T) {
	command := filepath.Join(t.TempDir(), "codex")
	if err := os.WriteFile(command, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$CODEX_ARGS\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	argsFile := filepath.Join(t.TempDir(), "args")
	t.Setenv("CODEX_ARGS", argsFile)
	uuid := "12345678-1234-1234-1234-123456789abc"
	unexpected := filepath.Join(t.TempDir(), "unexpected")
	event := `New PR comment by "$(touch ` + unexpected + `)": https://github.com/owner/repo/pull/42`
	if err := queueCodexEvent(context.Background(), command, uuid, event); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatal(err)
	}
	if want := "queue\n--thread\n" + uuid + "\n--message\n" + event + "\n"; string(got) != want {
		t.Errorf("codex args = %q, want %q", got, want)
	}
	if _, err := os.Stat(unexpected); !os.IsNotExist(err) {
		t.Errorf("event text executed as shell syntax: %v", err)
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
	if err := os.WriteFile(gh, []byte(ghReadyScript), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CHECK_BUCKET", "pass")
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
	if !strings.Contains(out.String(), "PR checks passed: https://github.com/owner/repo/pull/42 (head aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa)") {
		t.Errorf("watcher output = %q, want passed checks", out.String())
	}
}

func TestRunDeliversMergeReadyToAllStrategies(t *testing.T) {
	dir := t.TempDir()
	gh := filepath.Join(dir, "gh")
	if err := os.WriteFile(gh, []byte(ghReadyScript), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CHECK_BUCKET", "pass")
	command := filepath.Join(dir, "codex")
	if err := os.WriteFile(command, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" >> \"$CODEX_ARGS\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	argsFile := filepath.Join(dir, "args")
	t.Setenv("CODEX_ARGS", argsFile)
	notify := filepath.Join(dir, "zsb_tmux_agent_notification")
	if err := os.WriteFile(notify, []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$NOTIFICATION_ARGS\"\n"), 0700); err != nil {
		t.Fatal(err)
	}
	notificationArgs := filepath.Join(dir, "notification-args")
	t.Setenv("NOTIFICATION_ARGS", notificationArgs)
	t.Setenv("TMUX_PANE", "%42")
	t.Setenv("PATH", dir)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	uuid := "12345678-1234-1234-1234-123456789abc"
	var out bytes.Buffer
	if err := run(ctx, []string{"--tmux", "--codex-uuid", uuid, "https://github.com/owner/repo/pull/42"}, &out); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out.String(), "PR ready to merge: https://github.com/owner/repo/pull/42") {
		t.Errorf("watcher output = %q, want merge readiness", out.String())
	}
	gotNotification, err := os.ReadFile(notificationArgs)
	if err != nil {
		t.Fatal(err)
	}
	if string(gotNotification) != "--force-finished\n_\n%42\n" {
		t.Errorf("notification args = %q", gotNotification)
	}
	got, err := os.ReadFile(argsFile)
	if err != nil {
		t.Fatal(err)
	}
	if want := "queue\n--thread\n" + uuid + "\n--message\nPR checks passed: https://github.com/owner/repo/pull/42 (head aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa)\n" +
		"queue\n--thread\n" + uuid + "\n--message\nPR ready to merge: https://github.com/owner/repo/pull/42 (head aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa)\n"; string(got) != want {
		t.Errorf("codex args = %q, want %q", got, want)
	}
}

func TestRunReportsFailedChecksWithoutMergeReady(t *testing.T) {
	dir := t.TempDir()
	gh := filepath.Join(dir, "gh")
	if err := os.WriteFile(gh, []byte(ghReadyScript), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CHECK_BUCKET", "fail")
	t.Setenv("PATH", dir)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var out bytes.Buffer
	if err := run(ctx, []string{"https://github.com/owner/repo/pull/42"}, &out); err != nil {
		t.Fatal(err)
	}
	if got := out.String(); !strings.Contains(got, "PR checks failed: https://github.com/owner/repo/pull/42 (head aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa)") || strings.Contains(got, "PR ready to merge") {
		t.Errorf("watcher output = %q, want failed checks without merge readiness", got)
	}
}

func TestFetchCheckStatusRejectsChangedHead(t *testing.T) {
	dir := t.TempDir()
	gh := filepath.Join(dir, "gh")
	headFile := filepath.Join(dir, "head-called")
	script := "#!/bin/sh\nif [ \"$2\" = checks ]; then printf '[{\"bucket\":\"pass\"}]\\n'; elif [ -e \"$HEAD_FILE\" ]; then printf '{\"headRefOid\":\"bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb\"}\\n'; else : > \"$HEAD_FILE\"; printf '{\"headRefOid\":\"aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa\"}\\n'; fi\n"
	if err := os.WriteFile(gh, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("HEAD_FILE", headFile)
	head, result, err := fetchCheckStatus(context.Background(), gh, "https://github.com/owner/repo/pull/42")
	if err != nil || head != "" || result != "" {
		t.Fatalf("changed head status = (%q, %q, %v), want no result", head, result, err)
	}
}

func TestChecksStatus(t *testing.T) {
	for _, tc := range []struct {
		name   string
		checks []check
		want   string
	}{
		{"not started", nil, "pending"},
		{"pending", []check{{Bucket: "pass"}, {Bucket: "pending"}}, "pending"},
		{"success", []check{{Bucket: "pass"}, {Bucket: "skipping"}}, "passed"},
		{"failure", []check{{Bucket: "pass"}, {Bucket: "cancel"}}, "failed"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := checksStatus(tc.checks); got != tc.want {
				t.Fatalf("checksStatus() = %q, want %q", got, tc.want)
			}
		})
	}
}

func TestNotifyStrategiesContinuesAfterFailure(t *testing.T) {
	var calls []string
	strategies := []notificationStrategy{
		{notify: func(_ context.Context, event string, mergeReady bool) error {
			if event != "comment\nreview\nPR: url" || mergeReady {
				t.Errorf("first strategy got event %q, mergeReady %t", event, mergeReady)
			}
			calls = append(calls, "failed")
			return errors.New("failed")
		}, failure: "failed", retry: true},
		{notify: func(_ context.Context, event string, mergeReady bool) error {
			if event != "comment\nreview\nPR: url" || mergeReady {
				t.Errorf("second strategy got event %q, mergeReady %t", event, mergeReady)
			}
			calls = append(calls, "next")
			return nil
		}},
	}
	if retry := notifyStrategies(context.Background(), strategies, "comment\nreview\nPR: url", false); !retry || strings.Join(calls, ",") != "failed,next" {
		t.Errorf("retry = %t, calls = %v", retry, calls)
	}
}
