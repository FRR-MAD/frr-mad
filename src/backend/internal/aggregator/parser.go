package aggregator

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
	"strings"

	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"google.golang.org/protobuf/encoding/protojson"
)

func ParseGeneralOspfInformation(jsonData []byte) (*frrProto.GeneralOspfInformation, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(jsonData, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal OSPF JSON: %w", err)
	}

	result := &frrProto.GeneralOspfInformation{
		Areas: make(map[string]*frrProto.GeneralInfoOspfArea),
	}

	result.RouterId = getString(raw, "routerId")
	result.TosRoutesOnly = getBool(raw, "tosRoutesOnly")
	result.Rfc2328Conform = getBool(raw, "rfc2328Conform")
	result.SpfScheduleDelayMsecs = int32(getFloat(raw, "spfScheduleDelayMsecs"))
	result.HoldtimeMinMsecs = int32(getFloat(raw, "holdtimeMinMsecs"))
	result.HoldtimeMaxMsecs = int32(getFloat(raw, "holdtimeMaxMsecs"))
	result.HoldtimeMultiplier = int32(getFloat(raw, "holdtimeMultplier"))
	result.SpfLastExecutedMsecs = int64(getFloat(raw, "spfLastExecutedMsecs"))
	result.SpfLastDurationMsecs = int32(getFloat(raw, "spfLastDurationMsecs"))
	result.LsaMinIntervalMsecs = int32(getFloat(raw, "lsaMinIntervalMsecs"))
	result.LsaMinArrivalMsecs = int32(getFloat(raw, "lsaMinArrivalMsecs"))
	result.WriteMultiplier = int32(getFloat(raw, "writeMultiplier"))
	result.RefreshTimerMsecs = int32(getFloat(raw, "refreshTimerMsecs"))
	result.MaximumPaths = int32(getFloat(raw, "maximumPaths"))
	result.Preference = int32(getFloat(raw, "preference"))
	result.AsbrRouter = getString(raw, "asbrRouter")
	result.AbrType = getString(raw, "abrType")
	result.LsaExternalCounter = int32(getFloat(raw, "lsaExternalCounter"))
	result.LsaExternalChecksum = int64(getFloat(raw, "lsaExternalChecksum"))
	result.LsaAsopaqueCounter = int32(getFloat(raw, "lsaAsopaqueCounter"))
	result.LsaAsopaqueChecksum = int64(getFloat(raw, "lsaAsOpaqueChecksum"))
	result.AttachedAreaCounter = int32(getFloat(raw, "attachedAreaCounter"))

	if areasRaw, ok := raw["areas"].(map[string]interface{}); ok {
		for areaID, v := range areasRaw {
			areaMap, ok := v.(map[string]interface{})
			if !ok {
				continue
			}

			area := &frrProto.GeneralInfoOspfArea{
				Backbone:               getBool(areaMap, "backbone"),
				AreaIfTotalCounter:     int32(getFloat(areaMap, "areaIfTotalCounter")),
				AreaIfActiveCounter:    int32(getFloat(areaMap, "areaIfActiveCounter")),
				NbrFullAdjacentCounter: int32(getFloat(areaMap, "nbrFullAdjacentCounter")),
				Authentication:         getString(areaMap, "authentication"),
				SpfExecutedCounter:     int32(getFloat(areaMap, "spfExecutedCounter")),
				LsaNumber:              int32(getFloat(areaMap, "lsaNumber")),
				LsaRouterNumber:        int32(getFloat(areaMap, "lsaRouterNumber")),
				LsaRouterChecksum:      int64(getFloat(areaMap, "lsaRouterChecksum")),
				LsaNetworkNumber:       int32(getFloat(areaMap, "lsaNetworkNumber")),
				LsaNetworkChecksum:     int64(getFloat(areaMap, "lsaNetworkChecksum")),
				LsaSummaryNumber:       int32(getFloat(areaMap, "lsaSummaryNumber")),
				LsaSummaryChecksum:     int64(getFloat(areaMap, "lsaSummaryChecksum")),
				LsaAsbrNumber:          int32(getFloat(areaMap, "lsaAsbrNumber")),
				LsaAsbrChecksum:        int64(getFloat(areaMap, "lsaAsbrChecksum")),
				LsaNssaNumber:          int32(getFloat(areaMap, "lsaNssaNumber")),
				LsaNssaChecksum:        int64(getFloat(areaMap, "lsaNssaChecksum")),
				LsaOpaqueLinkNumber:    int32(getFloat(areaMap, "lsaOpaqueLinkNumber")),
				LsaOpaqueLinkChecksum:  int64(getFloat(areaMap, "lsaOpaqueLinkChecksum")),
				LsaOpaqueAreaNumber:    int32(getFloat(areaMap, "lsaOpaqueAreaNumber")),
				LsaOpaqueAreaChecksum:  int64(getFloat(areaMap, "lsaOpaqueAreaChecksum")),
			}

			result.Areas[areaID] = area
		}
	}

	return result, nil
}

func ParseOSPFRouterLSA(jsonData []byte) (*frrProto.OSPFRouterData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	if routerLinkStates, ok := jsonMap["Router Link States"].(map[string]interface{}); ok {
		transformedStates := make(map[string]interface{})
		for areaID, areaData := range routerLinkStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})

			for lsaID, lsaData := range areaDataMap {
				transformedLSAs[lsaID] = transformRouterLSA(lsaData.(map[string]interface{}))
			}

			transformedStates[areaID] = map[string]interface{}{
				"lsa_entries": transformedLSAs,
			}
		}
		transformedMap["router_states"] = transformedStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFRouterData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFRouterLSAAll(jsonData []byte) (*frrProto.OSPFRouterData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]
	type foo map[string]interface{}

	if routerLinkStates, ok := jsonMap["routerLinkStates"].(map[string]interface{}); ok {
		areaID := ""
		transformedStates := make(map[string]interface{})
		for _, areasData := range routerLinkStates {
			areasDataMap := areasData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})
			for areaId, areaData := range areasDataMap {
				areaID = areaId
				areaDataMap := areaData.([]interface{})
				for _, lsaData := range areaDataMap {
					tmpLSA := transformRouterLSA(lsaData.(map[string]interface{}))
					transformedLSAs[tmpLSA["advertising_router"].(string)] = transformRouterLSA(lsaData.(map[string]interface{}))
				}

				transformedStates[areaID] = map[string]interface{}{
					"lsa_entries": transformedLSAs,
				}
			}
		}
		transformedMap["router_states"] = transformedStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFRouterData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFNetworkLSA(jsonData []byte) (*frrProto.OSPFNetworkData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	if netStates, ok := jsonMap["Net Link States"].(map[string]interface{}); ok {
		transformedStates := make(map[string]interface{})
		for areaID, areaData := range netStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})

			for lsaID, lsaData := range areaDataMap {
				transformedLSAs[lsaID] = transformNetworkLSA(lsaData.(map[string]interface{}))
			}

			transformedStates[areaID] = map[string]interface{}{
				"lsa_entries": transformedLSAs,
			}
		}
		transformedMap["net_states"] = transformedStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFNetworkData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFNetworkLSAAll(jsonData []byte) (*frrProto.OSPFNetworkData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]
	type foo map[string]interface{}

	if netStates, ok := jsonMap["networkLinkStates"].(map[string]interface{}); ok {
		transformedStates := make(map[string]interface{})
		for _, areaData := range netStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})

			key := ""
			for lsaID, lsaData := range areaDataMap {
				for _, v := range lsaData.([]interface{}) {
					key = v.(map[string]interface{})["linkStateId"].(string)
					transformedLSAs[key] = transformNetworkLSA(v.(map[string]interface{}))
				}

				transformedStates[lsaID] = map[string]interface{}{
					"lsa_entries": transformedLSAs,
				}
			}
		}
		transformedMap["net_states"] = transformedStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFNetworkData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFSummaryLSA(jsonData []byte) (*frrProto.OSPFSummaryData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	if netStates, ok := jsonMap["Net Link States"].(map[string]interface{}); ok {
		transformedNetStates := make(map[string]interface{})
		for areaID, areaData := range netStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})

			for lsaID, lsaData := range areaDataMap {
				transformedLSAs[lsaID] = transformNetworkLSA(lsaData.(map[string]interface{}))
			}

			transformedNetStates[areaID] = map[string]interface{}{
				"lsa_entries": transformedLSAs,
			}
		}
		transformedMap["net_states"] = transformedNetStates
	}

	if summaryStates, ok := jsonMap["Summary Link States"].(map[string]interface{}); ok {
		transformedSummaryStates := make(map[string]interface{})
		for areaID, areaData := range summaryStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})

			for lsaID, lsaData := range areaDataMap {
				transformedLSAs[lsaID] = transformSummaryLSA(lsaData.(map[string]interface{}))
			}

			transformedSummaryStates[areaID] = map[string]interface{}{
				"lsa_entries": transformedLSAs,
			}
		}
		transformedMap["summary_states"] = transformedSummaryStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFSummaryData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFSummaryLSAAll(jsonData []byte) (*frrProto.OSPFSummaryData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	if sumStates, ok := jsonMap["summaryLinkStates"].(map[string]interface{}); ok {
		transformedNetStates := make(map[string]interface{})
		areaID := ""
		for _, areaData := range sumStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})
			for areaId, lsaData := range areaDataMap {
				areaID = areaId
				lsaDataList := lsaData.([]interface{})
				for _, lsa := range lsaDataList {
					tmpLSA := transformNetworkLSA(lsa.(map[string]interface{}))
					transformedLSAs[tmpLSA["link_state_id"].(string)] = tmpLSA
				}
			}
			transformedNetStates[areaID] = map[string]interface{}{
				"lsa_entries": transformedLSAs,
			}
		}
		transformedMap["summary_states"] = transformedNetStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFSummaryData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFAsbrSummaryLSA(jsonData []byte) (*frrProto.OSPFAsbrSummaryData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	if asbrStates, ok := jsonMap["ASBR-Summary Link States"].(map[string]interface{}); ok {
		transformedStates := make(map[string]interface{})
		for areaID, areaData := range asbrStates {
			areaDataMap := areaData.(map[string]interface{})
			transformedLSAs := make(map[string]interface{})

			for lsaID, lsaData := range areaDataMap {
				transformedLSAs[lsaID] = transformSummaryLSA(lsaData.(map[string]interface{}))
			}

			transformedStates[areaID] = map[string]interface{}{
				"lsa_entries": transformedLSAs,
			}
		}
		transformedMap["asbr_summary_states"] = transformedStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFAsbrSummaryData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFExternalLSA(jsonData []byte) (*frrProto.OSPFExternalData, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	if extStates, ok := jsonMap["AS External Link States"].(map[string]interface{}); ok {
		transformedStates := make(map[string]interface{})
		for lsaID, lsaData := range extStates {
			transformedStates[lsaID] = transformExternalLSA(lsaData.(map[string]interface{}))
		}
		transformedMap["as_external_link_states"] = transformedStates
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFExternalData
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFNssaExternalLSA(jsonData []byte) (*frrProto.OSPFNssaExternalData, error) {
	var rawData map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawData); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformed := make(map[string]interface{})

	if routerID, ok := rawData["routerId"]; ok {
		transformed["router_id"] = routerID
	}

	if nssaStates, ok := rawData["NSSA-external Link States"].(map[string]interface{}); ok {
		areas := make(map[string]interface{})

		for areaID, areaData := range nssaStates {
			areaDataMap, ok := areaData.(map[string]interface{})
			if !ok {
				continue
			}

			lsas := make(map[string]interface{})
			for lsaID, lsaData := range areaDataMap {
				lsaDataMap, ok := lsaData.(map[string]interface{})
				if !ok {
					continue
				}
				lsas[lsaID] = transformNssaExternalLSA(lsaDataMap)
			}

			areas[areaID] = map[string]interface{}{
				"data": lsas,
			}
		}

		transformed["nssa_external_link_states"] = areas
	}

	transformedJSON, err := json.Marshal(transformed)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed data: %w", err)
	}

	var result frrProto.OSPFNssaExternalData
	opts := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	if err := opts.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFNssaExternalAll(jsonData []byte) (*frrProto.OSPFNssaExternalAll, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformed := make(map[string]interface{})

	if routerID, ok := jsonMap["routerId"]; ok {
		transformed["router_id"] = routerID
	}

	if nssaStates, ok := jsonMap["nssaExternalLinkStates"].(map[string]interface{}); ok {
		areas := make(map[string]interface{})

		for _, routerAreas := range nssaStates {
			areaMap, ok := routerAreas.(map[string]interface{})
			areaID := ""
			if !ok {
				continue
			}
			lsas := make(map[string]interface{})
			for area, linkStates := range areaMap {
				areaID = area
				linkStates := linkStates.([]interface{})
				for _, lsaData := range linkStates {
					lsaDataMap, ok := lsaData.(map[string]interface{})
					lsaID := lsaDataMap["linkStateId"].(string)
					if !ok {
						continue
					}
					lsas[lsaID] = transformNssaExternalLSA(lsaDataMap)

				}
			}
			areas[areaID] = map[string]interface{}{
				"data": lsas,
			}
		}
		transformed["nssa_external_all_link_states"] = areas
	}

	transformedJSON, err := json.Marshal(transformed)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed data: %w", err)
	}

	var result frrProto.OSPFNssaExternalAll
	opts := protojson.UnmarshalOptions{
		AllowPartial:   true,
		DiscardUnknown: true,
	}

	if err := opts.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseFullOSPFDatabase(jsonData []byte) (*frrProto.OSPFDatabase, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	if areas, ok := jsonMap["areas"].(map[string]interface{}); ok {
		transformedAreas := make(map[string]interface{})
		for areaID, areaData := range areas {
			areaDataMap := areaData.(map[string]interface{})
			transformedArea := make(map[string]interface{})

			if routerLSAs, ok := areaDataMap["routerLinkStates"].([]interface{}); ok {
				transformedRouterLSAs := make([]interface{}, len(routerLSAs))
				for i, lsa := range routerLSAs {
					transformedRouterLSAs[i] = transformDatabaseRouterLSA(lsa.(map[string]interface{}))
				}
				transformedArea["router_link_states"] = transformedRouterLSAs
				transformedArea["router_link_states_count"] = areaDataMap["routerLinkStatesCount"]
			}

			if networkLSAs, ok := areaDataMap["networkLinkStates"].([]interface{}); ok {
				transformedNetworkLSAs := make([]interface{}, len(networkLSAs))
				for i, lsa := range networkLSAs {
					transformedNetworkLSAs[i] = transformDatabaseNetworkLSA(lsa.(map[string]interface{}))
				}
				transformedArea["network_link_states"] = transformedNetworkLSAs
				transformedArea["network_link_states_count"] = areaDataMap["networkLinkStatesCount"]
			}

			if summaryLSAs, ok := areaDataMap["summaryLinkStates"].([]interface{}); ok {
				transformedSummaryLSAs := make([]interface{}, len(summaryLSAs))
				for i, lsa := range summaryLSAs {
					transformedSummaryLSAs[i] = transformDatabaseSummaryLSA(lsa.(map[string]interface{}))
				}
				transformedArea["summary_link_states"] = transformedSummaryLSAs
				transformedArea["summary_link_states_count"] = areaDataMap["summaryLinkStatesCount"]
			}

			if asbrSummaryLSAs, ok := areaDataMap["asbrSummaryLinkStates"].([]interface{}); ok {
				transformedASBRLSAs := make([]interface{}, len(asbrSummaryLSAs))
				for i, lsa := range asbrSummaryLSAs {
					transformedASBRLSAs[i] = transformDatabaseASBRSummaryLSA(lsa.(map[string]interface{}))
				}
				transformedArea["asbr_summary_link_states"] = transformedASBRLSAs
				transformedArea["asbr_summary_link_states_count"] = areaDataMap["asbrSummaryLinkStatesCount"]
			}

			if nssaExternalLSAs, ok := areaDataMap["nssaExternalLinkStates"].([]interface{}); ok {
				transformedNSSALSAs := make([]interface{}, len(nssaExternalLSAs))
				for i, lsa := range nssaExternalLSAs {
					transformedNSSALSAs[i] = transformDatabaseNSSAExternalLSA(lsa.(map[string]interface{}))
				}
				transformedArea["nssa_external_link_states"] = transformedNSSALSAs
				transformedArea["nssa_external_link_states_count"] = areaDataMap["nssaExternalLinkStatesCount"]
			}

			transformedAreas[areaID] = transformedArea
		}
		transformedMap["areas"] = transformedAreas
	}

	if extLSAs, ok := jsonMap["asExternalLinkStates"].([]interface{}); ok {
		transformedExtLSAs := make([]interface{}, len(extLSAs))
		for i, lsa := range extLSAs {
			transformedExtLSAs[i] = transformDatabaseExternalLSA(lsa.(map[string]interface{}))
		}
		transformedMap["as_external_link_states"] = transformedExtLSAs
		transformedMap["as_external_count"] = jsonMap["asExternalLinkStatesCount"]
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFDatabase
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFExternalAll(jsonData []byte) (*frrProto.OSPFExternalAll, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	transformedMap["router_id"] = jsonMap["routerId"]

	if extLSAs, ok := jsonMap["asExternalLinkStates"].([]interface{}); ok {
		transformedExtLSAs := make([]interface{}, len(extLSAs))
		for i, lsa := range extLSAs {
			transformedExtLSAs[i] = transformExternalLinkState(lsa.(map[string]interface{}))
		}
		transformedMap["as_external_link_states"] = transformedExtLSAs
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFExternalAll
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseOSPFNeighbors(jsonData []byte) (*frrProto.OSPFNeighbors, error) {
	var jsonMap map[string]interface{}
	if err := json.Unmarshal(jsonData, &jsonMap); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	transformedMap := make(map[string]interface{})
	if neighbors, ok := jsonMap["neighbors"].(map[string]interface{}); ok {
		transformedNeighbors := make(map[string]interface{})
		for iface, neighborList := range neighbors {
			neighborsSlice := neighborList.([]interface{})
			transformedNeighborList := make([]interface{}, len(neighborsSlice))

			for i, neighbor := range neighborsSlice {
				transformedNeighborList[i] = transformNeighbor(neighbor.(map[string]interface{}))
			}

			transformedNeighbors[iface] = map[string]interface{}{
				"neighbors": transformedNeighborList,
			}
		}
		transformedMap["neighbors"] = transformedNeighbors
	}

	transformedJSON, err := json.Marshal(transformedMap)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal transformed map: %w", err)
	}

	var result frrProto.OSPFNeighbors
	unmarshaler := protojson.UnmarshalOptions{AllowPartial: true}
	if err := unmarshaler.Unmarshal(transformedJSON, &result); err != nil {
		return nil, fmt.Errorf("failed to unmarshal to protobuf: %w", err)
	}

	return &result, nil
}

func ParseInterfaceStatus(jsonData []byte) (*frrProto.InterfaceList, error) {
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	result := &frrProto.InterfaceList{
		Interfaces: make(map[string]*frrProto.SingleInterface),
	}

	for ifaceName, ifaceData := range rawResponse {
		ifaceMap, ok := ifaceData.(map[string]interface{})
		if !ok {
			continue
		}

		singleIface := &frrProto.SingleInterface{
			IpAddresses: make([]*frrProto.IpAddress, 0),
		}

		singleIface.AdministrativeStatus = getString(ifaceMap, "administrativeStatus")
		singleIface.OperationalStatus = getString(ifaceMap, "operationalStatus")
		singleIface.LinkDetection = getBool(ifaceMap, "linkDetection")
		singleIface.LinkUps = int32(getFloat(ifaceMap, "linkUps"))
		singleIface.LinkDowns = int32(getFloat(ifaceMap, "linkDowns"))
		singleIface.LastLinkUp = getString(ifaceMap, "lastLinkUp")
		singleIface.LastLinkDown = getString(ifaceMap, "lastLinkDown")
		singleIface.VrfName = getString(ifaceMap, "vrfName")
		singleIface.MplsEnabled = getBool(ifaceMap, "mplsEnabled")
		singleIface.LinkDown = getBool(ifaceMap, "linkDown")
		singleIface.LinkDownV6 = getBool(ifaceMap, "linkDownV6")
		singleIface.McForwardingV4 = getBool(ifaceMap, "mcForwardingV4")
		singleIface.McForwardingV6 = getBool(ifaceMap, "mcForwardingV6")
		singleIface.PseudoInterface = getBool(ifaceMap, "pseudoInterface")
		singleIface.Index = int32(getFloat(ifaceMap, "index"))
		singleIface.Metric = int32(getFloat(ifaceMap, "metric"))
		singleIface.Mtu = int32(getFloat(ifaceMap, "mtu"))
		singleIface.Speed = int32(getFloat(ifaceMap, "speed"))
		singleIface.Flags = getString(ifaceMap, "flags")
		singleIface.Type = getString(ifaceMap, "type")
		singleIface.HardwareAddress = getString(ifaceMap, "hardwareAddress")
		singleIface.InterfaceType = getString(ifaceMap, "interfaceType")
		singleIface.InterfaceSlaveType = getString(ifaceMap, "interfaceSlaveType")
		singleIface.LacpBypass = getBool(ifaceMap, "lacpBypass")
		singleIface.Protodown = getString(ifaceMap, "protodown")
		singleIface.ParentIfindex = int32(getFloat(ifaceMap, "parentIfindex"))

		if ipAddrs, ok := ifaceMap["ipAddresses"].([]interface{}); ok {
			for _, ipAddr := range ipAddrs {
				if ipMap, ok := ipAddr.(map[string]interface{}); ok {
					ip := &frrProto.IpAddress{
						Address:    getString(ipMap, "address"),
						Secondary:  getBool(ipMap, "secondary"),
						Unnumbered: getBool(ipMap, "unnumbered"),
					}
					singleIface.IpAddresses = append(singleIface.IpAddresses, ip)
				}
			}
		}

		if evpnData, ok := ifaceMap["evpnMh"].(map[string]interface{}); ok {
			singleIface.EvpnMh = &frrProto.EvpnMh{
				EthernetSegmentId: getString(evpnData, "ethernetSegmentId"),
				Esi:               getString(evpnData, "esi"),
				DfPreference:      int32(getFloat(evpnData, "dfPreference")),
				DfAlgorithm:       getString(evpnData, "dfAlgorithm"),
				DfStatus:          getString(evpnData, "dfStatus"),
				MultiHomingMode:   getString(evpnData, "multihomingMode"),
				ActiveMode:        getBool(evpnData, "activeMode"),
				BypassMode:        getBool(evpnData, "bypassMode"),
				LocalBias:         getBool(evpnData, "localBias"),
				FastFailover:      getBool(evpnData, "fastFailover"),
				UpTime:            getString(evpnData, "upTime"),
				BgpStatus:         getString(evpnData, "bgpStatus"),
				ProtocolStatus:    getString(evpnData, "protocolStatus"),
				ProtocolDown:      getBool(evpnData, "protocolDown"),
				MacCount:          int32(getFloat(evpnData, "macCount")),
				LocalIfindex:      int32(getFloat(evpnData, "localIfindex")),
				NetworkCount:      int32(getFloat(evpnData, "networkCount")),
				JoinCount:         int32(getFloat(evpnData, "joinCount")),
				LeaveCount:        int32(getFloat(evpnData, "leaveCount")),
			}
		}

		result.Interfaces[ifaceName] = singleIface
	}

	return result, nil
}

func ParseRib(jsonData []byte) (*frrProto.RoutingInformationBase, error) {
	var rawResponse map[string]interface{}
	if err := json.Unmarshal(jsonData, &rawResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	result := &frrProto.RoutingInformationBase{
		Routes: make(map[string]*frrProto.RouteEntry),
	}

	for prefix, routeData := range rawResponse {
		routeSlice, ok := routeData.([]interface{})
		if !ok {
			continue
		}

		routeEntry := &frrProto.RouteEntry{
			Routes: make([]*frrProto.Route, 0, len(routeSlice)),
		}

		for _, r := range routeSlice {
			routeMap, ok := r.(map[string]interface{})
			if !ok {
				continue
			}

			route := &frrProto.Route{
				Prefix:                   getString(routeMap, "prefix"),
				PrefixLen:                int32(getFloat(routeMap, "prefixLen")),
				Protocol:                 getString(routeMap, "protocol"),
				VrfId:                    int32(getFloat(routeMap, "vrfId")),
				VrfName:                  getString(routeMap, "vrfName"),
				Selected:                 getBool(routeMap, "selected"),
				DestSelected:             getBool(routeMap, "destSelected"),
				Distance:                 int32(getFloat(routeMap, "distance")),
				Metric:                   int32(getFloat(routeMap, "metric")),
				Installed:                getBool(routeMap, "installed"),
				Table:                    int32(getFloat(routeMap, "table")),
				InternalStatus:           int32(getFloat(routeMap, "internalStatus")),
				InternalFlags:            int32(getFloat(routeMap, "internalFlags")),
				InternalNextHopNum:       int32(getFloat(routeMap, "internalNextHopNum")),
				InternalNextHopActiveNum: int32(getFloat(routeMap, "internalNextHopActiveNum")),
				NexthopGroupId:           int32(getFloat(routeMap, "nexthopGroupId")),
				InstalledNexthopGroupId:  int32(getFloat(routeMap, "installedNexthopGroupId")),
				Uptime:                   getString(routeMap, "uptime"),
				Nexthops:                 make([]*frrProto.Nexthop, 0),
			}

			if nexthops, ok := routeMap["nexthops"].([]interface{}); ok {
				for _, nh := range nexthops {
					if nhMap, ok := nh.(map[string]interface{}); ok {
						nexthop := &frrProto.Nexthop{
							Flags:             int32(getFloat(nhMap, "flags")),
							Fib:               getBool(nhMap, "fib"),
							DirectlyConnected: getBool(nhMap, "directlyConnected"),
							Duplicate:         getBool(nhMap, "duplicate"),
							Ip:                getString(nhMap, "ip"),
							Afi:               getString(nhMap, "afi"),
							InterfaceIndex:    int32(getFloat(nhMap, "interfaceIndex")),
							InterfaceName:     getString(nhMap, "interfaceName"),
							Active:            getBool(nhMap, "active"),
							Weight:            int32(getFloat(nhMap, "weight")),
						}
						route.Nexthops = append(route.Nexthops, nexthop)
					}
				}
			}

			routeEntry.Routes = append(routeEntry.Routes, route)
		}

		result.Routes[prefix] = routeEntry
	}

	return result, nil
}

func ParseRibFibSummary(jsonData []byte) (*frrProto.RibFibSummaryRoutes, error) {
	var raw map[string]interface{}
	if err := json.Unmarshal(jsonData, &raw); err != nil {
		return nil, fmt.Errorf("failed to unmarshal summary JSON: %w", err)
	}

	result := &frrProto.RibFibSummaryRoutes{
		RouteSummaries: make([]*frrProto.RouteSummary, 0),
	}

	if routes, ok := raw["routes"].([]interface{}); ok {
		for _, r := range routes {
			routeMap, ok := r.(map[string]interface{})
			if !ok {
				continue
			}

			summary := &frrProto.RouteSummary{
				Fib:          int32(getFloat(routeMap, "fib")),
				Rib:          int32(getFloat(routeMap, "rib")),
				FibOffLoaded: int32(getFloat(routeMap, "fibOffLoaded")),
				FibTrapped:   int32(getFloat(routeMap, "fibTrapped")),
				Type:         getString(routeMap, "type"),
			}

			result.RouteSummaries = append(result.RouteSummaries, summary)
		}
	}

	result.RoutesTotal = int32(getFloat(raw, "routesTotal"))
	result.RoutesTotalFib = int32(getFloat(raw, "routesTotalFib"))

	return result, nil
}

func ParseStaticFRRConfig(path string) (*frrProto.StaticFRRConfiguration, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open config file: %w", err)
	}
	defer func(file *os.File) {
		if closeErr := file.Close(); closeErr != nil {
			err = fmt.Errorf("failed to close config file: %w", closeErr)
		}
	}(file)

	config := &frrProto.StaticFRRConfiguration{}
	scanner := bufio.NewScanner(file)

	var currentInterfacePointer *frrProto.Interface

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "!") {
			continue
		}

		if handled := parseMetadataLine(config, line); handled {
			continue
		}

		if newInterface := parseInterfaceLine(line); newInterface != nil {
			if currentInterfacePointer != nil {
				config.Interfaces = append(config.Interfaces, currentInterfacePointer)
			}
			currentInterfacePointer = newInterface
			continue
		}
		if currentInterfacePointer != nil && parseInterfaceSubLine(currentInterfacePointer, line) {
			continue
		}

		if handled := parseStaticRouteLine(config, line); handled {
			continue
		}

		if strings.HasPrefix(line, "router ospf") {
			parseRouterOSPFConfig(scanner, config)
			continue
		}

		if handled := parseAccessListLine(config, line); handled {
			continue
		}

		if handled := parseRouteMapLine(config, line); handled {
			continue
		}

		if handled := parseRouteMapMatchLine(config, line); handled {
			continue
		}
	}

	if currentInterfacePointer != nil {
		config.Interfaces = append(config.Interfaces, currentInterfacePointer)
	}

	return config, nil
}

func parseMetadataLine(config *frrProto.StaticFRRConfiguration, line string) bool {
	parts := strings.Fields(line)
	switch {
	case strings.HasPrefix(line, "hostname "):
		config.Hostname = parts[1]
		return true
	case strings.HasPrefix(line, "frr version "):
		config.FrrVersion = parts[2]
		return true
	case line == "no ipv6 forwarding": // todo: add to tests
		config.Ipv6Forwarding = false // todo: check if default value is true on router
		return true
	case line == "no ipv4 forwarding": // todo: add to tests
		config.Ipv4Forwarding = false
		return true
	case line == "service advanced-vty":
		config.ServiceAdvancedVty = true
		return true
	}
	return false
}

func parseInterfaceLine(line string) *frrProto.Interface {
	if !strings.HasPrefix(line, "interface ") {
		return nil
	}
	parts := strings.Fields(line)
	return &frrProto.Interface{Name: parts[1]}
}

func parseInterfaceSubLine(currentInterfacePointer *frrProto.Interface, line string) bool {
	parts := strings.Fields(line)
	switch {
	case strings.HasPrefix(line, "ip address "):
		if len(parts) == 3 {
			ip, ipNet, err := net.ParseCIDR(parts[2])
			if err != nil || ipNet == nil {
				log.Printf("bad CIDR %q: %v", parts[2], err)
				return true
			}
			prefixLength, _ := ipNet.Mask.Size()
			currentInterfacePointer.InterfaceIpPrefixes = append(currentInterfacePointer.InterfaceIpPrefixes, &frrProto.InterfaceIPPrefix{
				IpPrefix: &frrProto.IPPrefix{
					IpAddress:    ip.String(),
					PrefixLength: uint32(prefixLength),
				},
				Passive: false,
				HasPeer: false,
			})
			return true
		} else if parts[3] == "peer" {
			ip := parts[2]
			peerIp, ipNet, err := net.ParseCIDR(parts[4])
			if err != nil || ipNet == nil {
				log.Printf("bad CIDR %q: %v", parts[2], err)
				return true
			}
			peerIpPrefixLength, _ := ipNet.Mask.Size()
			currentInterfacePointer.InterfaceIpPrefixes = append(currentInterfacePointer.InterfaceIpPrefixes, &frrProto.InterfaceIPPrefix{
				IpPrefix: &frrProto.IPPrefix{
					IpAddress:    ip,
					PrefixLength: 32,
				},
				Passive: false,
				HasPeer: true,
				PeerIpPrefix: &frrProto.IPPrefix{
					IpAddress:    peerIp.String(),
					PrefixLength: uint32(peerIpPrefixLength),
				},
			})
			return true
		}
		return true
	case strings.HasPrefix(line, "ip ospf area "):
		currentInterfacePointer.Area = strings.Fields(line)[3]
		return true
	case strings.HasPrefix(line, "ip ospf passive"):
		if len(parts) < 4 {
			for _, interfaceIPPrefix := range currentInterfacePointer.InterfaceIpPrefixes {
				interfaceIPPrefix.Passive = true
				return true
			}
		} else {
			for _, interfaceIPPrefix := range currentInterfacePointer.InterfaceIpPrefixes {
				if interfaceIPPrefix.IpPrefix.IpAddress == parts[3] {
					interfaceIPPrefix.Passive = true
					return true
				}
			}
		}
		return true
	case line == "exit":
		return true
	}
	return false
}

func parseStaticRouteLine(config *frrProto.StaticFRRConfiguration, line string) bool {
	if !strings.HasPrefix(line, "ip route ") {
		return false
	}
	parts := strings.Fields(line)
	ip, ipNet, err := net.ParseCIDR(parts[2])
	if err != nil || ipNet == nil {
		log.Printf("bad static route CIDR %q", parts[2])
		return true
	}
	prefixLength, _ := ipNet.Mask.Size()
	config.StaticRoutes = append(config.StaticRoutes, &frrProto.StaticRoute{
		IpPrefix: &frrProto.IPPrefix{
			IpAddress:    ip.String(),
			PrefixLength: uint32(prefixLength),
		},
		NextHop: parts[3],
	})
	return true
}

func parseAccessListLine(config *frrProto.StaticFRRConfiguration, line string) bool {
	if !strings.HasPrefix(line, "access-list ") {
		return false
	}
	parts := strings.Fields(line)
	if len(parts) < 6 {
		log.Printf("short access-list line: %q", line)
		return true
	}
	name, sequence, action, target := parts[1], parts[3], parts[4], parts[5]
	seq, _ := strconv.Atoi(sequence)

	item := &frrProto.AccessListItem{
		Sequence:      uint32(seq),
		AccessControl: action,
	}
	if target == "any" {
		item.Destination = &frrProto.AccessListItem_Any{Any: true}
	} else if ip, ipnet, err := net.ParseCIDR(target); err == nil && ipnet != nil {
		prefixLength, _ := ipnet.Mask.Size()
		item.Destination = &frrProto.AccessListItem_IpPrefix{
			IpPrefix: &frrProto.IPPrefix{IpAddress: ip.String(),
				PrefixLength: uint32(prefixLength),
			},
		}
	} else {
		log.Printf("bad CIDR %q in ACL %q", target, line)
	}

	if config.AccessList == nil {
		config.AccessList = make(map[string]*frrProto.AccessList)
	}
	if _, ok := config.AccessList[name]; !ok {
		config.AccessList[name] = &frrProto.AccessList{}
	}
	config.AccessList[name].AccessListItems = append(config.AccessList[name].AccessListItems, item)
	return true
}

func parseRouteMapLine(config *frrProto.StaticFRRConfiguration, line string) bool {
	if !strings.HasPrefix(line, "route-map ") {
		return false
	}
	parts := strings.Fields(line)
	if len(parts) < 4 {
		log.Printf("short route-map line: %q", line)
		return true
	}
	name, action, sequence := parts[1], parts[2], parts[3]
	if config.RouteMap == nil {
		config.RouteMap = make(map[string]*frrProto.RouteMap)
	}
	config.RouteMap[name] = &frrProto.RouteMap{
		Permit:   action == "permit",
		Sequence: sequence,
	}
	return true
}

func parseRouteMapMatchLine(config *frrProto.StaticFRRConfiguration, line string) bool {
	if !strings.HasPrefix(line, "match ip address ") {
		return false
	}
	parts := strings.Fields(line)
	accessListName := parts[3]
	for _, rm := range config.RouteMap {
		if rm.AccessList == "" {
			rm.Match = "ip address"
			rm.AccessList = accessListName
			break
		}
	}
	return true
}

func parseRouterOSPFConfig(scanner *bufio.Scanner, config *frrProto.StaticFRRConfiguration) {
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "exit" {
			break
		}

		switch {
		case strings.HasPrefix(line, "ospf router-id "):
			if config.OspfConfig == nil {
				config.OspfConfig = &frrProto.OSPFConfig{}
			}
			parts := strings.Fields(line)
			config.OspfConfig.RouterId = parts[2]

		case strings.HasPrefix(line, "redistribute "):
			if config.OspfConfig == nil {
				config.OspfConfig = &frrProto.OSPFConfig{}
			}
			parts := strings.Fields(line)
			redistributionType := ""
			redistributionMetric := ""
			redistributionRouteMap := ""
			for i, part := range parts {
				if part == "redistribute" {
					redistributionType = parts[i+1]
				}
				if part == "metric-type" {
					redistributionMetric = parts[i+1]
				}
				if part == "route-map" {
					redistributionRouteMap = parts[i+1]
				}
			}
			config.OspfConfig.Redistribution = append(config.OspfConfig.Redistribution, &frrProto.Redistribution{
				Type:     redistributionType,
				Metric:   redistributionMetric,
				RouteMap: redistributionRouteMap,
			})

		case strings.HasPrefix(line, "area "):
			if config.OspfConfig == nil {
				config.OspfConfig = &frrProto.OSPFConfig{}
			}
			parts := strings.Fields(line)
			area := &frrProto.Area{Name: parts[1]}
			if len(parts) > 2 {
				area.Type = parts[2]
			}
			for i, part := range parts {
				if part == "virtual-link" && i+1 < len(parts) {
					area.Type = "transit (virtual-link)"
					config.OspfConfig.VirtualLinkNeighbor = parts[i+1]
					break
				}
			}
			config.OspfConfig.Area = append(config.OspfConfig.Area, area)
		}
	}
}

// TODO: is this needed?
func addNetworkToArea(config *frrProto.NetworkConfig, network, area string) {
	for i, a := range config.Areas {
		if a.Id == area {
			config.Areas[i].Networks = append(a.Networks, network)
			return
		}
	}
	config.Areas = append(config.Areas, &frrProto.OSPFArea{
		Id:       area,
		Networks: []string{network},
	})
}

func (c *Collector) ReadConfig() (string, error) {
	file, err := os.Open(c.configPath)
	if err != nil {
		return "", fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var staticConfig []string

	for scanner.Scan() {
		staticConfig = append(staticConfig, scanner.Text())
	}

	return strings.Join(staticConfig, "\n"), nil
}

func transformRouterLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"lsaAge":            "lsa_age",
		"options":           "options",
		"lsaFlags":          "lsa_flags",
		"flags":             "flags",
		"asbr":              "asbr",
		"lsaType":           "lsa_type",
		"linkStateId":       "link_state_id",
		"advertisingRouter": "advertising_router",
		"lsaSeqNumber":      "lsa_seq_number",
		"checksum":          "checksum",
		"length":            "length",
		"numOfLinks":        "num_of_links",
		"routerLinks":       "router_links",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := lsaData[origKey]; exists {
			if origKey == "routerLinks" {
				routerLinks := value.(map[string]interface{})
				transformedLinks := make(map[string]interface{})

				for linkID, linkData := range routerLinks {
					transformedLinks[linkID] = transformRouterLink(linkData.(map[string]interface{}))
				}

				transformed[newKey] = transformedLinks
			} else {
				transformed[newKey] = value
			}
		}
	}

	return transformed
}

func transformRouterLink(linkData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"linkType":                "link_type",
		"designatedRouterAddress": "designated_router_address",
		"neighborRouterId":        "neighbor_router_id",
		"routerInterfaceAddress":  "router_interface_address",
		"networkAddress":          "network_address",
		"networkMask":             "network_mask",
		"numOfTosMetrics":         "num_of_tos_metrics",
		"tos0Metric":              "tos0_metric",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := linkData[origKey]; exists {
			transformed[newKey] = value
		}
	}

	return transformed
}

func transformNetworkLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"lsaAge":            "lsa_age",
		"options":           "options",
		"lsaFlags":          "lsa_flags",
		"lsaType":           "lsa_type",
		"linkStateId":       "link_state_id",
		"advertisingRouter": "advertising_router",
		"lsaSeqNumber":      "lsa_seq_number",
		"checksum":          "checksum",
		"length":            "length",
		"networkMask":       "network_mask",
		"attachedRouters":   "attached_routers",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := lsaData[origKey]; exists {
			if origKey == "attachedRouters" {
				routers := value.(map[string]interface{})
				transformedRouters := make(map[string]interface{})

				for routerID, routerData := range routers {
					transformedRouters[routerID] = map[string]interface{}{
						"attached_router_id": routerData.(map[string]interface{})["attachedRouterId"],
					}
				}

				transformed[newKey] = transformedRouters
			} else {
				transformed[newKey] = value
			}
		}
	}

	return transformed
}

func transformSummaryLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"lsaAge":            "lsa_age",
		"options":           "options",
		"lsaFlags":          "lsa_flags",
		"lsaType":           "lsa_type",
		"linkStateId":       "link_state_id",
		"advertisingRouter": "advertising_router",
		"lsaSeqNumber":      "lsa_seq_number",
		"checksum":          "checksum",
		"length":            "length",
		"networkMask":       "network_mask",
		"tos0Metric":        "tos0_metric",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := lsaData[origKey]; exists {
			transformed[newKey] = value
		}
	}

	return transformed
}

func transformExternalLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"lsaAge":            "lsa_age",
		"options":           "options",
		"lsaFlags":          "lsa_flags",
		"lsaType":           "lsa_type",
		"linkStateId":       "link_state_id",
		"advertisingRouter": "advertising_router",
		"lsaSeqNumber":      "lsa_seq_number",
		"checksum":          "checksum",
		"length":            "length",
		"networkMask":       "network_mask",
		"metricType":        "metric_type",
		"tos":               "tos",
		"metric":            "metric",
		"forwardAddress":    "forward_address",
		"externalRouteTag":  "external_route_tag",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := lsaData[origKey]; exists {
			transformed[newKey] = value
		}
	}

	return transformed
}

func transformNssaExternalLSA(lsaData map[string]interface{}, isNssa ...bool) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"lsaAge":             "lsa_age",
		"options":            "options",
		"lsaFlags":           "lsa_flags",
		"lsaType":            "lsa_type",
		"linkStateId":        "link_state_id",
		"advertisingRouter":  "advertising_router",
		"lsaSeqNumber":       "lsa_seq_number",
		"checksum":           "checksum",
		"length":             "length",
		"networkMask":        "network_mask",
		"metricType":         "metric_type",
		"tos":                "tos",
		"metric":             "metric",
		"nssaForwardAddress": "nssa_forward_address",
		"externalRouteTag":   "external_route_tag",
	}

	// TODO: What is this?
	//if isNssa != nil {
	//	fmt.Println("=========== TEST ===========")
	//}
	for jsonKey, protoKey := range fieldMapping {
		//if isNssa != nil {
		//	fmt.Println(jsonKey)
		//	fmt.Println(lsaData)
		//}
		if value, exists := lsaData[jsonKey]; exists {
			transformed[protoKey] = value
		}
	}
	//if isNssa != nil {
	//	fmt.Println("=========== TEST ===========")
	//}

	return transformed
}

func transformDatabaseRouterLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	addDatabaseLSABaseParameters(transformed, lsaData)
	if v, ok := lsaData["numOfRouterLinks"]; ok {
		transformed["num_of_router_links"] = v
	}

	return transformed
}

func transformDatabaseNetworkLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	addDatabaseLSABaseParameters(transformed, lsaData)
	return transformed
}

func transformDatabaseSummaryLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	addDatabaseLSABaseParameters(transformed, lsaData)
	transformed["summary_address"] = lsaData["summaryAddress"]
	return transformed
}

func transformDatabaseASBRSummaryLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	addDatabaseLSABaseParameters(transformed, lsaData)
	return transformed
}

func transformDatabaseNSSAExternalLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	addDatabaseLSABaseParameters(transformed, lsaData)
	if v, ok := lsaData["metricType"]; ok {
		transformed["metric_type"] = v
	}

	if v, ok := lsaData["route"]; ok {
		transformed["route"] = v
	}

	if v, ok := lsaData["tag"]; ok {
		transformed["tag"] = v
	}

	return transformed
}

func transformDatabaseExternalLSA(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})
	addDatabaseLSABaseParameters(transformed, lsaData)
	if v, ok := lsaData["metricType"]; ok {
		transformed["metric_type"] = v
	}
	if v, ok := lsaData["route"]; ok {
		transformed["route"] = v
	}
	if v, ok := lsaData["tag"]; ok {
		transformed["tag"] = v
	}
	return transformed
}

func addDatabaseLSABaseParameters(transformed, lsaData map[string]interface{}) {
	base := make(map[string]interface{})

	if v, ok := lsaData["lsId"]; ok {
		base["ls_id"] = v
	}
	if v, ok := lsaData["advertisedRouter"]; ok {
		base["advertised_router"] = v
	}
	if v, ok := lsaData["lsaAge"]; ok {
		base["lsa_age"] = v
	}
	if v, ok := lsaData["sequenceNumber"]; ok {
		base["sequence_number"] = v
	}
	if v, ok := lsaData["checksum"]; ok {
		base["checksum"] = v
	}

	transformed["base"] = base
}

func transformExternalLinkState(lsaData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"lsaAge":            "lsa_age",
		"options":           "options",
		"lsaFlags":          "lsa_flags",
		"lsaType":           "lsa_type",
		"linkStateId":       "link_state_id",
		"advertisingRouter": "advertising_router",
		"lsaSeqNumber":      "lsa_seq_number",
		"checksum":          "checksum",
		"length":            "length",
		"networkMask":       "network_mask",
		"metricType":        "metric_type",
		"tos":               "tos",
		"metric":            "metric",
		"forwardAddress":    "forward_address",
		"externalRouteTag":  "external_route_tag",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := lsaData[origKey]; exists {
			transformed[newKey] = value
		}
	}

	return transformed
}

func transformNeighbor(neighborData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"priority":                           "priority",
		"state":                              "state",
		"nbrPriority":                        "nbr_priority",
		"nbrState":                           "nbr_state",
		"converged":                          "converged",
		"role":                               "role",
		"upTimeInMsec":                       "up_time_in_msec",
		"deadTimeMsecs":                      "dead_time_msecs",
		"routerDeadIntervalTimerDueMsec":     "router_dead_interval_timer_due_msec",
		"upTime":                             "up_time",
		"deadTime":                           "dead_time",
		"address":                            "address",
		"ifaceAddress":                       "iface_address",
		"ifaceName":                          "iface_name",
		"retransmitCounter":                  "retransmit_counter",
		"linkStateRetransmissionListCounter": "link_state_retransmission_list_counter",
		"requestCounter":                     "request_counter",
		"linkStateRequestListCounter":        "link_state_request_list_counter",
		"dbSummaryCounter":                   "db_summary_counter",
		"databaseSummaryListCounter":         "database_summary_list_counter",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := neighborData[origKey]; exists {
			transformed[newKey] = value
		}
	}

	return transformed
}

// TODO: is this needed?
func transformSingleInterface(ifaceData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"administrativeStatus": "administrative_status",
		"operationalStatus":    "operational_status",
		"linkDetection":        "link_detection",
		"linkUps":              "link_ups",
		"linkDowns":            "link_downs",
		"lastLinkUp":           "last_link_up",
		"lastLinkDown":         "last_link_down",
		"vrfName":              "vrf_name",
		"mplsEnabled":          "mpls_enabled",
		"linkDown":             "link_down",
		"linkDownV6":           "link_down_v6",
		"mcForwardingV4":       "mc_forwarding_v4",
		"mcForwardingV6":       "mc_forwarding_v6",
		"pseudoInterface":      "pseudo_interface",
		"index":                "index",
		"metric":               "metric",
		"mtu":                  "mtu",
		"speed":                "speed",
		"flags":                "flags",
		"type":                 "type",
		"hardwareAddress":      "hardware_address",
		"interfaceType":        "interface_type",
		"interfaceSlaveType":   "interface_slave_type",
		"lacpBypass":           "lacp_bypass",
		"protodown":            "protodown",
		"parentIfindex":        "parent_ifindex",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := ifaceData[origKey]; exists {
			transformed[newKey] = value
		}
	}

	if ipAddrs, ok := ifaceData["ipAddresses"].([]interface{}); ok {
		transformedIps := make([]interface{}, len(ipAddrs))
		for i, ip := range ipAddrs {
			ipMap := ip.(map[string]interface{})
			transformedIps[i] = map[string]interface{}{
				"address":    ipMap["address"],
				"secondary":  ipMap["secondary"],
				"unnumbered": ipMap["unnumbered"],
			}
		}
		transformed["ip_addresses"] = transformedIps
	}

	if evpnMh, ok := ifaceData["evpnMh"].(map[string]interface{}); ok {
		transformed["evpn_mh"] = transformEvpnMh(evpnMh)
	}

	return transformed
}

func transformEvpnMh(evpnData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"ethernetSegmentId": "ethernet_segment_id",
		"esi":               "esi",
		"dfPreference":      "df_preference",
		"dfAlgorithm":       "df_algorithm",
		"dfStatus":          "df_status",
		"multihomingMode":   "multihoming_mode",
		"activeMode":        "active_mode",
		"bypassMode":        "bypass_mode",
		"localBias":         "local_bias",
		"fastFailover":      "fast_failover",
		"upTime":            "up_time",
		"bgpStatus":         "bgp_status",
		"protocolStatus":    "protocol_status",
		"protocolDown":      "protocol_down",
		"macCount":          "mac_count",
		"localIfindex":      "local_ifindex",
		"networkCount":      "network_count",
		"joinCount":         "join_count",
		"leaveCount":        "leave_count",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := evpnData[origKey]; exists {
			transformed[newKey] = value
		}
	}

	return transformed
}

// TODO: is this needed?
func transformRoute(routeData map[string]interface{}) map[string]interface{} {
	transformed := make(map[string]interface{})

	fieldMapping := map[string]string{
		"prefix":                   "prefix",
		"prefixLen":                "prefix_len",
		"protocol":                 "protocol",
		"vrfId":                    "vrf_id",
		"vrfName":                  "vrf_name",
		"selected":                 "selected",
		"destSelected":             "dest_selected",
		"distance":                 "distance",
		"metric":                   "metric",
		"installed":                "installed",
		"table":                    "table",
		"internalStatus":           "internal_status",
		"internalFlags":            "internal_flags",
		"internalNextHopNum":       "internal_next_hop_num",
		"internalNextHopActiveNum": "internal_next_hop_active_num",
		"nexthopGroupID":           "nexthop_group_id",
		"installedNexthopGroupID":  "installed_nexthop_group_id",
		"uptime":                   "uptime",
	}

	for origKey, newKey := range fieldMapping {
		if value, exists := routeData[origKey]; exists {
			transformed[newKey] = value
		}
	}

	if nexthops, ok := routeData["nexthops"].([]interface{}); ok {
		transformedNexthops := make([]interface{}, len(nexthops))
		for i, nh := range nexthops {
			nhMap := nh.(map[string]interface{})
			transformedNexthops[i] = map[string]interface{}{
				"flags":              nhMap["flags"],
				"fib":                nhMap["fib"],
				"directly_connected": nhMap["directlyConnected"],
				"duplicate":          nhMap["duplicate"],
				"ip":                 nhMap["ip"],
				"afi":                nhMap["afi"],
				"interface_index":    nhMap["interfaceIndex"],
				"interface_name":     nhMap["interfaceName"],
				"active":             nhMap["active"],
				"weight":             nhMap["weight"],
			}
		}
		transformed["nexthops"] = transformedNexthops
	}

	return transformed
}

func getString(m map[string]interface{}, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}

func getBool(m map[string]interface{}, key string) bool {
	if val, ok := m[key].(bool); ok {
		return val
	}
	return false
}

func getFloat(m map[string]interface{}, key string) float64 {
	if val, ok := m[key].(float64); ok {
		return val
	}
	return 0
}
