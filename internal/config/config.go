package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Provider struct {
	APIKey string `yaml:"api_key"`
}

type ProviderConfig struct {
	Provider string `yaml:"provider"`
	Model    string `yaml:"model"`
}

type Config struct {
	Providers     map[string]Provider `yaml:"providers"`
	LargeProvider ProviderConfig      `yaml:"large_provider"`
	SmallProvider ProviderConfig      `yaml:"small_provider"`
}

func GetPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "micro-agent.yml"
	}
	return filepath.Join(home, ".config", "micro-agent.yml")
}

func Load() (*Config, error) {
	configPath := GetPath()
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, err
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, err
	}
	return &config, nil
}

func Save(config *Config) error {
	configPath := GetPath()

	// Ensure .config directory exists
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}
