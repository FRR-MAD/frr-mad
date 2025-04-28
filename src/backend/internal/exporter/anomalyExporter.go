package exporter

import (
	"sync"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/prometheus/client_golang/prometheus"
)

type AnomalyExporter struct {
	anomalies     *frrProto.Anomalies
	activeAlerts  map[string]bool
	gauges        map[string]prometheus.Gauge
	alertCounters map[string]prometheus.Gauge
	logger        *logger.Logger
	mutex         sync.Mutex
}

func NewAnomalyExporter(anomalies *frrProto.Anomalies, registry prometheus.Registerer, logger *logger.Logger) *AnomalyExporter {
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
		{"ospf_overadvertised_route_present", "1 if overadvertised routes exist, 0 otherwise"},
		{"ospf_underadvertised_route_present", "1 if underadvertised routes exist, 0 otherwise"},
		{"ospf_duplicate_route_present", "1 if duplicate routes exist, 0 otherwise"},
		{"ospf_misconfigured_route_present", "1 if misconfigured routes exist, 0 otherwise"},
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
		{"ospf_underadvertised_routes_total", "Total underadvertised routes detected"},
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
	overCount := len(a.anomalies.GetOveradvertisedRoutes())
	a.gauges["ospf_overadvertised_route_present"].Set(boolToFloat(overCount > 0))
	a.alertCounters["ospf_overadvertised_routes_total"].Set(float64(overCount))

	// Process underadvertised routes
	underCount := len(a.anomalies.GetUnderadvertisedRoutes())
	a.gauges["ospf_underadvertised_route_present"].Set(boolToFloat(underCount > 0))
	a.alertCounters["ospf_underadvertised_routes_total"].Set(float64(underCount))

	// Process duplicate routes
	dupCount := len(a.anomalies.GetDuplicateRoutes())
	a.gauges["ospf_duplicate_route_present"].Set(boolToFloat(dupCount > 0))
	a.alertCounters["ospf_duplicate_routes_total"].Set(float64(dupCount))

	// Process misconfigured routes
	misCount := len(a.anomalies.GetMisconfiguredRoutes())
	a.gauges["ospf_misconfigured_route_present"].Set(boolToFloat(misCount > 0))
	a.alertCounters["ospf_misconfigured_routes_total"].Set(float64(misCount))
}

func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}
