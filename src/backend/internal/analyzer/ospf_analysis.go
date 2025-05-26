package analyzer

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
)

func (a *Analyzer) AnomalyAnalysisFIB(fibMap map[string]*frrProto.RibPrefixes, receivedNetworkLSDB *frrProto.IntraAreaLsa, receivedSummaryLSDB *frrProto.InterAreaLsa, receivedExternalLSDB *frrProto.InterAreaLsa, receivedNssaExternalLSDB *frrProto.InterAreaLsa) {
	a.Logger.Debug("Starting FIB-LSDB consistency analysis")
	start := time.Now()

	a.Logger.WithAttrs(map[string]interface{}{
		"fib_entries":   len(fibMap),
		"network_lsas":  len(getLSDBMapAndList(receivedNetworkLSDB)),
		"summary_lsas":  len(getLSDBMapAndList(receivedSummaryLSDB)),
		"external_lsas": len(getLSDBMapAndList(receivedExternalLSDB)),
		"nssa_lsas":     len(getLSDBMapAndList(receivedNssaExternalLSDB)),
	}).Debug("Analysis input counts")

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

	if len(result.MissingEntries) > 0 {
		missingExamples := make([]string, 0, 3)
		for i, entry := range result.MissingEntries {
			// Limit to 3 because of storage
			if i >= 3 {
				break
			}
			missingExamples = append(missingExamples,
				fmt.Sprintf("%s/%s", entry.LinkStateId, entry.PrefixLength))
		}

		a.AnomalyLogger.WithAttrs(map[string]interface{}{
			"missing_count":    len(result.MissingEntries),
			"missing_examples": missingExamples,
			"analysis":         "LSDB contains prefixes not found in FIB",
		}).Warning("Found prefixes in LSDB missing from FIB")
	}
	a.Logger.WithAttrs(map[string]interface{}{
		"duration":      time.Since(start).String(),
		"missing_count": len(result.MissingEntries),
	}).Debug("Completed FIB-LSDB analysis")

}

func (a *Analyzer) RouterAnomalyAnalysisLSDB(accessList map[string]*frrProto.AccessListAnalyzer, shouldState *frrProto.IntraAreaLsa, isState *frrProto.IntraAreaLsa) (map[string]*frrProto.Advertisement, map[string]*frrProto.Advertisement) {
	a.Logger.Debug("Starting router LSDB analysis")
	start := time.Now()

	if isState == nil || shouldState == nil {
		a.Logger.Debug("Skipping router analysis - missing input data")
		return nil, nil
	}

	result := &frrProto.AnomalyDetection{
		SuperfluousEntries: []*frrProto.Advertisement{},
		MissingEntries:     []*frrProto.Advertisement{},
		DuplicateEntries:   []*frrProto.Advertisement{},
	}

	isStateCounter := make(map[string]int)
	shouldStateMap := getLsdbStateMap(shouldState)
	isStateMap := getLsdbStateMap(isState)

	for key, shouldLink := range shouldStateMap {
		if shouldLink.LinkType == strings.ToLower("unknown") {
			prefixLength := "/" + shouldLink.PrefixLength
			_, isTransitWithPrefix := isStateMap[shouldLink.LinkStateId+prefixLength]
			_, isTransit := isStateMap[shouldLink.LinkStateId]
			_, isStubWithPrefix := isStateMap[shouldLink.InterfaceAddress+prefixLength]
			_, isStub := isStateMap[shouldLink.InterfaceAddress]
			if !(isTransit || isTransitWithPrefix) && !(isStub || isStubWithPrefix) {
				result.MissingEntries = append(result.MissingEntries, shouldLink)
			}
		} else {
			if _, exists := isStateMap[key]; !exists {
				result.MissingEntries = append(result.MissingEntries, shouldLink)
			}
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

	if len(result.MissingEntries) > 0 {
		missingExamples := make([]map[string]interface{}, 0, 3)
		for i, entry := range result.MissingEntries {
			if i >= 3 {
				break
			}
			missingExamples = append(missingExamples, map[string]interface{}{
				"address": entry.InterfaceAddress,
				"type":    entry.LinkType,
			})
		}

		a.AnomalyLogger.WithAttrs(map[string]interface{}{
			"type":             "router",
			"count":            len(result.MissingEntries),
			"missing_examples": missingExamples,
			"analysis":         "Expected router links not found in operational state",
		}).Warning("Missing router LSAs detected")
	}

	if len(result.SuperfluousEntries) > 0 {
		extraExamples := make([]map[string]interface{}, 0, 3)
		for i, entry := range result.SuperfluousEntries {
			if i >= 3 {
				break
			}
			extraExamples = append(extraExamples, map[string]interface{}{
				"address": entry.InterfaceAddress,
				"type":    entry.LinkType,
			})
		}

		a.AnomalyLogger.WithAttrs(map[string]interface{}{
			"type":           "router",
			"count":          len(result.SuperfluousEntries),
			"extra_examples": extraExamples,
			"analysis":       "Unexpected router links found in operational state",
		}).Warning("Over-advertised router LSAs detected")
	}

	a.Logger.WithAttrs(map[string]interface{}{
		"duration":       time.Since(start).String(),
		"areas_analyzed": len(isState.Areas),
		"missing":        len(result.MissingEntries),
		"extra":          len(result.SuperfluousEntries),
		"duplicates":     len(result.DuplicateEntries),
	}).Info("Completed router LSDB analysis")

	a.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes = len(result.SuperfluousEntries) > 0
	a.AnalysisResult.RouterAnomaly.HasUnAdvertisedPrefixes = len(result.MissingEntries) > 0
	a.AnalysisResult.RouterAnomaly.MissingEntries = result.MissingEntries
	a.AnalysisResult.RouterAnomaly.SuperfluousEntries = result.SuperfluousEntries
	return isStateMap, shouldStateMap
}

func (a *Analyzer) ExternalAnomalyAnalysisLSDB(shouldState *frrProto.InterAreaLsa, isState *frrProto.InterAreaLsa) {
	a.Logger.Debug("Starting external LSDB analysis")
	start := time.Now()

	if isState == nil || shouldState == nil {
		a.Logger.Debug("Skipping external analysis - missing input data")
		return
	}

	result := &frrProto.AnomalyDetection{
		SuperfluousEntries: []*frrProto.Advertisement{},
		MissingEntries:     []*frrProto.Advertisement{},
		DuplicateEntries:   []*frrProto.Advertisement{},
	}

	isStateCounter := make(map[string]int)
	lsdbIsStateMap := getLsdbStateMap(isState)
	lsdbShouldStateMap := getLsdbStateMap(shouldState)

	for key, shouldLink := range lsdbShouldStateMap {
		if _, exists := lsdbIsStateMap[key]; !exists {
			result.MissingEntries = append(result.MissingEntries, shouldLink)
		}
	}

	for key, isLink := range lsdbIsStateMap {
		if _, exists := lsdbShouldStateMap[key]; !exists {
			result.SuperfluousEntries = append(result.SuperfluousEntries, isLink)
		}
	}

	for prefix, counter := range isStateCounter {
		if counter > 1 {
			result.DuplicateEntries = append(result.DuplicateEntries, lsdbIsStateMap[prefix])
		}
	}

	if len(result.MissingEntries) > 0 {
		missingPrefixes := make([]string, 0, 3)
		for i, entry := range result.MissingEntries {
			if i >= 3 {
				break
			}
			missingPrefixes = append(missingPrefixes,
				fmt.Sprintf("%s/%s", entry.LinkStateId, entry.PrefixLength))
		}

		a.AnomalyLogger.WithAttrs(map[string]interface{}{
			"type":             "external",
			"count":            len(result.MissingEntries),
			"missing_prefixes": missingPrefixes,
			"analysis":         "Expected external prefixes not advertised",
		}).Warning("Missing external LSAs detected")
	}

	if len(result.SuperfluousEntries) > 0 {
		extraPrefixes := make([]string, 0, 3)
		for i, entry := range result.SuperfluousEntries {
			if i >= 3 {
				break
			}
			extraPrefixes = append(extraPrefixes,
				fmt.Sprintf("%s/%s", entry.LinkStateId, entry.PrefixLength))
		}

		a.AnomalyLogger.WithAttrs(map[string]interface{}{
			"type":           "external",
			"count":          len(result.SuperfluousEntries),
			"extra_prefixes": extraPrefixes,
			"analysis":       "Unexpected external prefixes advertised",
		}).Warning("Over-advertised external LSAs detected")
	}

	a.Logger.WithAttrs(map[string]interface{}{
		"duration":       time.Since(start).String(),
		"areas_analyzed": len(isState.Areas),
		"missing":        len(result.MissingEntries),
		"extra":          len(result.SuperfluousEntries),
		"duplicates":     len(result.DuplicateEntries),
	}).Info("Completed external LSDB analysis")

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

func (a *Analyzer) NssaExternalAnomalyAnalysis(accessList map[string]*frrProto.AccessListAnalyzer, shouldState *frrProto.InterAreaLsa, isState *frrProto.InterAreaLsa, externalState *frrProto.InterAreaLsa) {
	a.Logger.Debug("Starting NSSA external analysis")
	start := time.Now()

	if isState == nil || shouldState == nil {
		a.Logger.Debug("Skipping NSSA analysis - missing input data")
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

	if len(result.MissingEntries) > 0 {
		missingDetails := make([]map[string]interface{}, 0, 3)
		for i, entry := range result.MissingEntries {
			if i >= 3 {
				break
			}
			missingDetails = append(missingDetails, map[string]interface{}{
				"prefix": fmt.Sprintf("%s/%s", entry.LinkStateId, entry.PrefixLength),
				"p-bit":  entry.PBit,
			})
		}

		a.AnomalyLogger.WithAttrs(map[string]interface{}{
			"type":           "nssa",
			"count":          len(result.MissingEntries),
			"missing_routes": missingDetails,
			"analysis":       "Expected NSSA external routes not advertised",
		}).Warning("Missing NSSA external LSAs detected")
	}

	if len(result.SuperfluousEntries) > 0 {
		extraDetails := make([]map[string]interface{}, 0, 3)
		for i, entry := range result.SuperfluousEntries {
			if i >= 3 {
				break
			}
			extraDetails = append(extraDetails, map[string]interface{}{
				"prefix": fmt.Sprintf("%s/%s", entry.LinkStateId, entry.PrefixLength),
			})
		}

		a.AnomalyLogger.WithAttrs(map[string]interface{}{
			"type":         "nssa",
			"count":        len(result.SuperfluousEntries),
			"extra_routes": extraDetails,
			"analysis":     "Unexpected NSSA external routes advertised",
		}).Warning("Over-advertised NSSA external LSAs detected")
	}

	// P-bit validation
	a.checkNssaPBitTranslation(isState, externalState, result)

	// Update analysis result
	a.AnalysisResult.NssaExternalAnomaly.HasOverAdvertisedPrefixes = len(result.SuperfluousEntries) > 0
	a.AnalysisResult.NssaExternalAnomaly.HasUnAdvertisedPrefixes = len(result.MissingEntries) > 0
	a.AnalysisResult.NssaExternalAnomaly.HasDuplicatePrefixes = len(result.DuplicateEntries) > 0
	a.AnalysisResult.NssaExternalAnomaly.MissingEntries = result.MissingEntries
	a.AnalysisResult.NssaExternalAnomaly.SuperfluousEntries = result.SuperfluousEntries
	a.AnalysisResult.NssaExternalAnomaly.DuplicateEntries = result.DuplicateEntries

	a.Logger.WithAttrs(map[string]interface{}{
		"duration":   time.Since(start).String(),
		"nssa_areas": len(isState.Areas),
		"missing":    len(result.MissingEntries),
		"extra":      len(result.SuperfluousEntries),
	}).Info("Completed NSSA external analysis")
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

func getAdvertisementKey(adv *frrProto.Advertisement) string {
	if adv.LinkType == "transit network" {
		return normalizeNetworkAddress(adv.InterfaceAddress)
	} else if strings.Contains(strings.ToLower(adv.LinkType), "virtual link") {
		return adv.InterfaceAddress + "/32"
	}
	return getKeyWithFallback(adv.InterfaceAddress, adv.LinkStateId, adv.PrefixLength)
}

func getKeyWithFallback(primary, fallback, prefixLength string) string {
	addr := normalizeNetworkAddress(primary)
	if addr == "" {
		addr = normalizeNetworkAddress(fallback)
	}
	return fmt.Sprintf("%s/%s", addr, normalizePrefixLength(prefixLength))
}

func normalizePrefixLength(prefixLength string) string {
	if prefixLength == "" {
		return "32"
	}
	i, err := strconv.Atoi(prefixLength)
	if err != nil {
		return "32"
	}
	return strconv.Itoa(i)
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

func getLsdbStateMap(lsdbState interface{}) map[string]*frrProto.Advertisement {

	result := make(map[string]*frrProto.Advertisement)
	var areas []*frrProto.AreaAnalyzer

	switch area := lsdbState.(type) {
	case *frrProto.IntraAreaLsa:
		areas = area.Areas
	case *frrProto.InterAreaLsa:
		areas = area.Areas
	}

	for _, area := range areas {
		for i := range area.Links {
			link := area.Links[i]
			key := getAdvertisementKey(link)
			adv := &frrProto.Advertisement{
				InterfaceAddress: link.InterfaceAddress,
				LinkStateId:      link.LinkStateId,
				LinkType:         link.LinkType,
				PrefixLength:     link.PrefixLength,
			}
			if strings.ToLower(link.LinkType) == "unknown" {
				prefixLength := "/" + link.PrefixLength
				result[link.InterfaceAddress+prefixLength] = adv
				result[link.LinkStateId+prefixLength] = adv
				result[link.InterfaceAddress] = adv
				result[link.LinkStateId] = adv
			} else {
				result[key] = adv
			}
		}
	}

	return result

}
