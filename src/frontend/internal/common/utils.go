package common

import frrProto "github.com/ba2025-ysmprc/frr-tui/pkg"

func ContainsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func HasAnyAnomaly(a *frrProto.AnomalyDetection) bool {
	return a.HasUnderAdvertisedPrefixes ||
		a.HasOverAdvertisedPrefixes ||
		a.HasDuplicatePrefixes ||
		a.HasMisconfiguredPrefixes
}
