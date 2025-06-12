package main

import (
	"fmt"
	"os"
	"time"

	"github.com/frr-mad/frr-mad/src/backend/internal/aggregator"
	"github.com/frr-mad/frr-mad/src/backend/internal/analyzer"
	"github.com/frr-mad/frr-mad/src/backend/internal/configs"
	"github.com/frr-mad/frr-mad/src/backend/internal/exporter"
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
)

func startAggregator(config configs.AggregatorConfig, logging *logger.Logger, pollInterval time.Duration) *aggregator.Collector {
	collector := aggregator.InitAggregator(config, logging)
	aggregator.StartAggregator(collector, pollInterval)
	logging.WithAttrs(map[string]interface{}{
		"poll_interval": pollInterval.String(),
		"config":        fmt.Sprintf("%+v", config),
	}).Info("Aggregator service started successfully")
	return collector
}

func startAnalyzer(logging *logger.Logger, anomalyLogger *logger.Logger, pollInterval time.Duration, aggregatorService *aggregator.Collector) *analyzer.Analyzer {
	detection := analyzer.InitAnalyzer(aggregatorService.FullFrrData, logging, anomalyLogger)
	analyzer.StartAnalyzer(detection, pollInterval)
	logging.WithAttrs(map[string]interface{}{
		"poll_interval": pollInterval.String(),
	}).Info("Analyzer service started successfully")
	return detection
}

func startExporter(config configs.ExporterConfig, logging *logger.Logger, pollInterval time.Duration, frrData *frrProto.FullFRRData, anomalyResult *frrProto.AnomalyAnalysis) *exporter.Exporter {
	metricsExporter := exporter.NewExporter(config, logging, pollInterval, frrData, anomalyResult)

	metricsExporter.Start()
	logging.WithAttrs(map[string]interface{}{
		"poll_interval": pollInterval.String(),
		"config":        fmt.Sprintf("%+v", config),
	}).Info("Exporter service started successfully")
	return metricsExporter
}

// Helper Functions
func createDirectories(config *configs.Config) {
	paths := []string{
		config.Default.TempFiles,
		config.Default.LogPath,
		config.Socket.UnixSocketLocation,
	}

	for _, path := range paths {
		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Printf("Error creating directory %s: %v\n", path, err)
			os.Exit(1)
		}
	}
}
