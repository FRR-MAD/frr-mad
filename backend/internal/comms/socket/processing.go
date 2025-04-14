package socket

import (
	"fmt"
	"time"

	"github.com/ba2025-ysmprc/frr-tui/backend/internal/aggregator"
	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"
)

// func processCommand(message *frrProto.Message) *frrProto.Response {
func processCommand(message *frrProto.Message) *frrProto.Response {
	var response frrProto.Response

	switch message.Command {
	case "ospf":
		response.Status = "success"
		response.Message = "Returning magical ospf data"
		return aggregator.OSPFNeighborDummyData()
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
