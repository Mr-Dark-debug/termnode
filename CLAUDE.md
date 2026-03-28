# TermNode

## Project Overview

TermNode is a unified, high-performance Terminal User Interface (TUI) for headless Android management in Termux. It compiles to a single Go binary that provides hardware monitoring, service management, network toggles, and an IoT webhook bridge — all from one terminal interface.

## Architecture

- **TUI Framework**: Bubble Tea (Elm Architecture) — Model/Update/View cycle
- **Styling**: Lip Gloss for terminal styling and layout
- **Components**: Bubbles for pre-built TUI components (viewport, spinner, table)
- **Database**: SQLite via modernc.org/sqlite (pure Go, no CGO)
- **Migrations**: pressly/goose with embedded SQL files
- **IoT Bridge**: Go 1.22+ net/http webhook receiver, optional MQTT via build tag
- **Module**: `github.com/Mr-Dark-debug/termnode`

### Root Model (Composite Pattern)

```
app.Model (root)
  ├── currentScreen (screen enum: dashboard | services | iotlog | help)
  ├── dashboard Model
  ├── services  Model
  ├── iotLog    Model
  ├── help      Model
  ├── repo      *db.Repository
  └── bridge    *iot.Bridge
```

### Data Flow Patterns

- **Hardware**: tea.Tick → exec.Command(termux-api) → JSON parse → hardwareUpdateMsg
- **IoT**: HTTP POST → channel → newIoTEventMsg → viewport refresh
- **Services**: keypress → daemon.Manager toggle → exec.Command → serviceStatusMsg

## Key Conventions

- All application code in `internal/` (compiler-enforced privacy)
- Screen models implement Init/Update/View
- Global keybindings in `internal/app/keys.go`
- Styling centralized in `internal/app/styles.go`
- Error messages flow as `errMsg` through the TUI, never fmt.Println
- Database auto-migrates on startup via embedded Goose migrations

## Operational Rules

- ALL project files must remain inside `/storage/emulated/0/Projects/termnode/`
- Database auto-created at `$HOME/.termnode/termnode.db`
- PID files for managed services at `$HOME/.termnode/*.pid`
- No CGO required — all dependencies are pure Go

## Build & Run

```bash
go build -o termnode .                    # Default build
go build -tags mqtt -o termnode .         # With MQTT support
./termnode                                # Run with defaults
./termnode -db path.db -port :9090        # Custom options
./termnode -debug                         # With debug logging
```

## Dependencies

- github.com/charmbracelet/bubbletea
- github.com/charmbracelet/lipgloss
- github.com/charmbracelet/bubbles
- modernc.org/sqlite
- github.com/pressly/goose/v3
- github.com/eclipse-paho/paho.mqtt.golang (optional, build tag: mqtt)

## Testing

```bash
go test ./...
```

Database tests use `:memory:` SQLite DSN.
