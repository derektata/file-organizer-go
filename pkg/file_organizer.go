package pkg

import (
	"file-organizer/config"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// FileOrganizer organizes files in a directory based on configuration and strategies.
type FileOrganizer struct {
	Path    string
	Config  *config.ConfigLoader
	Options OrganizerOptions
}

// OrganizerOptions contains options for the FileOrganizer.
type OrganizerOptions struct {
	PrependDate bool
	DryRun      bool
}

// NewFileOrganizer creates a new instance of the FileOrganizer struct.
//
// Parameters:
// - path: the path of the directory to organize files (string)
// - config: the configuration loader (pointer to ConfigLoader)
// - options: the options for the file organizer (OrganizerOptions)
//
// Returns:
// - a pointer to the newly created FileOrganizer (pointer to FileOrganizer)
func NewFileOrganizer(path string, config *config.ConfigLoader, options OrganizerOptions) *FileOrganizer {
	return &FileOrganizer{
		Path:    path,
		Config:  config,
		Options: options,
	}
}

// OrganizeFiles organizes files in the specified directory, optionally prepending the date to the filenames.
//
// Parameters:
// - prependDate: a boolean indicating whether to prepend the current date to the file names
// - tree: a pointer to a DirectoryTree object, used in dry-run mode to add the file to the tree
//
// Returns:
// - error: an error if the directory could not be read or the file could not be moved
func (o *FileOrganizer) OrganizeFiles(prependDate bool) error {
	files, err := os.ReadDir(o.Path)
	CheckErr(err, "failed to read directory %s: %v", o.Path, err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		srcPath := filepath.Join(o.Path, file.Name())
		destName := file.Name()

		if prependDate {
			destName = fmt.Sprintf("%s_%s", time.Now().Format("20060102"), destName)
		}

		err := o.MoveFile(srcPath, destName, prependDate)
		CheckErr(err, "failed to move file %s: %v", srcPath, err)
	}

	return nil
}

// MoveAndCreateDir creates the directory if it doesn't exist and moves the file to the new location.
//
// Parameters:
// - filePath: the path of the file to be moved (string)
// - categoryPath: the path of the directory where the file will be moved (string)
// - newFileName: the new name of the file (string)
//
// Returns:
// - error: an error if the directory could not be created or the file could not be moved (error)
func MoveAndCreateDir(filePath, categoryPath, newFileName string) error {
	if err := os.MkdirAll(categoryPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", categoryPath, err)
	}

	if err := os.Rename(filePath, filepath.Join(categoryPath, newFileName)); err != nil {
		return fmt.Errorf("failed to move file %s: %v", filePath, err)
	}
	return nil
}

// contains checks if a slice contains a specific element.
//
// Parameters:
// - slice: the slice to search in ([]T)
// - element: the element to search for (T)
//
// Returns:
// - bool: true if the element is found in the slice, false otherwise
func contains[T comparable](slice []T, element T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

// findCategory finds the category for a given file extension.
//
// Parameters:
// - extension: the file extension to search for (string)
//
// Returns:
// - category: the category found for the given extension (string)
// - found: a boolean indicating if a category was found (bool)
func (o *FileOrganizer) findCategory(extension string) (string, bool) {
	for category, extensions := range o.Config.FileExtensions {
		if contains(extensions, extension) {
			return category, true
		}
	}
	return "", false
}

// MoveFile moves a file to a categorized directory based on its extension and optionally prepends the date to the file name.
//
// Parameters:
// - filePath: the path of the file to be moved (string)
// - fileName: the name of the file to be moved (string)
// - prependDate: a boolean indicating whether to prepend the current date to the file name (bool)
//
// Returns:
// - error: an error if the file could not be moved or the directory could not be created (error)
func (o *FileOrganizer) MoveFile(filePath, fileName string, prependDate bool) error {
	fileExtension := strings.ToLower(filepath.Ext(fileName))

	category, found := o.findCategory(fileExtension)
	if !found {
		// fallback to mimetype if no category is found
		category = CategoryFromMimeType(filePath)

		switch category {
		case "":
			return nil
		case "others":
			return nil
		}

		category = strings.ToLower(category)
	}

	categoryPath := filepath.Join(o.Path, category)
	newFileName := fileName

	if prependDate {
		currentTime := time.Now().Format("2006-01-02")
		newFileName = fmt.Sprintf("%s_%s", currentTime, fileName)
	}

	return MoveAndCreateDir(filePath, categoryPath, newFileName)
}
