package analyzer

import (
	"fmt"
	"strconv"
	"strings"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

// lsa type 1 prediction parsing
func GetStaticFileRouterData(config *frrProto.StaticFRRConfiguration) (bool, *frrProto.IntraAreaLsa) {
	if config == nil || config.OspfConfig == nil {
		return false, nil
	}

	isNssa := false

	result := &frrProto.IntraAreaLsa{
		Hostname: config.Hostname,
		RouterId: config.OspfConfig.GetRouterId(),
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	// Map to store unique areas
	areaMap := make(map[string]*frrProto.AreaAnalyzer)

	// Process all interfaces
	for _, iface := range config.Interfaces {
		// Skip interfaces without an area
		peerInterface := false
		for _, peer := range iface.InterfaceIpPrefixes {
			if peer.PeerIpPrefix != nil {
				peerInterface = true
			}
		}
		if iface.Area == "" {
			continue
		}

		// Get or create area entry
		a, exists := areaMap[iface.Area]
		if !exists {
			newArea := frrProto.AreaAnalyzer{
				AreaName: iface.Area,
				LsaType:  "router-LSA",
				AreaType: "normal",
				Links:    []*frrProto.Advertisement{},
			}
			areaMap[iface.Area] = &newArea
			a = &newArea
		}

		// Create advertisements from IP addresses
		for _, ipPrefix := range iface.InterfaceIpPrefixes {
			if ipPrefix.IpPrefix == nil {
				continue
			}

			adv := frrProto.Advertisement{
				InterfaceAddress: ipPrefix.IpPrefix.IpAddress,
				PrefixLength:     strconv.Itoa(int(ipPrefix.IpPrefix.PrefixLength)),
			}

			// Determine link type based on interface properties
			if ipPrefix.Passive {
				adv.LinkType = "stub network"
				adv.InterfaceAddress = zeroLastOctetString(adv.InterfaceAddress)
			} else if strings.Contains(iface.Name, "lo") {
				// Loopback interfaces
				adv.LinkType = "stub network"
			} else {
				// Default to transit network unless we can determine it's point-to-point
				// In a real implementation, you would check for ospf network point-to-point configuration
				adv.LinkType = "transit network"

			}

			if peerInterface {
				adv.LinkType = "point-to-point"
			}

			a.Links = append(a.Links, &adv)
		}
	}

	// Process virtual links if present
	if config.OspfConfig != nil {
		for _, ospfArea := range config.OspfConfig.Area {
			if ospfArea.Type == "" {
				continue
			}

			// Check if this is the area containing virtual links
			if config.OspfConfig.VirtualLinkNeighbor != "" {
				// Get the transit area (where the virtual link is configured)
				a, exists := areaMap[ospfArea.Name]
				if !exists {
					// If area doesn't exist in the map yet, create it
					newArea := frrProto.AreaAnalyzer{
						AreaName: ospfArea.Name,
						LsaType:  "router-LSA",
						AreaType: "normal",
						Links:    []*frrProto.Advertisement{},
					}
					areaMap[ospfArea.Name] = &newArea
					a = &newArea
				}

				// Add virtual link advertisement
				adv := frrProto.Advertisement{
					LinkStateId: config.OspfConfig.VirtualLinkNeighbor,
					LinkType:    "virtual link",
				}

				a.Links = append(a.Links, &adv)
			}

			// check if is nssa
			switch ospfArea.Type {
			case "nssa":
				areaMap[ospfArea.Name].AreaType = "nssa"
				isNssa = true
			case "stub":
				areaMap[ospfArea.Name].AreaType = "stub"
			default:
				areaMap[ospfArea.Name].AreaType = "normal"
			}
		}
	}

	// Convert map to slice for the final result
	for _, a := range areaMap {
		if len(a.Links) > 0 {
			result.Areas = append(result.Areas, a)
		}
	}

	isASBR := false
	for _, redist := range config.OspfConfig.Redistribution {
		if redist.Type == "bgp" {
			isASBR = true
		}
	}

	// If no areas were found, return nil
	if len(result.Areas) == 0 {
		return false, nil
	} else if len(result.Areas) == 1 {
		result.RouterType = "normal"
	} else if isASBR {
		result.RouterType = "asbr"
	} else {
		result.RouterType = "abr"
	}

	return isNssa, result
}

func GetStaticFileExternalData(config *frrProto.StaticFRRConfiguration, accessList map[string]frrProto.AccessListAnalyzer, staticRouteMap map[string]*frrProto.StaticList) *frrProto.InterAreaLsa {
	if config == nil || config.OspfConfig == nil {
		return nil
	}

	result := &frrProto.InterAreaLsa{
		Hostname: config.Hostname,
		RouterId: config.OspfConfig.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	//fmt.Println(accessList)
	//fmt.Println(config)

	// Create a single AreaAnalyzer for all routes
	area := &frrProto.AreaAnalyzer{
		LsaType: "AS-external-LSA",
		Links:   []*frrProto.Advertisement{},
	}
	result.Areas = append(result.Areas, area)

	// Loop through static routes in the configuration
	for _, staticRoute := range config.StaticRoutes {
		ipAddr := staticRoute.IpPrefix.IpAddress
		prefixLen := staticRoute.IpPrefix.PrefixLength

		// Check if this static route is in the staticRouteMap
		if _, exists := staticRouteMap[ipAddr]; exists {
			// Check if this route is allowed by any access list
			isAllowed := false

			for _, aclAnalyzer := range accessList {

				for _, item := range aclAnalyzer.AclEntry {
					if item.IPAddress == ipAddr && item.IsPermit {
						isAllowed = true
						break
					}
				}
				if isAllowed {
					break
				}
			}

			if isAllowed {
				// Create an advertisement for this route
				advert := &frrProto.Advertisement{
					LinkStateId:  ipAddr,
					PrefixLength: fmt.Sprintf("%d", prefixLen),
					LinkType:     "external",
				}
				area.Links = append(area.Links, advert)
			}
		}
	}

	return result

}

func GetStaticFileExternalDataOld(config *frrProto.StaticFRRConfiguration) *frrProto.InterAreaLsa {
	if config == nil || config.OspfConfig == nil {
		return nil
	}
	// Create a new frrProto.InterAreaLsa instance
	result := &frrProto.InterAreaLsa{
		Hostname: config.Hostname,
		RouterId: config.OspfConfig.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{}}

	// Check for OSPF redistribution (potential external advertisements)
	hasRedistribution := false
	for _, redist := range config.OspfConfig.Redistribution {
		// BGP, connected, static, etc. redistribution means the router will advertise external routes
		if redist.Type != "" {
			hasRedistribution = true
			break
		}
	}

	// If no redistribution is configured, router won't generate external LSAs
	if !hasRedistribution {
		return nil
	}

	// Find NSSA areas (for type 7 LSAs)
	nssaAreas := make(map[string]bool)
	for _, ospfArea := range config.OspfConfig.Area {
		if ospfArea.Type == "nssa" {
			nssaAreas[ospfArea.Name] = true
		}
	}

	// Find all areas
	areaMap := make(map[string]bool)
	areaNssaMap := make(map[string]bool)
	areaList := []string{}
	for _, iface := range config.Interfaces {
		areaNssaMap[iface.Area] = false
		//if iface.Area != "" && !nssaAreas[iface.Area] {
		if iface.Area != "" {
			if _, exists := areaMap[iface.Area]; !exists {
				areaMap[iface.Area] = true
				areaList = append(areaList, iface.Area)
			}
		}

	}

	staticRedistMap := make(map[string]bool)
	routeMap := make(map[string]bool)

	for _, redist := range config.OspfConfig.Redistribution {
		if redist.Type != "" && redist.Type == "static" {
			if _, exists := config.RouteMap[redist.RouteMap]; exists {
				if _, exists := staticRedistMap[redist.Type]; !exists && config.RouteMap[redist.RouteMap].Permit {
					staticRedistMap[redist.Type] = true
					for _, access := range config.AccessList[config.RouteMap[redist.RouteMap].AccessList].AccessListItems {
						if access.AccessControl == "permit" {
							if ipPrefixDest, ok := access.Destination.(*frrProto.AccessListItem_IpPrefix); ok {
								routeMap[ipPrefixDest.IpPrefix.IpAddress] = true
							}
						}
					}
				}
			}

		}
	}

	for _, area := range config.OspfConfig.Area {
		if area.Type == "nssa" {
			areaNssaMap[area.Name] = true
		}
	}
	// fmt.Println("====================")
	// fmt.Println(areaNssaMap)

	// For regular AS-external-LSAs (type 5), only if we're not in a stub/nssa only router

	for _, area := range areaList {
		if areaNssaMap[area] {
			// fmt.Printf("FAIL: %v\n", area)
			continue
		}
		// TODO: I would like for AreaName to be useful, but I'm not sure it is
		externalArea := frrProto.AreaAnalyzer{
			//AreaName: area,
			LsaType: "AS-external-LSA", // Type 5
			Links:   []*frrProto.Advertisement{},
		}
		// Add static routes (will be advertised as type 5 in regular areas)
		for _, staticRoute := range config.StaticRoutes {
			if staticRoute.IpPrefix != nil && routeMap[staticRoute.IpPrefix.IpAddress] {
				adv := frrProto.Advertisement{
					LinkStateId:  staticRoute.IpPrefix.IpAddress,
					PrefixLength: strconv.Itoa(int(staticRoute.IpPrefix.PrefixLength)),
					LinkType:     "external",
				}
				externalArea.Links = append(externalArea.Links, &adv)
			}
		}

		//fmt.Println("Before externalArea")
		//fmt.Println(externalArea)

		// Add external area to the result if it has any links
		if len(externalArea.Links) > 0 {
			// Ensure no NSSA-only router (a router with only NSSA areas doesn't generate type 5 LSAs)
			// Check if router has any non-NSSA areas
			//hasNonNssaArea := false
			//for _, ospfArea := range config.OspfConfig.Area {
			//	if ospfArea.Type != "nssa" {
			//		hasNonNssaArea = true
			//		break
			//	}
			//}

			//// Only add type 5 LSAs if router has at least one non-NSSA area
			//if hasNonNssaArea || len(config.OspfConfig.Area) == 0 {
			//	result.Areas = append(result.Areas, &externalArea)
			//}
			result.Areas = append(result.Areas, &externalArea)
		}
	}

	//fmt.Println(result)
	// If no areas were added, return nil (no external LSAs predicted)
	if len(result.Areas) == 0 {
		return nil
	}

	// fmt.Println("===================")
	return result
}

func getStaticFileNssaExternalData(config *frrProto.StaticFRRConfiguration) *frrProto.InterAreaLsa {
	if config == nil || config.OspfConfig == nil {
		return nil
	}

	// Create a new frrProto.InterAreaLsa instance
	result := &frrProto.InterAreaLsa{
		Hostname: config.Hostname,
		RouterId: config.OspfConfig.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	// Check for OSPF redistribution (potential external advertisements)
	hasRedistribution := false
	for _, redist := range config.OspfConfig.Redistribution {
		// BGP, connected, static, etc. redistribution means the router will advertise external routes
		if redist.Type != "" {
			hasRedistribution = true
			break
		}
	}

	// If no redistribution is configured, router won't generate NSSA external LSAs
	if !hasRedistribution {
		return nil
	}

	// Find NSSA areas (for type 7 LSAs)
	nssaAreas := make(map[string]bool)
	for _, ospfArea := range config.OspfConfig.Area {
		if ospfArea.Type == "nssa" {
			nssaAreas[ospfArea.Name] = true
		}
	}

	// If no NSSA areas, router won't generate type 7 LSAs
	if len(nssaAreas) == 0 {
		return nil
	}

	// Find interfaces in NSSA areas
	interfacesByArea := make(map[string][]string)
	for _, iface := range config.Interfaces {
		if iface.Area != "" {
			interfacesByArea[iface.Area] = append(interfacesByArea[iface.Area], iface.Name)
		}
	}

	// Process each NSSA area
	for nssaArea := range nssaAreas {
		nssaAreaObj := frrProto.AreaAnalyzer{
			AreaName: nssaArea,
			LsaType:  "NSSA-LSA", // Type 7
			Links:    []*frrProto.Advertisement{},
		}

		// Add connected interfaces in this NSSA area
		for _, ifaceName := range interfacesByArea[nssaArea] {
			for _, iface := range config.Interfaces {
				if iface.Name == ifaceName {
					for _, ipPrefix := range iface.InterfaceIpPrefixes {
						if ipPrefix.IpPrefix != nil {
							adv := frrProto.Advertisement{
								LinkStateId:  ipPrefix.IpPrefix.IpAddress,
								PrefixLength: strconv.Itoa(int(ipPrefix.IpPrefix.PrefixLength)),
								LinkType:     "nssa-external",
							}
							nssaAreaObj.Links = append(nssaAreaObj.Links, &adv)
						}
					}
				}
			}
		}

		// Add static routes (will be advertised as type 7 in NSSA areas)
		for _, staticRoute := range config.StaticRoutes {
			if staticRoute.IpPrefix != nil {
				adv := frrProto.Advertisement{
					LinkStateId:  staticRoute.IpPrefix.IpAddress,
					PrefixLength: strconv.Itoa(int(staticRoute.IpPrefix.PrefixLength)),
					LinkType:     "nssa-external",
				}
				nssaAreaObj.Links = append(nssaAreaObj.Links, &adv)
			}
		}

		// Add NSSA area to the result if it has any links
		if len(nssaAreaObj.Links) > 0 {
			result.Areas = append(result.Areas, &nssaAreaObj)
		}
	}

	// If no areas were added, return nil (no NSSA external LSAs predicted)
	if len(result.Areas) == 0 {
		return nil
	}

	return result
}

func zeroLastOctetString(ipAddress string) string {
	parts := strings.Split(ipAddress, ".")

	//if len(parts) != 4 {
	//	return "", fmt.Errorf("invalid IP address format: %s", ipAddress)
	//}

	parts[3] = "0"

	return strings.Join(parts, ".")
}

// Get route map mapping

func convertStaticFileOspfRouteMap(config *frrProto.StaticFRRConfiguration) []*OspfRedistribution {

	// TODO: in general make the nil test better
	if config.OspfConfig == nil {
		return nil
	}

	redist := []*OspfRedistribution{}

	for _, redistribution := range config.OspfConfig.Redistribution {
		r := OspfRedistribution{
			Type:     redistribution.Type,
			RouteMap: redistribution.RouteMap,
			Metric:   redistribution.Metric,
		}

		redist = append(redist, &r)
	}

	return redist
}
