package common

import (
	"fmt"
	frrProto "github.com/ba2025-ysmprc/frr-tui/pkg"
)

func ContainsString(slice []string, s string) bool {
	for _, v := range slice {
		if v == s {
			return true
		}
	}
	return false
}

func HasAnyAnomaly(a *frrProto.AnomalyDetection) bool {
	if a == nil {
		return false
	}
	return a.HasUnderAdvertisedPrefixes ||
		a.HasOverAdvertisedPrefixes ||
		a.HasDuplicatePrefixes ||
		a.HasMisconfiguredPrefixes
}

func PrintBackendError(err error, functionName string) string {
	return fmt.Sprintf(
		"Error: %v\nNo data received from backend for '%s()'. Press 'r' to reload...",
		err, functionName,
	)
}
