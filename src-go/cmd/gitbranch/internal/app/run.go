package app

import (
	"os"

	tea "charm.land/bubbletea/v2"
)

func Run() (string, error) {
	m, err := tea.NewProgram(initialModel(), tea.WithOutput(os.Stderr)).Run()
	if err != nil {
		return "", err
	}
	return m.(model).selected, nil
}
