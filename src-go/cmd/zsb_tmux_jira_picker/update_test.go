package main

import (
	"slices"
	"testing"
)

func TestEditTicketsCmdStartsNvimAtTicketLine(t *testing.T) {
	cmd := editTicketsCmd("nvim --clean", "/tmp/tickets", 4)
	if got, want := cmd.Args, []string{"nvim", "--clean", "+4", "/tmp/tickets"}; !slices.Equal(got, want) {
		t.Fatalf("args = %q, want %q", got, want)
	}
}

func TestEditTicketsCmdKeepsOtherEditors(t *testing.T) {
	cmd := editTicketsCmd("vim", "/tmp/tickets", 4)
	if got, want := cmd.Args, []string{"vim", "/tmp/tickets"}; !slices.Equal(got, want) {
		t.Fatalf("args = %q, want %q", got, want)
	}
}
