package main

import "testing"

func TestBaseName(t *testing.T) {
	cases := []struct {
		name, in, want string
	}{
		{"no suffix", "myrepo", "myrepo"},
		{"finished only", "myrepo" + finishedSuffix, "myrepo"},
		{"working only", "myrepo" + workingSuffix, "myrepo"},
		{"both working then finished", "myrepo" + workingSuffix + finishedSuffix, "myrepo"},
		{"both finished then working", "myrepo" + finishedSuffix + workingSuffix, "myrepo"},
		{"double finished", "myrepo" + finishedSuffix + finishedSuffix, "myrepo"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := baseName(c.in); got != c.want {
				t.Errorf("baseName(%q) = %q, want %q", c.in, got, c.want)
			}
		})
	}
}
