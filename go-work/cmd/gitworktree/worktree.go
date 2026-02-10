package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// ParseWorktreeList parses the porcelain output of `git worktree list --porcelain`.
func ParseWorktreeList(output string) []Worktree {
	var result []Worktree
	var wt Worktree
	inBlock := false

	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "worktree ") {
			if inBlock {
				result = append(result, wt)
			}
			wt = Worktree{Path: strings.TrimPrefix(line, "worktree ")}
			inBlock = true
		} else if strings.HasPrefix(line, "branch ") {
			wt.Branch = strings.TrimPrefix(strings.TrimPrefix(line, "branch "), "refs/heads/")
		} else if line == "bare" {
			wt.IsBare = true
		} else if line == "detached" {
			wt.Branch = "(detached)"
		} else if line == "" && inBlock {
			result = append(result, wt)
			wt = Worktree{}
			inBlock = false
		}
	}
	if inBlock {
		result = append(result, wt)
	}
	return result
}

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
	return ParseWorktreeList(string(out)), nil
}

func createWorktree(mainPath, branch string) error {
	baseDir := WorktreeBaseDir(mainPath)
	wtPath := filepath.Join(baseDir, branch)
	out, err := exec.Command("git", "worktree", "add", wtPath, "-b", branch, "develop").CombinedOutput()
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

func currentWorktreeIndex(worktrees []Worktree) int {
	cwd, err := os.Getwd()
	if err != nil {
		return -1
	}
	return FindCurrentWorktree(worktrees, cwd)
}

func FindCurrentWorktree(worktrees []Worktree, cwd string) int {
	best := -1
	bestLen := 0
	for i, wt := range worktrees {
		if strings.HasPrefix(cwd, wt.Path) && len(wt.Path) > bestLen {
			best = i
			bestLen = len(wt.Path)
		}
	}
	return best
}
