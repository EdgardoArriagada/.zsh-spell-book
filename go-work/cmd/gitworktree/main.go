package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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

// ParseWorktreeList parses the porcelain output of `git worktree list --porcelain`.
func ParseWorktreeList(output string) []Worktree {
	var result []Worktree
	var wt Worktree
	inBlock := false

	for _, line := range strings.Split(output, "\n") {
		if strings.HasPrefix(line, "worktree ") {
			if inBlock {
				result = append(result, wt)
			}
			wt = Worktree{Path: strings.TrimPrefix(line, "worktree ")}
			inBlock = true
		} else if strings.HasPrefix(line, "branch ") {
			wt.Branch = strings.TrimPrefix(strings.TrimPrefix(line, "branch "), "refs/heads/")
		} else if line == "bare" {
			wt.IsBare = true
		} else if line == "detached" {
			wt.Branch = "(detached)"
		} else if line == "" && inBlock {
			result = append(result, wt)
			wt = Worktree{}
			inBlock = false
		}
	}
	if inBlock {
		result = append(result, wt)
	}
	return result
}

// WorktreeBaseDir returns the base directory for new worktrees.
// Convention: <parent_of_main>/<main_dirname>_gitworktree
func WorktreeBaseDir(mainWorktreePath string) string {
	parent := filepath.Dir(mainWorktreePath)
	base := filepath.Base(mainWorktreePath) + "_gitworktree"
	return filepath.Join(parent, base)
}

func listWorktrees() ([]Worktree, error) {
	out, err := exec.Command("git", "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil, fmt.Errorf("git worktree list: %w", err)
	}
	return ParseWorktreeList(string(out)), nil
}

func createWorktree(mainPath, branch string) error {
	baseDir := WorktreeBaseDir(mainPath)
	wtPath := filepath.Join(baseDir, branch)
	out, err := exec.Command("git", "worktree", "add", wtPath, "-b", branch).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func deleteWorktree(path string) error {
	out, err := exec.Command("git", "worktree", "remove", path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func currentWorktreeIndex(worktrees []Worktree) int {
	cwd, err := os.Getwd()
	if err != nil {
		return -1
	}
	for i, wt := range worktrees {
		if strings.HasPrefix(cwd, wt.Path) {
			return i
		}
	}
	return -1
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

	return model{
		worktrees: wts,
		err:       err,
		input:     ti,
		current:   cur,
	}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.mode {
	case addMode:
		return m.updateAdd(msg)
	case deleteConfirmMode:
		return m.updateDelete(msg)
	default:
		return m.updateList(msg)
	}
}

func (m model) updateList(msg tea.Msg) (tea.Model, tea.Cmd) {
	km, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	switch km.String() {
	case "q", "ctrl+c", "esc":
		return m, tea.Quit
	case "j", "down", "tab":
		if m.cursor < len(m.worktrees)-1 {
			m.cursor++
		}
	case "k", "up", "shift+tab":
		if m.cursor > 0 {
			m.cursor--
		}
	case "enter":
		if len(m.worktrees) > 0 {
			m.selected = m.worktrees[m.cursor].Path
		}
		return m, tea.Quit
	case "a":
		m.mode = addMode
		m.input.SetValue("")
		m.err = nil
		return m, m.input.Focus()
	case "d":
		if len(m.worktrees) > 0 && m.cursor != 0 {
			m.mode = deleteConfirmMode
		}
	}
	return m, nil
}

func (m model) updateAdd(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.mode = listMode
			m.input.Blur()
			m.err = nil
			return m, nil
		case "enter":
			branch := strings.TrimSpace(m.input.Value())
			if branch == "" {
				return m, nil
			}
			mainPath := m.worktrees[0].Path
			if err := createWorktree(mainPath, branch); err != nil {
				m.err = err
				return m, nil
			}
			wts, err := listWorktrees()
			if err != nil {
				m.err = err
				return m, nil
			}
			m.worktrees = wts
			m.current = currentWorktreeIndex(wts)
			m.mode = listMode
			m.input.Blur()
			m.err = nil
			for i, wt := range wts {
				if wt.Branch == branch {
					m.cursor = i
					break
				}
			}
			return m, nil
		}
	}
	var cmd tea.Cmd
	m.input, cmd = m.input.Update(msg)
	return m, cmd
}

func (m model) updateDelete(msg tea.Msg) (tea.Model, tea.Cmd) {
	if km, ok := msg.(tea.KeyMsg); ok {
		switch km.String() {
		case "y", "Y":
			if err := deleteWorktree(m.worktrees[m.cursor].Path); err != nil {
				m.err = err
				m.mode = listMode
				return m, nil
			}
			wts, err := listWorktrees()
			if err != nil {
				m.err = err
				m.mode = listMode
				return m, nil
			}
			m.worktrees = wts
			m.current = currentWorktreeIndex(wts)
			if m.cursor >= len(m.worktrees) {
				m.cursor = len(m.worktrees) - 1
			}
			m.mode = listMode
			m.err = nil
		default:
			m.mode = listMode
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.err != nil && len(m.worktrees) == 0 {
		return errStyle.Render(fmt.Sprintf("Error: %v", m.err)) + "\n"
	}

	var s strings.Builder
	s.WriteString(titleStyle.Render("Git Worktrees") + "\n\n")

	for i, wt := range m.worktrees {
		cursor := "  "
		if i == m.cursor {
			cursor = cursorStyle.Render("> ")
		}

		name := filepath.Base(wt.Path)
		branch := ""
		if wt.Branch != "" {
			branch = " " + branchStyle.Render("["+wt.Branch+"]")
		}
		if wt.IsBare {
			branch += dimStyle.Render(" (bare)")
		}

		var line string
		switch {
		case i == m.current:
			line = currentMark.Render(name) + branch + currentMark.Render(" ●")
		case i == m.cursor:
			line = activeStyle.Render(name) + branch
		default:
			line = dimStyle.Render(name) + branch
		}

		s.WriteString(cursor + line + "\n")
	}

	switch m.mode {
	case addMode:
		s.WriteString("\n" + promptStyle.Render("New branch: ") + m.input.View() + "\n")
		if m.err != nil {
			s.WriteString(errStyle.Render(fmt.Sprintf("  %v", m.err)) + "\n")
		}
	case deleteConfirmMode:
		wt := m.worktrees[m.cursor]
		s.WriteString("\n" + errStyle.Render(
			fmt.Sprintf("Delete worktree %s [%s]? (y/n)", filepath.Base(wt.Path), wt.Branch),
		) + "\n")
	default:
		if m.err != nil {
			s.WriteString("\n" + errStyle.Render(fmt.Sprintf("%v", m.err)) + "\n")
		}
	}

	s.WriteString("\n" + dimStyle.Render("j/k navigate • enter select • a add • d delete • esc/q quit") + "\n")
	return s.String()
}

func main() {
	p := tea.NewProgram(initialModel(), tea.WithAltScreen(), tea.WithOutput(os.Stderr))
	m, err := p.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if selected := m.(model).selected; selected != "" {
		fmt.Print(selected)
		os.Exit(0)
	}
	os.Exit(1)
}
