package main

import (
	"fmt"
	"testing"
)

func TestIsUnmergedBranchError(t *testing.T) {
	tests := []struct {
		err    error
		expect bool
	}{
		{nil, false},
		{fmt.Errorf("some other error"), false},
		{fmt.Errorf("error: The branch 'feat' is not fully merged"), true},
		{fmt.Errorf("is not fully merged"), true},
	}
	for _, tt := range tests {
		got := isUnmergedBranchError(tt.err)
		if got != tt.expect {
			t.Errorf("isUnmergedBranchError(%v) = %v, want %v", tt.err, got, tt.expect)
		}
	}
}
