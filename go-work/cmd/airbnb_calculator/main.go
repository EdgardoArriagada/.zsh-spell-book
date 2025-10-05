package main

import (
	"fmt"
	"os/exec"

	"example.com/workspace/lib/args"
	u "example.com/workspace/lib/utils"
)

// openFilesInTmux creates a new tmux tab with neovim in horizontal split
// filename at the top pane and outputFilename at the bottom pane
func openFilesInTmux(filename, outputFilename string) error {
	// Create a new tmux window with the first file
	cmd := exec.Command("tmux", "new-window", "nvim", filename)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to create tmux window: %v", err)
	}

	// Split the window horizontally and open the second file in the bottom pane
	cmd = exec.Command("tmux", "split-window", "-v", "nvim", outputFilename)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to split tmux window: %v", err)
	}

	return nil
}

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
