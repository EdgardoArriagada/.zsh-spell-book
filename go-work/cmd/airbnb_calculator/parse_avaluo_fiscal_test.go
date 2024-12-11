package main

import (
	"testing"
)

func TestParseAvaluoFiscal(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		hasError bool
	}{
		// valid
		{"$ 1.000.000", 1_000_000, false},
		{"1000000", 1_000_000, false},
		{"$1.000.000", 1_000_000, false},
		{"1.000.000", 1_000_000, false},
		{"$ 1 000 000", 1_000_000, false},
		// invalid
		{"invalid", 0, true},
		{"$ 1.000.000,50", 0, true},
	}

	for _, test := range tests {
		actual, err := ParseAvaluoFiscal(test.input)
		if (err != nil) != test.hasError {
			t.Errorf("ParseAvaluoFiscal(%q)\nexpected: %v\nactual: %v", test.input, test.hasError, err)
			continue
		}
		if actual != test.expected {
			t.Errorf("ParseAvaluoFiscal(%q)\nexpected: %v\nactual: %v", test.input, test.expected, actual)
		}
	}
}
