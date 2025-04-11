package socket

import (
	"fmt"
	"time"

	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"
)

type Message struct {
	Command string                 `json:"command"`
	Package string                 `json:"package"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type Command struct {
	Package string                 `json:"package"`
	Action  string                 `json:"action"`
	Params  map[string]interface{} `json:"params,omitempty"`
}

type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

func processCommand(message *frrProto.Message) *frrProto.Response {
	var response frrProto.Response

	switch message.Command {
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
