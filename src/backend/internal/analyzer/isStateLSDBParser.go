package analyzer

import (
	"strconv"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

// lsa type 5 parsing
func getRuntimeExternalRouterData(config *frrProto.OSPFExternalData, hostname string) *interAreaLsa {
	if config == nil {
		return nil
	}

	// TODO: check if redistribute has a route-map and only if compare to route-map lists
	result := &interAreaLsa{
		Hostname: hostname,
		RouterId: config.RouterId,
		Areas:    []area{},
	}

	// Since AS-external-LSA (type 5) doesn't belong to a specific area,
	// we'll create a single "area" to represent the AS external links
	externalArea := area{
		AreaName: "External",
		LsaType:  "AS-external-LSA", // Type 5
		Links:    []frrProto.Advertisement{},
	}

	for _, lsa := range config.AsExternalLinkStates {
		adv := frrProto.Advertisement{
			LinkStateId:  lsa.LinkStateId,
			PrefixLength: strconv.Itoa(int(lsa.NetworkMask)),
			LinkType:     "external",
		}

		externalArea.Links = append(externalArea.Links, adv)
	}

	result.Areas = append(result.Areas, externalArea)

	return result
}

// lsa type 7 parsing
func getNssaExternalRouterData(config *frrProto.OSPFNssaExternalData, hostname string) *interAreaLsa {
	if config == nil {
		return nil
	}

	result := &interAreaLsa{
		Hostname: hostname,
		RouterId: config.RouterId,
		Areas:    []area{},
	}

	for areaId, nssaArea := range config.NssaExternalLinkStates {
		nssaAreaObj := area{
			AreaName: areaId,
			LsaType:  "NSSA-LSA",
			Links:    []frrProto.Advertisement{},
		}

		for _, lsa := range nssaArea.Data {
			adv := frrProto.Advertisement{
				LinkStateId:  lsa.LinkStateId,
				PrefixLength: strconv.Itoa(int(lsa.NetworkMask)),
				LinkType:     "nssa-external",
			}

			nssaAreaObj.Links = append(nssaAreaObj.Links, adv)
		}

		result.Areas = append(result.Areas, nssaAreaObj)
	}

	return result
}

// lsa type 1 parsing
func getRuntimeRouterData(config *frrProto.OSPFRouterData, hostname string) *intraAreaLsa {
	result := intraAreaLsa{
		RouterId: config.RouterId,
		Areas:    []area{},
	}

	for areaName, routerArea := range config.RouterStates {
		for _, lsaEntry := range routerArea.LsaEntries {
			var currentArea *area
			for i := range result.Areas {
				if result.Areas[i].AreaName == areaName {
					currentArea = &result.Areas[i]
					break
				}
			}

			if currentArea == nil {
				newArea := area{
					AreaName: areaName,
					LsaType:  lsaEntry.LsaType,
					Links:    []frrProto.Advertisement{},
				}
				result.Areas = append(result.Areas, newArea)
				currentArea = &result.Areas[len(result.Areas)-1]
			}

			for _, routerLink := range lsaEntry.RouterLinks {
				var ipAddress, prefixLength string
				isStub := false
				if routerLink.LinkType == "Stub Network" {
					ipAddress = routerLink.NetworkAddress
					isStub = true
					prefixLength = maskToPrefixLength(routerLink.NetworkMask)
				} else if routerLink.LinkType == "a Transit Network" {
					ipAddress = routerLink.RouterInterfaceAddress
					//prefixLength = "24" // Assuming a /24 for transit links
				} else {
					if routerLink.RouterInterfaceAddress != "" {
						ipAddress = routerLink.RouterInterfaceAddress
					} else if routerLink.NetworkAddress != "" {
						ipAddress = routerLink.NetworkAddress
						//prefixLength = maskToPrefixLength(routerLink.NetworkMask)
					} else {
						continue
					}
				}

				adv := frrProto.Advertisement{}
				adv.InterfaceAddress = ipAddress
				adv.LinkType = routerLink.LinkType

				if isStub {
					adv.PrefixLength = prefixLength
				}

				currentArea.Links = append(currentArea.Links, adv)
			}
		}
	}

	result.Hostname = hostname

	return &result
}
