package analyzer

import (
	"sort"
	"strings"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

func (a *Analyzer) AnomalyAnalysisFIB(fibMap map[string]frrProto.RibPrefixes, receivedNetworkLSDB *frrProto.IntraAreaLsa, receivedSummaryLSDB *frrProto.InterAreaLsa, receivedExternalLSDB *frrProto.InterAreaLsa, receivedNssaExternalLSDB *frrProto.InterAreaLsa) {
	result := &frrProto.AnomalyDetection{
		SuperfluousEntries: []*frrProto.Advertisement{},
		MissingEntries:     []*frrProto.Advertisement{},
		DuplicateEntries:   []*frrProto.Advertisement{},
	}

	lsdbList := []string{}
	networkLsdbList := getLSDBMapAndList(receivedSummaryLSDB)
	summaryLsdbList := getLSDBMapAndList(receivedSummaryLSDB)
	externalLsdbList := getLSDBMapAndList(receivedExternalLSDB)
	nssaExternalLsdbList := getLSDBMapAndList(receivedNssaExternalLSDB)

	lsdbList = append(lsdbList, networkLsdbList...)
	lsdbList = append(lsdbList, summaryLsdbList...)
	lsdbList = append(lsdbList, externalLsdbList...)
	lsdbList = append(lsdbList, nssaExternalLsdbList...)
	lsdbList = filterUnique(lsdbList)

	setA := make(map[string]bool)
	for _, v := range fibMap {
		setA[v.Prefix] = true
	}
	fibList := []string{}
	for prefix, _ := range fibMap {
		fibList = append(fibList, prefix)
	}
	sort.Strings(fibList)
	sort.Strings(lsdbList)

	for _, entry := range lsdbList {
		_, exists := fibMap[entry]
		if !exists {
			result.HasUnAdvertisedPrefixes = true
			result.MissingEntries = append(result.MissingEntries, &frrProto.Advertisement{
				LinkStateId:  strings.Split(entry, "/")[0],
				PrefixLength: strings.Split(entry, "/")[1],
			})
		}
	}

	a.AnalysisResult.LsdbToRibAnomaly.HasUnAdvertisedPrefixes = result.HasUnAdvertisedPrefixes
	a.AnalysisResult.LsdbToRibAnomaly.MissingEntries = result.MissingEntries

}

func (a *Analyzer) RouterAnomalyAnalysisLSDB(accessList map[string]*frrProto.AccessListAnalyzer, shouldState *frrProto.IntraAreaLsa, isState *frrProto.IntraAreaLsa) (map[string]*frrProto.Advertisement, map[string]*frrProto.Advertisement) {
	if isState == nil || shouldState == nil {
		return nil, nil
	}

	result := &frrProto.AnomalyDetection{
		SuperfluousEntries: []*frrProto.Advertisement{},
		MissingEntries:     []*frrProto.Advertisement{},
		DuplicateEntries:   []*frrProto.Advertisement{},
	}

	isStateMap := make(map[string]*frrProto.Advertisement)
	isStateCounter := make(map[string]int)
	shouldStateMap := make(map[string]*frrProto.Advertisement)

	for _, area := range isState.Areas {
		for i := range area.Links {
			link := area.Links[i]
			key := getAdvertisementKey(link)
			isStateMap[key] = &frrProto.Advertisement{
				InterfaceAddress: link.InterfaceAddress,
				LinkType:         link.LinkType,
			}
			if link.LinkType == strings.ToLower("Stub Network") {
				isStateMap[key].PrefixLength = link.PrefixLength
			}
		}
	}

	for _, area := range shouldState.Areas {
		for i := range area.Links {
			link := area.Links[i]
			key := getAdvertisementKey(link)
			adv := &frrProto.Advertisement{
				InterfaceAddress: link.InterfaceAddress,
				LinkType:         link.LinkType,
			}
			if link.LinkType == strings.ToLower("Stub Network") {
				shouldStateMap[key] = adv
				shouldStateMap[key].PrefixLength = link.PrefixLength
			} else if link.LinkType == strings.ToLower("unknown") {
				adv.LinkStateId = link.LinkStateId
				shouldStateMap[link.InterfaceAddress] = adv
				shouldStateMap[link.LinkStateId] = adv
				shouldStateMap[link.InterfaceAddress].PrefixLength = link.PrefixLength
				shouldStateMap[link.LinkStateId].PrefixLength = link.PrefixLength
			} else {
				shouldStateMap[key] = adv
			}
		}
	}

	// unadvertised
	// is in shouldSate, but not in isState
	// loop through shouldState and check with is State
	for key, shouldLink := range shouldStateMap {
		isMissing := false
		if shouldLink.LinkType == strings.ToLower("unknown") {
			_, isTransit := isStateMap[shouldLink.LinkStateId]
			_, isStub := isStateMap[shouldLink.InterfaceAddress]
			if !isTransit && !isStub {
				isMissing = true
			}
		} else {
			if _, exists := isStateMap[key]; !exists {
				isMissing = true
			}
		}
		if isMissing {
			result.MissingEntries = append(result.MissingEntries, shouldLink)
		}
	}

	// overadvertised
	// is in isState, but not in shouldState
	// loop through isState and check with shouldState
	// because I don't know what kind of link type shouldState is (either transit or stub) I have to check both these networks differently
	for key, isLink := range isStateMap {
		if _, exists := shouldStateMap[key]; !exists {
			result.SuperfluousEntries = append(result.SuperfluousEntries, isLink)
		}
	}

	for prefix, counter := range isStateCounter {
		if counter > 1 {
			result.DuplicateEntries = append(result.DuplicateEntries, isStateMap[prefix])
		}
	}

	a.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes = len(result.SuperfluousEntries) > 0
	a.AnalysisResult.RouterAnomaly.HasUnAdvertisedPrefixes = len(result.MissingEntries) > 0
	a.AnalysisResult.RouterAnomaly.MissingEntries = result.MissingEntries
	a.AnalysisResult.RouterAnomaly.SuperfluousEntries = result.SuperfluousEntries
	return isStateMap, shouldStateMap

}

func getAdvertisementKey(adv *frrProto.Advertisement) string {
	if adv.LinkType == "transit network" {
		return adv.InterfaceAddress
	}
	return getKeyWithFallback(adv.InterfaceAddress, adv.LinkStateId, adv.PrefixLength)
}

func getKeyWithFallback(primary, fallback, prefixLength string) string {
	addr := normalizeNetworkAddress(primary)
	if addr == "" {
		addr = fallback
	}
	return addr
}

func isExcludedByAccessList(adv *frrProto.Advertisement, accessLists map[string]*frrProto.AccessListAnalyzer) bool {
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

func (a *Analyzer) ExternalAnomalyAnalysisLSDB(shouldState *frrProto.InterAreaLsa, isState *frrProto.InterAreaLsa) {
	if isState == nil || shouldState == nil {
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

	for _, area := range isState.Areas {
		for i := range area.Links {
			link := area.Links[i]
			key := getAdvertisementKey(link)
			isStateMap[key] = &frrProto.Advertisement{
				InterfaceAddress: link.InterfaceAddress,
				LinkStateId:      link.LinkStateId,
				LinkType:         link.LinkType,
				PrefixLength:     link.PrefixLength,
			}

		}
	}

	for _, area := range shouldState.Areas {
		for i := range area.Links {
			link := area.Links[i]
			key := getAdvertisementKey(link)
			shouldStateMap[key] = &frrProto.Advertisement{
				InterfaceAddress: link.InterfaceAddress,
				LinkStateId:      link.LinkStateId,
				LinkType:         link.LinkType,
				PrefixLength:     link.PrefixLength,
			}

		}
	}

	for key, shouldLink := range shouldStateMap {
		if _, exists := isStateMap[key]; !exists {
			result.MissingEntries = append(result.MissingEntries, shouldLink)
		}
	}

	for key, isLink := range isStateMap {
		if _, exists := shouldStateMap[key]; !exists {
			result.SuperfluousEntries = append(result.SuperfluousEntries, isLink)
		}
	}

	for prefix, counter := range isStateCounter {
		if counter > 1 {
			result.DuplicateEntries = append(result.DuplicateEntries, isStateMap[prefix])
		}
	}

	a.AnalysisResult.ExternalAnomaly.HasOverAdvertisedPrefixes = len(result.SuperfluousEntries) > 0
	a.AnalysisResult.ExternalAnomaly.HasUnAdvertisedPrefixes = len(result.MissingEntries) > 0
	a.AnalysisResult.ExternalAnomaly.HasDuplicatePrefixes = len(result.DuplicateEntries) > 0
	a.AnalysisResult.ExternalAnomaly.MissingEntries = result.MissingEntries
	a.AnalysisResult.ExternalAnomaly.SuperfluousEntries = result.SuperfluousEntries
	a.AnalysisResult.ExternalAnomaly.DuplicateEntries = result.DuplicateEntries

}

func normalizeNetworkAddress(address string) string {
	return strings.ToLower(strings.TrimSpace(address))
}

// TODO: Add missing analysis
func (a *Analyzer) NssaExternalAnomalyAnalysis(accessList map[string]*frrProto.AccessListAnalyzer, shouldState *frrProto.InterAreaLsa, isState *frrProto.InterAreaLsa, externalState *frrProto.InterAreaLsa) {
	if isState == nil || shouldState == nil {
		return
	}

	result := &frrProto.AnomalyDetection{
		HasMisconfiguredPrefixes: false,
		SuperfluousEntries:       []*frrProto.Advertisement{},
		MissingEntries:           []*frrProto.Advertisement{},
		DuplicateEntries:         []*frrProto.Advertisement{},
	}

	isStateMap := make(map[string]map[string]*frrProto.Advertisement)
	shouldStateMap := make(map[string]map[string]*frrProto.Advertisement)
	duplicateTracker := make(map[string]map[string]int)

	for _, area := range isState.Areas {
		if area.LsaType != "NSSA-LSA" {
			continue
		}

		if isStateMap[area.AreaName] == nil {
			isStateMap[area.AreaName] = make(map[string]*frrProto.Advertisement)
		}
		if duplicateTracker[area.AreaName] == nil {
			duplicateTracker[area.AreaName] = make(map[string]int)
		}

		for _, link := range area.Links {
			key := getAdvertisementKey(link)
			isStateMap[area.AreaName][key] = link

			duplicateTracker[area.AreaName][key]++
			if duplicateTracker[area.AreaName][key] > 1 {
				result.DuplicateEntries = append(result.DuplicateEntries, link)
			}
		}
	}

	for _, area := range shouldState.Areas {
		if area.LsaType != "NSSA-LSA" {
			continue
		}

		if shouldStateMap[area.AreaName] == nil {
			shouldStateMap[area.AreaName] = make(map[string]*frrProto.Advertisement)
		}

		for _, link := range area.Links {
			key := getAdvertisementKey(link)
			shouldStateMap[area.AreaName][key] = link
		}
	}

	for areaName, shouldRoutes := range shouldStateMap {
		for key, route := range shouldRoutes {
			if isStateMap[areaName] == nil || isStateMap[areaName][key] == nil {
				if !isExcludedByAccessList(route, accessList) {
					result.MissingEntries = append(result.MissingEntries, route)
				}
			}
		}
	}

	for areaName, isRoutes := range isStateMap {
		for key, route := range isRoutes {
			if shouldStateMap[areaName] == nil || shouldStateMap[areaName][key] == nil {
				result.SuperfluousEntries = append(result.SuperfluousEntries, route)
			}
		}
	}

	a.checkNssaPBitTranslation(isState, externalState, result)

	a.AnalysisResult.NssaExternalAnomaly.HasOverAdvertisedPrefixes = len(result.SuperfluousEntries) > 0
	a.AnalysisResult.NssaExternalAnomaly.HasUnAdvertisedPrefixes = len(result.MissingEntries) > 0
	a.AnalysisResult.NssaExternalAnomaly.HasDuplicatePrefixes = len(result.DuplicateEntries) > 0
	a.AnalysisResult.NssaExternalAnomaly.MissingEntries = result.MissingEntries
	a.AnalysisResult.NssaExternalAnomaly.SuperfluousEntries = result.SuperfluousEntries
	a.AnalysisResult.NssaExternalAnomaly.DuplicateEntries = result.DuplicateEntries
}

func getLSDBMapAndList(lsdb interface{}) []string {
	lsdbList := []string{}

	switch db := lsdb.(type) {
	case *frrProto.IntraAreaLsa:
		for _, area := range db.Areas {
			for _, lsa := range area.Links {
				lsdbList = append(lsdbList, lsa.LinkStateId+"/"+lsa.PrefixLength)
			}
		}

	case *frrProto.InterAreaLsa:
		for _, area := range db.Areas {
			for _, lsa := range area.Links {
				lsdbList = append(lsdbList, lsa.LinkStateId+"/"+lsa.PrefixLength)
			}
		}
	}
	return filterUnique(lsdbList)
}

func filterUnique(lsaList []string) []string {
	uniqueMap := make(map[string]bool)
	var uniqueNames []string

	for _, lsa := range lsaList {
		if _, exists := uniqueMap[lsa]; !exists {
			uniqueMap[lsa] = true
			uniqueNames = append(uniqueNames, lsa)
		}
	}
	return uniqueNames
}

func (a *Analyzer) checkNssaPBitTranslation(nssaState *frrProto.InterAreaLsa, externalState *frrProto.InterAreaLsa, result *frrProto.AnomalyDetection) {
	if externalState == nil {
		return
	}

	externalMap := make(map[string]bool)
	for _, area := range externalState.Areas {
		for _, link := range area.Links {
			key := getAdvertisementKey(link)
			externalMap[key] = true
		}
	}

	for _, area := range nssaState.Areas {
		if area.LsaType != "NSSA-LSA" {
			continue
		}
		for _, link := range area.Links {
			if link.PBit {
				key := getAdvertisementKey(link)
				if !externalMap[key] {
					// P-bit set, but no matching Type-5 exists
					result.MissingEntries = append(result.MissingEntries, link)
				}
			}
		}
	}
}
