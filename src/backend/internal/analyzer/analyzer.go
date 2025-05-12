package analyzer

import (
	"net"
	"strconv"
	"strings"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"google.golang.org/protobuf/proto"
)

type RedistributedRoute struct {
	IPPrefix     string `json:"ip_prefix"`
	PrefixLength int    `json:"prefix_length,omitempty"`
	NextHop      string `json:"next_hop,omitempty"`
	RouteMapName string `json:"route_map_name,omitempty"`
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

func (c *Analyzer) AnomalyAnalysis() {

	accessList := GetAccessList(c.metrics.StaticFrrConfiguration)
	staticRouteMap := GetStaticRouteList(c.metrics.StaticFrrConfiguration, accessList)
	peerInterfaceMap := GetPeerNetworkAddress(c.metrics.StaticFrrConfiguration)
	peerNeighborMap := GetPeerNeighbor(c.metrics.OspfNeighbors, peerInterfaceMap)
	hostname := c.metrics.StaticFrrConfiguration.Hostname

	isNssa, shouldRouterLSDB := GetStaticFileRouterData(c.metrics.StaticFrrConfiguration)
	shouldExternalLSDB := GetStaticFileExternalData(c.metrics.StaticFrrConfiguration, accessList, staticRouteMap)

	fibMap := GetFIB(c.metrics.RoutingInformationBase)
	receivedSummaryLSDB := GetRuntimeSummaryData(c.metrics.OspfSummaryDataAll, hostname)
	receivedNetworkLSDB := GetRuntimeNetworkData(c.metrics.OspfNetworkDataAll, hostname)
	receivedExternalLSDB := GetRuntimeExternalData(c.metrics.OspfExternalAll, hostname)
	receivedNssaExternalLSDB := GetRuntimeNssaExternalData(c.metrics.OspfNssaExternalAll, hostname)

	shouldNssaExternalLSDB := GetStaticFileNssaExternalData(c.metrics.StaticFrrConfiguration)

	isRouterLSDB, p2pMap := GetRuntimeRouterDataSelf(c.metrics.OspfRouterData, hostname, peerNeighborMap)

	isExternalLSDB := GetRuntimeExternalDataSelf(c.metrics.OspfExternalData, staticRouteMap, hostname)

	// TODO: testing, mino
	isNssaExternalLSDB := GetNssaExternalDataSelf(c.metrics.OspfNssaExternalData, staticRouteMap, hostname)

	c.RouterAnomalyAnalysisLSDB(accessList, shouldRouterLSDB, isRouterLSDB)
	c.ExternalAnomalyAnalysisLSDB(shouldExternalLSDB, isExternalLSDB)
	// TODO: implement, mino
	if isNssa {
		c.NssaExternalAnomalyAnalysis(accessList, shouldNssaExternalLSDB, isNssaExternalLSDB)
	}
	// TODO: implement ribMap -> fibMap analysis, if necessary?
	c.AnomalyAnalysisFIB(fibMap, receivedNetworkLSDB, receivedSummaryLSDB, receivedExternalLSDB, receivedNssaExternalLSDB)

	//c.UpdateMetrics(p2pMap)
	proto.Merge(c.P2pMap, &p2pMap)
}

func maskToPrefixLength(mask string) string {
	parts := strings.Split(mask, ".")
	if len(parts) != 4 {
		return "32"
	}

	octets := make([]byte, 4)
	for i, p := range parts {
		octet, err := strconv.Atoi(p)
		if err != nil {
			return "32"
		}
		octets[i] = byte(octet)
	}

	ipv4Mask := net.IPv4Mask(octets[0], octets[1], octets[2], octets[3])

	ones, _ := ipv4Mask.Size()

	return strconv.Itoa(ones)
}

func GetAccessList(config *frrProto.StaticFRRConfiguration) map[string]frrProto.AccessListAnalyzer {
	result := make(map[string]frrProto.AccessListAnalyzer)

	if config == nil || config.AccessList == nil {
		return result
	}

	for name, aclConfig := range config.AccessList {
		if aclConfig == nil {
			continue
		}

		var entries []*frrProto.ACLEntry

		for _, item := range aclConfig.AccessListItems {
			if item == nil {
				continue
			}

			entry := frrProto.ACLEntry{
				IsPermit: item.AccessControl == "permit",
				Sequence: int32(item.Sequence),
			}

			switch dest := item.Destination.(type) {
			case *frrProto.AccessListItem_IpPrefix:
				if dest != nil && dest.IpPrefix != nil {
					entry.IPAddress = dest.IpPrefix.IpAddress
					entry.PrefixLength = int32(dest.IpPrefix.PrefixLength)
				}
			case *frrProto.AccessListItem_Any:
				entry.IPAddress = "any"
				entry.Any = true
				entry.PrefixLength = 0
			}

			entries = append(entries, &entry)
		}

		result[name] = frrProto.AccessListAnalyzer{
			AccessList: name,
			AclEntry:   entries,
		}
	}

	return result
}

// Helper function to check if one prefix is a subnet of another
func isSubnetOf(subnet *frrProto.IPPrefix, network *frrProto.IPPrefix) bool {
	return subnet.IpAddress == network.IpAddress && subnet.PrefixLength >= network.PrefixLength
}

// TODO: check with accesslist if it is redistributed in ospf
func GetStaticRouteList(config *frrProto.StaticFRRConfiguration, accessList map[string]frrProto.AccessListAnalyzer) map[string]*frrProto.StaticList {
	if len(config.StaticRoutes) == 0 {
		return nil
	}

	result := map[string]*frrProto.StaticList{}

	for _, route := range config.StaticRoutes {
		//fmt.Println(route)
		result[route.IpPrefix.GetIpAddress()] = &frrProto.StaticList{
			IpAddress:    route.IpPrefix.GetIpAddress(),
			PrefixLength: int32(route.IpPrefix.GetPrefixLength()),
			NextHop:      route.NextHop,
		}
	}

	return result
}

func GetPeerNetworkAddress(config *frrProto.StaticFRRConfiguration) map[string]string {
	peerMap := make(map[string]string)

	for _, iface := range config.Interfaces {
		for _, i := range iface.InterfaceIpPrefixes {
			if i.HasPeer {
				peerMap[iface.Name] = i.IpPrefix.IpAddress
			}
		}
	}

	return peerMap
}

func GetPeerNeighbor(config *frrProto.OSPFNeighbors, peerInterface map[string]string) map[string]string {
	result := map[string]string{}

	for key, neighbors := range config.Neighbors {
		for _, neighbor := range neighbors.Neighbors {
			iface := strings.Split(neighbor.IfaceName, ":")
			if _, exists := peerInterface[iface[0]]; exists {
				result[key] = iface[1]
			}
		}
	}

	return result
}

func (a *Analyzer) UpdateMetrics(p2pMap frrProto.PeerInterfaceMap) {

	proto.Merge(a.P2pMap, &p2pMap)

}
