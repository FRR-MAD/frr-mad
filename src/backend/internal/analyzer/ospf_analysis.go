package analyzer

import (
	"fmt"
	"strings"
)

type AnomalyAnalysis struct {
	Type             string            `json:"type,omitempty"`
	AnomalyDetection *AnomalyDetection `json:"anomaly_detection,omitempty"`
	MissingEntries   []advertisment    `json:"missing_entries,omitempty"`
	ExtraEntries     []advertisment    `json:"extra_entries,omitempty"`
}

type AnomalyDetection struct {
	OverAdvertised  bool `json:"over_advertised,omitempty"`
	UnderAdvertised bool `json:"under_advertised,omitempty"`
}

func routerAnomalyAnalysis(accessList map[string]accessList, shouldState *intraAreaLsa, isState *intraAreaLsa) *AnomalyAnalysis {
	result := &AnomalyAnalysis{
		Type:             "router",
		AnomalyDetection: &AnomalyDetection{},
		MissingEntries:   []advertisment{},
		ExtraEntries:     []advertisment{},
	}

	if isState == nil || shouldState == nil {
		return result
	}

	isStateMap := make(map[string]advertisment)
	shouldStateMap := make(map[string]advertisment)

	for _, area := range isState.Areas {
		for _, link := range area.Links {
			key := getAdvertismentKey(link)
			isStateMap[key] = link
		}
	}

	for _, area := range shouldState.Areas {
		for _, link := range area.Links {
			key := getAdvertismentKey(link)
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
			result.ExtraEntries = append(result.ExtraEntries, isLink)
		}
	}

	result.AnomalyDetection.UnderAdvertised = len(result.MissingEntries) > 0
	result.AnomalyDetection.OverAdvertised = len(result.ExtraEntries) > 0

	fmt.Println()
	fmt.Println()
	fmt.Println("############### Router Anomaly Detection ###############")
	fmt.Println(result)
	fmt.Printf("%+v\n", result.AnomalyDetection)
	//fmt.Println(shouldState)
	//fmt.Println(isState)
	fmt.Println()
	fmt.Println()
	return result
}

func getAdvertismentKey(adv advertisment) string {
	if adv.InterfaceAddress != "" {
		return normalizeNetworkAddress(adv.InterfaceAddress)
	}
	return normalizeNetworkAddress(adv.LinkStateId)
}

func isExcludedByAccessList(adv advertisment, accessLists map[string]accessList) bool {
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

func externalAnomalyAnalysis(accessList map[string]accessList, isState *interAreaLsa, shouldState *interAreaLsa) *AnomalyAnalysis {
	result := &AnomalyAnalysis{
		Type:             "external",
		AnomalyDetection: &AnomalyDetection{},
		MissingEntries:   []advertisment{},
		ExtraEntries:     []advertisment{},
	}

	if isState == nil || shouldState == nil {
		return result
	}

	isStateMap := make(map[string]advertisment)
	shouldStateMap := make(map[string]advertisment)

	for _, area := range isState.Areas {
		for _, link := range area.Links {
			key := normalizeNetworkAddress(link.LinkStateId)
			isStateMap[key] = link
		}
	}

	for _, area := range shouldState.Areas {
		for _, link := range area.Links {
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

	result.AnomalyDetection.UnderAdvertised = len(result.MissingEntries) > 0

	fmt.Println()
	fmt.Println()
	fmt.Println("############### Router Anomaly Detection ###############")
	fmt.Println(result)
	fmt.Printf("%+v\n", result.AnomalyDetection)
	//fmt.Println(shouldState)
	//fmt.Println(isState)
	fmt.Println()
	fmt.Println()
	return result
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

	fmt.Println(accessList)
	fmt.Println(shouldState)
	fmt.Println(isState)

}
