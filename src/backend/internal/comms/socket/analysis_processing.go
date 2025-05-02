package socket

import (
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
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
