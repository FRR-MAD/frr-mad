package socket

import (
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

	var ospfMetrics frrProto.OSPFMetrics
	neighbors := aggregator.OSPFNeighborDummyData()
	ospfMetrics.Neighbors = neighbors
	// Create the Value with the OSPFMetrics field set
	value := &frrProto.Value{
		Kind: &frrProto.Value_OspfMetrics{
			OspfMetrics: &ospfMetrics,
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

	var networkConfig frrProto.NetworkConfig
	//networkConfig := aggregator.NetworkConfigDummyData()
	value := &frrProto.Value{
		Kind: &frrProto.Value_NetworkConfig{
			NetworkConfig: &networkConfig,
		},
	}
	response.Status = "success"
	response.Message = "Returning network config of host"
	response.Data = value

	return &response
}

func get() *frrProto.Response {
	var response frrProto.Response

	return &response
}
