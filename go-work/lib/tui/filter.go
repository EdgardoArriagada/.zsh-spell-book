package tui

import "github.com/sahilm/fuzzy"

// FuzzyFilter returns indices from targets that match term, sorted by score.
// If term is empty, returns all indices in order.
func FuzzyFilter(term string, targets []string) []int {
	if term == "" {
		idx := make([]int, len(targets))
		for i := range idx {
			idx[i] = i
		}
		return idx
	}
	matches := fuzzy.Find(term, targets)
	idx := make([]int, len(matches))
	for i, m := range matches {
		idx[i] = m.Index
	}
	return idx
}
