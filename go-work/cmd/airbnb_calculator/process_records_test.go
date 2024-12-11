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
