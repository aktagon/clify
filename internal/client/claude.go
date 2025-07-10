package client

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"clify/internal/models"
	"clify/internal/safety"
	"runtime"

	"github.com/aktagon/llmkit/anthropic"
)

//go:embed system_prompt.txt
var systemPrompt string

//go:embed command_response_schema.json
var commandResponseSchema string

type ClaudeClient struct {
	apiKey     string
	classifier *safety.Classifier
}

func NewClaudeClient(apiKey string) *ClaudeClient {
	classifier := safety.NewClassifier()
	
	return &ClaudeClient{
		apiKey:     apiKey,
		classifier: classifier,
	}
}

func (c *ClaudeClient) getOSInfo() string {
	osName := runtime.GOOS
	switch osName {
	case "darwin":
		return "macOS"
	case "linux":
		return "Linux"
	case "windows":
		return "Windows"
	case "freebsd":
		return "FreeBSD"
	case "openbsd":
		return "OpenBSD"
	case "netbsd":
		return "NetBSD"
	default:
		return "Unix-like"
	}
}

func (c *ClaudeClient) QueryCommands(ctx context.Context, query string) (*models.Response, error) {
	osInfo := c.getOSInfo()
	
	prompt := fmt.Sprintf(systemPrompt, osInfo, runtime.GOARCH, query, osInfo)

	schema := commandResponseSchema

	response, err := anthropic.Prompt("You are a helpful command-line assistant.", prompt, schema, c.apiKey)
	if err != nil {
		return nil, fmt.Errorf("failed to query Claude: %w", err)
	}

	if response == "" {
		return &models.Response{
			Explanation: "No commands found for the given query",
			Commands:    []models.Command{},
		}, nil
	}

	// Parse the full API response to extract the actual content
	var apiResponse struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	}
	
	if err := json.Unmarshal([]byte(response), &apiResponse); err != nil {
		return nil, fmt.Errorf("failed to parse API response: %w", err)
	}
	
	if len(apiResponse.Content) == 0 || apiResponse.Content[0].Text == "" {
		return &models.Response{
			Explanation: "No commands found for the given query",
			Commands:    []models.Command{},
		}, nil
	}
	
	// Extract the actual JSON content from the response
	jsonContent := apiResponse.Content[0].Text

	var result models.Response
	if err := json.Unmarshal([]byte(jsonContent), &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Commands) == 0 {
		return &models.Response{
			Explanation: "No commands found for the given query",
			Commands:    []models.Command{},
		}, nil
	}

	// Classify safety level for each command
	for i := range result.Commands {
		safetyLevel := c.classifier.ClassifyCommand(result.Commands[i].Text)
		result.Commands[i].SafetyLevel = string(safetyLevel)
	}

	return &result, nil
}

func (c *ClaudeClient) TestConnection(ctx context.Context) error {
	_, err := anthropic.Prompt("You are a helpful assistant.", "Respond with exactly: 'Connection successful'", "", c.apiKey)
	return err
}