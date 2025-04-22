package analyzer

import (
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/exporter"
)

// analyze the different ospf anomalies
// call ospf functions
func (c *Analyzer) AnomalyAnalysis() {
}

func (a *Analyzer) Foobar() string {
	return "mighty analyzer"
}

func (a *Analyzer) GenerateTestAlerts() {
	// Generate sample unadvertised route alert
	err := exporter.ReportOSPFUnadvertisedRoute(
		"test-router-1",  // Router ID
		"0.0.0.1",        // OSPF Area
		"192.168.1.0/24", // Missing prefix
		true,             // Critical
		map[string]string{
			"expected_source": "config",
			"last_seen":       time.Now().Format(time.RFC3339),
		},
	)
	if err != nil {
		a.Logger.Error("Failed to generate test alert: " + err.Error())
	}

	// Generate sample overadvertised route alert
	err = exporter.ReportOSPFOveradvertisedRoute(
		"test-router-2",
		"0.0.0.2",
		"10.100.0.0/16",
		false, // Warning
		map[string]string{
			"detected_in": "LSDB scan",
		},
	)
	if err != nil {
		a.Logger.Error("Failed to generate test alert: " + err.Error())
	}

	// Generate sample duplicated route alert
	err = exporter.ReportOSPFDuplicatedRoute(
		"test-router-3",
		"0.0.0.3",
		"172.16.10.0/24",
		"Different metrics from R1 and R2",
		true,
		map[string]string{
			"conflicting_routers": "R1, R2",
		},
	)
	if err != nil {
		a.Logger.Error("Failed to generate test alert: " + err.Error())
	}

	a.Logger.Info("Generated 3 test OSPF anomaly alerts")
}

func (a *Analyzer) CleanTestAlerts() {
	exporter.ResolveOSPFRouteAnomaly("test-router-1", "unadvertised", "192.168.1.0/24")
	exporter.ResolveOSPFRouteAnomaly("test-router-2", "overadvertised", "10.100.0.0/16")
	exporter.ResolveOSPFRouteAnomaly("test-router-3", "duplicated", "172.16.10.0/24")
	a.Logger.Info("Cleared test OSPF anomaly alerts")
}
