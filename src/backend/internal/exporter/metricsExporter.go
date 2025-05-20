package exporter

import (
	"strings"
	"sync"

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
}

func NewMetricExporter(
	data *frrProto.FullFRRData,
	registry prometheus.Registerer,
	logger *logger.Logger,
	flags map[string]*ParsedFlag,
) *MetricExporter {
	m := &MetricExporter{
		data:           data,
		metrics:        make(map[string]prometheus.Collector),
		enabledMetrics: make(map[string]bool),
		logger:         logger,
	}

	// Initialize all metrics based on config flags
	m.initializeRouterMetrics(flags)
	m.initializeNetworkMetrics(flags)
	m.initializeSummaryMetrics(flags)
	m.initializeASBRSummaryMetrics(flags)
	m.initializeExternalMetrics(flags)
	m.initializeNSSAExternalMetrics(flags)
	m.initializeDatabaseMetrics(flags)
	m.initializeNeighborMetrics(flags)
	m.initializeInterfaceMetrics(flags)
	m.initializeRouteMetrics(flags)

	// Register all enabledMetrics metrics
	for _, metric := range m.metrics {
		registry.MustRegister(metric)
	}

	return m
}

func (m *MetricExporter) initializeRouterMetrics(flags map[string]*ParsedFlag) {
	if flag, ok := flags["OSPFRouterData"]; ok && flag.Enabled {
		m.enabledMetrics["router"] = true
		m.metrics["ospf_router_links"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_mad_ospf_router_links_total",
				Help: "Number of router interfaces in OSPF",
			},
			[]string{"area_id", "link_state_id"},
		)
	}
}

func (m *MetricExporter) initializeNetworkMetrics(flags map[string]*ParsedFlag) {
	if flag, ok := flags["OSPFNetworkData"]; ok && flag.Enabled {
		m.enabledMetrics["network"] = true
		m.metrics["ospf_network_attached_routers"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_mad_ospf_network_attached_routers_total",
				Help: "Number of attached routers announced in network LSA",
			},
			[]string{"area_id", "link_state_id"},
		)
	}
}

func (m *MetricExporter) initializeSummaryMetrics(flags map[string]*ParsedFlag) {
	if flag, ok := flags["OSPFSummaryData"]; ok && flag.Enabled {
		m.enabledMetrics["summary"] = true
		m.metrics["ospf_summary_metric"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_mad_ospf_summary_metric",
				Help: "OSPF summary LSA metric",
			},
			[]string{"area_id", "link_state_id"},
		)
	}
}

func (m *MetricExporter) initializeASBRSummaryMetrics(flags map[string]*ParsedFlag) {
	if flag, ok := flags["OSPFAsbrSummaryData"]; ok && flag.Enabled {
		m.enabledMetrics["asbr_summary"] = true
		m.metrics["ospf_asbr_summary_metric"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_mad_ospf_asbr_summary_metric",
				Help: "OSPF ASBR summary LSA metric",
			},
			[]string{"area_id", "link_state_id"},
		)
	}
}

func (m *MetricExporter) initializeExternalMetrics(flags map[string]*ParsedFlag) {
	if flag, ok := flags["OSPFExternalData"]; ok && flag.Enabled {
		m.enabledMetrics["external"] = true
		m.metrics["ospf_external_metric"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_mad_ospf_external_metric",
				Help: "OSPF external LSA route metric",
			},
			[]string{"link_state_id", "metric_type"},
		)
	}
}

func (m *MetricExporter) initializeNSSAExternalMetrics(flags map[string]*ParsedFlag) {
	if flag, ok := flags["OSPFNssaExternalData"]; ok && flag.Enabled {
		m.enabledMetrics["nssa_external"] = true
		m.metrics["ospf_nssa_external_metric"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_mad_ospf_nssa_external_metric",
				Help: "OSPF NSSA external LSA route metric",
			},
			[]string{"area_id", "link_state_id", "metric_type"},
		)
	}
}

func (m *MetricExporter) initializeDatabaseMetrics(flags map[string]*ParsedFlag) {
	if flag, ok := flags["OSPFDatabase"]; ok && flag.Enabled {
		m.enabledMetrics["database"] = true
		m.metrics["ospf_database_counts"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_mad_ospf_database_lsa_count",
				Help: "Amount of LSDB entries for each LSA type",
			},
			[]string{"area_id", "lsa_type"},
		)
	}
}

func (m *MetricExporter) initializeNeighborMetrics(flags map[string]*ParsedFlag) {
	if flag, ok := flags["OSPFNeighbors"]; ok && flag.Enabled {
		m.enabledMetrics["neighbors"] = true
		m.metrics["ospf_neighbor_state"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_mad_ospf_neighbor_state",
				Help: "OSPF neighbor state (1=Full, 0.5=2-Way, 0=Down)",
			},
			[]string{"neighbor_id", "interface"},
		)
		m.metrics["ospf_neighbor_uptime"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_mad_ospf_neighbor_uptime_seconds",
				Help: "OSPF neighbor uptime in seconds",
			},
			[]string{"neighbor_id", "interface"},
		)
	}
}

func (m *MetricExporter) initializeInterfaceMetrics(flags map[string]*ParsedFlag) {
	if flag, ok := flags["InterfaceList"]; ok && flag.Enabled {
		m.enabledMetrics["interfaces"] = true
		m.metrics["interface_operational_status"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_mad_interface_operational_status",
				Help: "Network interface operational status (1=Up, 0=Down)",
			},
			[]string{"interface", "vrf"},
		)
		m.metrics["interface_admin_status"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_mad_interface_admin_status",
				Help: "Network interface administrative status (1=Up, 0=Down)",
			},
			[]string{"interface", "vrf"},
		)
	}
}

func (m *MetricExporter) initializeRouteMetrics(flags map[string]*ParsedFlag) {
	if flag, ok := flags["RouteList"]; ok && flag.Enabled {
		m.enabledMetrics["routes"] = true
		m.metrics["installed_ospf_route"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_mad_installed_ospf_route",
				Help: "Routing protocol metric for installed ospf routes",
			},
			[]string{"prefix", "protocol", "vrf"},
		)
		m.metrics["installed_ospf_routes_count"] = prometheus.NewGauge(
			prometheus.GaugeOpts{
				Name: "frr_mad_installed_ospf_routes_count",
				Help: "Number of installed ospf routes from RIB",
			},
		)
	}
}

func (m *MetricExporter) Update() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.data == nil {
		return
	}

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
