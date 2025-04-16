package aggregator

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"
)

func ParseConfig(path string) (*frrProto.NetworkConfig, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	config := &frrProto.NetworkConfig{}
	scanner := bufio.NewScanner(file)
	var currentIface *frrProto.OSPFInterfaceConfig

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "!") {
			continue
		}

		switch {
		case strings.HasPrefix(line, "interface "):
			if currentIface != nil {
				config.Interfaces = append(config.Interfaces, currentIface)
			}
			parts := strings.Fields(line)
			currentIface = &frrProto.OSPFInterfaceConfig{Name: parts[1]}

		case currentIface != nil && strings.HasPrefix(line, "ip address "):
			parts := strings.Fields(line)
			currentIface.IpAddress = parts[2]

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
		config.Interfaces = append(config.Interfaces, currentIface)
	}

	return config, nil
}

func parseOSPFGlobalConfig(scanner *bufio.Scanner, config *frrProto.NetworkConfig) {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "exit" {
			break
		}

		switch {
		case strings.HasPrefix(line, "ospf router-id "):
			parts := strings.Fields(line)
			config.RouterId = parts[2]

		case strings.HasPrefix(line, "network "):
			parts := strings.Fields(line)
			network, area := parts[1], parts[3]
			addNetworkToArea(config, network, area)
		}
	}
}

func addNetworkToArea(config *frrProto.NetworkConfig, network, area string) {
	for i, a := range config.Areas {
		if a.Id == area {
			config.Areas[i].Networks = append(a.Networks, network)
			return
		}
	}
	config.Areas = append(config.Areas, &frrProto.OSPFArea{
		Id:       area,
		Networks: []string{network},
	})
}

func (c *Collector) ReadConfig() (string, error) {
	file, err := os.Open(c.configPath)
	if err != nil {
		return "", fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var staticConfig []string

	for scanner.Scan() {
		staticConfig = append(staticConfig, scanner.Text())
	}

	return strings.Join(staticConfig, "\n"), nil
}
