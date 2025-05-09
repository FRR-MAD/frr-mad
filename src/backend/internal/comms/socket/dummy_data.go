package socket

import (
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

func getRouterAnomalyDummy1() *frrProto.Response {
	adv1 := []*frrProto.Advertisement{
		{
			InterfaceAddress: "10.0.0.0", // string
			LinkStateId:      "",         // string
			PrefixLength:     "24",       // string
			LinkType:         "",         // string
		},
		{
			InterfaceAddress: "10.200.155.254", // string
			LinkStateId:      "",               // string
			PrefixLength:     "32",             // string
			LinkType:         "Point-to-Point", // string
		},
	}

	adv2 := []*frrProto.Advertisement{
		{
			InterfaceAddress: "172.22.0.0", // string
			LinkStateId:      "",           // string
			PrefixLength:     "24",         // string
			LinkType:         "",           // string
		},
		{
			InterfaceAddress: "172.31.255.254", // string
			LinkStateId:      "",               // string
			PrefixLength:     "32",             // string
			LinkType:         "Point-to-Point", // string
		},
	}

	result := &frrProto.AnomalyAnalysis{
		RouterAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes:  true,
			HasUnderAdvertisedPrefixes: true,
			HasDuplicatePrefixes:       false,
			HasMisconfiguredPrefixes:   false,
			SuperfluousEntries:         adv1,
			MissingEntries:             adv2,
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
			InterfaceAddress: "",         // string
			LinkStateId:      "10.0.0.0", // string
			PrefixLength:     "24",       // string
			LinkType:         "",         // string
		},
		{
			InterfaceAddress: "",               // string
			LinkStateId:      "10.200.155.254", // string
			PrefixLength:     "32",             // string
			LinkType:         "Point-to-Point", // string
		},
	}

	adv2 := []*frrProto.Advertisement{
		{
			InterfaceAddress: "",           // string
			LinkStateId:      "172.22.0.0", // string
			PrefixLength:     "24",         // string
			LinkType:         "",           // string
		},
		{
			InterfaceAddress: "",               // string
			LinkStateId:      "172.31.255.254", // string
			PrefixLength:     "32",             // string
			LinkType:         "Point-to-Point", // string
		},
	}

	result := &frrProto.AnomalyAnalysis{
		RouterAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes:  true,
			HasUnderAdvertisedPrefixes: true,
			HasDuplicatePrefixes:       false,
			HasMisconfiguredPrefixes:   false,
			SuperfluousEntries:         adv1,
			MissingEntries:             adv2,
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
