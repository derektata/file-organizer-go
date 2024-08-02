package main

import (
	"fmt"
	"os"
	"path/filepath"
)

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
func MoveAndCreateDir(filePath, categoryPath, newFileName string) error {
	if err := os.MkdirAll(categoryPath, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", categoryPath, err)
	}

	if err := os.Rename(filePath, filepath.Join(categoryPath, newFileName)); err != nil {
		return fmt.Errorf("failed to move file %s: %v", filePath, err)
	}
	return nil
}
