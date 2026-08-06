package jira_test

import (
	"os"
	"path/filepath"
	"slices"
	"testing"

	"example.com/workspace/lib/jira"
)

func TestLoadTicketsKeepsSourceLine(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.Mkdir(filepath.Join(home, "temp"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, "temp", "tickets"), []byte("# note\n\nparent | JIRA-1 | First\nparent | JIRA-2 | Second\n"), 0644); err != nil {
		t.Fatal(err)
	}

	tickets, err := jira.LoadTickets()
	if err != nil {
		t.Fatal(err)
	}
	if tickets[1].Line != 4 {
		t.Fatalf("line = %d, want 4", tickets[1].Line)
	}
}

func TestLoadTicketsAssignsSectionTitles(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.Mkdir(filepath.Join(home, "temp"), 0755); err != nil {
		t.Fatal(err)
	}
	content := "parent | JIRA-1 | Untitled\n# comment\nPersonal\nparent | JIRA-2 | First\nWork\nparent | JIRA-3 | Second\n"
	if err := os.WriteFile(filepath.Join(home, "temp", "tickets"), []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	tickets, err := jira.LoadTickets()
	if err != nil {
		t.Fatal(err)
	}
	if got, want := []string{tickets[0].Title, tickets[1].Title, tickets[2].Title}, []string{"", "Personal", "Work"}; !slices.Equal(got, want) {
		t.Fatalf("titles = %q, want %q", got, want)
	}
}

func TestLoadTicketsKeepsPipesInLabel(t *testing.T) {
	home := t.TempDir()
	t.Setenv("HOME", home)
	if err := os.Mkdir(filepath.Join(home, "temp"), 0755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(home, "temp", "tickets"), []byte("parent | JIRA-1 | First | Second\n"), 0644); err != nil {
		t.Fatal(err)
	}

	tickets, err := jira.LoadTickets()
	if err != nil {
		t.Fatal(err)
	}
	if got := tickets[0].Label; got != "First | Second" {
		t.Fatalf("label = %q", got)
	}
}

func writeTickets(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), "tickets")
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestAppendTicketRowCreatesUnmanagedSection(t *testing.T) {
	path := writeTickets(t, "Current\nP|JIRA-1|First\n")

	if _, err := jira.AppendTicketRow(path, "P", "JIRA-2", "Second"); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "Current\nP|JIRA-1|First\nUnmanaged\nP|JIRA-2|Second\n"
	if string(got) != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAppendTicketRowInsertsAfterUnmanagedHeader(t *testing.T) {
	path := writeTickets(t, "Current\nP|JIRA-1|First\nUnmanaged\nOther\nP|JIRA-3|Third\n")

	if _, err := jira.AppendTicketRow(path, "P", "JIRA-2", "Second"); err != nil {
		t.Fatal(err)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	want := "Current\nP|JIRA-1|First\nUnmanaged\nP|JIRA-2|Second\nOther\nP|JIRA-3|Third\n"
	if string(got) != want {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestAppendTicketRowSkipsDuplicate(t *testing.T) {
	path := writeTickets(t, "P|JIRA-1|First\n")

	lnum, err := jira.AppendTicketRow(path, "P", "JIRA-1", "First")
	if err != nil {
		t.Fatal(err)
	}
	if lnum != 1 {
		t.Fatalf("lnum = %d, want 1", lnum)
	}

	got, _ := os.ReadFile(path)
	if string(got) != "P|JIRA-1|First\n" {
		t.Fatalf("file modified unexpectedly: %q", got)
	}
}
