package main

import (
	"fmt"
	"os"
	"os/exec"

	tea "charm.land/bubbletea/v2"
)

func tmux(args ...string) {
	exec.Command("tmux", args...).Run() //nolint:errcheck
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithOutput(os.Stderr))
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	mdl := m.(model)
	if mdl.selected == nil {
		os.Exit(1)
	}

	t := mdl.selected
	sessionName := t.SessionID

	if t.Current != os.Getenv("ZSB_CURRENT_TICKET") {
		home := os.Getenv("HOME")
		content := fmt.Sprintf(
			"export ZSB_PARENT_TICKET='%s'\nexport ZSB_CURRENT_TICKET='%s'\nexport ZSB_CURRENT_LABEL='%s'\n",
			t.Parent, t.Current, t.Label,
		)
		os.WriteFile(home+"/temp/current-ticket.zsh", []byte(content), 0644) //nolint:errcheck
	}

	if os.Getenv("TMUX") != "" {
		if exec.Command("tmux", "switch-client", "-t", "="+sessionName).Run() != nil {
			exec.Command("tmux", "new-session", "-d", "-s", sessionName).Run() //nolint:errcheck
			tmux("switch-client", "-t", "="+sessionName)
		}
	} else {
		if exec.Command("tmux", "has-session", "-t", "="+sessionName).Run() != nil {
			exec.Command("tmux", "new-session", "-d", "-s", sessionName).Run() //nolint:errcheck
		}
		tmux("attach-session", "-t", "="+sessionName)
	}
}
