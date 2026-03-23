package main

import (
	"fmt"
)

func GetOutputFilename(monthName string) string {
	return fmt.Sprintf("%s.csv", monthName)
}
