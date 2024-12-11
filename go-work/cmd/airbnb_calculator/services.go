package main

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

func ValidateFilename(filename string) error {
	base := filepath.Base(filename)
	ext := filepath.Ext(base)
	name := base[:len(base)-len(ext)]

	r, err := regexp.Compile(`^airbnb_\d{2}_\d{4}-\d{2}_\d{4}\.csv$`)
	if err != nil {
		return err
	}

	if !r.MatchString(base) {
		return fmt.Errorf("filename must be in the format airbnb_MM_YYYY-MM_YYYY.csv")
	}

	parts := strings.Split(name, "-")

	firstPart := parts[0]
	secondPart := parts[1]

	if firstPart != "airbnb_"+secondPart {
		return fmt.Errorf("months and years in the filename do not match")
	}

	return nil
}

func ProcessRecords(records [][]string) ([][]string, error) {
	headers := records[0]
	montoIdx := findIndex(headers[:], "Monto")
	nochesIdx := findIndex(headers[:], "Noches")
	if montoIdx == -1 || nochesIdx == -1 {
		return nil, fmt.Errorf("Monto or Noches column not found in CSV file")
	}

	lenRecordsWithoutHeaders := len(records) - 1
	filteredRecords := make([][]string, lenRecordsWithoutHeaders)
	i := 0
	for _, record := range records[1:] {
		if record[montoIdx] != "" {
			filteredRecords[i] = record
			i++
		}
	}
	filteredRecords = filteredRecords[:i]

	// Sum the "Noches" values
	var totalNights = 0
	for _, record := range filteredRecords {
		nights, err := strconv.Atoi(record[nochesIdx])
		if err != nil {
			return nil, err
		}

		totalNights += nights
	}

	// Sum the "Monto" values
	var total = 0
	for _, record := range filteredRecords {
		var rawAmount = record[montoIdx]
		rawAmount = strings.Split(rawAmount, ".")[0]

		amount, err := strconv.Atoi(rawAmount)
		if err != nil {
			return nil, err
		}
		total += amount
	}

	// add a new row with the total amount and nights under the corresponding headers
	totalRow := make([]string, len(headers))
	totalRow[0] = "Total:"
	totalRow[montoIdx] = strconv.Itoa(total) + ".00"
	totalRow[nochesIdx] = strconv.Itoa(totalNights)

	// append headers at the top
	var emptyRow []string
	filteredRecords = append(filteredRecords[:], emptyRow)
	copy(filteredRecords[1:], filteredRecords)
	filteredRecords[0] = headers

	// append total rows
	filteredRecords = append(filteredRecords[:], totalRow)
	return filteredRecords[:], nil
}

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
