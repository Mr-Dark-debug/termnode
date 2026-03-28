package hardware

import (
	"encoding/json"
	"os/exec"
)

// PollBattery runs termux-battery-status and parses the JSON output.
func PollBattery() (BatteryInfo, error) {
	out, err := exec.Command("termux-battery-status").Output()
	if err != nil {
		return BatteryInfo{}, err
	}
	var info BatteryInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return BatteryInfo{}, err
	}
	return info, nil
}
