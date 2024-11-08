package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func main() {
	pathToDotGit, err := getGitDir()
	if err != nil {
		return
	}

	if pathToDotGit == ".git" {
		fmt.Println(filepath.Base(getCurrentDir()))
		return
	}

	tail := filepath.Base(pathToDotGit)
	if tail != ".git" {
		fmt.Println(tail)
		return
	}

	parentRootDir := filepath.Dir(pathToDotGit)
	fmt.Println(filepath.Base(parentRootDir))
}

// getCurrentDir returns the path of the current working directory.
func getCurrentDir() string {
	dir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting the current directory:", err)
		os.Exit(1)
	}
	return dir
}

// getGitDir returns the .git directory path of the current git repository.
func getGitDir() (string, error) {
	cmd := exec.Command("git", "rev-parse", "--git-dir")
	cmd.Stderr = nil // silencing errors, similar to '2>/dev/null'

	output, err := cmd.Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}
