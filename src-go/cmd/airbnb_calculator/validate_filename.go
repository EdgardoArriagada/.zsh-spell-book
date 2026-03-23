package main

import (
	"fmt"
	"path/filepath"
	"regexp"
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
