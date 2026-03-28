# Phase 6: Polish, Makefile, and Future Features

## Goal

Final polish pass: graceful shutdown, error display, Makefile optimization, and preparation for optional MQTT support and future features.

---

## Completed in This Phase

### Makefile

Build shortcuts for common tasks:

| Command | Action |
|---------|--------|
| `make` | Build the binary |
| `make run` | Build and run |
| `make debug` | Build and run with debug logging |
| `make test` | Run all tests |
| `make clean` | Remove binary and logs |
| `make deps` | Tidy dependencies |
| `make build-mqtt` | Build with MQTT support |

### Graceful Shutdown (Reserved)

The `app.Model` struct has a `cancelFn context.CancelFunc` field reserved for coordinated shutdown:

```go
// Future implementation:
func (m Model) shutdown() tea.Cmd {
    return func() tea.Msg {
        if m.cancelFn != nil {
            m.cancelFn()  // signal all goroutines
        }
        m.bridge.Stop()  // close HTTP server
        // stop poller
        // clean up PID files
        return tea.Quit()
    }
}
```

### Error Display (Reserved)

The `errMsg` type is defined in `app/messages.go` and `screen/services.go`. Future implementation will add a status bar at the bottom showing transient errors:

```go
// Future: error bar in View()
if m.lastError != nil {
    errorBar := theme.StatusOffStyle.Render(fmt.Sprintf("Error: %v", m.lastError))
    content += "\n" + errorBar
}
```

---

## Future Features

### MQTT Support (Build Tag: `mqtt`)

File: `internal/iot/mqtt.go`

```go
//go:build mqtt

package iot

import (
    "github.com/eclipse-paho/paho.mqtt.golang"
)

type MQTTSubscriber struct {
    client  mqtt.Client
    broker  string
    topics  []string
    events  chan db.IoTEvent
}

func NewMQTTSubscriber(broker string, topics []string, events chan db.IoTEvent) *MQTTSubscriber
func (s *MQTTSubscriber) Connect() error
func (s *MQTTSubscriber) Disconnect()
```

**Build with MQTT:**
```bash
go build -tags mqtt -o termnode .
```

**Why behind a build tag?**
- MQTT adds significant binary size (~2MB) and dependency complexity
- Most users only need HTTP webhooks
- Default build stays lean
- Opt-in for users who need MQTT

**MQTT Configuration (Future):**
```bash
./termnode -mqtt-broker tcp://192.168.1.50:1883 -mqtt-topic sensors/#
```

### Tailscale/Ngrok Tunnel Management

```go
// Future service in daemon/
type Tunnel struct {
    Provider string  // "tailscale" or "ngrok"
    Status   string  // "connected", "disconnected"
    URL      string  // public URL
}
```

Would add a new row in the Services screen:
- Start/stop Tailscale: `tailscale up`
- Start/stop Ngrok: `ngrok http 8080`
- Display the public tunnel URL

### Process Supervision with Auto-Restart

Currently, if a managed process (sshd, file server) crashes, it stays dead. Future enhancement:

```go
type Manager struct {
    // ... existing fields ...
    autoRestart map[string]bool
    restartDelay time.Duration
}

// In background goroutine:
for {
    if m.autoRestart[key] && !isAlive(pid) {
        m.Start(key)
        time.Sleep(m.restartDelay)
    }
}
```

### Configuration File

```yaml
# ~/.termnode/config.yaml
database:
  path: ~/.termnode/termnode.db

http:
  webhook_port: :8080
  fileserver_port: :8081

polling:
  interval: 5s

mqtt:
  enabled: false
  broker: tcp://localhost:1883
  topics:
    - sensors/#
```

### GitHub Actions CI/CD

```yaml
# .github/workflows/release.yml
name: Release
on:
  push:
    tags: ['v*']
jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with: { go-version: '1.22' }
      - run: GOOS=linux GOARCH=arm64 go build -o termnode-arm64 .
      - uses: softprops/action-gh-release@v1
        with:
          files: termnode-arm64
```

### Real-Time Charts

Using `asciigraph` for in-terminal charts:

```go
import "github.com/guptarohit/asciigraph"

// CPU usage over last 60 readings
chart := asciigraph.Plot(cpuHistory, asciigraph.Width(50), asciigraph.Height(5))
```

Would show:
```
 12.3 ┤      ╭╮    ╭╮
 10.0 ┤   ╭╮╭╯╰╮  ╭╯╰╮
  7.5 ┤  ╭╯╰╯  ╰╮╭╯  ╰╮
  5.0 ┤ ╭╯      ╰╯    ╰
  2.5 ┤╭╯
  0.0 ┼──────────────────
```

---

## Performance Optimization Notes

### Binary Size

Current estimated binary size: ~15-20MB (with modernc.org/sqlite)

Reduction strategies:
```bash
# Strip debug symbols
go build -ldflags="-s -w" -o termnode .

# With UPX compression (if available)
upx --best termnode

# Expected final size: ~5-8MB
```

### Memory Usage

Expected memory footprint:
- Bubble Tea TUI: ~5MB
- SQLite connection: ~2MB
- HTTP server: ~1MB
- Total: ~10MB (negligible for a phone with 2-8GB RAM)

### CPU Usage

- Polling interval of 5s means minimal CPU impact
- `termux-api` commands are lightweight subprocess spawns
- Viewport rendering is only on change, not continuous

---

## Final Verification Checklist

- [ ] `go build -o termnode .` compiles without errors
- [ ] Binary size under 25MB
- [ ] `./termnode` launches TUI in under 1 second
- [ ] All 4 tabs switch correctly (F1-F4)
- [ ] Dashboard shows real data with termux-api installed
- [ ] Dashboard shows graceful fallback without termux-api
- [ ] Wake-lock toggle works
- [ ] SSH server starts and stops
- [ ] HTTP file server serves files
- [ ] Webhook receives and stores events
- [ ] IoT log displays events with auto-scroll
- [ ] Events persist across restarts
- [ ] `q` quits cleanly with no orphan processes
- [ ] `go test ./...` passes
- [ ] Database auto-creates on first run
- [ ] PID files cleaned up on exit
