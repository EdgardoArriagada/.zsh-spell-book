package jira

import (
	"os"
	"path/filepath"
	"testing"
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

	tickets, err := LoadTickets()
	if err != nil {
		t.Fatal(err)
	}
	if tickets[1].Line != 4 {
		t.Fatalf("line = %d, want 4", tickets[1].Line)
	}
}
