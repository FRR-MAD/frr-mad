package socket_test

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/frr-mad/frr-mad/src/backend/configs"
	"github.com/frr-mad/frr-mad/src/backend/internal/socket"
	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/stretchr/testify/assert"
)

// Test 1: Socket Creation and Initialization
func TestNewSocket(t *testing.T) {
	// Test configuration
	config := configs.SocketConfig{
		UnixSocketLocation: "/tmp",
		UnixSocketName:     "test-socket",
		SocketType:         "unix",
	}

	// Create mock dependencies
	mockLoggerInstance := &logger.Logger{}

	mockLoggerInstance, mockAnalyzerInstance, mockMetrics, parsedAnalyzerdata := getMockData()

	// Create socket
	socketInstance := socket.NewSocket(config, mockMetrics, mockAnalyzerInstance.AnalysisResult, mockLoggerInstance, parsedAnalyzerdata)

	assert.NotNil(t, socketInstance)
	os.Remove("/tmp/test-socket")
}

// Test 2: Socket Connection Handling
func TestSocketConnectionHandling(t *testing.T) {
	config := configs.SocketConfig{
		UnixSocketLocation: "/tmp",
		UnixSocketName:     "test-connection-socket",
		SocketType:         "unix",
	}

	os.Remove("/tmp/test-connection-socket")

	// Create mock dependencies
	mockLoggerInstance, mockAnalyzerInstance, mockMetrics, parsedAnalyzerdata := getMockData()

	// Create socket
	socketInstance := socket.NewSocket(config, mockMetrics, mockAnalyzerInstance.AnalysisResult, mockLoggerInstance, parsedAnalyzerdata)

	// Start socket server in a goroutine
	socketErrChan := make(chan error, 1)
	go func() {
		err := socketInstance.Start()
		socketErrChan <- err
	}()

	time.Sleep(100 * time.Millisecond)

	conn, err := net.Dial("unix", "/tmp/test-connection-socket")
	assert.NoError(t, err)
	assert.NotNil(t, conn)

	if conn != nil {
		conn.Close()
	}

	socketInstance.Close()

	select {
	case err := <-socketErrChan:
		if err != nil {
			assert.Fail(t, "Socket server returned error:", err.Error())
		}
	case <-time.After(500 * time.Millisecond):
	}

	_, err = os.Stat("/tmp/test-connection-socket")
	assert.True(t, os.IsNotExist(err), "Socket file should be removed during Close()")
}
