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

/*
 * interAreaLsa handles LSA type 5 and 7
 * this difference is handled because there are no area distinctions with these types
 */

/*
 * make protobuf structsjjjj
 */

/*
 * There are different kinds of routers
 * Backbone Router, Internal Router, ABR, ASBR
 * Which are important?
 * ABR: An ABR is an Area Boarder Router, it borders different routers
 * What kind of LSAs does it send?
	*
const (
	//BackboneRouter = iota // Has at least one interface in Area 0
	//InternalRouter // Has all of its interfaces in a single area
	Normal = iota // normal router, it's either a BackboneRouter or InternalRouter
	ABR // An OSPF Router that has one or more interfaces in teh backbone area and one or more interfaces in a non-backbone area
	ASBR // Connects to an area and also to an external AS
)

LSA received by area:
	Standar Area:
	- IN: 1,2,3,4,5
	Stub Area:
	- IN: 1,2,3
	Totally Stubby Area:
	- IN: 1,2
	NSSA:
	- IN: 1,2,3,7
	- OUT: 5
	Totally Stubby NSSA:
	- IN: 1,2,7
	- OUT: 5

LSA Types by Router type:
	Normal Router:
	- T1
	- T2 if DR
	ABR:
	- T1
	- T3
	- T4
	- T5
	- T7, if part of NSSA
	ASBR:
	- T1
	- T5
	- T7


Method of finding out what kind of Router we are dealing with:
	Standard:
	- Has all interfaces in only one area
	ABR:
	- Part of Backbone and has two Areas
	ASBR:
	- if it has external routing protocols alongside ospf


*/

type intraAreaLsa struct {
	Hostname   string `json:"hostname"`
	RouterId   string `json:"router_id"`
	RouterType string `json:"router_type"` // normal, asbr, asbr
	Areas      []area `json:"areas"`
}

type interAreaLsa struct {
	Hostname   string `json:"hostname"`
	RouterId   string `json:"router_id"`
	RouterType string `json:"router_type"` // normal, asbr, asbr
	Areas      []area `json:"areas"`
}

type area struct {
	AreaName string                   `json:"name"`
	LsaType  string                   `json:"lsa_type"`
	AreaType string                   `json:"area_type"` //
	Links    []frrProto.Advertisement `json:"links"`
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

type StaticList struct {
	IpAddress    string `json:"ip_address"`
	PrefixLength int    `json:"prefix_length"`
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

type OspfRedistribution struct {
	Type     string `json:"type"`
	RouteMap string `json:"route_map"`
	Metric   string `json:"metric"`
}

// DONE: Currently there is no system checking if a host is ABR, ASBR or just a normal router.
// DONE: add check if host has static routes
/*
	- Static Lists are used to check if LSA Type 5 are required to do
*/
// TODO: Add another router to core that would be a normal internal router without being ABR or ASBR, just a simple internal router.
// TODO: Properly check if the different network types per interface are correctly checked
/*
	- peer to peer is now correct
*/

// TODO: Logic for detecting anomalies is still very naive and early stage. This needs to be cross checked with theory.

// TODO: The current approach is to first analyze the different states as in should state according to configuration and is state according to what's in the LSDB. Is that enough?

// TODO: If an entry is in router, it CANNOT be in external or nssa-external LSDB
// TODO: Vice Versa: if entry is in external, or nssa-external, it CANNOT be in router LSDB

func (c *Analyzer) AnomalyAnalysis() {

	// required to know what routes are distributed
	accessList := getAccessLists(c.metrics.StaticFrrConfiguration)

	// required to know what static routes there are
	// very important for the stakeholder anomaly
	// lsa type 1, 5 and 7 relevant
	staticRouteMap := getStaticRouteList(c.metrics.StaticFrrConfiguration, accessList)

	// should state
	isNssa, predictedRouterLSDB := getStaticFileRouterData(c.metrics.StaticFrrConfiguration)
	predictedExternalLSDB := getStaticFileExternalData(c.metrics.StaticFrrConfiguration)
	predictedNssaExternalLSDB := getStaticFileNssaExternalData(c.metrics.StaticFrrConfiguration)

	// is state
	runtimeRouterLSDB := getRuntimeRouterData(c.metrics.OspfRouterData, c.metrics.StaticFrrConfiguration.Hostname)
	runtimeExternalLSDB := getRuntimeExternalRouterData(c.metrics.OspfExternalData, c.metrics.StaticFrrConfiguration.Hostname)
	runtimeNssaExternalLSDB := getNssaExternalRouterData(c.metrics.OspfNssaExternalData, c.metrics.StaticFrrConfiguration.Hostname)

	// lsa type 1 always needs to be checked
	c.routerAnomalyAnalysis(accessList, predictedRouterLSDB, runtimeRouterLSDB)

	// if router is an ABR/ASBR, lsa type 5 is important
	if len(staticRouteMap) > 0 || isNssa {
		fmt.Println("LSA Type 5 is checked")
		externalAnomalyAnalysis(accessList, predictedExternalLSDB, runtimeExternalLSDB)
	}

	// if router is in a NSSA area, this one is important
	// currently it does nothing
	if isNssa {
		nssaExternalAnomalyAnalysis(accessList, predictedNssaExternalLSDB, runtimeNssaExternalLSDB)
	}

}

//func parseIPAddress(ipWithPrefix string) []string {
//	fmt.Println(ipWithPrefix)
//	parts := strings.Split(ipWithPrefix, "/")
//	if len(parts) == 2 {
//		return parts
//	}
//	return []string{ipWithPrefix, "32"} // Default to /32 if no prefix is specified
//}

// todo: move to utils or use a library, i think https://pkg.go.dev/net should cover that. (i used it for file parser)
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

		var entries []ACLEntry

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
	var advertisementList []frrProto.Advertisement
	fmt.Println("########### Start New Print ###########")
	for _, area := range config.GetRouterStates() {
		fmt.Println("--------- Value ---------")
		//fmt.Println(value.LsaEntries["lsa_type"])
		for _, entry := range area.LsaEntries {
			fmt.Println("--------- entry ---------")
			for _, link := range entry.RouterLinks {
				if link.GetNetworkAddress() != "" {
					adv := frrProto.Advertisement{
						InterfaceAddress: link.GetNetworkAddress(),
						PrefixLength:     link.GetNetworkMask(),
						//LsaType:      entry.GetLsaType(),
						// Cost: int(link.GetTos0Metric()),
					}
					advertisementList = append(advertisementList, adv)

				}
				//fmt.Println(link)
			}
			//fmt.Println(entry)
			fmt.Println(advertisementList)
		}
	}
	//fmt.Println(advertisementList)
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

func getStaticRouteList(config *frrProto.StaticFRRConfiguration, accessList map[string]accessList) map[string]*StaticList {
	if len(config.StaticRoutes) == 0 {
		return nil
	}

	result := map[string]*StaticList{}

	for _, route := range config.StaticRoutes {
		fmt.Println(route)
		result[route.IpPrefix.GetIpAddress()] = &StaticList{
			IpAddress:    route.IpPrefix.GetIpAddress(),
			PrefixLength: int(route.IpPrefix.GetPrefixLength()),
		}
	}

	return result
}
