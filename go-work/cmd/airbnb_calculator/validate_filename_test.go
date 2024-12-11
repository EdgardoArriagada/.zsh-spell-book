package main

import (
	"testing"
)

func TestValidateFilename_WithValidFilename(t *testing.T) {
	testCases := []struct {
		filename string
	}{
		{"airbnb_01_2024-01_2024.csv"},
		{"airbnb_02_2024-02_2024.csv"},
		{"airbnb_03_2024-03_2024.csv"},
	}

	for _, tc := range testCases {
		err := ValidateFilename(tc.filename)
		if err != nil {
			t.Errorf("Failed,  %v", err)
		}
	}
}

func TestValidateFilename_WithInvalidFilenameFormat(t *testing.T) {
	testCases := []struct {
		filename string
	}{
		{"airbnb_01_2024.csv"},
		{"02_2024-02_2024.csv"},
	}

	for _, tc := range testCases {
		err := ValidateFilename(tc.filename)
		if err.Error() != "filename must be in the format airbnb_MM_YYYY-MM_YYYY.csv" {
			t.Errorf("Failed,  %v", err)
		}
	}
}

func TestValidateFilename_WithInvalidFilenameDates(t *testing.T) {
	testCases := []struct {
		filename string
	}{
		{"airbnb_01_2024-01_2025.csv"},
		{"airbnb_02_2024-12_2024.csv"},
	}

	for _, tc := range testCases {
		err := ValidateFilename(tc.filename)
		if err.Error() != "months and years in the filename do not match" {
			t.Errorf("Failed,  %v", err)
		}
	}
}
