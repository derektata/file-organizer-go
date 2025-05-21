package cmd

import (
	"file-organizer/config"
	"file-organizer/pkg"
	"fmt"

	flag "github.com/spf13/pflag"
)

func Run() {
	var (
		path       string
		configPath string
		prepend    bool
		dryRun     bool
		verbose    bool
		version    bool
	)

	flag.StringVarP(&configPath, "config", "c", config.GetConfigPath(), "Path to the YAML configuration file")
	flag.StringVarP(&path, "directory", "d", "", "Path to organize files")
	flag.BoolVarP(&prepend, "prepend-date", "", false, "Prepend the current date to the file name")
	flag.BoolVarP(&dryRun, "dry-run", "", false, "Perform a dry run without moving files")
	flag.BoolVarP(&verbose, "verbose", "v", false, "Show detailed output")
	flag.BoolVar(&version, "version", false, "Show the version number")

	flag.Parse()

	if version {
		fmt.Println(Version())
		return
	}
	if path == "" {
		fmt.Println("Error: you must specify a directory with -d")
		return
	}

	cfg, err := config.NewConfigLoader(configPath)
	pkg.CheckErr(err, "Failed to load configuration: %v", err)

	org := &pkg.FileOrganizer{
		Path:    path,
		Config:  cfg,
		Options: pkg.OrganizerOptions{PrependDate: prepend, DryRun: dryRun, Verbose: verbose},
	}
	err = org.OrganizeFiles()
	pkg.CheckErr(err, "Failed to organize files: %v", err)
}
