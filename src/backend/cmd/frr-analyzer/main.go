package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/analyzer"
	socket "github.com/ba2025-ysmprc/frr-mad/src/backend/internal/comms/socket"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/exporter"
	"github.com/ba2025-ysmprc/frr-mad/src/logger"
)

type Service struct {
	Name   string
	Active bool
}

type FrrMadApp struct {
	Analyzer     *analyzer.Analyzer
	Aggregator   *aggregator.Collector
	Exporter     *exporter.Exporter
	Socket       *socket.Socket
	Logger       *LoggerService
	Config       ServiceConfig
	Pid          int
	PidFile      string
	PollInterval time.Duration
	DebugLevel   int
}

type ServiceConfig struct {
	basis      configs.DefaultConfig
	socket     configs.SocketConfig
	aggregator configs.AggregatorConfig
	analyzer   configs.AnalyzerConfig
	exporter   configs.ExporterConfig
}

type LoggerService struct {
	Application *logger.Logger
}

func main() {
	cmdSet := flag.NewFlagSet("frr-mad", flag.ExitOnError)

	cmdSet.Usage = func() {
		fmt.Println("Usage: frr-mad [command] [options]")
		fmt.Println("\nCommands:")
		fmt.Println("  start   - Start the FRR Monitoring and Analysis Daemon")
		fmt.Println("  stop    - Stop a running FRR-MAD instance")
		fmt.Println("  reload  - Reload configuration for a running FRR-MAD")
		fmt.Println("  help    - Display this help message")
	}

	if len(os.Args) < 2 {
		cmdSet.Usage()
		os.Exit(1)
	}

	configRaw, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	createDirectories(configRaw)
	config := ServiceConfig{
		basis:      configRaw.Default,
		socket:     configRaw.Socket,
		aggregator: configRaw.Aggregator,
		analyzer:   configRaw.Analyzer,
		exporter:   configRaw.Exporter,
	}

	debugLevel := getDebugLevel(config.basis.DebugLevel)
	appLogger := createLogger("frr_mad", fmt.Sprintf("%v/frr_mad.log", config.basis.LogPath))
	appLogger.SetDebugLevel(debugLevel)
	appLogger.Info("Starting FRR Monitoring and Analysis Daemon")

	pollInterval := time.Duration(config.aggregator.PollInterval) * time.Second
	appLogger.Info(fmt.Sprintf("Setting poll interval to %v seconds", config.aggregator.PollInterval))

	logService := &LoggerService{
		Application: appLogger,
	}

	pidFile := fmt.Sprintf("%s/frr-mad.pid", configRaw.Socket.UnixSocketLocation)
	pid, _ := readPidFile(pidFile)
	app := &FrrMadApp{
		Logger:       logService,
		Pid:          pid,
		PollInterval: pollInterval,
		Config:       config,
		DebugLevel:   debugLevel,
		PidFile:      pidFile,
	}

	command := os.Args[1]

	switch command {
	case "start":
		if os.Getenv("FRR_MAD_DAEMON") != "1" {
			cmd := exec.Command(os.Args[0], os.Args[1:]...)
			cmd.Env = append(os.Environ(), "FRR_MAD_DAEMON=1")
			cmd.Start()

			app.Logger.Application.Info(fmt.Sprintf("FRR-MAD started with PID %d", cmd.Process.Pid))
			os.Exit(0)
		} else {
			app.startApp()
		}
	case "stop":
		app.stopApp()
	case "restart":
		fmt.Println("Restart FRR-MAD application...")
		fmt.Println("Not implemented yet. Please restart the application manually.")
	case "reload":
		fmt.Println("Reloading the FRR-MAD configuration...")
		fmt.Println("Not implemented yet. Please restart the application manually.")
	case "help":
		cmdSet.Usage()
	case "debug":
		app.startApp()
	default:
		fmt.Printf("Unknown command: %s\n", command)
		cmdSet.Usage()
		os.Exit(1)
	}
}

func (a *FrrMadApp) startApp() {
	if isProcessRunning(a.Pid) && a.Pid != 0 {
		fmt.Println("FRR-MAD is already running")
		a.Logger.Application.Error("FRR-MAD is already running")
		os.Exit(0)
	}

	pidFile := a.createPidFile()
	defer os.Remove(pidFile)

	services := []string{}
	services = append(services, "analyzer")
	services = append(services, "exporter")

	for _, service := range services {
		a.Logger.Application.Info(fmt.Sprintf("Starting %s service", service))
		switch service {
		case "analyzer":
			if a.Aggregator == nil {
				aggregatorLogger := createLogger("aggregator", fmt.Sprintf("%v/aggregator.log", a.Config.basis.LogPath))
				aggregatorLogger.SetDebugLevel(a.DebugLevel)
				a.Aggregator = startAggregator(a.Config.aggregator, aggregatorLogger, a.PollInterval)
			}

			analyzerLogger := createLogger("analyzer", fmt.Sprintf("%v/analyzer.log", a.Config.basis.LogPath))
			analyzerLogger.SetDebugLevel(a.DebugLevel)
			a.Analyzer = startAnalyzer(a.Config.analyzer, analyzerLogger, a.PollInterval, a.Aggregator)

		case "aggregator":
			if a.Aggregator == nil {
				aggregatorLogger := createLogger("aggregator", fmt.Sprintf("%v/aggregator.log", a.Config.basis.LogPath))
				aggregatorLogger.SetDebugLevel(a.DebugLevel)
				a.Aggregator = startAggregator(a.Config.aggregator, aggregatorLogger, a.PollInterval)
			}

		case "exporter":
			if a.Exporter == nil {
				exporterLogger := createLogger("exporter", fmt.Sprintf("%v/exporter.log", a.Config.basis.LogPath))
				exporterLogger.SetDebugLevel(a.DebugLevel)
				a.Exporter = startExporter(a.Config.exporter, exporterLogger, a.PollInterval, a.Aggregator.FullFrrData, a.Analyzer.AnalysisResult)
			}
		}
	}

	// TODO: Create a better handler for p2pMapping. This should ideally be part of FullFrrData and not a separate data object.
	if a.Aggregator != nil && a.Analyzer != nil && a.Exporter != nil {
		a.Socket = socket.NewSocket(a.Config.socket, a.Aggregator.FullFrrData, a.Analyzer.AnalysisResult, a.Analyzer.Logger, a.Analyzer.P2pMap)

		go func() {
			if err := a.Socket.Start(); err != nil {
				a.Logger.Application.Error(fmt.Sprintf("Error starting socket server: %s", err))
				os.Exit(1)
			}
		}()

		a.Logger.Application.Info(fmt.Sprintf("Socket server listening at %s/%s",
			a.Config.socket.UnixSocketLocation, a.Config.socket.UnixSocketName))
	} else {
		a.Logger.Application.Error("Cannot start socket server: required services not available")
		os.Exit(1)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGTERM)
	<-sigChan

	a.Logger.Application.Info("Received shutdown signal")

	if a.Socket != nil {
		a.Socket.Close()
	}

	a.Logger.Application.Info("FRR-MAD shutdown complete")
}

func (a *FrrMadApp) createPidFile() string {
	pid := os.Getpid()
	pidFile := fmt.Sprintf("%s/frr-mad.pid", a.Config.socket.UnixSocketLocation)
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		a.Logger.Application.Error(fmt.Sprintf("Failed to create PID file: %s", err))
		os.Exit(1)
	}
	return pidFile
}

func (a *FrrMadApp) stopApp() {
	if a.Pid == 0 {
		a.Logger.Application.Error("Service is not running or PID file not found.")
		os.Exit(1)
	}

	process, err := os.FindProcess(a.Pid)
	if err != nil {
		a.Logger.Application.Error(fmt.Sprintf("Process with PID %d not found: %v", a.Pid, err))
		os.Exit(1)
	}

	err = process.Signal(syscall.SIGTERM)
	if err != nil {
		a.Logger.Application.Error(fmt.Sprintf("Failed to send SIGTERM to process: %v", err))
		os.Exit(1)
	}

	time.Sleep(500 * time.Millisecond)

	if isProcessRunning(a.Pid) {
		a.Logger.Application.Error("Signal sent, but process is still running. It may take a moment to shut down...")
	} else {
		a.Logger.Application.Info("FRR-MAD successfully stopped")
		if _, err := os.Stat(a.PidFile); !os.IsNotExist(err) {
			os.Remove(a.PidFile)
		}
	}
}

func readPidFile(pidFile string) (int, error) {
	if _, err := os.Stat(pidFile); os.IsNotExist(err) {
		return 0, fmt.Errorf("PID file not found")
	}

	pidBytes, err := os.ReadFile(pidFile)
	if err != nil {
		return 0, fmt.Errorf("error reading PID file: %v", err)
	}

	pidStr := strings.TrimSpace(string(pidBytes))
	pid, err := strconv.Atoi(pidStr)
	if err != nil {
		return 0, fmt.Errorf("invalid PID in file: %v", err)
	}

	return pid, nil
}

func isProcessRunning(pid int) bool {
	process, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = process.Signal(syscall.Signal(0))
	return err == nil
}

func createLogger(name, filePath string) *logger.Logger {
	logger, err := logger.NewLogger(name, filePath)
	if err != nil {
		log.Fatalf("Failed to create logger %s: %v", name, err)
	}
	return logger
}

func getDebugLevel(level string) int {
	switch level {
	case "none":
		return 99
	case "debug":
		return 2
	case "error":
		return 1
	default:
		return 0
	}
}
