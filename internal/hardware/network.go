package hardware

import (
	"encoding/json"
	"os/exec"
)

// PollNetwork runs termux-wifi-connectioninfo and parses the JSON output.
func PollNetwork() (NetworkInfo, error) {
	out, err := exec.Command("termux-wifi-connectioninfo").Output()
	if err != nil {
		return NetworkInfo{}, err
	}
	var info NetworkInfo
	if err := json.Unmarshal(out, &info); err != nil {
		return NetworkInfo{}, err
	}
	return info, nil
}
