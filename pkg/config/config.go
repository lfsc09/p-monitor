package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Config holds the application configuration
type Config struct {
	UpdateInterval  int    `json:"update_interval"`
	TimeUnit        string `json:"time_unit"`        // "seconds" or "minutes"
	TemperatureUnit string `json:"temperature_unit"` // "celsius" or "fahrenheit"
}

// Default returns the default configuration
func Default() *Config {
	return &Config{
		UpdateInterval:  5,
		TimeUnit:        "seconds",
		TemperatureUnit: "celsius",
	}
}

// GetUpdateIntervalSeconds returns the update interval in seconds
func (c *Config) GetUpdateIntervalSeconds() int {
	if c.TimeUnit == "minutes" {
		return c.UpdateInterval * 60
	}
	return c.UpdateInterval
}

// Load loads configuration from file
func Load() (*Config, error) {
	configPath := getConfigPath()

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config if it doesn't exist
		cfg := Default()
		if err := cfg.Save(); err != nil {
			return nil, fmt.Errorf("failed to create default config: %v", err)
		}
		return cfg, nil
	}

	// Read existing config
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &cfg, nil
}

// Save saves configuration to file
func (c *Config) Save() error {
	configPath := getConfigPath()

	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %v", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

// getConfigPath returns the path to the configuration file
func getConfigPath() string {
	return filepath.Join(os.Getenv("HOME"), ".p-monitor", "config.json")
}
