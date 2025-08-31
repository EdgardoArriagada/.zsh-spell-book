package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"example.com/workspace/lib/utils"
)

var (
	zsb    string
	zsbDir = os.Getenv("HOME") + "/.zsh-spell-book"
)

// Bundle represents the content builder for zsh spell book
type Bundle struct {
	content strings.Builder
}

// NewBundler creates a new Bundle instance
func NewBundler() *Bundle {
	return &Bundle{}
}

// WriteString writes a string to the bundle content
func (b *Bundle) WriteString(s string) {
	b.content.WriteString(s)
}

// Write writes bytes to the bundle content
func (b *Bundle) Write(data []byte) {
	b.content.Write(data)
	b.content.WriteString("\n")
}

// BundleEnvFile processes the .env file if it exists
func (b *Bundle) BundleEnvFile() {
	envFile := filepath.Join(zsbDir, ".env")
	if _, err := os.Stat(envFile); err == nil {
		b.bundleFileAbsolute(envFile)
	}
}

// BundleFile processes a single file (relative to zsbDir) and adds its content to the bundle
func (b *Bundle) BundleFile(filename string) {
	fullPath := filepath.Join(zsbDir, filename)
	b.bundleFileAbsolute(fullPath)
}

// bundleFileAbsolute processes a single file with absolute path and adds its content to the bundle
func (b *Bundle) bundleFileAbsolute(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return // Skip files that don't exist
	}
	b.Write(data)
}

// BundleDir processes all .zsh files in a directory (relative to zsbDir) recursively
func (b *Bundle) BundleDir(basePath string) {
	fullPath := filepath.Join(zsbDir, basePath)
	b.bundleDirAbsolute(fullPath)
}

// bundleDirAbsolute processes all .zsh files in a directory with absolute path recursively
func (b *Bundle) bundleDirAbsolute(basePath string) {
	// Walk through directory and find matching files
	err := filepath.WalkDir(basePath, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil // Skip directories that don't exist
		}

		if d.IsDir() {
			return nil
		}

		// Check if file matches .zsh extension
		if strings.HasSuffix(path, ".zsh") {
			b.bundleFileAbsolute(path)
		}

		return nil
	})

	if err != nil {
		// Directory doesn't exist, skip silently
		return
	}
}

// ApplyTransformations applies text transformations and returns the final result
func (b *Bundle) ApplyTransformations() string {
	input := b.content.String()
	lines := strings.Split(input, "\n")
	var result []string

	// Set ZSB_TEMP_DIR
	zsbTempDir := filepath.Join(zsbDir, "src/temp")

	for _, line := range lines {
		// Remove comments (equivalent to sd '( |^)#.*' '')
		commentRegex := regexp.MustCompile(`( |^)#.*`)
		line = commentRegex.ReplaceAllString(line, "")

		// Replace variables
		line = strings.ReplaceAll(line, "${zsb}", zsb)
		line = strings.ReplaceAll(line, "$ZSB_DIR", zsbDir)
		line = strings.ReplaceAll(line, "$ZSB_TEMP_DIR", zsbTempDir)

		// Remove leading spaces (equivalent to sd '^ *' '')
		leadingSpaceRegex := regexp.MustCompile(`^ *`)
		line = leadingSpaceRegex.ReplaceAllString(line, "")

		// Skip empty lines
		if strings.TrimSpace(line) != "" {
			// Remove bash line breaks (equivalent to sd '\\\n' '')
			line = strings.ReplaceAll(line, "\\\n", "")
			result = append(result, line)
		}
	}

	return strings.Join(result, "\n")
}

func main() {
	// Get zsb prefix from environment or default to 'zsb'
	zsb = os.Getenv("zsb")
	if zsb == "" {
		zsb = "zsb"
	}

	bundler := NewBundler()

	// Add dynamic prefix function
	bundler.WriteString(fmt.Sprintf("%s.sourceFiles() for f in $*; do source $f; done\n", zsb))

	bundler.BundleEnvFile()

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
