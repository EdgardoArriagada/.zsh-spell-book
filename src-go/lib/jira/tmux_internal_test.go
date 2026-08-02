package jira

import "testing"

func TestParseNotifCounts(t *testing.T) {
	out := "sess-a\t1\n" +
		"sess-a\t2\n" +
		"sess-a\t0\n" +
		"sess-b\t2\n" +
		"sess-b\t2\n" +
		"sess-c\t\n" + // unset/empty → ignored
		"sess-c\t1\n"

	got := parseNotifCounts(out)
	want := map[string]NotifCounts{
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
	if got := parseNotifCounts(""); len(got) != 0 {
		t.Errorf("parseNotifCounts(\"\") = %v, want empty", got)
	}
}
