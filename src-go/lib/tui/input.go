package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
)

// ParseInputValue returns the textarea's current value with newlines removed and whitespace trimmed.
func ParseInputValue(input textarea.Model) string {
	return strings.TrimSpace(strings.ReplaceAll(input.Value(), "\n", ""))
}

func NewSearchInput() textinput.Model {
	ti := textinput.New()
	ti.Placeholder = "fuzzy search..."
	ti.CharLimit = 100
	ti.Prompt = "/ "
	return ti
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
