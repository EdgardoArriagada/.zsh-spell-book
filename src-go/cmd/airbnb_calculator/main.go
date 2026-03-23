package main

import (
	"fmt"

	"example.com/workspace/lib/args"
	u "example.com/workspace/lib/utils"
)

func main() {
	d := u.Must1(args.Parse())
	u.Expect(d.Len == 2, "Usage: airbnb_calculator <csv_filename> <avaluo_fiscal>")

	filename := d.Args[0]
	avaluoFiscal := u.Must1(ParseAvaluoFiscal(d.Args[1]))

	u.Must(ValidateFilename(filename))
	u.AssertFileExists(filename)

	daysInMonth, monthName := u.Must2(GetDateData(filename))

	records := u.Must1(getRecords(filename))

	records = u.Must1(ProcessRecords(records[:], avaluoFiscal, func(totalIncome, totalNights int) float64 {
		return CalculateIva(avaluoFiscal, totalIncome, totalNights, daysInMonth)
	},
	))

	outputFilename := GetOutputFilename(monthName)

	u.Must(writeToCsv(records[:], outputFilename))

	u.Must(openFilesInTmux(filename, outputFilename))

	fmt.Printf("Opening files '%s' and '%s' in tmux with neovim.\n", filename, outputFilename)
}
