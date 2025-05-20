package exporter

import (
	"fmt"
	"sync"

	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/prometheus/client_golang/prometheus"
)

type AnomalyExporter struct {
	anomalies          *frrProto.AnomalyAnalysis
	activeAlerts       map[string]bool
	ospfAnomalyDetails *prometheus.GaugeVec
	alertCounters      map[string]prometheus.Gauge
	logger             *logger.Logger
	mutex              sync.Mutex
	knownLabelSets     map[string]prometheus.Labels
}

func NewAnomalyExporter(anomalies *frrProto.AnomalyAnalysis, registry prometheus.Registerer, logger *logger.Logger) *AnomalyExporter {
	a := &AnomalyExporter{
		anomalies:      anomalies,
		activeAlerts:   make(map[string]bool),
		alertCounters:  make(map[string]prometheus.Gauge),
		logger:         logger,
		knownLabelSets: make(map[string]prometheus.Labels),
	}

	a.ospfAnomalyDetails = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "ospf_anomaly_details",
			Help: "Detailed information about OSPF anomalies (1=present, 0=absent)",
		},
		[]string{
			"anomaly_type",
			"source",
			"interface_address",
			"link_state_id",
			"prefix_length",
			"link_type",
		},
	)
	registry.MustRegister(a.ospfAnomalyDetails)

	counterTypes := []struct {
		name string
		help string
	}{
		{"ospf_overadvertised_routes_total", "Total overadvertised routes detected across all sources"},
		{"ospf_unadvertised_routes_total", "Total unadvertised routes detected across all sources"},
		{"ospf_duplicate_routes_total", "Total duplicate routes detected across all sources"},
		{"rib_to_fib_anomalies_total", "Total RIB to FIB anomalies detected"},
		{"lsdb_to_rib_anomalies_total", "Total LSDB to RIB anomalies detected"},
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
		a.ospfAnomalyDetails.With(labels).Set(0)
	}
	a.knownLabelSets = make(map[string]prometheus.Labels)

	for _, counter := range a.alertCounters {
		counter.Set(0)
	}

	if a.anomalies == nil {
		return
	}

	var (
		totalOver  int
		totalUnder int
		totalDup   int
	)

	processSource := func(source string, over, under, dup []*frrProto.Advertisement) {
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
	}

	fmt.Println(a.anomalies.RouterAnomaly)

	processSource("ExternalAnomaly",
		a.anomalies.ExternalAnomaly.GetSuperfluousEntries(),
		a.anomalies.ExternalAnomaly.GetMissingEntries(),
		a.anomalies.ExternalAnomaly.GetDuplicateEntries())

	processSource("RouterAnomaly",
		a.anomalies.RouterAnomaly.GetSuperfluousEntries(),
		a.anomalies.RouterAnomaly.GetMissingEntries(),
		a.anomalies.RouterAnomaly.GetDuplicateEntries())

	processSource("NssaExternalAnomaly",
		a.anomalies.NssaExternalAnomaly.GetSuperfluousEntries(),
		a.anomalies.NssaExternalAnomaly.GetMissingEntries(),
		a.anomalies.NssaExternalAnomaly.GetDuplicateEntries())

	a.alertCounters["ospf_overadvertised_routes_total"].Set(float64(totalOver))
	a.alertCounters["ospf_unadvertised_routes_total"].Set(float64(totalUnder))
	a.alertCounters["ospf_duplicate_routes_total"].Set(float64(totalDup))

	// RIB to FIB anomalies
	if ribToFib := a.anomalies.RibToFibAnomaly; ribToFib != nil {
		a.alertCounters["rib_to_fib_anomalies_total"].Set(
			float64(len(ribToFib.GetSuperfluousEntries()) + len(ribToFib.GetMissingEntries())),
		)
	}

	// LSDB to RIB anomalies
	if lsdbToRib := a.anomalies.LsdbToRibAnomaly; lsdbToRib != nil {
		a.alertCounters["lsdb_to_rib_anomalies_total"].Set(
			float64(len(lsdbToRib.GetSuperfluousEntries()) + len(lsdbToRib.GetMissingEntries())),
		)
	}
}

func (a *AnomalyExporter) setAnomalyDetail(anomalyType, source string, ad *frrProto.Advertisement) {
	labels := prometheus.Labels{
		"anomaly_type":      anomalyType,
		"source":            source,
		"interface_address": ad.GetInterfaceAddress(),
		"link_state_id":     ad.GetLinkStateId(),
		"prefix_length":     ad.GetPrefixLength(),
		"link_type":         ad.GetLinkType(),
	}

	key := anomalyType + source + ad.GetInterfaceAddress() + ad.GetLinkStateId() + ad.GetPrefixLength() + ad.GetLinkType()

	a.knownLabelSets[key] = labels

	a.ospfAnomalyDetails.With(labels).Set(1)
}
