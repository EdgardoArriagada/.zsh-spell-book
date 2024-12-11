package main

import (
	"testing"
)

func TestCalculateIva_ShouldCalculate(t *testing.T) {
	tests := []struct {
		avaluoFiscal int
		days         int
		totalIncome  int
		expected     float64
	}{
		{83_231_614, 10, 300_000, 7_294},
		{83_231_614, 15, 500_000, 18_924},
		{56_612_271, 10, 300_000, 20_281},
		{48_738_853, 3, 115716, 11_343},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			actual := CalculateIva(tt.avaluoFiscal, tt.days, tt.totalIncome)
			if actual != tt.expected {
				t.Errorf("\nexpected: %f \nactual: %f", tt.expected, actual)
			}
		})
	}
}
