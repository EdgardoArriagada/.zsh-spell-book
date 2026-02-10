package main_test

import (
	"testing"

	gitworktree "gitworktree"
)

func TestParseWorktreeList(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect []gitworktree.Worktree
	}{
		{
			name:   "empty output",
			input:  "",
			expect: nil,
		},
		{
			name: "single worktree",
			input: "worktree /home/user/repo\n" +
				"HEAD abc123def\n" +
				"branch refs/heads/main\n",
			expect: []gitworktree.Worktree{
				{Path: "/home/user/repo", Branch: "main"},
			},
		},
		{
			name: "multiple worktrees",
			input: "worktree /home/user/repo\n" +
				"HEAD abc123def\n" +
				"branch refs/heads/main\n" +
				"\n" +
				"worktree /home/user/repo_gitworktree/feature\n" +
				"HEAD def456abc\n" +
				"branch refs/heads/feature\n",
			expect: []gitworktree.Worktree{
				{Path: "/home/user/repo", Branch: "main"},
				{Path: "/home/user/repo_gitworktree/feature", Branch: "feature"},
			},
		},
		{
			name: "bare repo",
			input: "worktree /home/user/repo.git\n" +
				"HEAD abc123def\n" +
				"bare\n",
			expect: []gitworktree.Worktree{
				{Path: "/home/user/repo.git", IsBare: true},
			},
		},
		{
			name: "detached HEAD",
			input: "worktree /home/user/repo\n" +
				"HEAD abc123def\n" +
				"detached\n",
			expect: []gitworktree.Worktree{
				{Path: "/home/user/repo", Branch: "(detached)"},
			},
		},
		{
			name: "no trailing newline",
			input: "worktree /home/user/repo\n" +
				"HEAD abc123def\n" +
				"branch refs/heads/main",
			expect: []gitworktree.Worktree{
				{Path: "/home/user/repo", Branch: "main"},
			},
		},
		{
			name: "nested branch ref",
			input: "worktree /home/user/repo\n" +
				"HEAD abc123def\n" +
				"branch refs/heads/feature/my-branch\n",
			expect: []gitworktree.Worktree{
				{Path: "/home/user/repo", Branch: "feature/my-branch"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := gitworktree.ParseWorktreeList(tt.input)
			if len(got) != len(tt.expect) {
				t.Fatalf("got %d worktrees, want %d", len(got), len(tt.expect))
			}
			for i, wt := range got {
				if wt.Path != tt.expect[i].Path {
					t.Errorf("worktree[%d].Path = %q, want %q", i, wt.Path, tt.expect[i].Path)
				}
				if wt.Branch != tt.expect[i].Branch {
					t.Errorf("worktree[%d].Branch = %q, want %q", i, wt.Branch, tt.expect[i].Branch)
				}
				if wt.IsBare != tt.expect[i].IsBare {
					t.Errorf("worktree[%d].IsBare = %v, want %v", i, wt.IsBare, tt.expect[i].IsBare)
				}
			}
		})
	}
}

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
