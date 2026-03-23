package main_test

import (
	"testing"

	gitworktree "gitworktree"
)

func TestWorktreeBaseDir(t *testing.T) {
	tests := []struct {
		mainPath string
		expect   string
	}{
		{"/home/user/myrepo", "/home/user/myrepo_gitworktree"},
		{"/tmp/project", "/tmp/project_gitworktree"},
		{"/a/b/c", "/a/b/c_gitworktree"},
	}

	for _, tt := range tests {
		got := gitworktree.WorktreeBaseDir(tt.mainPath)
		if got != tt.expect {
			t.Errorf("WorktreeBaseDir(%q) = %q, want %q", tt.mainPath, got, tt.expect)
		}
	}
}
