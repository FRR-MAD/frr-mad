package socket

import (
	"fmt"
	"time"

	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"
)

// func processCommand(message *frrProto.Message) *frrProto.Response {
func processCommand(message *frrProto.Message) *frrProto.Response {
	var response frrProto.Response

	switch message.Command {
	case "ospf":
		//var ospfMetrics frrProto.OSPFMetrics
		//neighbors := aggregator.OSPFNeighborDummyData()
		//ospfMetrics.Neighbors = neighbors
		//// Create the Value with the OSPFMetrics field set
		//value := &frrProto.Value{
		//	Kind: &frrProto.Value_OspfMetrics{
		//		OspfMetrics: &ospfMetrics,
		//	},
		//}

		switch message.Package {
		case "neighbors":

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
