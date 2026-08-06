package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseLabels(t *testing.T) {
	tests := []struct {
		name  string
		input string
		want  []string
	}{
		{name: "trims and drops empties", input: " foo,bar,, baz ", want: []string{"foo", "bar", "baz"}},
		{name: "single label", input: "foo", want: []string{"foo"}},
		{name: "empty value", input: "", want: []string{}},
		{name: "only separators", input: ",,", want: []string{}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseLabels(tt.input)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseLabels(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestParseIssueTypeIDs(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []issueType
		wantErr bool
	}{
		{
			name:  "preserves order",
			input: "Epic=10000,Task=10101,Sub-task=10102",
			want: []issueType{
				{Label: "Epic", ID: "10000"},
				{Label: "Task", ID: "10101"},
				{Label: "Sub-task", ID: "10102"},
			},
		},
		{name: "trims whitespace", input: " Epic = 10000 ", want: []issueType{{Label: "Epic", ID: "10000"}}},
		{name: "skips entry without label", input: "=10,Task=10101", want: []issueType{{Label: "Task", ID: "10101"}}},
		{name: "skips entry without equals", input: "Epic,Task=10101", want: []issueType{{Label: "Task", ID: "10101"}}},
		{name: "skips entry without id", input: "Epic=,Task=10101", want: []issueType{{Label: "Task", ID: "10101"}}},
		{name: "errors on empty value", input: "", wantErr: true},
		{name: "errors when nothing valid", input: "Epic,Task", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseIssueTypeIDs(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("parseIssueTypeIDs(%q) = %v, want error", tt.input, got)
				}
				if want := "ZSB_JIRA_ISSUE_TYPE_IDS has no valid entries."; err.Error() != want {
					t.Errorf("err = %q, want %q", err.Error(), want)
				}
				return
			}
			if err != nil {
				t.Fatalf("parseIssueTypeIDs(%q) returned error %v", tt.input, err)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseIssueTypeIDs(%q) = %v, want %v", tt.input, got, tt.want)
			}
		})
	}
}

func TestControlPointFromTitle(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    string
		wantErr string
	}{
		{name: "single control point", input: "[Checkout] fix the thing", want: "Checkout"},
		{name: "trims inner whitespace", input: "[  Checkout  ] fix", want: "Checkout"},
		{name: "control point at the end", input: "fix the thing [Checkout]", want: "Checkout"},
		{name: "no control point", input: "fix the thing", wantErr: "Title must contain exactly one control point in square brackets."},
		{name: "two control points", input: "[A] [B] fix", wantErr: "Title must contain exactly one control point in square brackets."},
		{name: "unbalanced bracket", input: "[A] fix ] thing", wantErr: "Title must contain exactly one control point in square brackets."},
		{name: "empty brackets", input: "[] fix", wantErr: "Title must contain exactly one control point in square brackets."},
		{name: "whitespace only control point", input: "[   ] fix", wantErr: "Control point cannot be empty."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := controlPointFromTitle(tt.input)
			if tt.wantErr != "" {
				if err == nil {
					t.Fatalf("controlPointFromTitle(%q) = %q, want error", tt.input, got)
				}
				if err.Error() != tt.wantErr {
					t.Errorf("err = %q, want %q", err.Error(), tt.wantErr)
				}
				return
			}
			if err != nil {
				t.Fatalf("controlPointFromTitle(%q) returned error %v", tt.input, err)
			}
			if got != tt.want {
				t.Errorf("controlPointFromTitle(%q) = %q, want %q", tt.input, got, tt.want)
			}
		})
	}
}

func TestBuildADF(t *testing.T) {
	got, err := json.Marshal(buildADF("Line one\n\nLine three"))
	if err != nil {
		t.Fatal(err)
	}

	want := `{"type":"doc","version":1,"content":[` +
		`{"type":"paragraph","content":[{"type":"text","text":"Line one"}]},` +
		`{"type":"paragraph","content":[]},` +
		`{"type":"paragraph","content":[{"type":"text","text":"Line three"}]}]}`

	if string(got) != want {
		t.Errorf("buildADF() = %s, want %s", got, want)
	}
}

func TestBuildADFEmptyDescription(t *testing.T) {
	got, err := json.Marshal(buildADF(""))
	if err != nil {
		t.Fatal(err)
	}

	want := `{"type":"doc","version":1,"content":[{"type":"paragraph","content":[]}]}`
	if string(got) != want {
		t.Errorf("buildADF(\"\") = %s, want %s", got, want)
	}
}

func testConfig() *config {
	return &config{
		baseURL:      allowedJiraBaseURL,
		email:        "user@example.com",
		projectKey:   "ABC",
		reporterID:   "reporter-1",
		assigneeID:   "assignee-1",
		parentTicket: "ABC-1",
		priorityID:   "3",
		labels:       []string{"foo", "bar"},
		issueTypes:   []issueType{{Label: "Epic", ID: "10000"}, {Label: "Task", ID: "10101"}},
	}
}

func marshalPayload(t *testing.T, selected issueType) map[string]any {
	t.Helper()

	raw, err := json.Marshal(buildPayload(testConfig(), "[Checkout] fix the thing", "Checkout", "desc", selected))
	if err != nil {
		t.Fatal(err)
	}

	var decoded struct {
		Fields map[string]any `json:"fields"`
	}
	if err := json.Unmarshal(raw, &decoded); err != nil {
		t.Fatal(err)
	}
	return decoded.Fields
}

func TestBuildPayloadNonEpicIncludesParent(t *testing.T) {
	f := marshalPayload(t, issueType{Label: "Task", ID: "10101"})

	parent, ok := f["parent"].(map[string]any)
	if !ok {
		t.Fatalf("parent = %v, want an object", f["parent"])
	}
	if parent["key"] != "ABC-1" {
		t.Errorf("parent.key = %v, want %q", parent["key"], "ABC-1")
	}

	issuetype := f["issuetype"].(map[string]any)
	if issuetype["id"] != "10101" {
		t.Errorf("issuetype.id = %v, want %q", issuetype["id"], "10101")
	}
	if f["summary"] != "[Checkout] fix the thing" {
		t.Errorf("summary = %v, want the title verbatim", f["summary"])
	}
	if got := f[controlPointField]; !reflect.DeepEqual(got, []any{"Checkout"}) {
		t.Errorf("%s = %v, want [Checkout]", controlPointField, got)
	}
	if got := f["labels"]; !reflect.DeepEqual(got, []any{"foo", "bar"}) {
		t.Errorf("labels = %v, want [foo bar]", got)
	}
}

func TestBuildPayloadEpicOmitsParent(t *testing.T) {
	tests := []struct {
		name  string
		label string
	}{
		{name: "exact case", label: "Epic"},
		{name: "lowercase", label: "epic"},
		{name: "mixed case", label: "ePiC"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := marshalPayload(t, issueType{Label: tt.label, ID: "10000"})
			if _, ok := f["parent"]; ok {
				t.Errorf("parent = %v, want it omitted for %q", f["parent"], tt.label)
			}
		})
	}
}

func TestParseJiraError(t *testing.T) {
	tests := []struct {
		name string
		body string
		want string
	}{
		{name: "errorMessages joined", body: `{"errorMessages":["boom","bang"]}`, want: "boom; bang"},
		{name: "errors sorted by field", body: `{"errors":{"summary":"required","assignee":"invalid"}}`, want: "assignee: invalid; summary: required"},
		{name: "errorMessages wins over errors", body: `{"errorMessages":["boom"],"errors":{"summary":"required"}}`, want: "boom"},
		{name: "empty errorMessages falls through to errors", body: `{"errorMessages":[],"errors":{"summary":"required"}}`, want: "summary: required"},
		{name: "both empty", body: `{"errorMessages":[],"errors":{}}`, want: "Unexpected Jira response."},
		{name: "malformed json", body: `not json`, want: "Unexpected Jira response."},
		{name: "array body", body: `["nope"]`, want: "Unexpected Jira response."},
		{name: "empty body", body: ``, want: "Unexpected Jira response."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseJiraError([]byte(tt.body)); got != tt.want {
				t.Errorf("parseJiraError(%q) = %q, want %q", tt.body, got, tt.want)
			}
		})
	}
}

func TestExtractIssueKey(t *testing.T) {
	tests := []struct {
		name    string
		body    string
		want    string
		wantErr bool
	}{
		{name: "valid key", body: `{"key":"ABC-123"}`, want: "ABC-123"},
		{name: "alphanumeric project key", body: `{"key":"AB1C-9"}`, want: "AB1C-9"},
		{name: "missing key", body: `{"id":"1"}`, wantErr: true},
		{name: "lowercase key", body: `{"key":"abc-123"}`, wantErr: true},
		{name: "no number", body: `{"key":"ABC-"}`, wantErr: true},
		{name: "malformed json", body: `not json`, wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := extractIssueKey([]byte(tt.body))
			if tt.wantErr {
				if err == nil {
					t.Fatalf("extractIssueKey(%q) = %q, want error", tt.body, got)
				}
				if want := "Jira ticket creation succeeded but no valid key was returned."; err.Error() != want {
					t.Errorf("err = %q, want %q", err.Error(), want)
				}
				return
			}
			if err != nil {
				t.Fatalf("extractIssueKey(%q) returned error %v", tt.body, err)
			}
			if got != tt.want {
				t.Errorf("extractIssueKey(%q) = %q, want %q", tt.body, got, tt.want)
			}
		})
	}
}

func TestLoadConfig(t *testing.T) {
	setEnv := func(t *testing.T, overrides map[string]string) {
		t.Helper()
		defaults := map[string]string{
			"ZSB_JIRA_BASEURL":             allowedJiraBaseURL + "/",
			"ZSB_JIRA_EMAIL":               "user@example.com",
			"ZSB_JIRA_PROJECT_KEY":         "ABC",
			"ZSB_JIRA_ISSUE_TYPE_IDS":      "Epic=10000,Task=10101",
			"ZSB_JIRA_REPORTER_ACCOUNT_ID": "reporter-1",
			"ZSB_JIRA_ASSIGNEE_ACCOUNT_ID": "assignee-1",
			"ZSB_PARENT_TICKET":            "ABC-1",
			"ZSB_JIRA_PRIORITY_ID":         "3",
			"ZSB_JIRA_LABELS":              " foo,bar,, baz ",
		}
		for name, value := range overrides {
			defaults[name] = value
		}
		for name, value := range defaults {
			t.Setenv(name, value)
		}
	}

	t.Run("happy path normalizes trailing slashes", func(t *testing.T) {
		setEnv(t, map[string]string{"ZSB_JIRA_BASEURL": allowedJiraBaseURL + "///"})

		cfg, err := loadConfig()
		if err != nil {
			t.Fatalf("loadConfig() returned error %v", err)
		}
		if cfg.baseURL != allowedJiraBaseURL {
			t.Errorf("baseURL = %q, want %q", cfg.baseURL, allowedJiraBaseURL)
		}
		if want := []string{"foo", "bar", "baz"}; !reflect.DeepEqual(cfg.labels, want) {
			t.Errorf("labels = %v, want %v", cfg.labels, want)
		}
		if want := []issueType{{Label: "Epic", ID: "10000"}, {Label: "Task", ID: "10101"}}; !reflect.DeepEqual(cfg.issueTypes, want) {
			t.Errorf("issueTypes = %v, want %v", cfg.issueTypes, want)
		}
	})

	t.Run("missing env var", func(t *testing.T) {
		setEnv(t, map[string]string{"ZSB_JIRA_EMAIL": ""})

		_, err := loadConfig()
		if err == nil {
			t.Fatal("loadConfig() = nil error, want error")
		}
		if want := "You must set ZSB_JIRA_EMAIL first."; err.Error() != want {
			t.Errorf("err = %q, want %q", err.Error(), want)
		}
	})

	t.Run("disallowed base url", func(t *testing.T) {
		setEnv(t, map[string]string{"ZSB_JIRA_BASEURL": "https://evil.example.com"})

		_, err := loadConfig()
		if err == nil {
			t.Fatal("loadConfig() = nil error, want error")
		}
		if want := "ZSB_JIRA_BASEURL must be " + allowedJiraBaseURL; err.Error() != want {
			t.Errorf("err = %q, want %q", err.Error(), want)
		}
	})

	t.Run("invalid issue type ids", func(t *testing.T) {
		setEnv(t, map[string]string{"ZSB_JIRA_ISSUE_TYPE_IDS": "Epic,Task"})

		_, err := loadConfig()
		if err == nil {
			t.Fatal("loadConfig() = nil error, want error")
		}
		if want := "ZSB_JIRA_ISSUE_TYPE_IDS has no valid entries."; err.Error() != want {
			t.Errorf("err = %q, want %q", err.Error(), want)
		}
	})
}

func TestCreateJiraTicket(t *testing.T) {
	t.Run("posts the payload and returns the key", func(t *testing.T) {
		var gotMethod, gotPath, gotAuth, gotContentType string
		var gotBody map[string]any

		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			gotMethod, gotPath = r.Method, r.URL.Path
			gotAuth = r.Header.Get("Authorization")
			gotContentType = r.Header.Get("Content-Type")
			if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
				t.Errorf("decode request body: %v", err)
			}
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{"key":"ABC-123"}`)) //nolint:errcheck
		}))
		defer srv.Close()

		cfg := testConfig()
		cfg.baseURL = srv.URL
		key, err := createJiraTicket(cfg, http.DefaultClient, "s3cret", "[Checkout] fix", "Checkout", "desc", issueType{Label: "Task", ID: "10101"})
		if err != nil {
			t.Fatalf("createJiraTicket() returned error %v", err)
		}
		if key != "ABC-123" {
			t.Errorf("key = %q, want %q", key, "ABC-123")
		}
		if gotMethod != "POST" {
			t.Errorf("method = %q, want POST", gotMethod)
		}
		if gotPath != "/rest/api/3/issue" {
			t.Errorf("path = %q, want /rest/api/3/issue", gotPath)
		}
		// base64("user@example.com:s3cret")
		if want := "Basic dXNlckBleGFtcGxlLmNvbTpzM2NyZXQ="; gotAuth != want {
			t.Errorf("Authorization = %q, want %q", gotAuth, want)
		}
		if gotContentType != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", gotContentType)
		}
		if _, ok := gotBody["fields"]; !ok {
			t.Errorf("body = %v, want a fields object", gotBody)
		}
	})

	t.Run("non-201 surfaces the jira error", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(`{"errorMessages":["boom"]}`)) //nolint:errcheck
		}))
		defer srv.Close()

		cfg := testConfig()
		cfg.baseURL = srv.URL
		_, err := createJiraTicket(cfg, http.DefaultClient, "s3cret", "[Checkout] fix", "Checkout", "", issueType{Label: "Task", ID: "10101"})
		if err == nil {
			t.Fatal("createJiraTicket() = nil error, want error")
		}
		if want := "Jira ticket creation failed (400): boom"; err.Error() != want {
			t.Errorf("err = %q, want %q", err.Error(), want)
		}
	})

	t.Run("201 without a valid key", func(t *testing.T) {
		srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusCreated)
			w.Write([]byte(`{}`)) //nolint:errcheck
		}))
		defer srv.Close()

		cfg := testConfig()
		cfg.baseURL = srv.URL
		_, err := createJiraTicket(cfg, http.DefaultClient, "s3cret", "[Checkout] fix", "Checkout", "", issueType{Label: "Task", ID: "10101"})
		if err == nil {
			t.Fatal("createJiraTicket() = nil error, want error")
		}
		if want := "Jira ticket creation succeeded but no valid key was returned."; err.Error() != want {
			t.Errorf("err = %q, want %q", err.Error(), want)
		}
	})
}

func TestRunArgValidation(t *testing.T) {
	tests := []struct {
		name string
		args []string
		want string
	}{
		{name: "no args", args: []string{}, want: usage},
		{name: "three args", args: []string{"a", "b", "c"}, want: usage},
		{name: "empty title", args: []string{""}, want: "Title cannot be empty."},
		{name: "title without control point", args: []string{"fix the thing"}, want: "Title must contain exactly one control point in square brackets."},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := run(tt.args, io.Discard)
			if err == nil {
				t.Fatalf("run(%v) = nil error, want error", tt.args)
			}
			if err.Error() != tt.want {
				t.Errorf("err = %q, want %q", err.Error(), tt.want)
			}
		})
	}
}

func TestUsageString(t *testing.T) {
	if !strings.HasPrefix(usage, "Usage: create-jira-ticket ") {
		t.Errorf("usage = %q, want it to start with the command name", usage)
	}
}

// fakeBin creates a fake shell executable in a temp dir and prepends it to PATH.
func fakeBin(t *testing.T, name, script string) {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, name), []byte("#!/bin/sh\n"+script), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("PATH", dir+":"+os.Getenv("PATH"))
}

func TestReadJiraToken(t *testing.T) {
	t.Run("returns trimmed token", func(t *testing.T) {
		fakeBin(t, "pass", `printf "  mytoken\n"`)
		got, err := readJiraToken()
		if err != nil {
			t.Fatalf("readJiraToken() returned error %v", err)
		}
		if got != "mytoken" {
			t.Errorf("got %q, want %q", got, "mytoken")
		}
	})

	t.Run("empty output fails", func(t *testing.T) {
		fakeBin(t, "pass", `exit 0`)
		_, err := readJiraToken()
		if err == nil {
			t.Fatal("readJiraToken() = nil error, want error")
		}
		if want := "Token obtain failed."; err.Error() != want {
			t.Errorf("err = %q, want %q", err.Error(), want)
		}
	})

	t.Run("non-zero exit with stderr", func(t *testing.T) {
		fakeBin(t, "pass", `echo "Error: not found" >&2; exit 1`)
		_, err := readJiraToken()
		if err == nil {
			t.Fatal("readJiraToken() = nil error, want error")
		}
		if want := "Token obtain failed: Error: not found"; err.Error() != want {
			t.Errorf("err = %q, want %q", err.Error(), want)
		}
	})

	t.Run("non-zero exit without stderr uses fallback", func(t *testing.T) {
		fakeBin(t, "pass", `exit 1`)
		_, err := readJiraToken()
		if err == nil {
			t.Fatal("readJiraToken() = nil error, want error")
		}
		if want := "Token obtain failed: pass jira exited with a non-zero status."; err.Error() != want {
			t.Errorf("err = %q, want %q", err.Error(), want)
		}
	})
}

func TestSelectIssueType(t *testing.T) {
	types := []issueType{{Label: "Epic", ID: "10000"}, {Label: "Task", ID: "10101"}}

	t.Run("returns selected type", func(t *testing.T) {
		fakeBin(t, "fzf", `echo "Task"`)
		got, err := selectIssueType(types)
		if err != nil {
			t.Fatalf("selectIssueType() returned error %v", err)
		}
		if got != (issueType{Label: "Task", ID: "10101"}) {
			t.Errorf("got %+v, want {Task 10101}", got)
		}
	})

	t.Run("cancellation returns error", func(t *testing.T) {
		fakeBin(t, "fzf", `exit 130`)
		_, err := selectIssueType(types)
		if err == nil {
			t.Fatal("selectIssueType() = nil error, want error")
		}
		if want := "Issue type selection cancelled."; err.Error() != want {
			t.Errorf("err = %q, want %q", err.Error(), want)
		}
	})
}

func TestRunEndToEnd(t *testing.T) {
	var gotBody map[string]any

	srv := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := json.NewDecoder(r.Body).Decode(&gotBody); err != nil {
			t.Errorf("decode request body: %v", err)
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(`{"key":"ABC-123"}`)) //nolint:errcheck
	}))
	defer srv.Close()

	// Route all TCP to the test server; skip TLS verification so we can use allowedJiraBaseURL.
	origTransport := http.DefaultTransport
	defer func() { http.DefaultTransport = origTransport }()
	http.DefaultTransport = &http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (net.Conn, error) {
			return (&net.Dialer{}).DialContext(ctx, network, srv.Listener.Addr().String())
		},
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true}, //nolint:gosec // test-only transport
	}

	t.Setenv("ZSB_JIRA_BASEURL", allowedJiraBaseURL)
	t.Setenv("ZSB_JIRA_EMAIL", "user@example.com")
	t.Setenv("ZSB_JIRA_PROJECT_KEY", "ABC")
	t.Setenv("ZSB_JIRA_ISSUE_TYPE_IDS", "Epic=10000,Task=10101")
	t.Setenv("ZSB_JIRA_REPORTER_ACCOUNT_ID", "reporter-1")
	t.Setenv("ZSB_JIRA_ASSIGNEE_ACCOUNT_ID", "assignee-1")
	t.Setenv("ZSB_PARENT_TICKET", "ABC-1")
	t.Setenv("ZSB_JIRA_PRIORITY_ID", "3")
	t.Setenv("ZSB_JIRA_LABELS", "foo,bar")

	fakeBin(t, "pass", `echo "s3cret"`)
	fakeBin(t, "fzf", `echo "Task"`)

	var buf bytes.Buffer
	runErr := run([]string{"[Checkout] fix the thing", "desc"}, &buf)

	if runErr != nil {
		t.Fatalf("run() returned error %v", runErr)
	}
	if want := allowedJiraBaseURL + "/browse/ABC-123"; strings.TrimSpace(buf.String()) != want {
		t.Errorf("stdout = %q, want %q", strings.TrimSpace(buf.String()), want)
	}
	if _, ok := gotBody["fields"]; !ok {
		t.Errorf("request body = %v, want a fields object", gotBody)
	}
}
