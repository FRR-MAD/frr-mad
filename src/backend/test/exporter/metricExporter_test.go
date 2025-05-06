package exporter_test

import (
	"fmt"
	"testing"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/exporter"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"github.com/stretchr/testify/assert"
)

func TestMetricExporter_WithData(t *testing.T) {
	// Setup
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)

	flags := map[string]configs.ParsedFlag{
		"OSPFRouterData":       {Enabled: true},
		"OSPFNetworkData":      {Enabled: true},
		"OSPFSummaryData":      {Enabled: true},
		"OSPFAsbrSummaryData":  {Enabled: true},
		"OSPFExternalData":     {Enabled: true},
		"OSPFNssaExternalData": {Enabled: true},
		"OSPFDatabase":         {Enabled: true},
		"OSPFDuplicates":       {Enabled: true},
		"OSPFNeighbors":        {Enabled: true},
		"InterfaceList":        {Enabled: true},
		"RouteList":            {Enabled: true},
	}

	// Create test data
	attachedRouters := make(map[string]*frrProto.AttachedRouter)
	attachedRouters["1.1.1.1"] = &frrProto.AttachedRouter{AttachedRouterId: "1.1.1.1"}
	attachedRouters["2.2.2.2"] = &frrProto.AttachedRouter{AttachedRouterId: "2.2.2.2"}

	data := &frrProto.FullFRRData{
		OspfRouterData: &frrProto.OSPFRouterData{
			RouterStates: map[string]*frrProto.OSPFRouterArea{
				"0.0.0.0": {
					LsaEntries: map[string]*frrProto.OSPFRouterLSA{
						"1.1.1.1": {NumOfLinks: 3},
						"2.2.2.2": {NumOfLinks: 2},
					},
				},
			},
		},
		OspfNetworkData: &frrProto.OSPFNetworkData{
			NetStates: map[string]*frrProto.NetAreaState{
				"0.0.0.0": {
					LsaEntries: map[string]*frrProto.NetworkLSA{
						"192.168.1.1": {AttachedRouters: attachedRouters},
					},
				},
			},
		},
		OspfSummaryData: &frrProto.OSPFSummaryData{
			SummaryStates: map[string]*frrProto.SummaryAreaState{
				"0.0.0.0": {
					LsaEntries: map[string]*frrProto.SummaryLSA{
						"10.0.0.0": {Tos0Metric: 10},
					},
				},
			},
		},
		OspfAsbrSummaryData: &frrProto.OSPFAsbrSummaryData{
			AsbrSummaryStates: map[string]*frrProto.SummaryAreaState{
				"0.0.0.0": {
					LsaEntries: map[string]*frrProto.SummaryLSA{
						"3.3.3.3": {Tos0Metric: 20},
					},
				},
			},
		},
		OspfExternalData: &frrProto.OSPFExternalData{
			AsExternalLinkStates: map[string]*frrProto.ExternalLSA{
				"4.4.4.4": {Metric: 30, MetricType: "E2"},
			},
		},
		OspfNssaExternalData: &frrProto.OSPFNssaExternalData{
			NssaExternalLinkStates: map[string]*frrProto.NssaExternalArea{
				"0.0.0.1": {
					Data: map[string]*frrProto.NssaExternalLSA{
						"5.5.5.5": {Metric: 40, MetricType: "E1"},
					},
				},
			},
		},
		OspfDatabase: &frrProto.OSPFDatabase{
			Areas: map[string]*frrProto.OSPFDatabaseArea{
				"0.0.0.0": {
					RouterLinkStatesCount:      5,
					NetworkLinkStatesCount:     3,
					SummaryLinkStatesCount:     2,
					AsbrSummaryLinkStatesCount: 1,
				},
			},
			AsExternalCount: 4,
		},
		OspfDuplicates: &frrProto.OSPFDuplicates{
			AsExternalLinkStates: []*frrProto.ASExternalLinkState{
				{LinkStateId: "6.6.6.6"},
				{LinkStateId: "6.6.6.6"}, // Duplicate
			},
		},
		OspfNeighbors: &frrProto.OSPFNeighbors{
			Neighbors: map[string]*frrProto.NeighborList{
				"eth0": {
					Neighbors: []*frrProto.Neighbor{
						{Address: "7.7.7.7", NbrState: "Full", UpTimeInMsec: 60000},
						{Address: "8.8.8.8", NbrState: "2-Way", UpTimeInMsec: 30000},
					},
				},
			},
		},
		Interfaces: &frrProto.InterfaceList{
			Interfaces: map[string]*frrProto.SingleInterface{
				"eth0": {
					OperationalStatus:    "Up",
					AdministrativeStatus: "Up",
					VrfName:              "default",
				},
				"eth1": {
					OperationalStatus:    "Down",
					AdministrativeStatus: "Up",
					VrfName:              "default",
				},
			},
		},
		Routes: &frrProto.RouteList{
			Routes: map[string]*frrProto.RouteEntry{
				"default": {
					Routes: []*frrProto.Route{
						{Prefix: "10.0.0.0/24", Protocol: "ospf", Metric: 100, Installed: true},
						{Prefix: "192.168.1.0/24", Protocol: "bgp", Metric: 200, Installed: true},
					},
				},
			},
		},
	}

	// Create frrMadExporter
	frrMadExporter := exporter.NewMetricExporter(data, registry, testLogger, flags)

	// Test
	frrMadExporter.Update()

	// Verify metrics
	metrics, err := registry.Gather()
	assert.NoError(t, err)

	fmt.Println(metrics)

	// Helper function to find metric value
	getMetricValue := func(name string, labels map[string]string) float64 {
		for _, metric := range metrics {
			if *metric.Name == name {
				for _, m := range metric.Metric {
					match := true
					for k, v := range labels {
						found := false
						for _, l := range m.Label {
							if *l.Name == k && *l.Value == v {
								found = true
								break
							}
						}
						if !found {
							match = false
							break
						}
					}
					if match {
						return m.Gauge.GetValue()
					}
				}
			}
		}
		return 0
	}

	// Router metrics
	assert.Equal(t, 3.0, getMetricValue("frr_ospf_router_links_total", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "1.1.1.1",
	}))
	assert.Equal(t, 2.0, getMetricValue("frr_ospf_router_links_total", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "2.2.2.2",
	}))

	// Network metrics
	assert.Equal(t, 2.0, getMetricValue("frr_ospf_network_attached_routers_total", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "192.168.1.1",
	}))

	// Summary metrics
	assert.Equal(t, 10.0, getMetricValue("frr_ospf_summary_metric", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "10.0.0.0",
	}))

	// ASBR summary metrics
	assert.Equal(t, 20.0, getMetricValue("frr_ospf_asbr_summary_metric", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "3.3.3.3",
	}))

	// External metrics
	assert.Equal(t, 30.0, getMetricValue("frr_ospf_external_metric", map[string]string{
		"link_state_id": "4.4.4.4",
		"metric_type":   "E2",
	}))

	// NSSA external metrics
	assert.Equal(t, 40.0, getMetricValue("frr_ospf_nssa_external_metric", map[string]string{
		"area_id":       "0.0.0.1",
		"link_state_id": "5.5.5.5",
		"metric_type":   "E1",
	}))

	// Database metrics
	assert.Equal(t, 5.0, getMetricValue("frr_ospf_database_lsa_count", map[string]string{
		"area_id":  "0.0.0.0",
		"lsa_type": "router",
	}))
	assert.Equal(t, 3.0, getMetricValue("frr_ospf_database_lsa_count", map[string]string{
		"area_id":  "0.0.0.0",
		"lsa_type": "network",
	}))
	assert.Equal(t, 2.0, getMetricValue("frr_ospf_database_lsa_count", map[string]string{
		"area_id":  "0.0.0.0",
		"lsa_type": "summary",
	}))
	assert.Equal(t, 1.0, getMetricValue("frr_ospf_database_lsa_count", map[string]string{
		"area_id":  "0.0.0.0",
		"lsa_type": "asbr_summary",
	}))
	assert.Equal(t, 4.0, getMetricValue("frr_ospf_database_lsa_count", map[string]string{
		"area_id":  "0",
		"lsa_type": "external",
	}))

	// Duplicate metrics
	assert.Equal(t, 2.0, getMetricValue("frr_ospf_duplicate_lsa_count", map[string]string{
		"link_state_id": "6.6.6.6",
	}))

	// Neighbor metrics
	assert.Equal(t, 1.0, getMetricValue("frr_ospf_neighbor_state", map[string]string{
		"neighbor_id": "7.7.7.7",
		"interface":   "eth0",
	}))
	assert.Equal(t, 0.5, getMetricValue("frr_ospf_neighbor_state", map[string]string{
		"neighbor_id": "8.8.8.8",
		"interface":   "eth0",
	}))
	assert.Equal(t, 60.0, getMetricValue("frr_ospf_neighbor_uptime_seconds", map[string]string{
		"neighbor_id": "7.7.7.7",
		"interface":   "eth0",
	}))
	assert.Equal(t, 30.0, getMetricValue("frr_ospf_neighbor_uptime_seconds", map[string]string{
		"neighbor_id": "8.8.8.8",
		"interface":   "eth0",
	}))

	// Interface metrics
	assert.Equal(t, 1.0, getMetricValue("frr_interface_operational_status", map[string]string{
		"interface": "eth0",
		"vrf":       "default",
	}))
	assert.Equal(t, 0.0, getMetricValue("frr_interface_operational_status", map[string]string{
		"interface": "eth1",
		"vrf":       "default",
	}))
	assert.Equal(t, 1.0, getMetricValue("frr_interface_admin_status", map[string]string{
		"interface": "eth0",
		"vrf":       "default",
	}))
	assert.Equal(t, 1.0, getMetricValue("frr_interface_admin_status", map[string]string{
		"interface": "eth1",
		"vrf":       "default",
	}))

	// Route metrics
	assert.Equal(t, 100.0, getMetricValue("frr_route_metric", map[string]string{
		"prefix":   "10.0.0.0/24",
		"protocol": "ospf",
		"vrf":      "default",
	}))
	assert.Equal(t, 200.0, getMetricValue("frr_route_metric", map[string]string{
		"prefix":   "192.168.1.0/24",
		"protocol": "bgp",
		"vrf":      "default",
	}))
}

func TestMetricExporter_DisabledMetrics(t *testing.T) {
	// Setup
	registry := prometheus.NewRegistry()
	testLogger, err := logger.NewLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)

	flags := map[string]configs.ParsedFlag{
		"OSPFRouterData":       {Enabled: false},
		"OSPFNetworkData":      {Enabled: false},
		"OSPFSummaryData":      {Enabled: false},
		"OSPFAsbrSummaryData":  {Enabled: false},
		"OSPFExternalData":     {Enabled: false},
		"OSPFNssaExternalData": {Enabled: false},
		"OSPFDatabase":         {Enabled: false},
		"OSPFDuplicates":       {Enabled: false},
		"OSPFNeighbors":        {Enabled: false},
		"InterfaceList":        {Enabled: false},
		"RouteList":            {Enabled: false},
	}

	// Create test data
	data := &frrProto.FullFRRData{
		OspfRouterData: &frrProto.OSPFRouterData{
			RouterStates: map[string]*frrProto.OSPFRouterArea{
				"0.0.0.0": {
					LsaEntries: map[string]*frrProto.OSPFRouterLSA{
						"1.1.1.1": {NumOfLinks: 3},
					},
				},
			},
		},
	}

	// Create frrMadExporter
	frrMadExporter := exporter.NewMetricExporter(data, registry, testLogger, flags)

	// Test
	frrMadExporter.Update()

	// Verify no metrics are registered
	metrics, err := registry.Gather()
	assert.NoError(t, err)

	expectedMetrics := []string{
		"frr_ospf_router_links_total",
		"frr_ospf_network_attached_routers_total",
		"frr_ospf_summary_metric",
		"frr_ospf_asbr_summary_metric",
		"frr_ospf_external_metric",
		"frr_ospf_nssa_external_metric",
		"frr_ospf_database_lsa_count",
		"frr_ospf_duplicate_lsa_count",
		"frr_ospf_neighbor_state",
		"frr_ospf_neighbor_uptime_seconds",
		"frr_interface_operational_status",
		"frr_interface_admin_status",
		"frr_route_metric",
	}

	for _, expected := range expectedMetrics {
		found := false
		for _, metric := range metrics {
			if *metric.Name == expected {
				found = true
				break
			}
		}
		assert.False(t, found, "metric %s should not be registered", expected)
	}
}

func TestMetricExporter_OutageSimulation(t *testing.T) {
	// Setup
	testLogger, err := logger.NewLogger("test", "/tmp/frrMadExporter.log")
	assert.NoError(t, err)

	// We'll focus on the three requested data types
	flags := map[string]configs.ParsedFlag{
		"OSPFRouterData":  {Enabled: true},
		"OSPFNetworkData": {Enabled: true},
		"OSPFSummaryData": {Enabled: true},
	}

	// Helper function to find metric value
	getMetricValue := func(metrics []*dto.MetricFamily, name string, labels map[string]string) float64 {
		for _, metric := range metrics {
			if *metric.Name == name {
				for _, m := range metric.Metric {
					match := true
					for k, v := range labels {
						found := false
						for _, l := range m.Label {
							if *l.Name == k && *l.Value == v {
								found = true
								break
							}
						}
						if !found {
							match = false
							break
						}
					}
					if match {
						return m.Gauge.GetValue()
					}
				}
			}
		}
		return 0
	}

	// Helper function to check if a metric exists
	metricExists := func(metrics []*dto.MetricFamily, name string, labels map[string]string) bool {
		for _, metric := range metrics {
			if *metric.Name == name {
				if len(labels) == 0 {
					return true
				}
				for _, m := range metric.Metric {
					match := true
					for k, v := range labels {
						found := false
						for _, l := range m.Label {
							if *l.Name == k && *l.Value == v {
								found = true
								break
							}
						}
						if !found {
							match = false
							break
						}
					}
					if match {
						return true
					}
				}
			}
		}
		return false
	}

	// Helper function to count metrics in a metric family
	countMetrics := func(metrics []*dto.MetricFamily, name string) int {
		for _, metric := range metrics {
			if *metric.Name == name {
				return len(metric.Metric)
			}
		}
		return 0
	}

	// PART 1: Test with data
	registry := prometheus.NewRegistry()

	attachedRouters := make(map[string]*frrProto.AttachedRouter)
	attachedRouters["1.1.1.1"] = &frrProto.AttachedRouter{AttachedRouterId: "1.1.1.1"}
	attachedRouters["2.2.2.2"] = &frrProto.AttachedRouter{AttachedRouterId: "2.2.2.2"}

	data := &frrProto.FullFRRData{
		OspfRouterData: &frrProto.OSPFRouterData{
			RouterStates: map[string]*frrProto.OSPFRouterArea{
				"0.0.0.0": {
					LsaEntries: map[string]*frrProto.OSPFRouterLSA{
						"1.1.1.1": {NumOfLinks: 3},
						"2.2.2.2": {NumOfLinks: 2},
					},
				},
			},
		},
		OspfNetworkData: &frrProto.OSPFNetworkData{
			NetStates: map[string]*frrProto.NetAreaState{
				"0.0.0.0": {
					LsaEntries: map[string]*frrProto.NetworkLSA{
						"192.168.1.1": {AttachedRouters: attachedRouters},
					},
				},
			},
		},
		OspfSummaryData: &frrProto.OSPFSummaryData{
			SummaryStates: map[string]*frrProto.SummaryAreaState{
				"0.0.0.0": {
					LsaEntries: map[string]*frrProto.SummaryLSA{
						"10.0.0.0": {Tos0Metric: 10},
						"10.0.1.0": {Tos0Metric: 20},
					},
				},
			},
		},
	}

	frrMadExporter := exporter.NewMetricExporter(data, registry, testLogger, flags)

	frrMadExporter.Update()

	metrics, err := registry.Gather()
	assert.NoError(t, err)

	assert.True(t, metricExists(metrics, "frr_ospf_router_links_total", map[string]string{}))
	assert.Equal(t, 2, countMetrics(metrics, "frr_ospf_router_links_total"), "Should have exactly 2 router metrics")
	assert.Equal(t, 3.0, getMetricValue(metrics, "frr_ospf_router_links_total", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "1.1.1.1",
	}))
	assert.Equal(t, 2.0, getMetricValue(metrics, "frr_ospf_router_links_total", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "2.2.2.2",
	}))

	assert.True(t, metricExists(metrics, "frr_ospf_network_attached_routers_total", map[string]string{}))
	assert.Equal(t, 1, countMetrics(metrics, "frr_ospf_network_attached_routers_total"), "Should have exactly 1 network metric")
	assert.Equal(t, 2.0, getMetricValue(metrics, "frr_ospf_network_attached_routers_total", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "192.168.1.1",
	}))

	assert.True(t, metricExists(metrics, "frr_ospf_summary_metric", map[string]string{}))
	assert.Equal(t, 2, countMetrics(metrics, "frr_ospf_summary_metric"), "Should have exactly 2 summary metrics")
	assert.Equal(t, 10.0, getMetricValue(metrics, "frr_ospf_summary_metric", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "10.0.0.0",
	}))
	assert.Equal(t, 20.0, getMetricValue(metrics, "frr_ospf_summary_metric", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "10.0.1.0",
	}))

	// PART 2: Test with empty data
	registry = prometheus.NewRegistry()
	emptyData := &frrProto.FullFRRData{
		OspfRouterData:  &frrProto.OSPFRouterData{},
		OspfNetworkData: &frrProto.OSPFNetworkData{},
		OspfSummaryData: &frrProto.OSPFSummaryData{},
	}
	frrMadExporter = exporter.NewMetricExporter(emptyData, registry, testLogger, flags)

	frrMadExporter.Update()

	metrics, err = registry.Gather()
	assert.NoError(t, err)

	assert.False(t, metricExists(metrics, "frr_ospf_router_links_total", map[string]string{}),
		"Router metrics should not exist when there's no data")

	assert.False(t, metricExists(metrics, "frr_ospf_network_attached_routers_total", map[string]string{}),
		"Network metrics should not exist when there's no data")

	assert.False(t, metricExists(metrics, "frr_ospf_summary_metric", map[string]string{}),
		"Summary metrics should not exist when there's no data")

	// PART 3: Test with partial data (simulating partial outage)
	registry = prometheus.NewRegistry()
	partialData := &frrProto.FullFRRData{
		// Only include router data to simulate partial availability
		OspfRouterData: &frrProto.OSPFRouterData{
			RouterStates: map[string]*frrProto.OSPFRouterArea{
				"0.0.0.0": {
					LsaEntries: map[string]*frrProto.OSPFRouterLSA{
						"1.1.1.1": {NumOfLinks: 3},
					},
				},
			},
		},
		// Empty network data
		OspfNetworkData: &frrProto.OSPFNetworkData{},
		// Empty summary data
		OspfSummaryData: &frrProto.OSPFSummaryData{},
	}

	frrMadExporter = exporter.NewMetricExporter(partialData, registry, testLogger, flags)
	frrMadExporter.Update()

	metrics, err = registry.Gather()
	assert.NoError(t, err)

	assert.True(t, metricExists(metrics, "frr_ospf_router_links_total", map[string]string{}),
		"Router metrics should exist with partial data")
	assert.Equal(t, 1, countMetrics(metrics, "frr_ospf_router_links_total"),
		"Should have exactly 1 router metric with partial data")
	assert.Equal(t, 3.0, getMetricValue(metrics, "frr_ospf_router_links_total", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "1.1.1.1",
	}))

	assert.False(t, metricExists(metrics, "frr_ospf_network_attached_routers_total", map[string]string{}),
		"Network metrics should not exist with partial data")

	assert.False(t, metricExists(metrics, "frr_ospf_summary_metric", map[string]string{}),
		"Summary metrics should not exist with partial data")

	// PART 4: Test with new data after outage
	registry = prometheus.NewRegistry()

	newAttachedRouters := make(map[string]*frrProto.AttachedRouter)
	newAttachedRouters["3.3.3.3"] = &frrProto.AttachedRouter{AttachedRouterId: "3.3.3.3"}
	newAttachedRouters["4.4.4.4"] = &frrProto.AttachedRouter{AttachedRouterId: "4.4.4.4"}
	newAttachedRouters["5.5.5.5"] = &frrProto.AttachedRouter{AttachedRouterId: "5.5.5.5"}

	newData := &frrProto.FullFRRData{
		OspfRouterData: &frrProto.OSPFRouterData{
			RouterStates: map[string]*frrProto.OSPFRouterArea{
				"0.0.0.0": {
					LsaEntries: map[string]*frrProto.OSPFRouterLSA{
						"3.3.3.3": {NumOfLinks: 5},
					},
				},
				"0.0.0.1": {
					LsaEntries: map[string]*frrProto.OSPFRouterLSA{
						"4.4.4.4": {NumOfLinks: 1},
					},
				},
			},
		},
		OspfNetworkData: &frrProto.OSPFNetworkData{
			NetStates: map[string]*frrProto.NetAreaState{
				"0.0.0.0": {
					LsaEntries: map[string]*frrProto.NetworkLSA{
						"172.16.1.1": {AttachedRouters: newAttachedRouters},
					},
				},
			},
		},
		OspfSummaryData: &frrProto.OSPFSummaryData{
			SummaryStates: map[string]*frrProto.SummaryAreaState{
				"0.0.0.0": {
					LsaEntries: map[string]*frrProto.SummaryLSA{
						"192.168.0.0": {Tos0Metric: 30},
					},
				},
				"0.0.0.1": {
					LsaEntries: map[string]*frrProto.SummaryLSA{
						"172.16.0.0": {Tos0Metric: 15},
					},
				},
			},
		},
	}

	frrMadExporter = exporter.NewMetricExporter(newData, registry, testLogger, flags)

	frrMadExporter.Update()

	metrics, err = registry.Gather()
	assert.NoError(t, err)

	assert.True(t, metricExists(metrics, "frr_ospf_router_links_total", map[string]string{}))
	assert.Equal(t, 2, countMetrics(metrics, "frr_ospf_router_links_total"), "Should have exactly 2 router metrics")
	assert.Equal(t, 5.0, getMetricValue(metrics, "frr_ospf_router_links_total", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "3.3.3.3",
	}))
	assert.Equal(t, 1.0, getMetricValue(metrics, "frr_ospf_router_links_total", map[string]string{
		"area_id":       "0.0.0.1",
		"link_state_id": "4.4.4.4",
	}))

	assert.True(t, metricExists(metrics, "frr_ospf_network_attached_routers_total", map[string]string{}))
	assert.Equal(t, 1, countMetrics(metrics, "frr_ospf_network_attached_routers_total"), "Should have exactly 1 network metric")
	assert.Equal(t, 3.0, getMetricValue(metrics, "frr_ospf_network_attached_routers_total", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "172.16.1.1",
	}))

	assert.True(t, metricExists(metrics, "frr_ospf_summary_metric", map[string]string{}))
	assert.Equal(t, 2, countMetrics(metrics, "frr_ospf_summary_metric"), "Should have exactly 2 summary metrics")
	assert.Equal(t, 30.0, getMetricValue(metrics, "frr_ospf_summary_metric", map[string]string{
		"area_id":       "0.0.0.0",
		"link_state_id": "192.168.0.0",
	}))
	assert.Equal(t, 15.0, getMetricValue(metrics, "frr_ospf_summary_metric", map[string]string{
		"area_id":       "0.0.0.1",
		"link_state_id": "172.16.0.0",
	}))
}
