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
	AnomalyLogger              *logger.Logger
	config                     any
}

func InitAnalyzer(
	config any,
	metrics *frrProto.FullFRRData,
	logger *logger.Logger,
	anomalyLogger *logger.Logger,
) *Analyzer {
	logger.Info("Initializing analyzer")

	anomalyAnalysis := &frrProto.AnomalyAnalysis{
		RouterAnomaly:       initAnomalyDetection(),
		ExternalAnomaly:     initAnomalyDetection(),
		NssaExternalAnomaly: initAnomalyDetection(),
		RibToFibAnomaly:     initAnomalyDetection(),
		LsdbToRibAnomaly:    initAnomalyDetection(),
	}

	logger.Debug("Created empty anomaly detection structures")

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
		Logger:        logger,
		AnomalyLogger: anomalyLogger,
		config:        config,
	}
}

func StartAnalyzer(analyzer *Analyzer, pollInterval time.Duration) {
	analyzer.Logger.WithAttrs(map[string]any{
		"interval": pollInterval.String(),
	}).Info("Starting analyzer")

	ticker := time.NewTicker(pollInterval)

	go func() {
		defer ticker.Stop()
		for range ticker.C {
			start := time.Now()
			analyzer.AnomalyAnalysis()
			analyzer.Logger.WithAttrs(map[string]any{
				"duration": time.Since(start).String(),
			}).Debug("Completed analysis cycle")
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
