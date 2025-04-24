package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/analyzer"
	socket "github.com/ba2025-ysmprc/frr-mad/src/backend/internal/comms/socket"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
)

type MadService struct {
	//AnalysisService  ServiceActivator
	//ExportService    ServiceActivator
	//AggregateService ServiceActivator
	Analyzer   *analyzer.Analyzer
	Aggregator *aggregator.Collector
	Exporter   string
}

type ServiceActivator struct {
	Name    string
	Service ServiceSelector
	Active  bool
}

type ServiceSelector struct {
	Analyzer   *analyzer.Analyzer
	Aggregator *aggregator.Collector
	Exporter   string
}

func main() {

	// possible start arguments
	// serviceArgs := []string{"--aggregator", "--analyzer", "--exporter"}

	// load config
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	defaultConfig := config.Default
	socketConfig := config.Socket
	aggregatorConfig := config.Aggregator
	analyzerConfig := config.Analyzer
	exporterConfig := config.Exporter

	// set debug level
	var debugLevel int
	switch defaultConfig.DebugLevel {
	case "debug":
		debugLevel = 2
	case "error":
		debugLevel = 1
	default:
		debugLevel = 0
	}

	// start logger instances
	applicationLogger := createNewLogger("frr_mad", "/tmp/frr_mad.log")
	applicationLogger.SetDebugLevel(debugLevel)

	// poll interval
	pollInterval := time.Duration(aggregatorConfig.PollInterval) * time.Second

	// service manager
	var madService MadService

	var serviceList []string

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--analyzer":
			serviceList = append(serviceList, "analyzer")
		case "--exporter":
			serviceList = append(serviceList, "exporter")
		}
	}

	if len(serviceList) == 0 {
		serviceList = append(serviceList, "analyzer")
	}

	for _, service := range serviceList {
		if service == "analyzer" {
			aggregatorLogger := createNewLogger("aggregator", "/tmp/aggregator.log")
			aggregatorLogger.SetDebugLevel(debugLevel)
			madService.Aggregator = startAggregator(aggregatorConfig, aggregatorLogger, pollInterval)

			analyzerLogger := createNewLogger("analyzer", "/tmp/analyzer.log")
			analyzerLogger.SetDebugLevel(debugLevel)
			madService.Analyzer = startAnalyzer(analyzerConfig, analyzerLogger, pollInterval, madService.Aggregator)
		}
		if service == "aggregator" {
		}
		if service == "exporter" {
			exporterLogger := createNewLogger("exporter", "/tmp/exporter.log")
			exporterLogger.SetDebugLevel(debugLevel)
			madService.Exporter = startExporter(exporterConfig, exporterLogger, pollInterval)
			fmt.Println(exporterConfig)
		}
	}

	// start socket
	sockServer := socket.NewSocket(socketConfig, madService.Aggregator.FullFrrData, madService.Analyzer, applicationLogger)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := sockServer.Start(); err != nil {
			applicationLogger.Error(fmt.Sprintf("Error starting socket server: %s\n", err))
			//fmt.Printf("Error starting socket server: %s\n", err)
			os.Exit(1)
		}
	}()

	<-sigChan
	applicationLogger.Info("Shutting down...")

}

func startAggregator(config configs.AggregatorConfig, logging *logger.Logger, pollInterval time.Duration) *aggregator.Collector {
	collector := aggregator.InitAggregator(config, logging)
	aggregator.StartAggregator(collector, pollInterval)

	return collector
}

func startAnalyzer(config interface{}, logging *logger.Logger, pollInterval time.Duration, aggregatorService *aggregator.Collector) *analyzer.Analyzer {
	detection := analyzer.InitAnalyzer(config, aggregatorService.FullFrrData, logging)
	analyzer.StartAnalyzer(detection, pollInterval)
	detection.Foobar()
	return detection
}

func startExporter(config configs.ExporterConfig, logging *logger.Logger, pollInterval time.Duration) string {
	// exporter := exporter.InitExporter(config, logging)
	// exporter.StartExporter(exporter, pollInterval)
	// return exporter
	return "foobar"
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue

}

func strToInt(value string) int {
	retValue, err := strconv.Atoi(value)
	if err != nil {
		// TODO: do proper error handling and get a solution in case it doesn't work
		fmt.Println("Error turning string to int")
	}

	return retValue
}

func createNewLogger(name, filePath string) *logger.Logger {
	logger, err := logger.NewLogger(name, filePath)
	if err != nil {
		log.Fatal(err)
	}
	return logger
}

// create different files and folders
