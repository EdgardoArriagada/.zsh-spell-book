package main

import "math"

var daysPerMonth = map[int]int{
	1:  31,
	2:  28,
	3:  31,
	4:  30,
	5:  31,
	6:  30,
	7:  31,
	8:  31,
	9:  30,
	10: 31,
	11: 30,
	12: 31,
}

// make sure Feb has 29 days for leap years
func getDaysInMonth(month int, year int) int {
	if month == 2 && isLeapYear(year) {
		return 29
	}
	return daysPerMonth[month]
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func CalculateIva(avaluo_fiscal int, total_income int, days int, month int, year int) float64 {
	days_in_month := getDaysInMonth(month, year)
	with_iva := 0.11 * float64(avaluo_fiscal) / (float64(days_in_month) * 12) * float64(days)
	base_imponible := (float64(total_income) - with_iva) / 1.19

	iva := base_imponible * 0.19

	if iva < 0 {
		return 0
	}

	return math.Ceil(iva)
}
