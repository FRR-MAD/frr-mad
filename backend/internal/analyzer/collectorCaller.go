package analyzer

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/ba2025-ysmprc/frr-tui/backend/internal/aggregator"
)

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

func main() {
	// Defaults
	defaultMetricsURL := getEnv("FRR_METRICS_URL", "http://localhost:9342/metrics")
	defaultConfigPath := getEnv("FRR_CONFIG_PATH", "/etc/frr/frr.conf")

	// Flags
	metricsURL := flag.String("metrics-url", defaultMetricsURL, "FRR exporter metrics endpoint")
	configPath := flag.String("config-path", defaultConfigPath, "Path to FRR configuration")
	pollInterval := flag.Duration("poll-interval", 30*time.Second, "Metrics collection interval")
	flag.Parse()

	// Collector init
	collector := aggregator.NewCollector(*metricsURL, *configPath)

	// Collector loop
	ticker := time.NewTicker(*pollInterval)
	defer ticker.Stop()

	for range ticker.C {
		state, err := collector.Collect()
		if err != nil {
			log.Printf("Collection error: %v", err)
			continue
		}

		// TMP logging
		log.Printf("Collected state at %s", state.Timestamp.Format(time.RFC3339))
		log.Printf("OSPF Neighbors: %d", len(state.OSPF.Neighbors))
		log.Printf("OSPF Routes: %d", len(state.OSPF.Routes))
		log.Printf("System CPU: %.1f%%", state.System.CPUUsage)
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}
