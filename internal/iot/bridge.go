package iot

import (
	"io"
	"net/http"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Mr-Dark-debug/termnode/internal/db"
)

// Bridge runs an HTTP webhook server that receives IoT events and forwards
// them to the TUI via a channel.
type Bridge struct {
	server *http.Server
	repo   *db.Repository
	events chan db.IoTEvent
	port   string
}

// NewBridge creates a new IoT bridge.
func NewBridge(repo *db.Repository, port string, events chan db.IoTEvent) *Bridge {
	return &Bridge{
		repo:   repo,
		events: events,
		port:   port,
	}
}

// Start launches the HTTP webhook server in a goroutine.
func (b *Bridge) Start() error {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /webhook/{topic}", b.handleWebhook)
	mux.HandleFunc("GET /health", b.handleHealth)

	b.server = &http.Server{
		Addr:    b.port,
		Handler: mux,
	}

	go func() {
		if err := b.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			// Server failed to start; send error on channel
			b.events <- db.IoTEvent{
				Topic:   "_system",
				Payload: "HTTP server error: " + err.Error(),
				Source:  "system",
			}
		}
	}()

	return nil
}

// Stop gracefully shuts down the HTTP server.
func (b *Bridge) Stop() error {
	if b.server != nil {
		return b.server.Close()
	}
	return nil
}

// ListenCmd returns a tea.Cmd that blocks on the events channel,
// bridging IoT events from goroutines into the TUI update loop.
func (b *Bridge) ListenCmd() tea.Cmd {
	return func() tea.Msg {
		event, ok := <-b.events
		if !ok {
			return nil
		}
		return IoTEventMsg{Event: event}
	}
}

func (b *Bridge) handleWebhook(w http.ResponseWriter, r *http.Request) {
	topic := r.PathValue("topic")
	if topic == "" {
		http.Error(w, "missing topic", http.StatusBadRequest)
		return
	}

	body, err := io.ReadAll(io.LimitReader(r.Body, 1024*1024)) // 1MB limit
	if err != nil {
		http.Error(w, "failed to read body", http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()

	event := db.IoTEvent{
		Topic:     topic,
		Payload:   string(body),
		Source:    "http",
	}

	// Persist to database
	if b.repo != nil {
		if _, err := b.repo.Insert(event); err != nil {
			http.Error(w, "database error", http.StatusInternalServerError)
			return
		}
	}

	// Forward to TUI (non-blocking because channel is buffered)
	select {
	case b.events <- event:
	default:
		// Channel full, drop event to avoid blocking HTTP handler
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{"status":"ok"}`))
}

func (b *Bridge) handleHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

// IoTEventMsg is the tea.Msg sent when a new IoT event is received.
type IoTEventMsg struct {
	Event db.IoTEvent
}
