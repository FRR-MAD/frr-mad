package analyzer

import (
	"strconv"
	"strings"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

// lsa type 1 prediction parsing
func convertStaticFileRouterData(config *frrProto.StaticFRRConfiguration) *intraAreaLsa {
	if config == nil || config.OspfConfig == nil {
		return nil
	}

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

		// Get or create area entry
		a, exists := areaMap[iface.Area]
		if !exists {
			newArea := area{
				AreaName: iface.Area,
				LsaType:  "router-LSA",
				Links:    []advertisment{},
			}
			areaMap[iface.Area] = &newArea
			a = &newArea
		}

		// Create advertisements from IP addresses
		for _, ipPrefix := range iface.InterfaceIpPrefixes {
			if ipPrefix.IpPrefix == nil {
				continue
			}

			adv := advertisment{
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

				// Check for point-to-point configuration (simplified check)
				// In your full implementation, parse the actual OSPF network type
				if strings.Contains(iface.Name, "peer") || strings.Contains(iface.Name, "p2p") {
					adv.LinkType = "point-to-point"
				}
			}

			a.Links = append(a.Links, adv)
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
					newArea := area{
						AreaName: ospfArea.Name,
						LsaType:  "router-LSA",
						Links:    []advertisment{},
					}
					areaMap[ospfArea.Name] = &newArea
					a = &newArea
				}

				// Add virtual link advertisement
				adv := advertisment{
					LinkStateId: config.OspfConfig.VirtualLinkNeighbor,
					LinkType:    "virtual link",
				}
				a.Links = append(a.Links, adv)
			}
		}
	}

	// Convert map to slice for the final result
	for _, a := range areaMap {
		if len(a.Links) > 0 {
			result.Areas = append(result.Areas, *a)
		}
	}

	// If no areas were found, return nil
	if len(result.Areas) == 0 {
		return nil
	}

	return result
}

func convertStaticFileExternalData(config *frrProto.StaticFRRConfiguration) *interAreaLsa {
	if config == nil || config.OspfConfig == nil {
		return nil
	}

	// Create a new interAreaLsa instance
	result := &interAreaLsa{
		Hostname: config.Hostname,
		RouterId: config.OspfConfig.RouterId,
		Areas:    []area{},
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

	// For regular AS-external-LSAs (type 5), only if we're not in a stub/nssa only router
	externalArea := area{
		AreaName: "External",
		LsaType:  "AS-external-LSA", // Type 5
		Links:    []advertisment{},
	}

	// Add connected interfaces not in NSSA areas
	for _, iface := range config.Interfaces {
		// Skip interfaces in NSSA areas, they'll generate type 7 LSAs instead
		if nssaAreas[iface.Area] {
			continue
		}

		for _, ipPrefix := range iface.InterfaceIpPrefixes {
			if ipPrefix.IpPrefix != nil {
				adv := advertisment{
					LinkStateId:  ipPrefix.IpPrefix.IpAddress,
					PrefixLength: strconv.Itoa(int(ipPrefix.IpPrefix.PrefixLength)),
					LinkType:     "external",
				}
				externalArea.Links = append(externalArea.Links, adv)
			}
		}
	}

	// Add static routes (will be advertised as type 5 in regular areas)
	for _, staticRoute := range config.StaticRoutes {
		if staticRoute.IpPrefix != nil {
			adv := advertisment{
				LinkStateId:  staticRoute.IpPrefix.IpAddress,
				PrefixLength: strconv.Itoa(int(staticRoute.IpPrefix.PrefixLength)),
				LinkType:     "external",
			}
			externalArea.Links = append(externalArea.Links, adv)
		}
	}

	// Add external area to the result if it has any links
	if len(externalArea.Links) > 0 {
		// Ensure no NSSA-only router (a router with only NSSA areas doesn't generate type 5 LSAs)
		// Check if router has any non-NSSA areas
		hasNonNssaArea := false
		for _, ospfArea := range config.OspfConfig.Area {
			if ospfArea.Type != "nssa" {
				hasNonNssaArea = true
				break
			}
		}

		// Only add type 5 LSAs if router has at least one non-NSSA area
		if hasNonNssaArea || len(config.OspfConfig.Area) == 0 {
			result.Areas = append(result.Areas, externalArea)
		}
	}

	// If no areas were added, return nil (no external LSAs predicted)
	if len(result.Areas) == 0 {
		return nil
	}

	return result
}

func convertStaticFileNssaExternalData(config *frrProto.StaticFRRConfiguration) *interAreaLsa {
	if config == nil || config.OspfConfig == nil {
		return nil
	}

	// Create a new interAreaLsa instance
	result := &interAreaLsa{
		Hostname: config.Hostname,
		RouterId: config.OspfConfig.RouterId,
		Areas:    []area{},
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
		nssaAreaObj := area{
			AreaName: nssaArea,
			LsaType:  "NSSA-LSA", // Type 7
			Links:    []advertisment{},
		}

		// Add connected interfaces in this NSSA area
		for _, ifaceName := range interfacesByArea[nssaArea] {
			for _, iface := range config.Interfaces {
				if iface.Name == ifaceName {
					for _, ipPrefix := range iface.InterfaceIpPrefixes {
						if ipPrefix.IpPrefix != nil {
							adv := advertisment{
								LinkStateId:  ipPrefix.IpPrefix.IpAddress,
								PrefixLength: strconv.Itoa(int(ipPrefix.IpPrefix.PrefixLength)),
								LinkType:     "nssa-external",
							}
							nssaAreaObj.Links = append(nssaAreaObj.Links, adv)
						}
					}
				}
			}
		}

		// Add static routes (will be advertised as type 7 in NSSA areas)
		for _, staticRoute := range config.StaticRoutes {
			if staticRoute.IpPrefix != nil {
				adv := advertisment{
					LinkStateId:  staticRoute.IpPrefix.IpAddress,
					PrefixLength: strconv.Itoa(int(staticRoute.IpPrefix.PrefixLength)),
					LinkType:     "nssa-external",
				}
				nssaAreaObj.Links = append(nssaAreaObj.Links, adv)
			}
		}

		// Add NSSA area to the result if it has any links
		if len(nssaAreaObj.Links) > 0 {
			result.Areas = append(result.Areas, nssaAreaObj)
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
