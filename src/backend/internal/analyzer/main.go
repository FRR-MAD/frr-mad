package analyzer

import (
	"time"

	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
)

/*
 */
type Analyzer struct {
	AnalysisResult             *frrProto.AnomalyAnalysis
	AnalyserStateParserResults *frrProto.ParsedAnalyzerData
	metrics                    *frrProto.FullFRRData
	P2pMap                     *frrProto.PeerInterfaceMap
	Logger                     *logger.Logger
	config                     interface{}
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

	analyserStateParserResults := &frrProto.ParsedAnalyzerData{
		ShouldRouterLsdb:       &frrProto.IntraAreaLsa{},
		ShouldExternalLsdb:     &frrProto.InterAreaLsa{},
		ShouldNssaExternalLsdb: &frrProto.InterAreaLsa{},
		P2PMap: &frrProto.PeerInterfaceMap{
			PeerInterfaceToAddress: map[string]string{},
		},
	}

	return &Analyzer{
		AnalysisResult:             anomalyAnalysis,
		AnalyserStateParserResults: analyserStateParserResults,
		metrics:                    metrics,
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
