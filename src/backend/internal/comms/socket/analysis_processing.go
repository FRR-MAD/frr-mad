package socket

import (
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
)

func (s *Socket) getRouterAnomaly() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_Anomaly{
			Anomaly: s.anomalies.RouterAnomaly,
		},
	}

	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF Router Anomaly Analysis",
		Data:    value,
	}
}

func (s *Socket) getExternalAnomaly() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_Anomaly{
			Anomaly: s.anomalies.ExternalAnomaly,
		},
	}

	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF External Anomaly Analysis",
		Data:    value,
	}
}

func (s *Socket) getNssaExternalAnomaly() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_Anomaly{
			Anomaly: s.anomalies.NssaExternalAnomaly,
		},
	}

	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF Nssa External Anomaly Analysis",
		Data:    value,
	}
}

func (s *Socket) getLsdbToRibAnomaly() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_Anomaly{
			Anomaly: s.anomalies.LsdbToRibAnomaly,
		},
	}

	return &frrProto.Response{
		Status:  "success",
		Message: "Returning LSDB to RIB Anomaly Analysis",
		Data:    value,
	}
}

func (s *Socket) getRibToFibAnomaly() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_Anomaly{
			Anomaly: s.anomalies.RibToFibAnomaly,
		},
	}

	return &frrProto.Response{
		Status:  "success",
		Message: "Returning RIB to FIB Anomaly Analysis",
		Data:    value,
	}
}

func (s *Socket) getShouldParsedLsdb() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_ParsedAnalyzerData{
			ParsedAnalyzerData: s.parsedAnalyzerData,
		},
	}

	return &frrProto.Response{
		Status:  "success",
		Message: "",
		Data:    value,
	}
}
