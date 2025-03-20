package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/pterm/pterm"
)

var (
	activeTab    = 0 // 0: OSPF, 1: BGP, 2: Custom Commands
	inMonitoring = false
)

func main() {
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
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Println("OSPF Monitoring")
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
			getOSPFData()
			time.Sleep(2 * time.Second)
		}
	}
}

func bgpMonitoring() {
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightGreen)).WithMargin(10).Println("BGP Monitoring")
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
			getBGPData()
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

func getOSPFData() {
	// Fetch OSPF neighbor data using vtysh
	cmd := exec.Command("vtysh", "-c", "show ip ospf neighbor")
	output, err := cmd.Output()
	if err != nil {
		pterm.Error.Println("Error fetching OSPF neighbor data:", err)
		return
	}

	// Query OSPF route metrics from Prometheus
	prometheusURL := "http://mon:9090/api/v1/query"
	routeQuery := "frr_ospf_route_metric"
	routeData, err := queryPrometheus(prometheusURL, routeQuery)
	if err != nil {
		pterm.Error.Println("Error querying OSPF route metrics:", err)
		return
	}

	// Display OSPF data with pterm
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightBlue)).WithMargin(10).Println("OSPF Monitoring")

	// Display OSPF Neighbors
	pterm.Info.Println("OSPF Neighbors:")
	pterm.Println(string(output))

	// Display OSPF Route Metrics
	pterm.Info.Println("OSPF Route Metrics:")
	pterm.Println(routeData)
}

func getBGPData() {
	// Fetch BGP neighbor data using vtysh
	cmd := exec.Command("vtysh", "-c", "show ip bgp summary")
	output, err := cmd.Output()
	if err != nil {
		pterm.Error.Println("Error fetching BGP summary:", err)
		return
	}

	// Query BGP route anomalies from Prometheus
	prometheusURL := "http://mon:9090/api/v1/query"
	routeQuery := "frr_bgp_route_announced"
	routeData, err := queryPrometheus(prometheusURL, routeQuery)
	if err != nil {
		pterm.Error.Println("Error querying BGP route data:", err)
		return
	}

	// Display BGP data with pterm
	pterm.DefaultHeader.WithBackgroundStyle(pterm.NewStyle(pterm.BgLightGreen)).WithMargin(10).Println("BGP Monitoring")

	// Display BGP Summary
	pterm.Info.Println("BGP Summary:")
	pterm.Println(string(output))

	// Display BGP Route Data
	pterm.Info.Println("BGP Route Data:")
	pterm.Println(routeData)
}

// func getHostname() string {
// 	hostname, err := os.Hostname()
// 	if err != nil {
// 		pterm.Error.Println("Failed to get hostname:", err)
// 		return "localhost" // Fallback to localhost
// 	}
// 	return hostname
// }

func queryPrometheus(prometheusURL, query string) (string, error) {
	params := url.Values{}
	params.Add("query", query)

	resp, err := http.Get(prometheusURL + "?" + params.Encode())
	if err != nil {
		return "", fmt.Errorf("error querying Prometheus: %v", err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response body: %v", err)
	}

	return string(body), nil
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
