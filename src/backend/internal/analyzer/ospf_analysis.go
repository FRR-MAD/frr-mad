package analyzer

import (
	"strings"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	// "google.golang.org/protobuf/proto"
)

// TODO: Find a fix for point-to-point. The representation of p2p in frr is unclear to me. Thus it's removed from any testing
func (a *Analyzer) RouterAnomalyAnalysis(accessList map[string]frrProto.AccessListAnalyzer, shouldState *frrProto.IntraAreaLsa, isState *frrProto.IntraAreaLsa) {
	if isState == nil || shouldState == nil {
		//fmt.Println("nil!")
		return
	}

	result := &frrProto.AnomalyDetection{
		HasMisconfiguredPrefixes: false,
		SuperfluousEntries:       []*frrProto.Advertisement{},
		MissingEntries:           []*frrProto.Advertisement{},
		DuplicateEntries:         []*frrProto.Advertisement{},
	}

	isStateMap := make(map[string]*frrProto.Advertisement)
	isStateCounter := make(map[string]int)
	shouldStateMap := make(map[string]*frrProto.Advertisement)
	shouldStateCounter := make(map[string]int)

	for _, area := range isState.Areas {
		for i := range area.Links {
			if area.Links[i].LinkType != "point-to-point" {
				link := area.Links[i]
				key := getadvertisementKey(link)
				isStateMap[key] = &frrProto.Advertisement{
					InterfaceAddress: link.InterfaceAddress,
					LinkType:         link.LinkType,
				}
				if link.LinkType == "Stub Network" {
					isStateMap[key].PrefixLength = link.PrefixLength
				}
				if _, exist := isStateCounter[link.InterfaceAddress]; !exist {
					isStateCounter[link.InterfaceAddress] = 0
				}
				isStateCounter[link.InterfaceAddress] += 1
			}
		}
	}

	for _, area := range shouldState.Areas {
		for i := range area.Links {
			if area.Links[i].LinkType != "point-to-point" {
				link := area.Links[i]
				key := getadvertisementKey(link)
				shouldStateMap[key] = &frrProto.Advertisement{
					InterfaceAddress: link.InterfaceAddress,
					LinkType:         link.LinkType,
				}
				if link.LinkType == "Stub Network" {
					shouldStateMap[key].PrefixLength = link.PrefixLength
				}
				if _, exist := shouldStateCounter[link.InterfaceAddress]; !exist {
					shouldStateCounter[link.InterfaceAddress] = 0
				}
				shouldStateCounter[link.InterfaceAddress] += 1
			}
		}
	}

	// check for missing prefixes -> underadvertised
	for key, shouldLink := range shouldStateMap {
		if _, exists := isStateMap[key]; !exists {
			if !isExcludedByAccessList(shouldLink, accessList) {
				result.MissingEntries = append(result.MissingEntries, shouldLink)
			}
		}
	}

	// check for to many advertisements -> overadvertised
	for key, isLink := range isStateMap {
		if _, exists := shouldStateMap[key]; !exists {
			result.SuperfluousEntries = append(result.SuperfluousEntries, isLink)
		}
	}

	// check for duplicates
	for prefix, counter := range isStateCounter {
		if counter > 1 {
			result.SuperfluousEntries = append(result.DuplicateEntries, isStateMap[prefix])
		}
	}

	// fmt.Println(isStateMap)
	// fmt.Println(shouldStateMap)

	//a.AnalysisResult.RouterAnomaly.HasUnderAdvertisedPrefixes = writeBoolTarget(result.HasUnderAdvertisedPrefixes)
	//a.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes = writeBoolTarget(result.HasOverAdvertisedPrefixes)

	a.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes = len(result.MissingEntries) > 0
	a.AnalysisResult.RouterAnomaly.HasUnderAdvertisedPrefixes = len(result.SuperfluousEntries) > 0
	a.AnalysisResult.RouterAnomaly.HasDuplicatePrefixes = len(result.DuplicateEntries) > 0
	//writeBoolTarget(result.HasDuplicatePrefixes)
	//a.AnalysisResult.RouterAnomaly.HasMisconfiguredPrefixes = writeBoolTarget(result.HasMisconfiguredPrefixes)
	a.AnalysisResult.RouterAnomaly.MissingEntries = result.MissingEntries
	a.AnalysisResult.RouterAnomaly.SuperfluousEntries = result.SuperfluousEntries
	a.AnalysisResult.RouterAnomaly.DuplicateEntries = result.DuplicateEntries

}

func writeBoolTarget(source bool) bool {
	if source {
		return source
	}
	return false
}

func getadvertisementKey(adv *frrProto.Advertisement) string {
	if adv.InterfaceAddress != "" {
		return normalizeNetworkAddress(adv.InterfaceAddress)
	}
	return normalizeNetworkAddress(adv.LinkStateId)
}

func isExcludedByAccessList(adv *frrProto.Advertisement, accessLists map[string]frrProto.AccessListAnalyzer) bool {
	for _, acl := range accessLists {
		for _, entry := range acl.AclEntry {
			if !entry.IsPermit {
				if entry.Any {
					return true
				} else {
					networkAddr := normalizeNetworkAddress(adv.InterfaceAddress)
					if networkAddr == entry.IPAddress {
						return true
					}
				}
			}
		}
	}

	return false
}

// func ExternalAnomalyAnalysis(accessList map[string]frrProto.AccessListAnalyzer, isState *frrProto.InterAreaLsa, shouldState *frrProto.InterAreaLsa) {
// 	fmt.Println(accessList)
// 	fmt.Println(isState)
// 	fmt.Println(shouldState)
// }

// func (a *Analyzer) externalAnomalyAnalysis(accessList map[string]frrProto.AccessListAnalyzer, isState *frrProto.InterAreaLsa, shouldState *frrProto.InterAreaLsa) {

func (a *Analyzer) ExternalAnomalyAnalysis(accessList map[string]frrProto.AccessListAnalyzer, isState *frrProto.InterAreaLsa, shouldState *frrProto.InterAreaLsa) {
	if isState == nil || shouldState == nil {
		//fmt.Println("nil!")
		return
	}

	result := &frrProto.AnomalyDetection{
		HasMisconfiguredPrefixes: false,
		SuperfluousEntries:       []*frrProto.Advertisement{},
		MissingEntries:           []*frrProto.Advertisement{},
		DuplicateEntries:         []*frrProto.Advertisement{},
	}

	//fmt.Println(accessList)
	//fmt.Println(isState)
	//fmt.Println(shouldState)

	a.AnalysisResult.ExternalAnomaly.HasOverAdvertisedPrefixes = len(result.MissingEntries) > 0
	a.AnalysisResult.ExternalAnomaly.HasUnderAdvertisedPrefixes = len(result.SuperfluousEntries) > 0
	a.AnalysisResult.ExternalAnomaly.HasDuplicatePrefixes = len(result.DuplicateEntries) > 0
	//writeBoolTarget(result.HasDuplicatePrefixes)
	//a.AnalysisResult.RouterAnomaly.HasMisconfiguredPrefixes = writeBoolTarget(result.HasMisconfiguredPrefixes)
	a.AnalysisResult.ExternalAnomaly.MissingEntries = result.MissingEntries
	a.AnalysisResult.ExternalAnomaly.SuperfluousEntries = result.SuperfluousEntries
	a.AnalysisResult.ExternalAnomaly.DuplicateEntries = result.DuplicateEntries

	//result := &frrProto.AnomalyAnalysis{
	//	AnomalyDetection: &frrProto.AnomalyDetection{},
	//	MissingEntries:   []*frrProto.Advertisement{},
	//	SuperfluousEntries:     []*frrProto.Advertisement{},
	//}

	// if isState == nil || shouldState == nil {
	// 	return
	// }

	// isStateMap := make(map[string]*frrProto.Advertisement)
	// shouldStateMap := make(map[string]*frrProto.Advertisement)

	// for _, area := range isState.Areas {
	// 	for i := range area.Links {
	// 		link := area.Links[i]
	// 		key := normalizeNetworkAddress(link.LinkStateId)
	// 		isStateMap[key] = link
	// 	}
	// }

	// for _, area := range shouldState.Areas {
	// 	for i := range area.Links {
	// 		link := area.Links[i]
	// 		key := normalizeNetworkAddress(link.LinkStateId)
	// 		shouldStateMap[key] = link
	// 	}
	// }

	// accessListNetworks := make(map[string]bool)
	// for _, acl := range accessList {
	// 	for _, entry := range acl.AclEntry {
	// 		if entry.IsPermit {
	// 			network := fmt.Sprintf("%s/%d", entry.IPAddress, entry.PrefixLength)
	// 			accessListNetworks[normalizeNetworkAddress(network)] = true
	// 		}
	// 	}
	// }

	// for key, shouldLink := range shouldStateMap {
	// 	if _, exists := isStateMap[key]; !exists {
	// 		networkWithPrefix := fmt.Sprintf("%s/%s", shouldLink.LinkStateId, shouldLink.PrefixLength)
	// 		normalizedNetwork := normalizeNetworkAddress(networkWithPrefix)

	// 		if accessListNetworks[normalizedNetwork] || isInAccessList(shouldLink.LinkStateId, accessList) {
	// 			result.MissingEntries = append(result.MissingEntries, shouldLink)
	// 		}
	// 	}
	// }

	// result.HasUnderAdvertisedPrefixes = len(result.MissingEntries) > 0
	// result.HasOverAdvertisedPrefixes = false
	// result.HasDuplicatePrefixes = false
	// result.HasMisconfiguredPrefixes = true
}

func isInAccessList(network string, accessLists map[string]frrProto.AccessListAnalyzer) bool {
	ip := strings.Split(network, "/")[0]

	for _, acl := range accessLists {
		for _, entry := range acl.AclEntry {
			if entry.IsPermit && entry.IPAddress == ip {
				return true
			}
		}
	}

	return false
}

func normalizeNetworkAddress(address string) string {
	return strings.ToLower(strings.TrimSpace(address))
}

func NssaExternalAnomalyAnalysis(accessList map[string]frrProto.AccessListAnalyzer, shouldState *frrProto.InterAreaLsa, isState *frrProto.InterAreaLsa) {

	//fmt.Println(accessList)
	//fmt.Println(shouldState)
	//fmt.Println(isState)

}

// TODO: Find a fix for point-to-point. The representation of p2p in frr is unclear to me. Thus it's removed from any testing
func checkAdvertisement(accessList map[string]frrProto.AccessListAnalyzer, shouldState *frrProto.IntraAreaLsa, isState *frrProto.IntraAreaLsa) (bool, bool, bool) {

	overAdvertised := false
	underAdvertised := false
	duplicateAdvertised := false
	shouldLsaPrefixes := []string{}
	for _, lsa := range shouldState.Areas {
		for _, link := range lsa.Links {
			if link.LinkType != "point-to-point" {
				shouldLsaPrefixes = append(shouldLsaPrefixes, link.InterfaceAddress)
			}
		}
	}

	isLsaPrefixes := []string{}
	for _, lsa := range isState.Areas {
		for _, link := range lsa.Links {
			if link.LinkType != "point-to-point" {
				isLsaPrefixes = append(isLsaPrefixes, link.InterfaceAddress)
			}
		}
	}

	// Check for Overadvertisement of prefixes

	isOveradvertisedMap := make(map[string]bool)
	shouldPrefixMap := make(map[string]bool)

	// Convert shouldLsaPrefixes to a map for O(1) lookups
	for _, shouldPrefix := range shouldLsaPrefixes {
		shouldPrefixMap[shouldPrefix] = true
	}

	// Check if prefix is NOT in shouldLsaPrefixes
	for _, prefix := range isLsaPrefixes {
		if !shouldPrefixMap[prefix] {
			isOveradvertisedMap[prefix] = true
			overAdvertised = true
		}
	}

	// Check for Underadvertisement of prefixes
	isUnderdvertisedMap := make(map[string]bool)
	isPrefixMap := make(map[string]bool)

	// Convert shouldLsaPrefixes to a map for O(1) lookups
	for _, isPrefix := range isLsaPrefixes {
		isPrefixMap[isPrefix] = true
	}

	// Check if prefix is NOT in shouldLsaPrefixes
	for _, prefix := range shouldLsaPrefixes {
		if !isPrefixMap[prefix] {
			isUnderdvertisedMap[prefix] = true
			underAdvertised = true
		}
	}

	return overAdvertised, underAdvertised, duplicateAdvertised
}
