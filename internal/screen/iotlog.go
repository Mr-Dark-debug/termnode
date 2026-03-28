package screen

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/Mr-Dark-debug/termnode/internal/db"
	"github.com/Mr-Dark-debug/termnode/internal/theme"
)

// IoTLogModel shows a scrollable log of IoT events received via webhook.
type IoTLogModel struct {
	width    int
	height   int
	viewport viewport.Model
	repo     *db.Repository
	events   []db.IoTEvent
	ready    bool
}

// NewIoTLogModel creates a new IoT log screen model.
func NewIoTLogModel(repo *db.Repository) IoTLogModel {
	return IoTLogModel{
		repo: repo,
	}
}

// SetSize updates the available drawing area.
func (m *IoTLogModel) SetSize(w, h int) {
	m.width = w
	m.height = h
	if !m.ready {
		m.viewport = viewport.New(w, h-4)
		m.ready = true
	} else {
		m.viewport.Width = w
		m.viewport.Height = h - 4
	}
	m.refreshContent()
}

// AddEvent appends a new IoT event and refreshes the view.
func (m *IoTLogModel) AddEvent(event db.IoTEvent) {
	m.events = append(m.events, event)
	m.refreshContent()
	m.viewport.GotoBottom()
}

func (m *IoTLogModel) refreshContent() {
	if m.repo == nil {
		return
	}

	events, err := m.repo.Recent(100)
	if err == nil {
		m.events = events
	}

	var b strings.Builder
	b.WriteString(theme.PanelTitleStyle.Render("IoT Event Log"))
	b.WriteString("\n\n")

	if len(m.events) == 0 {
		b.WriteString(theme.LabelStyle.Render("No events received yet."))
		b.WriteString("\n\n")
		b.WriteString(theme.LabelStyle.Render("POST data to: http://<ip>:8080/webhook/<topic>"))
	} else {
		for _, e := range m.events {
			timestamp := e.Timestamp.Format("15:04:05")
			topicStyle := lipgloss.NewStyle().Foreground(theme.Theme.Accent).Render(e.Topic)
			b.WriteString(fmt.Sprintf("%s [%s] %s",
				theme.LabelStyle.Render(timestamp),
				topicStyle,
				truncate(e.Payload, m.width-30),
			))
			b.WriteString("\n")
		}
	}

	m.viewport.SetContent(b.String())
}

func (m IoTLogModel) Init() tea.Cmd { return nil }

func (m IoTLogModel) Update(msg tea.Msg) (IoTLogModel, tea.Cmd) {
	if !m.ready {
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			m.viewport.LineUp(1)
		case "down", "j":
			m.viewport.LineDown(1)
		}
	default:
		var cmd tea.Cmd
		m.viewport, cmd = m.viewport.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m IoTLogModel) View() string {
	if !m.ready || m.width == 0 {
		return "Waiting for window size..."
	}

	var statusBar strings.Builder
	count := len(m.events)
	statusBar.WriteString(theme.LabelStyle.Render(
		fmt.Sprintf("Events: %d  |  ↑/↓ scroll  |  Webhook: POST /webhook/<topic>", count),
	))

	return m.viewport.View() + "\n" + statusBar.String()
}

func truncate(s string, maxLen int) string {
	s = strings.ReplaceAll(s, "\n", " ")
	s = strings.TrimSpace(s)
	if len(s) > maxLen {
		return s[:maxLen] + "..."
	}
	return s
}
