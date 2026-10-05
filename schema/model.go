package schema

import (
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
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

// totalDuration is how long the whole scaffold should take, so the user can
// watch each step happen instead of it all flashing by at once.
const totalDuration = 5 * time.Second

// logLines is how many recent activity entries are shown under the progress bar.
const logLines = 6

type Model struct {
	WithUI bool
	Mode   Mode
	Cursor int

	Tasks     []steps.Task
	StepNames []string
	Done      int // number of completed tasks
	Log       []string
	Err       error

	spinner  spinner.Model
	progress progress.Model
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

// taskDoneMsg is sent after each task finishes.
type taskDoneMsg struct{ err error }

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

	logStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#6272A4")).
			PaddingLeft(4)

	errorStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FF5555")).
			PaddingTop(1)

	successStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#50FA7B")).
			PaddingTop(1)

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#FFB86C")).
			Padding(1, 2).
			Margin(1, 2).
			Width(80)
)

func NewModel() Model {
	s := spinner.New(spinner.WithSpinner(spinner.Dot))
	s.Style = highlightStyle
	return Model{
		Mode:     ModeConfirm,
		spinner:  s,
		progress: progress.New(progress.WithDefaultGradient(), progress.WithWidth(60)),
	}
}

// Init satisfies tea.Model
func (m Model) Init() tea.Cmd {
	return nil
}

// Update satisfies tea.Model
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
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
			case "q":
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

	case spinner.TickMsg:
		if m.Mode == ModeRunSteps && !m.finished() {
			var cmd tea.Cmd
			m.spinner, cmd = m.spinner.Update(msg)
			return m, cmd
		}

	case taskDoneMsg:
		if msg.err != nil {
			m.Err = fmt.Errorf("%s: %w", m.Tasks[m.Done].Label, msg.err)
			return m, tea.Quit
		}
		m.Log = append(m.Log, m.Tasks[m.Done].Label)
		m.Done++
		if m.finished() {
			return m, tea.Quit
		}
		return m, m.runTask()
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
		lines := []string{titleStyle.Render("🚀 Yamba Project Generator")}

		// Checklist of steps: done, in progress (spinner), or pending.
		current := len(m.StepNames)
		if !m.finished() {
			current = slices.Index(m.StepNames, m.Tasks[m.Done].Step)
		}
		for i, name := range m.StepNames {
			switch {
			case i < current:
				lines = append(lines, stepDoneStyle.Render("[✔] "+name))
			case i == current && m.Err != nil:
				lines = append(lines, errorStyle.UnsetPaddingTop().PaddingLeft(2).Render("[✘] "+name))
			case i == current:
				lines = append(lines, "  "+m.spinner.View()+" "+highlightStyle.Render(name))
			default:
				lines = append(lines, stepPendingStyle.Render("[ ] "+name))
			}
		}

		// Overall progress bar.
		percent := float64(m.Done) / float64(len(m.Tasks))
		lines = append(lines, "",
			"  "+m.progress.ViewAs(percent)+fmt.Sprintf("  %d/%d", m.Done, len(m.Tasks)),
			"",
		)

		// Most recent activity.
		start := max(0, len(m.Log)-logLines)
		for _, entry := range m.Log[start:] {
			lines = append(lines, logStyle.Render("✔ "+entry))
		}

		switch {
		case m.Err != nil:
			lines = append(lines, errorStyle.Render("✘ "+m.Err.Error()))
		case m.finished():
			lines = append(lines, successStyle.Render("✨ Project ready!"))
		default:
			lines = append(lines, instructionStyle.Render("Press 'q' to quit."))
		}
		return boxStyle.Render(lipgloss.JoinVertical(lipgloss.Left, lines...))
	}
	return ""
}

func (m Model) finished() bool {
	return m.Done >= len(m.Tasks)
}

func (m Model) chooseOption(opt projectOption) (tea.Model, tea.Cmd) {
	cwd, _ := os.Getwd()
	m.WithUI = opt.WithUI
	m.Mode = ModeRunSteps
	m.Tasks = steps.BuildTasks(cwd, filepath.Base(cwd), m.WithUI)
	m.StepNames = steps.Names(m.Tasks)
	return m, tea.Batch(m.spinner.Tick, m.runTask())
}

// runTask runs the next task, padding it out so that all tasks together
// take roughly totalDuration.
func (m Model) runTask() tea.Cmd {
	task := m.Tasks[m.Done]
	delay := totalDuration / time.Duration(len(m.Tasks))
	return func() tea.Msg {
		start := time.Now()
		err := task.Run()
		if remaining := delay - time.Since(start); remaining > 0 {
			time.Sleep(remaining)
		}
		return taskDoneMsg{err: err}
	}
}
