package backend

import (
	"encoding/binary"
	"fmt"
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
		return nil, fmt.Errorf("unable to connect to %q: %w", path, err)
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

func GetLSDB() (*frrProto.OSPFDatabase, error) {
	response, err := SendMessage("ospf", "database", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfDatabase(), nil
}

func GetRouterAnomalies() (*frrProto.AnomalyDetection, error) {
	response, err := SendMessage("analysis", "dummyRouterOne", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetAnomaly(), nil
}

func GetExternalAnomalies() (*frrProto.AnomalyDetection, error) {
	response, err := SendMessage("analysis", "dummyExternalOne", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetAnomaly(), nil
}

func GetNSSAExternalAnomalies() (*frrProto.AnomalyDetection, error) {
	response, err := SendMessage("analysis", "dummyNSSAExternalOne", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetAnomaly(), nil
}
