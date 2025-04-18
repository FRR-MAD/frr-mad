package aggregator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"os"
	"strings"

	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"
)

func ParseOSPFRouterLSA(jsonData []byte) (*frrProto.OSPFRouterData, error) {
	var response frrProto.OSPFRouterData
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OSPF router LSA json: %w", err)
	}
	return &response, nil
}

func ParseStaticFRRConfig(path string) (*frrProto.StaticFRRConfiguration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func(file *os.File) {
		if closeErr := file.Close(); closeErr != nil {
			// todo: this must be updated to write it to logger
			err = fmt.Errorf("failed to close config file: %w", closeErr)
		}
	}(file)

	config := &frrProto.StaticFRRConfiguration{}
	scanner := bufio.NewScanner(file)
	var currentInterfacePointer *frrProto.Interface

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "!") {
			continue
		}

		switch {
		case strings.HasPrefix(line, "hostname"):
			parts := strings.Fields(line)
			config.Hostname = parts[1]
		case strings.HasPrefix(line, "frr version"):
			parts := strings.Fields(line)
			config.FrrVersion = parts[2]
		case strings.HasPrefix(line, "no ipv6 forwarding"):
			config.Ipv6Forwarding = false
			// todo: r101:~# cat /proc/sys/net/ipv6/conf/all/forwarding
		case strings.HasPrefix(line, "no ipv4 forwarding"):
			config.Ipv4Forwarding = false
			// todo: r101:~# cat /proc/sys/net/ipv4/conf/all/forwarding
		case strings.HasPrefix(line, "service advanced-vty"):
			config.ServiceAdvancedVty = true

		case strings.HasPrefix(line, "interface "):
			if currentInterfacePointer != nil {
				config.Interfaces = append(config.Interfaces, currentInterfacePointer)
			}
			parts := strings.Fields(line)
			currentInterfacePointer = &frrProto.Interface{Name: parts[1]}
		case currentInterfacePointer != nil && strings.HasPrefix(line, "ip address "):
			parts := strings.Fields(line)
			ip, ipNet, _ := net.ParseCIDR(parts[2])
			prefixLength, _ := ipNet.Mask.Size()
			currentInterfacePointer.IpAddress = append(currentInterfacePointer.IpAddress, &frrProto.IPPrefix{
				IpAddress:    ip.String(),
				PrefixLength: uint32(prefixLength),
			})
		case currentInterfacePointer != nil && strings.HasPrefix(line, "ip ospf area "):
			parts := strings.Fields(line)
			currentInterfacePointer.Area = parts[3]
		case currentInterfacePointer != nil && strings.HasPrefix(line, "ip ospf passive"):
			currentInterfacePointer.Passive = true
		//case currentInterfacePointer != nil && strings.HasPrefix(line, "ip ospf cost "):
		//	parts := strings.Fields(line)
		//	fmt.Sscanf(parts[3], "%d", &currentInterfacePointer.Cost)

		case strings.HasPrefix(line, "ip route "):
			parts := strings.Fields(line)
			ip, ipNet, _ := net.ParseCIDR(parts[2])
			prefixLength, _ := ipNet.Mask.Size()
			nextHopIP := parts[3]
			config.StaticRoutes = append(config.StaticRoutes, &frrProto.StaticRoute{
				IpPrefix: &frrProto.IPPrefix{
					IpAddress:    ip.String(),
					PrefixLength: uint32(prefixLength),
				},
				NextHop: nextHopIP,
			})

		case strings.HasPrefix(line, "router ospf"):
			parseOSPFGlobalConfig(scanner, config)
		}
	}

	if currentInterfacePointer != nil {
		config.Interfaces = append(config.Interfaces, currentInterfacePointer)
	}

	return config, nil
}

func parseOSPFGlobalConfig(scanner *bufio.Scanner, config *frrProto.StaticFRRConfiguration) {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "exit" {
			break
		}

		switch {
		case strings.HasPrefix(line, "ospf router-id "):
			if config.OspfConfig == nil {
				config.OspfConfig = &frrProto.OSPFConfig{}
			}
			parts := strings.Fields(line)
			config.OspfConfig.RouterId = parts[2]

		//case strings.HasPrefix(line, "network "):
		//	parts := strings.Fields(line)
		//	network, area := parts[1], parts[3]
		//	addNetworkToArea(config, network, area)

		case strings.HasPrefix(line, "redistribute "):
			if config.OspfConfig == nil {
				config.OspfConfig = &frrProto.OSPFConfig{}
			}
			parts := strings.Fields(line)
			redistributionType := ""
			redistributionMetric := ""
			redistributionRouteMap := ""
			for i, part := range parts {
				if part == "redistribution" {
					redistributionType = parts[i+1]
				}
				if part == "metric-type" {
					redistributionMetric = parts[i+1]
				}
				if part == "route-map" {
					redistributionRouteMap = parts[i+1]
				}
			}
			config.OspfConfig.Redistribution = append(config.OspfConfig.Redistribution, &frrProto.Redistribution{
				Type:     redistributionType,
				Metric:   redistributionMetric,
				RouteMap: redistributionRouteMap,
			})
		case strings.HasPrefix(line, "area "):
			if config.OspfConfig == nil {
				config.OspfConfig = &frrProto.OSPFConfig{}
			}
			parts := strings.Fields(line)
			config.OspfConfig.Area = append(config.OspfConfig.Area, &frrProto.Area{Name: parts[1]})
			for i, part := range parts {
				if part == "virtual-link" {
				}
			}
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
