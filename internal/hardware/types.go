package hardware

// BatteryInfo holds parsed output from termux-battery-status.
type BatteryInfo struct {
	Percentage  int     `json:"percentage"`
	Status      string  `json:"status"`      // CHARGING, DISCHARGING, FULL, NOT_CHARGING
	Temperature float64 `json:"temperature"` // Celsius
	Health      string  `json:"health"`      // GOOD, OVERHEAT, etc.
	Current     int     `json:"current"`     // mA
	Plugged     string  `json:"plugged"`     // AC, USB, UNKNOWN
}

// NetworkInfo holds parsed output from termux-wifi-connectioninfo.
type NetworkInfo struct {
	IP    string `json:"ip"`
	SSID  string `json:"ssid"`
	BSSID string `json:"bssid"`
	MAC   string `json:"mac"`
}

// CPUStats holds CPU usage and temperature data.
type CPUStats struct {
	UsagePercent float64
	CoreCount    int
	Temperature  float64 // Celsius
	Arch         string  // e.g., "aarch64"
}
