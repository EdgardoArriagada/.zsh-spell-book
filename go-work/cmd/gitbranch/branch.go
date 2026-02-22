package main

import (
	"fmt"
	"os/exec"
	"strings"
)

var priorityBranches = []string{"master", "main", "develop"}

func sortBranches(branches []Branch) []Branch {
	priority := map[string]int{}
	for i, name := range priorityBranches {
		priority[name] = i
	}
	defaults := make([]*Branch, len(priorityBranches))
	var rest []Branch
	for i := range branches {
		if idx, ok := priority[branches[i].Name]; ok {
			defaults[idx] = &branches[i]
		} else {
			rest = append(rest, branches[i])
		}
	}
	result := make([]Branch, 0, len(branches))
	for _, b := range defaults {
		if b != nil {
			result = append(result, *b)
		}
	}
	return append(result, rest...)
}

func listBranches() ([]Branch, error) {
	out, err := exec.Command("git", "branch", "--format=%(refname:short) %(HEAD)").Output()
	if err != nil {
		return nil, fmt.Errorf("git branch: %w", err)
	}
	var branches []Branch
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, " ", 2)
		name := parts[0]
		isCurrent := len(parts) > 1 && parts[1] == "*"
		branches = append(branches, Branch{Name: name, IsCurrent: isCurrent})
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
