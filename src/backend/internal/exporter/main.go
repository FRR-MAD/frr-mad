package exporter

import (
	"fmt"
	"net/http"
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
) *Exporter {
	// Default port
	port := 9091

	if config.Port > 0 {
		if config.Port > 65535 {
			logger.Error(fmt.Sprintf("invalid port in config: port %d is too high", config.Port))
		} else {
			port = config.Port
		}
	}

	logger.WithAttrs(map[string]interface{}{
		"port":          port,
		"poll_interval": pollInterval.String(),
	}).Info("Initializing exporter")

	registry := prometheus.NewRegistry()

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
		metricExporter:  NewMetricExporter(frrData, registry, logger, config),
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
			e.logger.WithAttrs(map[string]interface{}{
				"error": err.Error(),
			}).Error("Metrics server failed")
		}
	}()

	//e.runExportLoop()
	go e.runExportLoop()
	e.logger.WithAttrs(map[string]interface{}{
		"address":  e.server.Addr,
		"interval": e.interval.String(),
	}).Info("Exporter successfully started")
}

func (e *Exporter) Stop() {
	e.logger.Info("Shutting down exporter")
	close(e.stopChan)
	if err := e.server.Close(); err != nil {
		e.logger.WithAttrs(map[string]interface{}{
			"error": err.Error(),
		}).Error("Error while shutting down server")
	}
	e.logger.Info("Exporter shutdown complete")
}

func (e *Exporter) runExportLoop() {
	e.logger.Info("Starting export loop")
	defer e.logger.Info("Export loop stopped")

	ticker := time.NewTicker(e.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			start := time.Now()
			e.exportData()
			e.logger.WithAttrs(map[string]interface{}{
				"duration": time.Since(start).String(),
			}).Debug("Completed export cycle")
		case <-e.stopChan:
			return
		}
	}
}

func (e *Exporter) exportData() {
	e.logger.Debug("Starting data export")

	err := tryUpdateWithRetry("AnomalyExporter", e.anomalyExporter.Update, e.logger)
	if err != nil {
		e.logger.WithAttrs(map[string]interface{}{
			"component": "AnomalyExporter",
			"error":     err.Error(),
		}).Error("Failed to update after retry")
	}

	if e.metricExporter != nil {
		err := tryUpdateWithRetry("MetricExporter", e.metricExporter.Update, e.logger)
		if err != nil {
			e.logger.WithAttrs(map[string]interface{}{
				"component": "MetricExporter",
				"error":     err.Error(),
			}).Error("Failed to update after retry")
		}
	}

	e.logger.Debug("Data export completed")
}

func tryUpdateWithRetry(name string, updateFunc func(), logger *logger.Logger) error {
	const retryDelay = 500 * time.Millisecond

	logger.WithAttrs(map[string]interface{}{
		"component": name,
	}).Debug("Attempting update")

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
		logger.WithAttrs(map[string]interface{}{
			"component":   name,
			"error":       err.Error(),
			"retry_delay": retryDelay.String(),
		}).Warning("Update failed, will retry")

		time.Sleep(retryDelay)

		// Retry
		if retryErr := try(); retryErr != nil {
			logger.WithAttrs(map[string]interface{}{
				"component": name,
				"error":     retryErr.Error(),
			}).Debug("Retry attempt failed")
			return retryErr
		}
		logger.WithAttrs(map[string]interface{}{
			"component": name,
		}).Debug("Retry succeeded")
	}

	logger.WithAttrs(map[string]interface{}{
		"component": name,
	}).Debug("Update completed successfully")
	return nil
}
