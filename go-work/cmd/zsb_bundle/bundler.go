package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

var (
	commentRegex = regexp.MustCompile(`( |^)#.*`)
)

// Bundler represents the content builder for zsh spell book
type Bundler struct {
	content strings.Builder
}

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
	b.content.WriteByte('\n')
}

// LoadFile processes a single file (relative to zsbDir) and adds its content to the bundle
func (b *Bundler) LoadFile(filename string) {
	fullPath := filepath.Join(ZSB_DIR, filename)
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

// LoadDir processes all .zsh files in a directory
func (b *Bundler) LoadDir(basePath string) {
	fullPath := filepath.Join(ZSB_DIR, basePath)
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

	var result strings.Builder
	// Pre-allocate capacity based on input size (estimate ~80% of original size after processing)
	result.Grow(len(input) * 4 / 5)

	varReplacer := strings.NewReplacer(
		"${zsb}", zsb,
		"$ZSB_DIR", ZSB_DIR,
		"$ZSB_TEMP_DIR", ZSB_TEMP_DIR,
		"\\\n", "", // Also handle bash line breaks here
	)

	for _, line := range lines {
		// Remove comments using pre-compiled regex
		line = commentRegex.ReplaceAllString(line, "")

		line = varReplacer.Replace(line)
		line = strings.TrimSpace(line)

		// Skip empty lines
		if line == "" {
			continue
		}

		result.WriteByte('\n')
		result.WriteString(line)
	}

	return result.String()[1:] // Remove leading newline
}
