package socket

import (
	"fmt"
	"time"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

// func processCommand(message *frrProto.Message) *frrProto.Response {
func (s *Socket) processCommand(message *frrProto.Message) *frrProto.Response {
	var response frrProto.Response

	switch message.Service {
	case "ospf":
		switch message.Command {
		case "dummy":
			//return s.ospfDummyData()

			response.Status = "success"
			response.Message = "Returning magical ospf dummy data"

			return &response
		default:
			response.Status = "error"
			response.Message = "There was an error"
			return &response
		}
	case "exit":
		response.Status = "success"
		response.Message = "Shutting system down"
		go func() {
			time.Sleep(100 * time.Millisecond)
			s.exitSocketServer()
		}()
		return &response
	case "system":
		switch message.Command {
		case "allResources":
			return s.getSystemResources()

		default:
			response.Status = "error"
			response.Message = "There was an error getting system resources"
			return &response
		}
	default:
		response.Status = "error"
		response.Message = fmt.Sprintf("Unknown command: %s", message.Command)
		return &response
	}
}
