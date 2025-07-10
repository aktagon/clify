package tui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Spinner struct {
	frames []string
	index  int
}

type spinnerTickMsg time.Time

func NewSpinner() *Spinner {
	return &Spinner{
		frames: []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		index:  0,
	}
}

func (s *Spinner) View() string {
	return s.frames[s.index]
}

func (s *Spinner) Tick() tea.Cmd {
	return tea.Tick(time.Millisecond*120, func(t time.Time) tea.Msg {
		return spinnerTickMsg(t)
	})
}

func (s *Spinner) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case spinnerTickMsg:
		s.index = (s.index + 1) % len(s.frames)
		return s.Tick()
	}
	return nil
}