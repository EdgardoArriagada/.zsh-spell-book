package main

import (
	"fmt"
	"os"
	"os/exec"

	"example.com/workspace/lib/argsLib"
	"example.com/workspace/lib/git"
)

func main() {
	args, err := argsLib.Parse()
	if err != nil || len(args) == 0 {
		fmt.Println(err)
		os.Exit(1)
	}

	repo_name, err := git.GetRepoName()

	if err != nil {
		os.Exit(1)
	}

	cmd := exec.Command("tmux", "rename-window", "-t", args[0], fmt.Sprintf(" %s", repo_name))

	err = cmd.Run()
}
