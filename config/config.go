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
	fileExtensions, err := loadConfig()
	if err != nil {
		return nil, err
	}
	return &ConfigLoader{
		Path:           configPath,
		FileExtensions: fileExtensions,
	}, nil
}

// getConfigPath returns the path to the configuration file.
//
// It uses the `os.Getenv` function to retrieve the value of the `HOME` environment variable,
// and then uses the `filepath.Join` function to concatenate the path to the configuration
// file.
//
// Returns a string representing the path to the configuration file.
func getConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".config/file-organizer/config.json")
}

// getDefaultConfig returns a default configuration of file extensions.
//
// It initializes and returns a types.FileExtensions struct with empty slices for each file type.
//
// Returns a types.FileExtensions struct.
func getDefaultConfig() types.FileExtensions {
	return types.FileExtensions{
		"3d-model":     []string{},
		"audio":        []string{},
		"video":        []string{},
		"image":        []string{},
		"document":     []string{},
		"archive":      []string{},
		"application":  []string{},
		"presentation": []string{},
		"spreadsheet":  []string{},
		"programming":  []string{},
	}
}

// configFileExists checks if a configuration file exists at the specified path.
//
// It takes a `configPath` parameter, which is the path to the configuration file.
//
// The function uses the `os.Stat` function to retrieve information about the file at the specified path.
// If the file exists, the function returns `true`, indicating that the configuration file exists.
// If the file does not exist, the function returns `false`, indicating that the configuration file does not exist.
//
// Returns a boolean value indicating whether the configuration file exists or not.
func configFileExists(configPath string) bool {
	_, err := os.Stat(configPath)
	return !os.IsNotExist(err)
}

// shouldCreateEmptyConfig prompts the user to generate an empty configuration file.
//
// It takes no parameters.
//
// The function creates a new `bufio.Reader` from the standard input.
// It then enters a loop that prompts the user to enter 'y' for yes or 'n' for no.
// If the user enters 'y', the function returns `true`.
// If the user enters 'n', the function returns `false`.
// If the user enters any other input, the function prints an error message and prompts the user again.
//
// Returns a boolean value indicating whether the user chose to generate an empty configuration file.
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

// createConfigDirectory creates the directory specified by the given configPath.
//
// Parameters:
// - configPath: the path of the directory to be created (string)
//
// Returns:
// - a boolean indicating whether the directory was successfully created (bool)
func createConfigDirectory(configPath string) bool {
	err := os.MkdirAll(filepath.Dir(configPath), os.ModePerm)
	return err == nil
}

// createConfigFile creates a JSON file at the specified path and serializes the given
// FileExtensions map into it.
//
// Parameters:
// - configPath: the path where the JSON file will be created.
// - extensions: the FileExtensions map to be serialized.
//
// Returns:
// - a boolean indicating whether the file was successfully created and written to.
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

// loadJSON reads and parses a JSON file from the specified configPath.
//
// It takes a `configPath` parameter, which is the path to the configuration file.
// It returns a `types.FileExtensions` and an error.
func loadJSON(configPath string) (types.FileExtensions, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %v", configPath, err)
	}

	var fileExtensions types.FileExtensions
	if err := json.Unmarshal(data, &fileExtensions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config file %s: %v", configPath, err)
	}

	return fileExtensions, nil
}

// LoadConfig loads the configuration file from the specified path.
//
// It first retrieves the configuration file path using the getConfigPath function.
// Then it checks if the configuration file exists using the configFileExists function.
// If the configuration file does not exist, it prompts the user to create an empty configuration.
// If the user chooses to create an empty configuration, it creates the configuration directory
// using the createConfigDirectory function and creates the configuration file using the
// createConfigFile function.
// Finally, it loads the JSON configuration file using the loadJSON function.
//
// It returns the loaded configuration as a types.FileExtensions and an error.
func loadConfig() (types.FileExtensions, error) {
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
	return loadJSON(configPath)
}
