package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	// Pre-compiled regular expressions for better performance
	commentRegex      = regexp.MustCompile(`( |^)#.*`)
	leadingSpaceRegex = regexp.MustCompile(`^ *`)
)

// Bundler represents the content builder for zsh spell book
type Bundler struct {
	content strings.Builder
}

// NewBundler creates a new Bundle instance
func NewBundler() *Bundler {
	b := &Bundler{}
	// Pre-allocate some capacity based on typical file sizes
	b.content.Grow(64 * 1024) // Start with 64KB capacity
	return b
}

// WriteString writes a string to the bundle content
func (b *Bundler) WriteString(s string) {
	b.content.WriteString(s)
}

// Write writes bytes to the bundle content
func (b *Bundler) Write(data []byte) {
	b.content.Write(data)
	b.content.WriteByte('\n') // More efficient than WriteString("\n")
}

// LoadFile processes a single file (relative to zsbDir) and adds its content to the bundle
func (b *Bundler) LoadFile(filename string) {
	fullPath := filepath.Join(zsbDir, filename)
	b.appendFileContents(fullPath)
}

// appendFileContents processes a single file with absolute path and adds its content to the bundle
func (b *Bundler) appendFileContents(filename string) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return // Skip files that don't exist
	}
	b.Write(data)
}

// LoadDir processes all .zsh files in a directory (relative to zsbDir) recursively
func (b *Bundler) LoadDir(basePath string) {
	fullPath := filepath.Join(zsbDir, basePath)
	b.appendDirContents(fullPath)
}

// appendDirContents processes all .zsh files in a directory with absolute path recursively
func (b *Bundler) appendDirContents(basePath string) {
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
			b.appendFileContents(path)
		}

		return nil
	})

	if err != nil {
		// Directory doesn't exist, skip silently
		return
	}
}

// Bundle applies text transformations and returns the final result
func (b *Bundler) Bundle() string {
	input := b.content.String()
	lines := strings.Split(input, "\n")

	// Use strings.Builder for result construction - more efficient than slice + Join
	var result strings.Builder
	// Pre-allocate capacity based on input size (estimate ~80% of original size after processing)
	result.Grow(len(input) * 4 / 5)

	// Create a replacer for variable substitutions - more efficient for multiple replacements
	replacer := strings.NewReplacer(
		"${zsb}", zsb,
		"$ZSB_DIR", zsbDir,
		"$ZSB_TEMP_DIR", zsbTempDir,
		"\\\n", "", // Also handle bash line breaks here
	)

	for i, line := range lines {
		// Remove comments using pre-compiled regex
		line = commentRegex.ReplaceAllString(line, "")

		// Apply all variable replacements in one pass
		line = replacer.Replace(line)

		// Remove leading spaces using pre-compiled regex
		line = leadingSpaceRegex.ReplaceAllString(line, "")

		// Skip empty lines
		if strings.TrimSpace(line) != "" {
			if i > 0 && result.Len() > 0 {
				result.WriteByte('\n')
			}
			result.WriteString(line)
		}
	}

	return result.String()
}
