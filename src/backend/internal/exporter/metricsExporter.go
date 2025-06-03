package exporter

import (
	"strings"
	"sync"
	"time"

	"github.com/frr-mad/frr-mad/src/backend/configs"
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/prometheus/client_golang/prometheus"
)

type MetricExporter struct {
	data           *frrProto.FullFRRData
	metrics        map[string]prometheus.Collector
	enabledMetrics map[string]bool
	logger         *logger.Logger
	mutex          sync.RWMutex
	config         configs.ExporterConfig
}

func NewMetricExporter(
	data *frrProto.FullFRRData,
	registry prometheus.Registerer,
	logger *logger.Logger,
	config configs.ExporterConfig,
) *MetricExporter {
	logger.Debug("Initializing metric exporter")

	m := &MetricExporter{
		data:           data,
		metrics:        make(map[string]prometheus.Collector),
		enabledMetrics: make(map[string]bool),
		logger:         logger,
		config:         config,
	}

	// Initialize all metrics based on config flags
	m.initializeMetrics()

	// Register all enabledMetrics metrics
	for _, metric := range m.metrics {
		registry.MustRegister(metric)
	}

	enabled := []string{}
	for k := range m.enabledMetrics {
		enabled = append(enabled, k)
	}
	logger.WithAttrs(map[string]interface{}{
		"enabled_metrics": enabled,
		"total_metrics":   len(m.metrics),
	}).Info("Metric exporter initialized")

	return m
}

func (m *MetricExporter) initializeMetrics() {
	// Use a struct to define metric configurations
	type metricConfig struct {
		configFlag   bool
		metricName   string
		metricKey    string
		helpText     string
		labels       []string
		extraMetrics []struct {
			name   string
			help   string
			labels []string
		}
	}

	metricsConfig := map[string]metricConfig{
		"router": {
			configFlag: m.config.OSPFRouterData,
			metricName: "frr_mad_ospf_router_links_total",
			metricKey:  "ospf_router_links",
			helpText:   "Number of router interfaces in OSPF",
			labels:     []string{"area_id", "link_state_id"},
		},
		"network": {
			configFlag: m.config.OSPFNetworkData,
			metricName: "frr_mad_ospf_network_attached_routers_total",
			metricKey:  "ospf_network_attached_routers",
			helpText:   "Number of attached routers announced in network LSA",
			labels:     []string{"area_id", "link_state_id"},
		},
		"summary": {
			configFlag: m.config.OSPFSummaryData,
			metricName: "frr_mad_ospf_summary_metric",
			metricKey:  "ospf_summary_metric",
			helpText:   "OSPF summary LSA metric",
			labels:     []string{"area_id", "link_state_id"},
		},
		"asbr_summary": {
			configFlag: m.config.OSPFAsbrSummaryData,
			metricName: "frr_mad_ospf_asbr_summary_metric",
			metricKey:  "ospf_asbr_summary_metric",
			helpText:   "OSPF ASBR summary LSA metric",
			labels:     []string{"area_id", "link_state_id"},
		},
		"external": {
			configFlag: m.config.OSPFExternalData,
			metricName: "frr_mad_ospf_external_metric",
			metricKey:  "ospf_external_metric",
			helpText:   "OSPF external LSA route metric",
			labels:     []string{"link_state_id", "metric_type"},
		},
		"nssa_external": {
			configFlag: m.config.OSPFNssaExternalData,
			metricName: "frr_mad_ospf_nssa_external_metric",
			metricKey:  "ospf_nssa_external_metric",
			helpText:   "OSPF NSSA external LSA route metric",
			labels:     []string{"area_id", "link_state_id", "metric_type"},
		},
		"database": {
			configFlag: m.config.OSPFDatabase,
			metricName: "frr_mad_ospf_database_lsa_count",
			metricKey:  "ospf_database_counts",
			helpText:   "Amount of LSDB entries for each LSA type",
			labels:     []string{"area_id", "lsa_type"},
		},
		"neighbors": {
			configFlag: m.config.OSPFNeighbors,
			metricKey:  "neighbors",
			extraMetrics: []struct {
				name   string
				help   string
				labels []string
			}{
				{
					name:   "ospf_neighbor_state",
					help:   "OSPF neighbor state (1=Full, 0.5=2-Way, 0=Down)",
					labels: []string{"neighbor_id", "interface"},
				},
				{
					name:   "ospf_neighbor_uptime",
					help:   "OSPF neighbor uptime in seconds",
					labels: []string{"neighbor_id", "interface"},
				},
			},
		},
		"interfaces": {
			configFlag: m.config.InterfaceList,
			metricKey:  "interfaces",
			extraMetrics: []struct {
				name   string
				help   string
				labels []string
			}{
				{
					name:   "interface_operational_status",
					help:   "Network interface operational status (1=Up, 0=Down)",
					labels: []string{"interface", "vrf"},
				},
				{
					name:   "interface_admin_status",
					help:   "Network interface administrative status (1=Up, 0=Down)",
					labels: []string{"interface", "vrf"},
				},
			},
		},
		"routes": {
			configFlag: m.config.RouteList,
			metricKey:  "routes",
			extraMetrics: []struct {
				name   string
				help   string
				labels []string
			}{
				{
					name:   "installed_ospf_route",
					help:   "Routing protocol metric for installed ospf routes",
					labels: []string{"prefix", "protocol", "vrf"},
				},
				{
					name: "installed_ospf_routes_count",
					help: "Number of installed ospf routes from RIB",
				},
			},
		},
	}

	for key, cfg := range metricsConfig {
		if !cfg.configFlag {
			continue
		}

		m.enabledMetrics[key] = true

		if cfg.metricName != "" {
			m.metrics[cfg.metricKey] = prometheus.NewGaugeVec(
				prometheus.GaugeOpts{
					Name: cfg.metricName,
					Help: cfg.helpText,
				},
				cfg.labels,
			)
		}

		for _, extra := range cfg.extraMetrics {
			if len(extra.labels) > 0 {
				m.metrics[extra.name] = prometheus.NewGaugeVec(
					prometheus.GaugeOpts{
						Name: "frr_mad_" + extra.name,
						Help: extra.help,
					},
					extra.labels,
				)
			} else {
				m.metrics[extra.name] = prometheus.NewGauge(
					prometheus.GaugeOpts{
						Name: "frr_mad_" + extra.name,
						Help: extra.help,
					},
				)
			}
		}
	}
}

func (m *MetricExporter) Update() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data == nil {
		m.logger.Warning("Skipping metric update - no data available")
		return
	}

	m.logger.Debug("Starting metric update")
	start := time.Now()

	// Update all enabledMetrics
	if m.enabledMetrics["router"] {
		m.updateRouterMetrics()
	}
	if m.enabledMetrics["network"] {
		m.updateNetworkMetrics()
	}
	if m.enabledMetrics["summary"] {
		m.updateSummaryMetrics()
	}
	if m.enabledMetrics["asbr_summary"] {
		m.updateASBRSummaryMetrics()
	}
	if m.enabledMetrics["external"] {
		m.updateExternalMetrics()
	}
	if m.enabledMetrics["nssa_external"] {
		m.updateNSSAExternalMetrics()
	}
	if m.enabledMetrics["database"] {
		m.updateDatabaseMetrics()
	}
	if m.enabledMetrics["neighbors"] {
		m.updateNeighborMetrics()
	}
	if m.enabledMetrics["interfaces"] {
		m.updateInterfaceMetrics()
	}
	if m.enabledMetrics["routes"] {
		m.updateRouteMetrics()
	}

	m.logger.WithAttrs(map[string]interface{}{
		"duration": time.Since(start).String(),
	}).Debug("Completed metric update")
}

func (m *MetricExporter) updateRouterMetrics() {
	if routerData := m.data.GetOspfRouterData(); routerData != nil {
		vec := m.metrics["ospf_router_links"].(*prometheus.GaugeVec)
		vec.Reset()

		for areaID, areaData := range routerData.RouterStates {
			for linkStateID, lsa := range areaData.LsaEntries {
				vec.WithLabelValues(areaID, linkStateID).Set(float64(lsa.NumOfLinks))
			}
		}
	}
}

func (m *MetricExporter) updateNetworkMetrics() {
	if networkData := m.data.GetOspfNetworkData(); networkData != nil {
		vec := m.metrics["ospf_network_attached_routers"].(*prometheus.GaugeVec)

		vec.Reset()

		for areaID, areaData := range networkData.NetStates {
			for linkStateID, lsa := range areaData.LsaEntries {
				vec.WithLabelValues(areaID, linkStateID).Set(float64(len(lsa.AttachedRouters)))
			}
		}
	}
}

func (m *MetricExporter) updateSummaryMetrics() {
	if summaryData := m.data.GetOspfSummaryData(); summaryData != nil {
		vec := m.metrics["ospf_summary_metric"].(*prometheus.GaugeVec)
		vec.Reset()

		for areaID, areaData := range summaryData.SummaryStates {
			for linkStateID, lsa := range areaData.LsaEntries {
				vec.WithLabelValues(areaID, linkStateID).Set(float64(lsa.Tos0Metric))
			}
		}
	}
}

func (m *MetricExporter) updateASBRSummaryMetrics() {
	if asbrData := m.data.GetOspfAsbrSummaryData(); asbrData != nil {
		vec := m.metrics["ospf_asbr_summary_metric"].(*prometheus.GaugeVec)
		vec.Reset()

		for areaID, areaData := range asbrData.AsbrSummaryStates {
			for linkStateID, lsa := range areaData.LsaEntries {
				vec.WithLabelValues(areaID, linkStateID).Set(float64(lsa.Tos0Metric))
			}
		}
	}
}

func (m *MetricExporter) updateExternalMetrics() {
	if externalData := m.data.GetOspfExternalData(); externalData != nil {
		vec := m.metrics["ospf_external_metric"].(*prometheus.GaugeVec)
		vec.Reset()

		for linkStateID, lsa := range externalData.AsExternalLinkStates {
			vec.WithLabelValues(linkStateID, lsa.MetricType).Set(float64(lsa.Metric))
		}
	}
}

func (m *MetricExporter) updateNSSAExternalMetrics() {
	if nssaData := m.data.GetOspfNssaExternalData(); nssaData != nil {
		vec := m.metrics["ospf_nssa_external_metric"].(*prometheus.GaugeVec)
		vec.Reset()

		for areaID, areaData := range nssaData.NssaExternalLinkStates {
			for linkStateID, lsa := range areaData.Data {
				vec.WithLabelValues(areaID, linkStateID, lsa.MetricType).Set(float64(lsa.Metric))
			}
		}
	}
}

func (m *MetricExporter) updateDatabaseMetrics() {
	if dbData := m.data.GetOspfDatabase(); dbData != nil {
		vec := m.metrics["ospf_database_counts"].(*prometheus.GaugeVec)
		vec.Reset()

		for areaID, area := range dbData.Areas {
			vec.WithLabelValues(areaID, "router").Set(float64(area.RouterLinkStatesCount))
			vec.WithLabelValues(areaID, "network").Set(float64(area.NetworkLinkStatesCount))
			vec.WithLabelValues(areaID, "summary").Set(float64(area.SummaryLinkStatesCount))
			vec.WithLabelValues(areaID, "asbr_summary").Set(float64(area.AsbrSummaryLinkStatesCount))
		}
		vec.WithLabelValues("0", "external").Set(float64(dbData.AsExternalCount))
	}
}

func (m *MetricExporter) updateNeighborMetrics() {
	if neighborData := m.data.GetOspfNeighbors(); neighborData != nil {
		stateVec := m.metrics["ospf_neighbor_state"].(*prometheus.GaugeVec)
		uptimeVec := m.metrics["ospf_neighbor_uptime"].(*prometheus.GaugeVec)
		stateVec.Reset()
		uptimeVec.Reset()

		for iface, neighborList := range neighborData.Neighbors {
			for _, neighbor := range neighborList.Neighbors {
				stateValue := 0.0
				switch {
				case strings.Contains(neighbor.NbrState, "Full"):
					stateValue = 1.0
				case strings.Contains(neighbor.NbrState, "2-Way"):
					stateValue = 0.5
				}

				stateVec.WithLabelValues(neighbor.Address, iface).Set(stateValue)
				uptimeVec.WithLabelValues(neighbor.Address, iface).Set(float64(neighbor.UpTimeInMsec / 1000))
			}
		}
	}
}

func (m *MetricExporter) updateInterfaceMetrics() {
	if ifaceData := m.data.GetInterfaces(); ifaceData != nil {
		operVec := m.metrics["interface_operational_status"].(*prometheus.GaugeVec)
		adminVec := m.metrics["interface_admin_status"].(*prometheus.GaugeVec)
		operVec.Reset()
		adminVec.Reset()

		for interfaceName, interfaceData := range ifaceData.Interfaces {
			operStatus := 0.0
			if interfaceData.OperationalStatus == "Up" {
				operStatus = 1.0
			}

			adminStatus := 0.0
			if interfaceData.AdministrativeStatus == "Up" {
				adminStatus = 1.0
			}

			operVec.WithLabelValues(interfaceName, interfaceData.VrfName).Set(operStatus)
			adminVec.WithLabelValues(interfaceName, interfaceData.VrfName).Set(adminStatus)
		}
	}
}

func (m *MetricExporter) updateRouteMetrics() {
	if routeData := m.data.GetRoutingInformationBase(); routeData != nil {
		vec := m.metrics["installed_ospf_route"].(*prometheus.GaugeVec)
		countMetric := m.metrics["installed_ospf_routes_count"].(prometheus.Gauge)
		vec.Reset()

		var counter float64
		for vrf, routeEntry := range routeData.Routes {
			for _, route := range routeEntry.Routes {
				if route.Installed && route.Protocol == "ospf" {
					counter++
					vec.WithLabelValues(route.Prefix, route.Protocol, vrf).Set(float64(route.Metric))
				}
			}
		}
		countMetric.Set(counter)
	}
}
