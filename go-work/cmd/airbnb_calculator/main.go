package main

import (
	"encoding/csv"
	"example.com/workspace/lib/args"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func main() {
	d, err := args.Parse()
	if err != nil {
		log.Fatal(err)
	}

	if d.Len < 1 {
		log.Fatalf("Usage: airbnb_calculator <csv_filename>")
	}

	filename := d.Args[1]

	// check if the file exists
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Fatalf("File %s does not exist\n", filename)
	}

	filteredRecords, err := extractRecords(filename)
	if err != nil {
		log.Fatalf("Failed to extract records: %s", err)
	}

	outputFilename := getOutputFilename(filename)

	err = writeToCsv(filteredRecords, outputFilename)
	if err != nil {
		log.Fatalf("Failed to write to CSV: %s", err)
	}

	fmt.Printf("Filtered CSV file created: %s\n", outputFilename)
}

func extractRecords(filename string) ([][]string, error) {
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

	// Filter records that have data in the 2nd column
	var filteredRecords [][]string
	for _, record := range records {
		if len(record) > 1 && record[1] != "" {
			filteredRecords = append(filteredRecords, record)
		}
	}
	return filteredRecords, nil
}

func getOutputFilename(filename string) string {
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]
	newFilename := name + "_filtered" + ext
	return newFilename
}

func writeToCsv(filteredRecords [][]string, outputFilename string) error {
	// Write the filtered records to the new CSV file
	newFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	err = writer.WriteAll(filteredRecords)
	if err != nil {
		return err
	}

	return nil
}
