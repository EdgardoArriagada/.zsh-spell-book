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
