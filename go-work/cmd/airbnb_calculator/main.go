package main

import (
	"fmt"

	"example.com/workspace/lib/args"
	"example.com/workspace/lib/open"
	u "example.com/workspace/lib/utils"
)

func main() {
	d := u.Must(args.Parse())
	u.Expect(d.Len == 1, "Usage: airbnb_calculator <csv_filename>")

	filename := d.Args[0]

	u.Assert(ValidateFilename(filename))
	u.AssertFileExists(filename)

	records := u.Must(getRecords(filename))
	records = u.Must(ProcessRecords(records[:]))
	outputFilename := u.Must(GetOutputFilename(filename))

	u.Assert(writeToCsv(records[:], outputFilename))

	open.File(filename)
	open.File(outputFilename)

	fmt.Printf("Opening files '%s' and '%s'.\n", filename, outputFilename)
}
