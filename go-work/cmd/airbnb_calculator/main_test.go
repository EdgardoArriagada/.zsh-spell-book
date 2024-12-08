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
	}
	expected := [][]string{
		{"ID", "Monto", "Noches"},
		{"2", "2000.00", "3"},
		{"Total:", "2000.00", "3"},
	}

	got, _ := ProcessRecords(records)
	if !equal(got, expected) {
		t.Errorf("ProcessRecords() = %v, expected %v", got, expected)
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

func TestGetOutputFileName(t *testing.T) {
	testCases := []struct {
		fileName string
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
		got, _ := GetOutputFilename(tc.fileName)
		if got != tc.expected {
			t.Errorf("GetOutputFilename(%v) = %v, expected %v", tc.fileName, got, tc.expected)
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
