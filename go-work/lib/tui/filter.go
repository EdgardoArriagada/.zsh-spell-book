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

// ApplyFilter returns items filtered and ranked by FuzzyFilter using key as the string projection.
// If term is empty, returns all items in order.
func ApplyFilter[T any](items []T, term string, key func(T) string) []T {
	keys := make([]string, len(items))
	for i, item := range items {
		keys[i] = key(item)
	}
	indices := FuzzyFilter(term, keys)
	result := make([]T, len(indices))
	for i, idx := range indices {
		result[i] = items[idx]
	}
	return result
}
