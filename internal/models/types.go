package models

import (
	"time"
)

// Response represents the AI response with commands and explanation
type Response struct {
	Explanation string    `json:"explanation"`
	Commands    []Command `json:"commands"`
}

// Command represents a single executable command
type Command struct {
	Text        string `json:"text"`
	Description string `json:"description"`
	SafetyLevel string `json:"safety_level"` // "safe", "warning", "dangerous"
}

// CacheEntry represents a cached query and response
type CacheEntry struct {
	Query     string    `json:"query"`
	Response  string    `json:"response"`
	Timestamp time.Time `json:"timestamp"`
}

// Config represents application configuration
type Config struct {
	APIKey    string `yaml:"api_key"`
	CacheFile string `yaml:"cache_file"`
	Model     string `yaml:"model"`
}

// SafetyLevel represents the safety classification of a command
type SafetyLevel string

const (
	SafetyLevelSafe      SafetyLevel = "safe"
	SafetyLevelWarning   SafetyLevel = "warning"
	SafetyLevelDangerous SafetyLevel = "dangerous"
)

// AppState represents the current state of the TUI application
type AppState struct {
	Mode            string
	Query           string
	Response        *Response
	SelectedCommand int
	ShowingHelp     bool
	ShowingTutorial bool
}

// TutorialStep represents a single step in the tutorial
type TutorialStep struct {
	Title       string
	Description string
	Example     string
	Action      string
}