package socket

import (
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

func (s *Socket) ospfDummyData() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_StaticFrrConfiguration{
			StaticFrrConfiguration: s.metrics.StaticFrrConfiguration,
		},
	}

	return &frrProto.Response{
		Status:  "success",
		Message: "Returning magical ospf data",
		Data:    value,
	}
}

func (s *Socket) getSystemResources() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_SystemMetrics{
			SystemMetrics: s.metrics.SystemMetrics,
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning magical ospf data",
		Data:    value,
	}
}
