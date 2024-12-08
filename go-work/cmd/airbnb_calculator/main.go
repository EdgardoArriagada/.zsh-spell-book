package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"example.com/workspace/lib/args"
)

func main() {
	d, err := args.Parse()
	if err != nil {
		log.Fatal(err)
	}

	if d.Len < 1 {
		log.Fatalf("Usage: airbnb_calculator <csv_filename>")
	}

	filename := d.Args[0]

	// check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("File %s does not exist\n", filename)
	}

	records, err := getRecords(filename)
	if err != nil {
		log.Fatalf("Failed to extract records: %s", err)
	}

	records, err = ProcessRecords(records)
	if err != nil {
		log.Fatalf("Failed to process records: %s", err)
	}

	outputFilename := getOutputFilename(filename)

	err = writeToCsv(records, outputFilename)
	if err != nil {
		log.Fatalf("Failed to write to CSV: %s", err)
	}

	fmt.Printf("Filtered CSV file created: %s\n", outputFilename)
}

func getRecords(filename string) ([][]string, error) {
	// Open the CSV file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Read the CSV file
	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	return records, nil
}

func findIndex(arr []string, name string) int {
	idx := -1
	for i, header := range arr {
		if header == name {
			idx = i
		}
	}
	return idx
}

func ProcessRecords(records [][]string) ([][]string, error) {
	headers := records[0]
	montoIdx := findIndex(headers, "Monto")
	nochesIdx := findIndex(headers, "Noches")
	if montoIdx == -1 || nochesIdx == -1 {
		return nil, fmt.Errorf("Monto or Noches column not found in CSV file")
	}

	// Filter the records to remove any rows with an empty "Monto" value
	var filteredRecords [][]string = nil
	for _, record := range records[1:] {
		if record[montoIdx] != "" {
			filteredRecords = append(filteredRecords, record)
		}
	}

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
	totalRow[montoIdx] = strconv.Itoa(total)
	totalRow[nochesIdx] = strconv.Itoa(totalNights)

	// append headers at the top
	var emptyRow []string
	filteredRecords = append(filteredRecords, emptyRow)
	copy(filteredRecords[1:], filteredRecords)
	filteredRecords[0] = headers

	// append total rows
	filteredRecords = append(filteredRecords, totalRow)
	return filteredRecords, nil
}

func getOutputFilename(filename string) string {
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	newFilename := name + "_filtered" + ext
	return newFilename
}

func writeToCsv(records [][]string, outputFilename string) error {
	// Write the filtered records to the new CSV file
	newFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	err = writer.WriteAll(records)
	if err != nil {
		return err
	}

	return nil
}
