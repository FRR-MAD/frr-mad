package exporter

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

/*--------------------------------------------------------------------------
- How to use:

err := exporter.ReportOSPFUnadvertisedRoute(
    "router-42",                      // Router ID
    "0.0.0.1",                        // OSPF Area ID
    "10.39.0.0/17",                   // Expected prefix that's missing
    true,                             // Critical issue
    map[string]string{                // Additional info
        "router_type": "ABR",
        "last_seen": "2025-04-09T15:30:45Z",
    },
)
--------------------------------------------------------------------------*/

type AlertType string

const (
	// MVP
	AlertTypeOSPFUnadvertisedRoute      AlertType = "ospf_unadvertised_route"
	AlertTypeOSPFOveradvertisedRoute    AlertType = "ospf_overadvertised_route"
	AlertTypeOSPFWronglyAdvertisedRoute AlertType = "ospf_wrongly_advertised_route"
	AlertTypeOSPFDuplicatedRoute        AlertType = "ospf_duplicated_route"

	// General Stuff
	AlertTypeOSPFNeighborDown      AlertType = "ospf_neighbor_down"
	AlertTypeOSPFAreaMisconfigured AlertType = "ospf_area_misconfigured"
)

type AlertSeverity string

const (
	SeverityCritical AlertSeverity = "critical"
	SeverityWarning  AlertSeverity = "warning"
	SeverityInfo     AlertSeverity = "info"
)

type RouterAlert struct {
	RouterID       string
	Type           AlertType
	Severity       AlertSeverity
	Message        string
	Timestamp      time.Time
	Labels         map[string]string
	ExpectedValue  string // For comparison
	ActualValue    string // What was actually found
	AffectedPrefix string // The network prefix affected
	AreaID         string // OSPF area identifier
}

type PrometheusAlerter struct {
	registry        *prometheus.Registry
	metrics         map[AlertType]*prometheus.GaugeVec
	alerts          map[string]RouterAlert
	mutex           sync.RWMutex
	server          *http.Server
	started         bool
	anomalyCounters map[AlertType]*prometheus.GaugeVec
	logger          *logger.Logger
}

var (
	singleton *PrometheusAlerter
	once      sync.Once
)

func GetPrometheusAlerter() *PrometheusAlerter {
	once.Do(func() {
		singleton = &PrometheusAlerter{
			registry:        prometheus.NewRegistry(),
			metrics:         make(map[AlertType]*prometheus.GaugeVec),
			alerts:          make(map[string]RouterAlert),
			anomalyCounters: make(map[AlertType]*prometheus.GaugeVec),
		}
		singleton.initMetrics()
	})
	return singleton
}

func (p *PrometheusAlerter) initMetrics() {
	ospfAnomalyTypes := []AlertType{
		AlertTypeOSPFUnadvertisedRoute,
		AlertTypeOSPFOveradvertisedRoute,
		AlertTypeOSPFWronglyAdvertisedRoute,
		AlertTypeOSPFDuplicatedRoute,
		AlertTypeOSPFNeighborDown,
		AlertTypeOSPFAreaMisconfigured,
	}

	for _, alertType := range ospfAnomalyTypes {
		metricName := fmt.Sprintf("router_%s", alertType)
		gauge := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: metricName,
				Help: fmt.Sprintf("Alert for %s", alertType),
			},
			[]string{"router_id", "area_id", "prefix", "severity", "message", "alert_type"},
		)
		p.registry.MustRegister(gauge)
		p.metrics[alertType] = gauge

		counterName := fmt.Sprintf("router_%s_total", alertType)
		counter := prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Name: counterName,
				Help: fmt.Sprintf("Total number of active %s anomalies", alertType),
			},
			[]string{"router_id", "area_id", "severity"},
		)
		p.registry.MustRegister(counter)
		p.anomalyCounters[alertType] = counter
	}
}

func (p *PrometheusAlerter) Start(port int) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	if p.started {
		return nil
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.HandlerFor(p.registry, promhttp.HandlerOpts{}))

	p.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: mux,
	}

	go func() {
		p.logger.Info(fmt.Sprintf("Starting Prometheus metrics server on %s", p.server.Addr))
		if err := p.server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			p.logger.Error(fmt.Sprintf("Error starting Prometheus metrics server: %v", err))
		}
	}()

	p.started = true
	return nil
}

func (p *PrometheusAlerter) CreateAlert(alert RouterAlert) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	key := fmt.Sprintf("%s_%s_%s", alert.RouterID, alert.Type, alert.AffectedPrefix)
	p.alerts[key] = alert

	gauge, exists := p.metrics[alert.Type]
	if !exists {
		return fmt.Errorf("unknown alert type: %s", alert.Type)
	}

	gauge.With(prometheus.Labels{
		"router_id":  alert.RouterID,
		"area_id":    alert.AreaID,
		"prefix":     alert.AffectedPrefix,
		"severity":   string(alert.Severity),
		"message":    alert.Message,
		"alert_type": string(alert.Type),
	}).Set(1)

	counter, exists := p.anomalyCounters[alert.Type]
	if exists {
		counter.With(prometheus.Labels{
			"router_id": alert.RouterID,
			"area_id":   alert.AreaID,
			"severity":  string(alert.Severity),
		}).Inc()
	}

	return nil
}

func (p *PrometheusAlerter) ResolveAlert(routerID string, alertType AlertType, prefix string) error {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	key := fmt.Sprintf("%s_%s_%s", routerID, alertType, prefix)

	alert, exists := p.alerts[key]
	if !exists {
		return fmt.Errorf("no active alert found for router %s with type %s and prefix %s", routerID, alertType, prefix)
	}

	// Remove alert from map
	delete(p.alerts, key)

	gauge, exists := p.metrics[alertType]
	if !exists {
		return fmt.Errorf("unknown alert type: %s", alertType)
	}

	gauge.With(prometheus.Labels{
		"router_id":  alert.RouterID,
		"area_id":    alert.AreaID,
		"prefix":     alert.AffectedPrefix,
		"severity":   string(alert.Severity),
		"message":    alert.Message,
		"alert_type": string(alert.Type),
	}).Set(0)

	counter, exists := p.anomalyCounters[alertType]
	if exists {
		counter.With(prometheus.Labels{
			"router_id": alert.RouterID,
			"area_id":   alert.AreaID,
			"severity":  string(alert.Severity),
		}).Dec()
	}

	return nil
}

func (p *PrometheusAlerter) GetActiveAlerts() []RouterAlert {
	p.mutex.RLock()
	defer p.mutex.RUnlock()

	alerts := make([]RouterAlert, 0, len(p.alerts))
	for _, alert := range p.alerts {
		alerts = append(alerts, alert)
	}

	return alerts
}

func ReportOSPFUnadvertisedRoute(routerID string, areaID string, expectedPrefix string, isCritical bool, additionalInfo map[string]string) error {
	alerter := GetPrometheusAlerter()

	if !alerter.started {
		return fmt.Errorf("exporter server not started")
	}

	severity := SeverityWarning
	if isCritical {
		severity = SeverityCritical
	}

	message := fmt.Sprintf("Expected route %s is not being advertised in area %s", expectedPrefix, areaID)

	alert := RouterAlert{
		RouterID:       routerID,
		Type:           AlertTypeOSPFUnadvertisedRoute,
		Severity:       severity,
		Message:        message,
		Timestamp:      time.Now(),
		Labels:         additionalInfo,
		ExpectedValue:  expectedPrefix,
		ActualValue:    "not advertised",
		AffectedPrefix: expectedPrefix,
		AreaID:         areaID,
	}

	return alerter.CreateAlert(alert)
}

func ReportOSPFOveradvertisedRoute(routerID string, areaID string, unexpectedPrefix string, isCritical bool, additionalInfo map[string]string) error {
	alerter := GetPrometheusAlerter()

	if !alerter.started {
		return fmt.Errorf("exporter server not started")
	}

	severity := SeverityWarning
	if isCritical {
		severity = SeverityCritical
	}

	message := fmt.Sprintf("Route %s is being advertised in area %s but should NOT be", unexpectedPrefix, areaID)

	alert := RouterAlert{
		RouterID:       routerID,
		Type:           AlertTypeOSPFOveradvertisedRoute,
		Severity:       severity,
		Message:        message,
		Timestamp:      time.Now(),
		Labels:         additionalInfo,
		ExpectedValue:  "not advertised",
		ActualValue:    unexpectedPrefix,
		AffectedPrefix: unexpectedPrefix,
		AreaID:         areaID,
	}

	return alerter.CreateAlert(alert)
}

func ReportOSPFWronglyAdvertisedRoute(routerID string, areaID string, advertizedPrefix string, expectedPrefix string, isCritical bool, additionalInfo map[string]string) error {
	alerter := GetPrometheusAlerter()

	if !alerter.started {
		return fmt.Errorf("exporter server not started")
	}

	severity := SeverityWarning
	if isCritical {
		severity = SeverityCritical
	}

	message := fmt.Sprintf("Route %s is incorrectly advertised (should be %s) in area %s", advertizedPrefix, expectedPrefix, areaID)

	alert := RouterAlert{
		RouterID:       routerID,
		Type:           AlertTypeOSPFWronglyAdvertisedRoute,
		Severity:       severity,
		Message:        message,
		Timestamp:      time.Now(),
		Labels:         additionalInfo,
		ExpectedValue:  expectedPrefix,
		ActualValue:    advertizedPrefix,
		AffectedPrefix: advertizedPrefix,
		AreaID:         areaID,
	}

	return alerter.CreateAlert(alert)
}

func ReportOSPFDuplicatedRoute(routerID string, areaID string, duplicatedPrefix string, conflictDetails string, isCritical bool, additionalInfo map[string]string) error {
	alerter := GetPrometheusAlerter()

	if !alerter.started {
		return fmt.Errorf("exporter server not started")
	}

	severity := SeverityWarning
	if isCritical {
		severity = SeverityCritical
	}

	message := fmt.Sprintf("Route %s is duplicated with conflicting attributes in area %s", duplicatedPrefix, areaID)

	alert := RouterAlert{
		RouterID:       routerID,
		Type:           AlertTypeOSPFDuplicatedRoute,
		Severity:       severity,
		Message:        message,
		Timestamp:      time.Now(),
		Labels:         additionalInfo,
		ExpectedValue:  "unique route",
		ActualValue:    conflictDetails,
		AffectedPrefix: duplicatedPrefix,
		AreaID:         areaID,
	}

	return alerter.CreateAlert(alert)
}

func ResolveOSPFRouteAnomaly(routerID string, anomalyType string, prefix string) error {
	alerter := GetPrometheusAlerter()

	var alertType AlertType

	switch anomalyType {
	case "unadvertised":
		alertType = AlertTypeOSPFUnadvertisedRoute
	case "overadvertised":
		alertType = AlertTypeOSPFOveradvertisedRoute
	case "wrongly_advertised":
		alertType = AlertTypeOSPFWronglyAdvertisedRoute
	case "duplicated":
		alertType = AlertTypeOSPFDuplicatedRoute
	default:
		return fmt.Errorf("unknown OSPF route anomaly type: %s", anomalyType)
	}

	return alerter.ResolveAlert(routerID, alertType, prefix)
}

func GetOSPFAnomalyStats(routerID string, areaID string) (map[string]int, error) {
	alerter := GetPrometheusAlerter()

	stats := make(map[string]int)

	activeAlerts := alerter.GetActiveAlerts()
	for _, alert := range activeAlerts {
		if alert.RouterID == routerID && (areaID == "" || alert.AreaID == areaID) {
			statKey := string(alert.Type)
			stats[statKey]++
		}
	}

	return stats, nil
}
