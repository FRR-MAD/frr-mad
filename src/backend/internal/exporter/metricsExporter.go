package exporter

import (
	"sync"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricExporter struct {
	data    *frrProto.FullFRRData
	metrics map[string]prometheus.Collector
	enabled map[string]bool
	logger  *logger.Logger
	mutex   sync.RWMutex
}

func newMetricExporter(
	data *frrProto.FullFRRData,
	registry prometheus.Registerer,
	logger *logger.Logger,
	flags map[string]configs.ParsedFlag,
) *MetricExporter {
	m := &MetricExporter{
		data:    data,
		metrics: make(map[string]prometheus.Collector),
		enabled: make(map[string]bool),
		logger:  logger,
	}

	// Initialize metrics based on enabled flags
	if flag, ok := flags["OSPFRouterData"]; ok && flag.Enabled {
		m.enabled["router"] = true
		m.metrics["ospf_router_links"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_ospf_router_links_total",
				Help: "Number of router links in OSPF",
			},
			[]string{"area_id", "link_state_id"},
		)
	}

	if flag, ok := flags["OSPFNetworkData"]; ok && flag.Enabled {
		m.enabled["network"] = true
		m.metrics["ospf_network_attached_routers"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_ospf_network_attached_routers_total",
				Help: "Number of routers attached to OSPF network",
			},
			[]string{"area_id", "link_state_id"},
		)
	}

	// ... initialize all other metrics based on config flags

	// Register all enabled metrics
	for _, metric := range m.metrics {
		registry.MustRegister(metric)
	}

	return m
}

func (m *MetricExporter) Update() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data == nil {
		return
	}

	// Update OSPF Router Data metrics
	if m.enabled["router"] && m.data.GetOspfRouterData() != nil {
		vec := m.metrics["ospf_router_links"].(*prometheus.GaugeVec)
		vec.Reset()

		for areaID, areaData := range m.data.GetOspfRouterData().RouterStates {
			for linkStateID, lsa := range areaData.LsaEntries {
				vec.WithLabelValues(areaID, linkStateID).Set(float64(lsa.NumOfLinks))
			}
		}
	}

	// Update OSPF Network Data metrics
	if m.enabled["network"] && m.data.GetOspfNetworkData() != nil {
		vec := m.metrics["ospf_network_attached_routers"].(*prometheus.GaugeVec)
		vec.Reset()

		for areaID, areaData := range m.data.GetOspfNetworkData().NetStates {
			for linkStateID, lsa := range areaData.LsaEntries {
				vec.WithLabelValues(areaID, linkStateID).Set(float64(len(lsa.AttachedRouters)))
			}
		}
	}

	// ... update all other enabled metrics
}
