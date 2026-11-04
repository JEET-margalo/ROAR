//go:build !windows

package api

import (
	"fmt"
	"os/exec"
	"syscall"
	"time"
)

func configureScanProcess(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
}

func stopScanProcess(cmd *exec.Cmd) {
	if cmd == nil || cmd.Process == nil {
		return
	}
	pid := cmd.Process.Pid
	if pgid, err := syscall.Getpgid(pid); err == nil {
		_ = syscall.Kill(-pgid, syscall.SIGTERM)
		time.AfterFunc(2*time.Second, func() { _ = syscall.Kill(-pgid, syscall.SIGKILL) })
		return
	}
	_ = cmd.Process.Kill()
}

func pauseScanProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return fmt.Errorf("scan has no active process")
	}
	if pgid, err := syscall.Getpgid(cmd.Process.Pid); err == nil {
		return syscall.Kill(-pgid, syscall.SIGSTOP)
	}
	return cmd.Process.Signal(syscall.SIGSTOP)
}

func resumeScanProcess(cmd *exec.Cmd) error {
	if cmd == nil || cmd.Process == nil {
		return fmt.Errorf("scan has no active process")
	}
	if pgid, err := syscall.Getpgid(cmd.Process.Pid); err == nil {
		return syscall.Kill(-pgid, syscall.SIGCONT)
	}
	return cmd.Process.Signal(syscall.SIGCONT)
}
