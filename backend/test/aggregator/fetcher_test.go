package aggregator_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/ba2025-ysmprc/frr-tui/backend/internal/aggregator"
)

func TestNewFetcher(t *testing.T) {
	fetcher := aggregator.NewFetcher("http://example.com")

	if fetcher.GetMetricURLForTesting() != "http://example.com" {
		t.Errorf("Expected metricsURL to be 'http://example.com', got '%s'", fetcher.GetMetricURLForTesting())
	}

	if fetcher.GetClientForTesting() == nil {
		t.Error("Expected HTTP client to be initialized")
	}

	if fetcher.GetClientForTesting().Timeout != 5*time.Second {
		t.Errorf("Expected timeout to be 5 seconds, got %v", fetcher.GetClientForTesting().Timeout)
	}
}

func TestFetchOSPF(t *testing.T) {
	// Valid response test
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
			],
			"hasRouteChanges": true
		}`))
	}))
	defer server.Close()

	fetcher := aggregator.NewFetcher(server.URL)
	metrics, err := fetcher.FetchOSPF()

	if err != nil {
		t.Fatalf("FetchOSPF failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("Expected non-nil metrics")
	}

	if !metrics.HasRouteChanges {
		t.Error("Expected HasRouteChanges to be true")
	}

	if len(metrics.Neighbors) != 1 {
		t.Errorf("Expected 1 neighbor, got %d", len(metrics.Neighbors))
	}

	if metrics.Neighbors[0].ID != "192.168.1.2" {
		t.Errorf("Expected neighbor ID '192.168.1.2', got '%s'", metrics.Neighbors[0].ID)
	}

	// Empty but valid response test
	emptyServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"neighbors":[],"routes":[],"interfaces":[],"lsas":[]}`))
	}))
	defer emptyServer.Close()

	fetcher = aggregator.NewFetcher(emptyServer.URL)
	metrics, err = fetcher.FetchOSPF()

	if err != nil {
		t.Fatalf("FetchOSPF failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("Expected non-nil metrics")
	}

	if len(metrics.Neighbors) != 0 {
		t.Errorf("Expected 0 neighbors, got %d", len(metrics.Neighbors))
	}
}

func TestFetchOSPFNoRouteChanges(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{
			"neighbors": [],
			"routes": [],
			"interfaces": [],
			"lsas": [],
			"hasRouteChanges": false
		}`))
	}))
	defer server.Close()

	fetcher := aggregator.NewFetcher(server.URL)
	metrics, err := fetcher.FetchOSPF()

	if err != nil {
		t.Fatalf("FetchOSPF failed: %v", err)
	}

	if metrics.HasRouteChanges {
		t.Error("Expected HasRouteChanges to be false")
	}
}

func TestFetchOSPFErrors(t *testing.T) {
	// Test case 1: Server error
	errorServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer errorServer.Close()

	fetcher := aggregator.NewFetcher(errorServer.URL)
	metrics, err := fetcher.FetchOSPF()

	if err == nil {
		t.Error("Expected error for server error response")
	}
	if metrics != nil {
		t.Error("Expected nil metrics on error")
	}

	// Test case 2: Invalid JSON
	badJSONServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"neighbors": [incomplete json`))
	}))
	defer badJSONServer.Close()

	fetcher = aggregator.NewFetcher(badJSONServer.URL)
	metrics, err = fetcher.FetchOSPF()

	if err == nil {
		t.Error("Expected error for invalid JSON")
	}
	if metrics != nil {
		t.Error("Expected nil metrics on error")
	}

	// Test case 3: Non-existent server
	fetcher = aggregator.NewFetcher("http://nonexistent.example.com")
	metrics, err = fetcher.FetchOSPF()

	if err == nil {
		t.Error("Expected error for non-existent server")
	}
	if metrics != nil {
		t.Error("Expected nil metrics on error")
	}
}

func TestCollectSystemMetrics(t *testing.T) {
	fetcher := aggregator.NewFetcher("http://example.com")
	metrics, err := fetcher.CollectSystemMetrics()

	if err != nil {
		t.Fatalf("CollectSystemMetrics failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("Expected non-nil system metrics")
	}

	// Note: Since the actual values depend on the system,
	// we just check that the function returns without error
}
