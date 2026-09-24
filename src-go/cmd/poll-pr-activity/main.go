package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"strings"
	"syscall"
	"time"
)

const interval = time.Minute

const help = `Usage: poll-pr-activity [-t|--tmux] [PR number|GitHub PR URL]

Watch a pull request for new comments and reviews every minute.
With no argument, watch the PR for the current branch. Press Ctrl+C to stop.

Options:
  -t, --tmux  Notify the current tmux pane when new activity appears
  -h, --help  Show this help
`

var (
	prNumber = regexp.MustCompile(`^[1-9][0-9]*$`)
	prPath   = regexp.MustCompile(`^/([A-Za-z0-9_.-]+)/([A-Za-z0-9_.-]+)/pull/([1-9][0-9]*)$`)
	tmuxPane = regexp.MustCompile(`^%[0-9]+$`)
)

type activity struct {
	ID    int64  `json:"id"`
	State string `json:"state"`
	User  struct {
		Login string `json:"login"`
	} `json:"user"`
	Kind string `json:"-"`
}

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if err := run(ctx, os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(ctx context.Context, args []string, out io.Writer) error {
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		_, err := fmt.Fprint(out, help)
		return err
	}
	var selection []string
	tmux := false
	for _, arg := range args {
		if arg == "-t" || arg == "--tmux" {
			if tmux {
				return errors.New("usage: poll-pr-activity [-t|--tmux] [PR number|GitHub PR URL]")
			}
			tmux = true
		} else {
			selection = append(selection, arg)
		}
	}
	if len(selection) > 1 || (len(selection) == 1 && !validPR(selection[0])) {
		return errors.New("usage: poll-pr-activity [-t|--tmux] [PR number|GitHub PR URL]")
	}
	var notify, pane string
	if tmux {
		pane = os.Getenv("TMUX_PANE")
		if !tmuxPane.MatchString(pane) {
			return errors.New("poll-pr-activity: run -t inside a tmux pane")
		}
		var err error
		notify, err = exec.LookPath("zsb_tmux_agent_notification")
		if err != nil {
			return errors.New("poll-pr-activity: zsb_tmux_agent_notification not found")
		}
	}
	gh, err := exec.LookPath("gh")
	if err != nil {
		return errors.New("poll-pr-activity: gh not found")
	}
	gh, err = filepath.Abs(gh)
	if err != nil {
		return errors.New("poll-pr-activity: cannot resolve gh path")
	}
	prURL, endpoint, err := resolvePR(ctx, gh, selection)
	if err != nil {
		if ctx.Err() != nil {
			return nil
		}
		return err
	}

	seen := make(map[string]bool)
	initialized := false
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		items, err := fetchActivity(ctx, gh, endpoint)
		if ctx.Err() != nil {
			return nil
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else if !initialized {
			newActivity(seen, items)
			initialized = true
			fmt.Fprintf(out, "Watching %s for new PR activity (every minute).\n", prURL)
		} else {
			fresh := newActivity(seen, items)
			for _, item := range fresh {
				actor := item.User.Login
				if actor == "" {
					actor = "someone"
				}
				fmt.Fprintf(out, "[%s] %s by %q: %s\n", time.Now().Format("15:04:05"), label(item), actor, prURL)
			}
			if tmux && len(fresh) > 0 {
				if err := tmuxNotification(ctx, notify, pane); err != nil {
					fmt.Fprintln(os.Stderr, "poll-pr-activity: notification failed")
				}
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func tmuxNotification(ctx context.Context, command, pane string) error {
	return exec.CommandContext(ctx, command, "--force-finished", "_", pane).Run()
}

func validPR(value string) bool {
	if prNumber.MatchString(value) {
		return true
	}
	_, _, err := parsePRURL(value)
	return err == nil
}

func parsePRURL(value string) (string, string, error) {
	u, err := url.Parse(value)
	if err != nil || u.Scheme != "https" || u.Host != "github.com" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", "", errors.New("poll-pr-activity: unsupported PR URL")
	}
	parts := prPath.FindStringSubmatch(u.Path)
	if parts == nil {
		return "", "", errors.New("poll-pr-activity: unsupported PR URL")
	}
	return u.String(), fmt.Sprintf("repos/%s/%s/pulls/%s", parts[1], parts[2], parts[3]), nil
}

func resolvePR(ctx context.Context, gh string, selection []string) (string, string, error) {
	args := []string{"pr", "view", "--json", "url"}
	args = append(args, selection...)
	output, err := exec.CommandContext(ctx, gh, args...).Output()
	if err != nil {
		return "", "", errors.New("poll-pr-activity: cannot find PR (check gh authentication and selection)")
	}
	var pr struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(output, &pr); err != nil {
		return "", "", errors.New("poll-pr-activity: invalid gh pr response")
	}
	return parsePRURL(pr.URL)
}

func fetchActivity(ctx context.Context, gh, endpoint string) ([]activity, error) {
	// ponytail: scans full history; add incremental cursors if large PRs hit rate limits.
	i := strings.LastIndex(endpoint, "/pulls/")
	sources := []struct {
		kind, path string
	}{
		{"comment", endpoint[:i] + "/issues" + endpoint[i+len("/pulls"):] + "/comments"},
		{"review", endpoint + "/reviews"},
		{"review comment", endpoint + "/comments"},
	}
	var items []activity
	for _, source := range sources {
		output, err := exec.CommandContext(ctx, gh, "api", source.path+"?per_page=100", "--paginate", "--slurp").Output()
		if err != nil {
			return nil, fmt.Errorf("poll-pr-activity: could not fetch %s activity", source.kind)
		}
		var pages [][]activity
		if err := json.Unmarshal(output, &pages); err != nil {
			return nil, fmt.Errorf("poll-pr-activity: invalid %s response", source.kind)
		}
		for _, page := range pages {
			for _, item := range page {
				item.Kind = source.kind
				items = append(items, item)
			}
		}
	}
	return items, nil
}

func newActivity(seen map[string]bool, items []activity) []activity {
	var fresh []activity
	for _, item := range items {
		if item.ID == 0 || (item.Kind == "review" && item.State == "PENDING") {
			continue
		}
		key := fmt.Sprintf("%s:%d:%s", item.Kind, item.ID, item.State)
		if !seen[key] {
			seen[key] = true
			fresh = append(fresh, item)
		}
	}
	return fresh
}

func label(item activity) string {
	switch item.Kind {
	case "comment":
		return "New PR comment"
	case "review comment":
		return "New review comment"
	case "review":
		switch item.State {
		case "APPROVED":
			return "PR approved"
		case "CHANGES_REQUESTED":
			return "Changes requested"
		case "DISMISSED":
			return "Review dismissed"
		}
	}
	return "New review"
}
