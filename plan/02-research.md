# Research Findings

## 1. Bubble Tea Framework

### The Elm Architecture

Bubble Tea is based on **The Elm Architecture** — a pattern for building interactive programs with three core components:

- **Model**: The application state (a Go struct)
- **Update**: A function that takes a message and the current model, and returns a new model + optional commands
- **View**: A function that renders the model as a string

The runtime loop:
```
Init() → Cmd → Msg → Update(Model, Msg) → (NewModel, Cmd) → Msg → ...
```

### Key Concepts

| Concept | Type | Purpose |
|---------|------|---------|
| `tea.Model` | Interface | `Init()`, `Update()`, `View()` methods |
| `tea.Msg` | Interface | Any type representing an event |
| `tea.Cmd` | Function | Returns a `tea.Msg`, executed in a goroutine |
| `tea.KeyMsg` | Struct | Keyboard input events |
| `tea.WindowSizeMsg` | Struct | Terminal resize events |

### Multi-Screen Pattern

The recommended pattern for multi-screen apps is the **Composite Root Model**:

```go
type Model struct {
    currentScreen Screen          // enum: dashboard, services, etc.
    dashboard     DashboardModel  // sub-model for dashboard
    services      ServicesModel   // sub-model for services
    // ...
}
```

The root `Update()` intercepts tab-switching keys, then delegates other messages to the active screen. The root `View()` renders the tab bar + active screen's view.

### Real-Time Updates

Use `tea.Tick(interval, callback)` to create periodic commands:

```go
func (m Model) Init() tea.Cmd {
    return tea.Tick(5 * time.Second, func(t time.Time) tea.Msg {
        return pollHardwareMsg{}
    })
}
```

After processing the tick message in `Update()`, return a new `tea.Tick` command to continue the loop.

### Background Processes

`tea.Cmd` functions run in goroutines. Results are delivered to `Update()` on the main goroutine — no mutexes needed for UI state.

For channel-based communication (e.g., IoT bridge → TUI):

```go
func listenForEvents(ch chan Event) tea.Cmd {
    return func() tea.Msg {
        event := <-ch
        return eventMsg{event}
    }
}
```

### Layout

Lip Gloss provides layout combinators:
- `lipgloss.JoinHorizontal(top, a, b, ...)` — side by side
- `lipgloss.JoinVertical(left, a, b, ...)` — stacked
- `lipgloss.Width(s)` / `lipgloss.Height(s)` — measure rendered size

### Companion Packages

| Package | Import | Purpose |
|---------|--------|---------|
| Bubble Tea | `github.com/charmbracelet/bubbletea` | TUI framework |
| Lip Gloss | `github.com/charmbracelet/lipgloss` | Styling & layout |
| Bubbles | `github.com/charmbracelet/bubbles` | Pre-built components |
| Huh | `github.com/charmbracelet/huh` | Forms & surveys |

---

## 2. Termux Environment

### Available termux-api Commands

| Command | Output | Fields |
|---------|--------|--------|
| `termux-battery-status` | JSON | `percentage`, `status` (CHARGING/DISCHARGING), `temperature`, `health`, `current`, `plugged` |
| `termux-wifi-connectioninfo` | JSON | `ip`, `ssid`, `bssid`, `mac` |
| `termux-sensor` | JSON | Various sensor readings |
| `termux-wake-lock` | (none) | Acquires wake lock |
| `termux-wake-unlock` | (none) | Releases wake lock |

### Service Management

- **sshd**: Installed via `pkg install openssh`, runs as `sshd` command
- **HTTP file server**: Can use Python's `http.server` or Go's `net/http.FileServer`
- **Background processes**: Termux supports standard Linux process management

### Filesystem Layout

```
/data/data/com.termux/files/
├── usr/              # $PREFIX — Termux's "system" directory
│   ├── bin/          # Executables (go, gcc, sshd, termux-*)
│   ├── lib/          # Libraries
│   ├── tmp/          # Temporary files
│   └── etc/          # Configuration
└── home/             # $HOME — User's home directory
```

External storage access: `/storage/emulated/0/`

### Go in Termux

- Install: `pkg install golang`
- Go compiles natively for `aarch64` in Termux
- No cross-compilation needed — build directly on device
- Pure Go packages work without issues
- CGO may require additional packages (`pkg install clang`)

---

## 3. Go SQLite Options

### modernc.org/sqlite (Recommended)

- **Pure Go** — no CGO, no native C library needed
- Implements `database/sql/driver` interface
- Supports `linux/arm64` (critical for Termux)
- Supports WAL mode, foreign keys, all standard SQLite features
- Current version tracks SQLite 3.51.3
- Used via `import _ "modernc.org/sqlite"` + `sql.Open("sqlite", dsn)`

### DSN Pragmas

```
file:path?_pragma=journal_mode(WAL)&_pragma=foreign_keys(1)
```

---

## 4. Go HTTP Server (Go 1.22+)

Go 1.22 introduced enhanced routing in `net/http`:

```go
mux := http.NewServeMux()
mux.HandleFunc("POST /webhook/{topic}", handler)  // method + path pattern
mux.HandleFunc("GET /health", healthHandler)

topic := r.PathValue("topic")  // extract path parameter
```

No external router needed. This keeps dependencies minimal.

---

## 5. Go Project Layout (Official Guidance)

From the official Go documentation:

- **`cmd/`** — Entry points (optional for single-binary projects)
- **`internal/`** — Compiler-enforced privacy (recommended for all app code)
- **`pkg/`** — Skip for single-binary applications
- **`main.go`** at root is acceptable for single-binary projects

The Go team's recommendation for server/binary projects:
> "Server projects typically won't have packages for export. Therefore, it's recommended to keep the Go packages implementing the server's logic in the `internal` directory."

---

## 6. Key Dependencies (All Pure Go)

| Dependency | Import Path | Purpose | CGO |
|-----------|-------------|---------|-----|
| Bubble Tea | `github.com/charmbracelet/bubbletea` | TUI framework | No |
| Lip Gloss | `github.com/charmbracelet/lipgloss` | Terminal styling | No |
| Bubbles | `github.com/charmbracelet/bubbles` | Viewport, keybindings | No |
| SQLite | `modernc.org/sqlite` | Database | No |
| Paho MQTT | `github.com/eclipse-paho/paho.mqtt.golang` | MQTT client | No |

All dependencies compile to pure Go binaries — zero native dependencies.
