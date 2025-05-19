package analyzer

import (
	"fmt"
	"net"
	"strconv"
	"strings"
	"time"

	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
)

func GetRuntimeRouterDataSelf(config *frrProto.OSPFRouterData, hostname string, peerNeighbor map[string]string, logger *logger.Logger) (*frrProto.IntraAreaLsa, frrProto.PeerInterfaceMap) {
	logger.Debug("Parsing router LSDB data")
	start := time.Now()

	intraAreaLsa := &frrProto.IntraAreaLsa{
		RouterId: config.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	// TODO: change to pointer
	p2pMap := frrProto.PeerInterfaceMap{
		PeerInterfaceToAddress: map[string]string{},
	}

	logger.WithAttrs(map[string]interface{}{
		"router_id":      config.RouterId,
		"areas":          len(config.RouterStates),
		"peer_neighbors": len(peerNeighbor),
	}).Debug("Starting router LSDB parsing")

	for areaName, routerArea := range config.RouterStates {
		for lsaName, lsaEntry := range routerArea.LsaEntries {
			var currentArea *frrProto.AreaAnalyzer
			for i := range intraAreaLsa.Areas {
				if intraAreaLsa.Areas[i].AreaName == areaName {
					currentArea = intraAreaLsa.Areas[i]
					break
				}
			}

			if currentArea == nil {
				newArea := frrProto.AreaAnalyzer{
					AreaName: areaName,
					LsaType:  lsaEntry.LsaType,
					Links:    []*frrProto.Advertisement{},
				}
				intraAreaLsa.Areas = append(intraAreaLsa.Areas, &newArea)
				currentArea = intraAreaLsa.Areas[len(intraAreaLsa.Areas)-1]
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

	intraAreaLsa.Hostname = hostname

	if len(p2pMap.PeerInterfaceToAddress) > 0 {
		logger.WithAttrs(map[string]interface{}{
			"p2p_mappings": len(p2pMap.PeerInterfaceToAddress),
		}).Debug("Created peer interface mappings")
	}

	logger.WithAttrs(map[string]interface{}{
		"duration":     time.Since(start).String(),
		"areas_parsed": len(intraAreaLsa.Areas),
		"total_links":  countTotalLinks(intraAreaLsa),
	}).Debug("Completed router LSDB parsing")

	return intraAreaLsa, p2pMap
}

func GetRuntimeRouterData(config *frrProto.OSPFRouterData, hostname string, logger *logger.Logger) *frrProto.IntraAreaLsa {
	if config == nil {
		logger.Debug("Skipping router data parsing - nil input")
		return nil
	}

	logger.Debug("Parsing full router LSDB data")
	start := time.Now()

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

	logger.WithAttrs(map[string]interface{}{
		"duration":    time.Since(start).String(),
		"lsas_parsed": len(routerLsdb.Links),
	}).Debug("Completed full router LSDB parsing")

	return result
}

func GetRuntimeNetworkData(config *frrProto.OSPFNetworkData, hostname string, logger *logger.Logger) *frrProto.IntraAreaLsa {
	if config == nil {
		logger.Debug("Skipping network data parsing - nil input")
		return nil
	}

	logger.Debug("Parsing network LSDB data")
	start := time.Now()

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

	logger.WithAttrs(map[string]interface{}{
		"duration":     time.Since(start).String(),
		"network_lsas": len(networkLsdb.Links),
	}).Debug("Completed network LSDB parsing")

	return result

}

func GetRuntimeSummaryData(config *frrProto.OSPFSummaryData, hostname string, logger *logger.Logger) *frrProto.InterAreaLsa {
	if config == nil {
		logger.Debug("Skipping summary data parsing - nil input")
		return nil
	}

	logger.Debug("Parsing summary LSDB data")
	start := time.Now()

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

	logger.WithAttrs(map[string]interface{}{
		"duration":     time.Since(start).String(),
		"summary_lsas": len(summaryLsdb.Links),
	}).Debug("Completed summary LSDB parsing")

	return result
}

// lsa type 5 parsing, this will only return static routes, as BGP routes aren't useful in ospf analysis
// Since AS-external-LSA (type 5) doesn't belong to a specific area,
// we'll create a single "area" to represent the AS external links
func GetRuntimeExternalDataSelf(config *frrProto.OSPFExternalData, staticRouteMap map[string]*frrProto.StaticList, hostname string, logger *logger.Logger) *frrProto.InterAreaLsa {
	if config == nil {
		logger.Debug("Skipping external data parsing - nil input")
		return nil
	}

	logger.Debug("Parsing self-originated external LSDB data")
	start := time.Now()

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

	logger.WithAttrs(map[string]interface{}{
		"duration":      time.Since(start).String(),
		"external_lsas": len(externalArea.Links),
		"static_routes": len(staticRouteMap),
	}).Debug("Completed self-originated external LSDB parsing")

	return result
}

func GetRuntimeExternalData(config *frrProto.OSPFExternalAll, hostname string, logger *logger.Logger) *frrProto.InterAreaLsa {
	if config == nil {
		logger.Debug("Skipping external data parsing - nil input")
		return nil
	}

	logger.Debug("Parsing all external LSDB data")
	start := time.Now()

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

	logger.WithAttrs(map[string]interface{}{
		"duration":      time.Since(start).String(),
		"external_lsas": len(externalArea.Links),
	}).Debug("Completed external LSDB parsing")

	return result
}

// lsa type 7 parsing
func GetNssaExternalData(config *frrProto.OSPFNssaExternalData, staticRouteMap map[string]*frrProto.StaticList, hostname string, logger *logger.Logger) *frrProto.InterAreaLsa {
	if config == nil {
		logger.Debug("Skipping NSSA external data parsing - nil input")
		return nil
	}

	logger.Debug("Parsing NSSA external LSDB data")
	start := time.Now()

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

	logger.WithAttrs(map[string]interface{}{
		"duration":   time.Since(start).String(),
		"nssa_areas": len(result.Areas),
		"nssa_lsas":  countTotalLinksNSSA(result),
	}).Debug("Completed NSSA external LSDB parsing")

	return result
}

func GetRuntimeNssaExternalData(config *frrProto.OSPFNssaExternalAll, hostname string, logger *logger.Logger) *frrProto.InterAreaLsa {
	if config == nil {
		logger.Debug("Skipping NSSA external data parsing - nil input")
		return nil
	}

	logger.Debug("Parsing all NSSA external LSDB data")
	start := time.Now()

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

	logger.WithAttrs(map[string]interface{}{
		"duration":  time.Since(start).String(),
		"nssa_lsas": len(externalArea.Links),
	}).Debug("Completed NSSA external LSDB parsing")

	return result

}

func GetFIB(rib *frrProto.RoutingInformationBase, logger *logger.Logger) map[string]frrProto.RibPrefixes {
	if rib == nil {
		logger.Debug("Skipping FIB parsing - nil RIB input")
		return nil
	}

	logger.Debug("Building FIB mapping from RIB")
	start := time.Now()

	OspfFibMap := map[string]frrProto.RibPrefixes{}
	ospfCount := 0

	for prefix, routes := range rib.Routes {
		for _, routeEntry := range routes.Routes {
			ospfCount++
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

	logger.WithAttrs(map[string]interface{}{
		"duration":          time.Since(start).String(),
		"ospf_routes":       ospfCount,
		"total_fib_entries": len(OspfFibMap),
	}).Debug("Completed FIB mapping")

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

// Helper functions
func countTotalLinks(lsa *frrProto.IntraAreaLsa) int {
	count := 0
	for _, area := range lsa.Areas {
		count += len(area.Links)
	}
	return count
}

func countTotalLinksNSSA(lsa *frrProto.InterAreaLsa) int {
	count := 0
	for _, area := range lsa.Areas {
		count += len(area.Links)
	}
	return count
}
