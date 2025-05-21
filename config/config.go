package config

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"file-organizer/types"

	"gopkg.in/yaml.v3"
)

type ConfigLoader struct {
	Path           string
	FileExtensions types.FileExtensions
}

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

// $HOME/.config/file-organizer/config.yaml
func GetConfigPath() string {
	return filepath.Join(os.Getenv("HOME"),
		".config/file-organizer/config.yaml")
}

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

func configFileExists(configPath string) bool {
	_, err := os.Stat(configPath)
	return !os.IsNotExist(err)
}

func shouldCreateEmptyConfig() bool {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Would you like to generate an empty version? (y/n): ")
		input, _ := reader.ReadString('\n')
		switch strings.TrimSpace(strings.ToLower(input)) {
		case "y":
			return true
		case "n":
			return false
		default:
			fmt.Println("Invalid input. Enter 'y' or 'n'.")
		}
	}
}

func createConfigDirectory(configPath string) bool {
	return os.MkdirAll(filepath.Dir(configPath), os.ModePerm) == nil
}

// --- YAML-specific helpers -----------------------------------------------------

func createConfigFile(configPath string, extensions types.FileExtensions) bool {
	file, err := os.Create(configPath)
	if err != nil {
		fmt.Println("Error creating config file:", err)
		return false
	}
	defer file.Close()

	enc := yaml.NewEncoder(file)
	enc.SetIndent(2)
	if err := enc.Encode(extensions); err != nil {
		fmt.Println("Error writing to config file:", err)
		return false
	}
	return true
}

func loadYAML(configPath string) (types.FileExtensions, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config %s: %v", configPath, err)
	}

	var fileExtensions types.FileExtensions
	if err := yaml.Unmarshal(data, &fileExtensions); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config %s: %v", configPath, err)
	}
	return fileExtensions, nil
}

// --- Public loader ------------------------------------------------------------

func loadConfig() (types.FileExtensions, error) {
	cfg := GetConfigPath()
	if !configFileExists(cfg) {
		fmt.Println("Configuration not found at:", cfg)
		if shouldCreateEmptyConfig() {
			if createConfigDirectory(cfg) &&
				createConfigFile(cfg, getDefaultConfig()) {
				fmt.Println("Empty YAML configuration created at:", cfg)
			}
		} else {
			fmt.Println("No configuration file created.")
		}
	}
	return loadYAML(cfg)
}
