package aggregator_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ba2025-ysmprc/frr-tui/backend/internal/aggregator"
)

// TestFullIntegration tests the complete flow from collection to state aggregation
func TestFullIntegration(t *testing.T) {
	// Create mock HTTP server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"neighbors": [
				{
					"neighborId": "192.168.1.2",
					"ipAddress": "192.168.1.2",
					"state": "FULL",
					"interface": "eth0",
					"area": "0.0.0.0"
				},
				{
					"neighborId": "192.168.2.2",
					"ipAddress": "192.168.2.2",
					"state": "INIT",
					"interface": "eth1",
					"area": "0.0.0.1"
				}
			],
			"routes": [
				{
					"prefix": "10.0.0.0/24",
					"nextHop": "192.168.1.2",
					"interface": "eth0",
					"cost": 10,
					"type": "intra-area",
					"area": "0.0.0.0"
				},
				{
					"prefix": "10.1.0.0/24",
					"nextHop": "192.168.2.2",
					"interface": "eth1",
					"cost": 20,
					"type": "external",
					"area": "0.0.0.1"
				}
			],
			"interfaces": [
				{
					"name": "eth0",
					"area": "0.0.0.0",
					"nbrCount": 1,
					"nbrAdjacentCount": 1,
					"passive": false
				},
				{
					"name": "eth1",
					"area": "0.0.0.1",
					"nbrCount": 1,
					"nbrAdjacentCount": 0,
					"passive": false
				},
				{
					"name": "lo0",
					"area": "0.0.0.0",
					"nbrCount": 0,
					"nbrAdjacentCount": 0,
					"passive": true
				}
			],
			"lsas": [
				{
					"type": "router",
					"lsId": "192.168.1.1",
					"advRouter": "192.168.1.1",
					"age": 3,
					"area": "0.0.0.0"
				},
				{
					"type": "network",
					"lsId": "192.168.1.0",
					"advRouter": "192.168.1.1",
					"age": 5,
					"area": "0.0.0.0"
				}
			]
		}`))
	}))
	defer server.Close()

	// Create temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "ospf.conf")
	configData := `
interface lo0
 ip ospf area 0.0.0.0
 ip ospf passive

interface eth0
 ip ospf area 0.0.0.0
 ip ospf cost 10

interface eth1
 ip ospf area 0.0.0.1
 ip ospf cost 20

router ospf
 ospf router-id 192.168.1.1
 network 192.168.1.0/24 area 0.0.0.0
 network 192.168.2.0/24 area 0.0.0.1
 network 127.0.0.0/8 area 0.0.0.0
exit
`
	err := os.WriteFile(configPath, []byte(configData), 0644)
	if err != nil {
		t.Fatalf("Failed to write test config file: %v", err)
	}

	// Create collector with mock server URL
	collector := aggregator.NewCollector(server.URL, configPath)

	// Test the Collect method
	state, err := collector.Collect()
	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	// Validate the integrated state
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	// Verify OSPF metrics
	if len(state.OSPF.Neighbors) != 2 {
		t.Errorf("Expected 2 neighbors, got %d", len(state.OSPF.Neighbors))
	}

	if len(state.OSPF.Routes) != 2 {
		t.Errorf("Expected 2 routes, got %d", len(state.OSPF.Routes))
	}

	if len(state.OSPF.Interfaces) != 3 {
		t.Errorf("Expected 3 interfaces, got %d", len(state.OSPF.Interfaces))
	}

	// Verify config parsing
	if state.Config.RouterID != "192.168.1.1" {
		t.Errorf("Expected RouterID to be '192.168.1.1', got '%s'", state.Config.RouterID)
	}

	if len(state.Config.Interfaces) != 3 {
		t.Errorf("Expected 3 interfaces in config, got %d", len(state.Config.Interfaces))
	}

	if len(state.Config.Areas) != 2 {
		t.Errorf("Expected 2 areas in config, got %d", len(state.Config.Areas))
	}

	// Verify system metrics are collected
	if state.System == nil {
		t.Fatal("Expected non-nil system metrics")
	}

	// Verify cache functionality
	cachedState := collector.GetCache()
	if cachedState != state {
		t.Error("Expected cache to be updated with the latest state")
	}

	// Verify timestamp is set
	if state.Timestamp.IsZero() {
		t.Error("Expected non-zero timestamp")
	}
}

// TestEdgeCaseIntegration tests edge cases that could occur in production
func TestEdgeCaseIntegration(t *testing.T) {
	// Test case 1: Server returns valid data but config is empty
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"neighbors":[],"routes":[],"interfaces":[],"lsas":[]}`))
	}))
	defer server.Close()

	tempDir := t.TempDir()
	emptyConfigPath := filepath.Join(tempDir, "empty.conf")

	err := os.WriteFile(emptyConfigPath, []byte(""), 0644)
	if err != nil {
		t.Fatalf("Failed to write empty config file: %v", err)
	}

	collector := aggregator.NewCollector(server.URL, emptyConfigPath)
	state, err := collector.Collect()

	if err != nil {
		t.Fatalf("Collect failed: %v", err)
	}

	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	if len(state.OSPF.Neighbors) != 0 {
		t.Errorf("Expected 0 neighbors, got %d", len(state.OSPF.Neighbors))
	}

	if len(state.Config.Interfaces) != 0 {
		t.Errorf("Expected 0 interfaces in config, got %d", len(state.Config.Interfaces))
	}

	// Test case 2: Config exists but server returns error
	badServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer badServer.Close()

	validConfigPath := filepath.Join(tempDir, "valid.conf")
	err = os.WriteFile(validConfigPath, []byte("interface eth0\n ip ospf area 0.0.0.0\n"), 0644)
	if err != nil {
		t.Fatalf("Failed to write valid config file: %v", err)
	}

	collector = aggregator.NewCollector(badServer.URL, validConfigPath)
	state, err = collector.Collect()

	if err == nil {
		t.Error("Expected error when server returns error")
	}

	if state != nil {
		t.Error("Expected nil state when error occurs")
	}
}
