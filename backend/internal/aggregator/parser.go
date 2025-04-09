package aggregator

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func ParseConfig(path string) (*NetworkConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := &NetworkConfig{}
	scanner := bufio.NewScanner(file)
	var currentIface *OSPFInterfaceConfig

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "!") {
			continue
		}

		switch {
		case strings.HasPrefix(line, "interface "):
			if currentIface != nil {
				config.Interfaces = append(config.Interfaces, *currentIface)
			}
			parts := strings.Fields(line)
			currentIface = &OSPFInterfaceConfig{Name: parts[1]}

		case currentIface != nil && strings.HasPrefix(line, "ip ospf area "):
			parts := strings.Fields(line)
			currentIface.Area = parts[3]

		case currentIface != nil && strings.HasPrefix(line, "ip ospf passive"):
			currentIface.Passive = true

		case currentIface != nil && strings.HasPrefix(line, "ip ospf cost "):
			parts := strings.Fields(line)
			fmt.Sscanf(parts[3], "%d", &currentIface.Cost)

		case strings.HasPrefix(line, "router ospf"):
			parseOSPFGlobalConfig(scanner, config)
		}
	}

	if currentIface != nil {
		config.Interfaces = append(config.Interfaces, *currentIface)
	}

	return config, nil
}

func parseOSPFGlobalConfig(scanner *bufio.Scanner, config *NetworkConfig) {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "exit" {
			break
		}

		switch {
		case strings.HasPrefix(line, "ospf router-id "):
			parts := strings.Fields(line)
			config.RouterID = parts[2]

		case strings.HasPrefix(line, "network "):
			parts := strings.Fields(line)
			network, area := parts[1], parts[3]
			addNetworkToArea(config, network, area)
		}
	}
}

func addNetworkToArea(config *NetworkConfig, network, area string) {
	for i, a := range config.Areas {
		if a.ID == area {
			config.Areas[i].Networks = append(a.Networks, network)
			return
		}
	}
	config.Areas = append(config.Areas, OSPFArea{
		ID:       area,
		Networks: []string{network},
	})
}
