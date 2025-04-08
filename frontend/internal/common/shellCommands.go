package common

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

func RunCommand(activeShell string, command string, timeout time.Duration) (string, error) {
	var cmd *exec.Cmd

	if activeShell == "vtysh" {
		cmd = exec.Command("vtysh", "-c", command)
	} else if activeShell == "bash" {
		cmd = exec.Command(command)
	} else {
		args := strings.Fields(command)
		if len(args) == 0 {
			return "", fmt.Errorf("no command provided")
		}
		cmd = exec.Command(args[0], args[1:]...)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Start(); err != nil {
		return "", err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		err := cmd.Process.Kill()
		if err != nil {
			return "", err
		}
		return "", fmt.Errorf("command timed out")
	case err := <-done:
		if err != nil {
			return "", fmt.Errorf("command error: %v\nOutput: %s", err, out.String())
		}
		return out.String(), nil
	}
}
