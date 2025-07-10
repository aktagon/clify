package config

import (
	"context"
	"fmt"
	"clify/internal/client"
	"clify/internal/models"
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	DefaultConfigFile = "~/.clify/config.yaml"
	DefaultModel      = "claude-3-sonnet-20240229"
)

// getConfigPath returns the full path to the config file
func getConfigPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(home, ".clify", "config.yaml"), nil
}

func LoadConfig() (*models.Config, error) {
	config := &models.Config{
		Model:     DefaultModel,
		CacheFile: DefaultCacheFile,
	}

	// Try to load from environment variable first
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}

	// Try to load from config file
	configPath, err := getConfigPath()
	if err != nil {
		return config, nil // Return default config if home dir not found
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return config, nil // Config file doesn't exist, use defaults
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		return config, fmt.Errorf("failed to read config file: %w", err)
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return config, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Environment variable takes precedence
	if apiKey := os.Getenv("ANTHROPIC_API_KEY"); apiKey != "" {
		config.APIKey = apiKey
	}

	return config, nil
}

func SaveConfig(config *models.Config) error {
	configPath, err := getConfigPath()
	if err != nil {
		return err
	}

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

func ValidateAPIKey(apiKey string) error {
	if apiKey == "" {
		return fmt.Errorf("API key is required")
	}

	if len(apiKey) < 10 {
		return fmt.Errorf("API key appears to be invalid (too short)")
	}

	// Test the API key by making a connection
	claudeClient := client.NewClaudeClient(apiKey)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := claudeClient.TestConnection(ctx); err != nil {
		return fmt.Errorf("API key validation failed: %w", err)
	}

	return nil
}
