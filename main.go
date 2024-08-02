package main

import (
	"os"
	"path/filepath"

	flag "github.com/spf13/pflag"
)

var configPath string = filepath.Join(os.Getenv("HOME"), ".config/file-organizer/config.json")

func main() {
	var path string
	var prependDate bool
	flag.StringVarP(&path, "path", "p", "", "Path to organize files")
	flag.BoolVarP(&prependDate, "prepend-date", "d", false, "Prepend the current date to the file name")
	flag.Parse()

	organizer, err := NewFileOrganizer(path)
	if err != nil {
		CheckErr(err, "Failed to initialize file organizer")
	}

	if err := organizer.OrganizeFiles(prependDate); err != nil {
		CheckErr(err, "Failed to organize files")
	}
}
