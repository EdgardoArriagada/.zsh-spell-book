package app

import (
	"fmt"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
)

func tmux(args ...string) { exec.Command("tmux", args...).Run() } //nolint:errcheck

func Run() (bool, error) {
	m, err := tea.NewProgram(initialModel(), tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		return false, err
	}
	t := m.(model).selected
	if t == nil {
		return false, nil
	}
	if t.Current != os.Getenv("ZSB_CURRENT_TICKET") {
		content := fmt.Sprintf("export ZSB_PARENT_TICKET='%s'\nexport ZSB_CURRENT_TICKET='%s'\nexport ZSB_CURRENT_LABEL='%s'\n", t.Parent, t.Current, t.Label)
		os.WriteFile(os.Getenv("HOME")+"/temp/current-ticket.zsh", []byte(content), 0644) //nolint:errcheck
	}
	if os.Getenv("TMUX") != "" {
		if exec.Command("tmux", "switch-client", "-t", "="+t.SessionID).Run() != nil {
			exec.Command("tmux", "new-session", "-d", "-s", t.SessionID).Run() //nolint:errcheck
			tmux("switch-client", "-t", "="+t.SessionID)
		}
	} else {
		if exec.Command("tmux", "has-session", "-t", "="+t.SessionID).Run() != nil {
			exec.Command("tmux", "new-session", "-d", "-s", t.SessionID).Run()
		} //nolint:errcheck
		tmux("attach-session", "-t", "="+t.SessionID)
	}
	return true, nil
}
