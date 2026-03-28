package daemon

import "os/exec"

// EnableWakeLock acquires the Termux wake lock to prevent Android from
// killing background processes.
func EnableWakeLock() error {
	return exec.Command("termux-wake-lock").Run()
}

// DisableWakeLock releases the Termux wake lock.
func DisableWakeLock() error {
	return exec.Command("termux-wake-unlock").Run()
}
