package main

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

func TestBaseBranch(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })
	if err := os.Chdir(dir); err != nil {
		t.Fatal(err)
	}

	run := func(args ...string) {
		t.Helper()
		if out, err := exec.Command("git", args...).CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\\n%s", args, err, out)
		}
	}
	run("init", "-q", "-b", "main")
	run("-c", "user.email=t@t", "-c", "user.name=t", "commit", "--allow-empty", "-q", "-m", "init")
	run("update-ref", "refs/remotes/origin/main", "HEAD")
	run("update-ref", "refs/remotes/origin/develop", "HEAD")

	for _, tt := range []struct {
		branch string
		want   string
	}{
		{"feature/new", "refs/remotes/origin/develop"},
		{"hotfix/urgent", "refs/remotes/origin/main"},
		{"fix/urgent", "refs/remotes/origin/main"},
	} {
		t.Run(tt.branch, func(t *testing.T) {
			got, err := baseBranch(tt.branch)
			if err != nil || got != tt.want {
				t.Errorf("baseBranch(%q) = %q, %v; want %q", tt.branch, got, err, tt.want)
			}
		})
	}

	run("update-ref", "refs/remotes/origin/master", "HEAD")
	for _, branch := range []string{"hotfix/urgent", "fix/urgent"} {
		if got, err := baseBranch(branch); err != nil || got != "refs/remotes/origin/master" {
			t.Errorf("baseBranch(%q) = %q, %v; want origin/master", branch, got, err)
		}
	}

	run("update-ref", "-d", "refs/remotes/origin/develop")
	run("branch", "develop")
	if got, err := baseBranch("feature/local"); err != nil || got != "develop" {
		t.Errorf("baseBranch(feature/local) = %q, %v; want local develop", got, err)
	}
	run("branch", "-D", "develop")
	if got, err := baseBranch("feature/fallback"); err != nil || got != "refs/remotes/origin/master" {
		t.Errorf("baseBranch(feature/fallback) = %q, %v; want origin/master", got, err)
	}
}

func TestWorktreePRStatus(t *testing.T) {
	dir := t.TempDir()
	script := `#!/bin/sh
printf '%s\n' "$PWD" "$@" > "$PR_ARGS"
if [ "$4" = --state=open ]; then
  printf '%s' "$PR_OPEN_RESPONSE"
  exit 0
fi
if [ "$PR_RESPONSE" = error ]; then
  echo 'secret-token' >&2
  exit 1
fi
printf '%s' "$PR_RESPONSE"
`
	if err := os.WriteFile(filepath.Join(dir, "gh"), []byte(script), 0700); err != nil {
		t.Fatal(err)
	}
	argsFile := filepath.Join(dir, "args")
	t.Setenv("PATH", dir+string(os.PathListSeparator)+os.Getenv("PATH"))
	t.Setenv("PR_ARGS", argsFile)
	wt := Worktree{Path: dir, Branch: "feature/pr-status"}
	for _, tt := range []struct {
		name, response, want string
		wantErr              bool
	}{
		{"none", `[]`, "no PR", false},
		{"open preferred", `[{"number":2,"state":"OPEN","reviewDecision":"APPROVED","url":"https://github.com/o/r/pull/2","statusCheckRollup":[{"status":"COMPLETED","conclusion":"SUCCESS"},{"status":"IN_PROGRESS"},{"state":"FAILURE"}]}]`, "#2 open\nchecks: 1 passed, 1 pending, 1 failed\nreview: approved\nhttps://github.com/o/r/pull/2", false},
		{"draft", `[{"number":3,"state":"OPEN","isDraft":true,"reviewDecision":"REVIEW_REQUIRED"}]`, "#3 draft\nchecks: none\nreview: review required\n", false},
		{"merged", `[{"number":4,"state":"MERGED","isDraft":true}]`, "#4 merged\nchecks: none\nreview: none\n", false},
		{"closed", `[{"number":5,"state":"CLOSED"}]`, "#5 closed\nchecks: none\nreview: none\n", false},
		{"checks", `[{"number":6,"state":"OPEN","statusCheckRollup":[{"state":"SUCCESS"},{"state":"PENDING"},{"status":"COMPLETED","conclusion":"NEUTRAL"},{"status":"COMPLETED","conclusion":"SKIPPED"},{"status":"COMPLETED","conclusion":"CANCELLED"},{"status":"COMPLETED","conclusion":"TIMED_OUT"},{"status":"COMPLETED","conclusion":"ACTION_REQUIRED"},{"state":"ERROR"}]}]`, "#6 open\nchecks: 3 passed, 1 pending, 4 failed\nreview: none\n", false},
		{"terminal controls", `[{"number":7,"state":"OPEN","url":"url\u001b\u0007\n"}]`, "#7 open\nchecks: none\nreview: none\nurl", false},
		{"failure", "error", "PR lookup failed; check gh authentication and repository access", true},
		{"invalid JSON", "broken", "invalid GitHub PR response", true},
	} {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("PR_RESPONSE", tt.response)
			open := "[]"
			if strings.Contains(tt.response, `"state":"OPEN"`) {
				open = tt.response
				t.Setenv("PR_RESPONSE", "error")
			}
			t.Setenv("PR_OPEN_RESPONSE", open)
			got, err := worktreePRStatus(context.Background(), wt)
			if tt.wantErr {
				if err == nil || err.Error() != tt.want {
					t.Fatalf("got %q, %v; want error %q", got, err, tt.want)
				}
			} else if err != nil || got != tt.want {
				t.Fatalf("got %q, %v; want %q", got, err, tt.want)
			}
		})
	}
	args, err := os.ReadFile(argsFile)
	if err != nil || !strings.HasPrefix(string(args), dir+"\npr\nlist\n--head="+wt.Branch+"\n--state=all\n") {
		t.Fatalf("wrong worktree or arguments: %q, %v", args, err)
	}
	for _, wt := range []Worktree{{IsBare: true}, {Branch: "(detached)"}, {}} {
		if got, err := worktreePRStatus(context.Background(), wt); err != nil || !strings.HasPrefix(got, "no PR") {
			t.Fatalf("branchless worktree: %q, %v", got, err)
		}
	}
	wt.Branch = "--repo=other"
	if _, err := worktreePRStatus(context.Background(), wt); err == nil {
		t.Fatal("accepted invalid branch")
	}
	wt.Branch = "feature/pr-status"
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if _, err := worktreePRStatus(ctx, wt); err == nil || !strings.Contains(err.Error(), "cancelled") {
		t.Fatalf("cancelled lookup: %v", err)
	}
	git, err := exec.LookPath("git")
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Symlink(git, filepath.Join(dir, "git")); err != nil {
		t.Fatal(err)
	}
	if err := os.Remove(filepath.Join(dir, "gh")); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir)
	if _, err := worktreePRStatus(context.Background(), wt); err == nil || !strings.Contains(err.Error(), "gh not found") {
		t.Fatalf("missing gh: %v", err)
	}
}

func TestCreateWorktreeFetchesRemoteBase(t *testing.T) {
	dir := t.TempDir()
	run := func(path string, args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = path
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git -C %s %v: %v\n%s", path, args, err, out)
		}
		return string(out)
	}

	remote := filepath.Join(dir, "remote.git")
	seed := filepath.Join(dir, "seed")
	mainPath := filepath.Join(dir, "repo")
	run(dir, "init", "--bare", "-q", remote)
	run(dir, "init", "-q", "-b", "main", seed)
	run(seed, "-c", "user.email=t@t", "-c", "user.name=t", "commit", "--allow-empty", "-q", "-m", "main")
	run(seed, "remote", "add", "origin", remote)
	run(seed, "push", "-q", "-u", "origin", "main")
	run(remote, "symbolic-ref", "HEAD", "refs/heads/main")
	run(seed, "checkout", "-q", "-b", "develop")
	run(seed, "push", "-q", "-u", "origin", "develop")
	run(dir, "clone", "-q", remote, mainPath)
	run(seed, "commit", "--allow-empty", "-q", "-m", "develop update")
	run(seed, "push", "-q")

	prev, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })
	if err := os.Chdir(mainPath); err != nil {
		t.Fatal(err)
	}
	if err := createWorktree(mainPath, "feature/latest"); err != nil {
		t.Fatal(err)
	}

	got := run(filepath.Join(WorktreeBaseDir(mainPath), "feature/latest"), "rev-parse", "HEAD")
	want := run(seed, "rev-parse", "HEAD")
	if got != want {
		t.Errorf("feature/latest is at %s, want fetched develop %s", got, want)
	}

	worktreePath := filepath.Join(WorktreeBaseDir(mainPath), "feature/latest")
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "@{upstream}")
	cmd.Dir = worktreePath
	if out, err := cmd.CombinedOutput(); err == nil {
		t.Errorf("feature/latest unexpectedly tracks %s", out)
	}
}
