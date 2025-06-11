package socket

import (
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
)

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

func (s *Socket) getGeneralOspfInformation() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_GeneralOspfInformation{
			GeneralOspfInformation: s.metrics.GetGeneralOspfInformation(),
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
		Message: "Returning OSPF network data self",
		Data:    value,
	}
}

func (s *Socket) getOspfNetworkDataAll() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_OspfNetworkData{
			OspfNetworkData: s.metrics.GetOspfNetworkDataAll(),
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
		Kind: &frrProto.ResponseValue_OspfExternalAll{
			OspfExternalAll: s.metrics.GetOspfExternalAll(),
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

func (s *Socket) getp2pMap() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_PeerInterfaceToAddress{
			PeerInterfaceToAddress: s.parsedAnalyzerData.P2PMap,
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning compounded P2P OSPF generated Interface Address to static Interface Address",
		Data:    value,
	}
}
