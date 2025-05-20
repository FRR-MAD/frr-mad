package analyzer_test

import (
	"strconv"
	"strings"

	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
)

func getIfaceMap(value []*frrProto.Advertisement) (map[string]*frrProto.Advertisement, []string) {
	result := make(map[string]*frrProto.Advertisement)
	keyResult := []string{}

	for _, iface := range value {
		keyResult = append(keyResult, zeroLastOctetString(iface.InterfaceAddress))
		result[zeroLastOctetString(iface.InterfaceAddress)] = &frrProto.Advertisement{
			InterfaceAddress: zeroLastOctetString(iface.InterfaceAddress),
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

func GetNssaExternalData(data *frrProto.OSPFNssaExternalData, hostname string) *frrProto.InterAreaLsa {
	result := &frrProto.InterAreaLsa{
		Hostname: hostname,
		RouterId: data.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	for areaID, areaData := range data.NssaExternalLinkStates {
		area := &frrProto.AreaAnalyzer{
			AreaName: areaID,
			LsaType:  "NSSA-LSA",
			Links:    []*frrProto.Advertisement{},
		}

		for _, lsa := range areaData.Data {
			area.Links = append(area.Links, &frrProto.Advertisement{
				LinkStateId:  lsa.LinkStateId,
				PrefixLength: strconv.Itoa(int(lsa.NetworkMask)),
				LinkType:     "nssa-external",
			})
		}

		result.Areas = append(result.Areas, area)
	}

	return result
}

func zeroLastOctetString(ipAddress string) string {
	parts := strings.Split(ipAddress, ".")

	parts[3] = "0"

	return strings.Join(parts, ".")
}
