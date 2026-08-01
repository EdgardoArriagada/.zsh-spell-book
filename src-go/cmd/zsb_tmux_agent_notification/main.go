package main

import (
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const notifOpt = "@zsb_agent_notif"
const notifSuffix = " 󰂟"

func windowID(pane string) string {
	out, _ := exec.Command("tmux", "display-message", "-p", "-t", pane, "#{window_id}").Output()
	return strings.TrimSpace(string(out))
}

func windowName(pane string) string {
	out, _ := exec.Command("tmux", "display-message", "-p", "-t", pane, "#{window_name}").Output()
	return strings.TrimSpace(string(out))
}

func setWindowName(pane, name string) {
	exec.Command("tmux", "rename-window", "-t", pane, name).Run() //nolint:errcheck
}

func windowHasNotif(winID string) bool {
	out, err := exec.Command("tmux", "list-panes", "-t", winID, "-F", "#{@zsb_agent_notif}").Output()
	if err != nil {
		return false
	}
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if n, _ := strconv.Atoi(strings.TrimSpace(line)); n > 0 {
			return true
		}
	}
	return false
}

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
		if !windowHasNotif(windowID(pane)) {
			name := windowName(pane)
			if strings.HasSuffix(name, notifSuffix) {
				setWindowName(pane, strings.TrimSuffix(name, notifSuffix))
			}
		}
		return
	}

	exec.Command("tmux", "set-option", "-p", "-t", pane, notifOpt, "1").Run() //nolint:errcheck
	name := windowName(pane)
	if !strings.HasSuffix(name, notifSuffix) {
		setWindowName(pane, name+notifSuffix)
	}
}
