package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"runtime"

	// for embedding default config
	_ "embed"
)

//go:embed config/default-config.json
var defaultConfig []byte

func checkConfig() {
	configPath, err := getConfigPath()
	if err != nil {
		fmt.Println("Error getting config path:", err)
		os.Exit(1)
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// fmt.Println("Configuration file not found. Please create a config.json file at", path)
		// Use the default config to create the config file, if possible create the Notifier folder too
		if err := os.MkdirAll(path.Dir(configPath), 0755); err != nil {
			fmt.Println("Error creating Notifier directory:", err)
			os.Exit(1)
		}
		if err := os.WriteFile(configPath, defaultConfig, 0644); err != nil {
			fmt.Println("Error creating config file:", err)
			os.Exit(1)
		}
	}
}

type Pattern struct {
	Name    string `json:"name"`
	Regex   string `json:"regex"`
	Beep    bool   `json:"beep"`
	Toast   bool   `json:"toast"`
	Message string `json:"message"`
}

type Config struct {
	Patterns []Pattern `json:"patterns"`
}

func importConfig() (*Config, error) {
	configPath, err := getConfigPath()
	if err != nil {
		return nil, fmt.Errorf("error getting config path: %w", err)
	}
	configData, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}
	var config Config
	if err := json.Unmarshal(configData, &config); err != nil {
		return nil, fmt.Errorf("error unmarshalling config: %w", err)
	}
	return &config, nil
}

func getConfigPath() (string, error) {
	// This function should return the path to the configuration file.
	switch runtime.GOOS {
	case "windows":
		// For Windows, return the path to the config file in the user's Documents/My Games/Path of Exile/Notifier directory.
		user, _ := user.Current()
		return path.Join(user.HomeDir, "Documents", "My Games", "Path of Exile", "Notifier", "notifier_config.json"), nil
	case "linux", "darwin":
		return "", errors.New("unsupported operating system")
	default:
		// return error if the OS is not supported
		return "", errors.New("unknown operating system")
	}
}
