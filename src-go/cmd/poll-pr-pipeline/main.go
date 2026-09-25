package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"regexp"
	"syscall"
	"time"
)

const interval = 30 * time.Second

const help = `Usage: poll-pr-pipeline [-t|--tmux] [PR number|GitHub PR URL]

Watch pull request checks until they complete. With no argument, watch the PR
for the current branch. Press Ctrl+C to stop.

Options:
  -t, --tmux  Notify the current tmux pane while checks run and when they finish
  -h, --help  Show this help
`

var (
	prNumber = regexp.MustCompile(`^[1-9][0-9]*$`)
	prURL    = regexp.MustCompile(`^/[^/]+/[^/]+/pull/[1-9][0-9]*$`)
)

type check struct {
	Bucket string `json:"bucket"`
}

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string) (result error) {
	if len(args) == 1 && (args[0] == "-h" || args[0] == "--help") {
		_, err := fmt.Print(help)
		return err
	}
	var selection []string
	tmux := false
	for _, arg := range args {
		if arg == "-t" || arg == "--tmux" {
			if tmux {
				return errors.New("usage: poll-pr-pipeline [-t|--tmux] [PR number|GitHub PR URL]")
			}
			tmux = true
		} else {
			selection = append(selection, arg)
		}
	}
	if len(selection) > 1 || (len(selection) == 1 && !validPR(selection[0])) {
		return errors.New("usage: poll-pr-pipeline [-t|--tmux] [PR number|GitHub PR URL]")
	}

	var notify, pane string
	if tmux {
		pane = os.Getenv("TMUX_PANE")
		if !regexp.MustCompile(`^%[0-9]+$`).MatchString(pane) {
			return errors.New("poll-pr-pipeline: run -t inside a tmux pane")
		}
	}
	gh, err := commandPath("gh")
	if err != nil {
		return err
	}
	if tmux {
		notify, err = commandPath("zsb_tmux_agent_notification")
		if err != nil {
			return err
		}
	}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	if tmux {
		if err := notification(notify, "--working", pane); err != nil {
			return err
		}
		defer func() {
			if err := notification(notify, "--force-finished", pane); err != nil && result == nil {
				result = err
			}
		}()
	}

	for {
		checks, err := fetchChecks(ctx, gh, selection)
		if err != nil {
			if ctx.Err() != nil {
				return errors.New("poll-pr-pipeline: interrupted")
			}
			return errors.New("poll-pr-pipeline: gh pr checks failed")
		}
		if done, failed := status(checks); done {
			if failed {
				return errors.New("poll-pr-pipeline: checks failed, cancelled, or timed out")
			}
			return nil
		}
		fmt.Fprintln(os.Stderr, "pipeline pending; checking again in 30s")
		select {
		case <-ctx.Done():
			return errors.New("poll-pr-pipeline: interrupted")
		case <-time.After(interval):
		}
	}
}

func commandPath(name string) (string, error) {
	path, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("poll-pr-pipeline: %s not found", name)
	}
	path, err = filepath.Abs(path)
	if err != nil {
		return "", errors.New("poll-pr-pipeline: cannot resolve command path")
	}
	return path, nil
}

func notification(command, state, pane string) error {
	if err := exec.Command(command, state, "_", pane).Run(); err != nil {
		return errors.New("poll-pr-pipeline: notification failed")
	}
	return nil
}

func fetchChecks(ctx context.Context, gh string, selection []string) ([]check, error) {
	args := []string{"pr", "checks", "--json", "bucket"}
	args = append(args, selection...)
	output, err := exec.CommandContext(ctx, gh, args...).Output()
	if err != nil {
		var exitErr *exec.ExitError
		if !errors.As(err, &exitErr) || exitErr.ExitCode() != 8 {
			return nil, err
		}
	}
	var checks []check
	if err := json.Unmarshal(output, &checks); err != nil {
		return nil, err
	}
	return checks, nil
}

func validPR(value string) bool {
	if prNumber.MatchString(value) {
		return true
	}
	parsed, err := url.Parse(value)
	return err == nil && parsed.Scheme == "https" && parsed.Host == "github.com" && prURL.MatchString(parsed.Path) && parsed.RawQuery == "" && parsed.Fragment == ""
}

func status(checks []check) (done, failed bool) {
	for _, check := range checks {
		switch check.Bucket {
		case "pass", "skipping":
		case "fail", "cancel":
			failed = true
		default:
			return false, false
		}
	}
	return true, failed
}
