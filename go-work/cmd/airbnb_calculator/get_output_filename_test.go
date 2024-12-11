package main

import (
	"fmt"
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
		got := GetOutputFilename(tc.monthName)
		if got != tc.expected {
			t.Errorf("GetOutputFilename(%v) = %v, expected %v", tc.monthName, got, tc.expected)
		}
	}
}

func TestGetDateData(t *testing.T) {
	testCases := []struct {
		filename string
		expected string
	}{
		{"airbnb_01_2024-01_2024.csv", "31-Jan"},
		{"airbnb_02_2027-02_2027.csv", "28-Feb"},
	}

	for _, tc := range testCases {
		daysInMonth, monthName, _ := GetDateData(tc.filename)
		got := fmt.Sprintf("%d-%s", daysInMonth, monthName)
		if got != tc.expected {
			t.Errorf("\nexpected: %v\nactual: %v", tc.expected, got)
		}
	}
}
