package analyzer

import (
	"fmt"

	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"
)

// TODO: import aggregator
// import "github.com/ba2025-ysmprc/frr-tui/backend/aggregator"

// only one function is necessary, all the rest will do the same?
func (m *messageList) updateOSPFMetrics() {
	var OSPFMetrics *frrProto.OSPFMetrics

	fmt.Println(OSPFMetrics)
}

func (m *messageList) getOSPFNeighbor() frrProto.OSPFNeighbor {
	var retValue frrProto.OSPFNeighbor
	retValue = *m.OSPFNeighbor
	// Create a deep copy of the OSPFNeighbor struct
	return retValue
}

// func (m *messageList) updateOSPFRoute() {
// var OSPFRoute *frrProto.OSPFRoute
// }
//
// func (m *messageList) updateOSPFInterface() {
// var OSPFInterface *frrProto.OSPFInterface
// }
//
// func (m *messageList) updateOSPFlsa() {
// return m.OSPFlsa
// var OSPFlsa frrProto.OSPFlsa
// }
//
// func (m *messageList) updateNetworkConfig() {
// var NetworkConfig *frrProto.NetworkConfig
// }
//
// func (m *messageList) updateOSPFArea() {
// var OSPFArea *frrProto.OSPFArea
// }
//
// func (m *messageList) updateOSPFInterfaceConfig() {
// var OSPFInterfaceConfig *frrProto.OSPFInterfaceConfig
// }
//
// func (m *messageList) updateSystemMetrics() {
// var SystemMetrics *frrProto.SystemMetrics
// }
//
// func (m *messageList) updateInterfaceStats() {
// var InterfaceStats *frrProto.InterfaceStats
// }
//
