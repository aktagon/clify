package main

import (
	"clify/internal/client"
	"clify/internal/commands"
	"clify/internal/config"
	"clify/internal/tui"
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/bubbletea"
)

func main() {
	if len(os.Args) < 2 {
		// Interactive mode
		runInteractiveMode()
		return
	}

	command := os.Args[1]

	switch command {
	case "setup":
		setupCmd := commands.NewSetupCommand()
		if err := setupCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Setup failed: %v\n", err)
			os.Exit(1)
		}

	case "tutorial":
		tutorialCmd := commands.NewTutorialCommand()
		if err := tutorialCmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "Tutorial failed: %v\n", err)
			os.Exit(1)
		}

	case "help", "--help", "-h":
		showHelp()

	case "version", "--version", "-v":
		showVersion()

	default:
		// Direct query mode
		query := strings.Join(os.Args[1:], " ")
		runDirectQuery(query)
	}
}

func setupAndRunTUI(query string) {
	// Check if setup is required
	setupCmd := commands.NewSetupCommand()
	if setupCmd.IsSetupRequired() {
		setupCmd.ShowSetupPrompt()
		return
	}

	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to load config: %v\n", err)
		os.Exit(1)
	}

	// Initialize clients
	claudeClient := client.NewClaudeClient(cfg.APIKey)
	cache := config.NewCacheManager()

	// Create TUI model
	model := tui.NewModel(claudeClient, cache)

	// Set initial query if provided
	if query != "" {
		model.SetInitialQuery(query)
	}

	// Run the TUI
	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error running program: %v\n", err)
		os.Exit(1)
	}
}

func runInteractiveMode() {
	setupAndRunTUI("")
}

func runDirectQuery(query string) {
	setupAndRunTUI(query)
}

func showHelp() {
	fmt.Println("clify - AI-powered command-line helper")
	fmt.Println()
	fmt.Println("USAGE:")
	fmt.Println("  clify [command|query]")
	fmt.Println()
	fmt.Println("COMMANDS:")
	fmt.Println("  setup     Interactive setup wizard")
	fmt.Println("  tutorial  Interactive tutorial")
	fmt.Println("  help      Show this help message")
	fmt.Println("  version   Show version information")
	fmt.Println()
	fmt.Println("EXAMPLES:")
	fmt.Println("  clify                           # Interactive mode")
	fmt.Println("  clify \"find all .txt files\"     # Direct query")
	fmt.Println("  clify \"kill process on port 8080\"")
	fmt.Println("  clify \"compress folder to zip\"")
	fmt.Println()
	fmt.Println("SAFETY:")
	fmt.Println("  Commands are color-coded for safety:")
	fmt.Println("  🟢 Safe commands (read-only)")
	fmt.Println("  🟡 Warning commands (make changes)")
	fmt.Println("  🔴 Dangerous commands (destructive)")
	fmt.Println()
	fmt.Println("CONFIGURATION:")
	fmt.Println("  Set ANTHROPIC_API_KEY environment variable")
	fmt.Println("  Or run 'clify setup' for interactive configuration")
	fmt.Println()
}

func showVersion() {
	fmt.Println("clify version 1.0.0")
	fmt.Println("https://github.com/aktagon/clify")
}
