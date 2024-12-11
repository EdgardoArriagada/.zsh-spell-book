package main

import (
	"fmt"
	"strconv"
	"strings"
)

func ProcessRecords(records [][]string, avaluoFiscal int, calculateIvaCallback func(totalAmount int, totalNights int) float64) ([][]string, error) {
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
	var totalAmount = 0
	for _, record := range filteredRecords {
		var rawAmount = record[montoIdx]
		rawAmount = strings.Split(rawAmount, ".")[0]

		amount, err := strconv.Atoi(rawAmount)
		if err != nil {
			return nil, err
		}
		totalAmount += amount
	}

	iva := calculateIvaCallback(totalAmount, totalNights)

	// append headers at the top
	var emptyRow []string
	filteredRecords = append(filteredRecords[:], emptyRow)
	copy(filteredRecords[1:], filteredRecords)
	filteredRecords[0] = headers

	// append separator row
	separatorRow := make([]string, len(headers))
	for i := range separatorRow {
		separatorRow[i] = "-------------"
	}
	filteredRecords = append(filteredRecords[:], separatorRow)

	// Sum row
	sumRow := make([]string, len(headers))
	sumRow[0] = "Sum:"
	sumRow[montoIdx] = strconv.Itoa(totalAmount) + ".00"
	sumRow[nochesIdx] = strconv.Itoa(totalNights)
	filteredRecords = append(filteredRecords[:], sumRow)

	// Avaluo Fiscal row
	avaluoFiscalRow := []string{"Avaluo Fiscal:", strconv.Itoa(avaluoFiscal)}
	filteredRecords = append(filteredRecords[:], avaluoFiscalRow)

	// Iva row
	ivaStr := strconv.FormatFloat(iva, 'f', 2, 64)
	ivaRow := []string{"Iva:", ivaStr}
	filteredRecords = append(filteredRecords[:], ivaRow)

	return filteredRecords[:], nil
}
