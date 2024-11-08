package git

import (
	"os/exec"
	"path/filepath"
	"strings"
)

func GetRepoRoot() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--show-toplevel")
	cmd.Stderr = nil // silencing errors, similar to '2>/dev/null'

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func GetRepoName() (string, error) {
	repoRoot, err := GetRepoRoot()
	if err != nil {
		return "", err
	}

	return filepath.Base(repoRoot), nil
}
