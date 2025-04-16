package socket

import (
	"log"

	"github.com/ba2025-ysmprc/frr-tui/backend/internal/aggregator"
	ospfAnalyzer "github.com/ba2025-ysmprc/frr-tui/backend/internal/analyzer/ospf"
	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"
)

func processOSPFCommand(action string, params map[string]interface{}) (interface{}, error) {
	result := ospfAnalyzer.Dummy()
	return map[string]string{
		"result": result,
	}, nil
}

func processBGPCommand(action string, params map[string]interface{}) (interface{}, error) {
	return map[string]string{
		"result": "Neighbor added successfully",
	}, nil
}

func getOSPFMetrics() *frrProto.Response {
	var response frrProto.Response

	ospfMetrics := aggregator.OSPFMetricsDummyData()
	value := &frrProto.Value{
		Kind: &frrProto.Value_OspfMetrics{
			OspfMetrics: ospfMetrics,
		},
	}
	response.Status = "success"
	response.Message = "Returning all OSPF Metrics"
	response.Data = value

	return &response
}

func getOSPFNeighbor() *frrProto.Response {
	var response frrProto.Response

	var ospfMetrics frrProto.OSPFMetrics
	neighbors := aggregator.OSPFNeighborDummyData()
	ospfMetrics.Neighbors = neighbors
	value := &frrProto.Value{
		Kind: &frrProto.Value_OspfMetrics{
			OspfMetrics: &ospfMetrics,
		},
	}
	response.Status = "success"
	response.Message = "Returning neighbors of host"
	response.Data = value

	return &response
}

func getOSPFRoute() *frrProto.Response {
	var response frrProto.Response

	var ospfMetrics frrProto.OSPFMetrics
	routes := aggregator.OSPFRouteDummyData()
	ospfMetrics.Routes = routes
	value := &frrProto.Value{
		Kind: &frrProto.Value_OspfMetrics{
			OspfMetrics: &ospfMetrics,
		},
	}
	response.Status = "success"
	response.Message = "Returning routes of host"
	response.Data = value

	return &response
}

func getOSPFInterface() *frrProto.Response {
	var response frrProto.Response

	var ospfMetrics frrProto.OSPFMetrics
	interfaces := aggregator.OSPFInterfaceDummyData()
	ospfMetrics.Interfaces = interfaces
	value := &frrProto.Value{
		Kind: &frrProto.Value_OspfMetrics{
			OspfMetrics: &ospfMetrics,
		},
	}
	response.Status = "success"
	response.Message = "Returning interfaces of host"
	response.Data = value

	return &response
}

func getOSPFlsa() *frrProto.Response {
	var response frrProto.Response

	var ospfMetrics frrProto.OSPFMetrics
	lsas := aggregator.OSPFlsaDummyData()
	ospfMetrics.Lsas = lsas
	value := &frrProto.Value{
		Kind: &frrProto.Value_OspfMetrics{
			OspfMetrics: &ospfMetrics,
		},
	}
	response.Status = "success"
	response.Message = "Returning lsa of host"
	response.Data = value

	return &response
}

func getNetworkConfig() *frrProto.Response {
	var response frrProto.Response

	var networkConfig *frrProto.NetworkConfig
	networkConfig = aggregator.NetworkConfigDummyData()
	value := &frrProto.Value{
		Kind: &frrProto.Value_NetworkConfig{
			NetworkConfig: networkConfig,
		},
	}
	response.Status = "success"
	response.Message = "Returning network config of host"
	response.Data = value

	return &response
}

func getOSPFArea() *frrProto.Response {
	var response frrProto.Response

	var networkConfig frrProto.NetworkConfig
	networkConfig.Areas = aggregator.OSPFAreaDummyData()
	value := &frrProto.Value{
		Kind: &frrProto.Value_NetworkConfig{
			NetworkConfig: &networkConfig,
		},
	}
	response.Status = "success"
	response.Message = "Returning ospf area of host"
	response.Data = value

	return &response
}

func getInterfaceConfig() *frrProto.Response {
	var response frrProto.Response

	var networkConfig frrProto.NetworkConfig
	networkConfig.Interfaces = aggregator.OSPFInterfaceConfigDummyData()
	value := &frrProto.Value{
		Kind: &frrProto.Value_NetworkConfig{
			NetworkConfig: &networkConfig,
		},
	}
	response.Status = "success"
	response.Message = "Returning ospf area of host"
	response.Data = value

	return &response
}

func getSystemMetrics() *frrProto.Response {
	var response frrProto.Response

	systemMetrics := aggregator.SystemMetricsDummyData()
	value := &frrProto.Value{
		Kind: &frrProto.Value_SystemMetrics{
			SystemMetrics: systemMetrics,
		},
	}
	response.Status = "success"
	response.Message = "Returning system metrics of host"
	response.Data = value

	return &response
}

func getInterfaceStats() *frrProto.Response {
	var response frrProto.Response

	interfaceStats := aggregator.GetInterfaceStats()
	value := &frrProto.Value{
		Kind: &frrProto.Value_InterfaceStats{
			InterfaceStats: interfaceStats[0],
		},
	}
	response.Status = "success"
	response.Message = "Returning interface stat of host"
	response.Data = value

	return &response
}

func getCombinedState() *frrProto.Response {
	var response frrProto.Response

	combinedState := aggregator.GetCombinedState()
	value := &frrProto.Value{
		Kind: &frrProto.Value_CombinedState{
			CombinedState: combinedState,
		},
	}
	response.Status = "success"
	response.Message = "Returning combined state of host"
	response.Data = value

	return &response
}

func (s *Socket) getTesting() *frrProto.Response {
	var response frrProto.Response

	state, err := s.collector.Collect()
	if err != nil {
		log.Printf("Collection error: %v", err)
	}
	value := &frrProto.Value{
		Kind: &frrProto.Value_CombinedState{
			CombinedState: state,
		},
	}
	response.Status = "success"
	response.Message = "Returning combined state of host"
	response.Data = value

	return &response
}

func (s *Socket) getTesting2() *frrProto.Response {
	var response frrProto.Response

	state, err := s.collector.Collect()
	if err != nil {
		log.Printf("Collection error: %v", err)
	}
	value := &frrProto.Value{
		Kind: &frrProto.Value_OspfMetrics{
			OspfMetrics: state.GetOspf(),
		},
	}
	response.Status = "success"
	response.Message = "Returning combined state of host"
	response.Data = value

	return &response
}
func (s *Socket) getTesting3() *frrProto.Response {
	var response frrProto.Response

	value := &frrProto.Value{
		Kind: &frrProto.Value_StringValue{
			StringValue: s.analyzer.Foobar(),
		},
	}
	response.Status = "success"
	response.Message = "Returning combined state of host"
	response.Data = value

	return &response
}

func (s *Socket) getTesting4() *frrProto.Response {
	var response frrProto.Response

	staticConfig, err := s.collector.ReadConfig()
	if err != nil {
		response.Status = "error"
		response.Message = err.Error()
		return &response
	}
	value := &frrProto.Value{
		Kind: &frrProto.Value_StringValue{
			StringValue: staticConfig,
		},
	}
	response.Status = "success"
	response.Message = "Returning string config fileof host"
	response.Data = value

	return &response
}
