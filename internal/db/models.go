package db

import "time"

// IoTEvent represents a single IoT event received via webhook or MQTT.
type IoTEvent struct {
	ID        int64
	Topic     string
	Payload   string
	Source    string    // "http" or "mqtt"
	Timestamp time.Time
}
