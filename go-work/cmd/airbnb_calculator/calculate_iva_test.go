package main

import (
	"testing"
)

func TestCalculateIva_ShouldCalculate(t *testing.T) {
	tests := []struct {
		avaluoFiscal int
		totalIncome  int
		days         int
		month        int
		year         int
		expected     float64
	}{
		{83_231_614, 300_000, 10, 11, 2024, 7_294},
		{83_231_614, 500_000, 15, 11, 2024, 18_924},
		{56_612_271, 300_000, 10, 1, 2024, 21_172},
		{48_738_853, 115716, 3, 11, 2024, 11_343},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := CalculateIva(tt.avaluoFiscal, tt.totalIncome, tt.days, tt.month, tt.year)
			if actual != tt.expected {
				t.Errorf("\nexpected: %f \nactual: %f", tt.expected, actual)
			}
		})
	}
}
