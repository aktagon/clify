package tui

import (
	"clify/internal/client"
	"clify/internal/config"
	"clify/internal/models"
	"clify/internal/safety"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"golang.design/x/clipboard"
)

type Model struct {
	state      *models.AppState
	client     *client.ClaudeClient
	cache      *config.CacheManager
	classifier *safety.Classifier
	textInput  textinput.Model
	viewport   struct {
		width  int
		height int
	}
	lastError    string
	showingModal bool
	modalMessage string
	spinner      *Spinner
	loading      bool
	historyIndex int
}

type msgResponse struct {
	response *models.Response
	err      error
}

type msgError struct {
	message string
}

type msgCopied struct{}

type msgQuitWithMessage struct {
	message string
}

func NewModel(client *client.ClaudeClient, cache *config.CacheManager) *Model {
	ti := textinput.New()
	ti.Placeholder = "Enter your query..."
	ti.Focus()
	ti.CharLimit = 200
	ti.Width = 50
	ti.ShowSuggestions = true

	// Set up search history for autocomplete
	searchHistory := cache.GetSearchHistory()
	if len(searchHistory) > 0 {
		ti.SetSuggestions(searchHistory)
	}

	return &Model{
		state: &models.AppState{
			Mode:            "input",
			SelectedCommand: 0,
		},
		client:       client,
		cache:        cache,
		classifier:   safety.NewClassifier(),
		textInput:    ti,
		spinner:      NewSpinner(),
		loading:      false,
		historyIndex: -1,
	}
}

func (m *Model) SetInitialQuery(query string) {
	m.textInput.SetValue(query)
	m.state.Query = query
}

func (m *Model) Init() tea.Cmd {
	return m.textInput.Focus()
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.width = msg.Width
		m.viewport.height = msg.Height
		m.textInput.Width = msg.Width - 20
		return m, nil

	case tea.KeyMsg:
		return m.handleKeyPress(msg)

	case msgResponse:
		m.loading = false
		if msg.err != nil {
			m.lastError = msg.err.Error()
			return m, nil
		}
		m.state.Response = msg.response
		m.state.Mode = "selection"
		m.state.SelectedCommand = 0
		m.textInput.Blur()
		m.lastError = ""

		// Refresh autocomplete suggestions with updated search history
		searchHistory := m.cache.GetSearchHistory()
		if len(searchHistory) > 0 {
			m.textInput.SetSuggestions(searchHistory)
		}

		return m, nil

	case msgError:
		m.loading = false
		m.lastError = msg.message
		return m, nil

	case msgCopied:
		m.showingModal = true
		m.modalMessage = "Copied to clipboard!"
		return m, nil

	case msgQuitWithMessage:
		fmt.Println(msg.message)
		return m, tea.Quit

	case spinnerTickMsg:
		if m.loading {
			return m, m.spinner.Update(msg)
		}
		return m, nil
	}

	// Update textinput when in input mode
	if m.state.Mode == "input" {
		m.textInput, cmd = m.textInput.Update(msg)
		return m, cmd
	}

	return m, nil
}

func (m *Model) handleKeyPress(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if m.showingModal {
		return m.handleModalMode(msg)
	}

	switch m.state.Mode {
	case "input":
		return m.handleInputMode(msg)
	case "selection":
		return m.handleSelectionMode(msg)
	case "tutorial":
		return m.handleTutorialMode(msg)
	}
	return m, nil
}

func (m *Model) handleModalMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "enter", "c":
		m.showingModal = false
		m.modalMessage = ""
		return m, nil
	case "esc", "q":
		m.showingModal = false
		m.modalMessage = ""
		return m, tea.Quit
	}
	return m, nil
}

func (m *Model) handleInputMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		return m, tea.Quit

	case "enter":
		query := strings.TrimSpace(m.textInput.Value())
		if query == "" {
			return m, nil
		}
		m.loading = true
		m.historyIndex = -1
		return m, tea.Batch(m.queryCommand(query), m.spinner.Tick())

	case "up":
		history := m.cache.GetSearchHistory()
		if len(history) > 0 {
			if m.historyIndex == -1 {
				m.historyIndex = len(history) - 1
			} else if m.historyIndex > 0 {
				m.historyIndex--
			}
			m.textInput.SetValue(history[m.historyIndex])
		}
		return m, nil

	case "down":
		history := m.cache.GetSearchHistory()
		if len(history) > 0 && m.historyIndex != -1 {
			if m.historyIndex < len(history)-1 {
				m.historyIndex++
				m.textInput.SetValue(history[m.historyIndex])
			} else {
				m.historyIndex = -1
				m.textInput.SetValue("")
			}
		}
		return m, nil
	}

	// Update the text input for regular typing
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return m, cmd
}

func (m *Model) handleSelectionMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc":
		m.state.Mode = "input"
		m.state.Response = nil
		m.textInput.Focus()
		m.lastError = ""
		return m, nil

	case "up":
		if m.state.SelectedCommand > 0 {
			m.state.SelectedCommand--
		}

	case "down":
		if m.state.Response != nil && m.state.SelectedCommand < len(m.state.Response.Commands)-1 {
			m.state.SelectedCommand++
		}

	case "enter":
		if m.state.Response != nil && m.state.SelectedCommand < len(m.state.Response.Commands) {
			cmd := m.state.Response.Commands[m.state.SelectedCommand]
			// Copy command to clipboard or execute based on safety level
			return m, m.executeCommand(cmd)
		}

	case "n":
		// New query
		m.state.Mode = "input"
		m.state.Response = nil
		m.textInput.SetValue("")
		m.textInput.Focus()
		m.lastError = ""

		// Refresh autocomplete suggestions
		searchHistory := m.cache.GetSearchHistory()
		if len(searchHistory) > 0 {
			m.textInput.SetSuggestions(searchHistory)
		}

		return m, nil
	}

	return m, nil
}

func (m *Model) handleTutorialMode(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "ctrl+c", "esc", "q":
		m.state.Mode = "input"
		m.state.ShowingTutorial = false
		return m, nil
	}
	return m, nil
}

func (m *Model) queryCommand(query string) tea.Cmd {
	return func() tea.Msg {
		// Check cache first
		if cached, found := m.cache.Get(query); found {
			var response models.Response
			if err := json.Unmarshal([]byte(cached), &response); err == nil {
				return msgResponse{response: &response}
			}
		}

		// Query Claude
		ctx := context.Background()
		response, err := m.client.QueryCommands(ctx, query)
		if err != nil {
			return msgResponse{err: err}
		}

		// Cache the response
		if data, err := json.Marshal(response); err == nil {
			m.cache.Set(query, string(data))
		}

		return msgResponse{response: response}
	}
}

func (m *Model) executeCommand(cmd models.Command) tea.Cmd {
	return func() tea.Msg {
		err := clipboard.Init()
		if err != nil {
			return msgError{message: fmt.Sprintf("Failed to initialize clipboard: %v", err)}
		}
		clipboard.Write(clipboard.FmtText, []byte(cmd.Text))
		return msgCopied{}
	}
}

func (m *Model) View() string {
	if m.viewport.width == 0 {
		return "Loading..."
	}

	var baseView string
	switch m.state.Mode {
	case "input":
		baseView = m.renderInputView()
	case "selection":
		baseView = m.renderSelectionView()
	case "tutorial":
		baseView = m.renderTutorialView()
	}

	if m.showingModal {
		return m.renderModalOverlay(baseView)
	}

	return baseView
}

func (m *Model) renderInputView() string {
	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("35")).
		Render("clify - AI Command Helper")

	b.WriteString(title)
	b.WriteString("\n\n")

	// Error display
	if m.lastError != "" {
		errorStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)
		b.WriteString(errorStyle.Render(fmt.Sprintf("Error: %s", m.lastError)))
		b.WriteString("\n\n")
	}

	// Loading state with spinner
	if m.loading {
		loadingStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("33"))
		b.WriteString(loadingStyle.Render(fmt.Sprintf("%s Searching...", m.spinner.View())))
		b.WriteString("\n\n")
	}

	// Text input component
	b.WriteString(m.textInput.View())
	b.WriteString("\n\n")

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	b.WriteString(helpStyle.Render("Press Enter to search • Tab for autocomplete • Ctrl+C to quit"))

	return b.String()
}

func (m *Model) renderSelectionView() string {
	if m.state.Response == nil {
		return "No response available"
	}

	var b strings.Builder

	// Title
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("35")).
		Render("Query Results")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Explanation
	if m.state.Response.Explanation != "" {
		explanationStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("245"))
		b.WriteString(explanationStyle.Render(m.state.Response.Explanation))
		b.WriteString("\n\n")
	}

	// Commands
	for i, cmd := range m.state.Response.Commands {
		selected := i == m.state.SelectedCommand

		// Safety icon
		safetyLevel := models.SafetyLevel(cmd.SafetyLevel)
		icon := m.classifier.GetSafetyIcon(safetyLevel)

		// Command styling
		var cmdStyle lipgloss.Style
		if selected {
			cmdStyle = lipgloss.NewStyle().
				Background(lipgloss.Color("240")).
				Foreground(lipgloss.Color("15")).
				Bold(true)
		} else {
			cmdStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("15"))
		}

		// Render command
		b.WriteString(fmt.Sprintf("%s ", icon))
		b.WriteString(cmdStyle.Render(cmd.Text))
		b.WriteString("\n")

		// Description
		if cmd.Description != "" {
			descStyle := lipgloss.NewStyle().
				Foreground(lipgloss.Color("245")).
				MarginLeft(4)
			b.WriteString(descStyle.Render(cmd.Description))
			b.WriteString("\n")
		}

		if i < len(m.state.Response.Commands)-1 {
			b.WriteString("\n")
		}
	}

	b.WriteString("\n\n")

	// Help text
	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240"))
	b.WriteString(helpStyle.Render("↑/↓ Navigate • Enter to select • N for new query • Esc to go back"))

	return b.String()
}

func (m *Model) renderTutorialView() string {
	var b strings.Builder

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("35")).
		Render("Tutorial")
	b.WriteString(title)
	b.WriteString("\n\n")

	b.WriteString("Welcome to clify! This tutorial will help you get started.\n\n")
	b.WriteString("Press Q or Esc to exit tutorial")

	return b.String()
}

func (m *Model) renderModalOverlay(baseView string) string {
	if m.modalMessage == "" {
		return baseView
	}

	modalWidth := 50
	if m.viewport.width < 60 {
		modalWidth = m.viewport.width - 10
	}

	modalStyle := lipgloss.NewStyle().
		Width(modalWidth).
		Padding(2, 4).
		Background(lipgloss.Color("237")).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("35")).
		Align(lipgloss.Center)

	messageStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("15")).
		Bold(true).
		Align(lipgloss.Center)

	helpStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("245")).
		Align(lipgloss.Center)

	modalContent := messageStyle.Render(m.modalMessage) + "\n\n" +
		helpStyle.Render("Press Enter to continue • Press Esc to exit")

	modal := modalStyle.Render(modalContent)

	return lipgloss.Place(
		m.viewport.width,
		m.viewport.height,
		lipgloss.Center,
		lipgloss.Center,
		modal,
		lipgloss.WithWhitespaceChars(" "),
		lipgloss.WithWhitespaceForeground(lipgloss.Color("240")),
	)
}
