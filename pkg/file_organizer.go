package pkg

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"file-organizer/config"
)

// FileOrganizer organizes files in a directory based on configuration and strategies.
type FileOrganizer struct {
	Path    string
	Config  *config.ConfigLoader
	Mover   *FileMover
	Options OrganizerOptions
}

// OrganizerOptions holds options for organizing files.
type OrganizerOptions struct {
	PrependDate bool
	DryRun      bool
}

// NewFileOrganizer creates a new instance of the FileOrganizer struct.
//
// It takes in the following parameters:
// - path: the path of the directory to organize
// - config: a pointer to the ConfigLoader object
// - options: the OrganizerOptions object
//
// It returns a pointer to the newly created FileOrganizer.
func NewFileOrganizer(path string, config *config.ConfigLoader, options OrganizerOptions) *FileOrganizer {
	return &FileOrganizer{
		Path:    path,
		Config:  config,
		Mover:   NewFileMover(),
		Options: options,
	}
}

// OrganizeFiles organizes the files in the specified directory.
//
// It reads the files in the directory specified by the `Path` field of the `FileOrganizer` struct.
// For each file, it checks if it is a directory. If it is, it skips it.
// Otherwise, it calls the `organizeFile` method to organize the file.
//
// Parameters:
// - directoryTree: a pointer to the `DirectoryTree` object.
//
// Returns:
// - error: an error if there was a problem reading the directory or organizing a file.
func (o *FileOrganizer) OrganizeFiles(directoryTree *DirectoryTree) error {
	files, err := os.ReadDir(o.Path)
	if err != nil {
		return fmt.Errorf("failed to read directory %s: %v", o.Path, err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err := o.organizeFile(file.Name(), directoryTree)
		if err != nil {
			return fmt.Errorf("failed to organize file %s: %v", file.Name(), err)
		}
	}
	return nil
}

// organizeFile organizes a file by determining its category and moving it to the appropriate directory.
//
// Parameters:
// - fileName: the name of the file to be organized.
// - directoryTree: a pointer to the DirectoryTree object.
//
// Returns:
// - error: an error if the file organization fails.
func (o *FileOrganizer) organizeFile(fileName string, directoryTree *DirectoryTree) error {
	filePath := filepath.Join(o.Path, fileName)
	category := o.Mover.determineCategory(filePath, o.Config.FileExtensions)

	// Fallback to MIME type if no category found based on extension
	if category == "" {
		category = CategoryFromMimeType(filePath)
		if category == "" {
			return nil // Skip file if no appropriate category found
		}
	}

	newFileName := fileName
	if o.Options.PrependDate {
		newFileName = fmt.Sprintf("%s-%s", time.Now().Format("2006-01-02"), fileName)
	}

	destinationPath := filepath.Join(o.Path, category, newFileName)
	err := o.Mover.Move(filePath, destinationPath, newFileName, o.Options.DryRun, directoryTree)
	if err != nil {
		return err
	}

	// Add the file to the directory tree
	directoryTree.AddFileToTree(destinationPath)

	return nil
}

// FileMover handles the logic of moving files.
type FileMover struct{}

// NewFileMover creates a new instance of the FileMover struct.
//
// Returns:
// *FileMover: A pointer to a new instance of the FileMover struct.
func NewFileMover() *FileMover {
	return &FileMover{}
}

// Move moves a file from the given filePath to the destination directory with the newFileName.
//
// Parameters:
// - filePath: the path of the file to be moved.
// - destinationDir: the directory where the file will be moved.
// - newFileName: the new name of the file.
// - dryRun: a boolean indicating if the move operation should be performed.
// - directoryTree: a pointer to the DirectoryTree object.
//
// Returns:
// - error: an error if the move operation fails.
func (m *FileMover) Move(filePath, destinationDir, newFileName string, dryRun bool, directoryTree *DirectoryTree) error {
	destPath := filepath.Join(destinationDir, newFileName)

	if dryRun {
		// In dry-run mode, do not print the dry run message, only add the file to the tree
		return nil
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(destinationDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", destinationDir, err)
	}

	// Move the file
	if err := os.Rename(filePath, destPath); err != nil {
		return fmt.Errorf("failed to move file %s: %v", filePath, err)
	}
	return nil
}

// determineCategory determines the category of a file based on its extension.
//
// Parameters:
// - filePath: the path of the file.
// - fileExtensions: a map of categories and their corresponding extensions.
//
// Returns:
// - string: the category of the file, or an empty string if no category is found.
func (m *FileMover) determineCategory(filePath string, fileExtensions map[string][]string) string {
	extension := strings.ToLower(filepath.Ext(filePath))

	for category, extensions := range fileExtensions {
		if contains(extensions, extension) {
			return category
		}
	}

	// No category found, return empty string to trigger fallback to MIME type.
	return ""
}

// contains checks if a slice contains a specific string.
//
// Parameters:
// - slice: the slice of strings to search in.
// - item: the string to search for.
//
// Returns:
// - bool: true if the slice contains the item, false otherwise.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
