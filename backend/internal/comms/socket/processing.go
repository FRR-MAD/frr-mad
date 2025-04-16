package socket

import (
	"fmt"
	"time"

	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"
)

// func processCommand(message *frrProto.Message) *frrProto.Response {
func (s *Socket) processCommand(message *frrProto.Message) *frrProto.Response {
	var response frrProto.Response

	switch message.Command {
	case "ospf":
		switch message.Package {
		case "metrics":
			return getOSPFMetrics()
		case "neighbor":
			return getOSPFNeighbor()
		case "route":
			return getOSPFRoute()
		case "interface":
			return getOSPFInterface()
		case "lsa":
			return getOSPFlsa()
		case "networkConfig":
			return getNetworkConfig()
		case "area":
			return getOSPFArea()
		case "interfaceConfig":
			return getInterfaceConfig()
		case "systemMetrics":
			return getSystemMetrics()
		case "interfaceStats":
			return getInterfaceStats()
		case "combinedState":
			return getCombinedState()
		case "testing":
			return s.getTesting()
		case "testing2":
			return s.getTesting2()
		case "testing3":
			return s.getTesting3()
		}

		response.Status = "success"
		response.Message = "Returning magical ospf data"
		//response.Data = value
		// aggregator.OSPFNeighborDummyData()
		return &response
		// return &response
	case "exit":
		response.Status = "success"
		response.Message = "Shutting system down"
		go func() {
			time.Sleep(100 * time.Millisecond)
			exitSocketServer()
		}()
		return &response
	default:
		response.Status = "error"
		response.Message = fmt.Sprintf("Unknown command: %s", message.Command)
		return &response
	}

}
