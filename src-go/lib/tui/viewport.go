package tui

import "github.com/charmbracelet/lipgloss"

const DefaultWidth = 80 // fallback width before first WindowSizeMsg

// AvailableRows returns the number of list rows that fit on screen after
// subtracting the rendered height of each section string.
// Returns 0 when windowHeight is 0 (no WindowSizeMsg yet).
func AvailableRows(windowHeight int, sections ...string) int {
	if windowHeight == 0 {
		return 0
	}
	overhead := 0
	for _, s := range sections {
		overhead += lipgloss.Height(s)
	}
	if rows := windowHeight - overhead; rows >= 1 {
		return rows
	}
	return 1
}

// Viewport tracks scrolling state for a list with more items than fit on screen.
type Viewport struct {
	Offset int
}

// MaxVisible returns the number of rows available for list items.
// availableRows is computed by the caller via lipgloss.Height on rendered sections.
// When availableRows is 0 (no WindowSizeMsg received yet), all items are considered visible.
func (v Viewport) MaxVisible(totalItems, availableRows int) int {
	if availableRows <= 0 {
		return totalItems
	}
	return max(1, availableRows)
}

// Clamp adjusts Offset so that cursor is always within the visible window.
func (v Viewport) Clamp(cursor, totalItems, availableRows int) Viewport {
	if availableRows <= 0 {
		return v
	}
	maxVis := v.MaxVisible(totalItems, availableRows)
	if cursor < v.Offset {
		v.Offset = cursor
	} else if cursor >= v.Offset+maxVis {
		v.Offset = cursor - maxVis + 1
	}
	if maxOffset := totalItems - maxVis; maxOffset >= 0 && v.Offset > maxOffset {
		v.Offset = maxOffset
	}
	if v.Offset < 0 {
		v.Offset = 0
	}
	return v
}
