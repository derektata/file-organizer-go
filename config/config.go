package config

import (
	"bufio"
	"encoding/json"
	"file-organizer/pkg"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func init() {
	// Register the default configuration
	RegisterDefaultConfig()
}

func RegisterDefaultConfig() {
	configPath := getConfigPath()
	config := getDefaultConfig(configPath)

	if configFileExists(configPath) {
		fmt.Println("Configuration file already exists at:", configPath)
		return
	}

	fmt.Println("Configuration file not found at:", configPath)
	if shouldCreateEmptyConfig() {
		if createConfigDirectory(configPath) && createConfigFile(configPath, config) {
			fmt.Println("Empty configuration created at:", configPath)
		}
	} else {
		fmt.Println("No configuration file created.")
	}
}

func getConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config/file-organizer/config.json")
}

func getDefaultConfig(configPath string) *pkg.FileOrganizer {
	return &pkg.FileOrganizer{
		Path: configPath,
		FileExtensions: map[string][]string{
			"3d-model":     {},
			"audio":        {},
			"video":        {},
			"image":        {},
			"document":     {},
			"archive":      {},
			"application":  {},
			"presentation": {},
			"spreadsheet":  {},
			"programming":  {},
		},
	}
}

func configFileExists(configPath string) bool {
	_, err := os.Stat(configPath)
	return !os.IsNotExist(err)
}

func shouldCreateEmptyConfig() bool {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Would you like to generate an empty version? (y/n): ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(strings.ToLower(input))

		if input == "y" {
			return true
		} else if input == "n" {
			return false
		} else {
			fmt.Println("Invalid input. Please enter 'y' for yes or 'n' for no.")
		}
	}
}

func createConfigDirectory(configPath string) bool {
	err := os.MkdirAll(filepath.Dir(configPath), 0755)
	if err != nil {
		fmt.Println("Error creating directories:", err)
		return false
	}
	return true
}

func createConfigFile(configPath string, config *pkg.FileOrganizer) bool {
	file, err := os.Create(configPath)
	if err != nil {
		fmt.Println("Error creating config file:", err)
		return false
	}
	defer file.Close()

	err = json.NewEncoder(file).Encode(config)
	if err != nil {
		fmt.Println("Error writing config to file:", err)
		return false
	}

	return true
}
