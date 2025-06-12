package socket_test

import (
	"net"
	"os"
	"testing"

	"github.com/frr-mad/frr-mad/src/backend/internal/configs"
	"github.com/frr-mad/frr-mad/src/backend/internal/socket"
	"github.com/stretchr/testify/assert"
)

// TestSocketShutdown tests socket's ability to properly shut down and clean up resources
func TestSocketShutdown(t *testing.T) {
	config := configs.SocketConfig{
		UnixSocketLocation: "/tmp",
		UnixSocketName:     "test-shutdown-socket",
		SocketType:         "unix",
	}

	socketPath := config.UnixSocketLocation + "/" + config.UnixSocketName

	os.Remove(socketPath)

	mockLoggerInstance, mockAnalyzerInstance, mockMetrics, parsedAnalyzerdata := getMockData()

	socketInstance := socket.NewSocket(config, mockMetrics, mockAnalyzerInstance.AnalysisResult, mockLoggerInstance, parsedAnalyzerdata)

	go socketInstance.Start()

	//_, err := os.Stat(socketPath)
	//assert.NoError(t, err, "Socket file should exist after Start()")

	// conn, err := net.Dial("unix", socketPath)
	// assert.NoError(t, err)
	// assert.NotNil(t, conn)
	// conn.Close()

	socketInstance.Close()

	_, err := os.Stat(socketPath)
	assert.True(t, os.IsNotExist(err), "Socket file should be removed during Close()")

	_, err = net.Dial("unix", socketPath)
	assert.Error(t, err, "Connection should fail after socket is closed")

}
