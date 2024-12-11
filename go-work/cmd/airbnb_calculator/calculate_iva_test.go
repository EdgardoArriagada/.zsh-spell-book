package main

import (
	"testing"
)

func TestCalculateIva_ShouldCalculate(t *testing.T) {
	avaluo_fiscal := 56612271
	days := 10
	total_income := 300000

	expected := 21171
	actual := CalculateIva(avaluo_fiscal, days, total_income)

	if actual != expected {
		t.Errorf("\nexpected: %d \nactual: %d", expected, actual)
	}
}
