package main

import (
	"testing"
)

func TestParseShortFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		valid    bool
	}{
		{"5s", 5, true},
		{"5S", 5, true},
		{"10m", 600, true},
		{"10M", 600, true},
		{"2h", 7200, true},
		{"2H", 7200, true},
		{"invalid", 0, false},
		{"5", 0, false},
		{"5x", 0, false},
		{"", 0, false},
	}

	for _, test := range tests {
		result, valid := parseShortFormat(test.input)
		if valid != test.valid {
			t.Errorf("parseShortFormat(%q) validity = %v, want %v", test.input, valid, test.valid)
		}
		if valid && result != test.expected {
			t.Errorf("parseShortFormat(%q) = %d, want %d", test.input, result, test.expected)
		}
	}
}

func TestParseTimeFormat(t *testing.T) {
	tests := []struct {
		input    string
		expected int
		valid    bool
	}{
		{"05:30", 330, true},      // 5 minutes 30 seconds
		{"1:05:30", 3930, true},   // 1 hour 5 minutes 30 seconds
		{"00:01", 1, true},        // 1 second
		{"59:59", 3599, true},     // 59 minutes 59 seconds
		{"1:00:00", 3600, true},   // 1 hour
		{"invalid", 0, false},
		{"60:00", 0, false},       // Invalid minutes
		{"00:60", 0, false},       // Invalid seconds
		{"", 0, false},
	}

	for _, test := range tests {
		result, valid := parseTimeFormat(test.input)
		if valid != test.valid {
			t.Errorf("parseTimeFormat(%q) validity = %v, want %v", test.input, valid, test.valid)
		}
		if valid && result != test.expected {
			t.Errorf("parseTimeFormat(%q) = %d, want %d", test.input, result, test.expected)
		}
	}
}

func TestFormatTime(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{5, "5"},
		{30, "30"},
		{60, "01:00"},
		{90, "01:30"},
		{3600, "01:00:00"},
		{3661, "01:01:01"},
		{7200, "02:00:00"},
	}

	for _, test := range tests {
		result := formatTime(test.input)
		if result != test.expected {
			t.Errorf("formatTime(%d) = %q, want %q", test.input, result, test.expected)
		}
	}
}

func TestGetCustomTimeMessage(t *testing.T) {
	tests := []struct {
		input    int
		expected string
	}{
		{1, "1 second"},
		{30, "30 seconds"},
		{60, "01:00"},
		{90, "01:30"},
		{3600, "01:00:00"},
	}

	for _, test := range tests {
		result := getCustomTimeMessage(test.input)
		if result != test.expected {
			t.Errorf("getCustomTimeMessage(%d) = %q, want %q", test.input, result, test.expected)
		}
	}
}
