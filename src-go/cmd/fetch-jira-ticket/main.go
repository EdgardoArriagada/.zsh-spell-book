package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"example.com/workspace/lib/jira"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: fetch-jira-ticket <jira-url>")
		os.Exit(1)
	}

	issueKey := extractIssueKey(os.Args[1])

	email := os.Getenv("ZSB_JIRA_EMAIL")
	baseURL := os.Getenv("ZSB_JIRA_BASEURL")
	if email == "" || baseURL == "" {
		fmt.Fprintln(os.Stderr, "ZSB_JIRA_EMAIL and ZSB_JIRA_BASEURL must be set")
		os.Exit(1)
	}

	tokenBytes, err := exec.Command("pass", "jira").Output()
	if err != nil {
		fmt.Fprintln(os.Stderr, "pass jira:", err)
		os.Exit(1)
	}
	token := strings.TrimSpace(string(tokenBytes))

	summary, parentKey, err := fetchTicket(baseURL, issueKey, email, token)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if parentKey == "" {
		parentKey = "xxxxxxxxxxxx"
	}

	fmt.Printf("%s|%s|%s\n", parentKey, issueKey, summary)

	home := os.Getenv("HOME")
	lnum, err := jira.AppendTicketRow(home+"/temp/tickets", parentKey, issueKey, summary)
	if err != nil {
		fmt.Fprintln(os.Stderr, "write tickets:", err)
		os.Exit(1)
	}
	if lnum > 0 {
		fmt.Printf("ticket already on line %d\n", lnum)
	}
}

func extractIssueKey(rawURL string) string {
	parts := strings.Split(strings.TrimRight(rawURL, "/"), "/")
	return parts[len(parts)-1]
}

func fetchTicket(baseURL, issueKey, email, token string) (summary, parentKey string, err error) {
	req, err := http.NewRequest("GET", baseURL+"/rest/api/3/issue/"+issueKey+"?fields=summary,parent", nil)
	if err != nil {
		return "", "", err
	}
	req.SetBasicAuth(email, token)
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", fmt.Errorf("jira returned %d", resp.StatusCode)
	}

	var body struct {
		Fields struct {
			Summary string `json:"summary"`
			Parent  struct {
				Key string `json:"key"`
			} `json:"parent"`
		} `json:"fields"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return "", "", fmt.Errorf("json decode: %w", err)
	}

	return body.Fields.Summary, body.Fields.Parent.Key, nil
}
