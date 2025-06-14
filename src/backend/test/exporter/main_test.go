package exporter_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/frr-mad/frr-mad/src/backend/internal/configs"
	"github.com/frr-mad/frr-mad/src/backend/internal/exporter"
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewExporter(t *testing.T) {
	logPath := "/tmp/exporter_test.log"
	tests := []struct {
		name          string
		configPort    int
		expectedPort  int
		expectedError bool
	}{
		{
			name:         "Default port",
			configPort:   0,
			expectedPort: 9091,
		},
		{
			name:         "Custom valid port",
			configPort:   9092,
			expectedPort: 9092,
		},
		{
			name:          "Invalid port",
			configPort:    65536,
			expectedPort:  9091,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testLogger, err := logger.NewApplicationLogger("test", logPath)
			require.NoError(t, err)

			config := configs.ExporterConfig{
				Port: tt.configPort,
			}

			e := exporter.NewExporter(
				config,
				testLogger,
				time.Minute,
				&frrProto.FullFRRData{},
				&frrProto.AnomalyAnalysis{},
			)

			port, err := strconv.Atoi(e.Server.Addr[1:])
			assert.Nil(t, err)
			assert.Equal(t, tt.expectedPort, port)
			if tt.expectedError {
				assert.True(t, checkLogForWarning(t, logPath, "invalid port in config"),
					"Expected warning message not found in log file")
			}
		})
	}
}

func TestExporterStartStop(t *testing.T) {
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/exporter_test.log")
	require.NoError(t, err)

	config := configs.ExporterConfig{
		Port: 0,
	}

	e := exporter.NewExporter(
		config,
		testLogger,
		100*time.Millisecond,
		&frrProto.FullFRRData{},
		&frrProto.AnomalyAnalysis{},
	)

	// Test server startup
	serverStopped := make(chan struct{})
	go func() {
		e.Start()
		close(serverStopped)
	}()

	// Give server time to start
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get(fmt.Sprintf("http://localhost%s", e.Server.Addr))
	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, resp.StatusCode)
		resp.Body.Close()
	}

	// Test stop functionality
	e.Stop()

	select {
	case <-serverStopped:
	case <-time.After(1 * time.Second):
		assert.Fail(t, "Server did not stop within timeout")
	}

	_, err = http.Get(fmt.Sprintf("http://localhost%s", e.Server.Addr))
	assert.Error(t, err)
}

/*func TestExportLoop(t *testing.T) {
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/exporter_test.log")
	require.NoError(t, err)

	anomalies := &frrProto.AnomalyAnalysis{}
	frrData := &frrProto.FullFRRData{}
	config := configs.ExporterConfig{}
	registry := prometheus.NewRegistry()

	anomalyExporter := exporter.NewAnomalyExporter(anomalies, registry, testLogger)
	metricExporter := exporter.NewMetricExporter(frrData, registry, testLogger, config)

	// Create test exporter with short interval
	e := &exporter.Exporter{
		Interval:        50 * time.Millisecond,
		AnomalyExporter: anomalyExporter,
		MetricExporter:  metricExporter,
		Logger:          testLogger,
		StopChan:        make(chan struct{}),
	}

	// Count number of updates
	var anomalyUpdates, metricUpdates int
	// e.AnomalyExporter.Update = func() { anomalyUpdates++ }
	// e.MetricExporter.Update = func() { metricUpdates++ }

	// Start export loop
	go e.RunExportLoop()

	// Let it run for a few intervals
	time.Sleep(300 * time.Millisecond)

	// Stop the loop
	e.Stop()

	// Verify updates occurred
	assert.Greater(t, anomalyUpdates, 0, "Anomaly exporter should have been updated")
	assert.Greater(t, metricUpdates, 0, "Metric exporter should have been updated")
}*/

func TestExportData(t *testing.T) {
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/exporter_test.log")
	require.NoError(t, err)

	tests := []struct {
		name               string
		anomalyUpdateFails bool
		metricUpdateFails  bool
		expectErrorLogs    int
	}{
		{
			name:               "Successful updates",
			anomalyUpdateFails: false,
			metricUpdateFails:  false,
			expectErrorLogs:    0,
		},
		{
			name:               "Anomaly update fails",
			anomalyUpdateFails: true,
			metricUpdateFails:  false,
			expectErrorLogs:    1,
		},
		{
			name:               "Metric update fails",
			anomalyUpdateFails: false,
			metricUpdateFails:  true,
			expectErrorLogs:    1,
		},
		{
			name:               "Both updates fail",
			anomalyUpdateFails: true,
			metricUpdateFails:  true,
			expectErrorLogs:    2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create mock exporters
			anomalyExporter := &exporter.AnomalyExporter{}
			metricExporter := &exporter.MetricExporter{}

			// Set up mock behavior
			// if tt.anomalyUpdateFails {
			// 	anomalyExporter.Update = func() { panic("anomaly update failed") }
			// } else {
			// 	anomalyExporter.Update = func() {}
			// }

			// if tt.metricUpdateFails {
			// 	metricExporter.Update = func() { panic("metric update failed") }
			// } else {
			// 	metricExporter.Update = func() {}
			// }

			e := &exporter.Exporter{
				AnomalyExporter: anomalyExporter,
				MetricExporter:  metricExporter,
				Logger:          testLogger,
			}

			// Clear log file before test
			if err := os.WriteFile("/tmp/exporter_test.log", []byte{}, 0644); err != nil {
				t.Fatalf("Failed to clear log file: %v", err)
			}

			// Run export
			e.ExportData()

			// Check error logs if expected
			if tt.expectErrorLogs > 0 {
				content, err := os.ReadFile("/tmp/exporter_test.log")
				require.NoError(t, err)
				assert.Contains(t, string(content), "Failed to update after retry")
			}
		})
	}
}

/*func TestTryUpdateWithRetry(t *testing.T) {
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/exporter_test.log")
	require.NoError(t, err)

	tests := []struct {
		name        string
		updateFunc  func()
		expectError bool
		expectRetry bool
		expectPanic bool
	}{
		{
			name:        "Successful update",
			updateFunc:  func() {},
			expectError: false,
			expectRetry: false,
		},
		{
			name:        "Failing update with successful retry",
			updateFunc:  func() { panic("simulated failure") },
			expectError: false,
			expectRetry: true,
			expectPanic: true,
		},
		{
			name:        "Persistent failure",
			updateFunc:  func() { panic("persistent failure") },
			expectError: true,
			expectRetry: true,
			expectPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear log file before test
			if err := os.WriteFile("/tmp/exporter_test.log", []byte{}, 0644); err != nil {
				t.Fatalf("Failed to clear log file: %v", err)
			}

			var actualRetry bool
			updateWrapper := func() {
				if tt.expectPanic && actualRetry {
					return // Succeed on retry
				}
				tt.updateFunc()
			}

			err := exporter.TryUpdateWithRetry("TestComponent", updateWrapper, testLogger)

			if tt.expectError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			// Check logs
			content, err := os.ReadFile("/tmp/exporter_test.log")
			require.NoError(t, err)
			logContent := string(content)

			if tt.expectRetry {
				assert.Contains(t, logContent, "Update failed, will retry")
				if tt.expectError {
					assert.Contains(t, logContent, "Retry attempt failed")
				} else {
					assert.Contains(t, logContent, "Retry succeeded")
				}
			}
		})
	}
}*/

func TestExporterEndpoints(t *testing.T) {
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/exporter_test.log")
	require.NoError(t, err)

	config := configs.ExporterConfig{}
	e := exporter.NewExporter(
		config,
		testLogger,
		time.Minute,
		&frrProto.FullFRRData{},
		&frrProto.AnomalyAnalysis{},
	)

	tests := []struct {
		name         string
		path         string
		expectedCode int
		expectedBody string
	}{
		{
			name:         "Root endpoint",
			path:         "/",
			expectedCode: http.StatusOK,
			expectedBody: "<title>FRR MAD Exporter</title>",
		},
		{
			name:         "Metrics endpoint",
			path:         "/metrics",
			expectedCode: http.StatusOK,
			expectedBody: "",
		},
		{
			name:         "Nonexistent endpoint",
			path:         "/nonexistent",
			expectedCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest("GET", tt.path, nil)
			w := httptest.NewRecorder()

			e.Server.Handler.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedCode, w.Code)
			if tt.expectedBody != "" {
				assert.Contains(t, w.Body.String(), tt.expectedBody)
			}
		})
	}
}

// Use this test maybe?
// func TestExporterWithNilMetricExporter(t *testing.T) {
// 	testLogger, err := logger.NewApplicationLogger("test", "/tmp/exporter_test.log")
// 	require.NoError(t, err)

// 	e := &exporter.Exporter{
// 		AnomalyExporter: &exporter.AnomalyExporter{},
// 		MetricExporter:  nil,
// 		Logger:          testLogger,
// 	}

// 	if err := os.WriteFile("/tmp/exporter_test.log", []byte{}, 0644); err != nil {
// 		t.Fatalf("Failed to clear log file: %v", err)
// 	}

// 	e.ExportData()

// 	content, err := os.ReadFile("/tmp/exporter_test.log")
// 	require.NoError(t, err)
// 	assert.NotContains(t, string(content), "error", "No errors should be logged")
// }

// Helper functions

func checkLogForWarning(t *testing.T, logPath string, expectedMessage string) bool {
	t.Helper()

	content, err := os.ReadFile(logPath)
	if err != nil {
		t.Fatalf("Failed to read log file: %v", err)
		return false
	}

	logContent := string(content)
	return strings.Contains(logContent, expectedMessage)
}
