package main

import (
	"fmt"
	"strconv"
	"strings"
)

var monthNames = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

func GetDateData(filename string) (int, int, string, error) {
	parts1 := strings.Split(filename, "_")
	parts := strings.Split(parts1[2], "-")

	monthStr := parts[1]
	yearStr := parts[0]

	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return 0, 0, "", err
	}

	year, err := strconv.Atoi(yearStr)
	if err != nil {
		return 0, 0, "", err
	}

	if month < 1 || month > 12 {
		return 0, 0, "", fmt.Errorf("Invalid month: %d", month)
	}

	return month, year, monthNames[month-1], nil
}

func GetOutputFilename(monthName string) string {
	return fmt.Sprintf("%s.csv", monthName)
}
