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
	d := u.Must(args.Parse())
	u.Expect(d.Len == 1 || d.Len == 2, "Usage: zsb_charm_tmux_renametab <win_id> <current_file_path>")

	winId := d.Args[0]
	currentFilePath := d.Get(1)

	repoRoot := u.Must(git.GetRepoRoot())

	if currentFilePath != "" && !strings.HasPrefix(currentFilePath, repoRoot) {
		os.Exit(1)
	}

	tabname := fmt.Sprintf("  %s", filepath.Base(repoRoot))

	cmd := exec.Command("tmux", "rename-window", "-t", winId, tabname)

	cmd.Run()
}
