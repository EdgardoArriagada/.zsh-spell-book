package main

import (
	"fmt"
	"os/exec"
	"strings"
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

func worktreeBranchSet() map[string]bool {
	out, err := exec.Command("git", "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil
	}
	set := map[string]bool{}
	firstBlock := true
	inBlock := false
	for _, line := range strings.Split(string(out), "\n") {
		if strings.HasPrefix(line, "worktree ") {
			if firstBlock {
				firstBlock = false
				inBlock = true
				continue
			}
			inBlock = true
			continue
		}
		if line == "" {
			inBlock = false
			continue
		}
		if !firstBlock && inBlock && strings.HasPrefix(line, "branch refs/heads/") {
			branch := strings.TrimPrefix(line, "branch refs/heads/")
			set[branch] = true
		}
	}
	return set
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

func listBranches() ([]Branch, error) {
	out, err := exec.Command("git", "branch", "--format=%(refname:short) %(HEAD)").Output()
	if err != nil {
		return nil, fmt.Errorf("git branch: %w", err)
	}
	worktrees := worktreeBranchSet()
	var branches []Branch
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		name := parts[0]
		isCurrent := len(parts) > 1 && parts[1] == "*"
		isWorktree := !isCurrent && worktrees[name]
		branches = append(branches, Branch{Name: name, IsCurrent: isCurrent, IsWorktree: isWorktree})
	}
	return sortBranches(branches), nil
}

func checkoutBranch(name string) error {
	out, err := exec.Command("git", "checkout", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func createBranch(name string) error {
	out, err := exec.Command("git", "branch", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func deleteBranch(name string) error {
	out, err := exec.Command("git", "branch", "-d", name).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func forceDeleteBranch(name string) error {
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
