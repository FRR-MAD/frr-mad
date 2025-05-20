package exporter

import (
	"sync"

	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/prometheus/client_golang/prometheus"
)

type AnomalyExporter struct {
	anomalies      *frrProto.AnomalyAnalysis
	activeAlerts   map[string]bool
	anomalyDetails *prometheus.GaugeVec
	anomalyFlags   *prometheus.GaugeVec
	alertCounters  map[string]prometheus.Gauge
	logger         *logger.Logger
	mutex          sync.Mutex
	knownLabelSets map[string]prometheus.Labels
}

func NewAnomalyExporter(anomalies *frrProto.AnomalyAnalysis, registry prometheus.Registerer, logger *logger.Logger) *AnomalyExporter {
	a := &AnomalyExporter{
		anomalies:      anomalies,
		activeAlerts:   make(map[string]bool),
		alertCounters:  make(map[string]prometheus.Gauge),
		logger:         logger,
		knownLabelSets: make(map[string]prometheus.Labels),
	}

	// Initialize anomaly details metric
	a.anomalyDetails = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "frr_mad_anomaly_details",
			Help: "Detailed information about anomalies (1=present, 0=absent)",
		},
		[]string{
			"anomaly_type", // overadvertised, unadvertised, duplicate, etc.
			"source",       // RouterAnomaly, ExternalAnomaly, NssaExternalAnomaly, RibToFib, LsdbToRib
			"interface_address",
			"link_state_id",
			"prefix_length",
			"link_type",
			"p_bit",
			"options",
		},
	)
	registry.MustRegister(a.anomalyDetails)

	a.anomalyFlags = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "frr_mad_anomaly_flags",
			Help: "Flag indicators for anomaly types (1=present, 0=absent)",
		},
		[]string{
			"source",
			"flag_type", // overadvertised, unadvertised, duplicate, misconfigured
		},
	)
	registry.MustRegister(a.anomalyFlags)

	counterTypes := []struct {
		name string
		help string
	}{
		{"frr_mad_ospf_overadvertised_routes_total", "Total overadvertised routes detected in OSPF"},
		{"frr_mad_ospf_unadvertised_routes_total", "Total unadvertised routes detected in OSPF"},
		{"frr_mad_ospf_duplicate_routes_total", "Total duplicate routes detected in OSPF"},
		{"frr_mad_ospf_misconfigured_routes_total", "Total misconfigured routes detected in OSPF"},
		{"frr_mad_rib_to_fib_anomalies_total", "Total RIB to FIB anomalies detected"},
		{"frr_mad_lsdb_to_rib_anomalies_total", "Total LSDB to RIB anomalies detected"},
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

	for _, labels := range a.knownLabelSets {
		a.anomalyDetails.With(labels).Set(0)
	}
	a.knownLabelSets = make(map[string]prometheus.Labels)

	for _, counter := range a.alertCounters {
		counter.Set(0)
	}

	for _, source := range []string{"RouterAnomaly", "ExternalAnomaly", "NssaExternalAnomaly", "RibToFib", "LsdbToRib"} {
		for _, flag := range []string{"overadvertised", "unadvertised", "duplicate", "misconfigured"} {
			a.anomalyFlags.WithLabelValues(source, flag).Set(0)
		}
	}

	if a.anomalies == nil {
		return
	}

	// OSPF anomalies
	a.processOspfSources()

	// RIB to FIB anomalies
	if ribToFib := a.anomalies.RibToFibAnomaly; ribToFib != nil {
		a.alertCounters["frr_mad_rib_to_fib_anomalies_total"].Set(float64(
			len(ribToFib.GetSuperfluousEntries()) + len(ribToFib.GetMissingEntries()) + len(ribToFib.GetDuplicateEntries()),
		))

		a.anomalyFlags.WithLabelValues("RibToFib", "overadvertised").Set(boolToFloat(ribToFib.GetHasOverAdvertisedPrefixes()))
		a.anomalyFlags.WithLabelValues("RibToFib", "unadvertised").Set(boolToFloat(ribToFib.GetHasUnAdvertisedPrefixes()))
		a.anomalyFlags.WithLabelValues("RibToFib", "duplicate").Set(boolToFloat(ribToFib.GetHasDuplicatePrefixes()))
		a.anomalyFlags.WithLabelValues("RibToFib", "misconfigured").Set(boolToFloat(ribToFib.GetHasMisconfiguredPrefixes()))

		for _, entry := range ribToFib.GetSuperfluousEntries() {
			a.setAnomalyDetail("overadvertised", "RibToFib", entry)
		}
		for _, entry := range ribToFib.GetMissingEntries() {
			a.setAnomalyDetail("unadvertised", "RibToFib", entry)
		}
		for _, entry := range ribToFib.GetDuplicateEntries() {
			a.setAnomalyDetail("duplicate", "RibToFib", entry)
		}
	}

	// LSDB to RIB anomalies
	if lsdbToRib := a.anomalies.LsdbToRibAnomaly; lsdbToRib != nil {
		a.alertCounters["frr_mad_lsdb_to_rib_anomalies_total"].Set(float64(
			len(lsdbToRib.GetSuperfluousEntries()) + len(lsdbToRib.GetMissingEntries()) + len(lsdbToRib.GetDuplicateEntries()),
		))

		a.anomalyFlags.WithLabelValues("LsdbToRib", "overadvertised").Set(boolToFloat(lsdbToRib.GetHasOverAdvertisedPrefixes()))
		a.anomalyFlags.WithLabelValues("LsdbToRib", "unadvertised").Set(boolToFloat(lsdbToRib.GetHasUnAdvertisedPrefixes()))
		a.anomalyFlags.WithLabelValues("LsdbToRib", "duplicate").Set(boolToFloat(lsdbToRib.GetHasDuplicatePrefixes()))
		a.anomalyFlags.WithLabelValues("LsdbToRib", "misconfigured").Set(boolToFloat(lsdbToRib.GetHasMisconfiguredPrefixes()))

		for _, entry := range lsdbToRib.GetSuperfluousEntries() {
			a.setAnomalyDetail("overadvertised", "LsdbToRib", entry)
		}
		for _, entry := range lsdbToRib.GetMissingEntries() {
			a.setAnomalyDetail("unadvertised", "LsdbToRib", entry)
		}
		for _, entry := range lsdbToRib.GetDuplicateEntries() {
			a.setAnomalyDetail("duplicate", "LsdbToRib", entry)
		}
	}
}

func (a *AnomalyExporter) processOspfSources() {
	var (
		totalOver      int
		totalUnder     int
		totalDup       int
		totalMisconfig int
	)

	processSource := func(source string, detection *frrProto.AnomalyDetection) {
		if detection == nil {
			return
		}

		a.anomalyFlags.WithLabelValues(source, "overadvertised").Set(boolToFloat(detection.GetHasOverAdvertisedPrefixes()))
		a.anomalyFlags.WithLabelValues(source, "unadvertised").Set(boolToFloat(detection.GetHasUnAdvertisedPrefixes()))
		a.anomalyFlags.WithLabelValues(source, "duplicate").Set(boolToFloat(detection.GetHasDuplicatePrefixes()))
		a.anomalyFlags.WithLabelValues(source, "misconfigured").Set(boolToFloat(detection.GetHasMisconfiguredPrefixes()))

		over := detection.GetSuperfluousEntries()
		under := detection.GetMissingEntries()
		dup := detection.GetDuplicateEntries()

		totalOver += len(over)
		totalUnder += len(under)
		totalDup += len(dup)

		for _, ad := range over {
			a.setAnomalyDetail("overadvertised", source, ad)
		}
		for _, ad := range under {
			a.setAnomalyDetail("unadvertised", source, ad)
		}
		for _, ad := range dup {
			a.setAnomalyDetail("duplicate", source, ad)
		}

		if detection.GetHasMisconfiguredPrefixes() {
			totalMisconfig++
		}
	}

	processSource("RouterAnomaly", a.anomalies.RouterAnomaly)
	processSource("ExternalAnomaly", a.anomalies.ExternalAnomaly)
	processSource("NssaExternalAnomaly", a.anomalies.NssaExternalAnomaly)

	a.alertCounters["frr_mad_ospf_overadvertised_routes_total"].Set(float64(totalOver))
	a.alertCounters["frr_mad_ospf_unadvertised_routes_total"].Set(float64(totalUnder))
	a.alertCounters["frr_mad_ospf_duplicate_routes_total"].Set(float64(totalDup))
	a.alertCounters["frr_mad_ospf_misconfigured_routes_total"].Set(float64(totalMisconfig))
}

func (a *AnomalyExporter) setAnomalyDetail(anomalyType, source string, ad *frrProto.Advertisement) {
	labels := prometheus.Labels{
		"anomaly_type":      anomalyType,
		"source":            source,
		"interface_address": ad.GetInterfaceAddress(),
		"link_state_id":     ad.GetLinkStateId(),
		"prefix_length":     ad.GetPrefixLength(),
		"link_type":         ad.GetLinkType(),
		"p_bit":             boolToString(ad.GetPBit()),
		"options":           ad.GetOptions(),
	}

	key := anomalyType + ":" + source + ":" + ad.GetInterfaceAddress() + ":" + ad.GetLinkStateId()
	a.knownLabelSets[key] = labels
	a.anomalyDetails.With(labels).Set(1)
}

func boolToFloat(b bool) float64 {
	if b {
		return 1
	}
	return 0
}

func boolToString(b bool) string {
	if b {
		return "true"
	}
	return "false"
}
