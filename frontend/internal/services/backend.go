package backend

import (
	"fmt"
	"net"
	"time"

	pb "github.com/ba2025-ysmprc/frr-tui/pkg"

	"google.golang.org/protobuf/proto"
)

const socketPath = "/tmp/backend.sock"

// SendMessage sends a Message and waits for a Response from the backend.
func SendMessage(msg *pb.Message) (*pb.Response, error) {
	return sendProtobuf(msg)
}

// sendProtobuf handles marshaling, sending, and receiving protobuf messages
func sendProtobuf(msg *pb.Message) (*pb.Response, error) {
	conn, err := net.DialTimeout("unix", socketPath, 2*time.Second)
	if err != nil {
		return nil, fmt.Errorf("connect error: %w", err)
	}
	defer conn.Close()

	// --- Send ---
	data, err := proto.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("marshal error: %w", err)
	}

	length := uint32(len(data))
	lenBuf := []byte{
		byte(length >> 24),
		byte(length >> 16),
		byte(length >> 8),
		byte(length),
	}

	_, err = conn.Write(lenBuf)
	if err != nil {
		return nil, fmt.Errorf("send length error: %w", err)
	}
	_, err = conn.Write(data)
	if err != nil {
		return nil, fmt.Errorf("send data error: %w", err)
	}

	// --- Receive ---
	respLenBuf := make([]byte, 4)
	_, err = conn.Read(respLenBuf)
	if err != nil {
		return nil, fmt.Errorf("read length error: %w", err)
	}

	respLen := (uint32(respLenBuf[0]) << 24) |
		(uint32(respLenBuf[1]) << 16) |
		(uint32(respLenBuf[2]) << 8) |
		uint32(respLenBuf[3])

	respBuf := make([]byte, respLen)
	_, err = conn.Read(respBuf)
	if err != nil {
		return nil, fmt.Errorf("read data error: %w", err)
	}

	var res pb.Response
	err = proto.Unmarshal(respBuf, &res)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error: %w", err)
	}

	return &res, nil
}

func GetOSPFAnomalies() [][]string {
	// Fetch OSPF Anomalies via protobuf

	// parse received protobuf data

	// parsed protobuf message should look something like this:
	anomalyRows := [][]string{
		{"10.0.12.0/23", "unadvertised route", "OSPF Monitoring Tab 5", "Start"},
		{"10.0.15.0/14", "wrongly advertised", "OSPF Monitoring Tab 3", "Start"},
		{"10.0.199.0/23", "overadvertised route", "OSPF Monitoring Tab 2", "Start"},
		{"10.0.12.0/23", "unadvertised route", "OSPF Monitoring Tab 5", "Start"},
		{"10.0.15.0/14", "wrongly advertised", "OSPF Monitoring Tab 3", "Start"},
		{"10.0.199.0/23", "overadvertised route", "OSPF Monitoring Tab 2", "Start"},
		{"10.0.12.0/23", "unadvertised route", "OSPF Monitoring Tab 5", "Start"},
		{"10.0.15.0/14", "wrongly advertised", "OSPF Monitoring Tab 3", "Start"},
		{"10.0.199.0/23", "overadvertised route", "OSPF Monitoring Tab 2", "Start"},
		{"10.0.12.0/23", "unadvertised route", "OSPF Monitoring Tab 5", "Start"},
		{"10.0.15.0/14", "wrongly advertised", "OSPF Monitoring Tab 3", "Start"},
		{"100.100.100.100/23", "overadvertised route", "OSPF Monitoring Tab 2", "Start"},
	}

	return anomalyRows
}

func GetOSPFMetrics() [][]string {
	// Fetch all metrics (maybe fetch periodically everything and with the Getter function only provide requested data

	// this getter provides the OSPF metrics for the dashboard if no anomaly is detected

	// Stub or Transit Network does only exist for Router (Type 1) LSAs
	allGoodRows := [][]string{
		{"10.0.0.0/23", "Stub Network"},
		{"10.0.12.0/24", "Transit Network"},
		{"10.0.13.0/24", "Transit Network"},
		{"10.0.14.0/24", "Transit Network"},
		{"10.0.15.0/24", "Transit Network"},
		{"10.0.16.0/24", "Transit Network"},
		{"10.0.17.0/24", "Transit Network"},
		{"10.0.18.0/24", "Transit Network"},
		{"10.0.19.0/24", "Transit Network"},
	}

	return allGoodRows
}
