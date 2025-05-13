package common

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/logger"
	tea "github.com/charmbracelet/bubbletea"
)

type OSPFMsg string

type RunningConfigMsg string

func RunCustomCommand(activeShell string, command string, timeout time.Duration, logger *logger.Logger) (string, error) {
	logger.WithAttrs(map[string]interface{}{
		"shell":   activeShell,
		"command": command,
		"timeout": timeout.String(),
	}).Info("Executing command")

	var cmd *exec.Cmd

	if activeShell == "vtysh" {
		cmd = exec.Command("vtysh", "-c", command)
	} else if activeShell == "bash" {
		args := strings.Fields(command)
		if len(args) == 0 {
			err := fmt.Errorf("no command provided")
			logger.WithAttrs(map[string]interface{}{
				"error": err.Error(),
			}).Error("Empty command provided")
			return "", err
		}
		cmd = exec.Command(args[0], args[1:]...)
	} else {
		args := strings.Fields(command)
		if len(args) == 0 {
			err := fmt.Errorf("no command provided")
			logger.WithAttrs(map[string]interface{}{
				"error": err.Error(),
			}).Error("Empty command provided")
			return "", err
		}
		cmd = exec.Command(args[0], args[1:]...)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	logger.Debug("Starting command execution")

	if err := cmd.Start(); err != nil {
		logger.WithAttrs(map[string]interface{}{
			"error": err.Error(),
		}).Error("Failed to start command")
		return "", err
	}

	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		killErr := cmd.Process.Kill()
		if killErr != nil {
			logger.WithAttrs(map[string]interface{}{
				"error":         killErr.Error(),
				"output_so_far": out.String(),
			}).Error("Failed to kill timed-out command")
			return "", killErr
		}

		timeoutErr := fmt.Errorf("command timed out")
		logger.WithAttrs(map[string]interface{}{
			"error":         timeoutErr.Error(),
			"output_so_far": out.String(),
		}).Warning("Command timed out")
		return "", timeoutErr

	case err := <-done:
		output := out.String()

		outputForLog := output
		if len(outputForLog) > 500 { // Limit output size in logs
			outputForLog = outputForLog[:500] + "... [truncated]"
		}

		if err != nil {
			logger.WithAttrs(map[string]interface{}{
				"error":     err.Error(),
				"output":    outputForLog,
				"exit_code": getExitCode(err),
			}).Error("Command execution failed")
			return "", fmt.Errorf("command error: %v\nOutput: %s", err, output)
		}

		logger.WithAttrs(map[string]interface{}{
			"output":        outputForLog,
			"output_length": len(output),
		}).Info("Command executed successfully")
		return output, nil
	}
}

// Helper function to extract exit code from error
func getExitCode(err error) int {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode()
	}
	return -1
}

// FetchRunningConfig msg.type = msg.RunningConfigMsg --> triggers function in update
func FetchRunningConfig(logger *logger.Logger) tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		data, err := GetRunningConfig(logger)
		if err != nil {
			return RunningConfigMsg(fmt.Sprintf("Error: %v", err))
		}
		return RunningConfigMsg(data)
	})
}

func GetRunningConfig(logger *logger.Logger) (string, error) {
	vtyshOutput, err := exec.Command("vtysh", "-c", "show running-config").Output()
	if err != nil {
		logger.Error(fmt.Sprintf("Error fetching OSPF neighbor data: %v", err))
		return "", fmt.Errorf("error fetching OSPF neighbor data: %v", err)
	}

	return string(vtyshOutput), nil
}

func ShowRunningConfig(data string) []string {
	if data == "" {
		return []string{"No OSPF data received"}
	}
	return strings.Split(data, "\n")
}

// FetchOSPFData msg.type = msg.OSPFMsg --> triggers function in update
func FetchOSPFData(logger *logger.Logger) tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		data, err := GetOSPFData(logger)
		if err != nil {
			logger.Error(fmt.Sprintf("Error on fetching ospf data: %v", err))
			return OSPFMsg(fmt.Sprintf("Error: %v", err))
		}
		return OSPFMsg(data)
	})
}

func GetOSPFData(logger *logger.Logger) (string, error) {
	vtyshOutput, err := exec.Command("vtysh", "-c", "show ip ospf neighbor").Output()
	if err != nil {
		logger.Error(fmt.Sprintf("Error on getting ospf data: %v", err))
		return "", fmt.Errorf("error fetching OSPF neighbor data: %v", err)
	}

	return fmt.Sprintf("OSPF Neighbors:\n%s", vtyshOutput), nil
}

func DetectOSPFAnomalies(data string) []string {
	if data == "" {
		return []string{"No OSPF data received"}
	}
	return strings.Split(data, "\n")
}
