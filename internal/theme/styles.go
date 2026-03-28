package theme

import "github.com/charmbracelet/lipgloss"

// Theme holds the global color/style definitions for the TUI.
var Theme = struct {
	Primary    lipgloss.Color
	Secondary  lipgloss.Color
	Accent     lipgloss.Color
	Success    lipgloss.Color
	Danger     lipgloss.Color
	Muted      lipgloss.Color
	Background lipgloss.Color
}{
	Primary:    "#7C3AED",
	Secondary:  "#6366F1",
	Accent:     "#06B6D4",
	Success:    "#10B981",
	Danger:     "#EF4444",
	Muted:      "#6B7280",
	Background: "#1F2937",
}

// Tab styles.
var (
	ActiveTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(Theme.Primary).
			Padding(0, 2)

	InactiveTabStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#9CA3AF")).
			Background(lipgloss.Color("#374151")).
			Padding(0, 2)

	TabGapStyle = lipgloss.NewStyle().
			Background(lipgloss.Color("#111827"))

	TabSeparatorStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#374151")).
				SetString("─────────────────────────────────────────")
)

// Panel styles.
var (
	PanelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#4B5563")).
			Padding(1, 2)

	PanelTitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Theme.Accent)

	LabelStyle = lipgloss.NewStyle().
			Foreground(Theme.Muted)

	ValueStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#F9FAFB"))

	StatusOnStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(Theme.Success)

	StatusOffStyle = lipgloss.NewStyle().
			Foreground(Theme.Danger)
)
