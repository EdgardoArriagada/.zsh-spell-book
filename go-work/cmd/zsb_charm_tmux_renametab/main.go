package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
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
	currentFilePath := argsLib.ReadArg(&args, 1)

	repoRoot, err := git.GetRepoRoot()

	if err != nil {
		os.Exit(1)
	}

	if currentFilePath != "" && !strings.HasPrefix(currentFilePath, repoRoot) {
		os.Exit(1)
	}

	repoName := filepath.Base(repoRoot)

	if err != nil {
		os.Exit(1)
	}

	cmd := exec.Command("tmux", "rename-window", "-t", win_id, fmt.Sprintf("  %s", repoName))

	cmd.Run()
}
