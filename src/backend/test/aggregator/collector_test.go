package aggregator_test

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator"
)

func TestNewCollector(t *testing.T) {
	collector := aggregator.NewCollector("http://example.com", "/path/to/config")

	if collector.GetFetcherForTesting() == nil {
		t.Error("Expected fetcher to be initialized")
	}

	if collector.GetConfigPathForTesting() != "/path/to/config" {
		t.Errorf("Expected configPath to be '/path/to/config', got '%s'", collector.GetConfigPathForTesting())
	}

	if collector.GetCacheForTesting() != nil {
		t.Error("Expected cache to be nil initially")
	}
}

func TestCollect(t *testing.T) {
	// Create a mock HTTP server
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
				}
			],
			"interfaces": [
				{
					"name": "eth0",
					"area": "0.0.0.0",
					"nbrCount": 1,
					"nbrAdjacentCount": 1,
					"passive": false
				}
			],
			"lsas": [
				{
					"type": "router",
					"lsId": "192.168.1.1",
					"advRouter": "192.168.1.1",
					"age": 3,
					"area": "0.0.0.0"
				}
			]
		}`))
	}))
	defer server.Close()

	// Create a temporary config file
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "ospf.conf")
	configData := `
					interface eth0
						ip address 10.0.16.1/24
						ip ospf area 0.0.0.0
						ip ospf cost 10

					router ospf
						ospf router-id 192.168.1.1
						network 192.168.1.0/24 area 0.0.0.0
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

	// Validate the collected data
	if state == nil {
		t.Fatal("Expected non-nil state")
	}

	if state.OSPF == nil {
		t.Fatal("Expected non-nil OSPF metrics")
	}

	if len(state.OSPF.Neighbors) != 1 {
		t.Errorf("Expected 1 neighbor, got %d", len(state.OSPF.Neighbors))
	}

	if len(state.OSPF.Routes) != 1 {
		t.Errorf("Expected 1 route, got %d", len(state.OSPF.Routes))
	}

	if state.Config == nil {
		t.Fatal("Expected non-nil config")
	}

	if len(state.Config.Interfaces) != 1 {
		t.Errorf("Expected 1 config, got %d", len(state.Config.Interfaces))
	}

	if state.Config.RouterID != "192.168.1.1" {
		t.Errorf("Expected RouterID to be '192.168.1.1', got '%s'", state.Config.RouterID)
	}

	// Test the cache functionality
	cachedState := collector.GetCache()
	if cachedState != state {
		t.Error("Expected cache to be updated with the latest state")
	}
}

func TestCollectErrors(t *testing.T) {
	// Test case 1: OSPF fetch fails
	badServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer badServer.Close()

	collector := aggregator.NewCollector(badServer.URL, "nonexistent.conf")
	state, err := collector.Collect()

	if err == nil {
		t.Error("Expected error for bad HTTP response")
	}
	if state != nil {
		t.Error("Expected nil state on error")
	}

	// Test case 2: Config parse fails
	goodServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"neighbors":[],"routes":[],"interfaces":[],"lsas":[]}`))
	}))
	defer goodServer.Close()

	collector = aggregator.NewCollector(goodServer.URL, "nonexistent.conf")
	state, err = collector.Collect()

	if err == nil {
		t.Error("Expected error for nonexistent config file")
	}
	if state != nil {
		t.Error("Expected nil state on error")
	}
}
