package main

import (
	"bufio"
	"context"
	"errors"
	"flag"
	"fmt"
	"io"
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

const help = `Usage: watch-pr-events [-t|--tmux] [-l|--logs] [--codex-thread UUID] [-p|--pull-request PR]
       watch-pr-events -h|--help

Watch a pull request for check results and new comments,
reviews, and merge readiness immediately and every 3 minutes.
Existing comments and reviews form the notification baseline; logs include initial history.
Merge readiness requires an open PR, CLEAN merge status, and passed checks on the same head.
With no --pull-request flag, watch the PR for the current branch in the current repository.
PR numbers select from the current repository. URLs must have the form
https://github.com/OWNER/REPO/pull/NUMBER, without a query, fragment, or trailing slash.
Requires authenticated gh. Events always print to the console. Press Ctrl+C to stop.
Place optional watch-pr-events.conf beside main.go to ignore comments and reviews
by GitHub username (case-insensitive; one per line; blank lines and # comment lines
are allowed). Check results and merge readiness are not filtered by username.

Options:
  -p, --pull-request PR  Select a PR number or GitHub PR URL
  -l, --logs           Append the PR description and full conversation from all users
                      to ~/temp/watch-pr-events/OWNER/REPO/NUMBER.log (JSON Lines)
                      Includes existing history, new messages, and observed edits;
                      preserves old versions and resumes without duplicating history
  -t, --tmux           Also notify the current tmux pane of all PR events
                      Requires a tmux pane and zsb_tmux_agent_notification on PATH
  --codex-thread UUID  Also queue PR events for this Codex thread; retry failed queues
                      Requires codex on PATH with the queue command
  -h, --help           Show this help
`

var (
	tmuxPane = regexp.MustCompile(`^%[0-9]+$`)
	codexID  = regexp.MustCompile(`^[0-9a-fA-F]{8}(-[0-9a-fA-F]{4}){3}-[0-9a-fA-F]{12}$`)
)

type notificationStrategy struct {
	notify       func(context.Context, string, bool) error
	conversation func(context.Context, []activity) error
	failure      string
	retry        bool
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
	var tmux, logs bool
	var uuid, pullRequest string
	flags := flag.NewFlagSet("watch-pr-events", flag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.BoolVar(&tmux, "t", false, "Also notify the current tmux pane")
	flags.BoolVar(&tmux, "tmux", false, "Also notify the current tmux pane")
	flags.BoolVar(&logs, "l", false, "Append the PR conversation to a log")
	flags.BoolVar(&logs, "logs", false, "Append the PR conversation to a log")
	flags.StringVar(&uuid, "codex-thread", "", "Queue PR events for a Codex thread UUID")
	flags.StringVar(&pullRequest, "p", "", "PR number or GitHub PR URL")
	flags.StringVar(&pullRequest, "pull-request", "", "PR number or GitHub PR URL")
	if err := flags.Parse(args); err != nil {
		if errors.Is(err, flag.ErrHelp) {
			_, err := fmt.Fprint(out, help)
			return err
		}
		return err
	}
	var selection []string
	threadSet := false
	flags.Visit(func(f *flag.Flag) {
		if f.Name == "p" || f.Name == "pull-request" {
			selection = []string{pullRequest}
		}
		if f.Name == "codex-thread" {
			threadSet = true
		}
	})
	if flags.NArg() != 0 {
		return errors.New("watch-pr-events: positional arguments are not supported; use -p or --pull-request")
	}
	if len(selection) != 0 && !validPR(pullRequest) {
		return errors.New("watch-pr-events: --pull-request requires a PR number or GitHub PR URL")
	}
	if threadSet && !codexID.MatchString(uuid) {
		return errors.New("watch-pr-events: --codex-thread requires a UUID")
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
	if logs {
		strategy, file, err := logStrategy(gh, endpoint, prURL)
		if err != nil {
			return err
		}
		defer file.Close()
		strategies = append(strategies, strategy)
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
		if err == nil {
			for _, strategy := range strategies {
				if strategy.conversation != nil && strategy.conversation(ctx, items) != nil {
					fmt.Fprintln(os.Stderr, strategy.failure)
				}
			}
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
		if strategy.notify == nil {
			continue
		}
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

func tmuxNotification(ctx context.Context, command, pane string) error {
	return exec.CommandContext(ctx, command, "--force-finished", "_", pane).Run()
}
