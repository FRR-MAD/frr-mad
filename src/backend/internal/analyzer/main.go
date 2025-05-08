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
	AnalysisResult *frrProto.AnomalyAnalysis
	metrics        *frrProto.FullFRRData
	Logger         *logger.Logger
	config         interface{}
}

func InitAnalyzer(
	config interface{},
	metrics *frrProto.FullFRRData,
	logger *logger.Logger,
) *Analyzer {
	//anomalies := &frrProto.Anomalies{}

	anomalyAnalysis := &frrProto.AnomalyAnalysis{
		RouterAnomaly:       initAnomalyDetection(),
		ExternalAnomaly:     initAnomalyDetection(),
		NssaExternalAnomaly: initAnomalyDetection(),
		FibAnomaly:          initAnomalyDetection(),
	}

	return &Analyzer{
		//Anomalies: anomalies,
		AnalysisResult: anomalyAnalysis,
		metrics:        metrics,
		Logger:         logger,
		config:         config,
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

func initAnomalyDetection() *frrProto.AnomalyDetection {
	return &frrProto.AnomalyDetection{
		HasOverAdvertisedPrefixes:  false,
		HasUnderAdvertisedPrefixes: false,
		HasDuplicatePrefixes:       false,
		HasMisconfiguredPrefixes:   false, // does nothing atm
		SuperfluousEntries:         []*frrProto.Advertisement{},
		MissingEntries:             []*frrProto.Advertisement{},
		DuplicateEntries:           []*frrProto.Advertisement{},
	}
}
