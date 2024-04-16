//go:build !windows
// +build !windows

package prepare

import (
	"os/exec"
)

func PrepareBackgroundCommand(cmd *exec.Cmd) {
}
