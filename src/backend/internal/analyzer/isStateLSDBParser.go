package analyzer

import (
	"fmt"
	"net"
	"strconv"
	"strings"

	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
)

func GetRuntimeRouterDataSelf(config *frrProto.OSPFRouterData, hostname string, peerNeighbor map[string]string) (*frrProto.IntraAreaLsa, *frrProto.PeerInterfaceMap) {
	if config == nil {
		return nil, nil
	}

	result := frrProto.IntraAreaLsa{
		RouterId: config.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	// TODO: change to pointer
	p2pMap := frrProto.PeerInterfaceMap{
		PeerInterfaceToAddress: map[string]string{},
	}

	for areaName, routerArea := range config.RouterStates {
		for lsaName, lsaEntry := range routerArea.LsaEntries {
			var currentArea *frrProto.AreaAnalyzer
			for i := range result.Areas {
				if result.Areas[i].AreaName == areaName {
					currentArea = result.Areas[i]
					break
				}
			}

			if currentArea == nil {
				newArea := frrProto.AreaAnalyzer{
					AreaName: areaName,
					LsaType:  lsaEntry.LsaType,
					Links:    []*frrProto.Advertisement{},
				}
				result.Areas = append(result.Areas, &newArea)
				currentArea = result.Areas[len(result.Areas)-1]
			}

			for routerName, routerLink := range lsaEntry.RouterLinks {
				var ipAddress, prefixLength string
				isStub := false
				if strings.EqualFold(routerLink.LinkType, "Stub Network") {
					routerLink.LinkType = "stub network"
					ipAddress = routerLink.NetworkAddress
					isStub = true
					prefixLength = maskToPrefixLength(routerLink.NetworkMask)
				} else if strings.EqualFold(routerLink.LinkType, "a Transit Network") {
					routerLink.LinkType = "transit network"
					ipAddress = routerLink.RouterInterfaceAddress
				} else {
					if routerLink.RouterInterfaceAddress != "" {
						ipAddress = routerLink.RouterInterfaceAddress
					} else if routerLink.NetworkAddress != "" {
						ipAddress = routerLink.NetworkAddress
					} else {
						continue
					}
				}

				adv := frrProto.Advertisement{}
				adv.InterfaceAddress = ipAddress
				if routerLink.LinkType == "another Router (point-to-point)" {
					adv.LinkType = "point-to-point"
					if strings.HasPrefix(ipAddress, "0") {
						if _, exists := peerNeighbor[routerLink.NeighborRouterId]; exists {
							adv.InterfaceAddress = peerNeighbor[routerLink.NeighborRouterId]
							config.RouterStates[areaName].LsaEntries[lsaName].RouterLinks[routerName].P2PInterfaceAddress =
								peerNeighbor[routerLink.NeighborRouterId]
							p2pMap.PeerInterfaceToAddress[ipAddress] = peerNeighbor[routerLink.NeighborRouterId]
						}
					}
				} else {
					adv.LinkType = routerLink.LinkType
				}

				if isStub {
					adv.PrefixLength = prefixLength
				}

				currentArea.Links = append(currentArea.Links, &adv)
			}
		}
	}

	result.Hostname = hostname

	return &result, &p2pMap
}

func GetRuntimeRouterData(config *frrProto.OSPFRouterData, hostname string) *frrProto.IntraAreaLsa {
	if config == nil {
		return nil
	}

	result := &frrProto.IntraAreaLsa{
		Hostname: hostname,
		RouterId: config.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	routerLsdb := &frrProto.AreaAnalyzer{
		LsaType: "router-LSA",
		Links:   []*frrProto.Advertisement{},
	}

	result.Areas = append(result.Areas, routerLsdb)

	for _, routerStates := range config.RouterStates {
		for _, lsa := range routerStates.LsaEntries {
			for _, routerLink := range lsa.RouterLinks {
				adv := &frrProto.Advertisement{
					Options: lsa.Options,
				}
				switch routerLink.LinkType {
				case "Stub Network":
					adv.InterfaceAddress = routerLink.NetworkAddress
					adv.PrefixLength = routerLink.NetworkMask
				default:
					adv.InterfaceAddress = routerLink.RouterInterfaceAddress
				}
				routerLsdb.Links = append(routerLsdb.Links, adv)
			}
		}
	}

	return result
}

func GetRuntimeNetworkData(config *frrProto.OSPFNetworkData, hostname string) *frrProto.IntraAreaLsa {
	if config == nil {
		return nil
	}
	result := &frrProto.IntraAreaLsa{
		Hostname: hostname,
		RouterId: config.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}
	networkLsdb := &frrProto.AreaAnalyzer{
		LsaType: "network-LSA",
		Links:   []*frrProto.Advertisement{},
	}

	result.Areas = append(result.Areas, networkLsdb)

	for _, netStates := range config.NetStates {
		for _, lsaEntry := range netStates.LsaEntries {
			adv := &frrProto.Advertisement{
				LinkStateId:  getNetworkAddress(lsaEntry.LinkStateId, lsaEntry.NetworkMask) + "/" + strconv.Itoa(int(lsaEntry.NetworkMask)),
				PrefixLength: strconv.Itoa(int(lsaEntry.NetworkMask)),
				Options:      lsaEntry.Options,
			}
			networkLsdb.Links = append(networkLsdb.Links, adv)
		}
	}

	return result

}

func GetRuntimeSummaryData(config *frrProto.OSPFSummaryData, hostname string) *frrProto.InterAreaLsa {
	if config == nil {
		return nil
	}
	result := &frrProto.InterAreaLsa{
		Hostname: hostname,
		RouterId: config.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}
	summaryLsdb := &frrProto.AreaAnalyzer{
		LsaType: "summary-LSA",
		Links:   []*frrProto.Advertisement{},
	}

	result.Areas = append(result.Areas, summaryLsdb)

	for _, sumStates := range config.SummaryStates {
		for _, lsaEntry := range sumStates.LsaEntries {
			adv := &frrProto.Advertisement{
				LinkStateId:  lsaEntry.LinkStateId,
				PrefixLength: strconv.Itoa(int(lsaEntry.NetworkMask)),
				Options:      lsaEntry.Options,
			}
			summaryLsdb.Links = append(summaryLsdb.Links, adv)
		}
	}

	return result
}

// lsa type 5 parsing, this will only return static routes, as BGP routes aren't useful in ospf analysis
// Since AS-external-LSA (type 5) doesn't belong to a specific area,
// we'll create a single "area" to represent the AS external links
func GetRuntimeExternalDataSelf(config *frrProto.OSPFExternalData, staticRouteMap map[string]*frrProto.StaticList, hostname string) *frrProto.InterAreaLsa {
	if config == nil {
		return nil
	}

	// TODO: check if redistribute has a route-map and only if compare to route-map lists
	result := &frrProto.InterAreaLsa{
		Hostname: hostname,
		RouterId: config.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	externalArea := frrProto.AreaAnalyzer{
		LsaType: "AS-external-LSA",
		Links:   []*frrProto.Advertisement{},
	}

	result.Areas = append(result.Areas, &externalArea)

	for key, lsa := range config.AsExternalLinkStates {
		if _, exists := staticRouteMap[key]; !exists {
			continue
		}
		adv := frrProto.Advertisement{
			LinkStateId:  lsa.LinkStateId,
			PrefixLength: strconv.Itoa(int(lsa.NetworkMask)),
			LinkType:     "external",
			Options:      lsa.Options,
		}

		externalArea.Links = append(externalArea.Links, &adv)
	}

	return result
}

func GetRuntimeExternalData(config *frrProto.OSPFExternalAll, hostname string) *frrProto.InterAreaLsa {
	if config == nil {
		return nil
	}

	result := &frrProto.InterAreaLsa{
		Hostname: hostname,
		RouterId: config.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	externalArea := frrProto.AreaAnalyzer{
		LsaType: "AS-external-LSA",
		Links:   []*frrProto.Advertisement{},
	}

	result.Areas = append(result.Areas, &externalArea)

	for _, linkState := range config.AsExternalLinkStates {
		adv := frrProto.Advertisement{
			LinkStateId:  linkState.LinkStateId,
			PrefixLength: strconv.Itoa(int(linkState.NetworkMask)),
			LinkType:     "external",
			Options:      linkState.Options,
		}
		externalArea.Links = append(externalArea.Links, &adv)
	}

	return result
}

// lsa type 7 parsing
func GetNssaExternalData(config *frrProto.OSPFNssaExternalData, staticRouteMap map[string]*frrProto.StaticList, hostname string, logger *logger.Logger) *frrProto.InterAreaLsa {
	if config == nil {
		return nil
	}

	result := &frrProto.InterAreaLsa{
		Hostname: hostname,
		RouterId: config.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	for areaId, nssaArea := range config.NssaExternalLinkStates {
		nssaAreaObj := frrProto.AreaAnalyzer{
			AreaName: areaId,
			LsaType:  "NSSA-LSA",
			Links:    []*frrProto.Advertisement{},
		}

		for key, lsa := range nssaArea.Data {
			if _, exists := staticRouteMap[key]; !exists {
				continue
			}

			pBitSet := false
			fields := strings.Split(lsa.Options, "|")
			if len(fields) > 4 && strings.Contains(fields[4], "P") {
				pBitSet = true
			}

			logger.Info(fmt.Sprintf("NSSA route %s/%s has P-bit: %v", lsa.LinkStateId, strconv.Itoa(int(lsa.NetworkMask)), pBitSet))

			adv := frrProto.Advertisement{
				LinkStateId:  lsa.LinkStateId,
				PrefixLength: strconv.Itoa(int(lsa.NetworkMask)),
				LinkType:     "nssa-external",
			}

			nssaAreaObj.Links = append(nssaAreaObj.Links, &adv)
		}

		result.Areas = append(result.Areas, &nssaAreaObj)
	}

	return result
}

func GetRuntimeNssaExternalData(config *frrProto.OSPFNssaExternalAll, hostname string) *frrProto.InterAreaLsa {
	if config == nil {
		return nil
	}

	result := &frrProto.InterAreaLsa{
		Hostname: hostname,
		RouterId: config.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	externalArea := frrProto.AreaAnalyzer{
		LsaType: "NSSA-LSA",
		Links:   []*frrProto.Advertisement{},
	}

	result.Areas = append(result.Areas, &externalArea)

	for _, linkStates := range config.NssaExternalAllLinkStates {
		for _, linkState := range linkStates.Data {
			adv := frrProto.Advertisement{
				LinkStateId:  linkState.LinkStateId,
				PrefixLength: strconv.Itoa(int(linkState.NetworkMask)),
				LinkType:     "nssa-external",
				Options:      linkState.Options,
			}
			externalArea.Links = append(externalArea.Links, &adv)
		}
	}

	return result

}

func GetFIB(rib *frrProto.RoutingInformationBase) map[string]frrProto.RibPrefixes {

	OspfFibMap := map[string]frrProto.RibPrefixes{}
	for prefix, routes := range rib.Routes {
		for _, routeEntry := range routes.Routes {
			for _, route := range routeEntry.Nexthops {
				if route.Fib {
					OspfFibMap[prefix] = frrProto.RibPrefixes{
						Prefix:         routeEntry.Prefix,
						PrefixLength:   strconv.FormatInt(int64(routeEntry.PrefixLen), 10),
						NextHopAddress: route.Ip,
						Protocol:       routeEntry.Protocol,
					}
				}
			}
		}
	}

	return OspfFibMap
}

func getNetworkAddress(prefix string, prefixLength int32) string {
	ip := net.ParseIP(prefix)

	tmpNet := &net.IPNet{
		IP:   ip,
		Mask: net.CIDRMask(int(prefixLength), 32),
	}

	network := tmpNet.IP.Mask(tmpNet.Mask)

	return network.String()

}
