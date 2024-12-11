package main

import (
	"testing"
)

func TestCalculateIva_ShouldCalculate(t *testing.T) {
	avaluo_fiscal := 83_231_614
	days := 10
	total_income := 300_000

	expected := float64(7_294)
	actual := CalculateIva(avaluo_fiscal, days, total_income)

	if actual != expected {
		t.Errorf("\nexpected: %f \nactual: %f", expected, actual)
	}
}
