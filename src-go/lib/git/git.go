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
	// Get remote origin URL
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Stderr = nil // silencing errors, similar to '2>/dev/null'

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	remoteURL := strings.TrimSpace(string(output))

	// Extract base name from URL
	baseName := filepath.Base(remoteURL)

	// Remove .git extension if present
	repoName := strings.TrimSuffix(baseName, ".git")

	return repoName, nil
}
