package socket_test

import (
	"encoding/binary"
	"net"
	"os"
	"testing"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/comms/socket"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/proto"
)

// TestMessageProcessing tests the socket's ability to process different message types
// and return appropriate responses
func TestMessageProcessing(t *testing.T) {
	config := configs.SocketConfig{
		UnixSocketLocation: "/tmp",
		UnixSocketName:     "test-message-socket",
		SocketType:         "unix",
	}

	os.Remove("/tmp/test-message-socket")

	mockLoggerInstance, mockAnalyzerInstance, mockMetrics := getMockData()

	socketInstance := socket.NewSocket(config, mockMetrics, mockAnalyzerInstance.AnalysisResult, mockLoggerInstance)

	go func() {
		socketInstance.Start()
	}()

	time.Sleep(100 * time.Millisecond)

	t.Run("TestOSPFDatabaseCommand", func(t *testing.T) {
		request := &frrProto.Message{
			Service: "ospf",
			Command: "database",
		}

		response := sendRequestAndGetResponse(t, request, "/tmp/test-message-socket")

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Returning OSPF database", response.Message)

		ospfDatabase := response.Data.GetOspfDatabase()
		assert.NotNil(t, ospfDatabase)
		assert.Equal(t, 1, len(ospfDatabase.Areas))
		assert.Equal(t, "192.168.1.1", ospfDatabase.Areas["0.0.0.0"].RouterLinkStates[0].Base.LsId)
	})

	t.Run("TestSystemResourcesCommand", func(t *testing.T) {
		request := &frrProto.Message{
			Service: "system",
			Command: "allResources",
		}

		response := sendRequestAndGetResponse(t, request, "/tmp/test-message-socket")

		assert.Equal(t, "success", response.Status)
		assert.Equal(t, "Returning magical system data", response.Message)

		systemMetrics := response.Data.GetSystemMetrics()
		assert.NotNil(t, systemMetrics)
		assert.Equal(t, 25.5, systemMetrics.CpuUsage)
		assert.Equal(t, 40.2, systemMetrics.MemoryUsage)
	})

	t.Run("TestInvalidCommand", func(t *testing.T) {
		request := &frrProto.Message{
			Service: "invalid",
			Command: "command",
		}

		response := sendRequestAndGetResponse(t, request, "/tmp/test-message-socket")

		assert.Equal(t, "error", response.Status)
		assert.Contains(t, response.Message, "Unknown command")
	})

	socketInstance.Close()
}

// Helper function to send a request to the socket and get the response
func sendRequestAndGetResponse(t *testing.T, request *frrProto.Message, socketPath string) *frrProto.Response {
	conn, err := net.Dial("unix", socketPath)
	assert.NoError(t, err)
	defer conn.Close()

	requestData, err := proto.Marshal(request)
	assert.NoError(t, err)

	sizeBuf := make([]byte, 4)
	binary.LittleEndian.PutUint32(sizeBuf, uint32(len(requestData)))

	_, err = conn.Write(sizeBuf)
	assert.NoError(t, err)
	_, err = conn.Write(requestData)
	assert.NoError(t, err)

	responseSizeBuf := make([]byte, 4)
	_, err = conn.Read(responseSizeBuf)
	assert.NoError(t, err)

	responseSize := binary.LittleEndian.Uint32(responseSizeBuf)
	responseData := make([]byte, responseSize)

	_, err = conn.Read(responseData)
	assert.NoError(t, err)

	response := &frrProto.Response{}
	err = proto.Unmarshal(responseData, response)
	assert.NoError(t, err)

	return response
}
