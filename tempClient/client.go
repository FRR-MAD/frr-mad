package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"os"
	"time"

	frrProto "temp/pkg"

	"google.golang.org/protobuf/proto"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: client <socket_path> [command] [package]")
		os.Exit(1)
	}

	// socketPath := os.Args[1]
	service := "PING"
	packageName := "system"

	if len(os.Args) > 1 {
		service = os.Args[1]
	}

	if len(os.Args) > 2 {
		packageName = os.Args[2]
	}

	// Connect to the Unix socket
	conn, err := net.Dial("unix", "/tmp/analyzer.sock")
	if err != nil {
		fmt.Printf("Failed to connect to socket: %s\n", err)
		os.Exit(1)
	}

	defer func(conn net.Conn) {
		err := conn.Close()
		if err != nil {
			fmt.Printf("Failed to close connection: %s\n", err)
		}
	}(conn)

	// Create a Message
	message := &frrProto.Message{
		Service: service,
		Command: packageName,
		Params: map[string]*frrProto.ResponseValue{
			"client_id": {
				Kind: &frrProto.ResponseValue_StringValue{
					StringValue: "example_client",
				},
			},
		},
	}

	/*
		package: ospf
		command: update
		subpackage: lsa

		package: ospf
		command: read
		subpackage: neighbor

		package: ospf
		command: update
		subpackage: lsa

	*/

	// Send the message
	if err := sendMessage(conn, message); err != nil {
		fmt.Printf("Failed to send message: %s\n", err)
		os.Exit(1)
	}

	time.Sleep(100 * time.Millisecond)

	// Receive the response
	response, err := receiveResponse(conn)
	if err != nil {
		fmt.Printf("Failed to receive response: %s\n", err)
		os.Exit(1)
	}

	// Print data if present
	if response.Data != nil {
		printResponseData(response.Data)
	}
}

// Send a message to the server
func sendMessage(conn net.Conn, message *frrProto.Message) error {
	// Marshal the message
	data, err := proto.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Prepare size buffer (4 bytes, little-endian)
	sizeBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizeBuf, uint32(len(data)))

	// Send size followed by data
	if _, err := conn.Write(sizeBuf); err != nil {
		return fmt.Errorf("failed to send message size: %w", err)
	}

	if _, err := conn.Write(data); err != nil {
		return fmt.Errorf("failed to send message data: %w", err)
	}

	return nil
}

// Receive a response from the server
func receiveResponse(conn net.Conn) (*frrProto.Response, error) {
	// Read message size (4 bytes)
	sizeBuf := make([]byte, 4)
	if _, err := io.ReadFull(conn, sizeBuf); err != nil {
		return nil, fmt.Errorf("failed to read response size: %w", err)
	}

	// Convert bytes to uint32
	messageSize := binary.LittleEndian.Uint32(sizeBuf)

	// Sanity check
	if messageSize > 10*1024*1024 { // 10MB limit
		return nil, fmt.Errorf("response too large: %d bytes", messageSize)
	}

	// Read the response data
	messageBuf := make([]byte, messageSize)

	if _, err := io.ReadFull(conn, messageBuf); err != nil {
		return nil, fmt.Errorf("failed to read response data: %w", err)
	}

	// Unmarshal the response
	response := &frrProto.Response{}
	if err := proto.Unmarshal(messageBuf, response); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return response, nil
}

// Helper function to print response data based on its type
func printResponseData(data *frrProto.ResponseValue) {
	if data == nil {
		return
	}

	fmt.Printf("Data: %v\n", data.String())
}
