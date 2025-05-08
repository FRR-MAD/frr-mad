package analyzer

import (
	"fmt"
	"strings"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

// WIP
// TODO: Implement FIB analysis
func (a *Analyzer) AnomalyAnalysisFIB(ribMap map[string]frrProto.RibPrefixes, isStateRouter *frrProto.IntraAreaLsa, isStateExternal *frrProto.InterAreaLsa, isStateNssaExternal *frrProto.InterAreaLsa) {
	result := &frrProto.AnomalyDetection{
		HasMisconfiguredPrefixes: false,
		SuperfluousEntries:       []*frrProto.Advertisement{},
		MissingEntries:           []*frrProto.Advertisement{},
		DuplicateEntries:         []*frrProto.Advertisement{},
	}

	ospfCounter := 0
	for _, entry := range ribMap {
		if entry.Protocol == "ospf" {
			ospfCounter += 1
		}
	}

	ospfIsStateRouterCounter := 0
	for _, area := range isStateRouter.Areas {
		ospfIsStateRouterCounter += len(area.Links)
	}

	ospfIsStateExternalCounter := 0
	for _, area := range isStateExternal.Areas {
		ospfIsStateExternalCounter += len(area.Links)
	}

	for _, foo := range isStateRouter.Areas {
		for _, i := range foo.Links {
			//fmt.Println(i.LinkStateId)
			fmt.Println(i.InterfaceAddress)
		}
	}

	// only  LinkStateId contains an entry
	for _, foo := range isStateExternal.Areas {
		for _, i := range foo.Links {
			fmt.Println(i.LinkStateId)
			//fmt.Println(i.InterfaceAddress)
		}
	}

	// only  InterfaceAddress contains an entry
	for _, foo := range isStateNssaExternal.Areas {
		for _, i := range foo.Links {
			//fmt.Println(i.LinkStateId)
			fmt.Println(i.InterfaceAddress)
		}
	}
	a.AnalysisResult.FibAnomaly.HasOverAdvertisedPrefixes = len(result.MissingEntries) > 0
	a.AnalysisResult.FibAnomaly.HasUnderAdvertisedPrefixes = len(result.SuperfluousEntries) > 0
	a.AnalysisResult.FibAnomaly.HasDuplicatePrefixes = len(result.DuplicateEntries) > 0
	a.AnalysisResult.FibAnomaly.MissingEntries = result.MissingEntries
	a.AnalysisResult.FibAnomaly.SuperfluousEntries = result.SuperfluousEntries
	a.AnalysisResult.FibAnomaly.DuplicateEntries = result.DuplicateEntries

}

func (a *Analyzer) RouterAnomalyAnalysisLSDB(accessList map[string]frrProto.AccessListAnalyzer, shouldState *frrProto.IntraAreaLsa, isState *frrProto.IntraAreaLsa) (map[string]*frrProto.Advertisement, map[string]*frrProto.Advertisement) {
	if isState == nil || shouldState == nil {
		//fmt.Println("nil!")
		return nil, nil
	}

	result := &frrProto.AnomalyDetection{
		HasUnderAdvertisedPrefixes: false,
		HasOverAdvertisedPrefixes:  false,
		HasDuplicatePrefixes:       false,
		HasMisconfiguredPrefixes:   false,
		SuperfluousEntries:         []*frrProto.Advertisement{},
		MissingEntries:             []*frrProto.Advertisement{},
		DuplicateEntries:           []*frrProto.Advertisement{},
	}

	isStateMap := make(map[string]*frrProto.Advertisement)
	isStateCounter := make(map[string]int)
	shouldStateMap := make(map[string]*frrProto.Advertisement)

	for _, area := range isState.Areas {
		for i := range area.Links {
			if area.Links[i].LinkType != "point-to-point" {
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
	}

	for _, area := range shouldState.Areas {
		for i := range area.Links {
			if area.Links[i].LinkType != "point-to-point" {
				link := area.Links[i]
				key := getAdvertisementKey(link)
				shouldStateMap[key] = &frrProto.Advertisement{
					InterfaceAddress: link.InterfaceAddress,
					LinkType:         link.LinkType,
				}
				if link.LinkType == strings.ToLower("Stub Network") {
					shouldStateMap[key].PrefixLength = link.PrefixLength
				}

			}
		}
	}

	// check for missing prefixes -> underadvertised
	for key, shouldLink := range shouldStateMap {
		if _, exists := isStateMap[key]; !exists {
			//fmt.Println(shouldLink)
			result.MissingEntries = append(result.MissingEntries, shouldLink)
			//if !isExcludedByAccessList(shouldLink, accessList) {
			//}
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
			result.DuplicateEntries = append(result.DuplicateEntries, isStateMap[prefix])
		}
	}

	a.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes = len(result.SuperfluousEntries) > 0
	a.AnalysisResult.RouterAnomaly.HasUnderAdvertisedPrefixes = len(result.MissingEntries) > 0
	a.AnalysisResult.RouterAnomaly.HasDuplicatePrefixes = len(result.DuplicateEntries) > 0
	a.AnalysisResult.RouterAnomaly.MissingEntries = result.MissingEntries
	a.AnalysisResult.RouterAnomaly.SuperfluousEntries = result.SuperfluousEntries
	a.AnalysisResult.RouterAnomaly.DuplicateEntries = result.DuplicateEntries
	return isStateMap, shouldStateMap

}

func writeBoolTarget(source bool) bool {
	if source {
		return source
	}
	return false
}

func getAdvertisementKey(adv *frrProto.Advertisement) string {
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

// TODO: Check what is received by BGP, cross check with FIB (not rib)
/*
	- if it is in FIB, it's good
	- if it is NOT in FIB, it's no good
*/
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

	// Process isState links
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

	// Process shouldState links
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

	// check for missing prefixes -> underadvertised
	for key, shouldLink := range shouldStateMap {
		if _, exists := isStateMap[key]; !exists {
			result.MissingEntries = append(result.MissingEntries, shouldLink)
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
			result.DuplicateEntries = append(result.DuplicateEntries, isStateMap[prefix])
		}
	}

	a.AnalysisResult.ExternalAnomaly.HasOverAdvertisedPrefixes = len(result.SuperfluousEntries) > 0
	a.AnalysisResult.ExternalAnomaly.HasUnderAdvertisedPrefixes = len(result.MissingEntries) > 0
	a.AnalysisResult.ExternalAnomaly.HasDuplicatePrefixes = len(result.DuplicateEntries) > 0
	a.AnalysisResult.ExternalAnomaly.MissingEntries = result.MissingEntries
	a.AnalysisResult.ExternalAnomaly.SuperfluousEntries = result.SuperfluousEntries
	a.AnalysisResult.ExternalAnomaly.DuplicateEntries = result.DuplicateEntries

}

func normalizeNetworkAddress(address string) string {
	return strings.ToLower(strings.TrimSpace(address))
}

// TODO: Add missing analysis
func (a *Analyzer) NssaExternalAnomalyAnalysis(accessList map[string]frrProto.AccessListAnalyzer, shouldState *frrProto.InterAreaLsa, isState *frrProto.InterAreaLsa) {
	if isState == nil || shouldState == nil {
		return
	}

	result := &frrProto.AnomalyDetection{
		HasMisconfiguredPrefixes: false,
		SuperfluousEntries:       []*frrProto.Advertisement{},
		MissingEntries:           []*frrProto.Advertisement{},
		DuplicateEntries:         []*frrProto.Advertisement{},
	}

	// Maps to track expected and actual NSSA-external routes
	isStateMap := make(map[string]map[string]*frrProto.Advertisement) // area -> prefix -> advertisement
	shouldStateMap := make(map[string]map[string]*frrProto.Advertisement)
	duplicateTracker := make(map[string]map[string]int) // area -> prefix -> count

	// Process actual NSSA-external routes (isState)
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
			key := link.LinkStateId + "/" + link.PrefixLength
			isStateMap[area.AreaName][key] = link

			// Track duplicates
			duplicateTracker[area.AreaName][key]++
			if duplicateTracker[area.AreaName][key] > 1 {
				result.DuplicateEntries = append(result.DuplicateEntries, link)
			}
		}
	}

	// Process expected NSSA-external routes (shouldState)
	for _, area := range shouldState.Areas {
		if area.LsaType != "NSSA-LSA" {
			continue // Skip non-NSSA areas
		}

		if shouldStateMap[area.AreaName] == nil {
			shouldStateMap[area.AreaName] = make(map[string]*frrProto.Advertisement)
		}

		for _, link := range area.Links {
			key := link.LinkStateId + "/" + link.PrefixLength
			shouldStateMap[area.AreaName][key] = link
		}
	}

	// Check for missing routes (under-advertised)
	for areaName, shouldRoutes := range shouldStateMap {
		for key, route := range shouldRoutes {
			// Check if route exists in isState for this area
			if isStateMap[areaName] == nil || isStateMap[areaName][key] == nil {
				// Check if route is excluded by access list
				if !isExcludedByAccessList(route, accessList) {
					result.MissingEntries = append(result.MissingEntries, route)
				}
			}
		}
	}

	// Check for superfluous routes (over-advertised)
	for areaName, isRoutes := range isStateMap {
		for key, route := range isRoutes {
			// Check if route exists in shouldState for this area
			if shouldStateMap[areaName] == nil || shouldStateMap[areaName][key] == nil {
				result.SuperfluousEntries = append(result.SuperfluousEntries, route)
			}
		}
	}

	// Check for P-bit issues (NSSA-external routes not being translated)
	// This requires comparing Type 7 LSAs in NSSA with Type 5 LSAs in backbone
	// You'll need to implement this separately by comparing NSSA and external databases

	// Update analysis result
	a.AnalysisResult.NssaExternalAnomaly.HasOverAdvertisedPrefixes = len(result.SuperfluousEntries) > 0
	a.AnalysisResult.NssaExternalAnomaly.HasUnderAdvertisedPrefixes = len(result.MissingEntries) > 0
	a.AnalysisResult.NssaExternalAnomaly.HasDuplicatePrefixes = len(result.DuplicateEntries) > 0
	a.AnalysisResult.NssaExternalAnomaly.MissingEntries = result.MissingEntries
	a.AnalysisResult.NssaExternalAnomaly.SuperfluousEntries = result.SuperfluousEntries
	a.AnalysisResult.NssaExternalAnomaly.DuplicateEntries = result.DuplicateEntries
}
