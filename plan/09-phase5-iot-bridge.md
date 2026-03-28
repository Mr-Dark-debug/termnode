# Phase 5: IoT Bridge + Log Viewer

## Goal

Implement the HTTP webhook receiver that accepts sensor data from external devices (ESP32, Arduino, etc.), persists it to SQLite, and displays it in a scrollable terminal log viewer.

---

## Architecture Overview

```
External Devices                TermNode
┌──────────┐                  ┌────────────────────────────────┐
│          │   HTTP POST       │                                │
│  ESP32   │─────────────────▶│  iot.Bridge                    │
│          │  /webhook/temp    │    ├── handleWebhook()          │
│          │  body: JSON       │    │   ├── Parse topic + body  │
└──────────┘                   │    │   ├── db.Insert()         │
                               │    │   └── events ch <- event │
┌──────────┐                   │    │                           │
│          │   HTTP POST       │    ├── ListenCmd()             │
│  Arduino │─────────────────▶│    │   └── <- events ch       │
│          │  /webhook/humid   │    │       → IoTEventMsg      │
│          │  body: JSON       │    │                           │
└──────────┘                   │    └── handleHealth()          │
                               │                                │
┌──────────┐                   │    ┌──────────────────────┐    │
│          │   HTTP POST       │    │  screen.IoTLogModel  │    │
│   curl   │─────────────────▶│    │    ├── viewport       │    │
│          │  /webhook/custom  │    │    ├── AddEvent()     │    │
│          │  body: any text   │    │    └── View()         │    │
└──────────┘                   │    └──────────────────────┘    │
                               │                                │
                               │    ┌──────────────────────┐    │
                               │    │  db.Repository       │    │
                               │    │    └── Insert()      │    │
                               │    │    └── Recent(100)   │    │
                               │    └──────────────────────┘    │
                               └────────────────────────────────┘
```

---

## Files

### internal/iot/bridge.go

The central piece connecting HTTP requests to the TUI.

**Bridge struct:**
```go
type Bridge struct {
    server *http.Server
    repo   *db.Repository
    events chan db.IoTEvent   // buffered channel (capacity 64)
    port   string
}
```

**Key Methods:**

| Method | Purpose |
|--------|---------|
| `NewBridge(repo, port, events)` | Constructor |
| `Start()` | Launch HTTP server in goroutine |
| `Stop()` | Graceful shutdown |
| `ListenCmd()` | Returns tea.Cmd that blocks on events channel |

**The Channel Bridge Pattern:**

This is the most important pattern in the IoT bridge:

```go
// HTTP handler writes to channel (non-blocking)
select {
case b.events <- event:
default:  // drop if channel is full (64 events buffered)
}

// tea.Cmd reads from channel (blocks until available)
func (b *Bridge) ListenCmd() tea.Cmd {
    return func() tea.Msg {
        event := <-b.events
        return IoTEventMsg{Event: event}
    }
}

// app.Update() processes the event and re-issues ListenCmd
case iot.IoTEventMsg:
    m.iotLog.AddEvent(msg.Event)
    return m, m.bridge.ListenCmd()  // keep listening!
```

**Why this works:**
- All UI state mutations happen on the Bubble Tea main goroutine
- No mutexes needed for UI state
- HTTP handlers never block (buffered channel with drop semantics)
- Events are processed sequentially, in order

### Webhook Handler

```go
func (b *Bridge) handleWebhook(w http.ResponseWriter, r *http.Request) {
    // 1. Extract topic from URL: r.PathValue("topic")
    // 2. Read body (1MB limit): io.ReadAll(io.LimitReader(r.Body, 1MB))
    // 3. Create IoTEvent{Topic, Payload, Source: "http"}
    // 4. Persist: b.repo.Insert(event)
    // 5. Forward to TUI: b.events <- event
    // 6. Respond: 201 Created {"status":"ok"}
}
```

**Go 1.22+ Routing:**
```go
mux := http.NewServeMux()
mux.HandleFunc("POST /webhook/{topic}", b.handleWebhook)
mux.HandleFunc("GET /health", b.handleHealth)
```

No external router needed! The `{topic}` path parameter is extracted via `r.PathValue("topic")`.

### Health Endpoint

```go
func (b *Bridge) handleHealth(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(200)
    w.Write([]byte(`{"status":"ok"}`))
}
```

### internal/screen/iotlog.go

**IoTLogModel** uses `bubbles/viewport` for scrollable content:

```go
type IoTLogModel struct {
    width    int
    height   int
    viewport viewport.Model
    repo     *db.Repository
    events   []db.IoTEvent
    ready    bool
}
```

**Key Behaviors:**
- `AddEvent()` appends event, refreshes view from DB, auto-scrolls to bottom
- `refreshContent()` queries `repo.Recent(100)` and formats display
- `View()` shows viewport content + status bar with event count
- Scroll with ↑/↓ or j/k

**Event Display Format:**
```
14:23:05 [temperature] {"sensor":"DHT22","value":23.5,"unit":"celsius"}
14:23:06 [humidity]    {"sensor":"DHT22","value":65.2,"unit":"percent"}
14:23:10 [motion]      {"zone":"kitchen","detected":true}
```

---

## API Reference

### POST /webhook/{topic}

Receives sensor data and stores it in SQLite.

**Request:**
```
POST /webhook/temperature HTTP/1.1
Content-Type: application/json

{"sensor":"DHT22","value":23.5,"unit":"celsius"}
```

**Response (201):**
```json
{"status":"ok"}
```

**Error Responses:**
| Status | Condition |
|--------|-----------|
| 400 | Missing topic in URL path |
| 500 | Failed to read request body |
| 500 | Database insert failed |

### GET /health

Health check for monitoring systems.

**Response (200):**
```json
{"status":"ok"}
```

---

## ESP32 Integration Example

```cpp
#include <WiFi.h>
#include <HTTPClient.h>

const char* server = "http://192.168.1.100:8080";

void sendSensorData(String topic, String payload) {
    HTTPClient http;
    http.begin(server + "/webhook/" + topic);
    http.POST(payload);
    http.end();
}

void loop() {
    float temp = dht.readTemperature();
    String json = "{\"sensor\":\"DHT22\",\"value\":" + String(temp) + "}";
    sendSensorData("temperature", json);
    delay(5000);
}
```

---

## Security Considerations

The current implementation has **no authentication** — the webhook endpoint is open to anyone on the local network.

This is acceptable for:
- Home/LAN use behind a router firewall
- Development and testing
- Networks you control

Future improvements:
- API key header validation
- IP whitelist
- TLS support
- Rate limiting

---

## Verification

```bash
# Build and run
go build -o termnode . && ./termnode

# Press F3 to go to IoT Log tab
# Shows: "No events received yet."
# Shows: "POST data to: http://<ip>:8080/webhook/<topic>"

# In another terminal, send test events:
curl -X POST http://localhost:8080/webhook/temperature -d '{"value":23.5}'
curl -X POST http://localhost:8080/webhook/humidity -d '{"value":65.2}'
curl -X POST http://localhost:8080/webhook/motion -d '{"zone":"kitchen"}'

# Check health endpoint:
curl http://localhost:8080/health
# Expected: {"status":"ok"}

# Back in TermNode:
# Events appear in the log, auto-scrolling to latest
# Each event shows: timestamp [topic] payload
# Status bar shows: "Events: 3"

# Verify persistence:
# Quit TermNode (q), restart it
# IoT Log tab shows the 3 previous events (loaded from SQLite)
```
