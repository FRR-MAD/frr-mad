package aggregator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

func ParseOSPFRouterLSA(jsonData []byte) (*frrProto.OSPFRouterData, error) {
	var response frrProto.OSPFRouterData
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OSPF router LSA json: %w", err)
	}
	return &response, nil
}

func ParseOSPFNetworkLSA(jsonData []byte) (*frrProto.OSPFNetworkData, error) {
	var response frrProto.OSPFNetworkData
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OSPF network LSA json: %w", err)
	}
	return &response, nil
}

func ParseOSPFSummaryLSA(jsonData []byte) (*frrProto.OSPFSummaryData, error) {
	var response frrProto.OSPFSummaryData
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OSPF summary LSA json: %w", err)
	}
	return &response, nil
}

func ParseOSPFAsbrSummaryLSA(jsonData []byte) (*frrProto.OSPFAsbrSummaryData, error) {
	var response frrProto.OSPFAsbrSummaryData
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OSPF asbr summary LSA json: %w", err)
	}
	return &response, nil
}

func ParseOSPFExternalLSA(jsonData []byte) (*frrProto.OSPFExternalData, error) {
	var response frrProto.OSPFExternalData
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OSPF external LSA json: %w", err)
	}
	return &response, nil
}

func ParseOSPFNssaExternalLSA(jsonData []byte) (*frrProto.OSPFNssaExternalData, error) {
	var response frrProto.OSPFNssaExternalData
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OSPF nssa external LSA json: %w", err)
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
			ip, ipNet, err := net.ParseCIDR(parts[2])
			if err != nil || ipNet == nil {
				log.Printf("invalid CIDR: %q on line: %q (err: %v)", parts[2], line, err)
				break
			}
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
			if err != nil || ipNet == nil {
				log.Printf("invalid CIDR: %q on line: %q (err: %v)", parts[2], line, err)
				break
			}
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

		case strings.HasPrefix(line, "access-list "):
			if config.AccessList == nil {
				config.AccessList = make(map[string]*frrProto.AccessList)
			}
			parts := strings.Fields(line)
			accessListName := parts[1]
			seq, _ := strconv.Atoi(parts[3])
			action := parts[4]
			ip, ipNet, _ := net.ParseCIDR(parts[5])
			if err != nil || ipNet == nil {
				log.Printf("invalid CIDR: %q on line: %q (err: %v)", parts[2], line, err)
				break
			}
			prefixLength, _ := ipNet.Mask.Size()
			accessListItem := &frrProto.AccessListItem{
				Sequence:      uint32(seq),
				AccessControl: action,
				Destination: &frrProto.AccessListItem_IpPrefix{
					IpPrefix: &frrProto.IPPrefix{
						IpAddress:    ip.String(),
						PrefixLength: uint32(prefixLength),
					},
				},
			}

			if _, exists := config.AccessList[accessListName]; !exists {
				config.AccessList[accessListName] = &frrProto.AccessList{}
			}
			config.AccessList[accessListName].AccessListItems = append(config.AccessList[accessListName].AccessListItems, accessListItem)

		case strings.HasPrefix(line, "route-map "):
			if config.RouteMap == nil {
				config.RouteMap = make(map[string]*frrProto.RouteMap)
			}
			parts := strings.Fields(line)
			if len(parts) < 4 {
				log.Printf("invalid route-map line: %q", line)
				break
			}
			routeMapName := parts[1]
			action := parts[2]
			sequence := parts[3]
			config.RouteMap[routeMapName] = &frrProto.RouteMap{
				Permit:   action == "permit",
				Sequence: sequence,
			}
		case strings.HasPrefix(line, "match ip address "):
			parts := strings.Fields(line)
			if len(parts) < 4 {
				log.Printf("invalid match line in route-map: %q", line)
				break
			}
			accessListName := parts[3]
			if len(config.RouteMap) > 0 {
				for _, rm := range config.RouteMap {
					if rm.AccessList == "" {
						rm.Match = "ip address"
						rm.AccessList = accessListName
						break
					}
				}
			}

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
				if part == "redistribute" {
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
			area := &frrProto.Area{Name: parts[1]}
			config.OspfConfig.Area = append(config.OspfConfig.Area, area)
			for i, part := range parts {
				if part == "virtual-link" && i+1 < len(parts) {
					area.Type = "transit"
					config.OspfConfig.VirtualLinkNeighbor = parts[i+1]
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
