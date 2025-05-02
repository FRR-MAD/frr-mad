package analyzer

import (
	"fmt"
	"strings"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	// "google.golang.org/protobuf/proto"
)

func (a *Analyzer) routerAnomalyAnalysis(accessList map[string]accessList, shouldState *intraAreaLsa, isState *intraAreaLsa) {
	if isState == nil || shouldState == nil {
		fmt.Println("nil!")
		return
	}

	result := &frrProto.AnomalyDetection{
		HasDuplicatePrefixes:     false,
		HasMisconfiguredPrefixes: false,
		SuperfluousEntries:       []*frrProto.Advertisement{},
		MissingEntries:           []*frrProto.Advertisement{},
	}

	isStateMap := make(map[string]*frrProto.Advertisement)
	shouldStateMap := make(map[string]*frrProto.Advertisement)

	for _, area := range isState.Areas {
		for i := range area.Links {
			link := &area.Links[i]
			key := getadvertisementKey(link)
			isStateMap[key] = link
		}
	}

	for _, area := range shouldState.Areas {
		for i := range area.Links {
			link := &area.Links[i]
			key := getadvertisementKey(link)
			shouldStateMap[key] = link
		}
	}

	for key, shouldLink := range shouldStateMap {
		if _, exists := isStateMap[key]; !exists {
			if !isExcludedByAccessList(shouldLink, accessList) {
				result.MissingEntries = append(result.MissingEntries, shouldLink)
			}
		}
	}

	for key, isLink := range isStateMap {
		if _, exists := shouldStateMap[key]; !exists {
			result.SuperfluousEntries = append(result.SuperfluousEntries, isLink)
		}
	}

	//a.AnalysisResult.RouterAnomaly.HasUnderAdvertisedPrefixes = writeBoolTarget(result.HasUnderAdvertisedPrefixes)
	//a.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes = writeBoolTarget(result.HasOverAdvertisedPrefixes)

	a.AnalysisResult.RouterAnomaly.HasUnderAdvertisedPrefixes = len(result.MissingEntries) > 0
	a.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes = len(result.SuperfluousEntries) > 0
	a.AnalysisResult.RouterAnomaly.HasDuplicatePrefixes = writeBoolTarget(result.HasDuplicatePrefixes)
	a.AnalysisResult.RouterAnomaly.HasMisconfiguredPrefixes = writeBoolTarget(result.HasMisconfiguredPrefixes)
	a.AnalysisResult.RouterAnomaly.MissingEntries = result.MissingEntries
	a.AnalysisResult.RouterAnomaly.SuperfluousEntries = result.SuperfluousEntries

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

func isExcludedByAccessList(adv *frrProto.Advertisement, accessLists map[string]accessList) bool {
	for _, acl := range accessLists {
		for _, entry := range acl.aclEntry {
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

func externalAnomalyAnalysis(accessList map[string]accessList, isState *interAreaLsa, shouldState *interAreaLsa) {
}

// func (a *Analyzer) externalAnomalyAnalysis(accessList map[string]accessList, isState *interAreaLsa, shouldState *interAreaLsa) {

func (a *Analyzer) externalAnomalyAnalysis(accessList map[string]accessList, isState *interAreaLsa, shouldState *interAreaLsa) {
	result := frrProto.AnomalyDetection{}

	//result := &frrProto.AnomalyAnalysis{
	//	AnomalyDetection: &frrProto.AnomalyDetection{},
	//	MissingEntries:   []*frrProto.Advertisement{},
	//	SuperfluousEntries:     []*frrProto.Advertisement{},
	//}

	if isState == nil || shouldState == nil {
		return
	}

	isStateMap := make(map[string]*frrProto.Advertisement)
	shouldStateMap := make(map[string]*frrProto.Advertisement)

	for _, area := range isState.Areas {
		for i := range area.Links {
			link := &area.Links[i]
			key := normalizeNetworkAddress(link.LinkStateId)
			isStateMap[key] = link
		}
	}

	for _, area := range shouldState.Areas {
		for i := range area.Links {
			link := &area.Links[i]
			key := normalizeNetworkAddress(link.LinkStateId)
			shouldStateMap[key] = link
		}
	}

	accessListNetworks := make(map[string]bool)
	for _, acl := range accessList {
		for _, entry := range acl.aclEntry {
			if entry.IsPermit {
				network := fmt.Sprintf("%s/%d", entry.IPAddress, entry.PrefixLength)
				accessListNetworks[normalizeNetworkAddress(network)] = true
			}
		}
	}

	for key, shouldLink := range shouldStateMap {
		if _, exists := isStateMap[key]; !exists {
			networkWithPrefix := fmt.Sprintf("%s/%s", shouldLink.LinkStateId, shouldLink.PrefixLength)
			normalizedNetwork := normalizeNetworkAddress(networkWithPrefix)

			if accessListNetworks[normalizedNetwork] || isInAccessList(shouldLink.LinkStateId, accessList) {
				result.MissingEntries = append(result.MissingEntries, shouldLink)
			}
		}
	}

	result.HasUnderAdvertisedPrefixes = len(result.MissingEntries) > 0
	result.HasOverAdvertisedPrefixes = false
	result.HasDuplicatePrefixes = false
	result.HasMisconfiguredPrefixes = true
}

func isInAccessList(network string, accessLists map[string]accessList) bool {
	ip := strings.Split(network, "/")[0]

	for _, acl := range accessLists {
		for _, entry := range acl.aclEntry {
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

func nssaExternalAnomalyAnalysis(accessList map[string]accessList, shouldState *interAreaLsa, isState *interAreaLsa) {

	//fmt.Println(accessList)
	//fmt.Println(shouldState)
	//fmt.Println(isState)

}
