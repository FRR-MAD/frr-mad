package exporter_test

import (
	"sync"
	"testing"

	"github.com/frr-mad/frr-mad/src/backend/internal/exporter"
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestAnomalyExporter_NoAnomalies(t *testing.T) {
	// Setup
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)
	anomalies := &frrProto.AnomalyAnalysis{}

	// Create exporter
	exp := exporter.NewAnomalyExporter(anomalies, registry, testLogger)

	// Test
	exp.Update()

	// Verify all metrics are 0
	metrics, err := registry.Gather()
	assert.NoError(t, err)

	for _, metric := range metrics {
		switch *metric.Name {
		case "frr_mad_anomaly_details":
			// Should have no metrics when no anomalies
			assert.Equal(t, 0, len(metric.Metric))
		case "frr_mad_anomaly_flags":
			// All flags should be 0
			for _, m := range metric.Metric {
				assert.Equal(t, 0.0, m.Gauge.GetValue())
			}
		case "frr_mad_ospf_overadvertised_routes_total":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		case "frr_mad_ospf_unadvertised_routes_total":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		case "frr_mad_ospf_duplicate_routes_total":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		case "frr_mad_ospf_misconfigured_routes_total":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		case "frr_mad_rib_to_fib_anomalies_total":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		case "frr_mad_lsdb_to_rib_anomalies_total":
			assert.Equal(t, 0.0, metric.Metric[0].Gauge.GetValue())
		}
	}
}

func TestAnomalyExporter_WithAnomalies(t *testing.T) {
	// Setup
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)

	// Create test data
	anomalyResult := &frrProto.AnomalyAnalysis{
		RouterAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes: true,
			HasUnAdvertisedPrefixes:   true,
			HasDuplicatePrefixes:      true,
			HasMisconfiguredPrefixes:  true,
			SuperfluousEntries: []*frrProto.Advertisement{
				{
					InterfaceAddress: "10.0.0.1",
					LinkStateId:      "1.1.1.1",
					PrefixLength:     "24",
					LinkType:         "Stub",
					PBit:             true,
					Options:          "0x02",
				},
			},
			MissingEntries: []*frrProto.Advertisement{
				{
					InterfaceAddress: "10.1.0.1",
					LinkStateId:      "2.2.2.2",
					PrefixLength:     "24",
					LinkType:         "Stub",
				},
			},
			DuplicateEntries: []*frrProto.Advertisement{
				{
					InterfaceAddress: "172.16.0.1",
					LinkStateId:      "3.3.3.3",
					PrefixLength:     "24",
					LinkType:         "Stub",
				},
			},
		},
		RibToFibAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes: true,
			SuperfluousEntries: []*frrProto.Advertisement{
				{
					InterfaceAddress: "192.168.1.1",
					LinkStateId:      "4.4.4.4",
					PrefixLength:     "24",
					LinkType:         "Stub",
				},
			},
		},
	}

	// Create exporter
	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)

	// Test
	exp.Update()

	// Verify metrics
	metrics, err := registry.Gather()
	assert.NoError(t, err)

	// Helper to get metric value
	getMetricValue := func(name string) float64 {
		for _, metric := range metrics {
			if *metric.Name == name {
				if len(metric.Metric) > 0 {
					return metric.Metric[0].Gauge.GetValue()
				}
				return 0
			}
		}
		return -1
	}

	// Check counters
	assert.Equal(t, 1.0, getMetricValue("frr_mad_ospf_overadvertised_routes_total"))
	assert.Equal(t, 1.0, getMetricValue("frr_mad_ospf_unadvertised_routes_total"))
	assert.Equal(t, 1.0, getMetricValue("frr_mad_ospf_duplicate_routes_total"))
	assert.Equal(t, 1.0, getMetricValue("frr_mad_ospf_misconfigured_routes_total"))
	assert.Equal(t, 1.0, getMetricValue("frr_mad_rib_to_fib_anomalies_total"))
	assert.Equal(t, 0.0, getMetricValue("frr_mad_lsdb_to_rib_anomalies_total"))

	// Check flags
	assert.Equal(t, 1.0, getMetricValue("frr_mad_anomaly_flags"))
}

func TestAnomalyExporter_ConcurrentUpdates(t *testing.T) {
	// Setup
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)

	anomalyResult := &frrProto.AnomalyAnalysis{
		RouterAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes: true,
			SuperfluousEntries: []*frrProto.Advertisement{
				{
					InterfaceAddress: "10.0.0.1",
					LinkStateId:      "1.1.1.1",
					PrefixLength:     "24",
					LinkType:         "Stub",
				},
			},
		},
	}

	// Create exporter
	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)

	// Run concurrent updates
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exp.Update()
		}()
	}
	wg.Wait()

	// Verify counter is correct (should be 1, not 10)
	metrics, err := registry.Gather()
	assert.NoError(t, err)

	for _, metric := range metrics {
		if *metric.Name == "frr_mad_ospf_overadvertised_routes_total" {
			assert.Equal(t, 1.0, metric.Metric[0].Gauge.GetValue())
		}
	}
}

func TestAnomalyExporter_NilAnomalies_DoesNothing(t *testing.T) {
	registry := prometheus.NewRegistry()
	testLogger, _ := logger.NewLogger("test", "")
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
	anomalyResult := &frrProto.AnomalyAnalysis{
		RouterAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes: true,
			SuperfluousEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "10.0.0.1", PrefixLength: "24"},
			},
		},
	}

	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)

	// First update: should see metrics
	exp.Update()
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

	// Check flags and counters
	overFlag, ok := getVal(metrics1, "frr_mad_anomaly_flags")
	assert.True(t, ok)
	assert.Equal(t, 1.0, overFlag)

	overCount, ok := getVal(metrics1, "frr_mad_ospf_overadvertised_routes_total")
	assert.True(t, ok)
	assert.Equal(t, 1.0, overCount)

	// Clear anomalies
	anomalyResult.RouterAnomaly.HasOverAdvertisedPrefixes = false
	anomalyResult.RouterAnomaly.SuperfluousEntries = nil

	// Second update: should reset to 0
	exp.Update()
	metrics2, err := registry.Gather()
	assert.NoError(t, err)

	overFlag, ok = getVal(metrics2, "frr_mad_anomaly_flags")
	assert.True(t, ok)
	assert.Equal(t, 0.0, overFlag)

	overCount, ok = getVal(metrics2, "frr_mad_ospf_overadvertised_routes_total")
	assert.True(t, ok)
	assert.Equal(t, 0.0, overCount)
}

func TestAnomalyExporter_NoAnomalies_Existence(t *testing.T) {
	// Setup
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewLogger("test", "/tmp/exporter_no_anom_exist.log")
	assert.NoError(t, err)

	// Empty AnomalyAnalysis struct
	anomalyResult := &frrProto.AnomalyAnalysis{}

	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)
	exp.Update()

	// Gather all metric families
	mfs, err := registry.Gather()
	assert.NoError(t, err)

	// We expect these metrics to be registered
	expected := []string{
		"frr_mad_anomaly_details",
		"frr_mad_anomaly_flags",
		"frr_mad_ospf_overadvertised_routes_total",
		"frr_mad_ospf_unadvertised_routes_total",
		"frr_mad_ospf_duplicate_routes_total",
		"frr_mad_ospf_misconfigured_routes_total",
		"frr_mad_rib_to_fib_anomalies_total",
		"frr_mad_lsdb_to_rib_anomalies_total",
	}

	for _, name := range expected {
		var found bool
		for _, mf := range mfs {
			if *mf.Name == name {
				found = true
				break
			}
		}
		assert.True(t, found, "expected metric %q to be registered", name)
	}
}
