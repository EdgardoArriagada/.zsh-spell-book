package git

import "testing"

func TestValidateBranchName(t *testing.T) {
	valid := []string{
		"main",
		"feature/my-branch",
		"fix-123",
		"release/v1.0",
		"my_branch",
	}
	for _, name := range valid {
		if err := ValidateBranchName(name); err != nil {
			t.Errorf("ValidateBranchName(%q) returned unexpected error: %v", name, err)
		}
	}

	invalid := []struct {
		name   string
		branch string
	}{
		{"space", "my branch"},
		{"tilde", "feat~1"},
		{"caret", "feat^0"},
		{"colon", "feat:base"},
		{"question mark", "feat?"},
		{"asterisk", "feat*"},
		{"open bracket", "feat[0]"},
		{"backslash", "feat\\bar"},
		{"double dot", "feat..bar"},
		{"at-brace", "feat@{0}"},
		{"starts with dash", "-mybranch"},
		{"starts with dot", ".mybranch"},
		{"ends with dot", "mybranch."},
		{"ends with slash", "mybranch/"},
		{"ends with .lock", "mybranch.lock"},
	}
	for _, tt := range invalid {
		t.Run(tt.name, func(t *testing.T) {
			if err := ValidateBranchName(tt.branch); err == nil {
				t.Errorf("ValidateBranchName(%q) expected error, got nil", tt.branch)
			}
		})
	}
}
