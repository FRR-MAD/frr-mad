package aggregator

import (
	"github.com/frr-mad/frr-mad/src/backend/internal/aggregator"
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
)

func getMockData() (*logger.Logger, *aggregator.Collector) {
	mockLoggerInstance, _ := logger.NewApplicationLogger("testing", "/tmp/testing.log")
	mockAggregatorInstance := &aggregator.Collector{
		FullFrrData: &frrProto.FullFRRData{},
	}

	return mockLoggerInstance, mockAggregatorInstance
}
