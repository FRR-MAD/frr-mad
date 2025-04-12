package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ba2025-ysmprc/frr-tui/backend/configs"
	socket "github.com/ba2025-ysmprc/frr-tui/backend/internal/comms/socket"
)

func main() {

	config := configs.LoadConfig()

	aggregatorConfig := config["aggregator"]
	socketConfig := config["socket"]

	sockServer := socket.NewSocket(config["socket"]["UnixSocketLocation"])
	fmt.Println(config["socket"]["UnixSocketLocation"])
	// sockServer := socket.NewSocket("config")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := sockServer.Start(); err != nil {
			fmt.Printf("Error starting socket server: %s\n", err)
			os.Exit(1)
		}
	}()
	runSocket(socketConfig)
	runAggregator(aggregatorConfig)

	// stopAnalyzer := make(chan bool)
	// go runAnalyzerProcess(config["analyzer"], stopAnalyzer)

	/*
				Start up all the different applica.
		t


			tions
				- aggregator
				- exporter
				- anlazyer
				- logger
	*/

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

func runAggregator(config map[string]string) {
	fmt.Println(config)
}

func runSocket(config map[string]string) {
	fmt.Println(config)
}
