package exporter

import (
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Exporter struct {
	interval        time.Duration
	anomalyExporter *AnomalyExporter
	metricExporter  *MetricExporter
	server          *http.Server
	stopChan        chan struct{}
	logger          *logger.Logger
}

func NewExporter(
	config map[string]string,
	logger *logger.Logger,
	pollInterval time.Duration,
	frrData *frrProto.FullFRRData,
	anomalies *frrProto.Anomalies,
) (*Exporter, error) {
	// Parse port
	port := 9091
	if portStr, exists := config["Port"]; exists {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			return nil, fmt.Errorf("invalid port in config: %v", err)
		}
	}

	registry := prometheus.NewRegistry()
	flags, err := configs.GetFlagConfigs(config)
	if err != nil {
		return nil, fmt.Errorf("error loading config flags: %v", err)
	}

	// Create exporter
	e := &Exporter{
		interval:        pollInterval,
		anomalyExporter: NewAnomalyExporter(anomalies, registry, logger),
		metricExporter:  NewMetricExporter(frrData, registry, logger, flags),
		stopChan:        make(chan struct{}),
		logger:          logger,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: promhttp.HandlerFor(registry, promhttp.HandlerOpts{}),
		},
	}

	return e, nil
}

func (e *Exporter) Start() {
	go func() {
		if err := e.server.ListenAndServe(); err != nil {
			e.logger.Error(fmt.Sprintf("Metrics server failed: %v", err))
		}
	}()

	go e.runExportLoop()
	e.logger.Info(fmt.Sprintf("Exporter started on port %s", e.server.Addr))
}

func (e *Exporter) Stop() {
	close(e.stopChan)
	_ = e.server.Close()
}

func (e *Exporter) runExportLoop() {
	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			e.exportData()
		case <-e.stopChan:
			return
		}
	}
}

func (e *Exporter) exportData() {
	// Always export anomalies
	e.anomalyExporter.Update()

	// Export metrics if exporter exists
	if e.metricExporter != nil {
		e.metricExporter.Update()
	}
}
