package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Bundler represents the content builder for zsh spell book
type Bundler struct {
	content strings.Builder
}

// NewBundler creates a new Bundle instance
func NewBundler() *Bundler {
	return &Bundler{}
}

// WriteString writes a string to the bundle content
func (b *Bundler) WriteString(s string) {
	b.content.WriteString(s)
}

// Write writes bytes to the bundle content
func (b *Bundler) Write(data []byte) {
	b.content.Write(data)
	b.content.WriteString("\n")
}

// BundleEnvFile processes the .env file if it exists
func (b *Bundler) BundleEnvFile() {
	envFile := filepath.Join(zsbDir, ".env")
	if _, err := os.Stat(envFile); err == nil {
		b.bundleFileAbsolute(envFile)
	}
}

// BundleFile processes a single file (relative to zsbDir) and adds its content to the bundle
func (b *Bundler) BundleFile(filename string) {
	fullPath := filepath.Join(zsbDir, filename)
	b.bundleFileAbsolute(fullPath)
}

// bundleFileAbsolute processes a single file with absolute path and adds its content to the bundle
func (b *Bundler) bundleFileAbsolute(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return // Skip files that don't exist
	}
	b.Write(data)
}

// BundleDir processes all .zsh files in a directory (relative to zsbDir) recursively
func (b *Bundler) BundleDir(basePath string) {
	fullPath := filepath.Join(zsbDir, basePath)
	b.bundleDirAbsolute(fullPath)
}

// bundleDirAbsolute processes all .zsh files in a directory with absolute path recursively
func (b *Bundler) bundleDirAbsolute(basePath string) {
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
func (b *Bundler) ApplyTransformations() string {
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
