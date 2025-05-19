package exporter

import (
	"sync"

	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/prometheus/client_golang/prometheus"
)

type AnomalyExporter struct {
	anomalies     *frrProto.AnomalyAnalysis
	activeAlerts  map[string]bool
	gauges        map[string]prometheus.Gauge
	alertCounters map[string]prometheus.Gauge
	logger        *logger.Logger
	mutex         sync.Mutex
}

func NewAnomalyExporter(anomalies *frrProto.AnomalyAnalysis, registry prometheus.Registerer, logger *logger.Logger) *AnomalyExporter {
	a := &AnomalyExporter{
		anomalies:     anomalies,
		activeAlerts:  make(map[string]bool),
		gauges:        make(map[string]prometheus.Gauge),
		alertCounters: make(map[string]prometheus.Gauge),
		logger:        logger,
	}

	// Initialize anomaly presence gauges
	anomalyTypes := []struct {
		name string
		help string
	}{
		{"ospf_overadvertised_route_present", "1: overadvertised routes exist, 0: otherwise"},
		{"ospf_unadvertised_route_present", "1: unadvertised routes exist, 0: otherwise"},
		{"ospf_duplicate_route_present", "1: duplicate routes exist, 0: otherwise"},
		{"ospf_misconfigured_route_present", "1: misconfigured routes exist, 0: otherwise"},
	}

	for _, at := range anomalyTypes {
		g := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: at.name,
			Help: at.help,
		})
		registry.MustRegister(g)
		a.gauges[at.name] = g
	}

	// Initialize anomaly counters
	counterTypes := []struct {
		name string
		help string
	}{
		{"ospf_overadvertised_routes_total", "Total overadvertised routes detected"},
		{"ospf_unadvertised_routes_total", "Total unadvertised routes detected"},
		{"ospf_duplicate_routes_total", "Total duplicate routes detected"},
		{"ospf_misconfigured_routes_total", "Total misconfigured routes detected"},
	}

	for _, ct := range counterTypes {
		c := prometheus.NewGauge(prometheus.GaugeOpts{
			Name: ct.name,
			Help: ct.help,
		})
		registry.MustRegister(c)
		a.alertCounters[ct.name] = c
	}

	return a
}

func (a *AnomalyExporter) Update() {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	if a.anomalies == nil {
		return
	}

	// Process overadvertised routes
	overCount :=
		len(a.anomalies.ExternalAnomaly.GetSuperfluousEntries()) +
			len(a.anomalies.RouterAnomaly.GetSuperfluousEntries()) +
			len(a.anomalies.NssaExternalAnomaly.GetSuperfluousEntries())
	a.gauges["ospf_overadvertised_route_present"].Set(boolToFloat(overCount > 0))
	a.alertCounters["ospf_overadvertised_routes_total"].Set(float64(overCount))

	// Process unadvertised routes
	underCount :=
		len(a.anomalies.RouterAnomaly.GetMissingEntries()) +
			len(a.anomalies.ExternalAnomaly.GetMissingEntries()) +
			len(a.anomalies.NssaExternalAnomaly.GetMissingEntries())
	a.gauges["ospf_unadvertised_route_present"].Set(boolToFloat(underCount > 0))
	a.alertCounters["ospf_unadvertised_routes_total"].Set(float64(underCount))

	// Process duplicate routes
	dupCount :=
		len(a.anomalies.RouterAnomaly.GetDuplicateEntries()) +
			len(a.anomalies.ExternalAnomaly.GetDuplicateEntries()) +
			len(a.anomalies.NssaExternalAnomaly.GetDuplicateEntries())
	a.gauges["ospf_duplicate_route_present"].Set(boolToFloat(dupCount > 0))
	a.alertCounters["ospf_duplicate_routes_total"].Set(float64(dupCount))

	// TODO: not implemented yet
	// Process misconfigured routes
	// misCount:=
	// len(a.anomalies.RouterAnomaly.()) +
	// len(a.anomalies.ExternalAnomaly.()) +
	// len(a.anomalies.NssaExternalAnomaly.())
	// a.gauges["ospf_misconfigured_route_present"].Set(boolToFloat(misCount > 0))
	// a.alertCounters["ospf_misconfigured_routes_total"].Set(float64(misCount))
}

func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
