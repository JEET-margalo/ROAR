//go:build windows

package api

import (
	"fmt"
	"os/exec"
)

// Windows does not expose Unix process groups or SIGSTOP/SIGCONT through the
// standard library. The web app remains fully usable; cancellation stops the
// direct scanner process and pause/resume reports a clear limitation.
func configureScanProcess(cmd *exec.Cmd) {}

func stopScanProcess(cmd *exec.Cmd) {
	if cmd != nil && cmd.Process != nil {
		_ = cmd.Process.Kill()
	}
}

func pauseScanProcess(cmd *exec.Cmd) error {
	return fmt.Errorf("pause is not supported on Windows; cancel the scan instead")
}

func resumeScanProcess(cmd *exec.Cmd) error {
	return fmt.Errorf("resume is not supported on Windows")
}
