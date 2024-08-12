package main

import (
	"file-organizer/config"
	"file-organizer/pkg"
	"log"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

var directoryTree = &pkg.DirectoryTree{
	Root: &pkg.DirectoryNode{
		Name:     ".",
		Children: make(map[string]*pkg.DirectoryNode),
	},
}

func main() {
	// Inline argument parsing using pflag
	directory := flag.StringP("directory", "d", "", "The path to the directory to organize")
	dryRun := flag.Bool("dry-run", false, "Print the actions without moving any files")
	prependDate := flag.Bool("prepend-date", false, "Prepend the current date to the filenames when moving files")
	configPath := flag.StringP("config", "c", filepath.Join(os.Getenv("HOME"), ".config/file-organizer/config.json"), "The path to the configuration file")
	flag.Parse()

	// Validate that the directory flag is provided
	if *directory == "" {
		log.Fatal("Please specify a directory to organize using the --directory or -d flag.")
	}

	// Load the configuration
	configLoader, err := config.NewConfigLoader(*configPath)
	pkg.CheckErr(err, "Error loading configuration")

	// Initialize the FileOrganizer
	organizer := pkg.NewFileOrganizer(*directory, configLoader, pkg.OrganizerOptions{
		PrependDate: *prependDate,
		DryRun:      *dryRun,
	})

	// Organize the files
	if *dryRun {
		log.Println("Running in dry-run mode. No files will be moved.")
	}

	err = organizer.OrganizeFiles(directoryTree)
	pkg.CheckErr(err, "Error organizing files")

	// Print the tree view of the specified directory
	log.Println("Tree view of the planned organization:")
	absolutePath, err := filepath.Abs(*directory)
	pkg.CheckErr(err, "Error getting absolute path of directory")
	directoryTree.PrintSubTree(absolutePath)
}
