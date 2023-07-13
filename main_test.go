package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// printFileTree prints the file tree starting from the given root directory.
//
// Parameters:
//   - root: the root directory to start printing the file tree from.
//   - fileIndent: the indentation string used for formatting the file tree.
//
// Return type: void.
func printFileTree(root string, fileIndent string) {
	files, err := os.ReadDir(root)
	if err != nil {
		fmt.Printf("Failed to read directory: %s\n", err)
		return
	}

	for _, file := range files {
		if file.IsDir() {
			fmt.Println(fileIndent + file.Name() + "/")
			printFileTree(filepath.Join(root, file.Name()), fileIndent+"  ")
		} else {
			fmt.Println(fileIndent + file.Name())
		}
	}
}

// createAndMoveFiles is a function that creates temporary files and moves them to a specific folder.
//
// The function takes two parameters:
// - t: a *testing.T object used for testing purposes.
// - prependDate: a boolean value that determines whether to prepend the date to the file name or not.
//
// The function returns a string representing the path of the temporary directory where the files are created.
func createAndMoveFiles(t *testing.T, prependDate bool) string {
	tmpDir, err := os.MkdirTemp("", "download_folder")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	org := NewDownloadOrganizer(tmpDir)

	testFiles := []struct {
		name     string
		category string
	}{
		{"test.mp3", "audio"},
		{"test.mp4", "video"},
		{"test.jpg", "image"},
		{"test.pdf", "document"},
		{"test.zip", "archive"},
		{"test.exe", "executable"},
		{"test.go", "programming"},
	}

	for _, tf := range testFiles {
		path := filepath.Join(tmpDir, tf.name)
		file, err := os.Create(path)
		if err != nil {
			t.Fatalf("Failed to create test file: %v", err)
		}
		file.Close()

		if err := org.MoveFile(path, tf.name, prependDate); err != nil {
			t.Fatalf("Failed to move test file: %v", err)
		}
	}

	return tmpDir
}

// TestMoveFileWithPrependDate is a test function that verifies the functionality of MoveFileWithPrependDate.
//
// This function creates temporary files and directories, and then moves the files to the directories with a prepended date in the filename.
// It validates if the date format is correct and if the filename contains the string "test.".
// The function prints the file tree after moving the files.
// The function does not return anything.
func TestMoveFileWithPrependDate(t *testing.T) {
	tmpDir := createAndMoveFiles(t, true)
	defer os.RemoveAll(tmpDir)

	files, err := os.ReadDir(tmpDir)
	if err != nil {
		t.Fatalf("Error reading directory: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			subFiles, err := os.ReadDir(filepath.Join(tmpDir, file.Name()))
			if err != nil {
				t.Fatalf("Error reading sub-directory: %v", err)
			}

			if len(subFiles) != 1 {
				t.Fatalf("Expected exactly one file in the %s directory", file.Name())
			}

			// Validate if the date format is correct and "test." is in filename
			filename := subFiles[0].Name()
			datePart := filename[0:10]
			_, err = time.Parse("2006-01-02", datePart) // Changed to =
			if err != nil || !strings.Contains(filename, "test.") {
				t.Errorf("Prepended date not found or invalid in filename")
			}
		}
	}

	fmt.Println("\nFile tree after moving files:")
	printFileTree(tmpDir, "")
}

// TestMoveFileWithoutPrependDate is a test function that tests the behavior of moving files without prepending the date.
//
// It creates a temporary directory and moves files into it. Then it checks if the files were moved to the correct category.
// The function prints the file tree after moving the files.
//
// Parameters:
//
//	t: A testing.T object that provides methods for testing and reporting test results.
//
// Return type: None.
func TestMoveFileWithoutPrependDate(t *testing.T) {
	tmpDir := createAndMoveFiles(t, false)
	defer os.RemoveAll(tmpDir)

	testFiles := []string{
		"test.mp3",
		"test.mp4",
		"test.jpg",
		"test.pdf",
		"test.zip",
		"test.exe",
		"test.go",
	}

	for _, tf := range testFiles {
		if _, err := os.Stat(filepath.Join(tmpDir, tf)); !os.IsNotExist(err) {
			t.Fatalf("File was not moved to the correct category: %v", err)
		}
	}

	fmt.Println("\nFile tree after moving files:")
	printFileTree(tmpDir, "")
}

// TestNewDownloadOrganizer is a unit test for the NewDownloadOrganizer function.
//
// It tests the functionality of the NewDownloadOrganizer function by creating a temporary directory and
// checking if the DownloadsPath of the NewDownloadOrganizer instance matches the temporary directory.
//
// Parameters:
// - t: The testing.T instance for running the unit test.
//
// Return Type(s):
// None
func TestNewDownloadOrganizer(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "download_folder")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	org := NewDownloadOrganizer(tmpDir)

	if org.DownloadsPath != tmpDir {
		t.Fatalf("Expected DownloadsPath to be %v but got %v", tmpDir, org.DownloadsPath)
	}

	fmt.Print("\n")
}
