package analyzer

import (
	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"
)

type Analyzer struct {
	//	anomalyDetection *AnomalyDetection
	Cache *frrProto.Anomalies
}

func newAnalyzer() *Analyzer {
	return &Analyzer{}
}
