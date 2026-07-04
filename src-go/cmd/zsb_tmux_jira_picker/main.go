package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func tmux(args ...string) {
	exec.Command("tmux", args...).Run() //nolint:errcheck
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithOutput(os.Stderr))
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
	labelKebab := strings.ToLower(strings.ReplaceAll(t.Label, " ", "-"))
	if len(labelKebab) > 25 {
		labelKebab = labelKebab[:25]
	}
	sessionName := t.Current + "-" + labelKebab

	sessionExists := exec.Command("tmux", "has-session", "-t", "="+sessionName).Run() == nil

	if !sessionExists {
		exec.Command("tmux", "new-session", "-d", "-s", sessionName).Run() //nolint:errcheck
		tmux("set-environment", "-t", sessionName, "ZSB_PARENT_TICKET", t.Parent)
		tmux("set-environment", "-t", sessionName, "ZSB_CURRENT_TICKET", t.Current)
		tmux("set-environment", "-t", sessionName, "ZSB_CURRENT_LABEL", t.Label)
	}

	home := os.Getenv("HOME")
	content := fmt.Sprintf(
		"declare ZSB_PARENT_TICKET='%s'\ndeclare ZSB_CURRENT_TICKET='%s'\ndeclare ZSB_CURRENT_LABEL='%s'\n",
		t.Parent, t.Current, t.Label,
	)
	os.WriteFile(home+"/temp/current-ticket.zsh", []byte(content), 0644) //nolint:errcheck

	if os.Getenv("TMUX") != "" {
		tmux("switch-client", "-t", "="+sessionName)
	} else {
		tmux("attach-session", "-t", "="+sessionName)
	}
}
