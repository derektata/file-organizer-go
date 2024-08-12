package main

import (
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

	// Initialize the FileOrganizer
	organizer, err := pkg.NewFileOrganizer(*directory, *configPath)
	pkg.CheckErr(err, "Error initializing FileOrganizer")

	// Organize the files
	if *dryRun {
		log.Println("Running in dry-run mode. No files will be moved.")
	}
	err = organizer.OrganizeFiles(*prependDate, *dryRun, directoryTree)
	pkg.CheckErr(err, "Error organizing files")

	// Print the tree view of the specified directory
	if *dryRun {
		log.Println("Tree view of the planned organization:")
		absolutePath, err := filepath.Abs(*directory)
		pkg.CheckErr(err, "Error getting absolute path of directory")
		directoryTree.PrintSubTree(absolutePath)
	}
}
