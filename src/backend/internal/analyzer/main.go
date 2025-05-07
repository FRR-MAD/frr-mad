package analyzer

import (
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

/*
 */
type Analyzer struct {
	//	anomalyDetection *AnomalyDetection
	Cache     *frrProto.Anomalies
	Collector *aggregator.Collector
	Logger    *logger.Logger
}

func InitAnalyzer(config map[string]string, collector *aggregator.Collector, logger *logger.Logger) *Analyzer {

	return &Analyzer{
		Cache:     &frrProto.Anomalies{},
		Collector: collector,
		Logger:    logger,
	}
}

func StartAnalyzer(analyzer *Analyzer, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			analyzer.AnomalyAnalysis()
			// fmt.Printf("Analyzer: \n%+v\n", analyzer.Collector.FullFrrData.OspfAsbrSummaryData)
		}
	}()

}
