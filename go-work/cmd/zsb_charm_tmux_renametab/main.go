package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"example.com/workspace/lib/args"
	"example.com/workspace/lib/git"
	u "example.com/workspace/lib/utils"
)

func main() {
	d := u.Must1(args.Parse())
	u.Expect(d.Len == 2, "Usage: zsb_charm_tmux_renametab <win_id> <current_file_path>")

	repoRoot := u.Must1(git.GetRepoRoot())

	winId := d.Args[0]
	currentFilePath := d.Args[1]

	// if a file is given when opening neovim, make sure it is in the repo
	if currentFilePath != "" && !strings.HasPrefix(currentFilePath, repoRoot) {
		os.Exit(1)
	}

	tabname := fmt.Sprintf("  %s", filepath.Base(repoRoot))

	cmd := exec.Command("tmux", "rename-window", "-t", winId, tabname)

	cmd.Run()
}
