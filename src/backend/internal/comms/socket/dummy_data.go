package socket

import (
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
)

func getRouterAnomalyDummy1() *frrProto.Response {
	adv1 := []*frrProto.Advertisement{
		{
			InterfaceAddress: "10.0.0.0",
			LinkStateId:      "",
			PrefixLength:     "24",
			LinkType:         "",
		},
		{
			InterfaceAddress: "10.200.155.254",
			LinkStateId:      "",
			PrefixLength:     "32",
			LinkType:         "Point-to-Point",
		},
	}

	adv2 := []*frrProto.Advertisement{
		{
			InterfaceAddress: "172.22.0.0",
			LinkStateId:      "",
			PrefixLength:     "24",
			LinkType:         "",
		},
		{
			InterfaceAddress: "172.31.255.254",
			LinkStateId:      "",
			PrefixLength:     "32",
			LinkType:         "Point-to-Point",
		},
	}

	result := &frrProto.AnomalyAnalysis{
		RouterAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes: true,
			HasUnAdvertisedPrefixes:   true,
			HasDuplicatePrefixes:      false,
			HasMisconfiguredPrefixes:  false,
			SuperfluousEntries:        adv1,
			MissingEntries:            adv2,
		},
	}

	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_Anomaly{
			Anomaly: result.RouterAnomaly,
		},
	}

	return &frrProto.Response{
		Status:  "success",
		Message: "Returning Router Anomaly DUMMY data 1",
		Data:    value,
	}

}

func getExternalAnomalyDummy1() *frrProto.Response {
	adv1 := []*frrProto.Advertisement{
		{
			InterfaceAddress: "",
			LinkStateId:      "10.0.0.0",
			PrefixLength:     "24",
			LinkType:         "",
		},
		{
			InterfaceAddress: "",
			LinkStateId:      "10.200.155.254",
			PrefixLength:     "32",
			LinkType:         "Point-to-Point",
		},
	}

	adv2 := []*frrProto.Advertisement{
		{
			InterfaceAddress: "",
			LinkStateId:      "172.22.0.0",
			PrefixLength:     "24",
			LinkType:         "",
		},
		{
			InterfaceAddress: "",
			LinkStateId:      "172.31.255.254",
			PrefixLength:     "32",
			LinkType:         "Point-to-Point",
		},
	}

	result := &frrProto.AnomalyAnalysis{
		RouterAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes: true,
			HasUnAdvertisedPrefixes:   true,
			HasDuplicatePrefixes:      false,
			HasMisconfiguredPrefixes:  false,
			SuperfluousEntries:        adv1,
			MissingEntries:            adv2,
		},
	}

	value := &frrProto.ResponseValue{
		Kind: &frrProto.ResponseValue_Anomaly{
			Anomaly: result.RouterAnomaly,
		},
	}

	return &frrProto.Response{
		Status:  "success",
		Message: "Returning Router Anomaly DUMMY data 1",
		Data:    value,
	}
}

func getNSSAExternalAnomalyDummy1() *frrProto.Response {
	return getExternalAnomalyDummy1()
}
