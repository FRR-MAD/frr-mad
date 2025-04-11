package aggregator_test

import (
	"testing"

	"github.com/ba2025-ysmprc/frr-tui/backend/internal/aggregator"
)

// MockFetcher is a mock implementation of the fetcher for testing
type MockFetcher struct {
	ospfMetrics    *aggregator.OSPFMetrics
	systemMetrics  *aggregator.SystemMetrics
	shouldFailOSPF bool
	shouldFailSys  bool
}

func NewMockFetcher() *MockFetcher {
	return &MockFetcher{
		ospfMetrics: &aggregator.OSPFMetrics{
			Neighbors: []aggregator.OSPFNeighbor{
				{
					ID:        "192.168.1.2",
					IP:        "192.168.1.2",
					State:     "FULL",
					Interface: "eth0",
					Area:      "0.0.0.0",
				},
			},
			Routes: []aggregator.OSPFRoute{
				{
					Prefix:    "10.0.0.0/24",
					NextHop:   "192.168.1.2",
					Interface: "eth0",
					Cost:      10,
					Type:      "intra-area",
					Area:      "0.0.0.0",
				},
			},
			Interfaces: []aggregator.OSPFInterface{
				{
					Name:     "eth0",
					Area:     "0.0.0.0",
					NbrCount: 1,
					NbrAdj:   1,
					Passive:  false,
				},
			},
		},
		systemMetrics: &aggregator.SystemMetrics{
			CPUUsage:    10.5,
			MemoryUsage: 30.2,
		},
	}
}

func (m *MockFetcher) FetchOSPF() (*aggregator.OSPFMetrics, error) {
	if m.shouldFailOSPF {
		return nil, &MockError{"OSPF fetch failed"}
	}
	return m.ospfMetrics, nil
}

func (m *MockFetcher) CollectSystemMetrics() (*aggregator.SystemMetrics, error) {
	if m.shouldFailSys {
		return nil, &MockError{"System metrics collection failed"}
	}
	return m.systemMetrics, nil
}

func (m *MockFetcher) SetShouldFailOSPF(fail bool) {
	m.shouldFailOSPF = fail
}

func (m *MockFetcher) SetShouldFailSys(fail bool) {
	m.shouldFailSys = fail
}

// MockError is a custom error for testing
type MockError struct {
	msg string
}

func (e *MockError) Error() string {
	return e.msg
}

// Test the mock fetcher to ensure it works as expected
func TestMockFetcher(t *testing.T) {
	mock := NewMockFetcher()

	// Test successful OSPF fetch
	metrics, err := mock.FetchOSPF()
	if err != nil {
		t.Fatalf("MockFetcher.FetchOSPF failed: %v", err)
	}

	if metrics == nil {
		t.Fatal("Expected non-nil metrics")
	}

	if len(metrics.Neighbors) != 1 {
		t.Errorf("Expected 1 neighbor, got %d", len(metrics.Neighbors))
	}

	// Test failed OSPF fetch
	mock.SetShouldFailOSPF(true)
	metrics, err = mock.FetchOSPF()

	if err == nil {
		t.Error("Expected error for OSPF fetch")
	}

	if metrics != nil {
		t.Error("Expected nil metrics when fetch fails")
	}

	// Test successful system metrics collection
	mock.SetShouldFailOSPF(false)
	sysMetrics, err := mock.CollectSystemMetrics()

	if err != nil {
		t.Fatalf("MockFetcher.CollectSystemMetrics failed: %v", err)
	}

	if sysMetrics == nil {
		t.Fatal("Expected non-nil system metrics")
	}

	if sysMetrics.CPUUsage != 10.5 {
		t.Errorf("Expected CPU usage 10.5, got %f", sysMetrics.CPUUsage)
	}

	// Test failed system metrics collection
	mock.SetShouldFailSys(true)
	sysMetrics, err = mock.CollectSystemMetrics()

	if err == nil {
		t.Error("Expected error for system metrics collection")
	}

	if sysMetrics != nil {
		t.Error("Expected nil system metrics when collection fails")
	}
}
