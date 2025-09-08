package main

import (
	"fmt"
	"os"
	"path/filepath"

	"example.com/workspace/lib/utils"
)

var (
	zsb             string
	ZSB_DIR         = os.Getenv("HOME") + "/.zsh-spell-book"
	ZSB_TEMP_DIR    = filepath.Join(ZSB_DIR, "src/temp")
	RESULT_FILEPATH = filepath.Join(ZSB_DIR, "result.zsh")
)

func main() {
	// Get zsb prefix from environment or default to 'zsb'
	zsb = os.Getenv("zsb")
	if zsb == "" {
		zsb = "zsb"
	}

	bundler := NewBundler()

	bundler.LoadFile(".env")
	bundler.LoadFile("src/zsh.config.zsh")
	bundler.LoadFile("src/globalVariables.zsh")

	bundler.LoadDir("src/utils")
	bundler.LoadDir("src/configurations")
	bundler.LoadDir("src/spells")
	bundler.LoadDir("src/temp/spells")
	bundler.LoadDir("src/automatic-calls")

	result := bundler.Bundle()

	utils.Must(os.WriteFile(RESULT_FILEPATH, []byte(result), 0644))

	fmt.Println("ℨ𝔰𝔥 𝔖𝔭𝔢𝔩𝔩𝔟𝔬𝔬𝔨 bundled!!")
}
