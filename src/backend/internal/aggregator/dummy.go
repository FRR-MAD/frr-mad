package aggregator

import (
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func DummyFunction() string {
	return "I am aggregator Aggregator"
}

func OSPFNeighborDummyData() []*frrProto.OSPFNeighbor {
	var neighbors []*frrProto.OSPFNeighbor

	neighborList := []frrProto.OSPFNeighbor{
		{
			Id:        "65.0.1.8",
			Ip:        "10.0.18.8",
			Interface: "eth10",
			Area:      "0.0.0.0",
		},
		{
			Id:        "65.0.1.9",
			Ip:        "10.0.19.9",
			Interface: "eth11",
			Area:      "0.0.0.0",
		},
		{
			Id:        "65.0.1.2",
			Ip:        "10.0.12.2",
			Interface: "eth2",
			Area:      "0.0.0.0",
		},
		{
			Id:        "65.0.1.3",
			Ip:        "10.0.13.3",
			Interface: "eth3",
			Area:      "0.0.0.0",
		},
		{
			Id:        "65.0.1.4",
			Ip:        "10.0.14.4",
			Interface: "eth6",
			Area:      "0.0.0.0",
		},
		{
			Id:        "65.0.1.5",
			Ip:        "10.0.15.5",
			Interface: "eth7",
			Area:      "0.0.0.0",
		},
		{
			Id:        "65.0.1.6",
			Ip:        "10.0.16.6",
			Interface: "eth8",
			Area:      "0.0.0.0",
		},
		{
			Id:        "65.0.1.7",
			Ip:        "10.0.17.7",
			Interface: "eth9",
			Area:      "0.0.0.0",
		},
	}
	for _, neighbor := range neighborList {
		neighbors = append(neighbors, &neighbor)
	}

	return neighbors
}

func OSPFRouteDummyData() []*frrProto.OSPFRoute {
	var routes []*frrProto.OSPFRoute
	routeList := []frrProto.OSPFRoute{
		{
			Prefix:    "192.168.8.0/24",
			NextHop:   "10.0.18.8",
			Interface: "eth10",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.8",
			NextHop:   "10.0.18.8",
			Interface: "eth10",
			Cost:      10,
			Type:      "R",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.8/32",
			NextHop:   "10.0.18.8",
			Interface: "eth10",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.0.18.0/24",
			NextHop:   "direct",
			Interface: "eth10",
			Cost:      10,
			Type:      "N",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "192.168.9.0/24",
			NextHop:   "10.0.19.9",
			Interface: "eth11",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.9",
			NextHop:   "10.0.19.9",
			Interface: "eth11",
			Cost:      10,
			Type:      "R",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.9/32",
			NextHop:   "10.0.19.9",
			Interface: "eth11",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.0.19.0/24",
			NextHop:   "direct",
			Interface: "eth11",
			Cost:      10,
			Type:      "N",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.0.23.0/24",
			NextHop:   "10.0.12.2",
			Interface: "eth2",
			Cost:      20,
			Type:      "N",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.1.0.0/24",
			NextHop:   "10.0.12.2",
			Interface: "eth2",
			Cost:      30,
			Type:      "N IA",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.1.12.0/24",
			NextHop:   "10.0.12.2",
			Interface: "eth2",
			Cost:      30,
			Type:      "N IA",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.1.21.0/24",
			NextHop:   "10.0.12.2",
			Interface: "eth2",
			Cost:      20,
			Type:      "N IA",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.30.12.0/24",
			NextHop:   "10.0.12.2",
			Interface: "eth2",
			Cost:      50,
			Type:      "N E1",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.2",
			NextHop:   "10.0.12.2",
			Interface: "eth2",
			Cost:      10,
			Type:      "R",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.3.1/32",
			NextHop:   "10.0.12.2",
			Interface: "eth2",
			Cost:      50,
			Type:      "N E1",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.0.12.0/24",
			NextHop:   "direct",
			Interface: "eth2",
			Cost:      10,
			Type:      "N",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.2.0.0/24",
			NextHop:   "10.0.13.3",
			Interface: "eth3",
			Cost:      30,
			Type:      "N IA",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.2.12.0/24",
			NextHop:   "10.0.13.3",
			Interface: "eth3",
			Cost:      30,
			Type:      "N IA",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.2.31.0/24",
			NextHop:   "10.0.13.3",
			Interface: "eth3",
			Cost:      20,
			Type:      "N IA",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.3.0.0/24",
			NextHop:   "10.0.13.3",
			Interface: "eth3",
			Cost:      50,
			Type:      "N IA",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.3.21.0/24",
			NextHop:   "10.0.13.3",
			Interface: "eth3",
			Cost:      40,
			Type:      "N IA",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.22",
			NextHop:   "10.0.13.3",
			Interface: "eth3",
			Cost:      30,
			Type:      "R",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.3",
			NextHop:   "10.0.13.3",
			Interface: "eth3",
			Cost:      10,
			Type:      "R",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.0.13.0/24",
			NextHop:   "direct",
			Interface: "eth3",
			Cost:      10,
			Type:      "N",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.0.0.0/23",
			NextHop:   "direct",
			Interface: "eth4",
			Cost:      10,
			Type:      "N",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "172.20.20.0/24",
			NextHop:   "10.0.14.4",
			Interface: "eth6",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "192.168.4.0/24",
			NextHop:   "10.0.14.4",
			Interface: "eth6",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.4",
			NextHop:   "10.0.14.4",
			Interface: "eth6",
			Cost:      10,
			Type:      "R",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.4/32",
			NextHop:   "10.0.14.4",
			Interface: "eth6",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.0.14.0/24",
			NextHop:   "direct",
			Interface: "eth6",
			Cost:      10,
			Type:      "N",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "192.168.5.0/24",
			NextHop:   "10.0.15.5",
			Interface: "eth7",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.5",
			NextHop:   "10.0.15.5",
			Interface: "eth7",
			Cost:      10,
			Type:      "R",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.5/32",
			NextHop:   "10.0.15.5",
			Interface: "eth7",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.0.15.0/24",
			NextHop:   "direct",
			Interface: "eth7",
			Cost:      10,
			Type:      "N",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "192.168.6.0/24",
			NextHop:   "10.0.16.6",
			Interface: "eth8",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.6",
			NextHop:   "10.0.16.6",
			Interface: "eth8",
			Cost:      10,
			Type:      "R",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.6/32",
			NextHop:   "10.0.16.6",
			Interface: "eth8",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.0.16.0/24",
			NextHop:   "direct",
			Interface: "eth8",
			Cost:      10,
			Type:      "N",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "192.168.7.0/24",
			NextHop:   "10.0.17.7",
			Interface: "eth9",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.7",
			NextHop:   "10.0.17.7",
			Interface: "eth9",
			Cost:      10,
			Type:      "R",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "65.0.1.7/32",
			NextHop:   "10.0.17.7",
			Interface: "eth9",
			Cost:      10,
			Type:      "N E2",
			Area:      "0.0.0.0",
		},
		{
			Prefix:    "10.0.17.0/24",
			NextHop:   "direct",
			Interface: "eth9",
			Cost:      10,
			Type:      "N",
			Area:      "0.0.0.0",
		},
	}
	for _, route := range routeList {
		routes = append(routes, &route)
	}

	return routes
}

func OSPFInterfaceDummyData() []*frrProto.OSPFInterface {
	var interfaces []*frrProto.OSPFInterface
	interfaceList := []frrProto.OSPFInterface{
		{
			Name:     "eth2",
			Area:     "0.0.0.0",
			NbrCount: 1,
			NbrAdj:   1,
			Passive:  false,
		},
		{
			Name:     "eth3",
			Area:     "0.0.0.0",
			NbrCount: 1,
			NbrAdj:   1,
			Passive:  false,
		},
		{
			Name:     "eth4",
			Area:     "0.0.0.0",
			NbrCount: 0,
			NbrAdj:   0,
			Passive:  true,
		},
		{
			Name:     "eth6",
			Area:     "0.0.0.0",
			NbrCount: 1,
			NbrAdj:   1,
			Passive:  false,
		},
		{
			Name:     "eth7",
			Area:     "0.0.0.0",
			NbrCount: 1,
			NbrAdj:   1,
			Passive:  false,
		},
		{
			Name:     "eth8",
			Area:     "0.0.0.0",
			NbrCount: 1,
			NbrAdj:   1,
			Passive:  false,
		},
		{
			Name:     "eth9",
			Area:     "0.0.0.0",
			NbrCount: 1,
			NbrAdj:   1,
			Passive:  false,
		},
		{
			Name:     "eth10",
			Area:     "0.0.0.0",
			NbrCount: 1,
			NbrAdj:   1,
			Passive:  false,
		},
		{
			Name:     "eth11",
			Area:     "0.0.0.0",
			NbrCount: 1,
			NbrAdj:   1,
			Passive:  false,
		},
	}

	for _, int := range interfaceList {
		interfaces = append(interfaces, &int)
	}

	return interfaces
}

func OSPFlsaDummyData() []*frrProto.OSPFlsa {
	var lsas []*frrProto.OSPFlsa

	lsaList := []frrProto.OSPFlsa{
		{
			Type:      "external",
			LsId:      "10.20.0.0",
			AdvRouter: "65.0.1.1",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "10.20.12.0",
			AdvRouter: "65.0.1.1",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "192.168.1.0",
			AdvRouter: "65.0.1.1",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "router",
			LsId:      "65.0.1.1",
			AdvRouter: "65.0.1.1",
			Sequence:  "80000031",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "65.0.2.1",
			AdvRouter: "65.0.1.1",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "network",
			LsId:      "10.0.12.2",
			AdvRouter: "65.0.1.2",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "summary",
			LsId:      "10.1.0.0",
			AdvRouter: "65.0.1.2",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "summary",
			LsId:      "10.1.12.0",
			AdvRouter: "65.0.1.2",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "summary",
			LsId:      "10.1.21.0",
			AdvRouter: "65.0.1.2",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "10.30.0.0",
			AdvRouter: "65.0.1.2",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "10.30.12.0",
			AdvRouter: "65.0.1.2",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "router",
			LsId:      "65.0.1.2",
			AdvRouter: "65.0.1.2",
			Sequence:  "80000015",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "65.0.3.1",
			AdvRouter: "65.0.1.2",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "summary",
			LsId:      "10.2.0.0",
			AdvRouter: "65.0.1.22",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "summary",
			LsId:      "10.2.12.0",
			AdvRouter: "65.0.1.22",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "summary",
			LsId:      "10.2.31.0",
			AdvRouter: "65.0.1.22",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "summary",
			LsId:      "10.3.0.0",
			AdvRouter: "65.0.1.22",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "summary",
			LsId:      "10.3.21.0",
			AdvRouter: "65.0.1.22",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "router",
			LsId:      "65.0.1.22",
			AdvRouter: "65.0.1.22",
			Sequence:  "80000011",
			Area:      "0.0.0.0",
		},
		{
			Type:      "network",
			LsId:      "10.0.13.3",
			AdvRouter: "65.0.1.3",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "network",
			LsId:      "10.0.23.3",
			AdvRouter: "65.0.1.3",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "summary",
			LsId:      "10.2.0.0",
			AdvRouter: "65.0.1.3",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "summary",
			LsId:      "10.2.12.0",
			AdvRouter: "65.0.1.3",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "summary",
			LsId:      "10.2.31.0",
			AdvRouter: "65.0.1.3",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "router",
			LsId:      "65.0.1.3",
			AdvRouter: "65.0.1.3",
			Sequence:  "80000019",
			Area:      "0.0.0.0",
		},
		{
			Type:      "network",
			LsId:      "10.0.14.4",
			AdvRouter: "65.0.1.4",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "172.20.20.0",
			AdvRouter: "65.0.1.4",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "192.168.4.0",
			AdvRouter: "65.0.1.4",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "65.0.1.4",
			AdvRouter: "65.0.1.4",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "router",
			LsId:      "65.0.1.4",
			AdvRouter: "65.0.1.4",
			Sequence:  "80000010",
			Area:      "0.0.0.0",
		},
		{
			Type:      "network",
			LsId:      "10.0.15.5",
			AdvRouter: "65.0.1.5",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "172.20.20.0",
			AdvRouter: "65.0.1.5",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "192.168.5.0",
			AdvRouter: "65.0.1.5",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "65.0.1.5",
			AdvRouter: "65.0.1.5",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "router",
			LsId:      "65.0.1.5",
			AdvRouter: "65.0.1.5",
			Sequence:  "80000011",
			Area:      "0.0.0.0",
		},
		{
			Type:      "network",
			LsId:      "10.0.16.6",
			AdvRouter: "65.0.1.6",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "172.20.20.0",
			AdvRouter: "65.0.1.6",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "192.168.6.0",
			AdvRouter: "65.0.1.6",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "65.0.1.6",
			AdvRouter: "65.0.1.6",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "router",
			LsId:      "65.0.1.6",
			AdvRouter: "65.0.1.6",
			Sequence:  "80000012",
			Area:      "0.0.0.0",
		},
		{
			Type:      "network",
			LsId:      "10.0.17.7",
			AdvRouter: "65.0.1.7",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "172.20.20.0",
			AdvRouter: "65.0.1.7",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "192.168.7.0",
			AdvRouter: "65.0.1.7",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "65.0.1.7",
			AdvRouter: "65.0.1.7",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "router",
			LsId:      "65.0.1.7",
			AdvRouter: "65.0.1.7",
			Sequence:  "80000011",
			Area:      "0.0.0.0",
		},
		{
			Type:      "network",
			LsId:      "10.0.18.8",
			AdvRouter: "65.0.1.8",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "172.20.20.0",
			AdvRouter: "65.0.1.8",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "192.168.8.0",
			AdvRouter: "65.0.1.8",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "65.0.1.8",
			AdvRouter: "65.0.1.8",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "router",
			LsId:      "65.0.1.8",
			AdvRouter: "65.0.1.8",
			Sequence:  "80000013",
			Area:      "0.0.0.0",
		},
		{
			Type:      "network",
			LsId:      "10.0.19.9",
			AdvRouter: "65.0.1.9",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "172.20.20.0",
			AdvRouter: "65.0.1.9",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "192.168.9.0",
			AdvRouter: "65.0.1.9",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "external",
			LsId:      "65.0.1.9",
			AdvRouter: "65.0.1.9",
			Sequence:  "0",
			Area:      "0.0.0.0",
		},
		{
			Type:      "router",
			LsId:      "65.0.1.9",
			AdvRouter: "65.0.1.9",
			Sequence:  "80000010",
			Area:      "0.0.0.0",
		},
	}

	for _, lsa := range lsaList {
		lsas = append(lsas, &lsa)
	}

	return lsas
}

func NetworkConfigDummyData() *frrProto.NetworkConfig {
	return &frrProto.NetworkConfig{
		RouterId:   "65.0.1.1",
		Interfaces: OSPFInterfaceConfigDummyData(),
		Areas:      OSPFAreaDummyData(),
	}
}

func OSPFInterfaceConfigDummyData() []*frrProto.OSPFInterfaceConfig {
	var interfaceConfig []*frrProto.OSPFInterfaceConfig

	interfaceConfigList := []frrProto.OSPFInterfaceConfig{
		{
			Name:      "eth1",
			IpAddress: "172.22.1.1/24",
		},
		{
			Name:      "eth2",
			Area:      "0.0.0.0",
			IpAddress: "10.0.12.1/24",
		},
		{
			Name:      "eth3",
			Area:      "0.0.0.0",
			IpAddress: "10.0.13.1/24",
		},
		{
			Name:      "eth4",
			Area:      "0.0.0.0",
			IpAddress: "10.0.0.1/23",
			Passive:   true,
		},
		{
			Name:      "eth5",
			IpAddress: "192.168.100.1/24",
		},
		{
			Name:      "eth6",
			Area:      "0.0.0.0",
			IpAddress: "10.0.14.1/24",
		},
		{
			Name:      "eth7",
			Area:      "0.0.0.0",
			IpAddress: "10.0.15.1/24",
		},
		{
			Name:      "eth8",
			Area:      "0.0.0.0",
			IpAddress: "10.0.16.1/24",
		},
		{
			Name:      "eth9",
			Area:      "0.0.0.0",
			IpAddress: "10.0.17.1/24",
		},
		{
			Name:      "eth10",
			Area:      "0.0.0.0",
			IpAddress: "10.0.18.1/24",
		},
		{
			Name:      "eth11",
			Area:      "0.0.0.0",
			IpAddress: "10.0.19.1/24",
		},
		{
			Name:      "lo",
			IpAddress: "65.0.1.1/32",
			Passive:   true,
		},
	}

	for _, intConfig := range interfaceConfigList {
		interfaceConfig = append(interfaceConfig, &intConfig)
	}

	return interfaceConfig
}

func OSPFAreaDummyData() []*frrProto.OSPFArea {
	var ospfArea []*frrProto.OSPFArea

	ospfAreaList := []frrProto.OSPFArea{
		{
			Id: "0.0.0.0",
			Networks: []string{
				"10.0.12.1/24",
				"10.0.13.1/24",
				"10.0.0.1/23",
				"192.168.100.1/24",
				"10.0.14.1/24",
				"10.0.15.1/24",
				"10.0.16.1/24",
				"10.0.17.1/24",
				"10.0.18.1/24",
				"10.0.19.1/24",
			},
		},
	}

	for _, areaList := range ospfAreaList {
		ospfArea = append(ospfArea, &areaList)
	}

	return ospfArea
}

func SystemMetricsDummyData() *frrProto.SystemMetrics {
	interfaceStats := GetInterfaceStats()

	systemMetrics := frrProto.SystemMetrics{
		CpuUsage:     33.3,
		MemoryUsage:  20.2,
		NetworkStats: interfaceStats,
	}

	return &systemMetrics
}

func GetInterfaceStats() []*frrProto.InterfaceStats {
	var interfaceStats []*frrProto.InterfaceStats

	interfaceStatsList := []frrProto.InterfaceStats{
		{
			Name:      "eth2",
			RxBytes:   10000,
			TxBytes:   20000,
			RxErrors:  200,
			TxErrors:  300,
			OperState: "up",
		},
		{
			Name:      "eth3",
			RxBytes:   20000,
			TxBytes:   30000,
			RxErrors:  600,
			TxErrors:  900,
			OperState: "up",
		},
	}

	for _, interfaceStat := range interfaceStatsList {
		interfaceStats = append(interfaceStats, &interfaceStat)
	}

	return interfaceStats
}

func GetCombinedState() *frrProto.CombinedState {
	var timestamp timestamppb.Timestamp
	timestamp.Reset()
	combinedState := frrProto.CombinedState{
		Timestamp: &timestamp,
		Ospf:      OSPFMetricsDummyData(),
		Config:    NetworkConfigDummyData(),
		System:    SystemMetricsDummyData(),
	}

	return &combinedState
}

func OSPFMetricsDummyData() *frrProto.OSPFMetrics {
	ospfMetrics := frrProto.OSPFMetrics{
		Neighbors:  OSPFNeighborDummyData(),
		Routes:     OSPFRouteDummyData(),
		Interfaces: OSPFInterfaceDummyData(),
		Lsas:       OSPFlsaDummyData(),
	}

	return &ospfMetrics
}
