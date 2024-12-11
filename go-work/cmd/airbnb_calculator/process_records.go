package main

import (
	"fmt"
	"strconv"
	"strings"
)

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
