package screen

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Mr-Dark-debug/termnode/internal/hardware"
	"github.com/Mr-Dark-debug/termnode/internal/theme"
)

// DashboardModel displays real-time hardware info in two panels.
type DashboardModel struct {
	width   int
	height  int
	battery hardware.BatteryInfo
	network hardware.NetworkInfo
	cpu     hardware.CPUStats
	err     error
}

// NewDashboardModel creates a new dashboard screen model.
func NewDashboardModel() DashboardModel {
	return DashboardModel{}
}

// SetSize updates the available drawing area.
func (m *DashboardModel) SetSize(w, h int) {
	m.width = w
	m.height = h
}

// SetHardware updates the displayed hardware data.
func (m *DashboardModel) SetHardware(b hardware.BatteryInfo, n hardware.NetworkInfo, c hardware.CPUStats, err error) {
	m.battery = b
	m.network = n
	m.cpu = c
	m.err = err
}

func (m DashboardModel) Init() tea.Cmd { return nil }

func (m DashboardModel) Update(msg tea.Msg) (DashboardModel, tea.Cmd) {
	return m, nil
}

func (m DashboardModel) View() string {
	if m.width == 0 {
		return "Waiting for window size..."
	}

	panelWidth := (m.width - 6) / 2

	leftPanel := m.renderBatteryPanel(panelWidth)
	rightPanel := m.renderSystemPanel(panelWidth)

	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, rightPanel)
}

func (m DashboardModel) renderBatteryPanel(w int) string {
	s := theme.PanelStyle.Width(w)

	var b strings.Builder
	b.WriteString(theme.PanelTitleStyle.Render("Battery"))
	b.WriteString("\n\n")

	if m.battery.Percentage == 0 && m.err != nil {
		b.WriteString(theme.LabelStyle.Render("Status: "))
		b.WriteString(theme.StatusOffStyle.Render("termux-api unavailable"))
		b.WriteString("\n")
	} else {
		b.WriteString(theme.LabelStyle.Render("Level:    "))
		b.WriteString(theme.ValueStyle.Render(fmt.Sprintf("%d%%", m.battery.Percentage)))
		b.WriteString("\n")

		b.WriteString(theme.LabelStyle.Render("Status:   "))
		statusStyle := theme.StatusOnStyle
		if m.battery.Status == "DISCHARGING" {
			statusStyle = theme.StatusOffStyle
		}
		b.WriteString(statusStyle.Render(m.battery.Status))
		b.WriteString("\n")

		b.WriteString(theme.LabelStyle.Render("Health:   "))
		b.WriteString(theme.ValueStyle.Render(m.battery.Health))
		b.WriteString("\n")

		b.WriteString(theme.LabelStyle.Render("Temp:     "))
		b.WriteString(theme.ValueStyle.Render(fmt.Sprintf("%.1f°C", m.battery.Temperature)))
		b.WriteString("\n")

		b.WriteString(theme.LabelStyle.Render("Current:  "))
		b.WriteString(theme.ValueStyle.Render(fmt.Sprintf("%d mA", m.battery.Current)))
		b.WriteString("\n")

		b.WriteString("\n")
		b.WriteString(m.renderBatteryBar(w - 6))
	}

	return s.Render(b.String())
}

func (m DashboardModel) renderBatteryBar(w int) string {
	if w <= 0 {
		w = 20
	}
	filled := m.battery.Percentage * w / 100
	if filled > w {
		filled = w
	}

	barStyle := theme.StatusOnStyle
	if m.battery.Percentage < 20 {
		barStyle = theme.StatusOffStyle
	} else if m.battery.Percentage < 50 {
		barStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#F59E0B"))
	}

	bar := barStyle.Render(strings.Repeat("█", filled)) +
		lipgloss.NewStyle().Foreground(lipgloss.Color("#374151")).Render(strings.Repeat("░", w-filled))
	return fmt.Sprintf("[%s]", bar)
}

func (m DashboardModel) renderSystemPanel(w int) string {
	s := theme.PanelStyle.Width(w)

	var b strings.Builder
	b.WriteString(theme.PanelTitleStyle.Render("System"))
	b.WriteString("\n\n")

	if m.network.IP == "" && m.err != nil {
		b.WriteString(theme.LabelStyle.Render("Network: "))
		b.WriteString(theme.StatusOffStyle.Render("unavailable"))
		b.WriteString("\n")
	} else {
		b.WriteString(theme.LabelStyle.Render("WiFi SSID: "))
		b.WriteString(theme.ValueStyle.Render(m.network.SSID))
		b.WriteString("\n")

		b.WriteString(theme.LabelStyle.Render("Local IP:  "))
		b.WriteString(theme.ValueStyle.Render(m.network.IP))
		b.WriteString("\n")

		b.WriteString(theme.LabelStyle.Render("BSSID:     "))
		b.WriteString(theme.ValueStyle.Render(m.network.BSSID))
		b.WriteString("\n\n")
	}

	b.WriteString(theme.PanelTitleStyle.Render("CPU"))
	b.WriteString("\n\n")
	b.WriteString(theme.LabelStyle.Render("Usage:     "))
	b.WriteString(theme.ValueStyle.Render(fmt.Sprintf("%.1f%%", m.cpu.UsagePercent)))
	b.WriteString("\n")

	b.WriteString(theme.LabelStyle.Render("Cores:     "))
	b.WriteString(theme.ValueStyle.Render(fmt.Sprintf("%d", m.cpu.CoreCount)))
	b.WriteString("\n")

	b.WriteString(theme.LabelStyle.Render("Temp:      "))
	b.WriteString(theme.ValueStyle.Render(fmt.Sprintf("%.1f°C", m.cpu.Temperature)))
	b.WriteString("\n")

	b.WriteString(theme.LabelStyle.Render("Arch:      "))
	b.WriteString(theme.ValueStyle.Render(m.cpu.Arch))
	b.WriteString("\n")

	return s.Render(b.String())
}
