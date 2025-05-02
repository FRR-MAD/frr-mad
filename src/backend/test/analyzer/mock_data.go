package analyzer_test

import (
	"fmt"
	"testing"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/analyzer"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/ba2025-ysmprc/frr-mad/src/logger"
)

type AnalyzerConfig struct {
	Foo string `mapstructure:"foo"`
}

// is that even needed? Yes.
func initAnalyzer() {

	config, metrics, logger := getMockData()

	analyzer := analyzer.InitAnalyzer(config, metrics, logger)
	fmt.Println(analyzer)
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

func getFullFRRData1() *frrProto.FullFRRData {

	staticFRRConfiguration := &frrProto.StaticFRRConfiguration{
		Hostname:           "router1",
		FrrVersion:         "8.1",
		Ipv6Forwarding:     true,
		Ipv4Forwarding:     true,
		ServiceAdvancedVty: true,
		Interfaces: []*frrProto.Interface{
			{
				Name: "eth0",
				InterfaceIpPrefixes: []*frrProto.InterfaceIPPrefix{
					{
						IpPrefix: &frrProto.IPPrefix{
							IpAddress:    "192.168.1.1",
							PrefixLength: 24,
						},
						Passive: false,
						HasPeer: true,
						PeerIpPrefix: &frrProto.IPPrefix{
							IpAddress:    "192.168.1.2",
							PrefixLength: 24,
						},
					},
				},
				Area: "0.0.0.0",
			},
		},
		StaticRoutes: []*frrProto.StaticRoute{
			{
				IpPrefix: &frrProto.IPPrefix{
					IpAddress:    "10.0.0.0",
					PrefixLength: 8,
				},
				NextHop: "192.168.1.254",
			},
		},
		RouteMap: map[string]*frrProto.RouteMap{
			"ROUTE_MAP_1": {
				Permit:     true,
				Sequence:   "10",
				Match:      "ip address",
				AccessList: "ACL_1",
			},
		},
		AccessList: map[string]*frrProto.AccessList{
			"ACL_1": {
				Name: "ACL_1",
				AccessListItems: []*frrProto.AccessListItem{
					{
						Sequence:      10,
						AccessControl: "permit",
						Destination: &frrProto.AccessListItem_IpPrefix{
							IpPrefix: &frrProto.IPPrefix{
								IpAddress:    "10.0.0.0",
								PrefixLength: 8,
							},
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

func getFullFRRData2() *frrProto.FullFRRData {
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

func getFullFRRData3() *frrProto.FullFRRData {

	staticFRRConfiguration := &frrProto.StaticFRRConfiguration{
		Hostname:           "r102",
		FrrVersion:         "8.5.4_git",
		ServiceAdvancedVty: true,
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

	metrics := &frrProto.FullFRRData{
		StaticFrrConfiguration: staticFRRConfiguration,
	}

	return metrics
}

func FoobarTesting(t *testing.T) {

}
