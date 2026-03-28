package screen

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Mr-Dark-debug/termnode/internal/theme"
)

// HelpModel displays keybindings and project information.
type HelpModel struct {
	width   int
	height  int
	version string
}

// NewHelpModel creates a new help screen model.
func NewHelpModel(ver string) HelpModel {
	return HelpModel{version: ver}
}

// SetSize updates the available drawing area.
func (m *HelpModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

func (m HelpModel) Init() tea.Cmd { return nil }

func (m HelpModel) Update(msg tea.Msg) (HelpModel, tea.Cmd) {
	return m, nil
}

func (m HelpModel) View() string {
	if m.width == 0 {
		return "Waiting for window size..."
	}

	var b strings.Builder

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(theme.Theme.Primary).
		Render("TermNode")
	subtitle := theme.LabelStyle.Render("Headless Android Management TUI")
	b.WriteString(fmt.Sprintf("  %s %s\n", title, subtitle))
	b.WriteString(fmt.Sprintf("  %s\n", theme.LabelStyle.Render(fmt.Sprintf("Version: %s", m.version))))
	b.WriteString("\n")

	b.WriteString(theme.PanelTitleStyle.Render("Keybindings"))
	b.WriteString("\n\n")

	bindings := []struct{ key, desc string }{
		{"F1 / 1", "Dashboard tab"},
		{"F2 / 2", "Services tab"},
		{"F3 / 3", "IoT Log tab"},
		{"F4 / ?", "Help tab"},
		{"q / Ctrl+C", "Quit"},
		{"↑/↓ or j/k", "Navigate"},
		{"Enter", "Toggle / Activate"},
		{"r", "Refresh (Dashboard)"},
	}

	for _, kb := range bindings {
		b.WriteString(fmt.Sprintf("  %-14s %s\n",
			theme.ValueStyle.Render(kb.key),
			theme.LabelStyle.Render(kb.desc),
		))
	}

	b.WriteString("\n")

	b.WriteString(theme.PanelTitleStyle.Render("Features"))
	b.WriteString("\n\n")

	features := []struct{ name, desc string }{
		{"Hardware Dashboard", "Real-time battery, CPU, network monitoring via termux-api"},
		{"Service Manager", "Toggle wake-lock, SSH server, HTTP file server"},
		{"IoT Bridge", "HTTP webhook receiver logging to SQLite"},
		{"MQTT (optional)", "Subscribe to MQTT topics (build with -tags mqtt)"},
	}

	for _, f := range features {
		b.WriteString(fmt.Sprintf("  %-22s %s\n",
			lipgloss.NewStyle().Foreground(theme.Theme.Accent).Render(f.name),
			theme.LabelStyle.Render(f.desc),
		))
	}

	b.WriteString("\n")

	b.WriteString(theme.PanelTitleStyle.Render("IoT Webhook Usage"))
	b.WriteString("\n\n")
	b.WriteString(fmt.Sprintf("  %s\n", theme.ValueStyle.Render("curl -X POST http://<ip>:8080/webhook/<topic> -d '<json>'")))
	b.WriteString(fmt.Sprintf("  %s\n", theme.LabelStyle.Render("ESP32/Arduino can POST sensor data to this endpoint")))

	panel := theme.PanelStyle.Width(m.width - 4).Render(b.String())
	return panel
}
