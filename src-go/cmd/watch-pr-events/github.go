package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
)

var (
	prNumber = regexp.MustCompile(`^[1-9][0-9]*$`)
	prPath   = regexp.MustCompile(`^/([A-Za-z0-9_.-]+)/([A-Za-z0-9_.-]+)/pull/([1-9][0-9]*)$`)
	headOID  = regexp.MustCompile(`^[0-9a-fA-F]{40}([0-9a-fA-F]{24})?$`)
)

type activity struct {
	ID    int64  `json:"id"`
	State string `json:"state"`
	User  struct {
		Login string `json:"login"`
	} `json:"user"`
	Kind                string `json:"kind,omitempty"`
	Title               string `json:"title,omitempty"`
	Body                string `json:"body"`
	HTMLURL             string `json:"html_url,omitempty"`
	CreatedAt           string `json:"created_at,omitempty"`
	UpdatedAt           string `json:"updated_at,omitempty"`
	SubmittedAt         string `json:"submitted_at,omitempty"`
	Path                string `json:"path,omitempty"`
	Line                int    `json:"line,omitempty"`
	OriginalLine        int    `json:"original_line,omitempty"`
	StartLine           int    `json:"start_line,omitempty"`
	Side                string `json:"side,omitempty"`
	DiffHunk            string `json:"diff_hunk,omitempty"`
	InReplyToID         int64  `json:"in_reply_to_id,omitempty"`
	PullRequestReviewID int64  `json:"pull_request_review_id,omitempty"`
}

type check struct {
	Bucket string `json:"bucket"`
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
