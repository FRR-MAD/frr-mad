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
	"google.golang.org/protobuf/encoding/protojson"
)

// func ParseOSPFRouterLSA(jsonData []byte) (*frrProto.OSPFRouterData, error) {
// 	var response frrProto.OSPFRouterData
// 	unmarshaler := protojson.UnmarshalOptions{
// 		DiscardUnknown: true,
// 	}

// 	if err := unmarshaler.Unmarshal(jsonData, &response); err != nil {
// 		return nil, fmt.Errorf("failed to parse OSPF router LSA json: %w", err)
// 	}

// 	return &response, nil
// }

// func ParseOSPFNetworkLSA(jsonData []byte) (*frrProto.OSPFNetworkData, error) {
// 	var response frrProto.OSPFNetworkData
// 	unmarshaler := protojson.UnmarshalOptions{DiscardUnknown: true}
// 	if err := unmarshaler.Unmarshal(jsonData, &response); err != nil {
// 		return nil, fmt.Errorf("failed to parse OSPF network LSA json: %w", err)
// 	}
// 	return &response, nil
// }

func ParseOSPFRouterLSA(jsonData []byte) (*frrProto.OSPFRouterData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	// Transform the JSON to match protobuf field names
	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	// Transform Router Link States
	if routerLinkStates, ok := jsonMap["Router Link States"].(map[string]interface{}); ok {
		transformedStates := make(map[string]interface{})
		for areaID, areaData := range routerLinkStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})

			for lsaID, lsaData := range areaDataMap {
				transformedLSAs[lsaID] = transformRouterLSA(lsaData.(map[string]interface{}))
			}

			transformedStates[areaID] = map[string]interface{}{
				"lsa_entries": transformedLSAs,
			}
		}
		transformedMap["router_states"] = transformedStates
	}

	// Convert back to JSON for protobuf unmarshaling
	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFRouterData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFNetworkLSA(jsonData []byte) (*frrProto.OSPFNetworkData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	if netStates, ok := jsonMap["Net Link States"].(map[string]interface{}); ok {
		transformedStates := make(map[string]interface{})
		for areaID, areaData := range netStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})

			for lsaID, lsaData := range areaDataMap {
				transformedLSAs[lsaID] = transformNetworkLSA(lsaData.(map[string]interface{}))
			}

			transformedStates[areaID] = map[string]interface{}{
				"lsa_entries": transformedLSAs,
			}
		}
		transformedMap["net_states"] = transformedStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFNetworkData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFSummaryLSA(jsonData []byte) (*frrProto.OSPFSummaryData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	// Handle Net States
	if netStates, ok := jsonMap["Net Link States"].(map[string]interface{}); ok {
		transformedNetStates := make(map[string]interface{})
		for areaID, areaData := range netStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})

			for lsaID, lsaData := range areaDataMap {
				transformedLSAs[lsaID] = transformNetworkLSA(lsaData.(map[string]interface{}))
			}

			transformedNetStates[areaID] = map[string]interface{}{
				"lsa_entries": transformedLSAs,
			}
		}
		transformedMap["net_states"] = transformedNetStates
	}

	// Handle Summary States
	if summaryStates, ok := jsonMap["Summary Link States"].(map[string]interface{}); ok {
		transformedSummaryStates := make(map[string]interface{})
		for areaID, areaData := range summaryStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})

			for lsaID, lsaData := range areaDataMap {
				transformedLSAs[lsaID] = transformSummaryLSA(lsaData.(map[string]interface{}))
			}

			transformedSummaryStates[areaID] = map[string]interface{}{
				"lsa_entries": transformedLSAs,
			}
		}
		transformedMap["summary_states"] = transformedSummaryStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFSummaryData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFAsbrSummaryLSA(jsonData []byte) (*frrProto.OSPFAsbrSummaryData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	if asbrStates, ok := jsonMap["ASBR-Summary Link States"].(map[string]interface{}); ok {
		transformedStates := make(map[string]interface{})
		for areaID, areaData := range asbrStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})

			for lsaID, lsaData := range areaDataMap {
				transformedLSAs[lsaID] = transformSummaryLSA(lsaData.(map[string]interface{}))
			}

			transformedStates[areaID] = map[string]interface{}{
				"lsa_entries": transformedLSAs,
			}
		}
		transformedMap["asbr_summary_states"] = transformedStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFAsbrSummaryData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFExternalLSA(jsonData []byte) (*frrProto.OSPFExternalData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	if extStates, ok := jsonMap["AS External Link States"].(map[string]interface{}); ok {
		transformedStates := make(map[string]interface{})
		for lsaID, lsaData := range extStates {
			transformedStates[lsaID] = transformExternalLSA(lsaData.(map[string]interface{}))
		}
		transformedMap["as_external_link_states"] = transformedStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFExternalData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFNssaExternalLSA(jsonData []byte) (*frrProto.OSPFNssaExternalData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	if nssaStates, ok := jsonMap["NSSA External Link States"].(map[string]interface{}); ok {
		transformedStates := make(map[string]interface{})
		for areaID, areaData := range nssaStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})

			for lsaID, lsaData := range areaDataMap {
				transformedLSAs[lsaID] = transformNssaExternalLSA(lsaData.(map[string]interface{}))
			}

			transformedStates[areaID] = map[string]interface{}{
				"data": transformedLSAs,
			}
		}
		transformedMap["nssa_external_link_states"] = transformedStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFNssaExternalData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
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

func transformRouterLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	// Map field names from camelCase to snake_case
	fieldMapping := map[string]string{
		"lsaAge":            "lsa_age",
		"options":           "options",
		"lsaFlags":          "lsa_flags",
		"flags":             "flags",
		"asbr":              "asbr",
		"lsaType":           "lsa_type",
		"linkStateId":       "link_state_id",
		"advertisingRouter": "advertising_router",
		"lsaSeqNumber":      "lsa_seq_number",
		"checksum":          "checksum",
		"length":            "length",
		"numOfLinks":        "num_of_links",
		"routerLinks":       "router_links",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := lsaData[origKey]; exists {
			if origKey == "routerLinks" {
				// Handle router links transformation
				routerLinks := value.(map[string]interface{})
				transformedLinks := make(map[string]interface{})

				for linkID, linkData := range routerLinks {
					transformedLinks[linkID] = transformRouterLink(linkData.(map[string]interface{}))
				}

				transformed[newKey] = transformedLinks
			} else {
				transformed[newKey] = value
			}
		}
	}

	return transformed
}

func transformRouterLink(linkData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"linkType":                "link_type",
		"designatedRouterAddress": "designated_router_address",
		"routerInterfaceAddress":  "router_interface_address",
		"networkAddress":          "network_address",
		"networkMask":             "network_mask",
		"numOfTosMetrics":         "num_of_tos_metrics",
		"tos0Metric":              "tos0_metric",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := linkData[origKey]; exists {
			transformed[newKey] = value
		}
	}

	return transformed
}

func transformNetworkLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"lsaAge":            "lsa_age",
		"options":           "options",
		"lsaFlags":          "lsa_flags",
		"lsaType":           "lsa_type",
		"linkStateId":       "link_state_id",
		"advertisingRouter": "advertising_router",
		"lsaSeqNumber":      "lsa_seq_number",
		"checksum":          "checksum",
		"length":            "length",
		"networkMask":       "network_mask",
		"attchedRouters":    "attached_routers",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := lsaData[origKey]; exists {
			if origKey == "attchedRouters" {
				// Handle attached routers transformation
				routers := value.(map[string]interface{})
				transformedRouters := make(map[string]interface{})

				for routerID, routerData := range routers {
					transformedRouters[routerID] = map[string]interface{}{
						"attached_router_id": routerData.(map[string]interface{})["attachedRouterId"],
					}
				}

				transformed[newKey] = transformedRouters
			} else {
				transformed[newKey] = value
			}
		}
	}

	return transformed
}

func transformSummaryLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"lsaAge":            "lsa_age",
		"options":           "options",
		"lsaFlags":          "lsa_flags",
		"lsaType":           "lsa_type",
		"linkStateId":       "link_state_id",
		"advertisingRouter": "advertising_router",
		"lsaSeqNumber":      "lsa_seq_number",
		"checksum":          "checksum",
		"length":            "length",
		"networkMask":       "network_mask",
		"tos0Metric":        "tos0_metric",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := lsaData[origKey]; exists {
			transformed[newKey] = value
		}
	}

	return transformed
}

func transformExternalLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"lsaAge":            "lsa_age",
		"options":           "options",
		"lsaFlags":          "lsa_flags",
		"lsaType":           "lsa_type",
		"linkStateId":       "link_state_id",
		"advertisingRouter": "advertising_router",
		"lsaSeqNumber":      "lsa_seq_number",
		"checksum":          "checksum",
		"length":            "length",
		"networkMask":       "network_mask",
		"metricType":        "metric_type",
		"tos":               "tos",
		"metric":            "metric",
		"forwardAddress":    "forward_address",
		"externalRouteTag":  "external_route_tag",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := lsaData[origKey]; exists {
			transformed[newKey] = value
		}
	}

	return transformed
}

func transformNssaExternalLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"lsaAge":             "lsa_age",
		"options":            "options",
		"lsaFlags":           "lsa_flags",
		"lsaType":            "lsa_type",
		"linkStateId":        "link_state_id",
		"advertisingRouter":  "advertising_router",
		"lsaSeqNumber":       "lsa_seq_number",
		"checksum":           "checksum",
		"length":             "length",
		"networkMask":        "network_mask",
		"metricType":         "metric_type",
		"tos":                "tos",
		"metric":             "metric",
		"nssaForwardAddress": "nssa_forward_address",
		"externalRouteTag":   "external_route_tag",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := lsaData[origKey]; exists {
			transformed[newKey] = value
		}
	}

	return transformed
}
