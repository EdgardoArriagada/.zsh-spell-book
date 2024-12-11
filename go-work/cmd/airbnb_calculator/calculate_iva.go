package main

import "math"

func CalculateIva(avaluo_fiscal int, days int, total_income int) float64 {
	with_iva := 0.11 * float64(avaluo_fiscal) / (30 * 12) * float64(days)
	base_imponible := (float64(total_income) - with_iva) / 1.19

	iva := base_imponible * 0.19

	// ceil iva
	return math.Ceil(iva)
}
