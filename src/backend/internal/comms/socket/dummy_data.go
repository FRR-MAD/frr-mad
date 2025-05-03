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
	}

	adv2 := []*frrProto.Advertisement{}
	result := &frrProto.AnomalyAnalysis{
		RouterAnomaly: &frrProto.AnomalyDetection{
			HasOverAdvertisedPrefixes:  true,
			HasUnderAdvertisedPrefixes: false,
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
