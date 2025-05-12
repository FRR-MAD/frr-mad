package analyzer

import (
	"fmt"
	"strconv"
	"strings"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

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

	areaMap := make(map[string]*frrProto.AreaAnalyzer)

	for _, iface := range config.Interfaces {
		peerInterface := false
		for _, peer := range iface.InterfaceIpPrefixes {
			if peer.PeerIpPrefix != nil {
				peerInterface = true
			}
		}
		if iface.Area == "" {
			continue
		}

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

		for _, interfaceIpPrefix := range iface.InterfaceIpPrefixes {
			if interfaceIpPrefix.IpPrefix == nil {
				continue
			}

			adv := frrProto.Advertisement{
				InterfaceAddress: interfaceIpPrefix.IpPrefix.IpAddress,
				PrefixLength:     strconv.Itoa(int(interfaceIpPrefix.IpPrefix.PrefixLength)),
			}

			if interfaceIpPrefix.Passive {
				adv.LinkType = "stub network"
				adv.InterfaceAddress = zeroLastOctetString(adv.InterfaceAddress)
			} else if peerInterface {
				adv.LinkType = "point-to-point"
			} else {
				adv.LinkType = "transit network"
			}

			a.Links = append(a.Links, &adv)
		}
	}

	if config.OspfConfig != nil {
		for _, ospfArea := range config.OspfConfig.Area {
			if ospfArea.Type == "" {
				continue
			}

			if config.OspfConfig.VirtualLinkNeighbor != "" {
				a, exists := areaMap[ospfArea.Name]
				if !exists {
					newArea := frrProto.AreaAnalyzer{
						AreaName: ospfArea.Name,
						LsaType:  "router-LSA",
						AreaType: "normal",
						Links:    []*frrProto.Advertisement{},
					}
					areaMap[ospfArea.Name] = &newArea
					a = &newArea
				}

				adv := frrProto.Advertisement{
					LinkStateId: config.OspfConfig.VirtualLinkNeighbor,
					LinkType:    "virtual link",
				}

				a.Links = append(a.Links, &adv)
			}

			switch ospfArea.Type {
			case "nssa":
				areaMap[ospfArea.Name].AreaType = "nssa"
				isNssa = true
			case "stub":
				areaMap[ospfArea.Name].AreaType = "stub"
			case "transit (virtual-link)":
				areaMap[ospfArea.Name].AreaType = "transit"
			default:
				areaMap[ospfArea.Name].AreaType = "normal"
			}
		}
	}

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

	if len(result.Areas) == 0 {
		return false, nil
	} else if len(result.Areas) == 1 {
		result.RouterType = "internal router"
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

	area := &frrProto.AreaAnalyzer{
		LsaType: "AS-external-LSA",
		Links:   []*frrProto.Advertisement{},
	}
	result.Areas = append(result.Areas, area)

	for _, staticRoute := range config.StaticRoutes {
		ipAddr := staticRoute.IpPrefix.IpAddress
		prefixLen := staticRoute.IpPrefix.PrefixLength

		if _, exists := staticRouteMap[ipAddr]; exists {
			isAllowed := false

			if len(accessList) == 0 {
				isAllowed = true
			} else {
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
			}

			if isAllowed {
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

// GetStaticFileNssaExternalData makes LSA type 7 prediction parsing
// TODO: finish this
func GetStaticFileNssaExternalData(config *frrProto.StaticFRRConfiguration) *frrProto.InterAreaLsa {
	if config == nil || config.OspfConfig == nil {
		return nil
	}

	result := &frrProto.InterAreaLsa{
		Hostname: config.Hostname,
		RouterId: config.OspfConfig.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	redistributionTypes := make(map[string]bool)
	for _, redist := range config.OspfConfig.Redistribution {
		if redist.Type != "" {
			redistributionTypes[redist.Type] = true
		}
	}

	// If no redistribution is configured, router won't generate NSSA external LSAs
	if len(redistributionTypes) == 0 {
		return nil
	}

	nssaAreas := make(map[string]bool)
	for _, ospfArea := range config.OspfConfig.Area {
		if ospfArea.Type == "nssa" {
			nssaAreas[ospfArea.Name] = true
		}
	}

	if len(nssaAreas) == 0 {
		return nil
	}

	// Find interfaces in NSSA areas for reference
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

		// Only include connected interfaces if "connected" is being redistributed
		if redistributionTypes["connected"] {
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
		}

		// Add static routes if static redistribution is enabled
		if redistributionTypes["static"] {
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

	parts[3] = "0"

	return strings.Join(parts, ".")
}
