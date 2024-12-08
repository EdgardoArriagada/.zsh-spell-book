package main

import (
	"testing"
)

func TestProcessRecords_ValidRecords(t *testing.T) {
	records := [][]string{
		{"ID", "Monto", "Noches"},
		{"1", "1000.00", "2"},
		{"2", "2000.00", "3"},
	}
	expected := [][]string{
		{"ID", "Monto", "Noches"},
		{"1", "1000.00", "2"},
		{"2", "2000.00", "3"},
		{"Total:", "3000.00", "5"},
	}
	got, _ := ProcessRecords(records)
	if !equal(got, expected) {
		t.Errorf("ProcessRecords() = %v, expected %v", got, expected)
	}
}

func TestProcessRecords_EmptyMontoValue(t *testing.T) {
	records := [][]string{
		{"ID", "Monto", "Noches"},
		{"1", "", "2"},
		{"2", "2000.00", "3"},
		{"", "", ""},
		{"", "", "10"},
	}
	expected := [][]string{
		{"ID", "Monto", "Noches"},
		{"2", "2000.00", "3"},
		{"Total:", "2000.00", "3"},
	}

	got, _ := ProcessRecords(records)
	if !equal(got, expected) {
		t.Errorf("got %v, expected %v", got, expected)
	}
}

func TestProcessRecords_InvalidNochesValue(t *testing.T) {
	records := [][]string{
		{"ID", "Monto", "Noches"},
		{"1", "1000", "two"},
	}
	_, err := ProcessRecords(records)

	if err == nil {
		t.Errorf("ProcessRecords() error = %v, wantErr %v", err, true)
	}
}

func TestProcessRecords_MissingMontoColumn(t *testing.T) {
	records := [][]string{
		{"ID", "Noches"},
		{"1", "2"},
	}

	_, err := ProcessRecords(records)

	if err.Error() != "Monto or Noches column not found in CSV file" {
		t.Errorf("Test failed")
	}
}

func TestProcessRecords_MissingNochesColumn(t *testing.T) {
	records := [][]string{
		{"ID", "Monto"},
		{"1", "1000"},
	}

	_, err := ProcessRecords(records)

	if err.Error() != "Monto or Noches column not found in CSV file" {
		t.Errorf("Test failed")
	}
}

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

func equal(a, b [][]string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if len(a[i]) != len(b[i]) {
			return false
		}
		for j := range a[i] {
			if a[i][j] != b[i][j] {
				return false
			}
		}
	}
	return true
}
