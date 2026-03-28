package app

// Re-export theme styles for backward compatibility.
// The canonical styles live in internal/theme.
import "github.com/Mr-Dark-debug/termnode/internal/theme"

var (
	ActiveTabStyle   = theme.ActiveTabStyle
	InactiveTabStyle = theme.InactiveTabStyle
	TabGapStyle      = theme.TabGapStyle
	TabSeparatorStyle = theme.TabSeparatorStyle
	PanelStyle       = theme.PanelStyle
	PanelTitleStyle  = theme.PanelTitleStyle
	LabelStyle       = theme.LabelStyle
	ValueStyle       = theme.ValueStyle
	StatusOnStyle    = theme.StatusOnStyle
	StatusOffStyle   = theme.StatusOffStyle
	Theme            = theme.Theme
)
