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

// NewFileOrganizer creates a new FileOrganizer with configurations loaded from a file.
func NewFileOrganizer(path string) (*FileOrganizer, error) {
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

// MoveFile moves the file to a categorized directory based on its extension and optionally prepends the date to the file name.
func (o *FileOrganizer) MoveFile(filePath, fileName string, prependDate bool) error {
	fileExtension := strings.ToLower(filepath.Ext(fileName))

	for category, extensions := range o.FileExtensions {
		for _, extension := range extensions {
			if fileExtension == extension {
				categoryPath := filepath.Join(o.Path, category)
				newFileName := fileName

				if prependDate {
					currentTime := time.Now().Format("2006-01-02")
					newFileName = fmt.Sprintf("%s_%s", currentTime, fileName)
				}

				return MoveAndCreateDir(filePath, categoryPath, newFileName)
			}
		}
	}
	return nil
}

// OrganizeFiles organizes files in the specified directory, optionally prepending the date to the filenames.
func (o *FileOrganizer) OrganizeFiles(prependDate bool) error {
	files, err := os.ReadDir(o.Path)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", o.Path, err)
	}

	for _, file := range files {
		if !file.IsDir() {
			if err := o.MoveFile(filepath.Join(o.Path, file.Name()), file.Name(), prependDate); err != nil {
				return fmt.Errorf("failed to move file %s: %v", file.Name(), err)
			}
		}
	}
	return nil
}
