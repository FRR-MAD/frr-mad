package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/eiannone/keyboard"
	"github.com/pterm/pterm"
)

var (
	activeTab    = 0 // 0: OSPF, 1: BGP, 2: Custom Commands
	inMonitoring = false
)

func main() {
	// Initialize keyboard input
	if err := keyboard.Open(); err != nil {
		panic(err)
	}
	defer keyboard.Close()

	// Start the main loop
	for {
		if !inMonitoring {
			showTabSelection()
		} else {
			switch activeTab {
			case 0:
				ospfMonitoring()
			case 1:
				bgpMonitoring()
			case 2:
				customCommands()
			}
		}
	}
}

func showTabSelection() {
	// Create a tabbed menu
	options := []string{"OSPF Monitoring", "BGP Monitoring", "Custom Commands", "Exit"}
	selectedOption, _ := pterm.DefaultInteractiveSelect.WithOptions(options).Show()

	switch selectedOption {
	case "OSPF Monitoring":
		activeTab = 0
		inMonitoring = true
	case "BGP Monitoring":
		activeTab = 1
		inMonitoring = true
	case "Custom Commands":
		activeTab = 2
		inMonitoring = true
	case "Exit":
		os.Exit(0)
	}
}

func ospfMonitoring() {
	pterm.Info.Println("OSPF Monitoring (Prototype)")
	pterm.Println("Press 'q+Enter' to go back to the tab selection.")

	// Start a goroutine to listen for key presses
	quitChan := make(chan struct{})
	go func() {
		for {
			command, _ := pterm.DefaultInteractiveTextInput.Show("INPUT")
			if command == "q" {
				close(quitChan)
				return
			}
		}
	}()

	for {
		select {
		case <-quitChan:
			inMonitoring = false
			return
		default:
			// Fetch OSPF data
			data := getOSPFData()
			pterm.Println(data)

			time.Sleep(2 * time.Second)
		}
	}
}

func bgpMonitoring() {
	pterm.Info.Println("BGP Monitoring")
	pterm.Println("Press 'q+Enter' to go back to the tab selection.")

	// Channel to signal quitting
	quitChan := make(chan struct{})

	// Goroutine to handle user input
	go func() {
		for {
			command, _ := pterm.DefaultInteractiveTextInput.Show("INPUT")
			if command == "q" {
				close(quitChan)
				return
			}
		}
	}()

	// Main loop to display BGP data
	for {
		select {
		case <-quitChan:
			inMonitoring = false
			return
		default:
			// Fetch BGP data
			data := getBGPData()
			pterm.Println(data)

			time.Sleep(2 * time.Second)
		}
	}
}

func customCommands() {
	pterm.Info.Println("Custom Commands")
	pterm.Println("Press 'q' to go back to the tab selection.")

	for {
		// Allow user to enter custom commands
		command, _ := pterm.DefaultInteractiveTextInput.Show()
		if command == "q" {
			inMonitoring = false
			return
		}

		// Execute the command with a timeout
		output, err := runCommandWithTimeout(command, 5*time.Second)
		if err != nil {
			pterm.Error.Println("Failed to run command:", err)
			continue
		}

		// Display the output
		pterm.Println(output)
	}
}

func getOSPFData() string {
	// Fetch OSPF data (e.g., using vtysh)
	cmd := exec.Command("vtysh", "-c", "show ip ospf neighbor")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return string(output)
}

func getBGPData() string {
	// Fetch BGP data using gobgp
	cmd := exec.Command("gobgp", "global", "rib", "-a", "ipv4")
	output, err := cmd.Output()
	if err != nil {
		return fmt.Sprintf("Error: %v", err)
	}
	return string(output)
}

func runCommandWithTimeout(command string, timeout time.Duration) (string, error) {
	args := strings.Fields(command)
	if len(args) == 0 {
		return "", fmt.Errorf("no command provided")
	}

	// Create a command with a timeout
	cmd := exec.Command(args[0], args[1:]...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	// Start the command
	if err := cmd.Start(); err != nil {
		return "", err
	}

	// Wait for the command to finish or timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()

	select {
	case <-time.After(timeout):
		cmd.Process.Kill()
		return "", fmt.Errorf("command timed out")
	case err := <-done:
		if err != nil {
			return "", err
		}
		return out.String(), nil
	}
}
