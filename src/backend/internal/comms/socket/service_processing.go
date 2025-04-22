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
			SystemMetrics: s.metrics.GetSystemMetrics(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning magical system data",
		Data:    value,
	}
}

func (s *Socket) getOspfDatabase() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_OspfDatabase{
			OspfDatabase: s.metrics.GetOspfDatabase(),
		},
	}

	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF database",
		Data:    value,
	}
}

func (s *Socket) getOspfRouterData() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_OspfRouterData{
			OspfRouterData: s.metrics.GetOspfRouterData(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF router data",
		Data:    value,
	}
}

func (s *Socket) getOspfNetworkData() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_OspfNetworkData{
			OspfNetworkData: s.metrics.GetOspfNetworkData(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF network data",
		Data:    value,
	}
}

func (s *Socket) getOspfSummaryData() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_OspfSummaryData{
			OspfSummaryData: s.metrics.GetOspfSummaryData(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF summary data",
		Data:    value,
	}
}

func (s *Socket) getOspfAsbrSummaryData() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_OspfAsbrSummaryData{
			OspfAsbrSummaryData: s.metrics.GetOspfAsbrSummaryData(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF ASBR summary data",
		Data:    value,
	}
}

func (s *Socket) getOspfExternalData() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_OspfExternalData{
			OspfExternalData: s.metrics.GetOspfExternalData(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF external data",
		Data:    value,
	}
}

func (s *Socket) getOspfNssaExternalData() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_OspfNssaExternalData{
			OspfNssaExternalData: s.metrics.GetOspfNssaExternalData(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF NSSA external data",
		Data:    value,
	}
}

func (s *Socket) getOspfDuplicates() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_OspfDuplicates{
			OspfDuplicates: s.metrics.GetOspfDuplicates(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF duplicates",
		Data:    value,
	}
}

func (s *Socket) getOspfNeighbors() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_OspfNeighbors{
			OspfNeighbors: s.metrics.GetOspfNeighbors(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning OSPF neighbors",
		Data:    value,
	}
}

func (s *Socket) getInterfaces() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_Interfaces{
			Interfaces: s.metrics.GetInterfaces(),
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning interfaces",
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
		Message: "Returning routes",
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
