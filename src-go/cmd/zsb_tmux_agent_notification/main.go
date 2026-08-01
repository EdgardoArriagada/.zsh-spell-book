package main

import (
	"os"
	"os/exec"
)

const notifOpt = "@zsb_agent_notif"

// zsb_tmux_agent_notification [--decrease] <session_name> <pane_id>
// pane_id identifies the pane (and its session); session_name matches the
// hook signature but is unused.
func main() {
	args := os.Args[1:]
	decrease := false
	if len(args) > 0 && args[0] == "--decrease" {
		decrease = true
		args = args[1:]
	}
	if len(args) < 2 {
		os.Exit(1)
	}
	pane := args[1]

	if decrease {
		exec.Command("tmux", "set-option", "-p", "-t", pane, notifOpt, "0").Run() //nolint:errcheck
		return
	}

	exec.Command("tmux", "set-option", "-p", "-t", pane, notifOpt, "1").Run() //nolint:errcheck
}
