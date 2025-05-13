package exporter_test

import (
	"sync"
	"testing"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/exporter"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/ba2025-ysmprc/frr-mad/src/logger"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func TestAnomalyExporter_NoAnomalies(t *testing.T) {
	// Setup
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)
	anomalies := &frrProto.AnomalyAnalysis{}

	// Create frrMadExporter
	frrMadExporter := exporter.NewAnomalyExporter(anomalies, registry, testLogger)

	// Test
	frrMadExporter.Update()

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
	testLogger, err := logger.NewLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)
	now := timestamppb.Now()

	anomalyResult := &frrProto.AnomalyAnalysis{}

	//anomalies := &frrProto.Anomalies{
	_ = &frrProto.Anomalies{
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
		UnadvertisedRoutes: []*frrProto.AnomalyUnadvertisedRoute{
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

	// Create frrMadExporter
	frrMadExporter := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)

	// Test
	frrMadExporter.Update()

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
	testLogger, err := logger.NewLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)
	now := timestamppb.Now()

	anomalyResult := &frrProto.AnomalyAnalysis{}
	//anomalies := &frrProto.Anomalies{
	_ = &frrProto.Anomalies{
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

	// Create frrMadExporter
	frrMadExporter := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)

	// Run concurrent updates
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			frrMadExporter.Update()
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

func TestAnomalyExporter_NilAnomalies_DoesNothing(t *testing.T) {
	registry := prometheus.NewRegistry()
	testLogger, _ := logger.NewLogger("test", "")
	//var anomalies *frrProto.Anomalies // nil
	var anomalyResult *frrProto.AnomalyAnalysis

	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)
	// Pre-gather to get the empty registry
	before, _ := registry.Gather()
	exp.Update()
	after, _ := registry.Gather()
	assert.Equal(t, before, after, "The Update on nil anomalies must not register or set anything")
}

func TestAnomalyExporter_ToggleAnomalies(t *testing.T) {
	// Setup
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewLogger("test", "/tmp/exporter_toggle.log")
	assert.NoError(t, err)

	// Start with some anomalies
	anomalyResult := &frrProto.AnomalyAnalysis{}
	//anomalies := &frrProto.Anomalies{
	_ = &frrProto.Anomalies{
		OveradvertisedRoutes: []*frrProto.AnomalyOveradvertisedRoute{{}, {}},
		UnadvertisedRoutes:   []*frrProto.AnomalyUnadvertisedRoute{{}},
		DuplicateRoutes:      []*frrProto.AnomalyDuplicateRoute{{}},
		MisconfiguredRoutes:  []*frrProto.AnomalyMisconfiguredRoute{{}},
	}

	frrMadExporter := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)

	// First update: should see presence=1 and correct counts
	frrMadExporter.Update()
	metrics1, err := registry.Gather()
	assert.NoError(t, err)

	getVal := func(metrics []*dto.MetricFamily, name string) (float64, bool) {
		for _, metricFamily := range metrics {
			if *metricFamily.Name == name {
				if len(metricFamily.Metric) > 0 {
					return metricFamily.Metric[0].Gauge.GetValue(), true
				}
				return 0, true
			}
		}
		return 0, false
	}

	// presence gauges should be 1
	expectedPresences := map[string]float64{
		"ospf_overadvertised_route_present":  1.0,
		"ospf_underadvertised_route_present": 1.0,
		"ospf_duplicate_route_present":       1.0,
		"ospf_misconfigured_route_present":   1.0,
	}
	for name, want := range expectedPresences {
		got, ok := getVal(metrics1, name)
		assert.True(t, ok, "expected %s to be registered", name)
		assert.Equal(t, want, got, "wrong value for %s", name)
	}

	// total counters should match slice lengths
	expectedTotals := map[string]float64{
		"ospf_overadvertised_routes_total":  2.0,
		"ospf_underadvertised_routes_total": 1.0,
		"ospf_duplicate_routes_total":       1.0,
		"ospf_misconfigured_routes_total":   1.0,
	}
	for name, want := range expectedTotals {
		got, ok := getVal(metrics1, name)
		assert.True(t, ok, "expected %s to be registered", name)
		assert.Equal(t, want, got, "wrong value for %s", name)
	}

	// Clear all anomalies
	//anomalies.OveradvertisedRoutes = nil
	//anomalies.UnderadvertisedRoutes = nil
	//anomalies.DuplicateRoutes = nil
	//anomalies.MisconfiguredRoutes = nil

	anomalyResult.ExternalAnomaly.Reset()
	anomalyResult.RouterAnomaly.Reset()
	anomalyResult.NssaExternalAnomaly.Reset()
	anomalyResult.RibToFibAnomaly.Reset()

	// Second update: everything should reset to 0
	frrMadExporter.Update()
	metrics2, err := registry.Gather()
	assert.NoError(t, err)

	for name := range expectedPresences {
		got, ok := getVal(metrics2, name)
		assert.True(t, ok, "metric %s should still be registered", name)
		assert.Equal(t, 0.0, got, "expected %s to reset to 0", name)
	}
	for name := range expectedTotals {
		got, ok := getVal(metrics2, name)
		assert.True(t, ok, "metric %s should still be registered", name)
		assert.Equal(t, 0.0, got, "expected %s to reset to 0", name)
	}
}

func TestAnomalyExporter_NoAnomalies_Existence(t *testing.T) {
	// Setup
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewLogger("test", "/tmp/exporter_no_anom_exist.log")
	assert.NoError(t, err)

	// Empty Anomalies struct -> all slices zero-length
	//anomalies := &frrProto.Anomalies{}
	anomalyResult := &frrProto.AnomalyAnalysis{}

	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)
	exp.Update()

	// Gather all metric families
	mfs, err := registry.Gather()
	assert.NoError(t, err)

	// We expect these eight gauges to be registered, each with one sample set to 0
	expected := []string{
		// presence gauges
		"ospf_overadvertised_route_present",
		"ospf_underadvertised_route_present",
		"ospf_duplicate_route_present",
		"ospf_misconfigured_route_present",
		// total counters
		"ospf_overadvertised_routes_total",
		"ospf_underadvertised_routes_total",
		"ospf_duplicate_routes_total",
		"ospf_misconfigured_routes_total",
	}

	for _, name := range expected {
		var fam *dto.MetricFamily
		for _, mf := range mfs {
			if *mf.Name == name {
				fam = mf
				break
			}
		}
		// 1) metric family is registered
		assert.NotNil(t, fam, "expected metric family %q to be registered", name)
		if fam == nil {
			continue
		}
		// 2) exactly one sample
		assert.Len(t, fam.Metric, 1,
			"expected exactly one Metric in family %q, got %d", name, len(fam.Metric))
		// 3) that sample's gauge value is zero
		val := fam.Metric[0].GetGauge().GetValue()
		assert.Equal(t, 0.0, val,
			"expected %q gauge value to be 0.0 when no anomalies, got %v", name, val)
	}
}
