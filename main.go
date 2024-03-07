package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	flag "github.com/spf13/pflag"
)

var configPath string = os.Getenv("HOME") + "/.config/file-organizer/config.json"

type FileOrganizer struct {
	FileExtensions map[string][]string
	Path           string
}

// MoveFile moves the file to a categorized directory based on its extension and optionally prepends the date to the file name.
//
// Parameters:
//
// - filePath: the path of the file to be moved
//
// - fileName: the name of the file to be moved
//
// - prependDate: a boolean indicating whether to prepend the date to the file name
//
// Returns error.
func (o *FileOrganizer) MoveFile(filePath, fileName string, prependDate bool) error {
	fileExtension := strings.ToLower(filepath.Ext(fileName))

	for category, extensions := range o.FileExtensions {
		for _, extension := range extensions {
			if fileExtension == extension {
				categoryPath := filepath.Join(o.Path, category)
				err := os.MkdirAll(categoryPath, os.ModePerm)
				checkErr(err, "Failed to create directory: %s", categoryPath)

				newFileName := fileName
				if prependDate {
					// Generate the current time and format it to string
					currentTime := time.Now().Format("2006-01-02")
					// Add the current time to the beginning of the file name
					newFileName = fmt.Sprintf("%s_%s", currentTime, fileName)
				}

				err = os.Rename(filePath, filepath.Join(categoryPath, newFileName))
				checkErr(err, "Failed to move file: %s", filePath)
				return nil
			}
		}
	}
	return nil
}

// OrganizeFiles organizes files in the specified directory, optionally prepending the date to the filenames.
//
// Parameters:
//
// - prependDate: a boolean indicating whether to prepend the date to the filenames bool
func (o *FileOrganizer) OrganizeFiles(prependDate bool) {
	files, err := os.ReadDir(o.Path)
	checkErr(err, "Failed to read directory: %s", o.Path)

	for _, file := range files {
		if !file.IsDir() {
			err := o.MoveFile(filepath.Join(o.Path, file.Name()), file.Name(), prependDate)
			checkErr(err, "Failed to move file: %s", file.Name())
		}
	}
}

func main() {
	// Define a command-line flag for the path
	var path string
	var prependDate bool
	flag.StringVarP(&path, "path", "p", "", "Path to organize files")
	flag.BoolVarP(&prependDate, "prepend-date", "d", false, "Prepend the current date to the file name")
	flag.Parse()

	// Read the JSON file
	data, err := os.ReadFile(configPath)
	checkErr(err, "Failed to read config file: %s", configPath)

	// Create a FileOrganizer instance
	organizer := FileOrganizer{
		FileExtensions: make(map[string][]string),
		Path:           path,
	}

	// Unmarshal the JSON data into the FileExtensions field
	err = json.Unmarshal(data, &organizer.FileExtensions)
	checkErr(err, "Failed to unmarshal config file: %s", configPath)

	// Organize the files
	organizer.OrganizeFiles(prependDate)
}

// checkErr checks for an error and prints a message with optional format and arguments before exiting the program if an error is found.
//
// Parameters:
//
// - err: the error to check
//
// - format: the format string for the error message
//
// - a: the arguments for the format string
func checkErr(err error, format string, a ...interface{}) {
	if err != nil {
		if format != "" {
			msg := fmt.Sprintf(format, a...)
			fmt.Println(msg)
		}
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
