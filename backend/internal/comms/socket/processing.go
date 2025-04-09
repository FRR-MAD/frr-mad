package socket

import (
	"encoding/json"
	"time"
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

func processCommand(message string) string {
	var cmd Command
	err := json.Unmarshal([]byte(message), &cmd)
	if err != nil {
		resp := Response{
			Status:  "error",
			Message: "Invalid JSON format: " + err.Error(),
		}
		respJSON, _ := json.Marshal(resp)
		return string(respJSON)
	}

	var result interface{}
	var errMsg string

	switch cmd.Package {
	case "bgp":
		result, err = processBGPCommand(cmd.Action, cmd.Params)
	case "ospf":
		result, err = processOSPFCommand(cmd.Action, cmd.Params)
	case "exit":
		result = map[string]string{"result": "Socket server shutting down..."}
		go func() {
			time.Sleep(100 * time.Millisecond)
			exitSocketServer()
		}()

	default:
		errMsg = "Unknown package: " + cmd.Package
	}
	resp := Response{}

	if err != nil || errMsg != "" {
		resp.Status = "error"
		if errMsg != "" {
			resp.Message = errMsg
		} else {
			resp.Message = err.Error()
		}
	} else {
		resp.Status = "success"
		resp.Data = result
	}

	respJSON, _ := json.MarshalIndent(resp, "", "  ")
	return string(respJSON)
}
