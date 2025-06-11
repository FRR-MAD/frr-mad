package socket

import (
	"fmt"
	"time"

	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
)

func (s *Socket) processCommand(message *frrProto.Message) *frrProto.Response {
	var response frrProto.Response

	switch message.Service {
	case "frr":
		return s.frrProcessing(message.Command)
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
		response.Message = fmt.Sprintf("Unknown service: %s", message.Service)
		return &response
	}
}

func (s *Socket) frrProcessing(command string) *frrProto.Response {
	var response frrProto.Response
	switch command {
	case "routerData":
		return s.getRouterName()
	case "rib":
		return s.getRoutingInformationBase()
	case "ribfibSummary":
		return s.getRibFibSummary()
	default:
		response.Status = "error"
		response.Message = "There was an error"
		return &response
	}
}

func (s *Socket) ospfProcessing(command string) *frrProto.Response {
	var response frrProto.Response
	switch command {
	case "database":
		return s.getOspfDatabase()
	case "generalInfo":
		return s.getGeneralOspfInformation()
	case "router":
		return s.getOspfRouterData()
	case "network":
		return s.getOspfNetworkData()
	case "networkAll":
		return s.getOspfNetworkDataAll()
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
	case "staticConfig": // TODO: should be added to case frr, because not only ospf data is contained
		return s.getStaticFrrConfiguration()
	case "peerMap":
		return s.getp2pMap()
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
	case "lsdbToRib":
		return s.getLsdbToRibAnomaly()
	case "ribToFib":
		return s.getRibToFibAnomaly()

	case "shouldParsedLsdb":
		return s.getShouldParsedLsdb()

	case "dummyRouterOne":
		return getRouterAnomalyDummy1()
	case "dummyExternalOne":
		return getExternalAnomalyDummy1()
	case "dummyNSSAExternalOne":
		return getNSSAExternalAnomalyDummy1()
	default:
		response.Status = "error"
		response.Message = fmt.Sprintf("Unknown command: %s", command)
	}

	return &response
}
