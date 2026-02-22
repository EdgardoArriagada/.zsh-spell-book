package tui

import "github.com/charmbracelet/bubbles/textarea"

func NewInput(placeholder string) textarea.Model {
	ti := textarea.New()
	ti.Placeholder = placeholder
	ti.CharLimit = 100
	ti.ShowLineNumbers = false
	ti.SetHeight(3)
	ti.SetWidth(80)
	ti.Prompt = "  "
	return ti
}
