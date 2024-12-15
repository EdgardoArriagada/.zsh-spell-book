package main

import (
	"encoding/csv"
	"fmt"
	"os"
)

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

	return records[:], nil
}

func findIndex(arr []string, item string) (int, error) {
	idx := -1
	for i, current := range arr {
		if current == item {
			idx = i
		}
	}

	if idx == -1 {
		return 0, fmt.Errorf("item %s not found in slice", item)
	}

	return idx, nil
}

func writeToCsv(records [][]string, outputFilename string) error {
	// Write the filtered records to the new CSV file
	newFile, err := os.Create(outputFilename)
	if err != nil {
		return err
	}
	defer newFile.Close()

	writer := csv.NewWriter(newFile)
	err = writer.WriteAll(records[:])
	if err != nil {
		return err
	}

	return nil
}
