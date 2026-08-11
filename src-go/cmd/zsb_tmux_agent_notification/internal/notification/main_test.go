package notification

import "testing"

func TestBaseName(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"no suffix", "myrepo", "myrepo"},
		{"finished only", "myrepo" + FINISHED_SUFFIX, "myrepo"},
		{"working only", "myrepo" + WORKING_SUFFIX, "myrepo"},
		{"manual only", "myrepo" + MANUAL_SUFFIX, "myrepo"},
		{"both working then finished", "myrepo" + WORKING_SUFFIX + FINISHED_SUFFIX, "myrepo"},
		{"both finished then working", "myrepo" + FINISHED_SUFFIX + WORKING_SUFFIX, "myrepo"},
		{"double finished", "myrepo" + FINISHED_SUFFIX + FINISHED_SUFFIX, "myrepo"},
		{"all three suffixes", "myrepo" + WORKING_SUFFIX + FINISHED_SUFFIX + MANUAL_SUFFIX, "myrepo"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := baseName(c.in); got != c.want {
				t.Errorf("baseName(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
