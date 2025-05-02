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
	case "frr":
		switch message.Command {
		case "routerData":
			return s.getRouterName()
		default:
			response.Status = "error"
			response.Message = "There was an error"
			return &response
		}
	case "ospf":
		return s.ospfProcessing(message.Command)
	case "analysis":
		return s.analysisProcessing(message.Command)
	case "system":
		switch message.Command {
		case "allResources":
			return s.getSystemResources()
		case "exit":
			response.Status = "success"
			response.Message = "Shutting system down"
			go func() {
				time.Sleep(100 * time.Millisecond)
				s.exitSocketServer()
			}()
			return &response
		default:
			response.Status = "error"
			response.Message = "There was an error getting system resources"
			return &response
		}
	default:
		response.Status = "error"
		response.Message = fmt.Sprintf("Unknown service: %s", message.Command)
		return &response
	}
}

func (s *Socket) ospfProcessing(command string) *frrProto.Response {
	var response frrProto.Response
	switch command {
	case "database":
		return s.getOspfDatabase()
	case "router":
		return s.getOspfRouterData()
	case "network":
		return s.getOspfNetworkData()
	case "summary":
		return s.getOspfSummaryData()
	case "asbrSummary":
		return s.getOspfAsbrSummaryData()
	case "externalData":
		return s.getOspfExternalData()
	case "nssaExternalData":
		return s.getOspfNssaExternalData()
	case "duplicates":
		return s.getOspfDuplicates()
	case "neighbors":
		return s.getOspfNeighbors()
	case "interfaces":
		return s.getInterfaces()
	case "rib":
		return s.getRoutingInformationBase()
	case "staticConfig":
		return s.getStaticFrrConfiguration()
	default:
		response.Status = "error"
		response.Message = fmt.Sprintf("Unknown command: %s", command)
		return &response
	}

}

func (s *Socket) analysisProcessing(command string) *frrProto.Response {
	var response frrProto.Response
	switch command {
	case "router":
		return s.getRouterAnomaly()
	case "external":
		return s.getExternalAnomaly()
	case "nssaExternal":
		return s.getNssaExternalAnomaly()

	case "dummyRouterOne":
		return getRouterAnomalyDummy1()
	default:
		response.Status = "error"
		response.Message = fmt.Sprintf("Unknown command: %s", command)
	}

	return &response
}
