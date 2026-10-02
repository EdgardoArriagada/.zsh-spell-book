package main

import (
	"bytes"
	"context"
	"encoding/json"
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
	for _, args := range [][]string{{"-h"}, {"--help"}, {"-t", "--help"}, {"--pull-request", "bad", "-h"}, {"--help", "--logs"}} {
		var out bytes.Buffer
		if err := run(context.Background(), args, &out); err != nil {
			t.Fatalf("run(%q) = %v", args, err)
		}
		if got := out.String(); !strings.Contains(got, "Usage: watch-pr-events [-t|--tmux]") || !strings.Contains(got, "-t, --tmux") || !strings.Contains(got, "-l, --logs") || !strings.Contains(got, "--codex-thread UUID") || !strings.Contains(got, "-h, --help") || !strings.Contains(got, "-p, --pull-request PR") {
			t.Errorf("run(%q) output = %q, want usage and supported options", args, got)
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
	for _, args := range [][]string{{"-t"}, {"--tmux"}, {"-p", "42", "-t"}, {"-t", "--tmux"}} {
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

func TestCodexThreadFlag(t *testing.T) {
	for _, args := range [][]string{{"--codex-thread", "bad"}, {"--codex-thread="}, {"--codex-thread", "12345678-1234-1234-1234-123456789abc", "--codex-thread", "bad"}} {
		if err := run(context.Background(), args, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "--codex-thread requires a UUID") {
			t.Errorf("run(%q) = %v, want UUID flag error", args, err)
		}
	}
}

func TestArgumentErrors(t *testing.T) {
	t.Setenv("PATH", "")
	for _, args := range [][]string{
		{"42"}, {"https://github.com/owner/repo/pull/42"},
		{"-p", "42", "extra"}, {"--unknown"},
		{"-p"}, {"--pull-request"}, {"--codex-thread"},
		{"-p="}, {"--pull-request", "bad"}, {"-p", "0"},
		{"-p=--repo=other/repo"},
		{"-p", "https://github.com.evil.test/owner/repo/pull/42"},
		{"-p", "https://github.com/owner/repo/pull/42?token=secret"},
	} {
		if err := run(context.Background(), args, &bytes.Buffer{}); err == nil || strings.Contains(err.Error(), "gh not found") {
			t.Errorf("run(%q) = %v, want argument error before looking up gh", args, err)
		}
	}
	if err := run(context.Background(), []string{"--codex-thread", "bad", "--codex-thread=12345678-1234-1234-1234-123456789abc"}, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "codex not found") {
		t.Errorf("duplicate thread flags = %v, want final UUID accepted", err)
	}
}

func TestPullRequestFlags(t *testing.T) {
	dir := t.TempDir()
	argsFile := filepath.Join(dir, "args")
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte("#!/bin/sh\nprintf '%s\\n' \"$@\" > \"$GH_ARGS\"\nexit 1\n"), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	t.Setenv("GH_ARGS", argsFile)
	for _, tc := range []struct {
		args []string
		want string
	}{
		{nil, ""},
		{[]string{"-p", "42"}, "42\n"},
		{[]string{"--pull-request", "https://github.com/owner/repo/pull/42"}, "https://github.com/owner/repo/pull/42\n"},
		{[]string{"--pull-request=42"}, "42\n"},
		{[]string{"-p=https://github.com/owner/repo/pull/42"}, "https://github.com/owner/repo/pull/42\n"},
		{[]string{"-p", "bad", "--pull-request", "43"}, "43\n"},
		{[]string{"--pull-request", "42", "-p", "43"}, "43\n"},
		{[]string{"-p", "42", "-t", "--tmux=false", "-l", "--logs=false"}, "42\n"},
	} {
		if err := run(context.Background(), tc.args, &bytes.Buffer{}); err == nil || !strings.Contains(err.Error(), "cannot find PR") {
			t.Fatalf("run(%q) = %v, want gh failure", tc.args, err)
		}
		got, err := os.ReadFile(argsFile)
		if err != nil || string(got) != "pr\nview\n--json\nurl\n"+tc.want {
			t.Errorf("run(%q) gh args = %q, error = %v", tc.args, got, err)
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
	if err := run(ctx, []string{"--pull-request", "https://github.com/owner/repo/pull/42"}, &out); err != nil {
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
	if err := run(ctx, []string{"--tmux", "--codex-thread", uuid, "-p", "https://github.com/owner/repo/pull/42"}, &out); err != nil {
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
	if err := run(ctx, []string{"--pull-request", "https://github.com/owner/repo/pull/42"}, &out); err != nil {
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

func TestConversationLog(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	gh := filepath.Join(t.TempDir(), "gh")
	script := "#!/bin/sh\nif [ \"$2\" = repos/owner/repo/pulls/42 ]; then printf '%s\\n' \"$PR_DATA\"; else printf '%s\\n' \"$COMMENT_DATA\"; fi\n"
	if err := os.WriteFile(gh, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PR_DATA", `{"id":42,"title":"PR title","body":"Description","user":{"login":"author"},"created_at":"2026-01-01T00:00:00Z"}`)
	t.Setenv("COMMENT_DATA", `[[{"id":7,"body":"First line\nSecond line $(touch nope)","user":{"login":"ignored"},"created_at":"2026-01-02T00:00:00Z","path":"main.go","line":12,"in_reply_to_id":3,"pull_request_review_id":9,"diff_hunk":"@@ code"}]]`)
	ctx := context.Background()
	endpoint := "repos/owner/repo/pulls/42"
	items, err := fetchActivity(ctx, gh, endpoint)
	if err != nil {
		t.Fatal(err)
	}
	strategy, file, err := logStrategy(gh, endpoint, "https://github.com/owner/repo/pull/42")
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if got := newActivity(make(map[string]bool), items, map[string]bool{"ignored": true}); len(got) != 0 {
		t.Fatal("notification exclusions failed")
	}
	if err := strategy.conversation(ctx, items); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(home, "temp", "watch-pr-events", "owner", "repo", "42.log")
	initial, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got := bytes.Count(initial, []byte("\n")); got != 4 {
		t.Fatalf("initial records = %d, want description and all three activity kinds", got)
	}
	var entry conversationEntry
	decoder := json.NewDecoder(bytes.NewReader(initial))
	if err := decoder.Decode(&entry); err != nil || entry.Activity.Kind != "description" || entry.Activity.Title != "PR title" {
		t.Fatalf("description = %+v, error = %v", entry, err)
	}
	for range 3 {
		if err := decoder.Decode(&entry); err != nil || entry.Activity.Body != "First line\nSecond line $(touch nope)" || entry.Activity.InReplyToID != 3 || entry.Activity.Path != "main.go" || entry.Activity.Line != 12 || entry.Activity.DiffHunk != "@@ code" {
			t.Fatalf("conversation content = %+v, error = %v", entry, err)
		}
	}
	if err := strategy.conversation(ctx, items); err != nil {
		t.Fatal(err)
	}
	unchanged, _ := os.ReadFile(path)
	if !bytes.Equal(initial, unchanged) {
		t.Fatal("unchanged poll duplicated history")
	}
	t.Setenv("PR_DATA", `{"id":42,"title":"PR title","body":"Description","user":{"login":"author"},"created_at":"2026-01-01T00:00:00Z","updated_at":"2026-01-03T00:00:00Z"}`)
	if err := strategy.conversation(ctx, items); err != nil {
		t.Fatal(err)
	}
	unchanged, _ = os.ReadFile(path)
	if !bytes.Equal(initial, unchanged) {
		t.Fatal("PR metadata change duplicated description")
	}
	if err := file.Close(); err != nil {
		t.Fatal(err)
	}
	strategy, resumed, err := logStrategy(gh, endpoint, "https://github.com/owner/repo/pull/42")
	if err != nil {
		t.Fatal(err)
	}
	defer resumed.Close()
	if err := strategy.conversation(ctx, items); err != nil {
		t.Fatal(err)
	}
	unchanged, _ = os.ReadFile(path)
	if !bytes.Equal(initial, unchanged) {
		t.Fatal("restart duplicated history")
	}
	items[0].Body = "Edited body"
	if err := strategy.conversation(ctx, items); err != nil {
		t.Fatal(err)
	}
	edited, _ := os.ReadFile(path)
	if !bytes.HasPrefix(edited, initial) || bytes.Count(edited, []byte("\n")) != 5 || !bytes.Contains(edited[len(initial):], []byte(`"change":"edited"`)) {
		t.Fatalf("edit did not preserve and append history: %s", edited)
	}
	if err := strategy.conversation(ctx, nil); err != nil {
		t.Fatal(err)
	}
	afterDeletion, _ := os.ReadFile(path)
	if !bytes.Equal(edited, afterDeletion) {
		t.Fatal("deletion changed stored history")
	}
	t.Setenv("PR_DATA", `{"id":42,"title":"Edited title","body":"Edited description","user":{"login":"author"},"created_at":"2026-01-01T00:00:00Z"}`)
	if err := strategy.conversation(ctx, nil); err != nil {
		t.Fatal(err)
	}
	descriptionEdit, _ := os.ReadFile(path)
	if !bytes.HasPrefix(descriptionEdit, edited) || bytes.Count(descriptionEdit, []byte("\n")) != 6 || !bytes.Contains(descriptionEdit[len(edited):], []byte("Edited description")) {
		t.Fatal("PR description edit not appended")
	}
	info, err := os.Stat(path)
	if err != nil || info.Mode().Perm() != 0600 {
		t.Fatalf("log must be private: %v, %v", info, err)
	}
}

func TestLogsFlags(t *testing.T) {
	dir := t.TempDir()
	gh := filepath.Join(dir, "gh")
	script := strings.Replace(ghReadyScript, "  printf '[[]]\\n'", "  if [ \"$2\" = repos/owner/repo/pulls/42 ]; then printf '{\"id\":42,\"body\":\"Description\"}\\n'; else printf '[[]]\\n'; fi", 1)
	if err := os.WriteFile(gh, []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	t.Setenv("CHECK_BUCKET", "pass")
	for _, args := range [][]string{{"-l"}, {"--logs"}, {"-l", "--logs"}} {
		t.Run(strings.Join(args, " "), func(t *testing.T) {
			home := t.TempDir()
			t.Setenv("HOME", home)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			var out bytes.Buffer
			if err := run(ctx, append(args, "-p", "42"), &out); err != nil {
				t.Fatal(err)
			}
			data, err := os.ReadFile(filepath.Join(home, "temp", "watch-pr-events", "owner", "repo", "42.log"))
			if err != nil || !bytes.Contains(data, []byte("Description")) || !strings.Contains(out.String(), "PR checks passed:") || !strings.Contains(out.String(), "PR ready to merge:") {
				t.Fatalf("log = %s, error = %v, console = %s", data, err, out.String())
			}
		})
	}
}

func TestLogRejectsPathEscape(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	for _, url := range []string{"https://github.com/../repo/pull/42", "https://github.com/owner/../pull/42"} {
		if _, _, err := logStrategy("gh", "unused", url); err == nil {
			t.Fatalf("accepted path escape: %s", url)
		}
	}
	base := filepath.Join(home, "temp", "watch-pr-events")
	if err := os.MkdirAll(base, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(t.TempDir(), filepath.Join(base, "owner")); err != nil {
		t.Fatal(err)
	}
	if _, _, err := logStrategy("gh", "unused", "https://github.com/owner/repo/pull/42"); err == nil {
		t.Fatal("accepted symlink escape")
	}
}
