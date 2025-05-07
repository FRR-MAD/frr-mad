package common

import (
	"bytes"
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"strings"
	"time"
)

type OSPFMsg string

type RunningConfigMsg string

func RunCustomCommand(activeShell string, command string, timeout time.Duration) (string, error) {
	var cmd *exec.Cmd

	if activeShell == "vtysh" {
		cmd = exec.Command("vtysh", "-c", command)
	} else if activeShell == "bash" {
		// cmd = exec.Command(command)
		args := strings.Fields(command)
		if len(args) == 0 {
			return "", fmt.Errorf("no command provided")
		}
		cmd = exec.Command(args[0], args[1:]...)
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

// FetchRunningConfig msg.type = msg.RunningConfigMsg --> triggers function in update
func FetchRunningConfig() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		data, err := GetRunningConfig()
		if err != nil {
			return RunningConfigMsg(fmt.Sprintf("Error: %v", err))
		}
		return RunningConfigMsg(data)
	})
}

func GetRunningConfig() (string, error) {
	vtyshOutput, err := exec.Command("vtysh", "-c", "show running-config").Output()
	if err != nil {
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
func FetchOSPFData() tea.Cmd {
	return tea.Tick(2*time.Second, func(time.Time) tea.Msg {
		data, err := GetOSPFData()
		if err != nil {
			return OSPFMsg(fmt.Sprintf("Error: %v", err))
		}
		return OSPFMsg(data)
	})
}

func GetOSPFData() (string, error) {
	vtyshOutput, err := exec.Command("vtysh", "-c", "show ip ospf neighbor").Output()
	if err != nil {
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
