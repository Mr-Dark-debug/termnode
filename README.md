# TermNode

A unified, high-performance Terminal User Interface (TUI) for headless Android management in Termux.

---

## What is TermNode?

TermNode is a single compiled Go binary that turns your old Android phone into a fully manageable headless server. Instead of juggling raw SSH commands, messy bash scripts, and scattered `termux-api` calls, TermNode gives you a beautiful terminal dashboard — right in your Termux session.

It was designed for the growing community of users repurposing Android devices as home servers, media hubs, network nodes, and IoT gateways.

### Why Go?

Go compiles to a **single, statically linked binary** for `aarch64`. Termux users can `wget` a release from GitHub and run it instantly — no Python environments, no dependency nightmares, no bloated runtimes. Just one file.

---

## Features

### Hardware Dashboard (`F1`)
Real-time monitoring of your device's critical hardware metrics, powered by `termux-api`:
- Battery percentage, status, health, and temperature
- CPU usage, core count, architecture, and temperature
- WiFi SSID, local IP address, and BSSID
- Auto-refreshes every 5 seconds
- Visual battery gauge with color-coded status

### Service Manager (`F2`)
One-toggle controls for critical background services:
- **Wake Lock** — Prevents Android's battery optimization from killing your processes
- **SSH Server** — Instantly start/stop the OpenSSH daemon for remote access
- **HTTP File Server** — Spin up a Go-powered file server on port 8081 to serve files from your home directory

No more memorizing port flags or digging through documentation.

### IoT Bridge (`F3`)
Turn your Android phone into a data collection hub for microcontrollers:
- Built-in HTTP webhook receiver on port 8080
- POST sensor data from any ESP32, Arduino, or IoT device:
  ```bash
  curl -X POST http://<phone-ip>:8080/webhook/temperature -d '{"value":23.5}'
  ```
- All events are persisted to a local **SQLite** database
- Scrollable, real-time event log in the terminal
- Health check endpoint at `GET /health`

### Help Screen (`F4`)
Full keybinding reference and feature documentation, always one keypress away.

---

## Screenshots

```
┌─────────────────────────────────────────────────┐
│ 1:Dashboard  2:Services  3:IoT Log  4:Help      │
│─────────────────────────────────────────────────│
│ ┌──────────────────┐ ┌──────────────────┐       │
│ │ Battery          │ │ System           │       │
│ │                  │ │                  │       │
│ │ Level:    85%    │ │ WiFi SSID: MyNet │       │
│ │ Status:   CHARG… │ │ Local IP:  192.… │       │
│ │ Health:   GOOD   │ │ BSSID:     AA:B… │       │
│ │ Temp:     32.1°C │ │                  │       │
│ │ Current:  450 mA │ │ CPU              │       │
│ │                  │ │                  │       │
│ │ [████████░░░░░]  │ │ Usage:  12.3%   │       │
│ └──────────────────┘ │ Cores:  8        │       │
│                      │ Temp:   41.2°C   │       │
│                      │ Arch:   aarch64  │       │
│                      └──────────────────┘       │
└─────────────────────────────────────────────────┘
```

---

## Quick Start

### Prerequisites

```bash
# In Termux
pkg install golang termux-api

# Grant termux-api permissions in Android settings
# Test it works:
termux-battery-status
```

### Install from Source

```bash
git clone https://github.com/Mr-Dark-debug/termnode.git
cd termnode
go mod tidy
go build -o termnode .
./termnode
```

### One-Line Install (Binary)

```bash
# Download the latest release (when available)
wget https://github.com/Mr-Dark-debug/termnode/releases/latest/download/termnode-arm64 -O termnode
chmod +x termnode
./termnode
```

### Command-Line Options

```
Usage: termnode [options]

Options:
  -db string     Path to SQLite database (default "$HOME/.termnode/termnode.db")
  -port string   HTTP webhook listen address (default ":8080")
  -debug         Enable debug logging to debug.log
  -version       Print version and exit
```

---

## Architecture

### Tech Stack

| Component | Technology | Why |
|-----------|-----------|-----|
| TUI Framework | Bubble Tea (Elm Architecture) | Elegant state management, composable |
| Styling | Lip Gloss | Beautiful terminal rendering |
| Components | Bubbles | Viewport, keybindings, etc. |
| Database | SQLite (modernc.org/sqlite) | Pure Go, no CGO, no native deps |
| HTTP Server | net/http (stdlib) | Go 1.22+ method routing built-in |
| Language | Go | Single binary, fast, cross-compile |

**All dependencies are pure Go — zero CGO required.** This is critical for Termux where the CGO toolchain may not be available.

### Project Structure

```
termnode/
├── main.go                         # Thin entry point
├── internal/
│   ├── app/                        # Root TUI model & routing
│   ├── daemon/                     # Wake-lock & process management
│   ├── hardware/                   # termux-api polling & parsing
│   ├── iot/                        # HTTP webhook server
│   ├── db/                         # SQLite repository layer
│   ├── screen/                     # Individual TUI screens
│   └── theme/                      # Lipgloss styles & colors
├── migrations/                     # SQL schema (embedded in binary)
└── plan/                           # Development plans & docs
```

### Data Flow

```
┌─────────────┐    tea.Tick(5s)    ┌──────────────┐
│  termux-api  │ ────────────────▶ │  Dashboard   │
│  /proc/stat  │   UpdateMsg       │  Screen      │
└─────────────┘                    └──────────────┘

┌─────────────┐    Enter key       ┌──────────────┐
│  daemon.     │ ────────────────▶ │  Services    │
│  Manager     │   Toggle()        │  Screen      │
└─────────────┘                    └──────────────┘

┌─────────────┐    HTTP POST       ┌───────┐    chan    ┌──────────────┐
│  ESP32 /    │ ────────────────▶ │  IoT  │ ────────▶ │  IoT Log     │
│  curl /     │   /webhook/{topic}│ Bridge│  Event    │  Screen      │
│  Arduino    │                    └───────┘           └──────────────┘
                                         │
                                         ▼
                                  ┌──────────────┐
                                  │   SQLite     │
                                  │   Database   │
                                  └──────────────┘
```

---

## IoT Webhook API

### POST /webhook/{topic}

Send sensor data to TermNode from any HTTP-capable device.

**Request:**
```bash
curl -X POST http://192.168.1.100:8080/webhook/temperature \
  -H "Content-Type: application/json" \
  -d '{"sensor":"DHT22","value":23.5,"unit":"celsius"}'
```

**Response:**
```json
{"status":"ok"}
```

**ESP32 Arduino Example:**
```cpp
HTTPClient http;
http.begin("http://192.168.1.100:8080/webhook/temperature");
http.POST("{\"sensor\":\"DHT22\",\"value\":23.5}");
http.end();
```

### GET /health

Health check endpoint for monitoring.

**Response:**
```json
{"status":"ok"}
```

---

## Keybindings

| Key | Action |
|-----|--------|
| `F1` or `1` | Dashboard tab |
| `F2` or `2` | Services tab |
| `F3` or `3` | IoT Log tab |
| `F4` or `?` | Help tab |
| `q` or `Ctrl+C` | Quit |
| `↑` / `k` | Move up |
| `↓` / `j` | Move down |
| `Enter` | Toggle / Activate |
| `r` | Manual refresh (Dashboard) |

---

## Building

```bash
# Default build
go build -o termnode .

# With version info
go build -ldflags="-X main.version=$(git describe --tags)" -o termnode .

# With MQTT support (future)
go build -tags mqtt -o termnode .

# Cross-compile from desktop
GOOS=android GOARCH=arm64 go build -o termnode-android .
```

### Using Make

```bash
make           # Build
make run       # Build and run
make test      # Run tests
make debug     # Build and run with debug logging
make clean     # Remove binary and logs
make deps      # Download and tidy dependencies
```

---

## Requirements

- **Termux** (from F-Droid or GitHub releases)
- **Go 1.22+** (`pkg install golang`)
- **termux-api** (`pkg install termux-api`) — for hardware dashboard
- **termux-api app** — install the Termux:API Android app for API access
- **OpenSSH** (`pkg install openssh`) — optional, for SSH server feature

---

## Roadmap

- [ ] MQTT subscriber (behind build tag)
- [ ] Tailscale/Ngrok tunnel management
- [ ] Process supervision with auto-restart
- [ ] Configuration file support
- [ ] GitHub Actions CI/CD for automated releases
- [ ] AUR/brew package
- [ ] Real-time charts in the terminal (asciigraph)

---

## Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to the branch (`git push origin feature/amazing`)
5. Open a Pull Request

---

## License

MIT License — see [LICENSE](LICENSE) for details.

---

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — The elegant TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — Terminal styling
- [modernc.org/sqlite](https://pkg.go.dev/modernc.org/sqlite) — Pure Go SQLite
- [Termux](https://termux.dev/) — The Android terminal emulator that makes this possible

---

Built with love for the r/termux community.
