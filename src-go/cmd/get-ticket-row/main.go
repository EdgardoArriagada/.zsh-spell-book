package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: get-ticket-row <jira-url>")
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

	row := fmt.Sprintf("%s|%s|%s", parentKey, issueKey, summary)
	fmt.Println(row)

	home := os.Getenv("HOME")
	ticketsPath := home + "/temp/tickets"

	if lnum := findTicketLine(ticketsPath, issueKey); lnum > 0 {
		fmt.Printf("ticket already on line %d\n", lnum)
		return
	}

	f, err := os.OpenFile(ticketsPath, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, "open tickets:", err)
		os.Exit(1)
	}
	defer f.Close()
	fmt.Fprintln(f, row)
}

func findTicketLine(path, issueKey string) int {
	f, err := os.Open(path)
	if err != nil {
		return 0
	}
	defer f.Close()
	needle := "|" + issueKey + "|"
	sc := bufio.NewScanner(f)
	for lnum := 1; sc.Scan(); lnum++ {
		if strings.Contains(sc.Text(), needle) {
			return lnum
		}
	}
	return 0
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
