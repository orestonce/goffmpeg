//+build windows

package goffmpeg

import (
	"os/exec"
	"syscall"
)

func setupCmd(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
}
