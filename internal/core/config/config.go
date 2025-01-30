package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v2"
)

// Config represents the structure of our YAML configuration
type Config struct {
	Game    GameConfig    `yaml:"game"`
	Logging LoggingConfig `yaml:"logging"`
}

type LoggingConfig struct {
	Level        string `yaml:"level"`
	Format       string `yaml:"format"`
	ReportCaller bool   `yaml:"reportCaller"`
}

// GameConfig holds game-specific settings
type GameConfig struct {
	TickRate     int          `yaml:"tickRate"`
	MapSize      MapSize      `yaml:"mapSize"`
	ServerConfig ServerConfig `yaml:"serverConfig"`
}

// MapSize defines the dimensions of the game map
type MapSize struct {
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

// ServerConfig defines the server settings
type ServerConfig struct {
	Port int `yaml:"port"`
}

// LoadConfig reads and parses the YAML configuration file
func LoadConfig(environment string) (*Config, error) {
	configPath := filepath.Join("configs", fmt.Sprintf("%s.yaml", environment))
	absPath, err := filepath.Abs(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve config path for environment '%s': %v",
			environment, err)
	}

	yamlFile, err := os.ReadFile(absPath)
	if err != nil {
		if os.IsPermission(err) {
			return nil, fmt.Errorf("check file permissions and ownership for: %s", absPath)
		}

		return nil, fmt.Errorf("failed to read config file for environment '%s': at path '%s': %v",
			environment, absPath, err)
	}

	var config Config

	// Parse the YAML into our Config struct
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		if len(yamlFile) > 0 {
			preview := string(yamlFile)
			if len(preview) > 200 {
				preview = preview[:200] + "..."
			}

			fmt.Printf("YAML content preview:\n%s", preview)
		}

		return nil, fmt.Errorf("failed to parse YAML config for environment '%s': %v",
			environment, err)
	}

	// Validate the configuration
	if err := validateConfig(&config); err != nil {
		fmt.Printf("Invalid configuration values: %+v", config)
		return nil, fmt.Errorf("configuration validation failed for environment '%s': %v",
			environment, err)
	}

	fmt.Printf("\n[ -- Successfully loaded YAML config for << %s >> environment -- ]\n\n", strings.ToUpper(environment))
	return &config, nil
}

// validateConfig ensures all required values are set and valid
func validateConfig(config *Config) error {
	if config.Game.TickRate <= 0 {
		return fmt.Errorf("tick rate must be positive, got %d", config.Game.TickRate)
	}

	if config.Game.MapSize.Width <= 0 {
		return fmt.Errorf("map width must be positive, got %d", config.Game.MapSize.Width)
	}

	if config.Game.MapSize.Height <= 0 {
		return fmt.Errorf("map height must be positive, got %d", config.Game.MapSize.Height)
	}

	if config.Game.ServerConfig.Port <= 0 || config.Game.ServerConfig.Port > 65535 {
		return fmt.Errorf("server port must be between 1-65535, got %d", config.Game.MapSize.Height)
	}

	return nil
}
