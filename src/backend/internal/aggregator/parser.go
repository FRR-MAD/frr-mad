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

func ParseOSPFRouterLSA(jsonData []byte) (*OSPFRouterData, error) {
	var response OSPFRouterData
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OSPF router LSA json: %w", err)
	}
	return &response, nil
}

func ParseOSPFNetworkLSA(jsonData []byte) (*OSPFNetworkData, error) {
	var response OSPFNetworkData
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OSPF network LSA json: %w", err)
	}
	return &response, nil
}

func ParseOSPFSummaryLSA(jsonData []byte) (*OSPFSummaryData, error) {
	var response OSPFSummaryData
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OSPF summary LSA json: %w", err)
	}
	return &response, nil
}

func ParseOSPFAsbrSummaryLSA(jsonData []byte) (*OSPFAsbrSummaryData, error) {
	var response OSPFAsbrSummaryData
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OSPF asbr summary LSA json: %w", err)
	}
	return &response, nil
}

func ParseOSPFExternalLSA(jsonData []byte) (*OSPFExternalData, error) {
	var response OSPFExternalData
	if err := json.Unmarshal(jsonData, &response); err != nil {
		return nil, fmt.Errorf("failed to parse OSPF external LSA json: %w", err)
	}
	return &response, nil
}

func ParseOSPFNssaExternalLSA(jsonData []byte) (*OSPFNssaExternalData, error) {
	var response OSPFNssaExternalData
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

		if handled := parseMetadataLine(config, line); handled {
			continue
		}

		if newInterface := parseInterfaceLine(line); newInterface != nil {
			// start of a new block: flush previous
			if currentInterfacePointer != nil {
				config.Interfaces = append(config.Interfaces, currentInterfacePointer)
			}
			currentInterfacePointer = newInterface
			continue
		}
		if currentInterfacePointer != nil && parseInterfaceSubLine(currentInterfacePointer, line) {
			continue
		}

		if handled := parseStaticRouteLine(config, line); handled {
			continue
		}

		if strings.HasPrefix(line, "router ospf") {
			parseRouterOSPFConfig(scanner, config)
			continue
		}

		if handled := parseAccessListLine(config, line); handled {
			continue
		}

		if handled := parseRouteMapLine(config, line); handled {
			continue
		}

		if handled := parseRouteMapMatchLine(config, line); handled {
			continue
		}
	}

	// flush last Interface block
	if currentInterfacePointer != nil {
		config.Interfaces = append(config.Interfaces, currentInterfacePointer)
	}

	return config, nil
}

func parseMetadataLine(config *frrProto.StaticFRRConfiguration, line string) bool {
	parts := strings.Fields(line)
	switch {
	case strings.HasPrefix(line, "hostname "):
		config.Hostname = parts[1]
		return true
	case strings.HasPrefix(line, "frr version "):
		config.FrrVersion = parts[2]
		return true
	case line == "no ipv6 forwarding": // todo: add to tests
		config.Ipv6Forwarding = false // todo: check if default value is true on router
		return true
	case line == "no ipv4 forwarding": // todo: add to tests
		config.Ipv4Forwarding = false
		return true
	case line == "service advanced-vty":
		config.ServiceAdvancedVty = true
		return true
	}
	return false
}

func parseInterfaceLine(line string) *frrProto.Interface {
	if !strings.HasPrefix(line, "interface ") {
		return nil
	}
	parts := strings.Fields(line)
	return &frrProto.Interface{Name: parts[1]}
}

func parseInterfaceSubLine(currentInterfacePointer *frrProto.Interface, line string) bool {
	switch {
	case strings.HasPrefix(line, "ip address "):
		parts := strings.Fields(line)
		ip, ipNet, err := net.ParseCIDR(parts[2])
		if err != nil || ipNet == nil {
			log.Printf("bad CIDR %q: %v", parts[2], err)
			return true
		}
		prefixLength, _ := ipNet.Mask.Size()
		currentInterfacePointer.IpAddress = append(currentInterfacePointer.IpAddress, &frrProto.IPPrefix{
			IpAddress:    ip.String(),
			PrefixLength: uint32(prefixLength),
		})
		return true
	case strings.HasPrefix(line, "ip ospf area "):
		currentInterfacePointer.Area = strings.Fields(line)[3]
		return true
	case line == "ip ospf passive":
		currentInterfacePointer.Passive = true
		return true
	case line == "exit":
		return true
	}
	return false
}

func parseStaticRouteLine(config *frrProto.StaticFRRConfiguration, line string) bool {
	if !strings.HasPrefix(line, "ip route ") {
		return false
	}
	parts := strings.Fields(line)
	ip, ipNet, err := net.ParseCIDR(parts[2])
	if err != nil || ipNet == nil {
		log.Printf("bad static route CIDR %q", parts[2])
		return true
	}
	prefixLength, _ := ipNet.Mask.Size()
	config.StaticRoutes = append(config.StaticRoutes, &frrProto.StaticRoute{
		IpPrefix: &frrProto.IPPrefix{
			IpAddress:    ip.String(),
			PrefixLength: uint32(prefixLength),
		},
		NextHop: parts[3],
	})
	return true
}

func parseAccessListLine(config *frrProto.StaticFRRConfiguration, line string) bool {
	if !strings.HasPrefix(line, "access-list ") {
		return false
	}
	parts := strings.Fields(line)
	if len(parts) < 6 {
		log.Printf("short access-list line: %q", line)
		return true
	}
	name, sequence, action, target := parts[1], parts[3], parts[4], parts[5]
	seq, _ := strconv.Atoi(sequence)

	item := &frrProto.AccessListItem{
		Sequence:      uint32(seq),
		AccessControl: action,
	}
	if target == "any" {
		item.Destination = &frrProto.AccessListItem_Any{Any: true}
	} else if ip, ipnet, err := net.ParseCIDR(target); err == nil && ipnet != nil {
		prefixLength, _ := ipnet.Mask.Size()
		item.Destination = &frrProto.AccessListItem_IpPrefix{
			IpPrefix: &frrProto.IPPrefix{IpAddress: ip.String(),
				PrefixLength: uint32(prefixLength),
			},
		}
	} else {
		log.Printf("bad CIDR %q in ACL %q", target, line)
	}

	if config.AccessList == nil {
		config.AccessList = make(map[string]*frrProto.AccessList)
	}
	if _, ok := config.AccessList[name]; !ok {
		config.AccessList[name] = &frrProto.AccessList{}
	}
	config.AccessList[name].AccessListItems = append(config.AccessList[name].AccessListItems, item)
	return true
}

func parseRouteMapLine(config *frrProto.StaticFRRConfiguration, line string) bool {
	if !strings.HasPrefix(line, "route-map ") {
		return false
	}
	parts := strings.Fields(line)
	if len(parts) < 4 {
		log.Printf("short route-map line: %q", line)
		return true
	}
	name, action, sequence := parts[1], parts[2], parts[3]
	if config.RouteMap == nil {
		config.RouteMap = make(map[string]*frrProto.RouteMap)
	}
	config.RouteMap[name] = &frrProto.RouteMap{
		Permit:   action == "permit",
		Sequence: sequence,
	}
	return true
}

func parseRouteMapMatchLine(config *frrProto.StaticFRRConfiguration, line string) bool {
	if !strings.HasPrefix(line, "match ip address ") {
		return false
	}
	parts := strings.Fields(line)
	accessListName := parts[3]
	for _, rm := range config.RouteMap {
		if rm.AccessList == "" {
			rm.Match = "ip address"
			rm.AccessList = accessListName
			break
		}
	}
	return true
}

func parseRouterOSPFConfig(scanner *bufio.Scanner, config *frrProto.StaticFRRConfiguration) {
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
			area.Type = parts[2]
			for i, part := range parts {
				if part == "virtual-link" && i+1 < len(parts) {
					area.Type = "transit"
					config.OspfConfig.VirtualLinkNeighbor = parts[i+1]
					break
				}
			}
			config.OspfConfig.Area = append(config.OspfConfig.Area, area)
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
