package analyzer

import (
	"strconv"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

// lsa type 5 parsing
// this will only return static routes, as BGP routes aren't useful in ospf analysis
func GetRuntimeExternalRouterData(config *frrProto.OSPFExternalData, staticRouteMap map[string]*frrProto.StaticList, hostname string) *frrProto.InterAreaLsa {
	if config == nil {
		return nil
	}

	//fmt.Println("staticRotueMap")
	//fmt.Println(staticRouteMap)

	// TODO: check if redistribute has a route-map and only if compare to route-map lists
	result := &frrProto.InterAreaLsa{
		Hostname: hostname,
		RouterId: config.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	// Since AS-external-LSA (type 5) doesn't belong to a specific area,
	// we'll create a single "area" to represent the AS external links
	externalArea := frrProto.AreaAnalyzer{
		//AreaName: "External",
		LsaType: "AS-external-LSA", // Type 5
		Links:   []*frrProto.Advertisement{},
	}

	for key, lsa := range config.AsExternalLinkStates {
		if _, exists := staticRouteMap[key]; !exists {
			continue
		}
		adv := frrProto.Advertisement{
			LinkStateId:  lsa.LinkStateId,
			PrefixLength: strconv.Itoa(int(lsa.NetworkMask)),
			LinkType:     "external",
		}

		externalArea.Links = append(externalArea.Links, &adv)
	}

	result.Areas = append(result.Areas, &externalArea)

	return result
}

// lsa type 7 parsing
func GetNssaExternalRouterData(config *frrProto.OSPFNssaExternalData, hostname string) *frrProto.InterAreaLsa {
	if config == nil {
		return nil
	}

	result := &frrProto.InterAreaLsa{
		Hostname: hostname,
		RouterId: config.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	for areaId, nssaArea := range config.NssaExternalLinkStates {
		nssaAreaObj := frrProto.AreaAnalyzer{
			AreaName: areaId,
			LsaType:  "NSSA-LSA",
			Links:    []*frrProto.Advertisement{},
		}

		for _, lsa := range nssaArea.Data {
			adv := frrProto.Advertisement{
				LinkStateId:  lsa.LinkStateId,
				PrefixLength: strconv.Itoa(int(lsa.NetworkMask)),
				LinkType:     "nssa-external",
			}

			nssaAreaObj.Links = append(nssaAreaObj.Links, &adv)
		}

		result.Areas = append(result.Areas, &nssaAreaObj)
	}

	return result
}

// lsa type 1 parsing
func GetRuntimeRouterData(config *frrProto.OSPFRouterData, hostname string) *frrProto.IntraAreaLsa {
	result := frrProto.IntraAreaLsa{
		RouterId: config.RouterId,
		Areas:    []*frrProto.AreaAnalyzer{},
	}

	//for _, value := range config.RouterStates {
	//	fmt.Println(value)
	//}

	for areaName, routerArea := range config.RouterStates {
		for _, lsaEntry := range routerArea.LsaEntries {
			var currentArea *frrProto.AreaAnalyzer
			for i := range result.Areas {
				if result.Areas[i].AreaName == areaName {
					currentArea = result.Areas[i]
					break
				}
			}

			if currentArea == nil {
				newArea := frrProto.AreaAnalyzer{
					AreaName: areaName,
					LsaType:  lsaEntry.LsaType,
					Links:    []*frrProto.Advertisement{},
				}
				result.Areas = append(result.Areas, &newArea)
				currentArea = result.Areas[len(result.Areas)-1]
			}

			for _, routerLink := range lsaEntry.RouterLinks {
				var ipAddress, prefixLength string
				isStub := false
				if routerLink.LinkType == "Stub Network" {
					routerLink.LinkType = "stub network"
					ipAddress = routerLink.NetworkAddress
					isStub = true
					prefixLength = maskToPrefixLength(routerLink.NetworkMask)
				} else if routerLink.LinkType == "a Transit Network" {
					routerLink.LinkType = "transit network"
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
				if routerLink.LinkType == "another Router (point-to-point)" {
					adv.LinkType = "point-to-point"
				} else {
					adv.LinkType = routerLink.LinkType
				}

				if isStub {
					adv.PrefixLength = prefixLength
				}

				currentArea.Links = append(currentArea.Links, &adv)
			}
		}
	}

	result.Hostname = hostname

	return &result
}
