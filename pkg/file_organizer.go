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
	Verbose     bool
}

// OrganizeFiles Organizes files in the directory based on the provided options.
//
// The prependDate parameter determines whether the current date should be prepended to the file names.
// The dryRun parameter determines whether to simulate the organization or perform the actual move operation.
// Returns an error if the operation fails.
func (o *FileOrganizer) OrganizeFiles() error {
	if o.Options.Verbose {
		fmt.Printf("Starting organization in '%s'...\n", o.Path)
		if o.Options.DryRun {
			dryLabel := chalk.Yellow.NewStyle().
				WithTextStyle(chalk.Bold).Style
			fmt.Printf("%s: No files will be moved.\n", dryLabel("Dry-run mode"))
		}
	}

	// Map to simulate the new structure if dry-run is enabled
	simulated := make(map[string][]string)

	files, err := os.ReadDir(o.Path)
	CheckErr(err, "failed to read directory %s: %v", o.Path, err)

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		src := filepath.Join(o.Path, file.Name())
		dest := o.prepareName(file.Name(), o.Options.PrependDate)

		// Determine the category (directory) this file will be placed in
		destPath, _ := o.destPath(src, dest)
		category := filepath.Base(filepath.Dir(destPath))

		if o.Options.DryRun {
			simLabel := chalk.Cyan.NewStyle().
				WithTextStyle(chalk.Bold).Style
			// Simulate the organization by adding to the simulated structure
			simulated[category] = append(simulated[category], dest)
			if o.Options.Verbose {
				fmt.Printf("%s: %s -> %s\n", simLabel("Simulated move"), src, destPath)
			}
		} else {
			// Actual move operation
			err := o.moveFile(src, destPath)
			CheckErr(err, "failed to move file %s: %v", src, err)
			if o.Options.Verbose {
				fmt.Printf("Moved: %s -> %s\n", src, destPath)
			}
		}
	}

	if o.Options.DryRun {
		// Display the simulated tree structure
		o.printSimulatedTree(simulated)
	}

	if o.Options.Verbose {
		fmt.Println("Organization completed.")
	}

	return nil
}

// prepareName Prepends the current date to a file name if specified.
//
// name: The original file name.
// prependDate: A boolean indicating whether to prepend the date.
// Returns a string representing the file name, either with the date prepended or the original name.
func (o *FileOrganizer) prepareName(name string, prependDate bool) string {
	if prependDate {
		if o.Options.Verbose {
			fmt.Printf("Prepending date to file: %s\n", name)
		}
		return fmt.Sprintf("%s_%s", time.Now().Format("20060102"), name)
	}
	return name
}

// destPath generates the destination path for a file based on its extension or mime type.
//
// src: The source path of the file.
// name: The name of the file.
// Returns a string representing the destination path and an error.
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
// src: The source path of the file to move.
// dest: The destination path of the file.
// Returns an error if the operation fails.
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
// ext: The file extension to find the category for.
// Returns the category name as a string and a boolean indicating whether the category was found.
func (o *FileOrganizer) findCategory(ext string) (string, bool) {
	for cat, exts := range o.Config.FileExtensions {
		if contains(exts, ext) {
			return cat, true
		}
	}
	return "", false
}

// printSimulatedTree prints a simulated tree structure of files organized within a given path.
//
// simulated: A map of categories to a slice of file names.
// This function does not return anything.
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

// contains Checks if an element exists in a slice.
//
// slice: The slice to search in.
// elem: The element to search for.
// Returns a boolean indicating whether the element was found.
func contains[T comparable](slice []T, elem T) bool {
	for _, e := range slice {
		if e == elem {
			return true
		}
	}
	return false
}
