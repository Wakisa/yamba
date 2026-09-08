package schema

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/wakisa/yamba/steps"
)

type Mode int

const (
	ModeConfirm Mode = iota
	ModeChoose
	ModeRunSteps
)

type Model struct {
	Steps  []string
	Done   int
	WithUI bool
	Mode   Mode
}

type StepDone struct{}

func NewModel(steps []string, withUI bool) Model {
	return Model{Steps: steps, Done: 0, WithUI: withUI, Mode: ModeConfirm}
}

// Init satisfies tea.Model
func (m Model) Init() tea.Cmd {
	// Start in confirmation mode, no steps yet
	return nil
}

// Update satisfies tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.Mode {
		case ModeConfirm:
			if msg.String() == "y" { // to continue
				m.Mode = ModeChoose

			}
			if msg.String() == "q" { // to quit
				return m, tea.Quit
			}

		case ModeChoose:
			if msg.String() == "b" { // backend only
				m.WithUI = false
				m.Mode = ModeRunSteps
				return m, m.nextStep()
			}
			if msg.String() == "u" { // backend + UI
				m.WithUI = true
				m.Mode = ModeRunSteps
				return m, m.nextStep()
			}
			if msg.String() == "q" { // to quit
				return m, tea.Quit
			}
		}
	case StepDone:
		if m.Mode == ModeRunSteps {
			m.Done++
			if m.Done < len(m.Steps) {
				return m, m.nextStep()
			}
			return m, tea.Quit
		}
	}
	return m, nil
}

// View satisfies tea.Model
func (m Model) View() string {
	switch m.Mode {
	case ModeConfirm:
		cwd, _ := os.Getwd()
		return fmt.Sprintf("Welcome to Yamba!\nYamba will generate files in the current directory:\n%s\n\nPress 'y' to continue or 'q' to cancel.\n", cwd)

	case ModeChoose:
		return "Do you want to generate:\n[b] Backend only\n[u] Backend + UI\nPress 'q' to quit.\n"

	case ModeRunSteps:
		s := "Yamba Project Generator\n\n"
		for i, step := range m.Steps {
			if i < m.Done {
				s += fmt.Sprintf("[✔] %s\n", step)
			} else {
				s += fmt.Sprintf("[ ] %s\n", step)
			}
		}
		s += "\nPress 'q' to quit.\n"
		return s
	}
	return ""
}

func (m Model) nextStep() tea.Cmd {
	step := m.Steps[m.Done]
	return func() tea.Msg {
		cwd, _ := os.Getwd()
		projectName := filepath.Base(cwd)
		steps.RunStep(step, cwd, projectName, m.WithUI)
		return StepDone{}
	}
}
