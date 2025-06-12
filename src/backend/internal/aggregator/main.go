package aggregator

import (
	"fmt"
	"time"

	"github.com/frr-mad/frr-mad/src/backend/internal/configs"
	"github.com/frr-mad/frr-mad/src/logger"
)

func InitAggregator(config configs.AggregatorConfig, logger *logger.Logger) *Collector {
	configPath := config.FRRConfigPath
	socketPath := config.SocketPath

	return newCollector(configPath, socketPath, logger)
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
