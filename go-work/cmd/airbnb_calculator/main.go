package main

import (
	"fmt"
	"strconv"

	"example.com/workspace/lib/args"
	"example.com/workspace/lib/open"
	u "example.com/workspace/lib/utils"
)

func main() {
	d := u.Must(args.Parse())
	u.Expect(d.Len == 2, "Usage: airbnb_calculator <csv_filename> <avaluo_fiscal>")

	filename := d.Args[0]
	avaluoFiscal := u.Must(strconv.Atoi(d.Args[1]))

	u.Assert(ValidateFilename(filename))
	u.AssertFileExists(filename)

	month, year, monthName := u.Must3(GetDateData(filename))
	records := u.Must(getRecords(filename))
	records = u.Must(ProcessRecords(records[:], avaluoFiscal, month, year))
	outputFilename := GetOutputFilename(monthName)

	u.Assert(writeToCsv(records[:], outputFilename))

	open.File(filename)
	open.File(outputFilename)

	fmt.Printf("Opening files '%s' and '%s'.\n", filename, outputFilename)
}
