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
