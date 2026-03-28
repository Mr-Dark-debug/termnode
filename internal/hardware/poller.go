package hardware

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

// Poller periodically polls hardware data using tea.Tick.
type Poller struct {
	interval time.Duration
}

// NewPoller creates a new hardware poller with the given interval in seconds.
func NewPoller(seconds int) *Poller {
	return &Poller{interval: time.Duration(seconds) * time.Second}
}

// Start returns a tea.Cmd that waits for the interval then polls all hardware.
func (p *Poller) Start() tea.Cmd {
	return tea.Tick(p.interval, func(_ time.Time) tea.Msg {
		return p.pollAll()
	})
}

func (p *Poller) pollAll() UpdateMsg {
	battery, bErr := PollBattery()
	network, nErr := PollNetwork()
	cpu, cErr := PollCPU()

	var err error
	if bErr != nil {
		err = bErr
	} else if nErr != nil {
		err = nErr
	} else if cErr != nil {
		err = cErr
	}

	return UpdateMsg{
		Battery: battery,
		Network: network,
		CPU:     cpu,
		Err:     err,
	}
}

// UpdateMsg is the tea.Msg sent when hardware data is refreshed.
type UpdateMsg struct {
	Battery BatteryInfo
	Network NetworkInfo
	CPU     CPUStats
	Err     error
}
