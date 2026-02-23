package ui

import (
	"strings"
	"tml-sync/client/internal/i18n"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// Styles
var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FAFAFA")).
			Background(lipgloss.Color("#7D56F4")).
			Padding(0, 1).
			MarginBottom(1)

	statusStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#04B575"))

	infoStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#BBBBBB"))

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FF0000"))

	boxStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7D56F4")).
			Padding(1).
			Width(60)

	selectedStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#7D56F4")).
			Bold(true)
)

type State int

const (
	StateIdle State = iota
	StateConfirmScan
	StateConnecting
	StateSyncing
	StateDone
)

type Model struct {
	State           State
	Status          string
	Logs            []string
	Progress        progress.Model
	Spinner         spinner.Model
	ConfirmSelected bool // true for YES, false for NO
	
	// External command trigger
	OnConfirm func(scan bool) tea.Cmd

	Host            string
	Port            int
	
	Quitting        bool
	Width           int
	Height          int
}

func NewModel(host string, port int, onConfirm func(bool) tea.Cmd) Model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	p := progress.New(progress.WithDefaultGradient())

	return Model{
		State:     StateIdle,
		Status:    i18n.T("press_to_start"),
		Spinner:   s,
		Progress:  p,
		Host:      host,
		Port:      port,
		OnConfirm: onConfirm,
		ConfirmSelected: true,
	}
}

type StatusMsg string
type LogMsg string
type ProgressMsg float64
type StateMsg State

func (m Model) Init() tea.Cmd {
	return m.Spinner.Tick
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, tea.Quit
		case "left", "right", "tab":
			if m.State == StateConfirmScan {
				m.ConfirmSelected = !m.ConfirmSelected
			}
		case "enter":
			if m.State == StateIdle {
				m.State = StateConfirmScan
				m.Status = i18n.T("scan_prompt")
				return m, nil
			}
			if m.State == StateConfirmScan {
				m.State = StateConnecting
				if m.OnConfirm != nil {
					return m, m.OnConfirm(m.ConfirmSelected)
				}
			}
		}

	case tea.WindowSizeMsg:
		m.Width = msg.Width
		m.Height = msg.Height
		m.Progress.Width = msg.Width - 20
		if m.Progress.Width > 40 {
			m.Progress.Width = 40
		}

	case StateMsg:
		m.State = State(msg)
		return m, nil

	case StatusMsg:
		m.Status = string(msg)
		return m, nil

	case LogMsg:
		m.Logs = append(m.Logs, string(msg))
		if len(m.Logs) > 8 {
			m.Logs = m.Logs[1:]
		}
		return m, nil

	case ProgressMsg:
		cmd = m.Progress.SetPercent(float64(msg))
		return m, cmd

	case spinner.TickMsg:
		m.Spinner, cmd = m.Spinner.Update(msg)
		return m, cmd

	case progress.FrameMsg:
		newModel, cmd := m.Progress.Update(msg)
		if newModel, ok := newModel.(progress.Model); ok {
			m.Progress = newModel
		}
		return m, cmd
	}

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.Quitting {
		return ""
	}

	var s strings.Builder

	s.WriteString(headerStyle.Render(i18n.T("welcome")))
	s.WriteString("\n\n")

	var content strings.Builder
	content.WriteString(statusStyle.Render(m.Status) + "\n\n")

	if m.State == StateConfirmScan {
		yes := i18n.T("confirm_yes")
		no := i18n.T("confirm_no")
		if m.ConfirmSelected {
			content.WriteString(selectedStyle.Render("[ " + yes + " ]") + "   [ " + no + " ]")
		} else {
			content.WriteString("[ " + yes + " ]   " + selectedStyle.Render("[ " + no + " ]"))
		}
	} else if m.State == StateSyncing || m.State == StateConnecting {
		content.WriteString(m.Spinner.View() + " ")
		content.WriteString(m.Progress.View() + "\n")
	}

	s.WriteString(boxStyle.Render(content.String()) + "\n\n")

	if len(m.Logs) > 0 {
		for _, log := range m.Logs {
			s.WriteString(infoStyle.Render("  " + log) + "\n")
		}
		s.WriteString("\n")
	}

	s.WriteString(infoStyle.Render(i18n.T("quit")))

	// Center everything
	return lipgloss.Place(m.Width, m.Height, lipgloss.Center, lipgloss.Center, s.String())
}
