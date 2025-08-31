package main

import (
	"fmt"
	"os"
	"path/filepath"

	"example.com/workspace/lib/utils"
)

var (
	zsb    string
	zsbDir = os.Getenv("HOME") + "/.zsh-spell-book"
)

func main() {
	// Get zsb prefix from environment or default to 'zsb'
	zsb = os.Getenv("zsb")
	if zsb == "" {
		zsb = "zsb"
	}

	bundler := NewBundler()

	// Add dynamic prefix function
	bundler.WriteString(fmt.Sprintf("%s.sourceFiles() for f in $*; do source $f; done\n", zsb))

	bundler.LoadFile(".env")
	bundler.LoadFile("src/zsh.config.zsh")
	bundler.LoadFile("src/globalVariables.zsh")

	bundler.LoadDir("src/utils")
	bundler.LoadDir("src/configurations")
	bundler.LoadDir("src/spells")
	bundler.LoadDir("src/temp/spells")
	bundler.LoadDir("src/automatic-calls")

	result := bundler.Bundle()

	resultPath := filepath.Join(zsbDir, "result.zsh")
	utils.Must(os.WriteFile(resultPath, []byte(result), 0644))

	fmt.Println("ℨ𝔰𝔥 𝔖𝔭𝔢𝔩𝔩𝔟𝔬𝔬𝔨 bundled!!")
}
