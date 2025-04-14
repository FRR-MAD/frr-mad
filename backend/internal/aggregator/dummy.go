package aggregator

import frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"

func DummyFunction() string {
	return "I am aggregator Aggregator"
}

func OSPFNeighborDummyData() []*frrProto.OSPFNeighbor {
	var neighbors []*frrProto.OSPFNeighbor
	neighbor1 := frrProto.OSPFNeighbor{
		Id:        "65.0.1.8",
		Ip:        "10.0.18.8",
		Interface: "eth10",
		Area:      "0.0.0.0",
	}
	neighbor2 := frrProto.OSPFNeighbor{
		Id:        "65.0.1.9",
		Ip:        "10.0.19.9",
		Interface: "eth11",
		Area:      "0.0.0.0",
	}
	neighbor3 := frrProto.OSPFNeighbor{
		Id:        "65.0.1.2",
		Ip:        "10.0.12.2",
		Interface: "eth2",
		Area:      "0.0.0.0",
	}
	neighbor4 := frrProto.OSPFNeighbor{
		Id:        "65.0.1.3",
		Ip:        "10.0.13.3",
		Interface: "eth3",
		Area:      "0.0.0.0",
	}
	neighbor5 := frrProto.OSPFNeighbor{
		Id:        "65.0.1.4",
		Ip:        "10.0.14.4",
		Interface: "eth6",
		Area:      "0.0.0.0",
	}
	neighbor6 := frrProto.OSPFNeighbor{
		Id:        "65.0.1.5",
		Ip:        "10.0.15.5",
		Interface: "eth7",
		Area:      "0.0.0.0",
	}
	neighbor7 := frrProto.OSPFNeighbor{
		Id:        "65.0.1.6",
		Ip:        "10.0.16.6",
		Interface: "eth8",
		Area:      "0.0.0.0",
	}
	neighbor8 := frrProto.OSPFNeighbor{
		Id:        "65.0.1.7",
		Ip:        "10.0.17.7",
		Interface: "eth9",
		Area:      "0.0.0.0",
	}
	neighbors = append(neighbors, &neighbor1)
	neighbors = append(neighbors, &neighbor2)
	neighbors = append(neighbors, &neighbor3)
	neighbors = append(neighbors, &neighbor4)
	neighbors = append(neighbors, &neighbor5)
	neighbors = append(neighbors, &neighbor6)
	neighbors = append(neighbors, &neighbor7)
	neighbors = append(neighbors, &neighbor8)

	return neighbors
}

func OSPFRouteDummyData() []*frrProto.OSPFRoute {
	var routes []*frrProto.OSPFRoute
	route1 := frrProto.OSPFRoute{
		Prefix:    "192.168.8.0/24",
		NextHop:   "10.0.18.8",
		Interface: "eth10",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route2 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.8",
		NextHop:   "10.0.18.8",
		Interface: "eth10",
		Cost:      10,
		Type:      "R",
		Area:      "0.0.0.0",
	}
	route3 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.8/32",
		NextHop:   "10.0.18.8",
		Interface: "eth10",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route4 := frrProto.OSPFRoute{
		Prefix:    "10.0.18.0/24",
		NextHop:   "direct",
		Interface: "eth10",
		Cost:      10,
		Type:      "N",
		Area:      "0.0.0.0",
	}
	route5 := frrProto.OSPFRoute{
		Prefix:    "192.168.9.0/24",
		NextHop:   "10.0.19.9",
		Interface: "eth11",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route6 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.9",
		NextHop:   "10.0.19.9",
		Interface: "eth11",
		Cost:      10,
		Type:      "R",
		Area:      "0.0.0.0",
	}
	route7 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.9/32",
		NextHop:   "10.0.19.9",
		Interface: "eth11",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route8 := frrProto.OSPFRoute{
		Prefix:    "10.0.19.0/24",
		NextHop:   "direct",
		Interface: "eth11",
		Cost:      10,
		Type:      "N",
		Area:      "0.0.0.0",
	}
	route9 := frrProto.OSPFRoute{
		Prefix:    "10.0.23.0/24",
		NextHop:   "10.0.12.2",
		Interface: "eth2",
		Cost:      20,
		Type:      "N",
		Area:      "0.0.0.0",
	}
	route10 := frrProto.OSPFRoute{
		Prefix:    "10.1.0.0/24",
		NextHop:   "10.0.12.2",
		Interface: "eth2",
		Cost:      30,
		Type:      "N IA",
		Area:      "0.0.0.0",
	}
	route11 := frrProto.OSPFRoute{
		Prefix:    "10.1.12.0/24",
		NextHop:   "10.0.12.2",
		Interface: "eth2",
		Cost:      30,
		Type:      "N IA",
		Area:      "0.0.0.0",
	}
	route12 := frrProto.OSPFRoute{
		Prefix:    "10.1.21.0/24",
		NextHop:   "10.0.12.2",
		Interface: "eth2",
		Cost:      20,
		Type:      "N IA",
		Area:      "0.0.0.0",
	}
	route13 := frrProto.OSPFRoute{
		Prefix:    "10.30.12.0/24",
		NextHop:   "10.0.12.2",
		Interface: "eth2",
		Cost:      50,
		Type:      "N E1",
		Area:      "0.0.0.0",
	}
	route14 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.2",
		NextHop:   "10.0.12.2",
		Interface: "eth2",
		Cost:      10,
		Type:      "R",
		Area:      "0.0.0.0",
	}
	route15 := frrProto.OSPFRoute{
		Prefix:    "65.0.3.1/32",
		NextHop:   "10.0.12.2",
		Interface: "eth2",
		Cost:      50,
		Type:      "N E1",
		Area:      "0.0.0.0",
	}
	route16 := frrProto.OSPFRoute{
		Prefix:    "10.0.12.0/24",
		NextHop:   "direct",
		Interface: "eth2",
		Cost:      10,
		Type:      "N",
		Area:      "0.0.0.0",
	}
	route17 := frrProto.OSPFRoute{
		Prefix:    "10.2.0.0/24",
		NextHop:   "10.0.13.3",
		Interface: "eth3",
		Cost:      30,
		Type:      "N IA",
		Area:      "0.0.0.0",
	}
	route18 := frrProto.OSPFRoute{
		Prefix:    "10.2.12.0/24",
		NextHop:   "10.0.13.3",
		Interface: "eth3",
		Cost:      30,
		Type:      "N IA",
		Area:      "0.0.0.0",
	}
	route19 := frrProto.OSPFRoute{
		Prefix:    "10.2.31.0/24",
		NextHop:   "10.0.13.3",
		Interface: "eth3",
		Cost:      20,
		Type:      "N IA",
		Area:      "0.0.0.0",
	}
	route20 := frrProto.OSPFRoute{
		Prefix:    "10.3.0.0/24",
		NextHop:   "10.0.13.3",
		Interface: "eth3",
		Cost:      50,
		Type:      "N IA",
		Area:      "0.0.0.0",
	}
	route21 := frrProto.OSPFRoute{
		Prefix:    "10.3.21.0/24",
		NextHop:   "10.0.13.3",
		Interface: "eth3",
		Cost:      40,
		Type:      "N IA",
		Area:      "0.0.0.0",
	}
	route22 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.22",
		NextHop:   "10.0.13.3",
		Interface: "eth3",
		Cost:      30,
		Type:      "R",
		Area:      "0.0.0.0",
	}
	route23 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.3",
		NextHop:   "10.0.13.3",
		Interface: "eth3",
		Cost:      10,
		Type:      "R",
		Area:      "0.0.0.0",
	}
	route24 := frrProto.OSPFRoute{
		Prefix:    "10.0.13.0/24",
		NextHop:   "direct",
		Interface: "eth3",
		Cost:      10,
		Type:      "N",
		Area:      "0.0.0.0",
	}
	route25 := frrProto.OSPFRoute{
		Prefix:    "10.0.0.0/23",
		NextHop:   "direct",
		Interface: "eth4",
		Cost:      10,
		Type:      "N",
		Area:      "0.0.0.0",
	}
	route26 := frrProto.OSPFRoute{
		Prefix:    "172.20.20.0/24",
		NextHop:   "10.0.14.4",
		Interface: "eth6",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route27 := frrProto.OSPFRoute{
		Prefix:    "192.168.4.0/24",
		NextHop:   "10.0.14.4",
		Interface: "eth6",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route28 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.4",
		NextHop:   "10.0.14.4",
		Interface: "eth6",
		Cost:      10,
		Type:      "R",
		Area:      "0.0.0.0",
	}
	route29 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.4/32",
		NextHop:   "10.0.14.4",
		Interface: "eth6",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route30 := frrProto.OSPFRoute{
		Prefix:    "10.0.14.0/24",
		NextHop:   "direct",
		Interface: "eth6",
		Cost:      10,
		Type:      "N",
		Area:      "0.0.0.0",
	}
	route31 := frrProto.OSPFRoute{
		Prefix:    "192.168.5.0/24",
		NextHop:   "10.0.15.5",
		Interface: "eth7",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route32 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.5",
		NextHop:   "10.0.15.5",
		Interface: "eth7",
		Cost:      10,
		Type:      "R",
		Area:      "0.0.0.0",
	}
	route33 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.5/32",
		NextHop:   "10.0.15.5",
		Interface: "eth7",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route34 := frrProto.OSPFRoute{
		Prefix:    "10.0.15.0/24",
		NextHop:   "direct",
		Interface: "eth7",
		Cost:      10,
		Type:      "N",
		Area:      "0.0.0.0",
	}
	route35 := frrProto.OSPFRoute{
		Prefix:    "192.168.6.0/24",
		NextHop:   "10.0.16.6",
		Interface: "eth8",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route36 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.6",
		NextHop:   "10.0.16.6",
		Interface: "eth8",
		Cost:      10,
		Type:      "R",
		Area:      "0.0.0.0",
	}
	route37 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.6/32",
		NextHop:   "10.0.16.6",
		Interface: "eth8",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route38 := frrProto.OSPFRoute{
		Prefix:    "10.0.16.0/24",
		NextHop:   "direct",
		Interface: "eth8",
		Cost:      10,
		Type:      "N",
		Area:      "0.0.0.0",
	}
	route39 := frrProto.OSPFRoute{
		Prefix:    "192.168.7.0/24",
		NextHop:   "10.0.17.7",
		Interface: "eth9",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route40 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.7",
		NextHop:   "10.0.17.7",
		Interface: "eth9",
		Cost:      10,
		Type:      "R",
		Area:      "0.0.0.0",
	}
	route41 := frrProto.OSPFRoute{
		Prefix:    "65.0.1.7/32",
		NextHop:   "10.0.17.7",
		Interface: "eth9",
		Cost:      10,
		Type:      "N E2",
		Area:      "0.0.0.0",
	}
	route42 := frrProto.OSPFRoute{
		Prefix:    "10.0.17.0/24",
		NextHop:   "direct",
		Interface: "eth9",
		Cost:      10,
		Type:      "N",
		Area:      "0.0.0.0",
	}
	routes = append(routes, &route1)
	routes = append(routes, &route1)
	routes = append(routes, &route2)
	routes = append(routes, &route3)
	routes = append(routes, &route4)
	routes = append(routes, &route5)
	routes = append(routes, &route6)
	routes = append(routes, &route7)
	routes = append(routes, &route8)
	routes = append(routes, &route9)
	routes = append(routes, &route10)
	routes = append(routes, &route11)
	routes = append(routes, &route12)
	routes = append(routes, &route13)
	routes = append(routes, &route14)
	routes = append(routes, &route15)
	routes = append(routes, &route16)
	routes = append(routes, &route17)
	routes = append(routes, &route18)
	routes = append(routes, &route19)
	routes = append(routes, &route20)
	routes = append(routes, &route21)
	routes = append(routes, &route22)
	routes = append(routes, &route23)
	routes = append(routes, &route24)
	routes = append(routes, &route25)
	routes = append(routes, &route26)
	routes = append(routes, &route27)
	routes = append(routes, &route28)
	routes = append(routes, &route29)
	routes = append(routes, &route30)
	routes = append(routes, &route31)
	routes = append(routes, &route32)
	routes = append(routes, &route33)
	routes = append(routes, &route34)
	routes = append(routes, &route35)
	routes = append(routes, &route36)
	routes = append(routes, &route37)
	routes = append(routes, &route38)
	routes = append(routes, &route39)
	routes = append(routes, &route40)
	routes = append(routes, &route41)
	routes = append(routes, &route42)

	return routes
}

func OSPFInterfaceDummyData() []*frrProto.OSPFInterface {
	var interfaces []*frrProto.OSPFInterface
	interface1 := frrProto.OSPFInterface{
		Name:     "eth2",
		Area:     "0.0.0.0",
		NbrCount: 1,
		NbrAdj:   1,
		Passive:  false,
	}
	interface2 := frrProto.OSPFInterface{
		Name:     "eth3",
		Area:     "0.0.0.0",
		NbrCount: 1,
		NbrAdj:   1,
		Passive:  false,
	}
	interface3 := frrProto.OSPFInterface{
		Name:     "eth4",
		Area:     "0.0.0.0",
		NbrCount: 0,
		NbrAdj:   0,
		Passive:  true,
	}
	interface4 := frrProto.OSPFInterface{
		Name:     "eth6",
		Area:     "0.0.0.0",
		NbrCount: 1,
		NbrAdj:   1,
		Passive:  false,
	}
	interface5 := frrProto.OSPFInterface{
		Name:     "eth7",
		Area:     "0.0.0.0",
		NbrCount: 1,
		NbrAdj:   1,
		Passive:  false,
	}
	interface6 := frrProto.OSPFInterface{
		Name:     "eth8",
		Area:     "0.0.0.0",
		NbrCount: 1,
		NbrAdj:   1,
		Passive:  false,
	}
	interface7 := frrProto.OSPFInterface{
		Name:     "eth9",
		Area:     "0.0.0.0",
		NbrCount: 1,
		NbrAdj:   1,
		Passive:  false,
	}
	interface8 := frrProto.OSPFInterface{
		Name:     "eth10",
		Area:     "0.0.0.0",
		NbrCount: 1,
		NbrAdj:   1,
		Passive:  false,
	}
	interface9 := frrProto.OSPFInterface{
		Name:     "eth11",
		Area:     "0.0.0.0",
		NbrCount: 1,
		NbrAdj:   1,
		Passive:  false,
	}

	interfaces = append(interfaces, &interface1)
	interfaces = append(interfaces, &interface2)
	interfaces = append(interfaces, &interface3)
	interfaces = append(interfaces, &interface4)
	interfaces = append(interfaces, &interface5)
	interfaces = append(interfaces, &interface6)
	interfaces = append(interfaces, &interface7)
	interfaces = append(interfaces, &interface8)
	interfaces = append(interfaces, &interface9)
	return interfaces
}


func OSPFlsaDummyData() []*frrProto.OSPFlsa {
	var interfaces []*frrProto.OSPFlsa

	type:"external"
	ls_id:"10.20.0.0"
	adv_router:"65.0.1.1"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"10.20.12.0"
	adv_router:"65.0.1.1"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"192.168.1.0"
	adv_router:"65.0.1.1"
	Sequence:"0"
	area:"0.0.0.0"
	type:"router"
	ls_id:"65.0.1.1"
	adv_router:"65.0.1.1"
	Sequence:"80000031"
	area:"0.0.0.0"
	type:"external"
	ls_id:"65.0.2.1"
	adv_router:"65.0.1.1"
	Sequence:"0"
	area:"0.0.0.0"
	type:"network"
	ls_id:"10.0.12.2"
	adv_router:"65.0.1.2"
	Sequence:"0"
	area:"0.0.0.0"
	type:"summary"
	ls_id:"10.1.0.0"
	adv_router:"65.0.1.2"
	Sequence:"0"
	area:"0.0.0.0"
	type:"summary"
	ls_id:"10.1.12.0"
	adv_router:"65.0.1.2"
	Sequence:"0"
	area:"0.0.0.0"
	type:"summary"
	ls_id:"10.1.21.0"
	adv_router:"65.0.1.2"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"10.30.0.0"
	adv_router:"65.0.1.2"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"10.30.12.0"
	adv_router:"65.0.1.2"
	Sequence:"0"
	area:"0.0.0.0"
	type:"router"
	ls_id:"65.0.1.2"
	adv_router:"65.0.1.2"
	Sequence:"80000015"
	area:"0.0.0.0"
	type:"external"
	ls_id:"65.0.3.1"
	adv_router:"65.0.1.2"
	Sequence:"0"
	area:"0.0.0.0"
	type:"summary"
	ls_id:"10.2.0.0"
	adv_router:"65.0.1.22"
	Sequence:"0"
	area:"0.0.0.0"
	type:"summary"
	ls_id:"10.2.12.0"
	adv_router:"65.0.1.22"
	Sequence:"0"
	area:"0.0.0.0"
	type:"summary"
	ls_id:"10.2.31.0"
	adv_router:"65.0.1.22"
	Sequence:"0"
	area:"0.0.0.0"
	type:"summary"
	ls_id:"10.3.0.0"
	adv_router:"65.0.1.22"
	Sequence:"0"
	area:"0.0.0.0"
	type:"summary"
	ls_id:"10.3.21.0"
	adv_router:"65.0.1.22"
	Sequence:"0"
	area:"0.0.0.0"
	type:"router"
	ls_id:"65.0.1.22"
	adv_router:"65.0.1.22"
	Sequence:"80000011"
	area:"0.0.0.0"
	type:"network"
	ls_id:"10.0.13.3"
	adv_router:"65.0.1.3"
	Sequence:"0"
	area:"0.0.0.0"
	type:"network"
	ls_id:"10.0.23.3"
	adv_router:"65.0.1.3"
	Sequence:"0"
	area:"0.0.0.0"
	type:"summary"
	ls_id:"10.2.0.0"
	adv_router:"65.0.1.3"
	Sequence:"0"
	area:"0.0.0.0"
	type:"summary"
	ls_id:"10.2.12.0"
	adv_router:"65.0.1.3"
	Sequence:"0"
	area:"0.0.0.0"
	type:"summary"
	ls_id:"10.2.31.0"
	adv_router:"65.0.1.3"
	Sequence:"0"
	area:"0.0.0.0"
	type:"router"
	ls_id:"65.0.1.3"
	adv_router:"65.0.1.3"
	Sequence:"80000019"
	area:"0.0.0.0"
	type:"network"
	ls_id:"10.0.14.4"
	adv_router:"65.0.1.4"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"172.20.20.0"
	adv_router:"65.0.1.4"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"192.168.4.0"
	adv_router:"65.0.1.4"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"65.0.1.4"
	adv_router:"65.0.1.4"
	Sequence:"0"
	area:"0.0.0.0"
	type:"router"
	ls_id:"65.0.1.4"
	adv_router:"65.0.1.4"
	Sequence:"80000010"
	area:"0.0.0.0"
	type:"network"
	ls_id:"10.0.15.5"
	adv_router:"65.0.1.5"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"172.20.20.0"
	adv_router:"65.0.1.5"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"192.168.5.0"
	adv_router:"65.0.1.5"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"65.0.1.5"
	adv_router:"65.0.1.5"
	Sequence:"0"
	area:"0.0.0.0"
	type:"router"
	ls_id:"65.0.1.5"
	adv_router:"65.0.1.5"
	Sequence:"80000011"
	area:"0.0.0.0"
	type:"network"
	ls_id:"10.0.16.6"
	adv_router:"65.0.1.6"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"172.20.20.0"
	adv_router:"65.0.1.6"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"192.168.6.0"
	adv_router:"65.0.1.6"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"65.0.1.6"
	adv_router:"65.0.1.6"
	Sequence:"0"
	area:"0.0.0.0"
	type:"router"
	ls_id:"65.0.1.6"
	adv_router:"65.0.1.6"
	Sequence:"80000012"
	area:"0.0.0.0"
	type:"network"
	ls_id:"10.0.17.7"
	adv_router:"65.0.1.7"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"172.20.20.0"
	adv_router:"65.0.1.7"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"192.168.7.0"
	adv_router:"65.0.1.7"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"65.0.1.7"
	adv_router:"65.0.1.7"
	Sequence:"0"
	area:"0.0.0.0"
	type:"router"
	ls_id:"65.0.1.7"
	adv_router:"65.0.1.7"
	Sequence:"80000011"
	area:"0.0.0.0"
	type:"network"
	ls_id:"10.0.18.8"
	adv_router:"65.0.1.8"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"172.20.20.0"
	adv_router:"65.0.1.8"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"192.168.8.0"
	adv_router:"65.0.1.8"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"65.0.1.8"
	adv_router:"65.0.1.8"
	Sequence:"0"
	area:"0.0.0.0"
	type:"router"
	ls_id:"65.0.1.8"
	adv_router:"65.0.1.8"
	Sequence:"80000013"
	area:"0.0.0.0"
	type:"network"
	ls_id:"10.0.19.9"
	adv_router:"65.0.1.9"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"172.20.20.0"
	adv_router:"65.0.1.9"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"192.168.9.0"
	adv_router:"65.0.1.9"
	Sequence:"0"
	area:"0.0.0.0"
	type:"external"
	ls_id:"65.0.1.9"
	adv_router:"65.0.1.9"
	Sequence:"0"
	area:"0.0.0.0"
	type:"router"
	ls_id:"65.0.1.9"
	adv_router:"65.0.1.9"
	Sequence:"80000010"
	area:"0.0.0.0"]

	return interfaces
}