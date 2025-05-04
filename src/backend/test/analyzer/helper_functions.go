package analyzer_test

import (
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

func getIfaceMap(value []*frrProto.Advertisement) (map[string]*frrProto.Advertisement, []string) {
	result := make(map[string]*frrProto.Advertisement)
	keyResult := []string{}

	for _, iface := range value {
		keyResult = append(keyResult, iface.InterfaceAddress)
		result[iface.InterfaceAddress] = &frrProto.Advertisement{
			InterfaceAddress: iface.InterfaceAddress,
			PrefixLength:     iface.PrefixLength,
			LinkType:         iface.LinkType,
		}
	}

	return result, keyResult
}

func uniqueNonEmptyElementsOf(s []string) []string {
	unique := make(map[string]bool, len(s))
	us := make([]string, len(unique))
	for _, elem := range s {
		if len(elem) != 0 {
			if !unique[elem] {
				us = append(us, elem)
				unique[elem] = true
			}
		}
	}

	return us

}
