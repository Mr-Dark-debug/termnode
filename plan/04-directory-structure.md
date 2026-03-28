# Directory Structure

## Complete File Layout

```
/storage/emulated/0/Projects/termnode/
│
├── CLAUDE.md                            # AI-assisted development guide
├── README.md                            # User-facing project documentation
├── main.go                              # Thin entry point (flags → db → app → run)
├── go.mod                               # Go module definition
├── go.sum                               # Dependency checksums (auto-generated)
├── Makefile                             # Build shortcuts
│
├── migrations/
│   └── 001_init.sql                     # SQL schema (events table + indexes)
│
├── plan/                                # Development planning documents
│   ├── 01-overview.md                   # Master plan and goals
│   ├── 02-research.md                   # Technology research
│   ├── 03-architecture.md              # System architecture
│   ├── 04-directory-structure.md       # This file
│   ├── 05-phase1-skeleton.md           # Phase 1 plan
│   ├── 06-phase2-database.md           # Phase 2 plan
│   ├── 07-phase3-dashboard.md          # Phase 3 plan
│   ├── 08-phase4-services.md           # Phase 4 plan
│   ├── 09-phase5-iot-bridge.md         # Phase 5 plan
│   ├── 10-phase6-polish.md             # Phase 6 plan
│   ├── 11-verification.md             # Testing plan
│   └── 12-decisions.md                # Architecture decisions
│
└── internal/                            # All application code (compiler-enforced privacy)
    │
    ├── app/                             # Root TUI model and routing
    │   ├── app.go                       # Root Bubble Tea model, screen router, Init/Update/View
    │   ├── keys.go                      # Global keybinding definitions (F1-F4, q, j/k, Enter, r)
    │   ├── messages.go                  # errMsg type for TUI error propagation
    │   └── styles.go                    # Re-exports theme styles for backward compatibility
    │
    ├── theme/                           # Centralized visual styling (no internal imports)
    │   └── styles.go                    # Lipgloss theme colors, panel/tab/status styles
    │
    ├── daemon/                          # Background service management
    │   ├── daemon.go                    # Wake-lock wrapper (termux-wake-lock/unlock)
    │   └── process.go                   # Process manager: sshd, HTTP file server, PID tracking
    │
    ├── hardware/                        # Hardware data polling from termux-api
    │   ├── types.go                     # BatteryInfo, NetworkInfo, CPUStats structs
    │   ├── battery.go                   # termux-battery-status JSON parser
    │   ├── network.go                   # termux-wifi-connectioninfo JSON parser
    │   ├── cpu.go                       # /proc/stat reader + sysfs thermal reader
    │   └── poller.go                    # tea.Tick-driven periodic polling + UpdateMsg type
    │
    ├── iot/                             # IoT webhook bridge
    │   ├── bridge.go                    # HTTP server lifecycle, webhook handler, channel bridge
    │   └── handler.go                   # Placeholder for future handler expansion
    │
    ├── db/                              # SQLite database layer
    │   ├── db.go                        # Connection open, WAL mode, embedded migration runner
    │   ├── models.go                    # IoTEvent struct definition
    │   ├── repository.go               # CRUD: Insert, Recent, ByTopic, Count, Purge
    │   └── migrations/
    │       └── 001_init.sql             # CREATE TABLE events + indexes
    │
    └── screen/                          # Individual TUI screen models
        ├── dashboard.go                 # Two-panel hardware display (battery + system)
        ├── services.go                  # Interactive service toggle list
        ├── iotlog.go                    # Scrollable IoT event log (bubbles/viewport)
        └── help.go                      # Static keybinding reference and features
```

---

## File Responsibilities

### Root Files

| File | Lines | Responsibility |
|------|-------|----------------|
| `main.go` | ~60 | Parse flags, open DB, create app model, run Bubble Tea |
| `CLAUDE.md` | ~80 | Project guide for AI-assisted development |
| `go.mod` | ~10 | Module definition and dependency versions |
| `Makefile` | ~20 | Build, run, test, clean shortcuts |

### internal/app/

| File | Responsibility |
|------|----------------|
| `app.go` | Root Model struct, screen enum, tab routing in Update(), tab bar rendering in View(), dependency wiring in New() |
| `keys.go` | KeyBindings struct, DefaultKeyBindings() with F1-F4, q, j/k, Enter, r |
| `messages.go` | errMsg type (for TUI error propagation) |
| `styles.go` | Re-exports from theme/ package for convenience |

### internal/theme/

| File | Responsibility |
|------|----------------|
| `styles.go` | Theme colors (Primary, Secondary, Accent, etc.), Panel/Tab/Status styles. **No imports from other internal packages** — breaks circular deps |

### internal/daemon/

| File | Responsibility |
|------|----------------|
| `daemon.go` | EnableWakeLock() / DisableWakeLock() wrapping termux commands |
| `process.go` | Manager struct: Start/Stop/Status for services, PID file tracking, Go HTTP file server, process signal handling |

### internal/hardware/

| File | Responsibility |
|------|----------------|
| `types.go` | BatteryInfo, NetworkInfo, CPUStats struct definitions (JSON-tagged) |
| `battery.go` | PollBattery() — exec termux-battery-status, parse JSON |
| `network.go` | PollNetwork() — exec termux-wifi-connectioninfo, parse JSON |
| `cpu.go` | PollCPU() — read /proc/stat for usage, sysfs for temperature, runtime for cores/arch |
| `poller.go` | Poller struct with interval, Start() returns tea.Cmd with tick loop, UpdateMsg type |

### internal/iot/

| File | Responsibility |
|------|----------------|
| `bridge.go` | Bridge struct with HTTP server, NewBridge(), Start/Stop, ListenCmd() for channel→TUI bridge, handleWebhook(), handleHealth(), IoTEventMsg type |
| `handler.go` | Placeholder for future handler expansion (e.g., MQTT config API) |

### internal/db/

| File | Responsibility |
|------|----------------|
| `db.go` | Open() function: ensure directory, sql.Open with WAL pragmas, embed + run migrations |
| `models.go` | IoTEvent struct: ID, Topic, Payload, Source, Timestamp |
| `repository.go` | Repository struct with Insert, Recent(limit), ByTopic(topic, limit), Count, Purge(olderThan) |

### internal/screen/

| File | Responsibility |
|------|----------------|
| `dashboard.go` | DashboardModel: holds hw data, two-panel View (battery bar + system stats), SetHardware() |
| `services.go` | ServicesModel: cursor, service list, Enter toggles via daemon.Manager, status indicators |
| `iotlog.go` | IoTLogModel: bubbles/viewport, AddEvent(), refreshContent() from DB, auto-scroll |
| `help.go` | HelpModel: static View with keybindings table, features list, webhook usage example |

---

## Import Graph (No Cycles)

```
main → app, db

app → screen, hardware, iot, db, theme

screen → theme, db, daemon, hardware

hardware → (only bubbletea + stdlib)

daemon → (only stdlib)

iot → db, (bubbletea + stdlib)

db → (only modernc.org/sqlite + stdlib)

theme → (only lipgloss)
```

**Zero circular dependencies.** Theme is a leaf package with no internal imports.
