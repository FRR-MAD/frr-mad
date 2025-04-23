package analyzer

import (
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

/*
 */
type Analyzer struct {
	//	anomalyDetection *AnomalyDetection
	Anomalies *frrProto.Anomalies
	metrics   *frrProto.FullFRRData
	Logger    *logger.Logger
}

func initAnomalies() *frrProto.Anomalies {
	return &frrProto.Anomalies{}
}

func InitAnalyzer(config map[string]string, metrics *frrProto.FullFRRData, logger *logger.Logger) *Analyzer {
	anomalies := initAnomalies()

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
