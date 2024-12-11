package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

var monthNames = []string{"Jan", "Feb", "Mar", "Apr", "May", "Jun", "Jul", "Aug", "Sep", "Oct", "Nov", "Dec"}

func GetOutputFilename(filename string) (string, error) {
	// Extract the base name without extension
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	parts := strings.Split(name, "_")
	monthStr := parts[1]

	// Convert month string to integer
	month, err := strconv.Atoi(monthStr)
	if err != nil {
		return "", err
	}
	if month < 1 || month > 12 {
		return "", fmt.Errorf("Invalid month: %d", month)
	}

	// Map month number to abbreviated month name
	monthName := monthNames[month-1]

	// Construct the new filename
	newFilename := fmt.Sprintf("%s.csv", monthName)
	return newFilename, nil
}
