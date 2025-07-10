package commands

import (
	"fmt"
	"clify/internal/models"
	"strings"
)

type TutorialCommand struct {
	steps []models.TutorialStep
}

func NewTutorialCommand() *TutorialCommand {
	return &TutorialCommand{
		steps: []models.TutorialStep{
			{
				Title:       "Welcome to clify!",
				Description: "clify converts natural language queries into command-line commands using AI.",
				Example:     "clify \"find all .txt files\"",
				Action:      "Press Enter to continue",
			},
			{
				Title:       "Basic Usage",
				Description: "Simply describe what you want to do in plain English.",
				Example:     "clify \"kill process on port 8080\"",
				Action:      "Press Enter to continue",
			},
			{
				Title:       "Safety Features",
				Description: "Commands are color-coded for safety:\n🟢 Safe commands (read-only)\n🟡 Warning commands (make changes)\n🔴 Dangerous commands (destructive)",
				Example:     "clify \"delete all files\" # Would show 🔴",
				Action:      "Press Enter to continue",
			},
			{
				Title:       "Interactive Mode",
				Description: "Run 'clify' without arguments to enter interactive mode.",
				Example:     "Use arrow keys to navigate, Enter to select",
				Action:      "Press Enter to continue",
			},
			{
				Title:       "Caching",
				Description: "Responses are cached locally for faster repeated queries.",
				Example:     "~/.clify/cache.json stores your query history",
				Action:      "Press Enter to finish",
			},
		},
	}
}

func (t *TutorialCommand) Run() error {
	fmt.Println("clify Tutorial")
	fmt.Println("===============")
	fmt.Println()

	for i, step := range t.steps {
		t.displayStep(i+1, step)
		
		// Wait for user input
		var input string
		fmt.Scanln(&input)
		
		if strings.ToLower(input) == "q" || strings.ToLower(input) == "quit" {
			fmt.Println("Tutorial cancelled.")
			return nil
		}
		
		fmt.Println()
	}

	fmt.Println("Tutorial completed!")
	fmt.Println("You're ready to use clify. Try running:")
	fmt.Println("  clify \"your question here\"")
	fmt.Println()

	return nil
}

func (t *TutorialCommand) displayStep(stepNum int, step models.TutorialStep) {
	fmt.Printf("Step %d: %s\n", stepNum, step.Title)
	fmt.Println(strings.Repeat("=", len(step.Title)+10))
	fmt.Println()
	
	fmt.Println(step.Description)
	fmt.Println()
	
	if step.Example != "" {
		fmt.Printf("Example: %s\n", step.Example)
		fmt.Println()
	}
	
	fmt.Printf("=> %s (or 'q' to quit)\n", step.Action)
}