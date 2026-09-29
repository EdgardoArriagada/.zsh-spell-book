package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
	"unicode"

	gitlib "example.com/workspace/lib/git"
)

// Worktree is an alias for the shared type in lib/git.
type Worktree = gitlib.Worktree

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
	return gitlib.ParseWorktreeList(string(out)), nil
}

func createWorktree(mainPath, branch string) error {
	if out, err := exec.Command("git", "-C", mainPath, "fetch", "--all", "--prune").CombinedOutput(); err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}

	baseDir := WorktreeBaseDir(mainPath)
	wtPath := filepath.Join(baseDir, branch)
	var args []string
	if gitlib.BranchExists(branch) || gitlib.RemoteBranchExists(branch) {
		args = []string{"worktree", "add", wtPath, branch}
	} else {
		base, err := baseBranch(branch)
		if err != nil {
			return err
		}
		args = []string{"worktree", "add", "--no-track", wtPath, "-b", branch, base}
	}
	out, err := exec.Command("git", append([]string{"-C", mainPath}, args...)...).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func baseBranch(branch string) (string, error) {
	names := []string{"develop", "master", "main"}
	if strings.HasPrefix(branch, "hotfix/") || strings.HasPrefix(branch, "fix/") {
		names = names[1:]
	}
	for _, name := range names {
		if ref, ok := gitlib.RemoteBranchRef(name); ok {
			return ref, nil
		}
		if gitlib.BranchExists(name) {
			return name, nil
		}
	}
	return "", fmt.Errorf("base branch not found")
}

func deleteWorktree(path string) error {
	out, err := exec.Command("git", "worktree", "remove", path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func deleteWorktreeForce(path string) error {
	out, err := exec.Command("git", "worktree", "remove", "--force", path).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func deleteWorktreeBranch(branch string) error {
	out, err := exec.Command("git", "branch", "-D", branch).CombinedOutput()
	if err != nil {
		return fmt.Errorf("%s", strings.TrimSpace(string(out)))
	}
	return nil
}

func isWorktreeDirtyError(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "contains modified or untracked files")
}

func currentWorktreeIndex(worktrees []Worktree) int {
	cwd, err := os.Getwd()
	if err != nil {
		return -1
	}
	return gitlib.FindCurrentWorktree(worktrees, cwd)
}

func worktreePRStatus(ctx context.Context, wt Worktree) (string, error) {
	if wt.IsBare || wt.Branch == "" || wt.Branch == "(detached)" {
		return "no PR (bare or detached worktree)", nil
	}
	if err := gitlib.ValidateBranchName(wt.Branch); err != nil {
		return "", errors.New("invalid worktree branch")
	}
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	var prs []struct {
		Number            int
		State             string
		IsDraft           bool
		ReviewDecision    string
		URL               string
		StatusCheckRollup []struct {
			Status, Conclusion, State string
		}
	}
	// Prefer an open PR; otherwise show the latest closed or merged PR.
	for _, state := range []string{"open", "all"} {
		cmd := exec.CommandContext(ctx, "gh", "pr", "list", "--head="+wt.Branch, "--state="+state, "--limit=1", "--json=number,state,isDraft,reviewDecision,statusCheckRollup,url")
		cmd.Dir = wt.Path
		out, err := cmd.Output()
		if err != nil {
			if ctx.Err() != nil {
				return "", errors.New("PR lookup cancelled or timed out")
			}
			if errors.Is(err, exec.ErrNotFound) {
				return "", errors.New("gh not found; install GitHub CLI")
			}
			return "", errors.New("PR lookup failed; check gh authentication and repository access")
		}
		if json.Unmarshal(out, &prs) != nil {
			return "", errors.New("invalid GitHub PR response")
		}
		if len(prs) > 0 {
			break
		}
	}
	if len(prs) == 0 {
		return "no PR", nil
	}
	pr := prs[0]
	state := strings.ToLower(pr.State)
	if pr.IsDraft && pr.State == "OPEN" {
		state = "draft"
	}
	review := strings.ToLower(strings.ReplaceAll(pr.ReviewDecision, "_", " "))
	if review == "" {
		review = "none"
	}
	passed, pending, failed := 0, 0, 0
	for _, check := range pr.StatusCheckRollup {
		result := check.State
		if result == "" {
			if check.Status != "COMPLETED" {
				pending++
				continue
			}
			result = check.Conclusion
		}
		switch result {
		case "SUCCESS", "NEUTRAL", "SKIPPED":
			passed++
		case "FAILURE", "ERROR", "CANCELLED", "TIMED_OUT", "ACTION_REQUIRED", "STARTUP_FAILURE", "STALE":
			failed++
		default:
			pending++
		}
	}
	checks := "none"
	if len(pr.StatusCheckRollup) > 0 {
		checks = fmt.Sprintf("%d passed, %d pending, %d failed", passed, pending, failed)
	}
	// Strip terminal controls from remote text before rendering it.
	clean := func(s string) string {
		return strings.Map(func(r rune) rune {
			if unicode.IsControl(r) {
				return -1
			}
			return r
		}, s)
	}
	return fmt.Sprintf("#%d %s\nchecks: %s\nreview: %s\n%s", pr.Number, clean(state), checks, clean(review), clean(pr.URL)), nil
}
