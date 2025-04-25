package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/analyzer"
	socket "github.com/ba2025-ysmprc/frr-mad/src/backend/internal/comms/socket"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
)

// Service represents a running service component of the application
type Service struct {
	Name   string
	Active bool
}

// FrrMadApp represents the main application with its service components
type FrrMadApp struct {
	Analyzer   *analyzer.Analyzer
	Aggregator *aggregator.Collector
	Exporter   string
	Socket     *socket.Socket
	Logger     *logger.Logger
}

func main() {
	// Create a custom flag set
	cmdSet := flag.NewFlagSet("frr-mad", flag.ExitOnError)

	// Define help text
	cmdSet.Usage = func() {
		fmt.Println("Usage: frr-mad [command] [options]")
		fmt.Println("\nCommands:")
		fmt.Println("  start   - Start the FRR Monitoring and Analysis Daemon")
		fmt.Println("  stop    - Stop a running FRR-MAD instance")
		fmt.Println("  reload  - Reload configuration for a running FRR-MAD")
		fmt.Println("  help    - Display this help message")
		// fmt.Println("\nOptions:")
		// fmt.Println("  --analyzer   - Run with analyzer service enabled")
		// fmt.Println("  --aggregator - Run with aggregator service enabled")
		// fmt.Println("  --exporter   - Run with exporter service enabled")
		// fmt.Println("\nIf no service options are specified, analyzer is enabled by default.")
	}

	// Check if a command was provided
	if len(os.Args) < 2 {
		cmdSet.Usage()
		os.Exit(1)
	}

	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Ensure required directories exist
	createDirectories(config)

	// Process the command
	command := os.Args[1]

	switch command {
	case "start":
		startApp(config)
	case "stop":
		pidFile := fmt.Sprintf("%s/frr-mad.pid", config.Socket.UnixSocketLocation)
		stopApp(pidFile)
	case "reload":
		fmt.Println("Reloading the FRR-MAD configuration...")
		fmt.Println("Not implemented yet. Please restart the application manually.")
	case "help":
		cmdSet.Usage()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		cmdSet.Usage()
		os.Exit(1)
	}
}

func startApp(config *configs.Config) {

	// Extract configuration sections
	defaultConfig := config.Default
	socketConfig := config.Socket
	aggregatorConfig := config.Aggregator
	analyzerConfig := config.Analyzer
	exporterConfig := config.Exporter

	// Configure logging
	debugLevel := getDebugLevel(defaultConfig.DebugLevel)
	appLogger := createLogger("frr_mad", fmt.Sprintf("%v/frr_mad.log", defaultConfig.LogPath))
	appLogger.SetDebugLevel(debugLevel)
	appLogger.Info("Starting FRR Monitoring and Analysis Daemon")

	// Configure polling interval
	pollInterval := time.Duration(aggregatorConfig.PollInterval) * time.Second
	appLogger.Info(fmt.Sprintf("Setting poll interval to %v seconds", aggregatorConfig.PollInterval))

	// Determine which services to start
	services := parseServiceFlags()
	if len(services) == 0 {
		services = append(services, "analyzer") // Default to analyzer if no services specified
	}

	// Initialize application
	app := &FrrMadApp{
		Logger: appLogger,
	}

	// Start requested services

	pidFile := createPid(socketConfig.UnixSocketLocation, appLogger)
	defer os.Remove(pidFile)

	for _, service := range services {
		appLogger.Info(fmt.Sprintf("Starting %s service", service))
		switch service {
		case "analyzer":
			// Start aggregator first (analyzer depends on it)
			if app.Aggregator == nil {
				aggregatorLogger := createLogger("aggregator", fmt.Sprintf("%v/aggregator.log", defaultConfig.LogPath))
				aggregatorLogger.SetDebugLevel(debugLevel)
				app.Aggregator = startAggregator(aggregatorConfig, aggregatorLogger, pollInterval)
			}

			// Then start analyzer
			analyzerLogger := createLogger("analyzer", fmt.Sprintf("%v/analyzer.log", defaultConfig.LogPath))
			analyzerLogger.SetDebugLevel(debugLevel)
			app.Analyzer = startAnalyzer(analyzerConfig, analyzerLogger, pollInterval, app.Aggregator)

		case "aggregator":
			if app.Aggregator == nil {
				aggregatorLogger := createLogger("aggregator", fmt.Sprintf("%v/aggregator.log", defaultConfig.LogPath))
				aggregatorLogger.SetDebugLevel(debugLevel)
				app.Aggregator = startAggregator(aggregatorConfig, aggregatorLogger, pollInterval)
			}

		case "exporter":
			exporterLogger := createLogger("exporter", fmt.Sprintf("%v/exporter.log", defaultConfig.LogPath))
			exporterLogger.SetDebugLevel(debugLevel)
			app.Exporter = startExporter(exporterConfig, exporterLogger, pollInterval)
		}
	}

	// Ensure aggregator is started if needed by other services
	if app.Analyzer != nil && app.Aggregator == nil {
		aggregatorLogger := createLogger("aggregator", fmt.Sprintf("%v/aggregator.log", defaultConfig.LogPath))
		aggregatorLogger.SetDebugLevel(debugLevel)
		app.Aggregator = startAggregator(aggregatorConfig, aggregatorLogger, pollInterval)
	}

	// Start the socket server if we have any active services
	if app.Aggregator != nil && app.Analyzer != nil {
		app.Socket = socket.NewSocket(socketConfig, app.Aggregator.FullFrrData, app.Analyzer, appLogger)

		go func() {
			if err := app.Socket.Start(); err != nil {
				appLogger.Error(fmt.Sprintf("Error starting socket server: %s", err))
				os.Exit(1)
			}
		}()

		appLogger.Info(fmt.Sprintf("Socket server listening at %s/%s",
			socketConfig.UnixSocketLocation, socketConfig.UnixSocketName))
	} else {
		appLogger.Error("Cannot start socket server: required services not available")
		os.Exit(1)
	}

	// Wait for termination signal
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	appLogger.Info("Received shutdown signal")

	// Cleanup
	if app.Socket != nil {
		app.Socket.Close()
	}

	appLogger.Info("FRR-MAD shutdown complete")
}

// Parse command line flags to determine which services to start
func parseServiceFlags() []string {
	var services []string

	for _, arg := range os.Args[2:] {
		switch arg {
		case "--analyzer":
			services = append(services, "analyzer")
		case "--aggregator":
			services = append(services, "aggregator")
		case "--exporter":
			services = append(services, "exporter")
		}
	}

	return services
}

// Create a new logger instance
func createLogger(name, filePath string) *logger.Logger {
	logger, err := logger.NewLogger(name, filePath)
	if err != nil {
		log.Fatalf("Failed to create logger %s: %v", name, err)
	}
	return logger
}

// Convert debug level string to int
func getDebugLevel(level string) int {
	switch level {
	case "debug":
		return 2
	case "error":
		return 1
	default:
		return 0
	}
}

// Start the aggregator service
func startAggregator(config configs.AggregatorConfig, logging *logger.Logger, pollInterval time.Duration) *aggregator.Collector {
	collector := aggregator.InitAggregator(config, logging)
	aggregator.StartAggregator(collector, pollInterval)
	logging.Info("Aggregator service started")
	return collector
}

// Start the analyzer service
func startAnalyzer(config interface{}, logging *logger.Logger, pollInterval time.Duration, aggregatorService *aggregator.Collector) *analyzer.Analyzer {
	detection := analyzer.InitAnalyzer(config, aggregatorService.FullFrrData, logging)
	analyzer.StartAnalyzer(detection, pollInterval)
	logging.Info("Analyzer service started")
	return detection
}

// Start the exporter service
func startExporter(config configs.ExporterConfig, logging *logger.Logger, pollInterval time.Duration) string {
	// Placeholder for exporter implementation
	logging.Info("Exporter service started (placeholder)")
	return "exporter"
}

// Create required directories
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

func createPid(socketPath string, appLogger *logger.Logger) string {
	pid := os.Getpid()
	pidFile := fmt.Sprintf("%s/frr-mad.pid", socketPath)
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		appLogger.Error(fmt.Sprintf("Failed to create PID file: %s", err))
		os.Exit(1)
	}
	// Register cleanup to remove PID file on exit
	return pidFile
}

func stopApp(pidFile string) {

	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		//appLogger.Error("No running instance found (PID file not found)")
		os.Exit(1)
	}

	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		//appLogger.Error(fmt.Sprintf("Error reading PID file: %v\n", err))
		os.Exit(1)
	}

	pidStr := strings.TrimSpace(string(pidBytes))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		//appLogger.Error(fmt.Sprintf("Invalid PID in file: %v\n", err))
		os.Exit(1)
	}

	process, err := os.FindProcess(pid)
	if err != nil {
		//appLogger.Error(fmt.Sprintf("Process with PID %d not found: %v\n", pid, err))
		os.Exit(1)
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		//appLogger.Error(fmt.Sprintf("Failed to send SIGTERM to process: %v\n", err))
		os.Exit(1)
	}

	time.Sleep(500 * time.Millisecond)

	err = process.Signal(syscall.Signal(0))
	if err == nil {
		//appLogger.Error("Signal sent, but process is still running. It may take a moment to shut down...")
	} else {
		//appLogger.Info("FRR-MAD successfully stopped")

		if _, err := os.Stat(pidFile); !os.IsNotExist(err) {
			os.Remove(pidFile)
		}
	}
}
