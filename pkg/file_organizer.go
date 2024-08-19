package pkg

import (
	"file-organizer/config"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/ttacon/chalk"
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

// OrganizeFiles Organizes files in the directory based on the provided options.
//
// The prependDate parameter determines whether the current date should be prepended to the file names.
// The dryRun parameter determines whether to simulate the organization or perform the actual move operation.
// Returns an error if the operation fails.
func (o *FileOrganizer) OrganizeFiles(prependDate bool, dryRun bool) error {
	// Map to simulate the new structure if dry-run is enabled
	simulated := make(map[string][]string)

	files, err := os.ReadDir(o.Path)
	CheckErr(err, "failed to read directory %s: %v", o.Path, err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		src := filepath.Join(o.Path, file.Name())
		dest := o.prepareName(file.Name(), prependDate)

		// Determine the category (directory) this file will be placed in
		destPath, _ := o.destPath(src, dest)
		category := filepath.Base(filepath.Dir(destPath))

		if dryRun {
			// Simulate the organization by adding to the simulated structure
			simulated[category] = append(simulated[category], dest)
		} else {
			// Actual move operation
			err := o.moveFile(src, destPath)
			CheckErr(err, "failed to move file %s: %v", src, err)
		}
	}

	if dryRun {
		// Display the simulated tree structure
		o.printSimulatedTree(simulated)
	}

	return nil
}

// prepareName Prepares a file name by optionally prepending the current date.
//
// The name parameter is the original file name.
// The prependDate parameter determines whether the current date should be prepended to the file name.
// Returns the prepared file name as a string.
func (o *FileOrganizer) prepareName(name string, prependDate bool) string {
	if prependDate {
		return fmt.Sprintf("%s_%s", time.Now().Format("20060102"), name)
	}
	return name
}

// destPath determines the destination path for a file based on its extension and category.
//
// Parameters:
//
//	src (string): The source path of the file.
//	name (string): The name of the file.
//
// Returns:
//
//	string: The destination path where the file will be placed.
//	error: An error if the destination path cannot be determined.
func (o *FileOrganizer) destPath(src, name string) (string, error) {
	ext := strings.ToLower(filepath.Ext(name))
	category, found := o.findCategory(ext)
	if !found {
		// fallback to mimetype if no category is found
		category = CategoryFromMimeType(src)
		if category == "" {
			category = "others"
		}
	}

	return filepath.Join(o.Path, category, name), nil // Keep files within the specified directory
}

// moveFile Moves a file from the source path to the destination path.
//
// src is the source path of the file to be moved.
// dest is the destination path where the file will be moved.
// Returns an error if the file cannot be moved.
func (o *FileOrganizer) moveFile(src, dest string) error {
	destDir := filepath.Dir(dest)
	if err := os.MkdirAll(destDir, os.ModePerm); err != nil {
		return fmt.Errorf("failed to create directory %s: %v", destDir, err)
	}

	if err := os.Rename(src, dest); err != nil {
		return fmt.Errorf("failed to move file %s: %v", src, err)
	}
	return nil
}

// findCategory Finds the category of a file extension.
//
// ext is the file extension to search for.
// Returns the category name and a boolean indicating whether the category was found.
func (o *FileOrganizer) findCategory(ext string) (string, bool) {
	for cat, exts := range o.Config.FileExtensions {
		if contains(exts, ext) {
			return cat, true
		}
	}
	return "", false
}

// printSimulatedTree Prints a simulated tree structure of the file organization.
//
// The simulated map parameter contains the file categories as keys and a slice of file names as values.
// No return value.
func (o *FileOrganizer) printSimulatedTree(simulated map[string][]string) {
	dryLabel := chalk.Yellow.NewStyle().
		WithTextStyle(chalk.Bold).Style
	pathLabel := chalk.Magenta.NewStyle().
		WithTextStyle(chalk.Bold).Style
	fmt.Printf("%s. Displaying where files will be organized within '%s':\n", dryLabel("Dry-run mode enabled"), pathLabel(o.Path))
	fmt.Println(filepath.Base(o.Path) + "/")

	// Get the keys (categories) in a sorted order for consistent output
	categories := make([]string, 0, len(simulated))
	for cat := range simulated {
		categories = append(categories, cat)
	}

	// Iterate over categories to display the tree structure
	for i, cat := range categories {
		// Determine if this is the last category for correct branching
		categoryPrefix := "├── "
		if i == len(categories)-1 {
			categoryPrefix = "└── "
		}

		catLabel := chalk.Blue.NewStyle().
			WithTextStyle(chalk.Bold).Style

		// Print the category branch
		fmt.Println(categoryPrefix + catLabel(cat+"/"))

		// Get the files in this category
		files := simulated[cat]
		for j, file := range files {
			filePrefix := "│   ├── "
			if i == len(categories)-1 {
				filePrefix = "    ├── "
			}
			if j == len(files)-1 {
				filePrefix = strings.Replace(filePrefix, "├──", "└──", 1)
			}
			fmt.Println(filePrefix + file)
		}
	}
}

// contains checks if a slice contains a specific element.
//
// The slice parameter is the list of elements to search through and the elem parameter is the element to search for.
// Returns true if the element is found in the slice, false otherwise.
func contains[T comparable](slice []T, elem T) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}
