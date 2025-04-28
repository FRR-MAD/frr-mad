package analyzer

import (
	"time"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/ba2025-ysmprc/frr-mad/src/logger"
)

/*
 */
type Analyzer struct {
	//	anomalyDetection *AnomalyDetection
	Anomalies *frrProto.Anomalies
	metrics   *frrProto.FullFRRData
	Logger    *logger.Logger
}

func InitAnalyzer(config interface{}, metrics *frrProto.FullFRRData, logger *logger.Logger) *Analyzer {
	anomalies := &frrProto.Anomalies{}

	return &Analyzer{
		Anomalies: anomalies,
		metrics:   metrics,
		Logger:    logger,
	}
}

func StartAnalyzer(analyzer *Analyzer, pollInterval time.Duration) {
	ticker := time.NewTicker(pollInterval)

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			analyzer.AnomalyAnalysis()
		}
	}()

}
