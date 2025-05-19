package exporter

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
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

type ParsedFlag struct {
	Name        string
	Description string
	Enabled     bool
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
		var err error
		port = config.Port
		if err != nil {
			logger.Error(fmt.Sprintf("invalid port in config: %v", err))
		}
	}

	logger.WithAttrs(map[string]interface{}{
		"port":          port,
		"poll_interval": pollInterval.String(),
	}).Info("Initializing exporter")

	registry := prometheus.NewRegistry()
	flags := getFlagConfigs(config)

	enabledFlags := []string{}
	for name, flag := range flags {
		if flag.Enabled {
			enabledFlags = append(enabledFlags, name)
		}
	}
	logger.WithAttrs(map[string]interface{}{
		"enabled_flags": enabledFlags,
	}).Debug("Exporter flag configuration")

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
		metricExporter:  NewMetricExporter(frrData, registry, logger, flags),
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
	e.logger.Debug("Starting export loop")
	defer e.logger.Debug("Export loop stopped")

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

// Use reflection to iterate over struct fields
func getFlagConfigs(config configs.ExporterConfig) map[string]*ParsedFlag {
	result := make(map[string]*ParsedFlag)

	val := reflect.ValueOf(config)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldName := typ.Field(i).Name

		if fieldName == "Port" {
			continue
		}

		if field.Kind() == reflect.String && field.String() != "" {
			result[fieldName] = parseFlagTuple(field.String())
		}
	}

	return result
}

func parseFlagTuple(tuple string) *ParsedFlag {
	tuple = strings.TrimSpace(tuple)
	if !strings.HasPrefix(tuple, "(") || !strings.HasSuffix(tuple, ")") {
		fmt.Errorf("invalid flag tuple format - must be enclosed in parentheses")
		return nil
	}

	tuple = tuple[1 : len(tuple)-1]

	parts := splitTupleComponents(tuple)
	if len(parts) != 3 {
		fmt.Errorf("flag tuple must have exactly 3 components")
		return nil
	}

	name := strings.TrimSpace(parts[0])
	description := strings.TrimSpace(parts[1])
	enabledStr := strings.TrimSpace(parts[2])

	if strings.HasPrefix(description, `"`) && strings.HasSuffix(description, `"`) {
		description = description[1 : len(description)-1]
	}

	enabled, err := strconv.ParseBool(enabledStr)
	if err != nil {
		fmt.Errorf("invalid boolean value in flag tuple: %v", err)
		return nil
	}

	return &ParsedFlag{
		Name:        name,
		Description: description,
		Enabled:     enabled,
	} //, nil
}

func splitTupleComponents(tuple string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false

	for _, r := range tuple {
		switch {
		case r == ',' && !inQuotes:
			parts = append(parts, current.String())
			current.Reset()
		case r == '"':
			inQuotes = !inQuotes
			current.WriteRune(r)
		default:
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		parts = append(parts, current.String())
	}

	return parts
}
