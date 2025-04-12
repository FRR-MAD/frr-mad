package analyzer

import (
	"fmt"

	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"
)

/*

 */

type messageList struct {
	OSPFMetrics         *frrProto.OSPFMetrics
	OSPFNeighbor        *frrProto.OSPFNeighbor
	OSPFRoute           *frrProto.OSPFRoute
	OSPFInterface       *frrProto.OSPFInterface
	OSPFlsa             *frrProto.OSPFlsa
	Networkconfig       *frrProto.NetworkConfig
	OSPFArea            *frrProto.OSPFArea
	OSPFInterfaceConfig *frrProto.OSPFInterfaceConfig
	SystemMetrics       *frrProto.SystemMetrics
	InterfaceStats      *frrProto.InterfaceStats
}

func initializeMessageList() *messageList {
	return &messageList{}
}

func (m *messageList) updateMessageList() {
	// call all fuctions to update values
	// should we do it goroutine and parallalize it?

}

func (m *messageList) updateMessageListSelected(messageList []string) {
	fmt.Printf("Create an interesting way, to update only individual metrics, maybe with case selector? Here is the messageList: %v\n", messageList)

}
