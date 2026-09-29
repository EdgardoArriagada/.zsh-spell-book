package main

import (
	"bufio"
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
	"runtime"
	"strings"
	"syscall"
	"time"
)

const interval = 3 * time.Minute

const help = `Usage: watch-pr-events [-t|--tmux] [--codex-uuid UUID] [PR number|GitHub PR URL]

Watch a pull request for check results and new comments,
reviews, and merge readiness every 3 minutes.
With no argument, watch the PR for the current branch. Press Ctrl+C to stop.
Place watch-pr-events.conf beside main.go to ignore activity by GitHub username
(one username per line; blank lines and # comments are allowed).

Options:
  -t, --tmux         Notify the current tmux pane of PR activity
  --codex-uuid UUID  Send PR events to this Codex thread
  -h, --help         Show this help
`

var (
	prNumber = regexp.MustCompile(`^[1-9][0-9]*$`)
	prPath   = regexp.MustCompile(`^/([A-Za-z0-9_.-]+)/([A-Za-z0-9_.-]+)/pull/([1-9][0-9]*)$`)
	tmuxPane = regexp.MustCompile(`^%[0-9]+$`)
	codexID  = regexp.MustCompile(`^[0-9a-fA-F]{8}(-[0-9a-fA-F]{4}){3}-[0-9a-fA-F]{12}$`)
	headOID  = regexp.MustCompile(`^[0-9a-fA-F]{40}([0-9a-fA-F]{24})?$`)
)

type activity struct {
	ID    int64  `json:"id"`
	State string `json:"state"`
	User  struct {
		Login string `json:"login"`
	} `json:"user"`
	Kind string `json:"-"`
}

type check struct {
	Bucket string `json:"bucket"`
}

type notificationStrategy struct {
	notify  func(context.Context, string, bool) error
	failure string
	retry   bool
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
	uuid := ""
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if arg == "-t" || arg == "--tmux" {
			if tmux {
				return errors.New("usage: watch-pr-events [-t|--tmux] [--codex-uuid UUID] [PR number|GitHub PR URL]")
			}
			tmux = true
		} else if arg == "--codex-uuid" {
			if uuid != "" || i+1 >= len(args) || !codexID.MatchString(args[i+1]) {
				return errors.New("watch-pr-events: --codex-uuid requires a UUID")
			}
			i++
			uuid = args[i]
		} else {
			selection = append(selection, arg)
		}
	}
	if len(selection) > 1 || (len(selection) == 1 && !validPR(selection[0])) {
		return errors.New("usage: watch-pr-events [-t|--tmux] [--codex-uuid UUID] [PR number|GitHub PR URL]")
	}
	path, err := configPath()
	if err != nil {
		return err
	}
	ignored, err := loadIgnoredUsers(path)
	if err != nil {
		return err
	}
	strategies := []notificationStrategy{{
		notify: func(_ context.Context, event string, _ bool) error {
			_, err := fmt.Fprintln(out, event)
			return err
		},
		failure: "watch-pr-events: console output failed",
	}}
	if tmux {
		pane := os.Getenv("TMUX_PANE")
		if !tmuxPane.MatchString(pane) {
			return errors.New("watch-pr-events: tmux notifications require running inside a tmux pane")
		}
		notify, err := exec.LookPath("zsb_tmux_agent_notification")
		if err != nil {
			return errors.New("watch-pr-events: zsb_tmux_agent_notification not found")
		}
		strategies = append(strategies, notificationStrategy{
			notify:  func(ctx context.Context, _ string, _ bool) error { return tmuxNotification(ctx, notify, pane) },
			failure: "watch-pr-events: notification failed",
		})
	}
	if uuid != "" {
		codex, err := exec.LookPath("codex")
		if err != nil {
			return errors.New("watch-pr-events: codex not found")
		}
		codex, err = filepath.Abs(codex)
		if err != nil {
			return errors.New("watch-pr-events: cannot resolve codex path")
		}
		strategies = append(strategies, notificationStrategy{
			notify: func(ctx context.Context, event string, _ bool) error {
				return queueCodexEvent(ctx, codex, uuid, event)
			},
			failure: "watch-pr-events: codex notification failed; will retry",
			retry:   true,
		})
	}
	gh, err := exec.LookPath("gh")
	if err != nil {
		return errors.New("watch-pr-events: gh not found")
	}
	gh, err = filepath.Abs(gh)
	if err != nil {
		return errors.New("watch-pr-events: cannot resolve gh path")
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
	mergeReadyHead := ""
	checkHead, checkStatus := "", ""
	notifiedCheckHead, notifiedCheckStatus := "", ""
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		head, result, err := fetchCheckStatus(ctx, gh, prURL)
		if ctx.Err() != nil {
			return nil
		}
		if err != nil {
			checkHead, checkStatus = "", ""
			fmt.Fprintln(os.Stderr, err)
		} else {
			checkHead, checkStatus = head, result
			if result == "pending" {
				notifiedCheckHead, notifiedCheckStatus = "", ""
			} else if result != "" && (head != notifiedCheckHead || result != notifiedCheckStatus) {
				event := fmt.Sprintf("PR checks %s: %s (head %s)", result, prURL, head)
				if !notifyStrategies(ctx, strategies, event, false) {
					notifiedCheckHead, notifiedCheckStatus = head, result
				}
			}
		}
		items, err := fetchActivity(ctx, gh, endpoint)
		if ctx.Err() != nil {
			return nil
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else if !initialized {
			newActivity(seen, items, ignored)
			initialized = true
			fmt.Fprintf(out, "Watching %s for PR checks and activity (every 3 minutes).\n", prURL)
		} else {
			fresh := newActivity(seen, items, ignored)
			var events []string
			for _, item := range fresh {
				actor := item.User.Login
				if actor == "" {
					actor = "someone"
				}
				event := fmt.Sprintf("[%s] %s by %q", time.Now().Format("15:04:05"), label(item), actor)
				events = append(events, event)
			}
			if len(events) > 0 && notifyStrategies(ctx, strategies, strings.Join(events, "\n")+"\nPR: "+prURL, false) {
				for _, item := range fresh {
					delete(seen, fmt.Sprintf("%s:%d:%s", item.Kind, item.ID, item.State))
				}
			}
		}
		ready, head, err := fetchMergeReady(ctx, gh, prURL)
		if ctx.Err() != nil {
			return nil
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
		} else if !ready || head != checkHead || checkStatus != "passed" {
			mergeReadyHead = ""
		} else if head != mergeReadyHead {
			event := fmt.Sprintf("PR ready to merge: %s (head %s)", prURL, head)
			if !notifyStrategies(ctx, strategies, event, true) {
				mergeReadyHead = head
			}
		}
		select {
		case <-ctx.Done():
			return nil
		case <-ticker.C:
		}
	}
}

func notifyStrategies(ctx context.Context, strategies []notificationStrategy, event string, mergeReady bool) (retry bool) {
	for _, strategy := range strategies {
		if err := strategy.notify(ctx, event, mergeReady); err != nil {
			fmt.Fprintln(os.Stderr, strategy.failure)
			retry = retry || strategy.retry
		}
	}
	return retry
}

func queueCodexEvent(ctx context.Context, codex, uuid, event string) error {
	return exec.CommandContext(ctx, codex, "queue", "--thread", uuid, "--message", event).Run()
}

func configPath() (string, error) {
	_, source, _, ok := runtime.Caller(0)
	if ok && filepath.IsAbs(source) {
		return filepath.Join(filepath.Dir(source), "watch-pr-events.conf"), nil
	}
	// Release builds use -trimpath; their binary lives in src-go/bin.
	executable, err := os.Executable()
	if err != nil {
		return "", errors.New("watch-pr-events: cannot locate config")
	}
	return filepath.Join(filepath.Dir(executable), "..", "cmd", "watch-pr-events", "watch-pr-events.conf"), nil
}

func loadIgnoredUsers(path string) (map[string]bool, error) {
	ignored := make(map[string]bool)
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return ignored, nil
	}
	if err != nil {
		return nil, fmt.Errorf("watch-pr-events: cannot open config: %w", err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		user := strings.TrimSpace(scanner.Text())
		if user != "" && !strings.HasPrefix(user, "#") {
			ignored[strings.ToLower(user)] = true
		}
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("watch-pr-events: cannot read config: %w", err)
	}
	return ignored, nil
}

func fetchHead(ctx context.Context, gh, prURL string) (string, error) {
	output, err := exec.CommandContext(ctx, gh, "pr", "view", prURL, "--json", "headRefOid").Output()
	if err != nil {
		return "", errors.New("watch-pr-events: could not fetch PR head")
	}
	var pr struct {
		HeadRefOID string `json:"headRefOid"`
	}
	if json.Unmarshal(output, &pr) != nil || !headOID.MatchString(pr.HeadRefOID) {
		return "", errors.New("watch-pr-events: invalid PR head response")
	}
	return pr.HeadRefOID, nil
}

func fetchCheckStatus(ctx context.Context, gh, prURL string) (string, string, error) {
	before, err := fetchHead(ctx, gh, prURL)
	if err != nil {
		return "", "", err
	}
	output, err := exec.CommandContext(ctx, gh, "pr", "checks", "--json", "bucket", prURL).Output()
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || (exitErr.ExitCode() != 1 && exitErr.ExitCode() != 8) {
			return "", "", errors.New("watch-pr-events: could not fetch PR checks")
		}
	}
	var checks []check
	if json.Unmarshal(output, &checks) != nil {
		return "", "", errors.New("watch-pr-events: invalid PR checks response")
	}
	result := checksStatus(checks)
	if result == "pending" {
		return before, result, nil
	}
	after, err := fetchHead(ctx, gh, prURL)
	if err != nil {
		return "", "", err
	}
	if before != after {
		return "", "", nil
	}
	return before, result, nil
}

func checksStatus(checks []check) string {
	if len(checks) == 0 {
		return "pending"
	}
	failed := false
	for _, check := range checks {
		switch check.Bucket {
		case "pass", "skipping":
		case "fail", "cancel":
			failed = true
		default:
			return "pending"
		}
	}
	if failed {
		return "failed"
	}
	return "passed"
}

func fetchMergeReady(ctx context.Context, gh, prURL string) (bool, string, error) {
	output, err := exec.CommandContext(ctx, gh, "pr", "view", prURL, "--json", "state,mergeStateStatus,headRefOid").Output()
	if err != nil {
		return false, "", errors.New("watch-pr-events: could not fetch merge status")
	}
	var pr struct {
		State            string `json:"state"`
		MergeStateStatus string `json:"mergeStateStatus"`
		HeadRefOID       string `json:"headRefOid"`
	}
	if json.Unmarshal(output, &pr) != nil || !headOID.MatchString(pr.HeadRefOID) {
		return false, "", errors.New("watch-pr-events: invalid merge status response")
	}
	return pr.State == "OPEN" && pr.MergeStateStatus == "CLEAN", pr.HeadRefOID, nil
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
		return "", "", errors.New("watch-pr-events: unsupported PR URL")
	}
	parts := prPath.FindStringSubmatch(u.Path)
	if parts == nil {
		return "", "", errors.New("watch-pr-events: unsupported PR URL")
	}
	return u.String(), fmt.Sprintf("repos/%s/%s/pulls/%s", parts[1], parts[2], parts[3]), nil
}

func resolvePR(ctx context.Context, gh string, selection []string) (string, string, error) {
	args := []string{"pr", "view", "--json", "url"}
	args = append(args, selection...)
	output, err := exec.CommandContext(ctx, gh, args...).Output()
	if err != nil {
		return "", "", errors.New("watch-pr-events: cannot find PR (check gh authentication and selection)")
	}
	var pr struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal(output, &pr); err != nil {
		return "", "", errors.New("watch-pr-events: invalid gh pr response")
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
			return nil, fmt.Errorf("watch-pr-events: could not fetch %s activity", source.kind)
		}
		var pages [][]activity
		if err := json.Unmarshal(output, &pages); err != nil {
			return nil, fmt.Errorf("watch-pr-events: invalid %s response", source.kind)
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

func newActivity(seen map[string]bool, items []activity, ignored map[string]bool) []activity {
	var fresh []activity
	for _, item := range items {
		if item.ID == 0 || ignored[strings.ToLower(item.User.Login)] || (item.Kind == "review" && item.State == "PENDING") {
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
