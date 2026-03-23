package main_test

import (
	"airbnb_calculator"
	"testing"
)

var CALCULATE_IVA_CALLBACK = func(totalAmount int, totalNights int) float64 {
	return 0.0
}
var AVALUO_FISCAL = 40_000_000

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
		{"-------------", "-------------", "-------------"},
		{"Sum:", "3000.00", "5"},
		{"Avaluo Fiscal:", "40000000"},
		{"Iva:", "0.00"},
	}
	actual, _ := main.ProcessRecords(records, AVALUO_FISCAL, CALCULATE_IVA_CALLBACK)
	if !equal(actual, expected) {
		t.Errorf("\nexpected: %v\nactual %v", expected, actual)
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
		{"-------------", "-------------", "-------------"},
		{"Sum:", "2000.00", "3"},
		{"Avaluo Fiscal:", "40000000"},
		{"Iva:", "0.00"},
	}

	actual, _ := main.ProcessRecords(records, AVALUO_FISCAL, CALCULATE_IVA_CALLBACK)
	if !equal(actual, expected) {
		t.Errorf("\nexpected: %v\nactual: %v", expected, actual)
	}
}

func TestProcessRecords_InvalidNochesValue(t *testing.T) {
	records := [][]string{
		{"ID", "Monto", "Noches"},
		{"1", "1000", "two"},
	}
	_, err := main.ProcessRecords(records, AVALUO_FISCAL, CALCULATE_IVA_CALLBACK)

	if err == nil {
		t.Errorf("ProcessRecords() error = %v, wantErr %v", err, true)
	}
}

func TestProcessRecords_MissingColumns(t *testing.T) {
	tests := []struct {
		records  [][]string
		expected string
	}{
		{
			records: [][]string{
				{"ID", "Noches"},
				{"1", "2"},
			},
			expected: "item Monto not found in slice",
		},
		{
			records: [][]string{
				{"ID", "Monto"},
				{"1", "1000"},
			},
			expected: "item Noches not found in slice",
		},
	}

	for _, tt := range tests {
		t.Run("Missing a column", func(t *testing.T) {
			_, err := main.ProcessRecords(tt.records, AVALUO_FISCAL, CALCULATE_IVA_CALLBACK)
			if err.Error() != tt.expected {
				t.Errorf("\nexpected: %v\nactual: %v", tt.expected, err)
			}
		})
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
