package exporter

import (
	"fmt"
	"sync"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/prometheus/client_golang/prometheus"
)

type OSPFMetricsExporter struct {
	registry *prometheus.Registry
	metrics  map[string]*prometheus.GaugeVec
	enabled  map[string]bool
	logger   *logger.Logger
	mutex    sync.RWMutex
}

var (
	ospfMetricsSingleton *OSPFMetricsExporter
	ospfMetricsOnce      sync.Once
)

// GetOSPFMetricsExporter returns a singleton instance of OSPFMetricsExporter
func GetOSPFMetricsExporter() *OSPFMetricsExporter {
	ospfMetricsOnce.Do(func() {
		ospfMetricsSingleton = &OSPFMetricsExporter{
			registry: prometheus.NewRegistry(),
			metrics:  make(map[string]*prometheus.GaugeVec),
			enabled:  make(map[string]bool),
		}
	})
	return ospfMetricsSingleton
}

func (e *OSPFMetricsExporter) Init(logger *logger.Logger, config map[string]map[string]string) error {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.logger = logger

	flags, err := configs.GetFlagConfigs(config, "exporter")
	if err != nil {
		return fmt.Errorf("error loading config flags: %v", err)
	}

	e.initializeMetrics(flags)
	return nil
}

func (e *OSPFMetricsExporter) initializeMetrics(flags map[string]configs.ParsedFlag) {
	// Router Metrics
	if flag, ok := flags["OSPFRouterData"]; ok && flag.Enabled {
		e.enabled["router"] = true
		e.metrics["ospf_router_links"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_ospf_router_links_total",
				Help: "Number of router links in OSPF",
			},
			[]string{"area_id", "link_state_id"},
		)
		e.registry.MustRegister(e.metrics["ospf_router_links"])
	}

	// Network Metrics
	if flag, ok := flags["OSPFNetworkData"]; ok && flag.Enabled {
		e.enabled["network"] = true
		e.metrics["ospf_network_attached_routers"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_ospf_network_attached_routers_total",
				Help: "Number of routers attached to OSPF network",
			},
			[]string{"area_id", "link_state_id"},
		)
		e.registry.MustRegister(e.metrics["ospf_network_attached_routers"])
	}

	// Summary Metrics
	if flag, ok := flags["OSPFSummaryData"]; ok && flag.Enabled {
		e.enabled["summary"] = true
		e.metrics["ospf_summary_metric"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_ospf_summary_metric",
				Help: "OSPF summary LSA metric",
			},
			[]string{"area_id", "link_state_id"},
		)
		e.registry.MustRegister(e.metrics["ospf_summary_metric"])
	}

	// ASBR Summary Metrics
	if flag, ok := flags["OSPFAsbrSummaryData"]; ok && flag.Enabled {
		e.enabled["asbr_summary"] = true
		e.metrics["ospf_asbr_summary_metric"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_ospf_asbr_summary_metric",
				Help: "OSPF ASBR summary LSA metric",
			},
			[]string{"area_id", "link_state_id"},
		)
		e.registry.MustRegister(e.metrics["ospf_asbr_summary_metric"])
	}

	// External Metrics
	if flag, ok := flags["OSPFExternalData"]; ok && flag.Enabled {
		e.enabled["external"] = true
		e.metrics["ospf_external_metric"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_ospf_external_metric",
				Help: "OSPF external route metric",
			},
			[]string{"link_state_id", "metric_type"},
		)
		e.registry.MustRegister(e.metrics["ospf_external_metric"])
	}

	// NSSA External Metrics
	if flag, ok := flags["OSPFNssaExternalData"]; ok && flag.Enabled {
		e.enabled["nssa_external"] = true
		e.metrics["ospf_nssa_external_metric"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_ospf_nssa_external_metric",
				Help: "OSPF NSSA external route metric",
			},
			[]string{"area_id", "link_state_id", "metric_type"},
		)
		e.registry.MustRegister(e.metrics["ospf_nssa_external_metric"])
	}

	// Database Metrics
	if flag, ok := flags["OSPFDatabase"]; ok && flag.Enabled {
		e.enabled["database"] = true
		e.metrics["ospf_database_counts"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_ospf_database_lsa_count",
				Help: "Counts of different LSA types in OSPF database",
			},
			[]string{"area_id", "lsa_type"},
		)
		e.registry.MustRegister(e.metrics["ospf_database_counts"])
	}

	// Duplicate Metrics
	if flag, ok := flags["OSPFDuplicates"]; ok && flag.Enabled {
		e.enabled["duplicates"] = true
		e.metrics["ospf_duplicate_lsas"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_ospf_duplicate_lsa_count",
				Help: "Count of duplicate OSPF LSAs",
			},
			[]string{"link_state_id"},
		)
		e.registry.MustRegister(e.metrics["ospf_duplicate_lsas"])
	}

	// Neighbor Metrics
	if flag, ok := flags["OSPFNeighbors"]; ok && flag.Enabled {
		e.enabled["neighbors"] = true
		e.metrics["ospf_neighbor_state"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_ospf_neighbor_state",
				Help: "OSPF neighbor state (1=Full, 0.5=2-Way, 0=Down)",
			},
			[]string{"neighbor_id", "interface"},
		)
		e.metrics["ospf_neighbor_uptime"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_ospf_neighbor_uptime_seconds",
				Help: "OSPF neighbor uptime in seconds",
			},
			[]string{"neighbor_id", "interface"},
		)
		e.registry.MustRegister(e.metrics["ospf_neighbor_state"])
		e.registry.MustRegister(e.metrics["ospf_neighbor_uptime"])
	}

	// Interface Metrics
	if flag, ok := flags["InterfaceList"]; ok && flag.Enabled {
		e.enabled["interfaces"] = true
		e.metrics["interface_operational_status"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_interface_operational_status",
				Help: "Network interface operational status (1=Up, 0=Down)",
			},
			[]string{"interface", "vrf"},
		)
		e.metrics["interface_admin_status"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_interface_admin_status",
				Help: "Network interface administrative status (1=Up, 0=Down)",
			},
			[]string{"interface", "vrf"},
		)
		e.registry.MustRegister(e.metrics["interface_operational_status"])
		e.registry.MustRegister(e.metrics["interface_admin_status"])
	}

	// Route Metrics
	if flag, ok := flags["RouteList"]; ok && flag.Enabled {
		e.enabled["routes"] = true
		e.metrics["route_metric"] = prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: "frr_route_metric",
				Help: "Routing protocol metric for installed routes",
			},
			[]string{"prefix", "protocol", "vrf"},
		)
		e.registry.MustRegister(e.metrics["route_metric"])
	}
}

func (e *OSPFMetricsExporter) GetRegistry() *prometheus.Registry {
	return e.registry
}

// UpdateRouterMetrics updates metrics for OSPF router data
func (e *OSPFMetricsExporter) UpdateRouterMetrics(data *frrProto.OSPFRouterData) {
	e.logger.Info("Dies ist ein Test")
	if !e.enabled["router"] {
		return
	}

	e.logger.Info("Test ist durchgegangen")

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for areaID, areaData := range data.RouterStates {
		for linkStateID, lsa := range areaData.LsaEntries {
			e.metrics["ospf_router_links"].With(prometheus.Labels{
				"area_id":       areaID,
				"link_state_id": linkStateID,
			}).Set(float64(lsa.NumOfLinks))
		}
	}
}

// UpdateNetworkMetrics updates metrics for OSPF network data
func (e *OSPFMetricsExporter) UpdateNetworkMetrics(data *frrProto.OSPFNetworkData) {
	if !e.enabled["network"] {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for areaID, areaData := range data.NetStates {
		for linkStateID, lsa := range areaData.LsaEntries {
			e.metrics["ospf_network_attached_routers"].With(prometheus.Labels{
				"area_id":       areaID,
				"link_state_id": linkStateID,
			}).Set(float64(len(lsa.AttachedRouters)))
		}
	}
}

// UpdateSummaryMetrics updates metrics for OSPF summary data
func (e *OSPFMetricsExporter) UpdateSummaryMetrics(data *frrProto.OSPFSummaryData) {
	if !e.enabled["summary"] {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for areaID, areaData := range data.SummaryStates {
		for linkStateID, lsa := range areaData.LsaEntries {
			e.metrics["ospf_summary_metric"].With(prometheus.Labels{
				"area_id":       areaID,
				"link_state_id": linkStateID,
			}).Set(float64(lsa.Tos0Metric))
		}
	}
}

// UpdateASBRSummaryMetrics updates metrics for OSPF ASBR summary data
func (e *OSPFMetricsExporter) UpdateASBRSummaryMetrics(data *frrProto.OSPFAsbrSummaryData) {
	if !e.enabled["asbr_summary"] {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for areaID, areaData := range data.AsbrSummaryStates {
		for linkStateID, lsa := range areaData.LsaEntries {
			e.metrics["ospf_asbr_summary_metric"].With(prometheus.Labels{
				"area_id":       areaID,
				"link_state_id": linkStateID,
			}).Set(float64(lsa.Tos0Metric))
		}
	}
}

// UpdateExternalMetrics updates metrics for OSPF external data
func (e *OSPFMetricsExporter) UpdateExternalMetrics(data *frrProto.OSPFExternalData) {
	if !e.enabled["external"] {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for linkStateID, lsa := range data.AsExternalLinkStates {
		e.metrics["ospf_external_metric"].With(prometheus.Labels{
			"link_state_id": linkStateID,
			"metric_type":   lsa.MetricType,
		}).Set(float64(lsa.Metric))
	}
}

// UpdateNSSAExternalMetrics updates metrics for OSPF NSSA external data
func (e *OSPFMetricsExporter) UpdateNSSAExternalMetrics(data *frrProto.OSPFNssaExternalData) {
	if !e.enabled["nssa_external"] {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for areaID, areaData := range data.NssaExternalLinkStates {
		for linkStateID, lsa := range areaData.Data {
			e.metrics["ospf_nssa_external_metric"].With(prometheus.Labels{
				"area_id":       areaID,
				"link_state_id": linkStateID,
				"metric_type":   lsa.MetricType,
			}).Set(float64(lsa.Metric))
		}
	}
}

// UpdateDatabaseMetrics updates metrics for OSPF database
func (e *OSPFMetricsExporter) UpdateDatabaseMetrics(data *frrProto.OSPFDatabase) {
	if !e.enabled["database"] {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for areaID, area := range data.Areas {
		e.metrics["ospf_database_counts"].With(prometheus.Labels{
			"area_id":  areaID,
			"lsa_type": "router",
		}).Set(float64(area.RouterLinkStatesCount))

		e.metrics["ospf_database_counts"].With(prometheus.Labels{
			"area_id":  areaID,
			"lsa_type": "network",
		}).Set(float64(area.NetworkLinkStatesCount))

		e.metrics["ospf_database_counts"].With(prometheus.Labels{
			"area_id":  areaID,
			"lsa_type": "summary",
		}).Set(float64(area.SummaryLinkStatesCount))

		e.metrics["ospf_database_counts"].With(prometheus.Labels{
			"area_id":  areaID,
			"lsa_type": "asbr_summary",
		}).Set(float64(area.AsbrSummaryLinkStatesCount))
	}

	e.metrics["ospf_database_counts"].With(prometheus.Labels{
		"area_id":  "0", // External LSAs aren't area-specific
		"lsa_type": "external",
	}).Set(float64(data.AsExternalCount))
}

// UpdateDuplicateMetrics updates metrics for OSPF duplicates
func (e *OSPFMetricsExporter) UpdateDuplicateMetrics(data *frrProto.OSPFDuplicates) {
	if !e.enabled["duplicates"] {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for _, lsa := range data.AsExternalLinkStates {
		e.metrics["ospf_duplicate_lsas"].With(prometheus.Labels{
			"link_state_id": lsa.LinkStateId,
		}).Inc()
	}
}

// UpdateNeighborMetrics updates metrics for OSPF neighbors
func (e *OSPFMetricsExporter) UpdateNeighborMetrics(data *frrProto.OSPFNeighbors) {
	if !e.enabled["neighbors"] {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for iface, neighborList := range data.Neighbors {
		for _, neighbor := range neighborList.Neighbors {
			// Convert state to numeric value
			stateValue := 0.0
			switch neighbor.NbrState {
			case "Full":
				stateValue = 1.0
			case "2-Way":
				stateValue = 0.5
			}

			e.metrics["ospf_neighbor_state"].With(prometheus.Labels{
				"neighbor_id": neighbor.Address,
				"interface":   iface,
			}).Set(stateValue)

			e.metrics["ospf_neighbor_uptime"].With(prometheus.Labels{
				"neighbor_id": neighbor.Address,
				"interface":   iface,
			}).Set(float64(neighbor.UpTimeInMsec / 1000)) // Convert to seconds
		}
	}
}

// UpdateInterfaceMetrics updates metrics for network interfaces
func (e *OSPFMetricsExporter) UpdateInterfaceMetrics(data *frrProto.InterfaceList) {
	if !e.enabled["interfaces"] {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for iface, intf := range data.Interfaces {
		// Operational status (1=Up, 0=Down)
		operStatus := 0.0
		if intf.OperationalStatus == "Up" {
			operStatus = 1.0
		}

		// Admin status (1=Up, 0=Down)
		adminStatus := 0.0
		if intf.AdministrativeStatus == "Up" {
			adminStatus = 1.0
		}

		e.metrics["interface_operational_status"].With(prometheus.Labels{
			"interface": iface,
			"vrf":       intf.VrfName,
		}).Set(operStatus)

		e.metrics["interface_admin_status"].With(prometheus.Labels{
			"interface": iface,
			"vrf":       intf.VrfName,
		}).Set(adminStatus)
	}
}

// UpdateRouteMetrics updates metrics for routing table
func (e *OSPFMetricsExporter) UpdateRouteMetrics(data *frrProto.RouteList) {
	if !e.enabled["routes"] {
		return
	}

	e.mutex.Lock()
	defer e.mutex.Unlock()

	for vrf, routeEntry := range data.Routes {
		for _, route := range routeEntry.Routes {
			if route.Installed {
				e.metrics["route_metric"].With(prometheus.Labels{
					"prefix":   route.Prefix,
					"protocol": route.Protocol,
					"vrf":      vrf,
				}).Set(float64(route.Metric))
			}
		}
	}
}
