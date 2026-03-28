package screen

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Mr-Dark-debug/termnode/internal/daemon"
	"github.com/Mr-Dark-debug/termnode/internal/theme"
)

// errMsg is a local error message type for the TUI.
type errMsg struct{ Err error }

// ServicesModel shows toggleable services (wake-lock, sshd, HTTP file server).
type ServicesModel struct {
	width    int
	height   int
	cursor   int
	services []daemon.Service
	manager  *daemon.Manager
}

// NewServicesModel creates a new services screen model.
func NewServicesModel() ServicesModel {
	mgr := daemon.NewManager()
	return ServicesModel{
		manager: mgr,
		services: []daemon.Service{
			{Name: "Wake Lock", Key: "wakelock", Desc: "Prevent Android battery optimization"},
			{Name: "SSH Server", Key: "sshd", Desc: "OpenSSH remote access"},
			{Name: "HTTP File Server", Key: "httpfs", Desc: "Serve files over HTTP"},
		},
	}
}

// SetSize updates the available drawing area.
func (m *ServicesModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m ServicesModel) Init() tea.Cmd { return nil }

func (m ServicesModel) Update(msg tea.Msg) (ServicesModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(m.services)-1 {
				m.cursor++
			}
		case "enter":
			svc := &m.services[m.cursor]
			running, _ := m.manager.Status(svc.Key)
			var err error
			if running {
				err = m.manager.Stop(svc.Key)
			} else {
				err = m.manager.Start(svc.Key)
			}
			if err != nil {
				return m, func() tea.Msg { return errMsg{Err: err} }
			}
			svc.Running, _ = m.manager.Status(svc.Key)
		}
	}

	// Refresh all statuses
	for i := range m.services {
		m.services[i].Running, _ = m.manager.Status(m.services[i].Key)
	}

	return m, nil
}

func (m ServicesModel) View() string {
	if m.width == 0 {
		return "Waiting for window size..."
	}

	var b strings.Builder
	b.WriteString(theme.PanelTitleStyle.Render("Service Management"))
	b.WriteString("\n\n")

	for i, svc := range m.services {
		cursor := "  "
		if i == m.cursor {
			cursor = lipgloss.NewStyle().Foreground(theme.Theme.Accent).Render("> ")
		}

		statusIcon := theme.StatusOffStyle.Render("● OFF")
		statusText := theme.StatusOffStyle.Render("Stopped")
		if svc.Running {
			statusIcon = theme.StatusOnStyle.Render("● ON ")
			statusText = theme.StatusOnStyle.Render("Running")
		}

		nameStyle := theme.ValueStyle
		if i == m.cursor {
			nameStyle = lipgloss.NewStyle().Bold(true).Foreground(theme.Theme.Primary)
		}

		line := fmt.Sprintf("%s %-20s %s  %s",
			cursor,
			nameStyle.Render(svc.Name),
			statusIcon,
			theme.LabelStyle.Render(svc.Desc),
		)
		b.WriteString(line)
		b.WriteString("\n")

		statusLine := fmt.Sprintf("    Status: %s", statusText)
		b.WriteString(statusLine)
		b.WriteString("\n\n")
	}

	b.WriteString("\n")
	b.WriteString(theme.LabelStyle.Render("Press Enter to toggle  |  ↑/↓ or j/k to navigate"))

	panel := theme.PanelStyle.Width(m.width - 4).Render(b.String())
	return panel
}
