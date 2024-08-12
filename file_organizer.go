package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileOrganizer struct {
	FileExtensions map[string][]string
	Path           string
}

// NewFileOrganizer creates a new FileOrganizer with configurations loaded from a specified config file.
func NewFileOrganizer(path string, configPath string) (*FileOrganizer, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", configPath, err)
	}

	var fileExtensions map[string][]string
	if err := json.Unmarshal(data, &fileExtensions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s: %v", configPath, err)
	}

	return &FileOrganizer{
		FileExtensions: fileExtensions,
		Path:           path,
	}, nil
}

// OrganizeFiles organizes files in the specified directory, optionally prepending the date to the filenames.
func (o *FileOrganizer) OrganizeFiles(prependDate bool, dryRun bool, directoryTree *DirectoryTree) error {
	files, err := os.ReadDir(o.Path)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", o.Path, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err := o.MoveFile(filepath.Join(o.Path, file.Name()), file.Name(), prependDate, dryRun, directoryTree)
		if err != nil {
			return fmt.Errorf("failed to move file %s: %v", file.Name(), err)
		}
	}
	return nil
}

// MoveFile moves the file to a categorized directory based on its extension or MIME type and optionally prepends the date to the file name.
func (o *FileOrganizer) MoveFile(filePath, fileName string, prependDate, dryRun bool, directoryTree *DirectoryTree) error {
	fileExtension := strings.ToLower(filepath.Ext(fileName))

	category, found := o.findCategory(fileExtension)
	if !found {
		// Fallback to MIME type detection if no category found in config
		category = o.categoryFromMimeType(filePath)
		if category == "" {
			return nil // No appropriate category found, skip file
		}
	}

	categoryPath := filepath.Join(o.Path, category)
	newFileName := fileName
	if prependDate {
		newFileName = fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), fileName)
	}

	return MoveAndCreateDir(filePath, categoryPath, newFileName, dryRun, directoryTree)
}

// findCategory finds the category of a file based on its extension.
func (o *FileOrganizer) findCategory(extension string) (string, bool) {
	for category, extensions := range o.FileExtensions {
		if contains(extensions, extension) {
			return category, true
		}
	}
	return "", false
}
