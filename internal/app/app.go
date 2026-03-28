package app

import (
	"context"
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Mr-Dark-debug/termnode/internal/db"
	"github.com/Mr-Dark-debug/termnode/internal/hardware"
	"github.com/Mr-Dark-debug/termnode/internal/iot"
	"github.com/Mr-Dark-debug/termnode/internal/screen"
)

// Screen represents the active TUI screen.
type Screen int

const (
	ScreenDashboard Screen = iota
	ScreenServices
	ScreenIoTLog
	ScreenHelp
)

// Model is the root Bubble Tea model that routes between screens.
type Model struct {
	currentScreen Screen
	dashboard     screen.DashboardModel
	services      screen.ServicesModel
	iotLog        screen.IoTLogModel
	help          screen.HelpModel
	keys          KeyBindings
	width         int
	height        int
	bridge        *iot.Bridge
	poller        *hardware.Poller
	cancelFn      context.CancelFunc
	version       string
}

// New creates the root model with all dependencies wired.
func New(repo *db.Repository, port string, ver string) Model {
	eventsCh := make(chan db.IoTEvent, 64)
	bridge := iot.NewBridge(repo, port, eventsCh)
	poller := hardware.NewPoller(5)

	return Model{
		currentScreen: ScreenDashboard,
		dashboard:     screen.NewDashboardModel(),
		services:      screen.NewServicesModel(),
		iotLog:        screen.NewIoTLogModel(repo),
		help:          screen.NewHelpModel(ver),
		keys:          DefaultKeyBindings(),
		bridge:        bridge,
		poller:        poller,
		version:       ver,
	}
}

// Init starts the TUI and background services.
func (m Model) Init() tea.Cmd {
	return tea.Batch(
		m.poller.Start(),
		m.bridge.ListenCmd(),
	)
}

// Update handles incoming messages and routes them to the active screen.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.dashboard.SetSize(msg.Width, msg.Height-4)
		m.services.SetSize(msg.Width, msg.Height-4)
		m.iotLog.SetSize(msg.Width, msg.Height-4)
		return m, nil

	case tea.KeyMsg:
		switch {
		case keyMatches(m.keys.Dashboard, msg):
			m.currentScreen = ScreenDashboard
			return m, m.poller.Start()
		case keyMatches(m.keys.Services, msg):
			m.currentScreen = ScreenServices
			return m, nil
		case keyMatches(m.keys.IoTLog, msg):
			m.currentScreen = ScreenIoTLog
			return m, nil
		case keyMatches(m.keys.Help, msg):
			m.currentScreen = ScreenHelp
			return m, nil
		case keyMatches(m.keys.Quit, msg):
			return m, tea.Quit
		}

	case hardware.UpdateMsg:
		m.dashboard.SetHardware(msg.Battery, msg.Network, msg.CPU, msg.Err)
		return m, m.poller.Start()

	case iot.IoTEventMsg:
		m.iotLog.AddEvent(msg.Event)
		return m, m.bridge.ListenCmd()
	}

	// Delegate to active screen
	var cmd tea.Cmd
	switch m.currentScreen {
	case ScreenDashboard:
		m.dashboard, cmd = m.dashboard.Update(msg)
	case ScreenServices:
		m.services, cmd = m.services.Update(msg)
	case ScreenIoTLog:
		m.iotLog, cmd = m.iotLog.Update(msg)
	case ScreenHelp:
		m.help, cmd = m.help.Update(msg)
	}
	return m, cmd
}

// View renders the tab bar and the active screen.
func (m Model) View() string {
	if m.width == 0 {
		return "Loading..."
	}

	tabBar := m.renderTabBar()
	sep := TabSeparatorStyle.String()

	var content string
	switch m.currentScreen {
	case ScreenDashboard:
		content = m.dashboard.View()
	case ScreenServices:
		content = m.services.View()
	case ScreenIoTLog:
		content = m.iotLog.View()
	case ScreenHelp:
		content = m.help.View()
	}

	return tabBar + "\n" + sep + "\n" + content
}

func (m Model) renderTabBar() string {
	tabs := []struct {
		name   string
		screen Screen
	}{
		{"Dashboard", ScreenDashboard},
		{"Services", ScreenServices},
		{"IoT Log", ScreenIoTLog},
		{"Help", ScreenHelp},
	}

	var rendered []string
	for i, t := range tabs {
		style := InactiveTabStyle
		if m.currentScreen == t.screen {
			style = ActiveTabStyle
		}
		label := fmt.Sprintf(" %d:%s ", i+1, t.name)
		rendered = append(rendered, style.Render(label))
	}

	tabRow := lipgloss.JoinHorizontal(lipgloss.Bottom, rendered...)
	gap := TabGapStyle.Render(strings.Repeat(" ", max(0, m.width-lipgloss.Width(tabRow))))
	return tabRow + gap
}

func keyMatches(b interface{ Matches(tea.KeyMsg) bool }, msg tea.KeyMsg) bool {
	return b.Matches(msg)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
