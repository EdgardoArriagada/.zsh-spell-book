package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	gitlib "example.com/workspace/lib/git"
)

// Worktree is an alias for the shared type in lib/git.
type Worktree = gitlib.Worktree

// WorktreeBaseDir returns the base directory for new worktrees.
// Convention: <parent_of_main>/<main_dirname>_gitworktree
func WorktreeBaseDir(mainWorktreePath string) string {
	parent := filepath.Dir(mainWorktreePath)
	base := filepath.Base(mainWorktreePath) + "_gitworktree"
	return filepath.Join(parent, base)
}

func listWorktrees() ([]Worktree, error) {
	out, err := exec.Command("git", "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil, fmt.Errorf("git worktree list: %w", err)
	}
	return gitlib.ParseWorktreeList(string(out)), nil
}

func createWorktree(mainPath, branch string) error {
	baseDir := WorktreeBaseDir(mainPath)
	wtPath := filepath.Join(baseDir, branch)
	var args []string
	if gitlib.BranchExists(branch) {
		args = []string{"worktree", "add", wtPath, branch}
	} else {
		args = []string{"worktree", "add", wtPath, "-b", branch, "develop"}
	}
	out, err := exec.Command("git", args...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func deleteWorktree(path string) error {
	out, err := exec.Command("git", "worktree", "remove", path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func deleteWorktreeForce(path string) error {
	out, err := exec.Command("git", "worktree", "remove", "--force", path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func isWorktreeDirtyError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "contains modified or untracked files")
}

func currentWorktreeIndex(worktrees []Worktree) int {
	cwd, err := os.Getwd()
	if err != nil {
		return -1
	}
	return gitlib.FindCurrentWorktree(worktrees, cwd)
}
