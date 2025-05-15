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
		peerIpPrefix := ""
		for _, peer := range iface.InterfaceIpPrefixes {
			if peer.PeerIpPrefix != nil {
				peerInterface = true
				peerIpPrefix = peer.PeerIpPrefix.IpAddress
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
				adv.InterfaceAddress = peerIpPrefix
			} else {
				adv.LinkType = "unknown"
				adv.InterfaceAddress = zeroLastOctetString(adv.InterfaceAddress)
				adv.LinkStateId = interfaceIpPrefix.IpPrefix.IpAddress

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

// GetStaticFileExternalData makes LSA type 5 prediction parsing
func GetStaticFileExternalData(config *frrProto.StaticFRRConfiguration, accessList map[string]*frrProto.AccessListAnalyzer, staticRouteMap map[string]*frrProto.StaticList) *frrProto.InterAreaLsa {
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
					if aclAnalyzer == nil {
						continue
					}
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
func GetStaticFileNssaExternalData(config *frrProto.StaticFRRConfiguration, accessList map[string]*frrProto.AccessListAnalyzer, staticRouteMap map[string]*frrProto.StaticList) *frrProto.InterAreaLsa {
	if config == nil || config.OspfConfig == nil {
		return nil
	}

	result := &frrProto.InterAreaLsa{
		Hostname: config.Hostname,
		RouterId: config.OspfConfig.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	// Find the NSSA area from OSPF configuration
	var nssaAreaID string
	for _, area := range config.OspfConfig.Area {
		if area.Type == "nssa" {
			nssaAreaID = area.Name
			break
		}
	}

	// Create a single AreaAnalyzer for all routes
	area := &frrProto.AreaAnalyzer{
		AreaName: nssaAreaID, // Set the area name
		LsaType:  "NSSA-LSA",
		Links:    []*frrProto.Advertisement{},
	}
	result.Areas = append(result.Areas, area)

	// Rest of the function remains the same...
	for _, staticRoute := range config.StaticRoutes {
		ipAddr := staticRoute.IpPrefix.IpAddress
		prefixLen := staticRoute.IpPrefix.PrefixLength

		if _, exists := staticRouteMap[ipAddr]; exists {
			isAllowed := false

			if len(accessList) == 0 {
				isAllowed = true
			} else {
				for _, aclAnalyzer := range accessList {
					if aclAnalyzer == nil {
						continue
					}
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
					LinkType:     "nssa-external",
				}
				area.Links = append(area.Links, advert)
			}
		}
	}

	return result
}

func zeroLastOctetString(ipAddress string) string {
	parts := strings.Split(ipAddress, ".")

	parts[3] = "0"

	return strings.Join(parts, ".")
}
