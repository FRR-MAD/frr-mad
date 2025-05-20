package exporter

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/frr-mad/frr-mad/src/backend/configs"
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
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
	config configs.ExporterConfig,
	logger *logger.Logger,
	pollInterval time.Duration,
	frrData *frrProto.FullFRRData,
	anomalies *frrProto.AnomalyAnalysis,
	flags map[string]bool,
) *Exporter {
	port := 9091

	if config.Port > 0 {
		var err error
		port = config.Port
		if err != nil {
			logger.Error(fmt.Sprintf("invalid port in config: %v", err))
		}
	}

	registry := prometheus.NewRegistry()

	finalFlags := mergeConfigAndFlags(config, flags)

	fmt.Println(config.OSPFNeighbors)
	fmt.Println(finalFlags)

	mux := http.NewServeMux()

	mux.Handle("/metrics", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html")
		fmt.Fprintf(w, `
			<html>
				<head><title>FRR MAD Exporter</title></head>
				<body>
					<h1>FRR MAD Exporter</h1>
					<p><a href="/metrics">Metrics</a></p>
				</body>
			</html>
			`)
	})

	e := &Exporter{
		interval:        pollInterval,
		anomalyExporter: NewAnomalyExporter(anomalies, registry, logger),
		metricExporter:  NewMetricExporter(frrData, registry, logger, finalFlags),
		stopChan:        make(chan struct{}),
		logger:          logger,
		server: &http.Server{
			Addr:    fmt.Sprintf(":%d", port),
			Handler: mux,
		},
	}

	return e
}

func (e *Exporter) Start() {
	go func() {
		if err := e.server.ListenAndServe(); err != nil {
			e.logger.Error(fmt.Sprintf("Metrics server failed: %v", err))
		}
	}()

	//e.runExportLoop()
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
	err := tryUpdateWithRetry("AnomalyExporter", e.anomalyExporter.Update, e.logger)
	if err != nil {
		e.logger.Error(fmt.Sprintf("Final failure: Anomaly Exporter Update function: %v", err))
	}

	if e.metricExporter != nil {
		err := tryUpdateWithRetry("MetricExporter", e.metricExporter.Update, e.logger)
		if err != nil {
			e.logger.Error(fmt.Sprintf("Final failure: Metric Exporter Update function: %v", err))
		}
	}
}

func tryUpdateWithRetry(name string, updateFunc func(), logger *logger.Logger) error {
	const retryDelay = 500 * time.Millisecond

	try := func() (err error) {
		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("panic: %v", r)
			}
		}()
		updateFunc()
		return nil
	}

	// First try
	if err := try(); err != nil {
		logger.Warning(fmt.Sprintf("%s update failed, retrying in %s: %v", name, retryDelay, err))
		time.Sleep(retryDelay)

		// Retry
		if retryErr := try(); retryErr != nil {
			return retryErr
		}
	}

	return nil
}

func mergeConfigAndFlags(config configs.ExporterConfig, flags map[string]bool) map[string]bool {
	merged := map[string]bool{
		"OSPFRouterData":       false,
		"OSPFNetworkData":      false,
		"OSPFSummaryData":      false,
		"OSPFAsbrSummaryData":  false,
		"OSPFExternalData":     false,
		"OSPFNssaExternalData": false,
		"OSPFDatabase":         false,
		"OSPFNeighbors":        false,
		"InterfaceList":        false,
		"RouteList":            false,
	}

	getConfigBool := func(value interface{}) bool {
		switch v := value.(type) {
		case int:
			return v == 1
		case int64:
			return v == 1
		case float64:
			return v == 1
		case string:
			return v == "1" || strings.EqualFold(v, "true")
		case bool:
			return v
		default:
			return false
		}
	}

	merged["OSPFRouterData"] = getConfigBool(config.OSPFRouterData)
	merged["OSPFNetworkData"] = getConfigBool(config.OSPFNetworkData)
	merged["OSPFSummaryData"] = getConfigBool(config.OSPFSummaryData)
	merged["OSPFAsbrSummaryData"] = getConfigBool(config.OSPFAsbrSummaryData)
	merged["OSPFExternalData"] = getConfigBool(config.OSPFExternalData)
	merged["OSPFNssaExternalData"] = getConfigBool(config.OSPFNssaExternalData)
	merged["OSPFDatabase"] = getConfigBool(config.OSPFDatabase)
	merged["OSPFNeighbors"] = getConfigBool(config.OSPFNeighbors)
	merged["InterfaceList"] = getConfigBool(config.InterfaceList)
	merged["RouteList"] = getConfigBool(config.RouteList)

	for metric, flagValue := range flags {
		if _, exists := merged[metric]; exists {
			merged[metric] = merged[metric] || flagValue
		}
	}

	return merged
}
