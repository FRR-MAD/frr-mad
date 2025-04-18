package aggregator

import (
	"log"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
)

func InitAggregator(config map[string]string, logger *logger.Logger) *Collector {
	metricsURL := config["FRRMetricsURL"]
	configPath := config["FRRConfigPath"]

	return NewCollector(metricsURL, configPath, logger)
}

func StartAggregator(collector *Collector, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)

	go func() {
		defer ticker.Stop()

		for range ticker.C {
			state, err := collector.Collect()
			if err != nil {
				log.Printf("Collection error: %v", err)
				continue
			}

			// TMP logging
			log.Printf("Collected state at %v", state.Timestamp.AsTime())
			log.Printf("OSPF Neighbors: %d\n", len(state.Ospf.Neighbors))
			log.Printf("OSPF Routes: %d\n", len(state.Ospf.Routes))
			log.Printf("System CPU: %.1f%%\n", state.System.CpuUsage)
		}
	}()
}
