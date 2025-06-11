package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"

	"github.com/frr-mad/frr-mad/src/backend/configs"
	"github.com/frr-mad/frr-mad/src/backend/internal/aggregator"
	"github.com/frr-mad/frr-mad/src/backend/internal/analyzer"
	socket "github.com/frr-mad/frr-mad/src/backend/internal/socket"
	"github.com/frr-mad/frr-mad/src/backend/internal/exporter"
	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/spf13/cobra"
)

var (
	DaemonVersion = "unknown"
	TUIVersion    = "unknown"
	GitCommit     = "unknown"
	BuildDate     = "unknown"
	RepoURL       = "https://github.com/frr-mad/frr-mad"
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
}

type ServiceConfig struct {
	basis      configs.DefaultConfig
	socket     configs.SocketConfig
	aggregator configs.AggregatorConfig
	exporter   configs.ExporterConfig
}

type LoggerService struct {
	Application *logger.Logger
	Anomaly     *logger.Logger
}

// TODO: create status command
func main() {
	var configFile string
	var rootCmd = &cobra.Command{
		Use:   os.Args[0],
		Short: "FRR-MAD application",
		Long:  `A CLI tool for managing the FRR-MAD application.`,
	}

	var startCmd = &cobra.Command{
		Use:   "start",
		Short: "Start the FRR-MAD application",
		Run: func(cmd *cobra.Command, args []string) {
			app := loadMadApplication(configFile)
			if os.Getenv("FRR_MAD_DAEMON") != "1" {
				command := exec.Command(os.Args[0], os.Args[1:]...)
				command.Env = append(os.Environ(), "FRR_MAD_DAEMON=1")
				command.Start()

				app.Logger.Application.WithAttrs(map[string]interface{}{
					"child_pid":   command.Process.Pid,
					"config_file": configFile,
				}).Info("FRR-MAD daemon started")
				os.Exit(0)
			} else {
				app.startApp(cmd)
			}
		},
	}

	var stopCmd = &cobra.Command{
		Use:   "stop",
		Short: "Stop the FRR-MAD application",
		Run: func(cmd *cobra.Command, args []string) {
			app := loadMadApplication(configFile)
			app.stopApp()
		},
	}

	var restartCmd = &cobra.Command{
		Use:   "restart",
		Short: "Restart the FRR-MAD application",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Restart FRR-MAD application...")
			fmt.Println("Not implemented yet. Please restart the application manually with stop and start.")
		},
	}

	var reloadCmd = &cobra.Command{
		Use:    "reload",
		Short:  "Reload the FRR-MAD configuration",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Reloading the FRR-MAD configuration...")
			fmt.Println("Not implemented yet. Please restart the application manually with stop and start.")
		},
	}

	var debugCmd = &cobra.Command{
		Use:    "debug",
		Short:  "Run the application in debug mode",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			app := loadMadApplication(configFile)
			app.startApp(cmd)
		},
	}
	var testCmd = &cobra.Command{
		Use:    "test",
		Short:  "testing",
		Hidden: true,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println(os.LookupEnv("FRR_MAD_CONFFILE"))
			fmt.Println(configs.ConfigLocation)
		},
	}

	var versionCmd = &cobra.Command{
		Use:   "version",
		Short: "show version number and exit",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Mad Analyzer Daemon Version: %s\n", DaemonVersion)
			fmt.Printf("Mad TUI Version: %s\n", TUIVersion)
			fmt.Printf("Commit: %s\n", GitCommit)
			fmt.Printf("Build Date: %s\n", BuildDate)
			fmt.Printf("Repository: %s\n", RepoURL)
			fmt.Printf("Config Location: %s\n", configs.ConfigLocation)
		},
	}

	startCmd.Flags().StringVarP(&configFile, "configFile", "c", "", "Provide path overwriting default configuration file location.")
	debugCmd.Flags().StringVarP(&configFile, "configFile", "c", "", "Provide path overwriting default configuration file location.")
	startCmd.Flags().Bool("ospf-router", false, "Enable OSPF router metrics")
	startCmd.Flags().Bool("ospf-network", false, "Enable OSPF network metrics")
	startCmd.Flags().Bool("ospf-summary", false, "Enable OSPF summary metrics")
	startCmd.Flags().Bool("ospf-asbr-summary", false, "Enable OSPF ASBR summary metrics")
	startCmd.Flags().Bool("ospf-external", false, "Enable OSPF external route metrics")
	startCmd.Flags().Bool("ospf-nssa-external", false, "Enable OSPF NSSA external route metrics")
	startCmd.Flags().Bool("ospf-database", false, "Enable OSPF database metrics")
	startCmd.Flags().Bool("ospf-neighbors", false, "Enable OSPF neighbor metrics")
	startCmd.Flags().Bool("interface-list", false, "Enable interface list metrics")
	startCmd.Flags().Bool("route-list", false, "Enable route list metrics")

	debugCmd.Flags().Bool("ospf-router", false, "Enable OSPF router metrics")
	debugCmd.Flags().Bool("ospf-network", false, "Enable OSPF network metrics")
	debugCmd.Flags().Bool("ospf-summary", false, "Enable OSPF summary metrics")
	debugCmd.Flags().Bool("ospf-asbr-summary", false, "Enable OSPF ASBR summary metrics")
	debugCmd.Flags().Bool("ospf-external", false, "Enable OSPF external route metrics")
	debugCmd.Flags().Bool("ospf-nssa-external", false, "Enable OSPF NSSA external route metrics")
	debugCmd.Flags().Bool("ospf-database", false, "Enable OSPF database metrics")
	debugCmd.Flags().Bool("ospf-neighbors", false, "Enable OSPF neighbor metrics")
	debugCmd.Flags().Bool("interface-list", false, "Enable interface list metrics")
	debugCmd.Flags().Bool("route-list", false, "Enable route list metrics")

	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(restartCmd)
	rootCmd.AddCommand(reloadCmd)
	rootCmd.AddCommand(debugCmd)
	rootCmd.AddCommand(versionCmd)
	rootCmd.AddCommand(testCmd)

	rootCmd.CompletionOptions.DisableDefaultCmd = false

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func (a *FrrMadApp) startApp(cmd *cobra.Command) {
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
		serviceLogger := a.Logger.Application.WithComponent(service)
		serviceLogger.Info("Starting service")

		switch service {
		case "analyzer":
			if a.Aggregator == nil {
				aggLogger := serviceLogger.WithComponent("aggregator")
				a.Aggregator = startAggregator(a.Config.aggregator, aggLogger, a.PollInterval)
			}

			analyzerLogger := serviceLogger.WithComponent("analyzer")
			a.Analyzer = startAnalyzer(analyzerLogger, a.Logger.Anomaly, a.PollInterval, a.Aggregator)

		case "exporter":
			if a.Exporter == nil {
				expLogger := serviceLogger.WithComponent("exporter")
				getFlagConfigsFromCmd(cmd, &a.Config.exporter)
				a.Exporter = startExporter(a.Config.exporter, expLogger, a.PollInterval, a.Aggregator.FullFrrData, a.Analyzer.AnalysisResult)
			}
		}
	}

	// TODO: Create a better handler for p2pMapping. This should ideally be part of FullFrrData and not a separate data object.
	if a.Aggregator != nil && a.Analyzer != nil && a.Exporter != nil {
		a.Socket = socket.NewSocket(a.Config.socket, a.Aggregator.FullFrrData, a.Analyzer.AnalysisResult, a.Logger.Application, a.Analyzer.AnalyserStateParserResults)

		go func() {
			a.Logger.Application.WithAttrs(map[string]interface{}{
				"socket_path": fmt.Sprintf("%s/%s",
					a.Config.socket.UnixSocketLocation,
					a.Config.socket.UnixSocketName),
			}).Info("Starting socket server")

			if err := a.Socket.Start(); err != nil {
				a.Logger.Application.WithAttrs(map[string]interface{}{
					"error": err.Error(),
					"socket_path": fmt.Sprintf("%s/%s",
						a.Config.socket.UnixSocketLocation,
						a.Config.socket.UnixSocketName),
				}).Error("Socket server failed")

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
	a.Logger.Application.WithAttrs(map[string]interface{}{
		"path": pidFile,
		"pid":  pid,
	}).Debug("Creating PID file")
	if err := os.WriteFile(pidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		a.Logger.Application.Error(fmt.Sprintf("Failed to create PID file: %s", err))
		os.Exit(1)
	}
	return pidFile
}

func (a *FrrMadApp) stopApp() {
	a.Logger.Application.WithAttrs(map[string]interface{}{
		"pid":      a.Pid,
		"pid_file": a.PidFile,
	}).Info("Initiating shutdown")

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
		a.Logger.Application.WithAttrs(map[string]any{
			"pid": a.Pid,
		}).Warning("Process still running after signal")
	} else {
		a.Logger.Application.WithAttrs(map[string]any{
			"pid": a.Pid,
		}).Info("Process terminated successfully")
		if _, err := os.Stat(a.PidFile); !os.IsNotExist(err) {
			os.Remove(a.PidFile)
		}
	}
}

func loadMadApplication(overwriteConfigPath string) *FrrMadApp {
	configRaw, err := configs.LoadConfig(overwriteConfigPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	createDirectories(configRaw)
	config := ServiceConfig{
		basis:      configRaw.Default,
		socket:     configRaw.Socket,
		aggregator: configRaw.Aggregator,
		exporter:   configRaw.Exporter,
	}

	logLevel := logger.ConvertLogLevelFromConfig(config.basis.DebugLevel)
	pollInterval := time.Duration(config.aggregator.PollInterval) * time.Second
	pidFile := fmt.Sprintf("%s/frr-mad.pid", configRaw.Socket.UnixSocketLocation)
	pid, _ := readPidFile(pidFile)

	appLogger, err := logger.NewApplicationLogger("frr-mad",
		fmt.Sprintf("%v/application.log", config.basis.LogPath))
	if err != nil {
		log.Fatalf("Failed to create application logger: %v", err)
	}
	appLogger.SetDebugLevel(logLevel)

	appLogger.Info(fmt.Sprintf("FRR-MAD initializing (version: %s)", DaemonVersion))
	appLogger.WithAttrs(map[string]interface{}{
		"config_path":              overwriteConfigPath,
		"debug_level":              config.basis.DebugLevel,
		"poll_interval in seconds": config.aggregator.PollInterval,
	}).Info("Configuration loaded")

	anomalyLogger, err := logger.NewApplicationLogger("frr-mad-anomaly",
		fmt.Sprintf("%v/anomalies.log", config.basis.LogPath))
	if err != nil {
		log.Fatalf("Failed to create anomaly logger: %v", err)
	}
	anomalyLogger.SetDebugLevel(logLevel)

	logService := &LoggerService{
		Application: appLogger,
		Anomaly:     anomalyLogger,
	}

	app := &FrrMadApp{
		Logger:       logService,
		Pid:          pid,
		PollInterval: pollInterval,
		Config:       config,
		PidFile:      pidFile,
	}

	return app
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

func getFlagConfigsFromCmd(cmd *cobra.Command, exporterConfig *configs.ExporterConfig) {
	if cmd.Flags().Changed("ospf-router") {
		exporterConfig.OSPFRouterData, _ = cmd.Flags().GetBool("ospf-router")
	}
	if cmd.Flags().Changed("ospf-network") {
		exporterConfig.OSPFNetworkData, _ = cmd.Flags().GetBool("ospf-network")
	}
	if cmd.Flags().Changed("ospf-summary") {
		exporterConfig.OSPFSummaryData, _ = cmd.Flags().GetBool("ospf-summary")
	}
	if cmd.Flags().Changed("ospf-asbr-summary") {
		exporterConfig.OSPFAsbrSummaryData, _ = cmd.Flags().GetBool("ospf-asbr-summary")
	}
	if cmd.Flags().Changed("ospf-external") {
		exporterConfig.OSPFExternalData, _ = cmd.Flags().GetBool("ospf-external")
	}
	if cmd.Flags().Changed("ospf-nssa-external") {
		exporterConfig.OSPFNssaExternalData, _ = cmd.Flags().GetBool("ospf-nssa-external")
	}
	if cmd.Flags().Changed("ospf-database") {
		exporterConfig.OSPFDatabase, _ = cmd.Flags().GetBool("ospf-database")
	}
	if cmd.Flags().Changed("ospf-neighbors") {
		exporterConfig.OSPFNeighbors, _ = cmd.Flags().GetBool("ospf-neighbors")
	}
	if cmd.Flags().Changed("interface-list") {
		exporterConfig.InterfaceList, _ = cmd.Flags().GetBool("interface-list")
	}
	if cmd.Flags().Changed("route-list") {
		exporterConfig.RouteList, _ = cmd.Flags().GetBool("route-list")
	}
}
