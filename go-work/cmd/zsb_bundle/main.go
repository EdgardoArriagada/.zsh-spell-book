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

	bundler.BundleFile(".env")
	bundler.BundleFile("src/zsh.config.zsh")
	bundler.BundleFile("src/globalVariables.zsh")

	bundler.BundleDir("src/utils")
	bundler.BundleDir("src/configurations")
	bundler.BundleDir("src/spells")
	bundler.BundleDir("src/temp/spells")
	bundler.BundleDir("src/automatic-calls")

	result := bundler.ApplyTransformations()

	resultPath := filepath.Join(zsbDir, "result.zsh")
	utils.Must(os.WriteFile(resultPath, []byte(result), 0644))

	fmt.Println("ℨ𝔰𝔥 𝔖𝔭𝔢𝔩𝔩𝔟𝔬𝔬𝔨 bundled!!")
}
