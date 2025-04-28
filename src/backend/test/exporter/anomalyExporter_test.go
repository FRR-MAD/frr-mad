package exporter_test

import (
	"sync"
	"testing"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/exporter"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAnomalyExporter_NoAnomalies(t *testing.T) {
	// Setup
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewLogger("test", "/tmp/exporter.log")
	assert.NoError(t, err)
	anomalies := &frrProto.Anomalies{}

	// Create exporter
	exporter := exporter.NewAnomalyExporter(anomalies, registry, testLogger)

	// Test
	exporter.Update()

	// Verify all gauges are 0
	metrics, err := registry.Gather()
	assert.NoError(t, err)

	for _, metric := range metrics {
		switch *metric.Name {
		case "ospf_overadvertised_route_present":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_underadvertised_route_present":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_duplicate_route_present":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_misconfigured_route_present":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_overadvertised_routes_total":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_underadvertised_routes_total":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_duplicate_routes_total":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_misconfigured_routes_total":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		}
	}
}

func TestAnomalyExporter_WithAnomalies(t *testing.T) {
	// Setup
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewLogger("test", "/tmp/exporter.log")
	assert.NoError(t, err)
	now := timestamppb.Now()

	anomalies := &frrProto.Anomalies{
		OveradvertisedRoutes: []*frrProto.AnomalyOveradvertisedRoute{
			{
				Timestamp:        now,
				Service:          "ospf",
				IsAdvertised:     "10.0.0.0/24",
				ShouldAdvertised: "10.0.0.0/16",
				Router: &frrProto.RouterAttribute{
					RouterName: "router1",
					RouterId:   "1.1.1.1",
				},
			},
			{
				Timestamp:        now,
				Service:          "ospf",
				IsAdvertised:     "192.168.1.0/24",
				ShouldAdvertised: "192.168.1.0/16",
				Router: &frrProto.RouterAttribute{
					RouterName: "router2",
					RouterId:   "2.2.2.2",
				},
			},
		},
		UnderadvertisedRoutes: []*frrProto.AnomalyUnderadvertisedRoute{
			{
				Timestamp:        now,
				Service:          "ospf",
				IsAdvertised:     "10.1.0.0/24",
				ShouldAdvertised: "10.1.0.0/16",
				Router: &frrProto.RouterAttribute{
					RouterName: "router3",
					RouterId:   "3.3.3.3",
				},
			},
		},
		DuplicateRoutes: []*frrProto.AnomalyDuplicateRoute{
			{
				Timestamp:        now,
				Service:          "ospf",
				IsAdvertised:     "172.16.0.0/24",
				ShouldAdvertised: "172.16.0.0/24",
				Router: &frrProto.RouterAttribute{
					RouterName: "router4",
					RouterId:   "4.4.4.4",
				},
			},
		},
		MisconfiguredRoutes: []*frrProto.AnomalyMisconfiguredRoute{
			{
				Timestamp:        now,
				Service:          "ospf",
				IsAdvertised:     "10.2.0.0/24",
				ShouldAdvertised: "10.2.0.0/24[tag:100]",
				Router: &frrProto.RouterAttribute{
					RouterName: "router5",
					RouterId:   "5.5.5.5",
				},
			},
			{
				Timestamp:        now,
				Service:          "ospf",
				IsAdvertised:     "10.3.0.0/24",
				ShouldAdvertised: "10.3.0.0/24[tag:200]",
				Router: &frrProto.RouterAttribute{
					RouterName: "router6",
					RouterId:   "6.6.6.6",
				},
			},
		},
	}

	// Create exporter
	exporter := exporter.NewAnomalyExporter(anomalies, registry, testLogger)

	// Test
	exporter.Update()

	// Verify metrics
	metrics, err := registry.Gather()
	assert.NoError(t, err)

	for _, metric := range metrics {
		switch *metric.Name {
		case "ospf_overadvertised_route_present":
			assert.Equal(t, 1.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_underadvertised_route_present":
			assert.Equal(t, 1.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_duplicate_route_present":
			assert.Equal(t, 1.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_misconfigured_route_present":
			assert.Equal(t, 1.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_overadvertised_routes_total":
			assert.Equal(t, 2.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_underadvertised_routes_total":
			assert.Equal(t, 1.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_duplicate_routes_total":
			assert.Equal(t, 1.0, metric.Metric[0].Gauge.GetValue())
		case "ospf_misconfigured_routes_total":
			assert.Equal(t, 2.0, metric.Metric[0].Gauge.GetValue())
		}
	}
}

func TestAnomalyExporter_ConcurrentUpdates(t *testing.T) {
	// Setup
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewLogger("test", "/tmp/exporter.log")
	assert.NoError(t, err)
	now := timestamppb.Now()

	anomalies := &frrProto.Anomalies{
		OveradvertisedRoutes: []*frrProto.AnomalyOveradvertisedRoute{
			{
				Timestamp:        now,
				Service:          "ospf",
				IsAdvertised:     "10.0.0.0/24",
				ShouldAdvertised: "10.0.0.0/16",
				Router: &frrProto.RouterAttribute{
					RouterName: "router1",
					RouterId:   "1.1.1.1",
				},
			},
		},
	}

	// Create exporter
	exporter := exporter.NewAnomalyExporter(anomalies, registry, testLogger)

	// Run concurrent updates
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exporter.Update()
		}()
	}
	wg.Wait()

	// Verify counter is correct (should be 1, not 10)
	metrics, err := registry.Gather()
	assert.NoError(t, err)

	for _, metric := range metrics {
		if *metric.Name == "ospf_overadvertised_routes_total" {
			assert.Equal(t, 1.0, metric.Metric[0].Gauge.GetValue())
		}
	}
}
