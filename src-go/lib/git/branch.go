package git

import (
	"fmt"
	"os/exec"
	"strings"
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

// RemoteBranchExists returns true if any remote-tracking branch matches the
// given branch name (e.g. refs/remotes/origin/<branch>).
func RemoteBranchExists(branch string) bool {
	out, err := exec.Command("git", "for-each-ref", "--format=%(refname:short)", "refs/remotes").Output()
	if err != nil {
		return false
	}
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		if _, rest, ok := strings.Cut(line, "/"); ok && rest == branch {
			return true
		}
	}
	return false
}
