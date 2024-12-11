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

func TestGetMonthAndYearFromFilename(t *testing.T) {
	testCases := []struct {
		filename string
		expected string
	}{
		{"airbnb_01_2024-01_2024.csv", "1-2024-Jan"},
		{"airbnb_02_2027-02_2027.csv", "2-2027-Feb"},
	}

	for _, tc := range testCases {
		month, year, monthName, _ := GetMonthAndYearFromFilename(tc.filename)
		got := fmt.Sprintf("%d-%d-%s", month, year, monthName)
		if got != tc.expected {
			t.Errorf("GetOutputFilename(%v) = %v, expected %v", tc.filename, got, tc.expected)
		}
	}
}
