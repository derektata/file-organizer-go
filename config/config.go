package config

import (
	"bufio"
	"encoding/json"
	"file-organizer/types"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// ConfigLoader handles the loading and management of the configuration.
type ConfigLoader struct {
	Path           string
	FileExtensions types.FileExtensions
}

// NewConfigLoader loads the file extensions configuration from a specified file.
//
// It takes a `configPath` parameter, which is the path to the configuration file.
// It returns a pointer to a `ConfigLoader` struct and an error.
func NewConfigLoader(configPath string) (*ConfigLoader, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", configPath, err)
	}

	var fileExtensions types.FileExtensions
	if err := json.Unmarshal(data, &fileExtensions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s: %v", configPath, err)
	}

	return &ConfigLoader{
		Path:           configPath,
		FileExtensions: fileExtensions,
	}, nil
}

func init() {
	// Register the default configuration
	RegisterDefaultConfig()
}

// RegisterDefaultConfig checks if the default configuration file exists and creates
// it if it doesn't.
//
// It first gets the path of the configuration file.
// Then it checks if the file exists.
// If the file does not exist, it prompts the user to create an empty configuration.
// If the user chooses to create the configuration, it creates the configuration directory
// and the configuration file.
// Finally, it prints a message indicating the creation of the empty configuration.
func RegisterDefaultConfig() {
	configPath := getConfigPath()
	configExists := configFileExists(configPath)

	if !configExists {
		fmt.Println("Configuration file not found at:", configPath)
		if shouldCreateEmptyConfig() {
			config := getDefaultConfig()
			if createConfigDirectory(configPath) && createConfigFile(configPath, config) {
				fmt.Println("Empty configuration created at:", configPath)
			}
		} else {
			fmt.Println("No configuration file created.")
		}
	}
}

// getConfigPath returns the path to the configuration file.
//
// It retrieves the value of the HOME environment variable and appends the
// ".config/file-organizer/config.json" path to it.
//
// Returns:
// - string: the path to the configuration file.
func getConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config/file-organizer/config.json")
}

// getDefaultConfig returns a FileExtensions map with default values.
//
// Returns:
// - types.FileExtensions: a map of file extensions to their corresponding category lists.
func getDefaultConfig() types.FileExtensions {
	return types.FileExtensions{
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
	}
}

// configFileExists checks if a file exists at the specified path.
//
// It takes a `configPath` parameter of type `string`, which is the path to the file.
//
// It returns a boolean value indicating whether the file exists or not.
func configFileExists(configPath string) bool {
	_, err := os.Stat(configPath)
	return !os.IsNotExist(err)
}

// shouldCreateEmptyConfig prompts the user to input a choice whether to generate an empty version or not.
//
// It uses the bufio.NewReader to read user input from the standard input.
// The function continues to prompt the user until a valid input is provided.
//
// Returns:
// - bool: true if the user chooses to generate an empty version, false otherwise.
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

// createConfigDirectory creates the directory for the given config path.
//
// Parameters:
// - configPath: the path to the config file.
//
// Returns:
// - bool: true if the directory was created successfully, false otherwise.
func createConfigDirectory(configPath string) bool {
	err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
	return err == nil
}

// createConfigFile creates a new configuration file at the specified path.
//
// Parameters:
// - configPath: the path to the configuration file.
// - extensions: the FileExtensions map.
//
// Returns:
// - a boolean indicating whether the file was created successfully.
func createConfigFile(configPath string, extensions types.FileExtensions) bool {
	file, err := os.Create(configPath)
	if err != nil {
		fmt.Println("Error creating config file:", err)
		return false
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "    ")
	// Serialize the FileExtensions map directly
	if err := encoder.Encode(extensions); err != nil {
		fmt.Println("Error writing to config file:", err)
		return false
	}
	return true
}
