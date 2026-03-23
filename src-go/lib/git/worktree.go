package git

import "strings"

// Worktree represents a single git worktree.
type Worktree struct {
	Path   string
	Branch string
	IsBare bool
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

// FindCurrentWorktree returns the index of the worktree whose path best matches cwd
// (longest prefix). Returns -1 if none match.
func FindCurrentWorktree(worktrees []Worktree, cwd string) int {
	best := -1
	bestLen := 0
	for i, wt := range worktrees {
		if strings.HasPrefix(cwd, wt.Path) && len(wt.Path) > bestLen {
			best = i
			bestLen = len(wt.Path)
		}
	}
	return best
}
