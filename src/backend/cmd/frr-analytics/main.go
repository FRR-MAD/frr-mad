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
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/exporter"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
)

type MadService struct {
	//AnalysisService  ServiceActivator
	//ExportService    ServiceActivator
	//AggregateService ServiceActivator
	Analyzer   *analyzer.Analyzer
	Aggregator *aggregator.Collector
	Exporter   *exporter.PrometheusAlerter
}

type ServiceActivator struct {
	Name    string
	Service ServiceSelector
	Active  bool
}

type ServiceSelector struct {
	Analyzer   *analyzer.Analyzer
	Aggregator *aggregator.Collector
	Exporter   *exporter.PrometheusAlerter
}

func main() {

	// possible start arguments
	// serviceArgs := []string{"--aggregator", "--analyzer", "--exporter"}

	// load config
	config := configs.LoadConfig()

	defaultConfig := config["default"]
	socketConfig := config["socket"]
	aggregatorConfig := config["aggregator"]
	analyzerConfig := config["analyzer"]
	exporterConfig := config["exporter"]

	// set debug level
	var debugLevel int
	switch defaultConfig["DebugLevel"] {
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
	pollInterval := time.Duration(strToInt(aggregatorConfig["PollInterval"])) * time.Second

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
		serviceList = append(serviceList, "exporter")
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
			exporter, err := startExporter(exporterConfig, exporterLogger)
			if err != nil {
				exporterLogger.Error(fmt.Sprintf("Failed to start exporter: %v", err))
				// if not working close everything
				//os.Exit(1)
			} else {
				madService.Exporter = exporter
			}
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

func startAggregator(config map[string]string, logging *logger.Logger, pollInterval time.Duration) *aggregator.Collector {
	collector := aggregator.InitAggregator(config, logging)
	aggregator.StartAggregator(collector, pollInterval)

	return collector
}

func startAnalyzer(config map[string]string, logging *logger.Logger, pollInterval time.Duration, aggregatorService *aggregator.Collector) *analyzer.Analyzer {
	detection := analyzer.InitAnalyzer(config, aggregatorService, logging)
	analyzer.StartAnalyzer(detection, pollInterval)

	// test the exporter with anomalies
	//detection.GenerateTestAlerts()

	// if the anomalie resolves:
	//detection.CleanTestAlerts()

	detection.Foobar()
	return detection
}

func startExporter(config map[string]string, logging *logger.Logger) (*exporter.PrometheusAlerter, error) {
	return exporter.InitExporter(config, logging)
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
