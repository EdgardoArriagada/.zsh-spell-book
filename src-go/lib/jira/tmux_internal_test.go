package jira_test

import (
	"testing"

	"example.com/workspace/lib/jira"
)

func TestParseNotifCounts(t *testing.T) {
	out := "sess-a\t1\n" +
		"sess-a\t2\n" +
		"sess-a\t0\n" +
		"sess-b\t2\n" +
		"sess-b\t2\n" +
		"sess-c\t\n" + // unset/empty → ignored
		"sess-c\t1\n"

	got := jira.ParseNotificationCounts(out)
	want := map[string]jira.NotifCounts{
		"sess-a": {Working: 1, Finished: 1},
		"sess-b": {Working: 2},
		"sess-c": {Finished: 1},
	}

	if len(got) != len(want) {
		t.Fatalf("got %d sessions, want %d: %v", len(got), len(want), got)
	}
	for k, w := range want {
		if got[k] != w {
			t.Errorf("session %q = %+v, want %+v", k, got[k], w)
		}
	}
}

func TestParseNotifCountsEmpty(t *testing.T) {
	if got := jira.ParseNotificationCounts(""); len(got) != 0 {
		t.Errorf("ParseNotificationCounts(\"\") = %v, want empty", got)
	}
}
