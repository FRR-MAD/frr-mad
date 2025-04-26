package analyzer

import (
	"fmt"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

// analyze the different ospf anomalies
// call ospf functions

/*
 * intraAreaLsa handles LSA type 1, 2, 3, and 4
 * this difference is handled because there are no area distinctions with these types
 */
type intraAreaLsa struct {
	Hostname string `json:"hostname"`
	RouterId string `json:"router_id"`
	Areas    []area `json:"areas"`
}

/*
 * interAreaLsa handles LSA type 5 and 7
 * this difference is handled because there are no area distinctions with these types
 */
type interAreaLsa struct {
	Hostname string `json:"hostname"`
	RouterId string `json:"router_id"`
	Areas    []area `json:"areas"`
}

type area struct {
	AreaName string         `json:"name"`
	LsaType  string         `json:"lsa_type"`
	Links    []advertisment `json:"links"`
}

type advertisment struct {
	// populated with type 1 lsa
	InterfaceAddress string `json:"interface_address,omitempty"`
	// populated with type 2, 3, 4, 5, 7
	LinkStateId  string `json:"link_state_id,omitempty"`
	PrefixLength string `json:"prefix_length,omitempty"`
	// Cost         int    `json:"cost,omitempty"`
	LinkType string `json:"link_type"`
}

type ACLEntry struct {
	IPAddress    string `json:"ip_address,omitempty"`
	PrefixLength int    `json:"prefix_length,omitempty"`
	IsPermit     bool   `json:"is_permit"`
	Any          bool   `json:"any,omitempty"`
	Sequence     int    `json:"sequence"`
}

type accessList struct {
	accessListName string     `json:"access_list_name"`
	aclEntry       []ACLEntry `json:"acl_entries"`
}

type RedistributedRoute struct {
	IPPrefix     string `json:"ip_prefix"`
	PrefixLength int    `json:"prefix_length,omitempty"`
	NextHop      string `json:"next_hop,omitempty"`
	RouteMapName string `json:"route_map_name,omitempty"`
	Metric       string `json:"metric,omitempty"`
	MetricType   string `json:"metric_type,omitempty"`
}

type RedistributionList struct {
	StaticRoutes []RedistributedRoute `json:"static_routes,omitempty"`
	BGPRoutes    []RedistributedRoute `json:"bgp_routes,omitempty"`
}

func (c *Analyzer) AnomalyAnalysis() {

	// fmt.Println("#################### File Configuration Access List Enhanced ####################") accessListEnhanced := getStaticRedistributionList(c.metrics.StaticFrrConfiguration)
	// fmt.Printf("\n%+v\n", accessListEnhanced)
	// fmt.Println()
	// fmt.Println()
	//
	// access list
	// fmt.Println("#################### File Configuration Access List ####################")
	accessList := getAccessLists(c.metrics.StaticFrrConfiguration)
	// fmt.Printf("\n%+v\n", accessList)
	// fmt.Println()
	// fmt.Println()
	//
	// static file parsing
	// fmt.Println("#################### File Configuration Router LSDB Prediction ####################")
	predictedRouterLSDB := convertStaticFileRouterData(c.metrics.StaticFrrConfiguration)
	// if predictedRouterLSDB != nil {
	// for _, area := range predictedRouterLSDB.Areas {
	// fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
	// }
	// }
	// fmt.Printf("\n%+v\n", predictedRouterLSDB)
	// fmt.Println()
	// fmt.Println()
	//
	// fmt.Println("#################### File Configuration External LSDB Prediction ####################")
	predictedExternalLSDB := convertStaticFileExternalData(c.metrics.StaticFrrConfiguration)
	// if predictedExternalLSDB != nil {
	// for _, area := range predictedExternalLSDB.Areas {
	// fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
	// }
	// }
	// fmt.Printf("\n%+v\n", predictedExternalLSDB)
	// fmt.Println()
	// fmt.Println()
	//
	// fmt.Println("#################### File Configuration NSSA External LSDB Prediction ####################")
	predictedNssaExternalLSDB := convertStaticFileNssaExternalData(c.metrics.StaticFrrConfiguration)
	// if predictedNssaExternalLSDB != nil {
	//
	// for _, area := range predictedNssaExternalLSDB.Areas {
	// fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
	// }
	// }
	// fmt.Printf("\n%+v\n", predictedNssaExternalLSDB)
	// fmt.Println()
	// fmt.Println()
	//
	// runtime parsing
	// fmt.Println("#################### Runtime Configuration Router LSDB IS_STATE ####################")
	runtimeRouterLSDB := convertRuntimeRouterData(c.metrics.OspfRouterData, c.metrics.StaticFrrConfiguration.Hostname)
	// for _, area := range runtimeRouterLSDB.Areas {
	// fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
	// }
	// fmt.Printf("\n%+v\n", runtimeRouterLSDB)
	// fmt.Println()
	// fmt.Println()
	//
	// fmt.Println("#################### Runtime Configuration External LSDB IS_STATE ####################")
	runtimeExternalLSDB := convertRuntimeExternalRouterData(c.metrics.OspfExternalData, c.metrics.StaticFrrConfiguration.Hostname)
	// for _, area := range runtimeExternalLSDB.Areas {
	// fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
	// }
	// fmt.Printf("\n%+v\n", runtimeExternalLSDB)
	// fmt.Println()
	// fmt.Println()
	//
	// fmt.Println("#################### Runtime Configuration NSSA External LSDB IS_STATE ####################")
	runtimeNssaExternalLSDB := convertNssaExternalRouterData(c.metrics.OspfNssaExternalData, c.metrics.StaticFrrConfiguration.Hostname)
	// for _, area := range runtimeNssaExternalLSDB.Areas {
	// fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
	// }
	// fmt.Printf("\n%+v\n", runtimeNssaExternalLSDB)
	// fmt.Println()
	// fmt.Println()

	routerAnomalyAnalysis(accessList, predictedRouterLSDB, runtimeRouterLSDB)

	externalAnomalyAnalysis(accessList, predictedExternalLSDB, runtimeExternalLSDB)

	nssaExternalAnomalyAnalysis(accessList, predictedNssaExternalLSDB, runtimeNssaExternalLSDB)

}

//func parseIPAddress(ipWithPrefix string) []string {
//	fmt.Println(ipWithPrefix)
//	parts := strings.Split(ipWithPrefix, "/")
//	if len(parts) == 2 {
//		return parts
//	}
//	return []string{ipWithPrefix, "32"} // Default to /32 if no prefix is specified
//}

func maskToPrefixLength(mask string) string {
	maskMap := map[string]string{
		"255.255.255.255": "32",
		"255.255.255.254": "31",
		"255.255.255.252": "30",
		"255.255.255.248": "29",
		"255.255.255.240": "28",
		"255.255.255.224": "27",
		"255.255.255.192": "26",
		"255.255.255.128": "25",
		"255.255.255.0":   "24",
		"255.255.254.0":   "23",
		"255.255.252.0":   "22",
		"255.255.248.0":   "21",
		"255.255.240.0":   "20",
		"255.255.224.0":   "19",
		"255.255.192.0":   "18",
		"255.255.128.0":   "17",
		"255.255.0.0":     "16",
		"255.254.0.0":     "15",
		"255.252.0.0":     "14",
		"255.248.0.0":     "13",
		"255.240.0.0":     "12",
		"255.224.0.0":     "11",
		"255.192.0.0":     "10",
		"255.128.0.0":     "9",
		"255.0.0.0":       "8",
		"254.0.0.0":       "7",
		"252.0.0.0":       "6",
		"248.0.0.0":       "5",
		"240.0.0.0":       "4",
		"224.0.0.0":       "3",
		"192.0.0.0":       "2",
		"128.0.0.0":       "1",
		"0.0.0.0":         "0",
	}

	if prefix, ok := maskMap[mask]; ok {
		return prefix
	}

	return "32" // Default to /32 if mask is unknown
}

func getAccessLists(config *frrProto.StaticFRRConfiguration) map[string]accessList {
	result := make(map[string]accessList)

	if config == nil || config.AccessList == nil {
		return result
	}

	for name, aclConfig := range config.AccessList {
		if aclConfig == nil {
			continue
		}

		entries := []ACLEntry{}

		for _, item := range aclConfig.AccessListItems {
			if item == nil {
				continue
			}

			entry := ACLEntry{
				IsPermit: item.AccessControl == "permit",
				Sequence: int(item.Sequence),
			}

			switch dest := item.Destination.(type) {
			case *frrProto.AccessListItem_IpPrefix:
				if dest != nil && dest.IpPrefix != nil {
					entry.IPAddress = dest.IpPrefix.IpAddress
					entry.PrefixLength = int(dest.IpPrefix.PrefixLength)
				}
			case *frrProto.AccessListItem_Any:
				entry.IPAddress = "any"
				entry.Any = true
				entry.PrefixLength = 0
			}

			entries = append(entries, entry)
		}

		result[name] = accessList{
			accessListName: name,
			aclEntry:       entries,
		}
	}

	return result
}

func convertToMagicalStateRuntime(config *frrProto.OSPFRouterData) {
	//fmt.Printf("%+v\n", config)
	//result := &magicalState{
	//Hostname: config.GetRouterId(),
	//Areas:    []area{},
	//}

	//areaMap := make(map[string]*area)

	//fmt.Println(result)
	//fmt.Printf("%+v\n", config.GetRouterStates())
	var advertismentList []advertisment
	fmt.Println("########### Start New Print ###########")
	for _, area := range config.GetRouterStates() {
		fmt.Println("--------- Value ---------")
		//fmt.Println(value.LsaEntries["lsa_type"])
		for _, entry := range area.LsaEntries {
			fmt.Println("--------- entry ---------")
			for _, link := range entry.RouterLinks {
				if link.GetNetworkAddress() != "" {
					adv := advertisment{
						InterfaceAddress: link.GetNetworkAddress(),
						PrefixLength:     link.GetNetworkMask(),
						//LsaType:      entry.GetLsaType(),
						// Cost: int(link.GetTos0Metric()),
					}
					advertismentList = append(advertismentList, adv)

				}
				//fmt.Println(link)
			}
			//fmt.Println(entry)
			fmt.Println(advertismentList)
		}
	}
	//fmt.Println(advertismentList)
}

func getStaticRedistributionList(config *frrProto.StaticFRRConfiguration) RedistributionList {
	result := RedistributionList{
		StaticRoutes: []RedistributedRoute{},
		BGPRoutes:    []RedistributedRoute{},
	}

	if config == nil || config.StaticRoutes == nil || config.OspfConfig == nil {
		return result
	}

	// Find static redistribution configuration in OSPF
	var staticRedistConfig *frrProto.Redistribution
	for _, redistribution := range config.OspfConfig.Redistribution {
		if redistribution != nil && redistribution.Type == "static" {
			staticRedistConfig = redistribution
			break
		}
	}

	// If no static redistribution is configured, return empty list
	if staticRedistConfig == nil {
		return result
	}

	// If no route-map is specified, all static routes will be redistributed
	if staticRedistConfig.RouteMap == "" {
		for _, staticRoute := range config.StaticRoutes {
			if staticRoute != nil && staticRoute.IpPrefix != nil {
				result.StaticRoutes = append(result.StaticRoutes, RedistributedRoute{
					IPPrefix:     staticRoute.IpPrefix.IpAddress,
					PrefixLength: int(staticRoute.IpPrefix.PrefixLength),
					NextHop:      staticRoute.NextHop,
					Metric:       staticRedistConfig.Metric,
					MetricType:   "E1", // Default assuming metric-type 1
				})
			}
		}
		return result
	}

	// Get the route-map specified in the redistribution
	routeMapName := staticRedistConfig.RouteMap
	routeMap, exists := config.RouteMap[routeMapName]
	if !exists || routeMap == nil {
		return result
	}

	// Check if the route-map permits and what access list it uses
	if !routeMap.Permit {
		return result // If the route-map is deny, no routes will be redistributed
	}

	accessListName := routeMap.AccessList
	accessList, exists := config.AccessList[accessListName]
	if !exists || accessList == nil {
		return result
	}

	// Filter static routes based on the access list
	for _, staticRoute := range config.StaticRoutes {
		if staticRoute == nil || staticRoute.IpPrefix == nil {
			continue
		}

		// Check if static route matches any permit rule in the access list
		isPermitted := false
		for _, item := range accessList.AccessListItems {
			if item == nil {
				continue
			}

			if item.AccessControl != "permit" {
				continue
			}

			// Check for IP prefix match or "any" match
			switch dest := item.Destination.(type) {
			case *frrProto.AccessListItem_IpPrefix:
				if dest != nil && dest.IpPrefix != nil {
					// Check if the static route prefix is contained within the access list prefix
					if isSubnetOf(staticRoute.IpPrefix, dest.IpPrefix) {
						isPermitted = true
						break
					}
				}
			case *frrProto.AccessListItem_Any:
				isPermitted = true
				break
			}
		}

		if isPermitted {
			result.StaticRoutes = append(result.StaticRoutes, RedistributedRoute{
				IPPrefix:     staticRoute.IpPrefix.IpAddress,
				PrefixLength: int(staticRoute.IpPrefix.PrefixLength),
				NextHop:      staticRoute.NextHop,
				RouteMapName: routeMapName,
				Metric:       staticRedistConfig.Metric,
				MetricType:   "E1", // Assuming metric-type 1 as default
			})
		}
	}

	// Find BGP redistribution configuration in OSPF
	var bgpRedistConfig *frrProto.Redistribution
	for _, redistribution := range config.OspfConfig.Redistribution {
		if redistribution != nil && redistribution.Type == "bgp" {
			bgpRedistConfig = redistribution
			break
		}
	}

	// For demonstration purposes, we're just returning a placeholder for BGP routes
	// since we don't have BGP route information in the config structure
	if bgpRedistConfig != nil {
		result.BGPRoutes = append(result.BGPRoutes, RedistributedRoute{
			RouteMapName: bgpRedistConfig.RouteMap,
			Metric:       bgpRedistConfig.Metric,
			MetricType:   "E1", // Assuming metric-type 1 as default
		})
	}

	return result
}

// Helper function to check if one prefix is a subnet of another
func isSubnetOf(subnet *frrProto.IPPrefix, network *frrProto.IPPrefix) bool {
	// This is a simplified implementation. In a real-world scenario,
	// you would need to convert IP addresses to binary and compare them properly.
	// For now, we'll just check if they have the same IP address and subnet contains network.
	return subnet.IpAddress == network.IpAddress && subnet.PrefixLength >= network.PrefixLength
}
