package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	gitlib "example.com/workspace/lib/git"
)

var priorityBranches = []string{"develop", "master", "main"}

func isDefaultBranch(name string) bool {
	for _, d := range priorityBranches {
		if name == d {
			return true
		}
	}
	return false
}

func sortBranches(branches []Branch) []Branch {
	priority := map[string]int{}
	for i, name := range priorityBranches {
		priority[name] = i
	}
	defaults := make([]*Branch, len(priorityBranches))
	var regular []Branch
	var worktrees []Branch
	for i := range branches {
		if branches[i].IsWorktree {
			worktrees = append(worktrees, branches[i])
		} else if idx, ok := priority[branches[i].Name]; ok {
			defaults[idx] = &branches[i]
		} else {
			regular = append(regular, branches[i])
		}
	}
	result := make([]Branch, 0, len(branches))
	for _, b := range defaults {
		if b != nil {
			result = append(result, *b)
		}
	}
	result = append(result, regular...)
	return append(result, worktrees...)
}

// listBranchesWithWorktrees parses branches using pre-fetched worktree data.
func listBranchesWithWorktrees(wts []gitlib.Worktree) ([]Branch, error) {
	out, err := exec.Command("git", "branch", "--format=%(refname:short) %(HEAD)").Output()
	if err != nil {
		return nil, fmt.Errorf("git branch: %w", err)
	}
	worktreeSet := map[string]bool{}
	for i := 1; i < len(wts); i++ {
		if wts[i].Branch != "" && wts[i].Branch != "(detached)" {
			worktreeSet[wts[i].Branch] = true
		}
	}
	var branches []Branch
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		name := parts[0]
		isCurrent := len(parts) > 1 && parts[1] == "*"
		isWorktree := !isCurrent && worktreeSet[name]
		branches = append(branches, Branch{Name: name, IsCurrent: isCurrent, IsWorktree: isWorktree})
	}
	return sortBranches(branches), nil
}

var listBranches = func() ([]Branch, error) {
	out, _ := exec.Command("git", "worktree", "list", "--porcelain").Output()
	return listBranchesWithWorktrees(gitlib.ParseWorktreeList(string(out)))
}

var checkoutBranch = func(name string) error {
	out, err := exec.Command("git", "checkout", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

var createBranch = func(name string) error {
	out, err := exec.Command("git", "branch", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

var deleteBranch = func(name string) error {
	out, err := exec.Command("git", "branch", "-d", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

var forceDeleteBranch = func(name string) error {
	out, err := exec.Command("git", "branch", "-D", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func isUnmergedBranchError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "is not fully merged")
}

// initBranchData fetches all startup data with a single git worktree subprocess call.
func initBranchData() (branches []Branch, inWorktree bool, err error) {
	wtOut, _ := exec.Command("git", "worktree", "list", "--porcelain").Output()
	wts := gitlib.ParseWorktreeList(string(wtOut))
	branches, err = listBranchesWithWorktrees(wts)
	if err != nil {
		return
	}
	cwd, cwdErr := os.Getwd()
	if cwdErr == nil {
		inWorktree = gitlib.FindCurrentWorktree(wts, cwd) > 0
	}
	return
}
