# Phase 3: Hardware Dashboard

## Goal

Implement real-time hardware monitoring by polling `termux-api` commands and displaying the data in a two-panel TUI dashboard with auto-refresh.

---

## Data Sources

| Data | Source | Command / File | Output |
|------|--------|---------------|--------|
| Battery | termux-api | `termux-battery-status` | JSON |
| WiFi | termux-api | `termux-wifi-connectioninfo` | JSON |
| CPU Usage | procfs | `/proc/stat` | Text |
| CPU Temp | sysfs | `/sys/class/thermal/thermal_zone*/temp` | Number |
| CPU Cores | Go runtime | `runtime.NumCPU()` | Integer |
| Architecture | Go runtime | `runtime.GOARCH` | String |

---

## Files

### internal/hardware/types.go

```go
type BatteryInfo struct {
    Percentage  int     `json:"percentage"`
    Status      string  `json:"status"`      // CHARGING, DISCHARGING, FULL, NOT_CHARGING
    Temperature float64 `json:"temperature"` // Celsius
    Health      string  `json:"health"`      // GOOD, OVERHEAT, etc.
    Current     int     `json:"current"`     // mA
    Plugged     string  `json:"plugged"`     // AC, USB, UNKNOWN
}

type NetworkInfo struct {
    IP    string `json:"ip"`
    SSID  string `json:"ssid"`
    BSSID string `json:"bssid"`
    MAC   string `json:"mac"`
}

type CPUStats struct {
    UsagePercent float64
    CoreCount    int
    Temperature  float64
    Arch         string  // "aarch64"
}
```

### internal/hardware/battery.go

```go
func PollBattery() (BatteryInfo, error)
```
- Executes `termux-battery-status` via `exec.Command`
- Captures JSON stdout
- Unmarshals into `BatteryInfo`
- Returns zero-value + error if termux-api not installed

### internal/hardware/network.go

```go
func PollNetwork() (NetworkInfo, error)
```
- Executes `termux-wifi-connectioninfo` via `exec.Command`
- Same JSON parse pattern as battery

### internal/hardware/cpu.go

```go
func PollCPU() (CPUStats, error)
```
- Reads `/proc/stat` line 1 for total/idle CPU ticks
- Calculates usage percentage from tick values
- Reads `/sys/class/thermal/thermal_zone0/temp` for temperature
- Handles millidegree values (divide by 1000 if > 1000)
- Gets core count and arch from `runtime` package

### internal/hardware/poller.go

```go
type Poller struct { interval time.Duration }
type UpdateMsg struct {
    Battery BatteryInfo
    Network NetworkInfo
    CPU     CPUStats
    Err     error
}

func NewPoller(seconds int) *Poller
func (p *Poller) Start() tea.Cmd  // returns tick + poll command
```

The poller creates a continuous loop:
1. `tea.Tick(5s)` fires a tick message
2. The tick triggers `PollAll()` which runs all three poll functions
3. Returns `UpdateMsg` with collected data
4. `app.Model.Update()` passes data to dashboard and re-issues `poller.Start()`

### internal/screen/dashboard.go

```go
type DashboardModel struct { /* hw data fields */ }

func (m DashboardModel) View() string {
    // Two-panel layout:
    // Left: Battery (percentage, status, health, temp, current, bar)
    // Right: System (WiFi info) + CPU (usage, cores, temp, arch)
}
```

**Visual Elements:**
- Battery bar: `[████████░░░░░░░░]` with color coding (green > 50%, yellow > 20%, red < 20%)
- Status indicators: Green for CHARGING, red for DISCHARGING
- Temperature in Celsius
- IP address prominent for headless server use

---

## Polling Loop Sequence

```
┌──────────┐     tea.Tick(5s)     ┌───────────┐
│ Poller   │ ──────────────────▶  │ goroutine │
│ .Start() │                      │  PollAll  │
└──────────┘                      └─────┬─────┘
      ▲                                 │
      │ re-issue                        ▼
      │ after                    ┌───────────┐
      │ processing               │ UpdateMsg │
      │                          └─────┬─────┘
      │                                │
      └────────────────────────────────┘
         app.Update() dispatches
         to dashboard, then
         calls poller.Start() again
```

---

## Error Handling

- If `termux-api` is not installed: Dashboard shows "termux-api unavailable" instead of crashing
- If WiFi is disconnected: Network fields show empty/unavailable
- If `/proc/stat` can't be read: CPU usage shows 0%
- Individual poll failures don't block other polls — partial data is shown

---

## Termux-API Installation

Users must install:
```bash
pkg install termux-api
```

And install the **Termux:API** Android app from F-Droid (not Play Store — the Play Store version is outdated).

After installation, grant the necessary permissions when prompted.

---

## Verification

```bash
# Test termux-api directly
termux-battery-status
# Expected: {"percentage":85,"status":"CHARGING","temperature":32.1,...}

termux-wifi-connectioninfo
# Expected: {"ip":"192.168.1.100","ssid":"MyWiFi",...}

# Build and run TermNode
go build -o termnode . && ./termnode

# Expected: Dashboard tab shows live data, refreshing every 5 seconds
# Battery bar changes color based on percentage
# CPU usage fluctuates
# WiFi shows current connection info
```
