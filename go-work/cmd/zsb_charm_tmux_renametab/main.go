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
	if err != nil || args.Len == 0 {
		fmt.Println(err)
		os.Exit(1)
	}

	win_id := args.Args[0]
	currentFilePath := args.Get(1)

	repoRoot, err := git.GetRepoRoot()

	if err != nil {
		os.Exit(1)
	}

	if currentFilePath != "" && !strings.HasPrefix(currentFilePath, repoRoot) {
		os.Exit(1)
	}

	tabname := fmt.Sprintf("  %s", filepath.Base(repoRoot))

	cmd := exec.Command("tmux", "rename-window", "-t", win_id, tabname)

	cmd.Run()
}
