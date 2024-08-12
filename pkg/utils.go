package pkg

import (
	"fmt"
	"os"
	"path/filepath"
)

// contains checks if a slice contains a specific element.
func contains[T comparable](slice []T, element T) bool {
	for _, e := range slice {
		if e == element {
			return true
		}
	}
	return false
}

// CheckErr checks for an error and prints a message with optional format and arguments before exiting the program if an error is found.
func CheckErr(err error, format string, a ...interface{}) {
	if err != nil {
		if format != "" {
			msg := fmt.Sprintf(format, a...)
			fmt.Println(msg)
		}
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}

// MoveAndCreateDir creates the directory if it doesn't exist and moves the file to the new location.
func MoveAndCreateDir(filePath, categoryPath, newFileName string, dryRun bool, directoryTree *DirectoryTree) error {
	destPath := filepath.Join(categoryPath, newFileName)
	if dryRun {
		// Add to the directory tree for dry-run visualization
		directoryTree.AddFileToTree(destPath)
		return nil
	}

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(categoryPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", categoryPath, err)
	}

	// Move the file
	if err := os.Rename(filePath, destPath); err != nil {
		return fmt.Errorf("failed to move file %s: %v", filePath, err)
	}
	return nil
}
