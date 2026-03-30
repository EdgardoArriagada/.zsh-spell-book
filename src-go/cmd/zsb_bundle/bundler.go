package main

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Bundler represents the content builder for zsh spell book
type Bundler struct {
	content strings.Builder
}

func NewBundler() *Bundler {
	return &Bundler{}
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

// Bundle applies variable replacements and returns the final result
func (b *Bundler) Bundle() string {
	replacer := strings.NewReplacer(
		"${zsb}", zsb,
		"$ZSB_DIR", ZSB_DIR,
		"$ZSB_TEMP_DIR", ZSB_TEMP_DIR,
	)
	return replacer.Replace(b.content.String())
}
