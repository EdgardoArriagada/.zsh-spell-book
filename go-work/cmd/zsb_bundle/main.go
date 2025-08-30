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

func main() {
	// Get current working directory (equivalent to ZSB_DIR)
	zsbDir := os.Getenv("HOME") + "/.zsh-spell-book"

	// Get zsb prefix from environment or default to 'zsb'
	zsb := os.Getenv("zsb")
	if zsb == "" {
		zsb = "zsb"
	}

	// Remove existing files
	removeFileIfExists("transpiled.zsh")
	removeFileIfExists("temp_transpilation.zsh")

	var content strings.Builder

	// Add dynamic prefix function
	content.WriteString(fmt.Sprintf("%s.sourceFiles() for f in $*; do source $f; done\n", zsb))

	// Process files in the specific order defined in bundle.zsh
	processEnvFile(zsbDir, &content)
	processFile(filepath.Join(zsbDir, "src/zsh.config.zsh"), &content)
	processFile(filepath.Join(zsbDir, "src/globalVariables.zsh"), &content)

	// Process utils files
	processGlobPattern(filepath.Join(zsbDir, "src/utils"), "**/*.zsh", &content)

	// Process configuration files
	processGlobPattern(filepath.Join(zsbDir, "src/configurations"), "**/*.zsh", &content)

	// Process spell files
	processGlobPattern(filepath.Join(zsbDir, "src/spells"), "**/*.zsh", &content)

	// Process temporary spell files
	processGlobPattern(filepath.Join(zsbDir, "src/temp/spells"), "**/*.zsh", &content)

	// Process automatic calls
	processGlobPattern(filepath.Join(zsbDir, "src/automatic-calls"), "**/*.zsh", &content)

	// Apply text transformations
	result := applyTransformations(content.String(), zsb, zsbDir)

	// Write result to file
	utils.Must(os.WriteFile("result.zsh", []byte(result), 0644))

	fmt.Println("ℨ𝔰𝔥 𝔖𝔭𝔢𝔩𝔩𝔟𝔬𝔬𝔨 bundled!!")
}

func removeFileIfExists(filename string) {
	if _, err := os.Stat(filename); err == nil {
		os.Remove(filename)
	}
}

func processEnvFile(zsbDir string, content *strings.Builder) {
	envFile := filepath.Join(zsbDir, ".env")
	if _, err := os.Stat(envFile); err == nil {
		processFile(envFile, content)
	}
}

func processFile(filename string, content *strings.Builder) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return // Skip files that don't exist
	}
	content.Write(data)
	content.WriteString("\n")
}

func processGlobPattern(basePath, pattern string, content *strings.Builder) {
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
			processFile(path, content)
		}

		return nil
	})

	if err != nil {
		// Directory doesn't exist, skip silently
		return
	}
}

func applyTransformations(input, zsb, zsbDir string) string {
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
