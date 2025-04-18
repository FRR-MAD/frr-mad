package main

import (
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/analyzer"
	socket "github.com/ba2025-ysmprc/frr-mad/src/backend/internal/comms/socket"
)

func main() {

	// load config
	config := configs.LoadConfig()

	aggregatorConfig := config["aggregator"]
	socketConfig := config["socket"]
	analyzerConfig := config["analyzer"]

	// start collector
	collector := aggregator.InitAggregator(aggregatorConfig)

	// Start collection in a goroutine
	pollInterval := time.Duration(strToInt(aggregatorConfig["PollInterval"])) * time.Second
	aggregator.StartAggregator(collector, pollInterval)

	// Start analyzer
	anomalyDetection := analyzer.InitAnalyzer(analyzerConfig)

	analyzer.StartAnalyzer(anomalyDetection, pollInterval)

	// start socket
	sockServer := socket.NewSocket(socketConfig, collector, anomalyDetection)

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := sockServer.Start(); err != nil {
			fmt.Printf("Error starting socket server: %s\n", err)
			os.Exit(1)
		}
	}()

	<-sigChan
	fmt.Println("\nShutting down...")

}

/* ---------------------------------------------------------------------------------------------------+
| With default values																				  |
| ./frr-analytics																					  |
|																									  |
| With custom configuration																			  |
| ./frr-analytics -metrics-url http://192.168.0.5:9342/metrics -config-path /mnt/configs/frr.conf	  |
|																									  |
| Using environment variables																		  |
| export FRR_METRICS_URL="http://10.0.0.5:9342/metrics"												  |
| export FRR_CONFIG_PATH="/etc/frr/frr-prod.conf"													  |
| ./frr-analytics																					  |
+---------------------------------------------------------------------------------------------------- */

// func main() {
// 	// Defaults
// 	defaultMetricsURL := getEnv("FRR_METRICS_URL", "http://localhost:9342/metrics")
// 	defaultConfigPath := getEnv("FRR_CONFIG_PATH", "/etc/frr/frr.conf")
//
// 	// Flags
// 	metricsURL := flag.String("metrics-url", defaultMetricsURL, "FRR exporter metrics endpoint")
// 	configPath := flag.String("config-path", defaultConfigPath, "Path to FRR configuration")
// 	pollInterval := flag.Duration("poll-interval", 30*time.Second, "Metrics collection interval")
// 	flag.Parse()
//
// 	// Collector init
// 	collector := aggregator.NewCollector(*metricsURL, *configPath)
//
// 	// Collector loop
// 	ticker := time.NewTicker(*pollInterval)
// 	defer ticker.Stop()
//
// 	for range ticker.C {
// 		state, err := collector.Collect()
// 		if err != nil {
// 			log.Printf("Collection error: %v", err)
// 			continue
// 		}
//
// 		// TMP logging
// 		log.Printf("Collected state at %s", state.Timestamp.Format(time.RFC3339))
// 		log.Printf("OSPF Neighbors: %d", len(state.OSPF.Neighbors))
// 		log.Printf("OSPF Routes: %d", len(state.OSPF.Routes))
// 		log.Printf("System CPU: %.1f%%", state.System.CPUUsage)
// 	}
// }
//
// func getEnv(key, defaultValue string) string {
// 	if value, exists := os.LookupEnv(key); exists {
// 		return value
// 	}
// 	return defaultValue
// }

//func runAggregator(config map[string]string) *aggregator.Collector {
//
//	metricsURL := config["FRRMetricsURL"]
//	configPath := config["FRRConfigPath"]
//	pollInterval := time.Duration(strToInt(config["PollInterval"])) * time.Second
//
//	// Collector init
//	collector := aggregator.NewCollector(metricsURL, configPath)
//
//	// Collector loop
//	ticker := time.NewTicker(pollInterval)
//	defer ticker.Stop()
//
//	for range ticker.C {
//		//fmt.Println(aggregator.OSPFNeighborDummyData())
//		//_, _ = collector.Collect()
//		state, err := collector.Collect()
//		if err != nil {
//			log.Printf("Collection error: %v", err)
//			continue
//		}
//
//		// TMP logging
//		log.Printf("System metrics: %v", state.System.GetMemoryUsage())
//		log.Printf("System metrics: %v", state.System.GetNetworkStats())
//		//log.Printf("NetworkConfig %v", state.GetConfig())
//		log.Printf("Collected state at %v", state.Timestamp.AsTime())
//		log.Printf("OSPF Neighbors: %d\n", len(state.Ospf.Neighbors))
//		log.Printf("OSPF Routes: %d\n", len(state.Ospf.Routes))
//		log.Printf("System CPU: %.1f%%\n", state.System.CpuUsage)
//	}
//
//	fmt.Println("Yes, it really works")
//
//	return collector
//}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue

}

func runAnalyzer(config map[string]string) {

}

func runSocket(config map[string]string) {
	fmt.Println(config)
}

func strToInt(value string) int {
	retValue, err := strconv.Atoi(value)
	if err != nil {
		// TODO: do proper error handling and get a solution in case it doesn't work
		fmt.Println("Error turning string to int")
	}

	return retValue
}
