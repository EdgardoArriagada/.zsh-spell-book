package tui

const ViewportOverhead = 6  // title(2) + status/confirm(2) + footer(2)
const DefaultWidth      = 80 // fallback width before first WindowSizeMsg

// Viewport tracks scrolling state for a list with more items than fit on screen.
type Viewport struct {
	Height int
	Offset int
}

// MaxVisible returns the number of items that fit in the terminal.
// When Height is 0 (no WindowSizeMsg received yet), all items are considered visible.
func (v Viewport) MaxVisible(totalItems int) int {
	if v.Height == 0 {
		return totalItems
	}
	n := v.Height - ViewportOverhead
	if n < 1 {
		return 1
	}
	return n
}

// Clamp adjusts Offset so that cursor is always within the visible window.
func (v Viewport) Clamp(cursor, totalItems int) Viewport {
	if v.Height == 0 {
		return v
	}
	maxVis := v.MaxVisible(totalItems)
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
