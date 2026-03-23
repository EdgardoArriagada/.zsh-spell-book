package main

import (
	"regexp"
	"strconv"
)

func cleanAvaluoFiscal(input string) string {
	re := regexp.MustCompile(`[\$\. ]`)
	return re.ReplaceAllString(input, "")
}

func ParseAvaluoFiscal(rawAvaluoFiscal string) (int, error) {
	return strconv.Atoi(cleanAvaluoFiscal(rawAvaluoFiscal))
}
