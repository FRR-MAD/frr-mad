package backend

import (
	"encoding/binary"
	"fmt"
	"google.golang.org/protobuf/encoding/protojson"
	"io"
	"net"
	"time"

	frrProto "github.com/ba2025-ysmprc/frr-tui/pkg"

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
) (*frrProto.Response, error) {
	// Build top‐level Message
	message := &frrProto.Message{
		Service: service,
		Command: command,
		Params:  params,
	}

	// Open the Unix socket
	conn, err := openSocket(socketPath)
	if err != nil {
		return nil, err
	}
	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			// todo: log to logger
			fmt.Printf("Failed to close connection: %s\n", err)
		}
	}(conn)

	if err := sendProto(conn, message); err != nil {
		return nil, err
	}

	res, err := receiveProto(conn)
	if err != nil {
		return nil, err
	}

	if res.Status != "success" {
		return nil, fmt.Errorf("backend error: %s", res.Message)
	}

	return res, nil
}

// openSocket dials the Unix‐domain socket at path and returns a live connection.
func openSocket(path string) (net.Conn, error) {
	conn, err := net.DialTimeout("unix", path, socketDialTimeout)
	if err != nil {
		return nil, fmt.Errorf("unable to connect to %q:\n\nBackend message:\n%w", path, err)
	}
	return conn, nil
}

// sendProto marshals the given protobuf message, prefixes it with a 4‑byte
// length header (little endian), and writes both to conn.
func sendProto(conn net.Conn, message *frrProto.Message) error {
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	var header [4]byte
	binary.LittleEndian.PutUint32(header[:], uint32(len(data)))
	if _, err := conn.Write(header[:]); err != nil {
		return fmt.Errorf("failed sending length header: %w", err)
	}

	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("failed sending payload: %w", err)
	}

	return nil
}

// receiveProto reads a 4‑byte length header, then that many bytes,
// and unmarshals them into a Response.
func receiveProto(conn net.Conn) (*frrProto.Response, error) {
	// Read length header
	var header [4]byte
	if _, err := io.ReadFull(conn, header[:]); err != nil {
		return nil, fmt.Errorf("failed reading length header: %w", err)
	}
	length := binary.LittleEndian.Uint32(header[:])
	if length > maxResponseSize {
		return nil, fmt.Errorf("response too big: %d bytes", length)
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(conn, buf); err != nil {
		return nil, fmt.Errorf("failed reading payload: %w", err)
	}

	res := &frrProto.Response{}
	if err := proto.Unmarshal(buf, res); err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}
	return res, nil
}

func GetRouterName() (string, string, error) {
	response, err := SendMessage("frr", "routerData", nil)
	if err != nil {
		return "", "", err
	}

	routerData := response.Data.GetFrrRouterData()

	routerName := routerData.RouterName
	ospfRouterId := routerData.OspfRouterId

	return routerName, ospfRouterId, nil
}

func GetSystemResources() (int64, float64, float64, error) {

	response, err := SendMessage("system", "allResources", nil)
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

func GetRIB() (*frrProto.RoutingInformationBase, error) {
	response, err := SendMessage("frr", "rib", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetRoutingInformationBase(), nil
}

func GetRibFibSummary() (*frrProto.RibFibSummaryRoutes, error) {
	response, err := SendMessage("frr", "ribfibSummary", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetRibFibSummaryRoutes(), nil
}

func GetLSDB() (*frrProto.OSPFDatabase, error) {
	response, err := SendMessage("ospf", "database", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfDatabase(), nil
}

func GetOspfRouterDataSelf() (*frrProto.OSPFRouterData, error) {
	response, err := SendMessage("ospf", "router", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfRouterData(), nil
}

func GetOSPF() (*frrProto.GeneralOspfInformation, error) {
	response, err := SendMessage("ospf", "generalInfo", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetGeneralOspfInformation(), nil
}

func GetOspfP2PInterfaceMapping() (*frrProto.PeerInterfaceMap, error) {
	response, err := SendMessage("ospf", "peerMap", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetPeerInterfaceToAddress(), nil
}

func GetOspfNetworkDataSelf() (*frrProto.OSPFNetworkData, error) {
	response, err := SendMessage("ospf", "network", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfNetworkData(), nil
}

func GetOspfNeighbors() (*frrProto.OSPFNeighbors, error) {
	response, err := SendMessage("ospf", "neighbors", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfNeighbors(), nil
}

func GetOspfNeighborInterfaces() ([]string, error) {
	response, err := SendMessage("ospf", "neighbors", nil)
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

func GetOspfExternalDataSelf() (*frrProto.OSPFExternalData, error) {
	response, err := SendMessage("ospf", "externalData", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfExternalData(), nil
}

func GetOspfNssaExternalDataSelf() (*frrProto.OSPFNssaExternalData, error) {
	response, err := SendMessage("ospf", "nssaExternalData", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfNssaExternalData(), nil
}

func GetStaticFRRConfigurationPretty() (string, error) {
	response, err := SendMessage("ospf", "staticConfig", nil)
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

func GetRouterAnomalies() (*frrProto.AnomalyDetection, error) {
	response, err := SendMessage("analysis", "router", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetAnomaly(), nil
}

func GetExternalAnomalies() (*frrProto.AnomalyDetection, error) {
	response, err := SendMessage("analysis", "external", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetAnomaly(), nil
}

func GetNSSAExternalAnomalies() (*frrProto.AnomalyDetection, error) {
	response, err := SendMessage("analysis", "nssaExternal", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetAnomaly(), nil
}
