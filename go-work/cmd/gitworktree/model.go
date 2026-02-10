package main

import (
	"os"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Worktree represents a single git worktree.
type Worktree struct {
	Path   string
	Branch string
	IsBare bool
}

type mode int

const (
	listMode mode = iota
	addMode
	deleteConfirmMode
)

type model struct {
	worktrees []Worktree
	cursor    int
	mode      mode
	input     textinput.Model
	selected  string
	err       error
	current   int // index of current worktree, -1 if none
}

var (
	cursorStyle lipgloss.Style
	activeStyle lipgloss.Style
	currentMark lipgloss.Style
	dimStyle    lipgloss.Style
	branchStyle lipgloss.Style
	titleStyle  lipgloss.Style
	promptStyle lipgloss.Style
	errStyle    lipgloss.Style
)

func init() {
	lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stderr))

	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Bold(true)
	activeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	currentMark = lipgloss.NewStyle().Foreground(lipgloss.Color("114"))
	dimStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	branchStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("39"))
	titleStyle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("99"))
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	errStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
}

func initialModel() model {
	wts, err := listWorktrees()
	ti := textinput.New()
	ti.Placeholder = "branch-name"
	ti.CharLimit = 100
	ti.Width = 40

	cur := -1
	if err == nil {
		cur = currentWorktreeIndex(wts)
	}

	cursor := 0
	if cur > 0 {
		cursor = cur
	}

	return model{
		worktrees: wts,
		cursor:    cursor,
		err:       err,
		input:     ti,
		current:   cur,
	}
}

func (m model) Init() tea.Cmd { return nil }
