package socket

import (
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

func (s *Socket) getRouterName() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_FrrRouterData{
			FrrRouterData: s.metrics.GetFrrRouterData(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning FRR meta data of router itself",
		Data:    value,
	}
}

func (s *Socket) getSystemResources() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_SystemMetrics{
			SystemMetrics: s.metrics.GetSystemMetrics(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning system metrics including CPU and memory",
		Data:    value,
	}
}

func (s *Socket) getRoutingInformationBase() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_RoutingInformationBase{
			RoutingInformationBase: s.metrics.GetRoutingInformationBase(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning all routes (RIB)",
		Data:    value,
	}
}

func (s *Socket) getRibFibSummary() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_RibFibSummaryRoutes{
			RibFibSummaryRoutes: s.metrics.GetRibFibSummaryRoutes(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning route summaries of RIB and FIB",
		Data:    value,
	}
}

func (s *Socket) getStaticFrrConfiguration() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_StaticFrrConfiguration{
			StaticFrrConfiguration: s.metrics.GetStaticFrrConfiguration(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning static FRR configuration",
		Data:    value,
	}
}
