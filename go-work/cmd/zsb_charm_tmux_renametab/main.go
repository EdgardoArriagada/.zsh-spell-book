package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"example.com/workspace/lib/argsLib"
	"example.com/workspace/lib/git"
)

func main() {
	args, err := argsLib.Parse()
	if err != nil || len(args) == 0 {
		fmt.Println(err)
		os.Exit(1)
	}

	win_id := args[0]

	var currentFilePath string
	if len(args) > 1 {
		currentFilePath = args[1]
	}

	if currentFilePath != "" {
		repoRoot, err := git.GetRepoRoot()
		if err != nil {
			fmt.Println("Error getting git repo root:", err)
			os.Exit(1)
		}

		if !strings.HasPrefix(currentFilePath, repoRoot) {
			fmt.Println("Current file is not within the git repository")
			os.Exit(1)
		}
	}

	repoName, err := git.GetRepoName()

	if err != nil {
		os.Exit(1)
	}

	cmd := exec.Command("tmux", "rename-window", "-t", win_id, fmt.Sprintf(" %s", repoName))

	cmd.Run()
}
