package analyzer

import (
	"fmt"
	"strings"

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
	IPAddress    string
	PrefixLength int
	IsPermit     bool
	Any          bool
}

type accessList struct {
	accessListName string
	aclEntry       []ACLEntry
}

func (c *Analyzer) AnomalyAnalysis() {
	//fmt.Println(c.metrics.StaticFrrConfiguration.GetOspfConfig())
	//redistributionList := []string{
	//	"ospf",
	//}
	//redistributionList = append(redistributionList, "static")
	//c.Logger.Debug(fmt.Sprintf("%v", staticConfig))

	// basic StaticFRRConfiguration
	//fmt.Println(c.metrics.StaticFrrConfiguration)

	// access list stuff
	fmt.Println(getAccessLists(c.metrics.StaticFrrConfiguration))

	// analysis of ospf networks
	// still needs to do cross checking with a potential access list?
	// for that it's better to create a variety of different configurations
	fmt.Println("#################### File Configuration Router LSDB Prediction ####################")
	configuredRouterLSDB := convertStaticFileRouterData(c.metrics.StaticFrrConfiguration)
	for _, area := range configuredRouterLSDB.Areas {
		fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
	}
	fmt.Printf("\n%v\n", configuredRouterLSDB)
	fmt.Println()
	fmt.Println()

	// analysis of runtime configuration

	//fmt.Println(c.metrics.OspfRouterData)
	fmt.Println("#################### Runtime Configuration Router LSDB IS_STATE ####################")
	runtimeRouterLSDB := convertOSPFRouterData(c.metrics.OspfRouterData, c.metrics.StaticFrrConfiguration.Hostname)
	for _, area := range runtimeRouterLSDB.Areas {
		fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
	}
	fmt.Printf("\n%v\n", runtimeRouterLSDB)
	//for _, area := range runTimeRouterLSDB {
	//}

	//getFileStaticAdvertisment(c.metrics.StaticFrrConfiguration)

	//	for _, redist := range c.metrics.StaticFrrConfiguration.OspfConfig.GetRedistribution() { fmt.Println(redist.Type)
	//		switch redist.Type {
	//		case "ospf":
	//
	//		case "static":
	//		case "bgp":
	//			//if c.metrics.StaticFrrConfiguration.OspfConfig.Area
	//			// check if area is nssa
	//			// if area is nssa, some bgp distribution into ospf should be visible
	//		default:
	//			continue
	//		}
	//	}

}

func convertOSPFRouterData(config *frrProto.OSPFRouterData, hostname string) *intraAreaLsa {
	result := intraAreaLsa{
		RouterId: config.RouterId,
		Areas:    []area{},
	}

	for areaName, routerArea := range config.RouterStates {
		for _, lsaEntry := range routerArea.LsaEntries {
			var currentArea *area
			for i := range result.Areas {
				if result.Areas[i].AreaName == areaName {
					currentArea = &result.Areas[i]
					break
				}
			}

			if currentArea == nil {
				newArea := area{
					AreaName: areaName,
					LsaType:  lsaEntry.LsaType,
					Links:    []advertisment{},
				}
				result.Areas = append(result.Areas, newArea)
				currentArea = &result.Areas[len(result.Areas)-1]
			}

			for _, routerLink := range lsaEntry.RouterLinks {
				var ipAddress, prefixLength string
				isStub := false
				if routerLink.LinkType == "Stub Network" {
					ipAddress = routerLink.NetworkAddress
					isStub = true
					prefixLength = maskToPrefixLength(routerLink.NetworkMask)
				} else if routerLink.LinkType == "a Transit Network" {
					ipAddress = routerLink.RouterInterfaceAddress
					//prefixLength = "24" // Assuming a /24 for transit links
				} else {
					if routerLink.RouterInterfaceAddress != "" {
						ipAddress = routerLink.RouterInterfaceAddress
					} else if routerLink.NetworkAddress != "" {
						ipAddress = routerLink.NetworkAddress
						//prefixLength = maskToPrefixLength(routerLink.NetworkMask)
					} else {
						continue
					}
				}

				adv := advertisment{}
				adv.InterfaceAddress = ipAddress
				adv.LinkType = routerLink.LinkType
				fmt.Println(isStub)
				if isStub {
					adv.PrefixLength = prefixLength
				}
				//adv := advertisment{
				//	InterfaceAddress: ipAddress,
				//	PrefixLength:     prefixLength,
				//	//Cost:             int(routerLink.Tos0Metric),
				//	LinkType:         routerLink.LinkType,
				//}

				currentArea.Links = append(currentArea.Links, adv)
			}
		}
	}

	result.Hostname = hostname

	return &result
}

func parseIPAddress(ipWithPrefix string) []string {
	fmt.Println(ipWithPrefix)
	parts := strings.Split(ipWithPrefix, "/")
	if len(parts) == 2 {
		return parts
	}
	return []string{ipWithPrefix, "32"} // Default to /32 if no prefix is specified
}

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

func getFileStaticAdvertisment(config *frrProto.StaticFRRConfiguration) {

}

func printAccessLists(accessLists map[string]accessList) {
	for _, acl := range accessLists {
		fmt.Printf("Access List: %s\n", acl.accessListName)

		for i, entry := range acl.aclEntry {
			action := "deny"
			if entry.IsPermit {
				action = "permit"
			}

			if entry.Any {
				fmt.Printf("  Rule %d: %s any\n", i+1, action)
			} else {
				fmt.Printf("  Rule %d: %s %s/%d\n", i+1, action, entry.IPAddress, entry.PrefixLength)
			}
		}
		fmt.Println()
	}
}

func getAccessLists(config *frrProto.StaticFRRConfiguration) map[string]accessList {
	// Create result map: access list name -> AccessList struct
	result := make(map[string]accessList)

	// Iterate through all access lists in the configuration
	for name, aclConfig := range config.AccessList {
		entries := []ACLEntry{}

		// Process each access list item
		for _, item := range aclConfig.AccessListItems {
			// Create a new ACL entry
			entry := ACLEntry{
				IsPermit: item.AccessControl == "permit",
			}

			// Check which type of destination we have
			switch dest := item.Destination.(type) {
			case *frrProto.AccessListItem_IpPrefix:
				entry.IPAddress = dest.IpPrefix.IpAddress
				entry.PrefixLength = int(dest.IpPrefix.PrefixLength)
			case *frrProto.AccessListItem_Any:
				entry.IPAddress = "any"
				//dest.IpPrefix.IpAddress
				entry.Any = true
			}

			entries = append(entries, entry)
		}

		// Create AccessList struct and add to result map
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
	fmt.Println("###########  End New Print  ###########")
}

// ConvertToMagicalState converts a StaticFRRConfiguration to a magicalState
func convertStaticFileRouterData(config *frrProto.StaticFRRConfiguration) *intraAreaLsa {
	result := &intraAreaLsa{
		Hostname: config.Hostname,
		RouterId: config.OspfConfig.GetRouterId(),
		Areas:    []area{},
	}

	// Map to store unique areas
	areaMap := make(map[string]*area)

	// Process all interfaces
	for _, iface := range config.Interfaces {
		// Skip interfaces without an area
		if iface.Area == "" {
			continue
		}
		linkType := "a Transit Network"
		if iface.Passive {
			linkType = "Stub Network"
		}

		// Get or create area
		a, exists := areaMap[iface.Area]
		advertismentList := make([]advertisment, 0)
		if !exists {
			newArea := area{
				AreaName: iface.Area,
				LsaType:  "router-lsa", // Default LSA type for areas
				Links:    advertismentList,
			}
			areaMap[iface.Area] = &newArea
			a = &newArea
		}

		// Create advertisements from IP addresses
		var adv advertisment
		for _, ip := range iface.IpAddress {
			adv.InterfaceAddress = ip.IpAddress
			//			adv.Cost = 10
			adv.LinkType = linkType
			if iface.Passive {
				adv.PrefixLength = fmt.Sprintf("%d", ip.PrefixLength)
			}

			a.Links = append(a.Links, adv)
		}

	}

	// Convert map to slice for the final result
	for _, a := range areaMap {
		result.Areas = append(result.Areas, *a)
	}

	return result
}

func Example() {
	// Convert to magicalState
}

func (a *Analyzer) Foobar() string {
	return "mighty analyzer"
}
