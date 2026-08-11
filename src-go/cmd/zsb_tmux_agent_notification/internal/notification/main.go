package notification

import (
	"os/exec"
	"strings"
)

const NOTIF_VAR = "@zsb_agent_notif"
const FINISHED_SUFFIX = " 󰂚" // bell — agent finished / needs attention
const WORKING_SUFFIX = " 󰔟"  // hourglass — agent still working
const MANUAL_SUFFIX = " 󰹇"   // flag — manually flagged

const (
	CLEAR    = "0"
	FINISHED = "1"
	WORKING  = "2"
	MANUAL   = "3"
)

func paneNotif(pane string) string {
	out, _ := exec.Command("tmux", "show-options", "-pqv", "-t", pane, NOTIF_VAR).Output()
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
		case strings.HasSuffix(name, FINISHED_SUFFIX):
			name = strings.TrimSuffix(name, FINISHED_SUFFIX)
		case strings.HasSuffix(name, WORKING_SUFFIX):
			name = strings.TrimSuffix(name, WORKING_SUFFIX)
		case strings.HasSuffix(name, MANUAL_SUFFIX):
			name = strings.TrimSuffix(name, MANUAL_SUFFIX)
		default:
			return name
		}
	}
}

// windowStates reports whether any pane in the window is finished (1), working (2), or manual (3).
func windowStates(winID string) (finished, working, manual bool) {
	out, err := exec.Command("tmux", "list-panes", "-t", winID, "-F", "#{@zsb_agent_notif}").Output()
	if err != nil {
		return false, false, false
	}
	for line := range strings.SplitSeq(strings.TrimSpace(string(out)), "\n") {
		switch strings.TrimSpace(line) {
		case FINISHED:
			finished = true
		case WORKING:
			working = true
		case MANUAL:
			manual = true
		}
	}
	return finished, working, manual
}

// refreshWindowName recomputes the window-name suffixes from the window's pane
// states: working glyph first, then finished bell. Renames only if changed.
func refreshWindowName(pane string) {
	finished, working, manual := windowStates(windowID(pane))
	cur := windowName(pane)
	want := baseName(cur)
	if working {
		want += WORKING_SUFFIX
	}
	if finished {
		want += FINISHED_SUFFIX
	}
	if manual {
		want += MANUAL_SUFFIX
	}
	if want != cur {
		exec.Command("tmux", "rename-window", "-t", pane, want).Run() //nolint:errcheck
	}
}

func setPaneState(pane string, state string) {
	exec.Command("tmux", "set-option", "-p", "-t", pane, NOTIF_VAR, state).Run() //nolint:errcheck
}

// zsb_tmux_agent_notification [--finished|--working|--clear-finished] <session_name> <pane_id>
// pane_id identifies the pane (and its session); session_name matches the hook
// signature but is unused. No flag defaults to --finished.
func Run(args []string) int {
	flag := ""
	if len(args) > 0 && strings.HasPrefix(args[0], "--") {
		flag = args[0]
		args = args[1:]
	}
	if len(args) < 2 {
		return 1
	}
	pane := args[1]

	switch flag {
	case "", "--finished":
		// Already watching: clear working state without ringing the bell.
		if paneIsFocused(pane) {
			setPaneState(pane, CLEAR)
		} else {
			setPaneState(pane, FINISHED)
		}
	case "--force-finished":
		setPaneState(pane, FINISHED)
	case "--working":
		// working fires while the pane is still focused (at submit time); don't
		// skip on focus — the focus-out/next --clear resets it.
		setPaneState(pane, WORKING)
	case "--clear-finished":
		// Only clear the finished (1) state; leave working (2) and manual (3) alone.
		if paneNotif(pane) == FINISHED {
			setPaneState(pane, CLEAR)
		}
	case "--manual":
		// Toggle: 3 → 0, 0 → 3. Ignore if pane is working (2) or finished (1).
		switch paneNotif(pane) {
		case MANUAL:
			setPaneState(pane, CLEAR)
		case CLEAR, "":
			setPaneState(pane, MANUAL)
		}
	default:
		return 1
	}
	refreshWindowName(pane)
	return 0
}
