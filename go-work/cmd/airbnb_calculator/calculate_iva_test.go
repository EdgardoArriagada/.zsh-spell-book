package main

import (
	"testing"
)

func TestCalculateIva_ShouldCalculate(t *testing.T) {
	avaluo_fiscal := 1000000
	days := 30
	total_income := 1000000

	expected := 1000000
	actual := CalculateIva(avaluo_fiscal, days, total_income)

	if actual != expected {
		t.Errorf("\nexpected: %d \nactual: %d", expected, actual)
	}
}
