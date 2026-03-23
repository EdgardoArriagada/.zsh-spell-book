package main_test

import (
	"airbnb_calculator"
	"fmt"
	"testing"
)

func TestGetDateData(t *testing.T) {
	testCases := []struct {
		filename string
		expected string
	}{
		{"airbnb_01_2024-01_2024.csv", "31-Jan"},
		{"airbnb_02_2027-02_2027.csv", "28-Feb"},
	}

	for _, tc := range testCases {
		daysInMonth, monthName, _ := main.GetDateData(tc.filename)
		actual := fmt.Sprintf("%d-%s", daysInMonth, monthName)
		if actual != tc.expected {
			t.Errorf("\nexpected: %v\nactual: %v", tc.expected, actual)
		}
	}
}
