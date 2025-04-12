package socket

import (
	ospfAnalyzer "github.com/ba2025-ysmprc/frr-tui/backend/internal/analyzer/ospf"
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
