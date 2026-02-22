package main

import (
	"fmt"
	"os/exec"
	"strings"
)

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
	return branches, nil
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
