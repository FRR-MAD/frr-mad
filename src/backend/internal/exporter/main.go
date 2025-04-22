package exporter

import (
	"fmt"
	"strconv"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
)

func InitExporter(config map[string]string, logger *logger.Logger) (*PrometheusAlerter, error) {
	port := 9091
	if configPort, exists := config["Port"]; exists {
		var err error
		port, err = strconv.Atoi(configPort)
		if err != nil {
			return nil, fmt.Errorf("invalid port in config: %v", err)
		}
	}

	alerter := GetPrometheusAlerter()
	alerter.logger = logger

	if err := alerter.Start(port); err != nil {
		return nil, err
	}
	return alerter, nil
}
