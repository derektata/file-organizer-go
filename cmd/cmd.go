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

	flag.StringVarP(&configPath, "config", "c", filepath.Join(os.Getenv("HOME"), ".config/file-organizer/config.json"), "The path to the configuration file")
	flag.StringVarP(&path, "directory", "d", "", "Path to organize files")
	flag.BoolVarP(&prependDate, "prepend-date", "", false, "Prepend the current date to the file name")
	flag.BoolVarP(&dryRun, "dry-run", "", false, "Perform a dry run without moving files")
	flag.Parse()

	if path == "" {
		fmt.Println("Error: You must specify a directory to organize using the -d flag.")
		return
	}

	configLoader, err := config.NewConfigLoader(configPath)
	pkg.CheckErr(err, "Failed to load configuration: %v", err)

	fileOrganizer := &pkg.FileOrganizer{
		Path:    path,
		Config:  configLoader,
		Options: pkg.OrganizerOptions{PrependDate: prependDate},
	}

	err = fileOrganizer.OrganizeFiles(prependDate, dryRun)
	pkg.CheckErr(err, "Failed to organize files: %v", err)
}
