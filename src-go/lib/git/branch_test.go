package git

import (
	"os"
	"os/exec"
	"testing"
)

func TestValidateBranchName(t *testing.T) {
	valid := []string{
		"main",
		"feature/my-branch",
		"fix-123",
		"release/v1.0",
		"my_branch",
	}
	for _, name := range valid {
		if err := ValidateBranchName(name); err != nil {
			t.Errorf("ValidateBranchName(%q) returned unexpected error: %v", name, err)
		}
	}

	invalid := []struct {
		name   string
		branch string
	}{
		{"empty", ""},
		{"space", "my branch"},
		{"tilde", "feat~1"},
		{"caret", "feat^0"},
		{"colon", "feat:base"},
		{"question mark", "feat?"},
		{"asterisk", "feat*"},
		{"open bracket", "feat[0]"},
		{"backslash", "feat\\bar"},
		{"double dot", "feat..bar"},
		{"at-brace", "feat@{0}"},
		{"starts with dash", "-mybranch"},
		{"starts with dot", ".mybranch"},
		{"ends with dot", "mybranch."},
		{"ends with slash", "mybranch/"},
		{"ends with .lock", "mybranch.lock"},
	}
	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateBranchName(tt.branch); err == nil {
				t.Errorf("ValidateBranchName(%q) expected error, got nil", tt.branch)
			}
		})
	}
}

func TestRemoteBranchExists(t *testing.T) {
	dir := t.TempDir()
	prev, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	t.Cleanup(func() { _ = os.Chdir(prev) })
	if err := os.Chdir(dir); err != nil {
		t.Fatalf("chdir: %v", err)
	}

	run := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", args...)
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v\n%s", args, err, out)
		}
	}

	run("init", "-q", "-b", "main")
	run("-c", "user.email=t@t", "-c", "user.name=t", "commit", "--allow-empty", "-q", "-m", "init")
	run("update-ref", "refs/remotes/origin/foo", "HEAD")
	run("update-ref", "refs/remotes/upstream/feature/nested", "HEAD")

	tests := []struct {
		branch string
		want   bool
	}{
		{"foo", true},
		{"feature/nested", true},
		{"missing", false},
		{"origin", false}, // remote name itself must not match
	}
	for _, tt := range tests {
		t.Run(tt.branch, func(t *testing.T) {
			if got := RemoteBranchExists(tt.branch); got != tt.want {
				t.Errorf("RemoteBranchExists(%q) = %v, want %v", tt.branch, got, tt.want)
			}
		})
	}
}
