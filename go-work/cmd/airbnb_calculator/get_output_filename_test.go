package main

import (
	"testing"
)

func TestGetOutputFilename(t *testing.T) {
	testCases := []struct {
		filename string
		expected string
	}{
		{"airbnb_01_2024-01_2024.csv", "Jan.csv"},
		{"airbnb_02_2024-02_2024.csv", "Feb.csv"},
		{"airbnb_03_2024-03_2024.csv", "Mar.csv"},
		{"airbnb_04_2024-04_2024.csv", "Apr.csv"},
		{"airbnb_05_2024-05_2024.csv", "May.csv"},
		{"airbnb_06_2024-06_2024.csv", "Jun.csv"},
		{"airbnb_07_2024-07_2024.csv", "Jul.csv"},
		{"airbnb_08_2024-08_2024.csv", "Aug.csv"},
		{"airbnb_09_2024-09_2024.csv", "Sep.csv"},
		{"airbnb_10_2024-10_2024.csv", "Oct.csv"},
		{"airbnb_11_2024-11_2024.csv", "Nov.csv"},
		{"airbnb_12_2024-12_2024.csv", "Dec.csv"},
	}

	for _, tc := range testCases {
		got, _ := GetOutputFilename(tc.filename)
		if got != tc.expected {
			t.Errorf("GetOutputFilename(%v) = %v, expected %v", tc.filename, got, tc.expected)
		}
	}
}
