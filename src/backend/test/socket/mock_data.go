package socket_test

import (
	"github.com/frr-mad/frr-mad/src/backend/internal/analyzer"
	"github.com/frr-mad/frr-mad/src/backend/internal/configs"
	"github.com/frr-mad/frr-mad/src/backend/internal/socket"
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
)

// CreateMockOSPFDatabase creates a simple mock OSPFDatabase for testing
func CreateMockOSPFDatabase() *frrProto.OSPFDatabase {
	return &frrProto.OSPFDatabase{
		RouterId:        "192.168.1.1",
		AsExternalCount: 1,
		Areas: map[string]*frrProto.OSPFDatabaseArea{
			"0.0.0.0": {
				RouterLinkStatesCount: 1,
				RouterLinkStates: []*frrProto.RouterDataLSA{
					{
						Base: &frrProto.BaseLSA{
							LsId:             "192.168.1.1",
							AdvertisedRouter: "192.168.1.1",
							LsaAge:           60,
							SequenceNumber:   "80000001",
							Checksum:         "A1B2",
						},
						NumOfRouterLinks: 1,
					},
				},
			},
		},
		AsExternalLinkStates: []*frrProto.ASExternalLSA{
			{
				Base: &frrProto.BaseLSA{
					LsId:             "192.168.1.0",
					AdvertisedRouter: "192.168.1.1",
					LsaAge:           60,
					SequenceNumber:   "80000001",
					Checksum:         "A1B2",
				},
				MetricType: "E2",
				Route:      "192.168.1.0/24",
			},
		},
	}
}

func CreateMockFRRRouterData() *frrProto.FRRRouterData {
	return &frrProto.FRRRouterData{
		RouterName:   "r101",
		OspfRouterId: "192.168.1.1",
	}
}

// CreateMockOSPFRouterData creates a simple mock OSPFRouterData for testing
func CreateMockOSPFRouterData() *frrProto.OSPFRouterData {
	return &frrProto.OSPFRouterData{
		RouterId: "192.168.1.1",
		RouterStates: map[string]*frrProto.OSPFRouterArea{
			"0.0.0.0": {
				LsaEntries: map[string]*frrProto.OSPFRouterLSA{
					"192.168.1.1": {
						LsaAge:            60,
						Options:           "E",
						LsaType:           "router-LSA",
						LinkStateId:       "192.168.1.1",
						AdvertisingRouter: "192.168.1.1",
						LsaSeqNumber:      "80000001",
						Checksum:          "A1B2",
						Length:            36,
						NumOfLinks:        1,
						RouterLinks: map[string]*frrProto.OSPFRouterLSALink{
							"stub": {
								LinkType:        "stub",
								NetworkAddress:  "192.168.1.0",
								NetworkMask:     "255.255.255.0",
								NumOfTosMetrics: 0,
								Tos0Metric:      10,
							},
						},
					},
				},
			},
		},
	}
}

// CreateMockOSPFNeighbors creates a simple mock OSPFNeighbors for testing
func CreateMockOSPFNeighbors() *frrProto.OSPFNeighbors {
	return &frrProto.OSPFNeighbors{
		Neighbors: map[string]*frrProto.NeighborList{
			"eth0": {
				Neighbors: []*frrProto.Neighbor{
					{
						Priority:      1,
						State:         "Full",
						NbrPriority:   1,
						NbrState:      "Full",
						Converged:     "Complete",
						Role:          "DR",
						UpTimeInMsec:  3600000,
						DeadTimeMsecs: 40000,
						UpTime:        "01:00:00",
						DeadTime:      "00:00:40",
						Address:       "192.168.1.2",
						IfaceAddress:  "192.168.1.1",
						IfaceName:     "eth0",
					},
				},
			},
		},
	}
}

func CreateMockAnomalyDetectionRouter() *frrProto.AnomalyDetection {
	return &frrProto.AnomalyDetection{
		HasOverAdvertisedPrefixes: true,
		HasUnAdvertisedPrefixes:   false,
		HasDuplicatePrefixes:      false,
		HasMisconfiguredPrefixes:  false,
		SuperfluousEntries: []*frrProto.Advertisement{
			{
				InterfaceAddress: "89.207.132.170",
				LinkStateId:      "89.207.132.170/24",
				PrefixLength:     "24",
				LinkType:         "transit network",
			},
		},
		MissingEntries: []*frrProto.Advertisement{},
	}
}

func CreateMockAnomalyDetectionExternal() *frrProto.AnomalyDetection {
	return &frrProto.AnomalyDetection{
		HasOverAdvertisedPrefixes: false,
		HasUnAdvertisedPrefixes:   true,
		HasDuplicatePrefixes:      false,
		HasMisconfiguredPrefixes:  false,
		SuperfluousEntries:        []*frrProto.Advertisement{},
		MissingEntries: []*frrProto.Advertisement{
			{
				InterfaceAddress: "89.207.132.170",
				LinkStateId:      "89.207.132.170/24",
				PrefixLength:     "24",
				LinkType:         "transit network",
			},
		},
		DuplicateEntries: []*frrProto.Advertisement{},
	}
}

func CreateMockAnomalyDetectionNssaExternal() *frrProto.AnomalyDetection {
	return &frrProto.AnomalyDetection{
		HasOverAdvertisedPrefixes: true,
		HasUnAdvertisedPrefixes:   false,
		HasDuplicatePrefixes:      false,
		HasMisconfiguredPrefixes:  false,
		SuperfluousEntries: []*frrProto.Advertisement{
			{
				InterfaceAddress: "89.207.132.170",
				LinkStateId:      "89.207.132.170/24",
				PrefixLength:     "24",
				LinkType:         "transit network",
			},
		},
		MissingEntries:   []*frrProto.Advertisement{},
		DuplicateEntries: []*frrProto.Advertisement{},
	}
}

func CreateMockAnomalyDetectionLsdbToRib() *frrProto.AnomalyDetection {
	return &frrProto.AnomalyDetection{
		HasOverAdvertisedPrefixes: false,
		HasUnAdvertisedPrefixes:   false,
		HasDuplicatePrefixes:      false,
		HasMisconfiguredPrefixes:  false,
		SuperfluousEntries:        []*frrProto.Advertisement{},
		MissingEntries:            []*frrProto.Advertisement{},
		DuplicateEntries:          []*frrProto.Advertisement{},
	}
}

// CreateMockSystemMetrics creates a simple mock SystemMetrics for testing
func CreateMockSystemMetrics() *frrProto.SystemMetrics {
	return &frrProto.SystemMetrics{
		CpuAmount:   4,
		CpuUsage:    25.5,
		MemoryUsage: 40.2,
	}
}

// CreateMockFullFRRData creates a simple mock FullFRRData for testing
func CreateMockFullFRRData() *frrProto.FullFRRData {
	return &frrProto.FullFRRData{
		OspfDatabase:   CreateMockOSPFDatabase(),
		OspfRouterData: CreateMockOSPFRouterData(),
		OspfNeighbors:  CreateMockOSPFNeighbors(),
		SystemMetrics:  CreateMockSystemMetrics(),
		FrrRouterData:  CreateMockFRRRouterData(),
	}
}

func getMockData() (*logger.Logger, *analyzer.Analyzer, *frrProto.FullFRRData, *frrProto.ParsedAnalyzerData) {
	mockLoggerInstance, _ := logger.NewApplicationLogger("testing", "/tmp/testing.log")
	mockMetrics := CreateMockFullFRRData()
	mockAnalyzerInstance := &analyzer.Analyzer{
		AnalysisResult: &frrProto.AnomalyAnalysis{},
	}
	mockParsedAnalyzerdata := &frrProto.ParsedAnalyzerData{}

	return mockLoggerInstance, mockAnalyzerInstance, mockMetrics, mockParsedAnalyzerdata
}

func getEmptyMockSocket() *socket.Socket {
	socket := socket.NewSocket(
		configs.SocketConfig{},
		&frrProto.FullFRRData{},
		&frrProto.AnomalyAnalysis{},
		&logger.Logger{},
		&frrProto.ParsedAnalyzerData{})
	return socket
}
