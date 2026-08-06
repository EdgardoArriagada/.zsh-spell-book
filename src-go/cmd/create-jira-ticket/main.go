package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"

	"example.com/workspace/lib/jira"
)

const (
	allowedJiraBaseURL = "https://mercadolibre.atlassian.net"
	usage              = `Usage: create-jira-ticket "<title>" ["<description>"]`
	controlPointField  = "customfield_25390"
)

var (
	jiraKeyRe       = regexp.MustCompile(`^[A-Z][A-Z0-9]+-\d+$`)
	controlPointRe  = regexp.MustCompile(`\[[^\[\]]+\]`)
	requiredEnvVars = []string{
		"ZSB_JIRA_BASEURL",
		"ZSB_JIRA_EMAIL",
		"ZSB_JIRA_PROJECT_KEY",
		"ZSB_JIRA_ISSUE_TYPE_IDS",
		"ZSB_JIRA_REPORTER_ACCOUNT_ID",
		"ZSB_JIRA_ASSIGNEE_ACCOUNT_ID",
		"ZSB_PARENT_TICKET",
		"ZSB_JIRA_PRIORITY_ID",
		"ZSB_JIRA_LABELS",
	}
)

type issueType struct {
	Label string
	ID    string
}

type config struct {
	baseURL      string
	email        string
	projectKey   string
	reporterID   string
	assigneeID   string
	parentTicket string
	priorityID   string
	labels       []string
	issueTypes   []issueType
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run(args []string, out io.Writer) error {
	if len(args) < 1 || len(args) > 2 {
		return errors.New(usage)
	}

	title := args[0]
	description := ""
	if len(args) == 2 {
		description = args[1]
	}
	if title == "" {
		return errors.New("Title cannot be empty.")
	}

	controlPoint, err := controlPointFromTitle(title)
	if err != nil {
		return err
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	// fzf runs before the token fetch so a cancelled pick never touches pass nor the network.
	selected, err := selectIssueType(cfg.issueTypes)
	if err != nil {
		return err
	}

	token, err := readJiraToken()
	if err != nil {
		return err
	}

	key, err := createJiraTicket(cfg, http.DefaultClient, token, title, controlPoint, description, selected)
	if err != nil {
		return err
	}

	fmt.Fprintln(out, cfg.baseURL+"/browse/"+key)

	parentKey := cfg.parentTicket
	if strings.EqualFold(selected.Label, "epic") {
		parentKey = "xxxxxxxxxxxx"
	}
	home := os.Getenv("HOME")
	if _, err := jira.AppendTicketRow(home+"/temp/tickets", parentKey, key, title); err != nil {
		fmt.Fprintf(os.Stderr, "warning: could not write ticket row: %v\n", err)
	}

	return nil
}

func requireEnv(name string) (string, error) {
	value := os.Getenv(name)
	if value == "" {
		return "", fmt.Errorf("You must set %s first.", name)
	}
	return value, nil
}

func loadConfig() (*config, error) {
	env := make(map[string]string, len(requiredEnvVars))
	for _, name := range requiredEnvVars {
		value, err := requireEnv(name)
		if err != nil {
			return nil, err
		}
		env[name] = value
	}

	baseURL := strings.TrimRight(env["ZSB_JIRA_BASEURL"], "/")
	if baseURL != allowedJiraBaseURL {
		return nil, fmt.Errorf("ZSB_JIRA_BASEURL must be %s", allowedJiraBaseURL)
	}

	issueTypes, err := parseIssueTypeIDs(env["ZSB_JIRA_ISSUE_TYPE_IDS"])
	if err != nil {
		return nil, err
	}

	return &config{
		baseURL:      baseURL,
		email:        env["ZSB_JIRA_EMAIL"],
		projectKey:   env["ZSB_JIRA_PROJECT_KEY"],
		reporterID:   env["ZSB_JIRA_REPORTER_ACCOUNT_ID"],
		assigneeID:   env["ZSB_JIRA_ASSIGNEE_ACCOUNT_ID"],
		parentTicket: env["ZSB_PARENT_TICKET"],
		priorityID:   env["ZSB_JIRA_PRIORITY_ID"],
		labels:       parseLabels(env["ZSB_JIRA_LABELS"]),
		issueTypes:   issueTypes,
	}, nil
}

func parseLabels(value string) []string {
	labels := []string{}
	for entry := range strings.SplitSeq(value, ",") {
		if label := strings.TrimSpace(entry); label != "" {
			labels = append(labels, label)
		}
	}
	return labels
}

// Ordered slice, not a map: fzf --no-sort shows the labels in ZSB_JIRA_ISSUE_TYPE_IDS order.
func parseIssueTypeIDs(value string) ([]issueType, error) {
	var types []issueType
	for entry := range strings.SplitSeq(value, ",") {
		rawLabel, rawID, found := strings.Cut(entry, "=")
		if !found {
			continue
		}
		label, id := strings.TrimSpace(rawLabel), strings.TrimSpace(rawID)
		if label == "" || id == "" {
			continue
		}
		types = append(types, issueType{Label: label, ID: id})
	}
	if len(types) == 0 {
		return nil, errors.New("ZSB_JIRA_ISSUE_TYPE_IDS has no valid entries.")
	}
	return types, nil
}

func controlPointFromTitle(title string) (string, error) {
	matches := controlPointRe.FindAllString(title, -1)
	if len(matches) != 1 || strings.Count(title, "[") != 1 || strings.Count(title, "]") != 1 {
		return "", errors.New("Title must contain exactly one control point in square brackets.")
	}

	controlPoint := strings.TrimSpace(matches[0][1 : len(matches[0])-1])
	if controlPoint == "" {
		return "", errors.New("Control point cannot be empty.")
	}
	return controlPoint, nil
}

func readJiraToken() (string, error) {
	var stderr bytes.Buffer
	cmd := exec.Command("pass", "jira")
	cmd.Stderr = &stderr

	stdout, err := cmd.Output()
	if err != nil {
		message := strings.TrimSpace(stderr.String())
		if message == "" {
			message = "pass jira exited with a non-zero status."
		}
		return "", fmt.Errorf("Token obtain failed: %s", message)
	}

	token := strings.TrimSpace(string(stdout))
	if token == "" {
		return "", errors.New("Token obtain failed.")
	}
	return token, nil
}

func selectIssueType(issueTypes []issueType) (issueType, error) {
	labels := make([]string, len(issueTypes))
	for i, t := range issueTypes {
		labels[i] = t.Label
	}

	cmd := exec.Command("fzf", "--prompt=Issue type: ", "--no-sort")
	cmd.Stdin = strings.NewReader(strings.Join(labels, "\n"))
	cmd.Stderr = os.Stderr // fzf draws on the terminal

	stdout, err := cmd.Output()
	if err != nil {
		return issueType{}, errors.New("Issue type selection cancelled.")
	}

	selected := strings.TrimSpace(string(stdout))
	for _, t := range issueTypes {
		if t.Label == selected {
			return t, nil
		}
	}
	return issueType{}, fmt.Errorf("Unknown issue type: %s", selected)
}

type adfText struct {
	Type string `json:"type"`
	Text string `json:"text"`
}

type adfParagraph struct {
	Type    string    `json:"type"`
	Content []adfText `json:"content"`
}

type adfDoc struct {
	Type    string         `json:"type"`
	Version int            `json:"version"`
	Content []adfParagraph `json:"content"`
}

func buildADF(description string) adfDoc {
	lines := strings.Split(description, "\n")
	content := make([]adfParagraph, 0, len(lines))
	for _, line := range lines {
		// Non-nil empty slice so an empty line marshals as "content": [], not null.
		texts := []adfText{}
		if line != "" {
			texts = append(texts, adfText{Type: "text", Text: line})
		}
		content = append(content, adfParagraph{Type: "paragraph", Content: texts})
	}
	return adfDoc{Type: "doc", Version: 1, Content: content}
}

type keyRef struct {
	Key string `json:"key"`
}

type idRef struct {
	ID string `json:"id"`
}

type accountRef struct {
	AccountID string `json:"accountId"`
}

type fields struct {
	Project      keyRef     `json:"project"`
	IssueType    idRef      `json:"issuetype"`
	Reporter     accountRef `json:"reporter"`
	Assignee     accountRef `json:"assignee"`
	Parent       *keyRef    `json:"parent,omitempty"`
	Priority     idRef      `json:"priority"`
	Labels       []string   `json:"labels"`
	ControlPoint []string   `json:"customfield_25390"`
	Summary      string     `json:"summary"`
	Description  adfDoc     `json:"description"`
}

type payload struct {
	Fields fields `json:"fields"`
}

func buildPayload(cfg *config, title, controlPoint, description string, selected issueType) payload {
	f := fields{
		Project:      keyRef{Key: cfg.projectKey},
		IssueType:    idRef{ID: selected.ID},
		Reporter:     accountRef{AccountID: cfg.reporterID},
		Assignee:     accountRef{AccountID: cfg.assigneeID},
		Priority:     idRef{ID: cfg.priorityID},
		Labels:       cfg.labels,
		ControlPoint: []string{controlPoint},
		Summary:      title,
		Description:  buildADF(description),
	}
	if !strings.EqualFold(selected.Label, "epic") {
		f.Parent = &keyRef{Key: cfg.parentTicket}
	}
	return payload{Fields: f}
}

func parseJiraError(body []byte) string {
	var parsed struct {
		ErrorMessages []string          `json:"errorMessages"`
		Errors        map[string]string `json:"errors"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil {
		return "Unexpected Jira response."
	}

	if len(parsed.ErrorMessages) > 0 {
		return strings.Join(parsed.ErrorMessages, "; ")
	}

	if len(parsed.Errors) > 0 {
		// ponytail: sorted keys, since Go map order is random. Only affects message ordering.
		names := make([]string, 0, len(parsed.Errors))
		for name := range parsed.Errors {
			names = append(names, name)
		}
		slices.Sort(names)

		messages := make([]string, len(names))
		for i, name := range names {
			messages[i] = name + ": " + parsed.Errors[name]
		}
		return strings.Join(messages, "; ")
	}

	return "Unexpected Jira response."
}

func extractIssueKey(body []byte) (string, error) {
	var parsed struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(body, &parsed); err != nil || !jiraKeyRe.MatchString(parsed.Key) {
		return "", errors.New("Jira ticket creation succeeded but no valid key was returned.")
	}
	return parsed.Key, nil
}

func createJiraTicket(cfg *config, client *http.Client, token, title, controlPoint, description string, selected issueType) (string, error) {
	body, err := json.Marshal(buildPayload(cfg, title, controlPoint, description, selected))
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", cfg.baseURL+"/rest/api/3/issue", bytes.NewReader(body))
	if err != nil {
		return "", err
	}
	req.SetBasicAuth(cfg.email, token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusCreated {
		return "", fmt.Errorf("Jira ticket creation failed (%d): %s", resp.StatusCode, parseJiraError(respBody))
	}

	return extractIssueKey(respBody)
}
