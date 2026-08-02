package tui

import (
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

var (
	CursorStyle   lipgloss.Style
	ActiveStyle   lipgloss.Style
	CurrentMark   lipgloss.Style
	DimStyle      lipgloss.Style
	BranchStyle   lipgloss.Style
	WorktreeStyle lipgloss.Style
	TitleStyle    lipgloss.Style
	PromptStyle   lipgloss.Style
	ErrStyle      lipgloss.Style
	WarnStyle     lipgloss.Style
	StatusStyle   lipgloss.Style
	KeyStyle      lipgloss.Style
	NotifBadge    lipgloss.Style
	NotifWorking  lipgloss.Style
)

func init() {
	r := lipgloss.NewRenderer(os.Stderr)
	// ponytail: pin profile/background so termenv skips its terminal-capability
	// probe (env parse + possible OSC round-trip) on the keybind-hot boot path.
	// This stack is truecolor + dark (COLORTERM=truecolor, tmux Tc). All styles
	// here are fixed hex, so no adaptive-color detection is lost.
	r.SetColorProfile(termenv.TrueColor)
	r.SetHasDarkBackground(true)
	lipgloss.SetDefaultRenderer(r)

	CursorStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#EBCB8B")).Bold(true) // yellow — warm cursor
	ActiveStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#D8DEE9"))             // white1 — hovered item
	CurrentMark   = lipgloss.NewStyle().Foreground(lipgloss.Color("#A3BE8C"))             // green — current marker
	DimStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#60728A"))             // gray5 — inactive text
	BranchStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1"))             // blue1 — branch names
	WorktreeStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#4C566A"))             // gray4 — worktree branches (darker)
	TitleStyle    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#88C0D0"))  // blue2 — title
	PromptStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#EBCB8B"))             // yellow — input prompt
	ErrStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#BF616A"))             // red — errors
	WarnStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("#D08770"))             // orange — destructive confirm
	StatusStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("#8FBCBB"))             // cyan — info messages
	KeyStyle      = lipgloss.NewStyle().Foreground(lipgloss.Color("#EBCB8B")).Bold(true)  // yellow bold — key hints
	NotifBadge    = lipgloss.NewStyle().Foreground(lipgloss.Color("#BF616A")).Bold(true)  // red — finished notifications
	NotifWorking  = lipgloss.NewStyle().Foreground(lipgloss.Color("#81A1C1")).Bold(true)  // blue1 — agents still working
}

func Title(text string) string {
	return "\n" + TitleStyle.Render("  "+text) + "\n\n"
}

func Hint(key, desc string) string {
	return KeyStyle.Render(key) + DimStyle.Render(" "+desc)
}

func Sep() string {
	return DimStyle.Render("  ·  ")
}
