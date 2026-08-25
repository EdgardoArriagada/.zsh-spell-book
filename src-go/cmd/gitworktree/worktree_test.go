package main

import (
	"os"
	"os/exec"
	"path/filepath"
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
