package backend

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"time"

	"google.golang.org/protobuf/encoding/protojson"

	"github.com/frr-mad/frr-mad/src/logger"
	frrProto "github.com/frr-mad/frr-tui/pkg"

	"google.golang.org/protobuf/proto"
)

// todo: take the path from the config file
const (
	// Path to the Unix domain socket your analyzer listens on.
	socketPath        = "/var/run/frr-mad/analyzer.sock"
	socketDialTimeout = 2 * time.Second

	// Maximum response size we’re willing to read (for sanity checking).
	maxResponseSize = 10 * 1024 * 1024 // 10 MB
)

// SendMessage sends a Message and waits for a Response from the backend.

func SendMessage(
	service string,
	command string,
	params map[string]*frrProto.ResponseValue,
	logger *logger.Logger,
) (*frrProto.Response, error) {
	// Build top‐level Message
	message := &frrProto.Message{
		Service: service,
		Command: command,
		Params:  params,
	}

	// Open the Unix socket
	conn, err := openSocket(socketPath, logger)
	if err != nil {
		return nil, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			// todo: log to logger
			logger.Error(fmt.Sprintf("Failed to close connection: %s\n", err))
			fmt.Printf("Failed to close connection: %s\n", err)
		}
	}(conn)

	if err := sendProto(conn, message, logger); err != nil {
		return nil, err
	}

	res, err := receiveProto(conn, logger)
	if err != nil {
		return nil, err
	}

	if res.Status != "success" {
		logger.Error(fmt.Sprintf("Error backend error sending and receiving proto status: %v message: %v", res.Status, res.Message))
		return nil, fmt.Errorf("backend error: %s", res.Message)
	}

	return res, nil
}

// openSocket dials the Unix‐domain socket at path and returns a live connection.
func openSocket(path string, logger *logger.Logger) (net.Conn, error) {
	conn, err := net.DialTimeout("unix", path, socketDialTimeout)
	if err != nil {
		logger.Error(fmt.Sprintf("unable to connect to %q:\n\nBackend message:\n%w", path, err))
		return nil, fmt.Errorf("unable to connect to %q:\n\nBackend message:\n%w", path, err)
	}
	return conn, nil
}

// sendProto marshals the given protobuf message, prefixes it with a 4‑byte
// length header (little endian), and writes both to conn.
func sendProto(conn net.Conn, message *frrProto.Message, logger *logger.Logger) error {
	data, err := proto.Marshal(message)
	if err != nil {
		logger.Error(fmt.Sprintf("Error on mashal message (Proto): %v", err))
		return fmt.Errorf("marshal error: %w", err)
	}

	var header [4]byte
	binary.LittleEndian.PutUint32(header[:], uint32(len(data)))
	if _, err := conn.Write(header[:]); err != nil {
		logger.Error(fmt.Sprintf("Error sending length header (Proto): %v", err))
		return fmt.Errorf("failed sending length header: %w", err)
	}

	if _, err := conn.Write(data); err != nil {
		logger.Error(fmt.Sprintf("Error sending payload (Proto): %v", err))
		return fmt.Errorf("failed sending payload: %w", err)
	}

	return nil
}

// receiveProto reads a 4‑byte length header, then that many bytes,
// and unmarshals them into a Response.
func receiveProto(conn net.Conn, logger *logger.Logger) (*frrProto.Response, error) {
	// Read length header
	var header [4]byte
	if _, err := io.ReadFull(conn, header[:]); err != nil {
		logger.Error(fmt.Sprintf("Error reading length header (Proto): %v", err))
		return nil, fmt.Errorf("failed reading length header: %w", err)
	}
	length := binary.LittleEndian.Uint32(header[:])
	if length > maxResponseSize {
		logger.Error(fmt.Sprintf("Error response too big (Proto): %v bytes", length))
		return nil, fmt.Errorf("response too big: %d bytes", length)
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(conn, buf); err != nil {
		logger.Error(fmt.Sprintf("Error failed reading payload (Proto): %v", err))
		return nil, fmt.Errorf("failed reading payload: %w", err)
	}

	res := &frrProto.Response{}
	if err := proto.Unmarshal(buf, res); err != nil {
		logger.Error(fmt.Sprintf("Error on unmarshal message (Proto): %v", err))
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	return res, nil
}

func GetRouterName(logger *logger.Logger) (string, string, error) {
	response, err := SendMessage("frr", "routerData", nil, logger)
	if err != nil {
		return "", "", err
	}

	routerData := response.Data.GetFrrRouterData()

	routerName := routerData.RouterName
	ospfRouterId := routerData.OspfRouterId

	return routerName, ospfRouterId, nil
}

func GetSystemResources(logger *logger.Logger) (int64, float64, float64, error) {

	response, err := SendMessage("system", "allResources", nil, logger)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("rpc error: %w", err)
	}
	if response.Status != "success" {
		return 0, 0, 0, fmt.Errorf("backend returned status %q: %s", response.Status, response.Message)
	}

	systemMetrics := response.Data.GetSystemMetrics()

	cores := systemMetrics.CpuAmount
	cpuUsage := systemMetrics.CpuUsage
	memoryUsage := systemMetrics.MemoryUsage

	return cores, cpuUsage, memoryUsage, nil
}

func GetRIB(logger *logger.Logger) (*frrProto.RoutingInformationBase, error) {
	response, err := SendMessage("frr", "rib", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetRoutingInformationBase(), nil
}

func GetRibFibSummary(logger *logger.Logger) (*frrProto.RibFibSummaryRoutes, error) {
	response, err := SendMessage("frr", "ribfibSummary", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetRibFibSummaryRoutes(), nil
}

func GetLSDB(logger *logger.Logger) (*frrProto.OSPFDatabase, error) {
	response, err := SendMessage("ospf", "database", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfDatabase(), nil
}

func GetOSPF(logger *logger.Logger) (*frrProto.GeneralOspfInformation, error) {
	response, err := SendMessage("ospf", "generalInfo", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetGeneralOspfInformation(), nil
}

func GetOspfRouterDataSelf(logger *logger.Logger) (*frrProto.OSPFRouterData, error) {
	response, err := SendMessage("ospf", "router", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfRouterData(), nil
}

func GetOspfP2PInterfaceMapping(logger *logger.Logger) (*frrProto.PeerInterfaceMap, error) {
	response, err := SendMessage("ospf", "peerMap", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetPeerInterfaceToAddress(), nil
}

func GetOspfNetworkDataSelf(logger *logger.Logger) (*frrProto.OSPFNetworkData, error) {
	response, err := SendMessage("ospf", "network", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfNetworkData(), nil
}

func GetOspfNeighbors(logger *logger.Logger) (*frrProto.OSPFNeighbors, error) {
	response, err := SendMessage("ospf", "neighbors", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfNeighbors(), nil
}

func GetOspfNeighborInterfaces(logger *logger.Logger) ([]string, error) {
	response, err := SendMessage("ospf", "neighbors", nil, logger)
	if err != nil {
		return nil, err
	}
	ospfNeighbors := response.Data.GetOspfNeighbors()

	var neighborAddresses []string
	for _, neighborGroup := range ospfNeighbors.Neighbors {
		for _, neighbor := range neighborGroup.Neighbors {
			neighborAddresses = append(neighborAddresses, neighbor.IfaceAddress)
		}
	}

	return neighborAddresses, nil
}

func GetOspfSummaryDataSelf(logger *logger.Logger) (*frrProto.OSPFSummaryData, error) {
	response, err := SendMessage("ospf", "summary", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfSummaryData(), nil
}

func GetOspfAsbrSummaryDataSelf(logger *logger.Logger) (*frrProto.OSPFAsbrSummaryData, error) {
	response, err := SendMessage("ospf", "asbrSummary", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfAsbrSummaryData(), nil
}

func GetOspfExternalDataSelf(logger *logger.Logger) (*frrProto.OSPFExternalData, error) {
	response, err := SendMessage("ospf", "externalData", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfExternalData(), nil
}

func GetOspfNssaExternalDataSelf(logger *logger.Logger) (*frrProto.OSPFNssaExternalData, error) {
	response, err := SendMessage("ospf", "nssaExternalData", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfNssaExternalData(), nil
}

func GetStaticFRRConfiguration(logger *logger.Logger) (*frrProto.StaticFRRConfiguration, error) {
	response, err := SendMessage("ospf", "staticConfig", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetStaticFrrConfiguration(), nil
}

func GetStaticFRRConfigurationPretty(logger *logger.Logger) (string, error) {
	response, err := SendMessage("ospf", "staticConfig", nil, logger)
	if err != nil {
		return "", err
	}

	var prettyJson string

	// Pretty‑print the protobuf into nice indented JSON
	marshaler := protojson.MarshalOptions{
		Multiline:     true,
		Indent:        "  ",
		UseProtoNames: true,
	}
	pretty, perr := marshaler.Marshal(response.Data)
	if perr != nil {
		prettyJson = response.Data.String()
	} else {
		prettyJson = string(pretty)
	}

	return prettyJson, nil
}

func GetRouterAnomalies(logger *logger.Logger) (*frrProto.AnomalyDetection, error) {
	response, err := SendMessage("analysis", "router", nil, logger) // router / dummyRouterOne
	if err != nil {
		return nil, err
	}

	return response.Data.GetAnomaly(), nil
}

func GetExternalAnomalies(logger *logger.Logger) (*frrProto.AnomalyDetection, error) {
	response, err := SendMessage("analysis", "external", nil, logger) // external / dummyExternalOne
	if err != nil {
		return nil, err
	}

	return response.Data.GetAnomaly(), nil
}

func GetNSSAExternalAnomalies(logger *logger.Logger) (*frrProto.AnomalyDetection, error) {
	response, err := SendMessage("analysis", "nssaExternal", nil, logger) // nssaExternal / dummyNSSAExternalOne
	if err != nil {
		return nil, err
	}

	return response.Data.GetAnomaly(), nil
}

func GetLSDBToRibAnomalies(logger *logger.Logger) (*frrProto.AnomalyDetection, error) {
	response, err := SendMessage("analysis", "lsdbToRib", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetAnomaly(), nil
}

func GetRibToFibAnomalies(logger *logger.Logger) (*frrProto.AnomalyDetection, error) {
	response, err := SendMessage("analysis", "ribToFib", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetAnomaly(), nil
}

func GetParsedShouldStates(logger *logger.Logger) (*frrProto.ParsedAnalyzerData, error) {
	response, err := SendMessage("analysis", "shouldParsedLsdb", nil, logger)
	if err != nil {
		return nil, err
	}

	return response.Data.GetParsedAnalyzerData(), nil
}
