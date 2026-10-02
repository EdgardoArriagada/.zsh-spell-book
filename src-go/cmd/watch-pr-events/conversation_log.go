package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

type conversationEntry struct {
	ObservedAt string   `json:"observed_at"`
	Change     string   `json:"change"`
	Activity   activity `json:"activity"`
}

func logStrategy(gh, endpoint, prURL string) (notificationStrategy, *os.File, error) {
	failure := errors.New("watch-pr-events: cannot open conversation log")
	_, endpoint, err := parsePRURL(prURL)
	if err != nil {
		return notificationStrategy{}, nil, err
	}
	parts := strings.Split(endpoint, "/")
	if parts[1] == "." || parts[1] == ".." || parts[2] == "." || parts[2] == ".." {
		return notificationStrategy{}, nil, failure
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return notificationStrategy{}, nil, failure
	}
	base := filepath.Join(home, "temp", "watch-pr-events")
	if err := os.MkdirAll(base, 0700); err != nil {
		return notificationStrategy{}, nil, failure
	}
	root, err := os.OpenRoot(base)
	if err != nil {
		return notificationStrategy{}, nil, failure
	}
	defer root.Close()
	dir := filepath.Join(parts[1], parts[2])
	if err := root.MkdirAll(dir, 0700); err != nil {
		return notificationStrategy{}, nil, failure
	}
	file, err := root.OpenFile(filepath.Join(dir, parts[4]+".log"), os.O_CREATE|os.O_RDWR|os.O_APPEND, 0600)
	if err != nil {
		return notificationStrategy{}, nil, failure
	}
	if err := file.Chmod(0600); err != nil {
		file.Close()
		return notificationStrategy{}, nil, failure
	}
	// ponytail: one watcher per PR; add file locking if concurrent writers are needed.
	latest := make(map[string]activity)
	key := func(item activity) string { return fmt.Sprintf("%s:%d", item.Kind, item.ID) }
	decoder := json.NewDecoder(file)
	for {
		var entry conversationEntry
		if err := decoder.Decode(&entry); err != nil {
			if err == io.EOF {
				break
			}
			file.Close()
			return notificationStrategy{}, nil, errors.New("watch-pr-events: invalid conversation log; preserve or repair it before retrying")
		}
		latest[key(entry.Activity)] = entry.Activity
	}
	return notificationStrategy{
		failure: "watch-pr-events: conversation logging failed; will retry",
		conversation: func(ctx context.Context, items []activity) error {
			output, err := exec.CommandContext(ctx, gh, "api", endpoint).Output()
			if err != nil {
				return errors.New("could not fetch PR description")
			}
			var description activity
			if err := json.Unmarshal(output, &description); err != nil || description.ID == 0 {
				return errors.New("invalid PR description")
			}
			description.Kind = "description"
			conversation := append([]activity{description}, items...)
			sort.SliceStable(conversation, func(i, j int) bool {
				timestamp := func(item activity) string {
					if item.CreatedAt != "" {
						return item.CreatedAt
					}
					return item.SubmittedAt
				}
				return timestamp(conversation[i]) < timestamp(conversation[j])
			})
			for _, item := range conversation {
				if item.ID == 0 || (item.Kind == "review" && item.State == "PENDING") {
					continue
				}
				previous, exists := latest[key(item)]
				if item.Kind == "description" {
					// PR updated_at also changes when comments arrive.
					previous.UpdatedAt = item.UpdatedAt
				}
				if exists && previous == item {
					continue
				}
				change := "new"
				if exists {
					change = "edited"
				}
				entry := conversationEntry{time.Now().UTC().Format(time.RFC3339Nano), change, item}
				if err := json.NewEncoder(file).Encode(entry); err != nil {
					return err
				}
				if err := file.Sync(); err != nil {
					return err
				}
				latest[key(item)] = item
			}
			return nil
		},
	}, file, nil
}
