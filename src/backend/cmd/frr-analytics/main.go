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
	config := configs.LoadConfig()

	socketConfig := config["socket"]
	aggregatorConfig := config["aggregator"]
	analyzerConfig := config["analyzer"]
	exporterConfig := config["exporter"]

	// start logger instances
	applicationLogger := createNewLogger("frr_mad", "/tmp/frr_mad.log")

	// poll interval
	pollInterval := time.Duration(strToInt(aggregatorConfig["PollInterval"])) * time.Second

	// service manager
	var madService MadService
	//madService := MadService{
	//AggregateService: ServiceActivator{
	//	Name:   "aggregator",
	//	Active: false,
	//},
	//AnalysisService: ServiceActivator{
	//	Name:   "analyzer",
	//	Active: false,
	//},
	//ExportService: ServiceActivator{
	//	Name:   "exporter",
	//	Active: false,
	//},
	//}

	var serviceList []string

	for _, arg := range os.Args[1:] {
		switch arg {
		case "--aggregator":
			//madService.AggregateService.Logger = createNewLogger("aggregator", "/tmp/aggregator.log")
			//madService.AggregateService.Service.Aggregator = startAggregator(aggregatorConfig, aggregatorLogger, pollInterval)

			// simplified
			//aggregatorLogger := createNewLogger("aggregator", "/tmp/aggregator.log")
			//madService.Aggregator = startAggregator(aggregatorConfig, aggregatorLogger, pollInterval)
			serviceList = append(serviceList, "aggregator")
		case "--analyzer":
			//madService.AnalysisService.Logger = createNewLogger("analyzer", "/tmp/analyzer.log")
			//madService.AnalysisService.Service.Analyzer = startAnalyzer(analyzerConfig, analyzerLogger, pollInterval)

			// simplified
			//analyzerLogger := createNewLogger("analyzer", "/tmp/analyzer.log")
			//madService.Analyzer = startAnalyzer(analyzerConfig, analyzerLogger, pollInterval)
			serviceList = append(serviceList, "analyzer")
		case "--exporter":
			//madService.ExportService.Logger = createNewLogger("exporter", "/tmp/exporter.log")
			//madService.ExportService.Service.Exporter = startExporter(analyzerConfig, exporterLogger, pollInterval)

			// simplified
			//exporterLogger := createNewLogger("exporter", "/tmp/exporter.log")
			//madService.Exporter = startExporter(analyzerConfig, exporterLogger, pollInterval)
			//fmt.Println(exporterConfig)
			serviceList = append(serviceList, "exporter")
		}
	}

	if len(serviceList) == 0 {
		serviceList = append(serviceList, "analyzer")
		serviceList = append(serviceList, "aggregator")
	}

	for _, service := range serviceList {
		if service == "analyzer" {
			analyzerLogger := createNewLogger("analyzer", "/tmp/analyzer.log")
			madService.Analyzer = startAnalyzer(analyzerConfig, analyzerLogger, pollInterval)
		}
		if service == "aggregator" {
			aggregatorLogger := createNewLogger("aggregator", "/tmp/aggregator.log")
			madService.Aggregator = startAggregator(aggregatorConfig, aggregatorLogger, pollInterval)
		}
		if service == "exporter" {
			exporterLogger := createNewLogger("exporter", "/tmp/exporter.log")
			madService.Exporter = startExporter(analyzerConfig, exporterLogger, pollInterval)
			fmt.Println(exporterConfig)
		}
	}

	// start collector
	// aggregatorLogger := createNewLogger("aggregator", "/tmp/aggregator.log")
	// collector := startAggregator(aggregatorConfig, aggregatorLogger, pollInterval)

	// Start analyzer
	// analyzerLogger := createNewLogger("analyzer", "/tmp/analyzer.log")
	// anomalyDetection := startAnalyzer(analyzerConfig, analyzerLogger, pollInterval)

	// start socket
	sockServer := socket.NewSocket(socketConfig, madService.Aggregator, madService.Analyzer, applicationLogger)

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

func startAnalyzer(config map[string]string, logging *logger.Logger, pollInterval time.Duration) *analyzer.Analyzer {
	detection := analyzer.InitAnalyzer(config)
	analyzer.StartAnalyzer(detection, pollInterval)
	return detection
}

func startExporter(config map[string]string, logging *logger.Logger, pollInterval time.Duration) string {
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
