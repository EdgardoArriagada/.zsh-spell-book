package git

import (
	"fmt"
	"os/exec"
	"strings"
)

// ValidateBranchName returns an error if the branch name contains characters
// that git-check-ref-format would reject.
func ValidateBranchName(branch string) error {
	if branch == "" {
		return fmt.Errorf("branch name cannot be empty")
	}
	banned := []string{" ", "~", "^", ":", "?", "*", "[", "\\", "..", "@{"}
	for _, b := range banned {
		if strings.Contains(branch, b) {
			return fmt.Errorf("invalid branch name: contains %q", b)
		}
	}
	if strings.HasPrefix(branch, "-") || strings.HasPrefix(branch, ".") {
		return fmt.Errorf("invalid branch name: cannot start with %q", string(branch[0]))
	}
	if strings.HasSuffix(branch, ".") || strings.HasSuffix(branch, "/") || strings.HasSuffix(branch, ".lock") {
		return fmt.Errorf("invalid branch name: invalid suffix")
	}
	return nil
}

// BranchExists returns true if the given local branch exists.
func BranchExists(branch string) bool {
	err := exec.Command("git", "show-ref", "--verify", "--quiet", "refs/heads/"+branch).Run()
	return err == nil
}
