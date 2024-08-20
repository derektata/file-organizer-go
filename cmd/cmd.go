package cmd

import (
	"file-organizer/config"
	"file-organizer/pkg"
	"fmt"
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

func Run() {
	var path string
	var configPath string
	var prependDate bool
	var dryRun bool
	var verbose bool
	var version bool

	flag.StringVarP(&configPath, "config", "c", filepath.Join(os.Getenv("HOME"), ".config/file-organizer/config.json"), "The path to the configuration file")
	flag.StringVarP(&path, "directory", "d", "", "Path to organize files")
	flag.BoolVarP(&prependDate, "prepend-date", "", false, "Prepend the current date to the file name")
	flag.BoolVarP(&dryRun, "dry-run", "", false, "Perform a dry run without moving files")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Show detailed output")
	flag.BoolVar(&version, "version", false, "Show the version number")

	flag.Parse()

	if version {
		fmt.Println(Version())
		return
	}

	if path == "" {
		fmt.Println("Error: You must specify a directory to organize using the -d flag.")
		return
	}

	configLoader, err := config.NewConfigLoader(configPath)
	pkg.CheckErr(err, "Failed to load configuration: %v", err)

	fileOrganizer := &pkg.FileOrganizer{
		Path:    path,
		Config:  configLoader,
		Options: pkg.OrganizerOptions{PrependDate: prependDate, DryRun: dryRun, Verbose: verbose},
	}

	err = fileOrganizer.OrganizeFiles()
	pkg.CheckErr(err, "Failed to organize files: %v", err)
}
