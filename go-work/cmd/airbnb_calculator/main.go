package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"

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

	records = processRecords(records)
	fmt.Println("le records", records)

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
	montoIndex := -1
	for i, header := range arr {
		if header == name {
			montoIndex = i
		}
	}
	return montoIndex
}

func processRecords(records [][]string) [][]string {
	headers := records[0]
	montoIndex := findIndex(headers, "Monto")

	fmt.Println("le montoIndex", montoIndex)

	return records
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
