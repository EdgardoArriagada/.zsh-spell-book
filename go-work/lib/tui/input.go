package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
)

// ParseInputValue returns the textarea's current value with newlines removed and whitespace trimmed.
func ParseInputValue(input textarea.Model) string {
	return strings.TrimSpace(strings.ReplaceAll(input.Value(), "\n", ""))
}

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
