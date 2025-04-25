package socket_test

import (
	"encoding/binary"
	"net"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/comms/socket"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

// TestConcurrentConnections tests the socket's ability to handle multiple
// simultaneous connections from different clients
func TestConcurrentConnections(t *testing.T) {
	config := configs.SocketConfig{
		UnixSocketLocation: "/tmp",
		UnixSocketName:     "test-concurrent-socket",
		SocketType:         "unix",
	}

	socketPath := "/tmp/test-concurrent-socket"

	os.Remove(socketPath)

	mockLoggerInstance, mockAnalyzerInstance, mockFullFRRData := getMockData()

	socketInstance := socket.NewSocket(config, mockFullFRRData, mockAnalyzerInstance, mockLoggerInstance)

	go func() {
		err := socketInstance.Start()
		if err != nil {
			t.Logf("Socket server returned error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	numClients := 10

	var wg sync.WaitGroup
	wg.Add(numClients)

	type clientResult struct {
		clientID         int
		success          bool
		responseReceived bool
		responseValid    bool
		err              error
	}

	results := make(chan clientResult, numClients)

	for i := 0; i < numClients; i++ {
		go func(clientID int) {
			defer wg.Done()

			result := clientResult{clientID: clientID}

			request := &frrProto.Message{
				Service: "system",
				Command: "allResources",
			}

			response, err := sendRequest(socketPath, request)
			if err != nil {
				result.err = err
				results <- result
				return
			}

			result.success = true
			result.responseReceived = true

			if response.Status == "success" &&
				response.Data != nil &&
				response.Data.GetSystemMetrics() != nil {
				result.responseValid = true
			}

			results <- result
		}(i)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out waiting for clients to complete")
	}

	close(results)

	successCount := 0
	responseCount := 0
	validCount := 0

	for result := range results {
		if result.err != nil {
			t.Logf("Client %d error: %v", result.clientID, result.err)
		}

		if result.success {
			successCount++
		}

		if result.responseReceived {
			responseCount++
		}

		if result.responseValid {
			validCount++
		}
	}

	assert.Equal(t, numClients, successCount, "All clients should connect successfully")

	assert.Equal(t, numClients, responseCount, "All clients should receive a response")

	assert.Equal(t, numClients, validCount, "All responses should be valid")

	socketInstance.Close()
}

// TestConcurrentCommandProcessing tests if multiple different commands can be processed
// concurrently without race conditions or data corruption
func TestConcurrentCommandProcessing(t *testing.T) {
	config := configs.SocketConfig{
		UnixSocketLocation: "/tmp",
		UnixSocketName:     "test-concurrent-commands-socket",
		SocketType:         "unix",
	}

	socketPath := "/tmp/test-concurrent-commands-socket"

	os.Remove(socketPath)

	mockLoggerInstance, mockAnalyzerInstance, mockFullFRRData := getMockData()

	socketInstance := socket.NewSocket(config, mockFullFRRData, mockAnalyzerInstance, mockLoggerInstance)

	go func() {
		err := socketInstance.Start()
		if err != nil {
			t.Logf("Socket server returned error: %v", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	commands := []struct {
		service string
		command string
	}{
		{service: "ospf", command: "database"},
		{service: "ospf", command: "router"},
		{service: "ospf", command: "network"},
		{service: "ospf", command: "neighbors"},
		{service: "system", command: "allResources"},
	}

	numClients := len(commands)

	var wg sync.WaitGroup
	wg.Add(numClients)

	type commandResult struct {
		service string
		command string
		success bool
		err     error
	}

	results := make(chan commandResult, numClients)

	for i, cmd := range commands {
		go func(idx int, svc, cmd string) {
			defer wg.Done()

			result := commandResult{
				service: svc,
				command: cmd,
			}

			request := &frrProto.Message{
				Service: svc,
				Command: cmd,
			}

			response, err := sendRequest(socketPath, request)
			if err != nil {
				result.err = err
				results <- result
				return
			}

			if response.Status == "success" {
				result.success = true
			}

			results <- result
		}(i, cmd.service, cmd.command)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out waiting for clients to complete")
	}

	close(results)

	successCount := 0

	for result := range results {
		if result.err != nil {
			t.Logf("Command %s/%s error: %v", result.service, result.command, result.err)
		}

		if result.success {
			successCount++
		}
	}

	assert.Equal(t, numClients, successCount, "All commands should be processed successfully")

	socketInstance.Close()
}

// Helper function to send a request and get a response
func sendRequest(socketPath string, request *frrProto.Message) (*frrProto.Response, error) {
	conn, err := net.Dial("unix", socketPath)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	requestData, err := proto.Marshal(request)
	if err != nil {
		return nil, err
	}

	sizeBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizeBuf, uint32(len(requestData)))

	_, err = conn.Write(sizeBuf)
	if err != nil {
		return nil, err
	}
	_, err = conn.Write(requestData)
	if err != nil {
		return nil, err
	}

	responseSizeBuf := make([]byte, 4)
	_, err = conn.Read(responseSizeBuf)
	if err != nil {
		return nil, err
	}

	responseSize := binary.LittleEndian.Uint32(responseSizeBuf)
	responseData := make([]byte, responseSize)

	_, err = conn.Read(responseData)
	if err != nil {
		return nil, err
	}

	response := &frrProto.Response{}
	err = proto.Unmarshal(responseData, response)
	if err != nil {
		return nil, err
	}

	return response, nil
}
