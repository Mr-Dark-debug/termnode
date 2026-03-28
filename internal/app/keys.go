package app

import "github.com/charmbracelet/bubbles/key"

// KeyBindings defines global and per-screen keybindings.
type KeyBindings struct {
	Dashboard key.Binding
	Services  key.Binding
	IoTLog    key.Binding
	Help      key.Binding
	Quit      key.Binding
	Up        key.Binding
	Down      key.Binding
	Enter     key.Binding
	Refresh   key.Binding
}

// DefaultKeyBindings returns the default keybindings.
func DefaultKeyBindings() KeyBindings {
	return KeyBindings{
		Dashboard: key.NewBinding(
			key.WithKeys("f1", "1"),
			key.WithHelp("F1/1", "dashboard"),
		),
		Services: key.NewBinding(
			key.WithKeys("f2", "2"),
			key.WithHelp("F2/2", "services"),
		),
		IoTLog: key.NewBinding(
			key.WithKeys("f3", "3"),
			key.WithHelp("F3/3", "iot log"),
		),
		Help: key.NewBinding(
			key.WithKeys("f4", "?"),
			key.WithHelp("F4/?", "help"),
		),
		Quit: key.NewBinding(
			key.WithKeys("q", "ctrl+c"),
			key.WithHelp("q/Ctrl+C", "quit"),
		),
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "down"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "toggle"),
		),
		Refresh: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "refresh"),
		),
	}
}
