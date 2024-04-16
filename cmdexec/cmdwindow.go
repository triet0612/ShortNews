//go:build windows
// +build windows

package prepare

import (
	"os/exec"
	"syscall"
)

func PrepareBackgroundCommand(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
}
