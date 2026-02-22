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
	forceDeleteConfirmMode
)

type model struct {
	worktrees    []Worktree
	cursor       int
	mode         mode
	input        textinput.Model
	selected     string
	fallbackPath string // set when current worktree is deleted, so quit still cd's somewhere valid
	err          error
	statusMsg    string // transient info message shown in list mode
	current      int    // index of current worktree, -1 if none
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
	warnStyle   lipgloss.Style
	statusStyle lipgloss.Style
	keyStyle    lipgloss.Style
)

func init() {
	lipgloss.SetDefaultRenderer(lipgloss.NewRenderer(os.Stderr))

	cursorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EBCB8B")).Bold(true) // yellow — warm cursor
	activeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))             // white1 — hovered item
	currentMark = lipgloss.NewStyle().Foreground(lipgloss.Color("#A3BE8C"))             // green — current worktree
	dimStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#60728A"))             // gray5 — inactive text
	branchStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1"))             // blue1 — branch names
	titleStyle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#88C0D0")) // blue2 — title
	promptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#EBCB8B"))             // yellow — input prompt
	errStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#BF616A"))             // red — errors
	warnStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#D08770"))             // orange — destructive confirm
	statusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#8FBCBB"))             // cyan — info messages
	keyStyle    = lipgloss.NewStyle().Foreground(lipgloss.Color("#EBCB8B")).Bold(true)  // yellow bold — key hints
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
