package analyzer

import (
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

type Analyzer struct {
	//	anomalyDetection *AnomalyDetection
	Cache *frrProto.Anomalies
}

func newAnalyzer() *Analyzer {
	return &Analyzer{}
}

// analyze the different ospf anomalies
// call ospf functions
func (c *Analyzer) analyzeAnomalies() {

}
