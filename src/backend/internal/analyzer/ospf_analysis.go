package analyzer

import "strings"

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

func routerAnomalyAnalysis(accessList map[string]accessList, isState *intraAreaLsa, shouldState *intraAreaLsa) *AnomalyAnalysis {
	result := &AnomalyAnalysis{
		Type:             "router",
		AnomalyDetection: &AnomalyDetection{},
		MissingEntries:   []advertisment{},
		ExtraEntries:     []advertisment{},
	}

	if isState == nil || shouldState == nil {
		return result
	}

	// Create maps to track advertisements by network address for easier comparison
	isStateMap := make(map[string]advertisment)
	shouldStateMap := make(map[string]advertisment)

	// Process the "is state" LSAs
	for _, area := range isState.Areas {
		for _, link := range area.Links {
			// Use interface address or link state ID as the key
			key := getAdvertismentKey(link)
			isStateMap[key] = link
		}
	}

	// Process the "should state" LSAs
	for _, area := range shouldState.Areas {
		for _, link := range area.Links {
			key := getAdvertismentKey(link)
			shouldStateMap[key] = link
		}
	}

	// Check for missing entries (in should but not in is)
	for key, shouldLink := range shouldStateMap {
		if _, exists := isStateMap[key]; !exists {
			// Check if this is a valid entry considering access lists
			if !isExcludedByAccessList(shouldLink, accessList) {
				result.MissingEntries = append(result.MissingEntries, shouldLink)
			}
		}
	}

	// Check for extra entries (in is but not in should)
	for key, isLink := range isStateMap {
		if _, exists := shouldStateMap[key]; !exists {
			result.ExtraEntries = append(result.ExtraEntries, isLink)
		}
	}

	// Set anomaly detection flags
	result.AnomalyDetection.UnderAdvertised = len(result.MissingEntries) > 0
	result.AnomalyDetection.OverAdvertised = len(result.ExtraEntries) > 0

	return result
}

// Helper function to get a unique key for an advertisement
func getAdvertismentKey(adv advertisment) string {
	// Prefer interface address if available, otherwise use link state ID
	if adv.InterfaceAddress != "" {
		return normalizeNetworkAddress(adv.InterfaceAddress)
	}
	return normalizeNetworkAddress(adv.LinkStateId)
}

// Helper function to normalize network addresses for comparison
func normalizeNetworkAddress(address string) string {
	// This is a simplified implementation
	// In a real-world scenario, you'd want to normalize IP addresses
	// by removing spaces and converting to a canonical form
	return strings.TrimSpace(address)
}

// Helper function to check if an advertisement should be excluded based on access lists
func isExcludedByAccessList(adv advertisment, accessLists map[string]accessList) bool {
	// This is a simplified implementation
	// In a real-world scenario, you'd want to check if the network address
	// is denied by any of the relevant access lists

	// For demonstration purposes, we'll assume that if a network appears in a "deny" entry
	// in any access list, it should be excluded
	for _, acl := range accessLists {
		for _, entry := range acl.aclEntry {
			if !entry.IsPermit {
				// If this is a deny entry, check if the advertisement matches
				if entry.Any {
					// If the entry denies any, check if the advertisement is not explicitly permitted elsewhere
					// This is a simplification; in reality, you'd need more sophisticated logic
					return true
				} else {
					// Check if the advertisement's network matches this deny entry
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

func nssaExternalAnomalyAnalysis(accessList map[string]accessList, isState *interAreaLsa, shouldState *interAreaLsa) {

}
