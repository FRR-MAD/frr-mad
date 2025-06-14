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
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)
	anomalies := &frrProto.AnomalyAnalysis{}

	exp := exporter.NewAnomalyExporter(anomalies, registry, testLogger)
	exp.Update()

	metrics, err := registry.Gather()
	assert.NoError(t, err)

	// Verify all counters are 0
	for _, name := range []string{
		"frr_mad_ospf_overadvertised_routes_total",
		"frr_mad_ospf_unadvertised_routes_total",
		"frr_mad_ospf_duplicate_routes_total",
		"frr_mad_ospf_misconfigured_routes_total",
		"frr_mad_rib_to_fib_anomalies_total",
		"frr_mad_lsdb_to_rib_anomalies_total",
	} {
		val := getMetricValue(metrics, name)
		assert.Equal(t, 0.0, val, "expected %s to be 0", name)
	}

	// Verify all flags are 0
	flagMetrics := getMetricFamily(metrics, "frr_mad_anomaly_flags")
	if assert.NotNil(t, flagMetrics, "anomaly_flags metric should exist") {
		for _, m := range flagMetrics.Metric {
			assert.Equal(t, 0.0, m.GetGauge().GetValue(), "flag should be 0")
		}
	}
}

func TestAnomalyExporter_WithAnomalies(t *testing.T) {
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)

	anomalyResult := &frrProto.AnomalyAnalysis{
		RouterAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes: true,
			HasUnAdvertisedPrefixes:   true,
			HasDuplicatePrefixes:      true,
			HasMisconfiguredPrefixes:  true,
			SuperfluousEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "10.0.0.1"},
				{InterfaceAddress: "192.168.1.1"},
			},
			MissingEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "10.1.0.1"},
			},
			DuplicateEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "172.16.0.1"},
			},
		},
	}

	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)
	exp.Update()

	metrics, err := registry.Gather()
	assert.NoError(t, err)

	// Check counters
	assert.Equal(t, 2.0, getMetricValue(metrics, "frr_mad_ospf_overadvertised_routes_total"))
	assert.Equal(t, 1.0, getMetricValue(metrics, "frr_mad_ospf_unadvertised_routes_total"))
	assert.Equal(t, 1.0, getMetricValue(metrics, "frr_mad_ospf_duplicate_routes_total"))
	assert.Equal(t, 1.0, getMetricValue(metrics, "frr_mad_ospf_misconfigured_routes_total"))

	// Check flags
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "RouterAnomaly", "flag_type": "overadvertised"}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "RouterAnomaly", "flag_type": "unadvertised"}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "RouterAnomaly", "flag_type": "duplicate"}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "RouterAnomaly", "flag_type": "misconfigured"}))

	// Check details
	details := getMetricFamily(metrics, "frr_mad_anomaly_details")
	assert.NotNil(t, details)
	assert.Greater(t, len(details.Metric), 0, "should have anomaly details")
}

func TestAnomalyExporter_ConcurrentUpdates(t *testing.T) {
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)

	anomalyResult := &frrProto.AnomalyAnalysis{
		RouterAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes: true,
			SuperfluousEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "10.0.0.1"},
			},
		},
	}

	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			exp.Update()
		}()
	}
	wg.Wait()

	metrics, err := registry.Gather()
	assert.NoError(t, err)
	assert.Equal(t, 1.0, getMetricValue(metrics, "frr_mad_ospf_overadvertised_routes_total"))
}

func TestAnomalyExporter_NilAnomalies_DoesNothing(t *testing.T) {
	registry := prometheus.NewRegistry()
	testLogger, _ := logger.NewApplicationLogger("test", "/tmp/exporter_toggle.log")
	var anomalyResult *frrProto.AnomalyAnalysis

	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)

	// First update to establish initialized state
	exp.Update()

	// Gather metrics after first update
	before, err := registry.Gather()
	assert.NoError(t, err)

	// Map metrics to their values for comparison
	beforeValues := getMetricValues(before)

	// Update again with nil anomalies
	exp.Update()

	// Gather metrics after second update
	after, err := registry.Gather()
	assert.NoError(t, err)
	afterValues := getMetricValues(after)

	// All metrics should remain at 0
	for name, val := range beforeValues {
		afterVal, exists := afterValues[name]
		assert.True(t, exists, "metric %s should exist after update", name)
		assert.Equal(t, val, afterVal, "metric %s value should remain the same", name)
	}
}

func TestAnomalyExporter_ToggleAnomalies(t *testing.T) {
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/exporter_toggle.log")
	assert.NoError(t, err)

	anomalyResult := &frrProto.AnomalyAnalysis{
		RouterAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes: true,
			SuperfluousEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "10.0.0.1"},
			},
		},
	}

	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)
	exp.Update()

	// Verify anomaly is present
	metrics, err := registry.Gather()
	assert.NoError(t, err)
	assert.Equal(t, 1.0, getMetricValue(metrics, "frr_mad_ospf_overadvertised_routes_total"))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "RouterAnomaly", "flag_type": "overadvertised"}))

	// Clear anomalies
	anomalyResult.RouterAnomaly.HasOverAdvertisedPrefixes = false
	anomalyResult.RouterAnomaly.SuperfluousEntries = nil
	exp.Update()

	// Verify anomaly is cleared
	metrics, err = registry.Gather()
	assert.NoError(t, err)
	assert.Equal(t, 0.0, getMetricValue(metrics, "frr_mad_ospf_overadvertised_routes_total"))
	assert.Equal(t, 0.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "RouterAnomaly", "flag_type": "overadvertised"}))
}

func TestAnomalyExporter_NoAnomalies_Existence(t *testing.T) {
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/exporter_no_anom_exist.log")
	assert.NoError(t, err)

	anomalyResult := &frrProto.AnomalyAnalysis{}
	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)
	exp.Update()

	metrics, err := registry.Gather()
	assert.NoError(t, err)

	// Map of metric names to expected metric type
	requiredMetrics := map[string]dto.MetricType{
		"frr_mad_anomaly_details":                  dto.MetricType_GAUGE,
		"frr_mad_anomaly_flags":                    dto.MetricType_GAUGE,
		"frr_mad_ospf_overadvertised_routes_total": dto.MetricType_GAUGE,
		"frr_mad_ospf_unadvertised_routes_total":   dto.MetricType_GAUGE,
		"frr_mad_ospf_duplicate_routes_total":      dto.MetricType_GAUGE,
		"frr_mad_ospf_misconfigured_routes_total":  dto.MetricType_GAUGE,
		"frr_mad_rib_to_fib_anomalies_total":       dto.MetricType_GAUGE,
		"frr_mad_lsdb_to_rib_anomalies_total":      dto.MetricType_GAUGE,
	}

	// Check for existence of each required metric
	for name, expectedType := range requiredMetrics {
		metricFamily := getMetricFamily(metrics, name)
		assert.NotNil(t, metricFamily, "metric %s should exist", name)
		if metricFamily != nil {
			assert.Equal(t, expectedType, *metricFamily.Type, "metric %s should be of type %v", name, expectedType)
		}
	}

	// Check that all flag combinations exist when there are no anomalies
	flagMetrics := getMetricFamily(metrics, "frr_mad_anomaly_flags")
	if assert.NotNil(t, flagMetrics, "anomaly_flags metric should exist") {
		// We should have 5 sources × 4 flag types = 20 metrics
		assert.Equal(t, 20, len(flagMetrics.Metric),
			"should have metrics for all source/flag combinations")

		// All flags should be 0 as there are no anomalies
		for _, m := range flagMetrics.Metric {
			assert.Equal(t, 0.0, m.GetGauge().GetValue(), "flag should be 0")
		}
	}
}

func TestAnomalyExporter_MixedAnomalies(t *testing.T) {
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)

	anomalyResult := &frrProto.AnomalyAnalysis{
		RouterAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes: true,
			SuperfluousEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "10.0.0.1"},
			},
		},
		LsdbToRibAnomaly: &frrProto.AnomalyDetection{
			HasUnAdvertisedPrefixes: true,
			MissingEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "192.168.1.1"},
			},
		},
	}

	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)
	exp.Update()

	metrics, err := registry.Gather()
	assert.NoError(t, err)

	// Verify router overadvertised anomalies
	assert.Equal(t, 1.0, getMetricValue(metrics, "frr_mad_ospf_overadvertised_routes_total"))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "RouterAnomaly", "flag_type": "overadvertised"}))

	// Verify LSDB-to-RIB unadvertised anomalies
	assert.Equal(t, 1.0, getMetricValue(metrics, "frr_mad_lsdb_to_rib_anomalies_total"))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "LsdbToRib", "flag_type": "unadvertised"}))

	// Check details metrics
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_details",
		map[string]string{
			"anomaly_type":      "overadvertised",
			"source":            "RouterAnomaly",
			"interface_address": "10.0.0.1",
		}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_details",
		map[string]string{
			"anomaly_type":      "unadvertised",
			"source":            "LsdbToRib",
			"interface_address": "192.168.1.1",
		}))
}

func TestAnomalyExporter_RibToFibAnomalies(t *testing.T) {
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)

	anomalyResult := &frrProto.AnomalyAnalysis{
		RibToFibAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes: true,
			HasUnAdvertisedPrefixes:   true,
			HasDuplicatePrefixes:      true,
			HasMisconfiguredPrefixes:  true,
			SuperfluousEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "10.0.0.1", LinkStateId: "1.1.1.1", PrefixLength: "24"},
				{InterfaceAddress: "10.0.0.2", LinkStateId: "2.2.2.2", PrefixLength: "24"},
			},
			MissingEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "192.168.1.1", LinkStateId: "3.3.3.3", PrefixLength: "24"},
			},
			DuplicateEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "172.16.0.1", LinkStateId: "4.4.4.4", PrefixLength: "24"},
			},
		},
	}

	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)
	exp.Update()

	metrics, err := registry.Gather()
	assert.NoError(t, err)

	// Check total anomalies counter
	assert.Equal(t, 4.0, getMetricValue(metrics, "frr_mad_rib_to_fib_anomalies_total"))

	// Check flags
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "RibToFib", "flag_type": "overadvertised"}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "RibToFib", "flag_type": "unadvertised"}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "RibToFib", "flag_type": "duplicate"}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "RibToFib", "flag_type": "misconfigured"}))

	// Check details for each type of anomaly
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_details",
		map[string]string{
			"anomaly_type":      "overadvertised",
			"source":            "RibToFib",
			"interface_address": "10.0.0.1",
			"link_state_id":     "1.1.1.1",
		}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_details",
		map[string]string{
			"anomaly_type":      "overadvertised",
			"source":            "RibToFib",
			"interface_address": "10.0.0.2",
			"link_state_id":     "2.2.2.2",
		}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_details",
		map[string]string{
			"anomaly_type":      "unadvertised",
			"source":            "RibToFib",
			"interface_address": "192.168.1.1",
			"link_state_id":     "3.3.3.3",
		}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_details",
		map[string]string{
			"anomaly_type":      "duplicate",
			"source":            "RibToFib",
			"interface_address": "172.16.0.1",
			"link_state_id":     "4.4.4.4",
		}))
}

func TestAnomalyExporter_LsdbToRibAnomalies(t *testing.T) {
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewApplicationLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)

	anomalyResult := &frrProto.AnomalyAnalysis{
		LsdbToRibAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes: true,
			HasDuplicatePrefixes:      true,
			SuperfluousEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "10.1.0.1", LinkStateId: "5.5.5.5", PrefixLength: "24"},
			},
			DuplicateEntries: []*frrProto.Advertisement{
				{InterfaceAddress: "10.2.0.1", LinkStateId: "6.6.6.6", PrefixLength: "24"},
				{InterfaceAddress: "10.2.0.2", LinkStateId: "7.7.7.7", PrefixLength: "24"},
			},
		},
	}

	exp := exporter.NewAnomalyExporter(anomalyResult, registry, testLogger)
	exp.Update()

	metrics, err := registry.Gather()
	assert.NoError(t, err)

	// Check total anomalies counter
	assert.Equal(t, 3.0, getMetricValue(metrics, "frr_mad_lsdb_to_rib_anomalies_total"))

	// Check flags
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "LsdbToRib", "flag_type": "overadvertised"}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_flags",
		map[string]string{"source": "LsdbToRib", "flag_type": "duplicate"}))

	// Check details
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_details",
		map[string]string{
			"anomaly_type":      "overadvertised",
			"source":            "LsdbToRib",
			"interface_address": "10.1.0.1",
			"link_state_id":     "5.5.5.5",
		}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_details",
		map[string]string{
			"anomaly_type":      "duplicate",
			"source":            "LsdbToRib",
			"interface_address": "10.2.0.1",
			"link_state_id":     "6.6.6.6",
		}))
	assert.Equal(t, 1.0, getMetricValueWithLabels(metrics, "frr_mad_anomaly_details",
		map[string]string{
			"anomaly_type":      "duplicate",
			"source":            "LsdbToRib",
			"interface_address": "10.2.0.2",
			"link_state_id":     "7.7.7.7",
		}))
}

func TestBoolToString(t *testing.T) {
	assert.Equal(t, "true", exporter.BoolToString(true))
	assert.Equal(t, "false", exporter.BoolToString(false))
}

// Helper functions

func getMetricValue(metrics []*dto.MetricFamily, name string) float64 {
	for _, mf := range metrics {
		if mf.GetName() == name {
			if len(mf.Metric) > 0 {
				return mf.Metric[0].GetGauge().GetValue()
			}
			return 0
		}
	}
	return -1
}

func getMetricValueWithLabels(metrics []*dto.MetricFamily, name string, labels map[string]string) float64 {
	for _, mf := range metrics {
		if mf.GetName() == name {
			for _, m := range mf.Metric {
				match := true
				for _, lp := range m.Label {
					if expected, ok := labels[lp.GetName()]; ok {
						if lp.GetValue() != expected {
							match = false
							break
						}
					}
				}
				if match {
					return m.GetGauge().GetValue()
				}
			}
		}
	}
	return -1
}

func getMetricFamily(metrics []*dto.MetricFamily, name string) *dto.MetricFamily {
	for _, mf := range metrics {
		if mf.GetName() == name {
			return mf
		}
	}
	return nil
}

func getMetricValues(metrics []*dto.MetricFamily) map[string]float64 {
	result := make(map[string]float64)

	for _, mf := range metrics {
		name := mf.GetName()

		// For simple gauges with no labels
		if len(mf.Metric) == 1 && len(mf.Metric[0].Label) == 0 {
			result[name] = mf.Metric[0].GetGauge().GetValue()
			continue
		}

		// For metrics with labels
		for _, m := range mf.Metric {
			labelStr := ""
			for _, l := range m.Label {
				labelStr += l.GetName() + "=" + l.GetValue() + ","
			}
			if labelStr != "" {
				result[name+"{"+labelStr+"}"] = m.GetGauge().GetValue()
			} else {
				result[name] = m.GetGauge().GetValue()
			}
		}
	}

	return result
}
