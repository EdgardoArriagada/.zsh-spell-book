package main_test

import (
	"airbnb_calculator"
	"testing"
)

func TestGetOutputFilename(t *testing.T) {
	testCases := []struct {
		monthName string
		expected  string
	}{
		{"Jan", "Jan.csv"},
		{"Feb", "Feb.csv"},
	}

	for _, tc := range testCases {
		got := main.GetOutputFilename(tc.monthName)
		if got != tc.expected {
			t.Errorf("GetOutputFilename(%v) = %v, expected %v", tc.monthName, got, tc.expected)
		}
	}
}
