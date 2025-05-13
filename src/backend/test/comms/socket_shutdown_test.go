package socket_test

import (
	"net"
	"os"
	"testing"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/comms/socket"
	"github.com/stretchr/testify/assert"
)

// TestSocketShutdown tests socket's ability to properly shut down and clean up resources
func TestSocketShutdown(t *testing.T) {
	config := configs.SocketConfig{
		UnixSocketLocation: "/tmp",
		UnixSocketName:     "test-shutdown-socket",
		SocketType:         "unix",
	}

	socketPath := "/tmp/test-shutdown-socket"

	os.Remove(socketPath)

	mockLoggerInstance, mockAnalyzerInstance, mockMetrics, p2pMap := getMockData()

	socketInstance := socket.NewSocket(config, mockMetrics, mockAnalyzerInstance.AnalysisResult, mockLoggerInstance, p2pMap)

	socketErrChan := make(chan error, 1)
	go func() {
		err := socketInstance.Start()
		socketErrChan <- err
	}()

	time.Sleep(100 * time.Millisecond)

	_, err := os.Stat(socketPath)
	assert.NoError(t, err, "Socket file should exist after Start()")

	conn, err := net.Dial("unix", socketPath)
	assert.NoError(t, err)
	assert.NotNil(t, conn)
	conn.Close()

	socketInstance.Close()

	_, err = os.Stat(socketPath)
	assert.True(t, os.IsNotExist(err), "Socket file should be removed during Close()")

	_, err = net.Dial("unix", socketPath)
	assert.Error(t, err, "Connection should fail after socket is closed")

	select {
	case err := <-socketErrChan:
		if err != nil {
			t.Logf("Socket server returned error: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
	}
}
