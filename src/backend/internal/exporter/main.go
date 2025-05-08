package exporter

import (
	"fmt"
	"net/http"
	"reflect"
	"strconv"
	"strings"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/ba2025-ysmprc/frr-mad/src/logger"
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
	config          configs.Config
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
	// Parse port
	port := 9091

	if config.Port > 0 {
		var err error
		port = config.Port
		if err != nil {
			logger.Error(fmt.Sprintf("invalid port in config: %v", err))
		}
	}

	fmt.Println(port)

	registry := prometheus.NewRegistry()
	flags := getFlagConfigs(config)
	//flags, err := configs.GetFlagConfigs(config)
	// if err != nil {
	// 	return nil, fmt.Errorf("error loading config flags: %v", err)
	// }

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
	fmt.Sprintf("Exporter started on port %s", e.server.Addr)
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
	e.anomalyExporter.Update()

	if e.metricExporter != nil {
		e.metricExporter.Update()
	}
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
