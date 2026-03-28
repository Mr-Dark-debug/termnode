package hardware

import (
	"os"
	"runtime"
	"strconv"
	"strings"
)

// PollCPU reads CPU usage from /proc/stat and temperature from sysfs.
func PollCPU() (CPUStats, error) {
	stats := CPUStats{
		CoreCount: runtime.NumCPU(),
		Arch:      runtime.GOARCH,
	}

	// CPU usage from /proc/stat
	usage, err := readCPUUsage()
	if err == nil {
		stats.UsagePercent = usage
	}

	// Temperature from sysfs (Android thermal zone)
	temp, err := readCPUTemp()
	if err == nil {
		stats.Temperature = temp
	}

	return stats, nil
}

func readCPUUsage() (float64, error) {
	data, err := os.ReadFile("/proc/stat")
	if err != nil {
		return 0, err
	}

	line := strings.Split(string(data), "\n")[0]
	fields := strings.Fields(line)
	if len(fields) < 5 {
		return 0, nil
	}

	var total, idle float64
	for i, f := range fields[1:] {
		v, _ := strconv.ParseFloat(f, 64)
		total += v
		if i == 3 { // idle is the 4th field
			idle = v
		}
	}

	if total == 0 {
		return 0, nil
	}
	return ((total - idle) / total) * 100, nil
}

func readCPUTemp() (float64, error) {
	// Try common Android thermal zone paths
	paths := []string{
		"/sys/class/thermal/thermal_zone0/temp",
		"/sys/class/thermal/thermal_zone1/temp",
	}

	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			continue
		}
		val, err := strconv.ParseFloat(strings.TrimSpace(string(data)), 64)
		if err != nil {
			continue
		}
		// Some zones report in millidegrees
		if val > 1000 {
			return val / 1000, nil
		}
		return val, nil
	}
	return 0, nil
}
