package commands

import (
	"context"
	"fmt"
	"clify/internal/client"
	"clify/internal/config"
	"clify/internal/models"
	"strings"
)

type SetupCommand struct {
	config *models.Config
}

func NewSetupCommand() *SetupCommand {
	return &SetupCommand{}
}

func (s *SetupCommand) Run() error {
	fmt.Println("clify Setup Wizard")
	fmt.Println("===================")
	fmt.Println()

	// Check if API key already exists
	currentConfig, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load current config: %w", err)
	}

	if currentConfig.APIKey != "" {
		fmt.Println("API key already configured")
		fmt.Print("Would you like to update it? (y/N): ")

		var response string
		fmt.Scanln(&response)

		if strings.ToLower(response) != "y" && strings.ToLower(response) != "yes" {
			fmt.Println("Setup cancelled.")
			return nil
		}
	}

	// Get API key from user
	apiKey, err := s.promptForAPIKey()
	if err != nil {
		return fmt.Errorf("failed to get API key: %w", err)
	}

	// Validate API key
	fmt.Print("Validating API key...")
	if err := s.validateAPIKey(apiKey); err != nil {
		fmt.Println(" FAILED")
		return fmt.Errorf("API key validation failed: %w", err)
	}
	fmt.Println(" SUCCESS")

	// Update config
	currentConfig.APIKey = apiKey
	if err := config.SaveConfig(currentConfig); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println()
	fmt.Println("Setup completed successfully!")
	fmt.Println("You can now use clify to query commands.")
	fmt.Println()
	fmt.Println("Examples:")
	fmt.Println("  clify \"find all .txt files\"")
	fmt.Println("  clify \"kill process on port 8080\"")
	fmt.Println("  clify \"compress folder to zip\"")
	fmt.Println()

	return nil
}

func (s *SetupCommand) promptForAPIKey() (string, error) {
	fmt.Println("You need an Anthropic API key to use clify.")
	fmt.Println("Get one at: https://console.anthropic.com/")
	fmt.Println()
	fmt.Print("Enter your Anthropic API key: ")

	var apiKey string
	_, err := fmt.Scanln(&apiKey)
	if err != nil {
		return "", fmt.Errorf("failed to read API key: %w", err)
	}

	apiKey = strings.TrimSpace(apiKey)
	if apiKey == "" {
		return "", fmt.Errorf("API key cannot be empty")
	}

	return apiKey, nil
}

func (s *SetupCommand) validateAPIKey(apiKey string) error {
	// Basic validation
	if err := config.ValidateAPIKey(apiKey); err != nil {
		return err
	}

	// Test connection
	client := client.NewClaudeClient(apiKey)
	ctx := context.Background()

	return client.TestConnection(ctx)
}

func (s *SetupCommand) IsSetupRequired() bool {
	config, err := config.LoadConfig()
	if err != nil {
		return true
	}

	return config.APIKey == ""
}

func (s *SetupCommand) ShowSetupPrompt() {
	fmt.Println("clify is not configured yet.")
	fmt.Println("Run 'clify setup' to get started.")
	fmt.Println()
}
