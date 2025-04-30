package socket

import (
	"fmt"

	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

func (s *Socket) getRouterAnomaly() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_Anomaly{
			Anomaly: s.anomalies.RouterAnomaly,
		},
	}

	fmt.Println(s.anomalies.RouterAnomaly)
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF Router Anomaly Analysis",
		Data:    value,
	}
}
