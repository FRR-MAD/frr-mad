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
	// static file parsing
	fmt.Println("#################### File Configuration Router LSDB Prediction ####################")
	predictedRouterLSDB := convertStaticFileRouterData(c.metrics.StaticFrrConfiguration)
	if predictedRouterLSDB != nil {
		for _, area := range predictedRouterLSDB.Areas {
			fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
		}
	}
	fmt.Printf("\n%+v\n", predictedRouterLSDB)
	fmt.Println()
	fmt.Println()

	fmt.Println("#################### File Configuration External LSDB Prediction ####################")
	predictedExternalLSDB := convertStaticFileExternalData(c.metrics.StaticFrrConfiguration)
	if predictedExternalLSDB != nil {
		for _, area := range predictedExternalLSDB.Areas {
			fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
		}
	}
	fmt.Printf("\n%+v\n", predictedExternalLSDB)
	fmt.Println()
	fmt.Println()

	fmt.Println("#################### File Configuration NSSA External LSDB Prediction ####################")
	predictedNssaExternalLSDB := convertStaticFileNssaExternalData(c.metrics.StaticFrrConfiguration)
	if predictedNssaExternalLSDB != nil {

		for _, area := range predictedNssaExternalLSDB.Areas {
			fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
		}
	}
	fmt.Printf("\n%+v\n", predictedNssaExternalLSDB)
	fmt.Println()
	fmt.Println()

	// runtime parsing
	fmt.Println("#################### Runtime Configuration Router LSDB IS_STATE ####################")
	runtimeRouterLSDB := convertRuntimeRouterData(c.metrics.OspfRouterData, c.metrics.StaticFrrConfiguration.Hostname)
	for _, area := range runtimeRouterLSDB.Areas {
		fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
	}
	fmt.Printf("\n%+v\n", runtimeRouterLSDB)
	fmt.Println()
	fmt.Println()

	fmt.Println("#################### Runtime Configuration External LSDB IS_STATE ####################")
	runtimeExternalLSDB := convertRuntimeExternalRouterData(c.metrics.OspfExternalData, c.metrics.StaticFrrConfiguration.Hostname)
	for _, area := range runtimeExternalLSDB.Areas {
		fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
	}
	fmt.Printf("\n%+v\n", runtimeExternalLSDB)
	fmt.Println()
	fmt.Println()

	fmt.Println("#################### Runtime Configuration NSSA External LSDB IS_STATE ####################")
	runtimeNssaExternalLSDB := convertNssaExternalRouterData(c.metrics.OspfNssaExternalData, c.metrics.StaticFrrConfiguration.Hostname)
	for _, area := range runtimeNssaExternalLSDB.Areas {
		fmt.Printf("Length of area %s: %d\n", area.AreaName, len(area.Links))
	}
	fmt.Printf("\n%+v\n", runtimeNssaExternalLSDB)
	fmt.Println()
	fmt.Println()

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
}
