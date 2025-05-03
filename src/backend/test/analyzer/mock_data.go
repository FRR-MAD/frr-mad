package analyzer_test

import (
	"testing"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/analyzer"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/ba2025-ysmprc/frr-mad/src/logger"
)

type AnalyzerConfig struct {
	Foo string `mapstructure:"foo"`
}

// is that even needed? Yes.
func initAnalyzer() *analyzer.Analyzer {

	config, metrics, logger := getMockData()

	return analyzer.InitAnalyzer(config, metrics, logger)
}

func getMockData() (map[string]interface{}, *frrProto.FullFRRData, *logger.Logger) {

	config := map[string]interface{}{
		"key1": "value1",
		"key2": 42,
		"key3": true,
	}

	metrics := &frrProto.FullFRRData{}
	logger, _ := logger.NewLogger("testing", "/tmp/testing.log")

	return config, metrics, logger

}

func getR101FRRdata() *frrProto.FullFRRData {

	staticFRRConfiguration := &frrProto.StaticFRRConfiguration{
		Hostname:           "r101",
		FrrVersion:         "8.5.4_git",
		ServiceAdvancedVty: true,
		Interfaces: []*frrProto.Interface{
			{
				Name: "eth1",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "172.22.1.1",
							PrefixLength: 24,
						},
					},
				},
			},
			{
				Name: "eth2",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.12.1",
							PrefixLength: 24,
						},
					},
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.2.1",
							PrefixLength: 24,
						},
						Passive: true,
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth3",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.13.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth4",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.0.1",
							PrefixLength: 23,
						},
						Passive: true,
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth5",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "192.168.100.1",
							PrefixLength: 24,
						},
					},
				},
			},
			{
				Name: "eth6",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.14.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth7",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.15.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth8",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.16.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth9",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.17.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth10",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.18.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth11",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.19.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "lo",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "65.0.1.1",
							PrefixLength: 32,
						},
						Passive: true,
					},
				},
			},
		},
		StaticRoutes: []*frrProto.StaticRoute{
			{
				IpPrefix: &frrProto.IPPrefix{
					IpAddress:    "192.168.1.0",
					PrefixLength: 24,
				},
				NextHop: "192.168.100.91",
			},
		},
		OspfConfig: &frrProto.OSPFConfig{
			RouterId: "65.0.1.1",
			Redistribution: []*frrProto.Redistribution{
				{
					Type:     "static",
					Metric:   "1",
					RouteMap: "lanroutes",
				},
				{
					Type:   "bgp",
					Metric: "1",
				},
			},
		},
		RouteMap: map[string]*frrProto.RouteMap{
			"lanroutes": {
				Permit:     true,
				Sequence:   "10",
				Match:      "ip address",
				AccessList: "localsite",
			},
		},
		AccessList: map[string]*frrProto.AccessList{
			"localsite": {
				AccessListItems: []*frrProto.AccessListItem{
					{
						Sequence:      15,
						AccessControl: "permit",
						Destination: &frrProto.AccessListItem_IpPrefix{
							IpPrefix: &frrProto.IPPrefix{
								IpAddress:    "192.168.1.0",
								PrefixLength: 24,
							},
						},
					},
				},
			},
			"term": {
				AccessListItems: []*frrProto.AccessListItem{
					{
						Sequence:      5,
						AccessControl: "permit",
						Destination: &frrProto.AccessListItem_IpPrefix{
							IpPrefix: &frrProto.IPPrefix{
								IpAddress:    "127.0.0.1",
								PrefixLength: 32,
							},
						},
					},
					{
						Sequence:      10,
						AccessControl: "deny",
						Destination: &frrProto.AccessListItem_Any{
							Any: true,
						},
					},
				},
			},
		},
	}

	ospfRouterData := &frrProto.OSPFRouterData{
		RouterId: "65.0.1.1",
		RouterStates: map[string]*frrProto.OSPFRouterArea{
			"0.0.0.0": {
				LsaEntries: map[string]*frrProto.OSPFRouterLSA{
					"65.0.1.1": {
						LsaAge:            40,
						Options:           "*|-|-|-|-|-|E|-",
						LsaFlags:          3,
						Flags:             2,
						Asbr:              true,
						LsaType:           "router-LSA",
						LinkStateId:       "65.0.1.1",
						AdvertisingRouter: "65.0.1.1",
						LsaSeqNumber:      "8000002a",
						Checksum:          "b021",
						Length:            144,
						NumOfLinks:        10,
						RouterLinks: map[string]*frrProto.OSPFRouterLSALink{
							"link0": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.0.12.2",
								RouterInterfaceAddress:  "10.0.12.1",
								Tos0Metric:              10,
							},
							"link1": {
								LinkType:       "Stub Network",
								NetworkAddress: "10.0.2.0",
								NetworkMask:    "255.255.255.0",
								Tos0Metric:     10,
							},
							"link2": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.0.14.4",
								RouterInterfaceAddress:  "10.0.14.1",
								Tos0Metric:              10,
							},
							"link3": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.0.16.6",
								RouterInterfaceAddress:  "10.0.16.1",
								Tos0Metric:              10,
							},
							"link4": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.0.18.8",
								RouterInterfaceAddress:  "10.0.18.1",
								Tos0Metric:              10,
							},
							"link5": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.0.15.5",
								RouterInterfaceAddress:  "10.0.15.1",
								Tos0Metric:              10,
							},
							"link6": {
								LinkType:       "Stub Network",
								NetworkAddress: "10.0.0.0",
								NetworkMask:    "255.255.254.0",
								Tos0Metric:     10,
							},
							"link7": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.0.17.7",
								RouterInterfaceAddress:  "10.0.17.1",
								Tos0Metric:              10,
							},
							"link8": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.0.13.3",
								RouterInterfaceAddress:  "10.0.13.1",
								Tos0Metric:              10,
							},
							"link9": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.0.19.9",
								RouterInterfaceAddress:  "10.0.19.1",
								Tos0Metric:              10,
							},
						},
					},
				},
			},
		},
	}

	metrics := &frrProto.FullFRRData{
		StaticFrrConfiguration: staticFRRConfiguration,
		OspfRouterData:         ospfRouterData,
	}

	return metrics
}

// type accessList struct {
// 	accessListName string     `json:"access_list_name"`
// 	aclEntry       []ACLEntry `json:"acl_entries"`
// }

// type ACLEntry struct {
// 	IPAddress    string `json:"ip_address,omitempty"`
// 	PrefixLength int    `json:"prefix_length,omitempty"`
// 	IsPermit     bool   `json:"is_permit"`
// 	Any          bool   `json:"any,omitempty"`
// 	Sequence     int    `json:"sequence"`
// }

func getR102FRRdata() *frrProto.FullFRRData {
	staticFRRConfiguration := &frrProto.StaticFRRConfiguration{
		Hostname:   "r102",
		FrrVersion: "8.5.4_git",
		Interfaces: []*frrProto.Interface{
			{
				Name: "eth1",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.12.2",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth2",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.23.2",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth3",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.1.21.2",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.1",
			},
			{
				Name: "eth4",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "192.168.101.2",
							PrefixLength: 24,
						},
					},
				},
			},
			{
				Name: "lo",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "65.0.1.2",
							PrefixLength: 32,
						},
						Passive: true,
					},
				},
			},
		},
		StaticRoutes: []*frrProto.StaticRoute{
			{
				IpPrefix: &frrProto.IPPrefix{
					IpAddress:    "192.168.11.0",
					PrefixLength: 24,
				},
				NextHop: "192.168.101.93",
			},
		},
		OspfConfig: &frrProto.OSPFConfig{
			RouterId: "65.0.1.2",
			Redistribution: []*frrProto.Redistribution{
				{
					Type:     "static",
					Metric:   "1",
					RouteMap: "lanroutes",
				},
			},
			Area: []*frrProto.Area{
				{
					Name: "0.0.0.1",
					Type: "nssa",
				},
			},
		},
		RouteMap: map[string]*frrProto.RouteMap{
			"lanroutes": {
				Permit:     true,
				Sequence:   "10",
				Match:      "ip address",
				AccessList: "localsite",
			},
		},
		AccessList: map[string]*frrProto.AccessList{
			"localsite": {
				AccessListItems: []*frrProto.AccessListItem{
					{
						Sequence:      15,
						AccessControl: "permit",
						Destination: &frrProto.AccessListItem_IpPrefix{
							IpPrefix: &frrProto.IPPrefix{
								IpAddress:    "192.168.11.0",
								PrefixLength: 24,
							},
						},
					},
				},
			},
			"term": {
				AccessListItems: []*frrProto.AccessListItem{
					{
						Sequence:      5,
						AccessControl: "permit",
						Destination: &frrProto.AccessListItem_IpPrefix{
							IpPrefix: &frrProto.IPPrefix{
								IpAddress:    "127.0.0.1",
								PrefixLength: 32,
							},
						},
					},
					{
						Sequence:      10,
						AccessControl: "deny",
						Destination: &frrProto.AccessListItem_Any{
							Any: true,
						},
					},
				},
			},
		},
	}

	ospfRouterData := &frrProto.OSPFRouterData{
		RouterId: "65.0.1.2",
		RouterStates: map[string]*frrProto.OSPFRouterArea{
			"0.0.0.0": {
				LsaEntries: map[string]*frrProto.OSPFRouterLSA{
					"65.0.1.2": {
						LsaAge:            609,
						Options:           "*|-|-|-|-|-|E|-",
						LsaFlags:          3,
						Flags:             3,
						Asbr:              true,
						LsaType:           "router-LSA",
						LinkStateId:       "65.0.1.2",
						AdvertisingRouter: "65.0.1.2",
						LsaSeqNumber:      "8000000d",
						Checksum:          "45dc",
						Length:            48,
						NumOfLinks:        2,
						RouterLinks: map[string]*frrProto.OSPFRouterLSALink{
							"link0": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.0.12.2",
								RouterInterfaceAddress:  "10.0.12.2",
								Tos0Metric:              10,
							},
							"link1": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.0.23.3",
								RouterInterfaceAddress:  "10.0.23.2",
								Tos0Metric:              10,
							},
						},
					},
				},
			},
			"0.0.0.1": {
				LsaEntries: map[string]*frrProto.OSPFRouterLSA{
					"65.0.1.2": {
						LsaAge:            639,
						Options:           "*|-|-|-|N/P|-|-|-",
						LsaFlags:          3,
						Flags:             3,
						Asbr:              true,
						LsaType:           "router-LSA",
						LinkStateId:       "65.0.1.2",
						AdvertisingRouter: "65.0.1.2",
						LsaSeqNumber:      "8000000c",
						Checksum:          "5e02",
						Length:            36,
						NumOfLinks:        1,
						RouterLinks: map[string]*frrProto.OSPFRouterLSALink{
							"link0": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.1.21.11",
								RouterInterfaceAddress:  "10.1.21.2",
								Tos0Metric:              10,
							},
						},
					},
				},
			},
		},
	}

	metrics := &frrProto.FullFRRData{
		StaticFrrConfiguration: staticFRRConfiguration,
		OspfRouterData:         ospfRouterData,
	}

	return metrics
}

func backup() *frrProto.FullFRRData {
	staticFRRConfiguration := &frrProto.StaticFRRConfiguration{
		Hostname:           "r101",
		FrrVersion:         "8.5.4_git",
		ServiceAdvancedVty: true,
		Interfaces: []*frrProto.Interface{
			{
				Name: "eth1",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "172.22.1.1",
							PrefixLength: 24,
						},
					},
				},
			},
			{
				Name: "eth2",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.12.1",
							PrefixLength: 24,
						},
					},
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.2.1",
							PrefixLength: 24,
						},
						Passive: true,
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth3",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.13.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth4",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.0.1",
							PrefixLength: 23,
						},
						Passive: true,
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth5",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "192.168.100.1",
							PrefixLength: 24,
						},
					},
				},
			},
			{
				Name: "eth6",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.14.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth7",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.15.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth8",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.16.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth9",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.17.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth10",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.18.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth11",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.19.1",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "lo",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "65.0.1.1",
							PrefixLength: 32,
						},
						Passive: true,
					},
				},
			},
		},
		StaticRoutes: []*frrProto.StaticRoute{
			{
				IpPrefix: &frrProto.IPPrefix{
					IpAddress:    "192.168.1.0",
					PrefixLength: 24,
				},
				NextHop: "192.168.100.91",
			},
		},
		OspfConfig: &frrProto.OSPFConfig{
			RouterId: "65.0.1.1",
			//		route_map:"lanroutes"} redistribution:{type:"bgp" metric:"1"}}
			Redistribution: []*frrProto.Redistribution{
				{
					Type:   "static",
					Metric: "1",
				},
				{
					Type:   "bgp",
					Metric: "1",
				},
			},
		},
		RouteMap: map[string]*frrProto.RouteMap{
			"lanroutes": {
				Permit:     true,
				Sequence:   "10",
				Match:      "ip address",
				AccessList: "localsite",
			},
		},
		AccessList: map[string]*frrProto.AccessList{
			"localsite": {
				Name: "localsite",
				AccessListItems: []*frrProto.AccessListItem{
					{
						Sequence:      15,
						AccessControl: "permit",
						Destination: &frrProto.AccessListItem_IpPrefix{
							IpPrefix: &frrProto.IPPrefix{
								IpAddress:    "192.168.1.0",
								PrefixLength: 24,
							},
						},
					},
				},
			},
			"term": {
				AccessListItems: []*frrProto.AccessListItem{
					{
						Sequence:      5,
						AccessControl: "permit",
						Destination: &frrProto.AccessListItem_IpPrefix{
							IpPrefix: &frrProto.IPPrefix{
								IpAddress:    "127.0.0.1",
								PrefixLength: 32,
							},
						},
					},
					{
						Sequence:      10,
						AccessControl: "deny",
						Destination: &frrProto.AccessListItem_Any{
							Any: true,
						},
					},
				},
			},
		},
	}

	metrics := &frrProto.FullFRRData{
		StaticFrrConfiguration: staticFRRConfiguration,
	}

	return metrics
}

func getR103FRRdata() *frrProto.FullFRRData {

	staticFRRConfiguration := &frrProto.StaticFRRConfiguration{
		Hostname:   "r103",
		FrrVersion: "8.5.4_git",
		Interfaces: []*frrProto.Interface{
			{
				Name: "eth1",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.13.3",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth2",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.0.23.3",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
			{
				Name: "eth3",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "10.2.31.3",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.2",
			},
			{
				Name: "lo",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "65.0.1.3",
							PrefixLength: 32,
						},
						Passive: true,
					},
				},
			},
		},
		OspfConfig: &frrProto.OSPFConfig{
			RouterId: "65.0.1.3",
			Area: []*frrProto.Area{
				{

					Name: "0.0.0.2",
					Type: "transit",
				},
			},
			VirtualLinkNeighbor: "65.0.1.22",
		},
	}

	ospfRouterData := &frrProto.OSPFRouterData{
		RouterId: "65.0.1.3",
		RouterStates: map[string]*frrProto.OSPFRouterArea{
			"0.0.0.0": {
				LsaEntries: map[string]*frrProto.OSPFRouterLSA{
					"65.0.1.3": {
						LsaAge:            99,
						Options:           "*|-|-|-|-|-|E|-",
						LsaFlags:          3,
						Flags:             1,
						LsaType:           "router-LSA",
						LinkStateId:       "65.0.1.3",
						AdvertisingRouter: "65.0.1.3",
						LsaSeqNumber:      "80000010",
						Checksum:          "204e",
						Length:            60,
						NumOfLinks:        3,
						RouterLinks: map[string]*frrProto.OSPFRouterLSALink{
							"link0": {
								LinkType:               "a Virtual Link",
								RouterInterfaceAddress: "10.2.31.3",
								Tos0Metric:             20,
							},
							"link1": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.0.13.3",
								RouterInterfaceAddress:  "10.0.13.3",
								Tos0Metric:              10,
							},
							"link2": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.0.23.3",
								RouterInterfaceAddress:  "10.0.23.3",
								Tos0Metric:              10,
							},
						},
					},
				},
			},
			"0.0.0.2": {
				LsaEntries: map[string]*frrProto.OSPFRouterLSA{
					"65.0.1.3": {
						LsaAge:            149,
						Options:           "*|-|-|-|-|-|E|-",
						LsaFlags:          3,
						Flags:             5,
						LsaType:           "router-LSA",
						LinkStateId:       "65.0.1.3",
						AdvertisingRouter: "65.0.1.3",
						LsaSeqNumber:      "8000000c",
						Checksum:          "122f",
						Length:            36,
						NumOfLinks:        1,
						RouterLinks: map[string]*frrProto.OSPFRouterLSALink{
							"link0": {
								LinkType:                "a Transit Network",
								DesignatedRouterAddress: "10.2.31.21",
								RouterInterfaceAddress:  "10.2.31.3",
								Tos0Metric:              10,
							},
						},
					},
				},
			},
		},
	}

	metrics := &frrProto.FullFRRData{
		StaticFrrConfiguration: staticFRRConfiguration,
		OspfRouterData:         ospfRouterData,
	}

	return metrics
}

func template() *frrProto.FullFRRData {

	staticFRRConfiguration := &frrProto.StaticFRRConfiguration{}

	ospfRouterData := &frrProto.OSPFRouterData{}

	metrics := &frrProto.FullFRRData{
		StaticFrrConfiguration: staticFRRConfiguration,
		OspfRouterData:         ospfRouterData,
	}

	return metrics
}

func getRXXXFRRData() *frrProto.FullFRRData {

	staticFRRConfiguration := &frrProto.StaticFRRConfiguration{}

	ospfRouterData := &frrProto.OSPFRouterData{}

	metrics := &frrProto.FullFRRData{
		StaticFrrConfiguration: staticFRRConfiguration,
		OspfRouterData:         ospfRouterData,
	}

	return metrics
}

func FoobarTesting(t *testing.T) {

}
