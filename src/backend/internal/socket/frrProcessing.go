package socket

import (
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
)

func (s *Socket) getRouterName() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_FrrRouterData{
			FrrRouterData: s.Metrics.GetFrrRouterData(),
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
			SystemMetrics: s.Metrics.GetSystemMetrics(),
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
			RoutingInformationBase: s.Metrics.GetRoutingInformationBase(),
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
			RibFibSummaryRoutes: s.Metrics.GetRibFibSummaryRoutes(),
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
			StaticFrrConfiguration: s.Metrics.GetStaticFrrConfiguration(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning static FRR configuration",
		Data:    value,
	}
}
