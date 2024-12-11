package main

import "math"

func CalculateIva(avaluoFiscal int, totalIncome int, nights int, daysInMonth int) float64 {
	with_iva := 0.11 * float64(avaluoFiscal) / (float64(daysInMonth) * 12) * float64(nights)
	base_imponible := (float64(totalIncome) - with_iva) / 1.19

	iva := base_imponible * 0.19

	if iva < 0 {
		return 0
	}

	return math.Ceil(iva)
}
