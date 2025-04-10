package aggregator

import "time"

type CombinedState struct {
	Timestamp time.Time
	OSPF      *OSPFMetrics
	Config    *NetworkConfig
	System    *SystemMetrics
}

// Maybe also has changes for route metrics
type OSPFMetrics struct {
	Neighbors  []OSPFNeighbor
	Routes     []OSPFRoute
	Interfaces []OSPFInterface
	LSAs       []OSPFlsa
}

type OSPFNeighbor struct {
	ID        string `json:"neighborId"`
	IP        string `json:"ipAddress"`
	State     string `json:"state"`
	Interface string `json:"interface"`
	Area      string `json:"area"`
}

type OSPFRoute struct {
	Prefix    string `json:"prefix"`
	NextHop   string `json:"nextHop"`
	Interface string `json:"interface"`
	Cost      int    `json:"cost"`
	Type      string `json:"type"`
	Area      string `json:"area"`
}

type OSPFInterface struct {
	Name     string `json:"name"`
	Area     string `json:"area"`
	NbrCount int    `json:"nbrCount"`
	NbrAdj   int    `json:"nbrAdjacentCount"`
	Passive  bool   `json:"passive"`
}

type OSPFlsa struct {
	Type      string `json:"type"`
	LSID      string `json:"lsId"`
	AdvRouter string `json:"advRouter"`
	Age       int    `json:"age"`
	Area      string `json:"area"`
}

type NetworkConfig struct {
	RouterID   string
	Areas      []OSPFArea
	Interfaces []OSPFInterfaceConfig
}

type OSPFArea struct {
	ID       string
	Networks []string
}

type OSPFInterfaceConfig struct {
	Name    string
	Area    string
	Passive bool
	Cost    int
}

type SystemMetrics struct {
	CPUUsage     float64
	MemoryUsage  float64
	NetworkStats []InterfaceStats
}

type InterfaceStats struct {
	Name      string
	RxBytes   uint64
	TxBytes   uint64
	RxErrors  uint64
	TxErrors  uint64
	OperState string
}
