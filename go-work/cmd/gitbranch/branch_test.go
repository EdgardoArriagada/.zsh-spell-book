package main

import (
	"fmt"
	"testing"
)

func TestParseWorktreeBranches(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		expect map[string]bool
	}{
		{
			name:   "empty output",
			input:  "",
			expect: map[string]bool{},
		},
		{
			name: "only main worktree",
			input: "worktree /repo\nHEAD abc123\nbranch refs/heads/main\n\n",
			expect: map[string]bool{},
		},
		{
			name: "main + one linked worktree",
			input: "worktree /repo\nHEAD abc123\nbranch refs/heads/main\n\nworktree /repo/feat\nHEAD def456\nbranch refs/heads/feature\n\n",
			expect: map[string]bool{"feature": true},
		},
		{
			name: "main + two linked worktrees",
			input: "worktree /repo\nHEAD abc\nbranch refs/heads/main\n\nworktree /repo/a\nHEAD bbb\nbranch refs/heads/feat-a\n\nworktree /repo/b\nHEAD ccc\nbranch refs/heads/feat-b\n\n",
			expect: map[string]bool{"feat-a": true, "feat-b": true},
		},
		{
			name: "linked worktree in detached HEAD state",
			input: "worktree /repo\nHEAD abc\nbranch refs/heads/main\n\nworktree /repo/detached\nHEAD def\ndetached\n\n",
			expect: map[string]bool{},
		},
		{
			name: "main detached + linked worktree with branch",
			input: "worktree /repo\nHEAD abc\ndetached\n\nworktree /repo/feat\nHEAD def\nbranch refs/heads/feat\n\n",
			expect: map[string]bool{"feat": true},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseWorktreeBranches(tt.input)
			if len(got) != len(tt.expect) {
				t.Errorf("len = %d, want %d; got %v", len(got), len(tt.expect), got)
				return
			}
			for k := range tt.expect {
				if !got[k] {
					t.Errorf("missing key %q in result %v", k, got)
				}
			}
		})
	}
}

func TestSortBranches(t *testing.T) {
	tests := []struct {
		name   string
		input  []Branch
		expect []string // expected order of branch names
	}{
		{
			name:   "empty input",
			input:  []Branch{},
			expect: []string{},
		},
		{
			name: "only regular branches",
			input: []Branch{
				{Name: "feat-a"},
				{Name: "fix-b"},
			},
			expect: []string{"feat-a", "fix-b"},
		},
		{
			name: "priority branches appear first",
			input: []Branch{
				{Name: "feat"},
				{Name: "main"},
				{Name: "develop"},
			},
			expect: []string{"develop", "main", "feat"},
		},
		{
			name: "master before regular",
			input: []Branch{
				{Name: "feat"},
				{Name: "master"},
			},
			expect: []string{"master", "feat"},
		},
		{
			name: "priority order: develop < master < main",
			input: []Branch{
				{Name: "main"},
				{Name: "master"},
				{Name: "develop"},
			},
			expect: []string{"develop", "master", "main"},
		},
		{
			name: "worktrees sorted last",
			input: []Branch{
				{Name: "feat"},
				{Name: "wt-branch", IsWorktree: true},
				{Name: "main"},
			},
			expect: []string{"main", "feat", "wt-branch"},
		},
		{
			name: "priority → regular → worktrees",
			input: []Branch{
				{Name: "regular-b"},
				{Name: "wt-one", IsWorktree: true},
				{Name: "main"},
				{Name: "regular-a"},
				{Name: "wt-two", IsWorktree: true},
			},
			expect: []string{"main", "regular-b", "regular-a", "wt-one", "wt-two"},
		},
		{
			name: "missing priority branches are skipped",
			input: []Branch{
				{Name: "feat"},
				{Name: "main"},
			},
			expect: []string{"main", "feat"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := sortBranches(tt.input)
			if len(got) != len(tt.expect) {
				t.Fatalf("len = %d, want %d; got %v", len(got), len(tt.expect), got)
			}
			for i, name := range tt.expect {
				if got[i].Name != name {
					t.Errorf("position %d: got %q, want %q (full order: %v)", i, got[i].Name, name, got)
				}
			}
		})
	}
}

func TestParseWorktreePaths(t *testing.T) {
	tests := []struct {
		name       string
		input      string
		wantMain   string
		wantLinked []string
	}{
		{
			name:       "empty output",
			input:      "",
			wantMain:   "",
			wantLinked: nil,
		},
		{
			name:       "only main worktree",
			input:      "worktree /repo\nHEAD abc123\nbranch refs/heads/main\n\n",
			wantMain:   "/repo",
			wantLinked: nil,
		},
		{
			name:       "main + one linked worktree",
			input:      "worktree /repo\nHEAD abc123\nbranch refs/heads/main\n\nworktree /repo/feat\nHEAD def456\nbranch refs/heads/feature\n\n",
			wantMain:   "/repo",
			wantLinked: []string{"/repo/feat"},
		},
		{
			name:       "main + two linked worktrees",
			input:      "worktree /repo\nHEAD abc\nbranch refs/heads/main\n\nworktree /repo/a\nHEAD bbb\nbranch refs/heads/feat-a\n\nworktree /repo/b\nHEAD ccc\nbranch refs/heads/feat-b\n\n",
			wantMain:   "/repo",
			wantLinked: []string{"/repo/a", "/repo/b"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotMain, gotLinked := parseWorktreePaths(tt.input)
			if gotMain != tt.wantMain {
				t.Errorf("main = %q, want %q", gotMain, tt.wantMain)
			}
			if len(gotLinked) != len(tt.wantLinked) {
				t.Errorf("linked len = %d, want %d; got %v", len(gotLinked), len(tt.wantLinked), gotLinked)
				return
			}
			for i, p := range tt.wantLinked {
				if gotLinked[i] != p {
					t.Errorf("linked[%d] = %q, want %q", i, gotLinked[i], p)
				}
			}
		})
	}
}

func TestIsLinkedWorktreeIn(t *testing.T) {
	output := "worktree /repo\nHEAD abc\nbranch refs/heads/main\n\nworktree /repo/feat\nHEAD def\nbranch refs/heads/feature\n\n"
	tests := []struct {
		name   string
		cwd    string
		expect bool
	}{
		{"cwd is main worktree", "/repo", false},
		{"cwd is subdirectory of main", "/repo/src", false},
		{"cwd is linked worktree", "/repo/feat", true},
		{"cwd is subdirectory of linked worktree", "/repo/feat/subdir", true},
		{"cwd is outside all worktrees", "/other", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := isLinkedWorktreeIn(output, tt.cwd)
			if got != tt.expect {
				t.Errorf("isLinkedWorktreeIn(%q) = %v, want %v", tt.cwd, got, tt.expect)
			}
		})
	}
}

func TestIsUnmergedBranchError(t *testing.T) {
	tests := []struct {
		err    error
		expect bool
	}{
		{nil, false},
		{fmt.Errorf("some other error"), false},
		{fmt.Errorf("error: The branch 'feat' is not fully merged"), true},
		{fmt.Errorf("is not fully merged"), true},
	}
	for _, tt := range tests {
		got := isUnmergedBranchError(tt.err)
		if got != tt.expect {
			t.Errorf("isUnmergedBranchError(%v) = %v, want %v", tt.err, got, tt.expect)
		}
	}
}
