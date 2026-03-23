package main

import (
	"fmt"
	"strconv"
	"strings"
)

var monthNames = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

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
func GetDaysInMonth(month int, year int) int {
	if month == 2 && isLeapYear(year) {
		return 29
	}
	return daysPerMonth[month]
}

func isLeapYear(year int) bool {
	return year%4 == 0 && (year%100 != 0 || year%400 == 0)
}

func GetDateData(filename string) (int, string, error) {
	parts1 := strings.Split(filename, "_")
	parts := strings.Split(parts1[2], "-")

	monthStr := parts[1]
	yearStr := parts[0]

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return 0, "", err
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, "", err
	}

	if month < 1 || month > 12 {
		return 0, "", fmt.Errorf("Invalid month: %d", month)
	}

	daysInMonth := GetDaysInMonth(month, year)

	return daysInMonth, monthNames[month-1], nil
}
