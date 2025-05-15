package analyzer

import (
	"time"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/ba2025-ysmprc/frr-mad/src/logger"
)

/*
 */
type Analyzer struct {
	AnalysisResult *frrProto.AnomalyAnalysis
	metrics        *frrProto.FullFRRData
	P2pMap         *frrProto.PeerInterfaceMap
	Logger         *logger.Logger
	config         interface{}
}

func InitAnalyzer(
	config interface{},
	metrics *frrProto.FullFRRData,
	logger *logger.Logger,
) *Analyzer {

	anomalyAnalysis := &frrProto.AnomalyAnalysis{
		RouterAnomaly:       initAnomalyDetection(),
		ExternalAnomaly:     initAnomalyDetection(),
		NssaExternalAnomaly: initAnomalyDetection(),
		RibToFibAnomaly:     initAnomalyDetection(),
		LsdbToRibAnomaly:    initAnomalyDetection(),
	}

	return &Analyzer{
		AnalysisResult: anomalyAnalysis,
		metrics:        metrics,
		P2pMap: &frrProto.PeerInterfaceMap{
			PeerInterfaceToAddress: map[string]string{},
		},
		Logger: logger,
		config: config,
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

// TODO: implement misconfiguredPrefixes functionality
func initAnomalyDetection() *frrProto.AnomalyDetection {
	return &frrProto.AnomalyDetection{
		HasOverAdvertisedPrefixes: false,
		HasUnAdvertisedPrefixes:   false,
		HasDuplicatePrefixes:      false,
		HasMisconfiguredPrefixes:  false,
		SuperfluousEntries:        []*frrProto.Advertisement{},
		MissingEntries:            []*frrProto.Advertisement{},
		DuplicateEntries:          []*frrProto.Advertisement{},
	}
}
