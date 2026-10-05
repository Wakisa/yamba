package schema

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/wakisa/yamba/steps"
)

// Mode represents the current screen/state of the UI
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
	Cursor int
}

// projectOption is a selectable entry on the "choose project type" screen
type projectOption struct {
	Key    string
	Label  string
	WithUI bool
}

var projectOptions = []projectOption{
	{Key: "b", Label: "Backend only", WithUI: false},
	{Key: "u", Label: "Backend + UI", WithUI: true},
}

type StepDone struct{}

// --- Lip Gloss Styles ---
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF79C6")).
			Padding(1, 2).
			Border(lipgloss.DoubleBorder()).
			BorderForeground(lipgloss.Color("#BD93F9"))

	stepDoneStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#50FA7B")).
			PaddingLeft(2)

	stepPendingStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#6272A4")).
				PaddingLeft(2)

	highlightStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#F1FA8C"))

	instructionStyle = lipgloss.NewStyle().
				Italic(true).
				Foreground(lipgloss.Color("#8BE9FD")).
				PaddingTop(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFB86C")).
			Padding(1, 2).
			Margin(1, 2).
			Width(80)
)

func NewModel(steps []string, withUI bool) Model {
	return Model{Steps: steps, Done: 0, WithUI: withUI, Mode: ModeConfirm}
}

// Init satisfies tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update satisfies tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.Mode {
		case ModeConfirm:
			if msg.String() == "y" {
				m.Mode = ModeChoose
			}
			if msg.String() == "q" {
				return m, tea.Quit
			}

		case ModeChoose:
			switch key := msg.String(); key {
			case "up", "k":
				m.Cursor = (m.Cursor - 1 + len(projectOptions)) % len(projectOptions)
			case "down", "j":
				m.Cursor = (m.Cursor + 1) % len(projectOptions)
			case "enter", " ":
				return m.chooseOption(projectOptions[m.Cursor])
			case "q", "ctrl+c":
				return m, tea.Quit
			default:
				for _, opt := range projectOptions {
					if key == opt.Key {
						return m.chooseOption(opt)
					}
				}
			}

		case ModeRunSteps:
			if msg.String() == "q" {
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
		content := lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("✨ Welcome to Yamba!"),
			fmt.Sprintf("Yamba will generate files in: %s", highlightStyle.Render(cwd)),
			instructionStyle.Render("Press 'y' to continue or 'q' to cancel."),
		)
		return boxStyle.Render(content)

	case ModeChoose:
		lines := []string{titleStyle.Render("Choose your project type:")}
		for i, opt := range projectOptions {
			label := fmt.Sprintf("[%s] %s", opt.Key, opt.Label)
			if i == m.Cursor {
				lines = append(lines, highlightStyle.PaddingLeft(1).Render("> "+label))
			} else {
				lines = append(lines, stepPendingStyle.Render("  "+label))
			}
		}
		lines = append(lines, instructionStyle.Render("Use ↑/↓ and Enter, or press a letter. Press 'q' to quit."))
		return boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))

	case ModeRunSteps:
		var stepsUI string
		for i, step := range m.Steps {
			if i < m.Done {
				stepsUI += stepDoneStyle.Render("[✔] "+step) + "\n"
			} else {
				stepsUI += stepPendingStyle.Render("[ ] "+step) + "\n"
			}
		}
		content := lipgloss.JoinVertical(lipgloss.Left,
			titleStyle.Render("🚀 Yamba Project Generator"),
			stepsUI,
			instructionStyle.Render("Press 'q' to quit."),
		)
		return boxStyle.Render(content)
	}
	return ""
}

func (m Model) chooseOption(opt projectOption) (tea.Model, tea.Cmd) {
	m.WithUI = opt.WithUI
	m.Mode = ModeRunSteps
	return m, m.nextStep()
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
