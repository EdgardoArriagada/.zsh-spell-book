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

// RemoteBranchRef returns the full remote-tracking ref for branch.
func RemoteBranchRef(branch string) (string, bool) {
	out, err := exec.Command("git", "for-each-ref", "--format=%(refname)", "refs/remotes").Output()
	if err != nil {
		return "", false
	}
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		const prefix = "refs/remotes/"
		if remoteBranch, ok := strings.CutPrefix(line, prefix); ok {
			if _, rest, ok := strings.Cut(remoteBranch, "/"); ok && rest == branch {
				return line, true
			}
		}
	}
	return "", false
}

// RemoteBranchExists returns true if any remote-tracking branch matches the
// given branch name (e.g. refs/remotes/origin/<branch>).
func RemoteBranchExists(branch string) bool {
	_, ok := RemoteBranchRef(branch)
	return ok
}
