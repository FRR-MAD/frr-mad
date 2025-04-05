package common

import (
	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"os/exec"
	"strings"
	"time"
)

type OSPFMsg string

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
