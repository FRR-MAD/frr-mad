package aggregator

import (
	"fmt"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/configs"
	"github.com/ba2025-ysmprc/frr-mad/src/logger"
)

func InitAggregator(config configs.AggregatorConfig, logger *logger.Logger) *Collector {
	metricsURL := config.FRRMetricsURL
	configPath := config.FRRConfigPath
	socketPath := config.SocketPath

	return newCollector(metricsURL, configPath, socketPath, logger)
}

func StartAggregator(collector *Collector, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)
	go func() {
		defer ticker.Stop()

		for range ticker.C {
			err := collector.Collect()
			if err != nil {
				collector.logger.Error(fmt.Sprintf("Collection error: %v", err))
				continue
			}
		}
	}()
}
