package main

import (
	"os"
	"os/exec"
	"strings"
)

const notifOpt = "@zsb_agent_notif"
const finishedSuffix = " 󰂚" // bell — agent finished / needs attention
const workingSuffix = " 󰒓" // cog — agent still working

func paneNotif(pane string) string {
	out, _ := exec.Command("tmux", "show-options", "-pqv", "-t", pane, notifOpt).Output()
	return strings.TrimSpace(string(out))
}

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

// paneIsFocused reports whether the user is currently looking at this pane:
// its session is attached, its window is active, and the pane is active.
func paneIsFocused(pane string) bool {
	out, err := exec.Command("tmux", "display-message", "-p", "-t", pane,
		"#{&&:#{session_attached},#{&&:#{window_active},#{pane_active}}}").Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "1"
}

// baseName strips any trailing known notif suffixes (in either order).
func baseName(name string) string {
	for {
		switch {
		case strings.HasSuffix(name, finishedSuffix):
			name = strings.TrimSuffix(name, finishedSuffix)
		case strings.HasSuffix(name, workingSuffix):
			name = strings.TrimSuffix(name, workingSuffix)
		default:
			return name
		}
	}
}

// windowStates reports whether any pane in the window is finished (1) or working (2).
func windowStates(winID string) (finished, working bool) {
	out, err := exec.Command("tmux", "list-panes", "-t", winID, "-F", "#{@zsb_agent_notif}").Output()
	if err != nil {
		return false, false
	}
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		switch strings.TrimSpace(line) {
		case "1":
			finished = true
		case "2":
			working = true
		}
	}
	return finished, working
}

// refreshWindowName recomputes the window-name suffixes from the window's pane
// states: working glyph first, then finished bell. Renames only if changed.
func refreshWindowName(pane string) {
	finished, working := windowStates(windowID(pane))
	cur := windowName(pane)
	want := baseName(cur)
	if working {
		want += workingSuffix
	}
	if finished {
		want += finishedSuffix
	}
	if want != cur {
		setWindowName(pane, want)
	}
}

// zsb_tmux_agent_notification [--finished|--working|--clear-finished] <session_name> <pane_id>
// pane_id identifies the pane (and its session); session_name matches the hook
// signature but is unused. No flag defaults to --finished.
func main() {
	args := os.Args[1:]
	mode := "finished"
	if len(args) > 0 && strings.HasPrefix(args[0], "--") {
		switch args[0] {
		case "--finished":
			mode = "finished"
		case "--working":
			mode = "working"
		case "--clear-finished":
			mode = "clear-finished"
		default:
			os.Exit(1)
		}
		args = args[1:]
	}
	if len(args) < 2 {
		os.Exit(1)
	}
	pane := args[1]

	switch mode {
	case "clear-finished":
		// Only clear the finished (1) state; leave a working (2) pane alone so
		// the blue "still working" badge persists until the agent finishes.
		if paneNotif(pane) == "1" {
			exec.Command("tmux", "set-option", "-p", "-t", pane, notifOpt, "0").Run() //nolint:errcheck
		}
	case "working":
		// working fires while the pane is still focused (at submit time); don't
		// skip on focus — the focus-out/next --clear resets it.
		exec.Command("tmux", "set-option", "-p", "-t", pane, notifOpt, "2").Run() //nolint:errcheck
	default: // finished
		if paneIsFocused(pane) {
			return
		}
		exec.Command("tmux", "set-option", "-p", "-t", pane, notifOpt, "1").Run() //nolint:errcheck
	}
	refreshWindowName(pane)
}
