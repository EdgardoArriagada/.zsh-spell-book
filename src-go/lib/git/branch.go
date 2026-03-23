package git

import (
	"fmt"
	"os/exec"
)

// ValidateBranchName returns an error if the branch name is rejected by git.
func ValidateBranchName(branch string) error {
	if branch == "" {
		return fmt.Errorf("branch name cannot be empty")
	}
	if err := exec.Command("git", "check-ref-format", "--branch", branch).Run(); err != nil {
		return fmt.Errorf("invalid branch name %q", branch)
	}
	return nil
}

// BranchExists returns true if the given local branch exists.
func BranchExists(branch string) bool {
	err := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch).Run()
	return err == nil
}
