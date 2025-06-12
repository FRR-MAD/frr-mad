package socket

import (
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
)

func (s *Socket) getOspfDatabase() *frrProto.Response {
	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_OspfDatabase{
			OspfDatabase: s.Metrics.GetOspfDatabase(),
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
			GeneralOspfInformation: s.Metrics.GetGeneralOspfInformation(),
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
			OspfRouterData: s.Metrics.GetOspfRouterData(),
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
			OspfNetworkData: s.Metrics.GetOspfNetworkData(),
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
			OspfNetworkData: s.Metrics.GetOspfNetworkDataAll(),
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
			OspfSummaryData: s.Metrics.GetOspfSummaryData(),
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
			OspfAsbrSummaryData: s.Metrics.GetOspfAsbrSummaryData(),
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
			OspfExternalData: s.Metrics.GetOspfExternalData(),
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
			OspfNssaExternalData: s.Metrics.GetOspfNssaExternalData(),
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
			OspfExternalAll: s.Metrics.GetOspfExternalAll(),
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
			OspfNeighbors: s.Metrics.GetOspfNeighbors(),
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
			Interfaces: s.Metrics.GetInterfaces(),
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
			PeerInterfaceToAddress: s.ParsedAnalyzerData.P2PMap,
		},
	}
	return &frrProto.Response{
		Status:  "success",
		Message: "Returning compounded P2P OSPF generated Interface Address to static Interface Address",
		Data:    value,
	}
}
