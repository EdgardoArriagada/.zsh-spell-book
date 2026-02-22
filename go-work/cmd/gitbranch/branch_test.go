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
