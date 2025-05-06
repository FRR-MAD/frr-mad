package analyzer_test

import (
	"fmt"
	"testing"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/analyzer"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestRouterLsa1(t *testing.T) {

	ana := initAnalyzer()

	frrMetrics := getR101FRRdata()

	expectedAccessList := map[string]frrProto.AccessListAnalyzer{
		"localsite": {
			AccessList: "localsite",
			AclEntry: []*frrProto.ACLEntry{
				{
					IPAddress:    "192.168.1.0",
					PrefixLength: 24,
					IsPermit:     true,
					Sequence:     15,
				},
			},
		},
		"term": {
			AccessList: "term",
			AclEntry: []*frrProto.ACLEntry{
				{
					IPAddress:    "127.0.0.1",
					PrefixLength: 32,
					IsPermit:     true,
					Sequence:     5,
				},
				{
					IPAddress:    "any",
					PrefixLength: 0,
					IsPermit:     false,
					Any:          true,
					Sequence:     10,
				},
			},
		},
	}

	expectedRuntimeRouterLSDB := frrProto.InterAreaLsa{
		RouterId: "65.0.1.1",
		Hostname: "r101",
		Areas: []*frrProto.AreaAnalyzer{
			{
				AreaName: "0.0.0.0",
				LsaType:  "router-LSA",
				Links: []*frrProto.Advertisement{
					{
						InterfaceAddress: "10.0.12.1",
						LinkType:         "a Transit Network",
					},
					{
						InterfaceAddress: "10.0.2.0",
						PrefixLength:     "24",
						LinkType:         "Stub Network",
					},
					{
						InterfaceAddress: "10.0.14.1",
						LinkType:         "a Transit Network",
					},
					{
						InterfaceAddress: "10.0.16.1",
						LinkType:         "a Transit Network",
					},
					{
						InterfaceAddress: "10.0.18.1",
						LinkType:         "a Transit Network",
					},
					{
						InterfaceAddress: "10.0.15.1",
						LinkType:         "a Transit Network",
					},
					{
						InterfaceAddress: "10.0.0.0",
						PrefixLength:     "23",
						LinkType:         "Stub Network",
					},
					{
						InterfaceAddress: "10.0.17.1",
						LinkType:         "a Transit Network",
					},
					{
						InterfaceAddress: "10.0.13.1",
						LinkType:         "a Transit Network",
					},
					{
						InterfaceAddress: "10.0.19.1",
						LinkType:         "a Transit Network",
					},
				},
			},
		},
	}

	expectedPredictedRouterLSDB := frrProto.IntraAreaLsa{
		Hostname: "r101",
		RouterId: "65.0.1.1",
		Areas: []*frrProto.AreaAnalyzer{
			{
				AreaName: "0.0.0.0",    //  string
				LsaType:  "router-LSA", //     string
				AreaType: "normal",     //     string
				Links: []*frrProto.Advertisement{
					{
						InterfaceAddress: "10.0.12.1",
						PrefixLength:     "24",
						LinkType:         "transit network",
					},
					{
						InterfaceAddress: "10.0.2.0",
						PrefixLength:     "24",
						LinkType:         "stub network",
					},
					{
						InterfaceAddress: "10.0.13.1",
						PrefixLength:     "24",
						LinkType:         "transit network",
					},
					{
						InterfaceAddress: "10.0.0.0",
						PrefixLength:     "23",
						LinkType:         "stub network",
					},
					{
						InterfaceAddress: "10.0.14.1",
						PrefixLength:     "24",
						LinkType:         "transit network",
					},
					{
						InterfaceAddress: "10.0.15.1",
						PrefixLength:     "24",
						LinkType:         "transit network",
					},
					{
						InterfaceAddress: "10.0.16.1",
						PrefixLength:     "24",
						LinkType:         "transit network",
					},
					{
						InterfaceAddress: "10.0.17.1",
						PrefixLength:     "24",
						LinkType:         "transit network",
					},
					{
						InterfaceAddress: "10.0.18.1",
						PrefixLength:     "24",
						LinkType:         "transit network",
					},
					{
						InterfaceAddress: "10.0.19.1",
						PrefixLength:     "24",
						LinkType:         "transit network",
					},
				},
			},
		},
	}

	//fmt.Println(expectedPredictedRouterLSDB)

	actualAccessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)

	less := func(a, b string) bool { return a < b }

	var actualAccessListKeys []string
	for k, _ := range actualAccessList {
		actualAccessListKeys = append(actualAccessListKeys, k)
	}
	var expectedAccessListKeys []string
	for k, _ := range actualAccessList {
		expectedAccessListKeys = append(expectedAccessListKeys, k)
	}

	expectedStaticList := map[string]*frrProto.StaticList{
		"192.168.1.0": {
			IpAddress:    "192.168.1.0",
			PrefixLength: 24,
			NextHop:      "192.168.100.91",
		},
	}

	actualStaticList := analyzer.GetStaticRouteList(frrMetrics.StaticFrrConfiguration, actualAccessList)

	var actualStaticListKeys []string
	for k, _ := range actualStaticList {
		actualStaticListKeys = append(actualStaticListKeys, k)
	}
	var expectedStaticListKeys []string
	for k, _ := range actualStaticList {
		expectedStaticListKeys = append(expectedStaticListKeys, k)
	}

	//t.Logf("%v\n", expectedStaticList)
	//t.Logf("%v\n", actualStaticList)

	// Runtime parsing of router

	actualRuntimeRouterLSDB := analyzer.GetRuntimeRouterData(frrMetrics.OspfRouterData, frrMetrics.StaticFrrConfiguration.Hostname)

	expectedRuntimeRouterLSDBAreaLength := len(expectedRuntimeRouterLSDB.Areas)
	actualRuntimeRouterLSDBAreaLength := len(actualRuntimeRouterLSDB.Areas)

	expectedRuntimeRouterLSDBAreaMap := make(map[string][]*frrProto.Advertisement)
	expectedRuntimeRouterLSDBAreas := []string{}
	expectedRuntimeRouterLSDBLsaType := map[string]string{}
	for _, area := range expectedRuntimeRouterLSDB.Areas {
		expectedRuntimeRouterLSDBAreas = append(expectedRuntimeRouterLSDBAreas, area.AreaName)
		expectedRuntimeRouterLSDBLsaType[area.AreaName] = area.LsaType
		for _, link := range area.Links {
			expectedRuntimeRouterLSDBAreaMap[area.AreaName] = append(expectedRuntimeRouterLSDBAreaMap[area.AreaName], link)
		}
	}

	actualRuntimeRouterLSDBAreaMap := make(map[string][]*frrProto.Advertisement)
	actualRuntimeRouterLSDBAreas := []string{}
	actualRuntimeRouterLSDBLsaType := map[string]string{}
	for _, area := range actualRuntimeRouterLSDB.Areas {
		actualRuntimeRouterLSDBAreas = append(actualRuntimeRouterLSDBAreas, area.AreaName)
		actualRuntimeRouterLSDBLsaType[area.AreaName] = area.LsaType
		for _, link := range area.Links {
			actualRuntimeRouterLSDBAreaMap[area.AreaName] = append(actualRuntimeRouterLSDBAreaMap[area.AreaName], link)
		}
	}

	t.Run("TestRuntimeRouterParsing", func(t *testing.T) {
		assert.Equal(t, expectedRuntimeRouterLSDB.RouterId, actualRuntimeRouterLSDB.RouterId)
		assert.Equal(t, expectedRuntimeRouterLSDB.Hostname, actualRuntimeRouterLSDB.Hostname)
		assert.Equal(t, expectedRuntimeRouterLSDBAreaLength, actualRuntimeRouterLSDBAreaLength)
		assert.True(t, cmp.Diff(expectedRuntimeRouterLSDBAreas, actualRuntimeRouterLSDBAreas, cmpopts.SortSlices(less)) == "")
		for _, key := range expectedRuntimeRouterLSDBAreas {
			assert.Equal(t, expectedRuntimeRouterLSDBLsaType[key], actualRuntimeRouterLSDBLsaType[key])
		}
		//for _, key := range expectedRuntimeRouterLSDBAreas {
		//	t.Logf("%v\n", actualRuntimeRouterLSDBAreaMap[key])

		//}

	})

	// test should state Static File Router Data
	//expectedPredictedRouterLSDB := &frrProto.IntraAreaLsa{}

	isNssa, actualPredictedRouterLSDB := analyzer.GetStaticFileRouterData(frrMetrics.StaticFrrConfiguration)

	// Write Router Testing now, because parsing of static config, router config, static list and access list is successful

	ana.RouterAnomalyAnalysis(actualAccessList, actualPredictedRouterLSDB, actualRuntimeRouterLSDB)

	t.Run("TestRouterAdvertisment", func(t *testing.T) {
		//ana.RouterAnomalyAnalysis(actualAccessList, )
	})

	//analyzer.RouterAnomalyAnalysis(accessList, shouldState, isState)
	t.Run("TestGetAccessList", func(t *testing.T) {
		assert.Equal(t, len(expectedAccessList), len(actualAccessList))
		assert.True(t, cmp.Diff(expectedAccessListKeys, actualAccessListKeys, cmpopts.SortSlices(less)) == "")
		for _, v := range actualAccessListKeys {
			assert.Equal(t, actualAccessList[v], expectedAccessList[v])
		}
	})

	t.Run("TestGetStaticRouteList", func(t *testing.T) {
		assert.Equal(t, len(expectedStaticList), len(actualStaticList))
		assert.Equal(t, expectedStaticListKeys, actualStaticListKeys)
		for _, v := range actualStaticListKeys {
			assert.Equal(t, actualStaticList[v], expectedStaticList[v])
		}
	})

	// test GetStaticFileRouterData
	expectedPredictedRouterLSDBAreas := []string{} //done
	actualPredictedRouterLSDBAreas := []string{}

	expectedPredictedRouterLSDBLsaTypePerArea := make(map[string]string) //done
	actualPredictedRouterLSDBLsaTypePerArea := make(map[string]string)

	expectedPredictedRouterLSDBIntPerArea := make(map[string][]*frrProto.Advertisement)
	actualPredictedRouterLSDBIntPerArea := make(map[string][]*frrProto.Advertisement)

	for _, area := range expectedPredictedRouterLSDB.Areas {
		tmp := []*frrProto.Advertisement{}
		expectedPredictedRouterLSDBAreas = append(expectedPredictedRouterLSDBAreas, area.AreaName)
		expectedPredictedRouterLSDBLsaTypePerArea[area.GetAreaName()] = area.LsaType
		for _, iface := range area.Links {
			tmp = append(tmp, iface)
		}
		expectedPredictedRouterLSDBIntPerArea[area.AreaName] = tmp
	}

	for _, area := range actualPredictedRouterLSDB.Areas {
		tmp := []*frrProto.Advertisement{}
		actualPredictedRouterLSDBAreas = append(actualPredictedRouterLSDBAreas, area.AreaName)
		actualPredictedRouterLSDBLsaTypePerArea[area.GetAreaName()] = area.LsaType
		for _, iface := range area.Links {
			tmp = append(tmp, iface)
		}
		actualPredictedRouterLSDBIntPerArea[area.AreaName] = tmp
	}

	t.Run("TestStaticFileRouterDataFunction", func(t *testing.T) {
		assert.False(t, isNssa)
		assert.Equal(t, expectedPredictedRouterLSDB.Hostname, actualPredictedRouterLSDB.Hostname)
		assert.Equal(t, expectedPredictedRouterLSDB.RouterId, actualPredictedRouterLSDB.RouterId)
		assert.Equal(t, len(expectedPredictedRouterLSDBAreas), len(actualPredictedRouterLSDBAreas))
		assert.True(t, cmp.Diff(expectedPredictedRouterLSDBAreas, actualPredictedRouterLSDBAreas, cmpopts.SortSlices(less)) == "")

		for _, key := range expectedPredictedRouterLSDBAreas {
			assert.Equal(t, expectedPredictedRouterLSDBLsaTypePerArea[key], actualPredictedRouterLSDBLsaTypePerArea[key])
		}
		// expectedPredictedRouterLSDBIntPerArea := make(map[string][]*frrProto.Advertisement)
		// actualPredictedRouterLSDBIntPerArea := make(map[string][]*frrProto.Advertisement)

		for _, key := range expectedPredictedRouterLSDBAreas {
			assert.Equal(t, expectedPredictedRouterLSDBLsaTypePerArea[key], actualPredictedRouterLSDBLsaTypePerArea[key])
		}

		for _, key := range expectedPredictedRouterLSDBAreas {
			expectedIfaceMap, expectedIfaceList := getIfaceMap(expectedPredictedRouterLSDBIntPerArea[key])
			actualIfaceMap, actualIfaceList := getIfaceMap(actualPredictedRouterLSDBIntPerArea[key])
			assert.Equal(t, len(expectedIfaceList), len(actualIfaceList))
			assert.True(t, cmp.Diff(expectedIfaceList, actualIfaceList, cmpopts.SortSlices(less)) == "")

			for _, ifaceKey := range expectedIfaceList {
				assert.Equal(t, expectedIfaceMap[ifaceKey], actualIfaceMap[ifaceKey])
			}
		}
	})

	// test GetStaticFileExternalData

	// test GetStaticFileNssaExternalData

	t.Run("TestStaticListEqualACL", func(t *testing.T) {

	})

	t.Run("TestStaticRoute", func(t *testing.T) {

	})

}

// what do I need?
/*
	- frrMetrics.OspfExternalData
	- frrMetrics.StaticFrrConfiguration
	- accessList
	- predictedExternalLSDB
	- runtimeExternalLSDB

	runtimeExternalLSDB := GetRuntimeExternalRouterData(c.metrics.OspfExternalData, c.metrics.StaticFrrConfiguration.Hostname)
	predictedExternalLSDB := getStaticFileExternalData(c.metrics.StaticFrrConfiguration)

	if len(staticRouteMap) > 0 || isNssa {
		ExternalAnomalyAnalysis(accessList, predictedExternalLSDB, runtimeExternalLSDB)
	}
*/
// what will be tested?
/*
	- actualExternalLSDB -> GetRuntimeExternalRouterData
	- predictedExternalLSDB

	- runtimeExternalLSDB
*/
func TestExternalLsa1(t *testing.T) {

	less := func(a, b string) bool { return a < b }
	frrMetrics := getR101FRRdata()
	accessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)
	staticList := analyzer.GetStaticRouteList(frrMetrics.StaticFrrConfiguration, accessList)

	expectedPredictedExternalLSDB := &frrProto.InterAreaLsa{
		Hostname: frrMetrics.StaticFrrConfiguration.Hostname,
		RouterId: frrMetrics.StaticFrrConfiguration.OspfConfig.RouterId,
		Areas: []*frrProto.AreaAnalyzer{

			{
				//AreaName: "0.0.0.0",
				LsaType: "AS-external-LSA",
				//AreaType: "",
				Links: []*frrProto.Advertisement{
					{
						LinkStateId:  "192.168.1.0", //   str
						PrefixLength: "24",          //  str
						LinkType:     "external",    // str
					},
				},
			},
		},
	}
	actualPredictedExternalLSDB := analyzer.GetStaticFileExternalData(frrMetrics.StaticFrrConfiguration, accessList, staticList)

	// - actualExternalLSDB -> GetRuntimeExternalRouterData
	// - predictedExternalLSDB
	t.Run("TestExternalDataStaticShouldAndIs", func(t *testing.T) {

		assert.Equal(t, expectedPredictedExternalLSDB.Hostname, actualPredictedExternalLSDB.Hostname)
		assert.Equal(t, expectedPredictedExternalLSDB.RouterId, actualPredictedExternalLSDB.RouterId)
		expectedAreaMapTmp := make(map[string][]string)
		expectedAreaListTmp := []string{}
		for _, area := range expectedPredictedExternalLSDB.Areas {
			expectedAreaListTmp = append(expectedAreaListTmp, area.AreaName)
			for _, adv := range area.Links {
				expectedAreaMapTmp[area.AreaName] = append(expectedAreaMapTmp[area.AreaName], adv.LinkStateId)
			}
		}
		actualAreaMapTmp := make(map[string][]string)
		actualAreaListTmp := []string{}
		for _, area := range expectedPredictedExternalLSDB.Areas {
			actualAreaListTmp = append(actualAreaListTmp, area.AreaName)
			for _, adv := range area.Links {
				actualAreaMapTmp[area.AreaName] = append(actualAreaMapTmp[area.AreaName], adv.LinkStateId)
			}
		}

		assert.Equal(t, len(expectedAreaListTmp), len(actualAreaListTmp))
		for _, value := range expectedAreaListTmp {
			assert.True(t, cmp.Diff(expectedAreaMapTmp[value], actualAreaMapTmp[value], cmpopts.SortSlices(less)) == "")
		}
	})

	//- runtimeExternalLSDB
	expectedRuntimeExternalLSDB := &frrProto.InterAreaLsa{
		Hostname: "r101",
		RouterId: "65.0.1.1",
		Areas: []*frrProto.AreaAnalyzer{
			{
				LsaType: "AS-external-LSA",
				Links: []*frrProto.Advertisement{
					{
						LinkStateId:  "192.168.1.0",
						PrefixLength: "24",
						LinkType:     "external",
					},
				},
			},
		},
	}
	actualRuntimeExternalLSDB := analyzer.GetRuntimeExternalData(frrMetrics.OspfExternalData, staticList, frrMetrics.StaticFrrConfiguration.Hostname)

	// TODO: maybe add AreaName testing? For that area assignment needs to be done. It doesn't seem too easy and it's not really necessary. Considering that static and connected redistributions happen via LSA Type 5 anyway and if it's connected to an NSSA it will still show a type 5 lsa but in type 7 lsa testing it will correctly show the correct static and connected redistributions.
	t.Run("TestExternalDataRuntimeShouldAndIs", func(t *testing.T) {
		assert.Equal(t, expectedRuntimeExternalLSDB.Hostname, actualRuntimeExternalLSDB.Hostname)
		assert.Equal(t, expectedRuntimeExternalLSDB.RouterId, actualRuntimeExternalLSDB.RouterId)

		assert.Equal(t, len(expectedRuntimeExternalLSDB.Areas), len(actualRuntimeExternalLSDB.Areas))
		expectedTotalLinks := 0
		actualTotalLinks := 0

		for _, area := range expectedRuntimeExternalLSDB.Areas {
			expectedTotalLinks += len(area.Links)
		}
		for _, area := range actualRuntimeExternalLSDB.Areas {
			actualTotalLinks += len(area.Links)
		}

		assert.Equal(t, expectedTotalLinks, actualTotalLinks)

		// Create maps with LinkStateId as keys for comparison
		expectedTmp := map[string][]*frrProto.Advertisement{}
		actualTmp := map[string][]*frrProto.Advertisement{}

		// Populate the map for expected data
		for _, area := range expectedRuntimeExternalLSDB.Areas {
			for _, link := range area.Links {
				expectedTmp[link.LinkStateId] = append(expectedTmp[link.LinkStateId], link)
			}
		}

		// Populate the map for actual data
		for _, area := range actualRuntimeExternalLSDB.Areas {
			for _, link := range area.Links {
				actualTmp[link.LinkStateId] = append(actualTmp[link.LinkStateId], link)
			}
		}

		// Assert that both maps have the same keys
		assert.Equal(t, len(expectedTmp), len(actualTmp), "Expected and actual maps should have the same number of LinkStateIds")

		// Assert that for each key, both maps have the same advertisements
		for linkStateId, expectedAdvs := range expectedTmp {
			actualAdvs, exists := actualTmp[linkStateId]
			assert.True(t, exists, "LinkStateId %s should exist in actual data", linkStateId)
			assert.Equal(t, len(expectedAdvs), len(actualAdvs), "Expected and actual should have same number of advertisements for LinkStateId %s", linkStateId)

			// Additional assertions could be added here to compare specific fields of each advertisement
			// Create maps to compare advertisements by PrefixLength and LinkType
			for _, expectedAdv := range expectedAdvs {
				foundMatch := false
				for _, actualAdv := range actualAdvs {
					if expectedAdv.PrefixLength == actualAdv.PrefixLength &&
						expectedAdv.LinkType == actualAdv.LinkType {
						foundMatch = true
						break
					}
				}
				assert.True(t, foundMatch, "No matching advertisement found for LinkStateId %s with PrefixLength %s and LinkType %s",
					linkStateId, expectedAdv.PrefixLength, expectedAdv.LinkType)
			}

		}
	})
}

func TestAnomalyAnalysis1(t *testing.T) {

	ana := initAnalyzer()
	frrMetrics := getR101FRRdata()
	accessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)
	staticRouteMap := analyzer.GetStaticRouteList(frrMetrics.StaticFrrConfiguration, accessList)

	runtimeRouterLSDB := analyzer.GetRuntimeRouterData(frrMetrics.OspfRouterData, frrMetrics.StaticFrrConfiguration.Hostname)

	_, predictedRouterLSDB := analyzer.GetStaticFileRouterData(frrMetrics.StaticFrrConfiguration)
	ana.RouterAnomalyAnalysis(accessList, predictedRouterLSDB, runtimeRouterLSDB)

	t.Run("TestRouterLSAAnomalyTesting", func(t *testing.T) {
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasUnderAdvertisedPrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasDuplicatePrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasMisconfiguredPrefixes)
		assert.Empty(t, ana.AnalysisResult.RouterAnomaly.MissingEntries)
		assert.Empty(t, ana.AnalysisResult.RouterAnomaly.SuperfluousEntries)
		assert.Empty(t, ana.AnalysisResult.RouterAnomaly.DuplicateEntries)
	})

	//

	predictedExternalLSDB := analyzer.GetStaticFileExternalData(frrMetrics.StaticFrrConfiguration, accessList, staticRouteMap)
	runtimeExternalLSDB := analyzer.GetRuntimeExternalData(frrMetrics.OspfExternalData, staticRouteMap, frrMetrics.StaticFrrConfiguration.Hostname)
	ana.ExternalAnomalyAnalysis(predictedExternalLSDB, runtimeExternalLSDB)

	t.Run("TestExternalLSAAnomalyTesting", func(t *testing.T) {
		assert.False(t, ana.AnalysisResult.ExternalAnomaly.HasOverAdvertisedPrefixes)
		assert.False(t, ana.AnalysisResult.ExternalAnomaly.HasUnderAdvertisedPrefixes)
		assert.False(t, ana.AnalysisResult.ExternalAnomaly.HasDuplicatePrefixes)
		assert.False(t, ana.AnalysisResult.ExternalAnomaly.HasMisconfiguredPrefixes)
		assert.Empty(t, ana.AnalysisResult.ExternalAnomaly.MissingEntries)
		assert.Empty(t, ana.AnalysisResult.ExternalAnomaly.SuperfluousEntries)
		assert.Empty(t, ana.AnalysisResult.ExternalAnomaly.DuplicateEntries)
	})

}

func TestAnomalyAnalysisLsaFive1(t *testing.T) {

}

// TODO: TestNssaExternalLsa1

// func TestGetAccessList2(t *testing.T) {

// 	frrMetrics := getR102FRRdata()

// 	expectedResult := map[string]frrProto.AccessListAnalyzer{
// 		"localsite": {
// 			AccessList: "localsite",
// 			AclEntry: []*frrProto.ACLEntry{
// 				{
// 					IPAddress:    "192.168.11.0",
// 					PrefixLength: 24,
// 					IsPermit:     true,
// 					Sequence:     15,
// 				},
// 			},
// 		},
// 		"term": {
// 			AccessList: "term",
// 			AclEntry: []*frrProto.ACLEntry{
// 				{
// 					IPAddress:    "127.0.0.1",
// 					PrefixLength: 32,
// 					IsPermit:     true,
// 					Sequence:     5,
// 				},
// 				{
// 					IPAddress:    "any",
// 					PrefixLength: 0,
// 					IsPermit:     true,
// 					Any:          true,
// 					Sequence:     10,
// 				},
// 			},
// 		},
// 	}

// 	result := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)

// 	//assert.Equal(t, expectedResult, result)
// 	t.Log()
// 	t.Logf("%v\n", expectedResult)
// 	t.Log()
// 	t.Logf("%v\n", result)

// }

// Add these test cases to your analyzer_test.go file

func TestNssaExternalLsaHappy1(t *testing.T) {
	// Setup test data for NSSA-External analysis
	ana := initAnalyzer()
	frrMetrics := getNssaRouterFRRdataHappy1()

	accessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)
	staticRouteMap := analyzer.GetStaticRouteList(frrMetrics.StaticFrrConfiguration, accessList)

	// Get predicted and runtime NSSA-external LSDBs
	predictedNssaExternalLSDB := analyzer.GetStaticFileNssaExternalData(frrMetrics.StaticFrrConfiguration)
	runtimeNssaExternalLSDB := analyzer.GetNssaExternalData(frrMetrics.OspfNssaExternalData, staticRouteMap, frrMetrics.StaticFrrConfiguration.Hostname)

	// Run the analysis
	ana.NssaExternalAnomalyAnalysis(accessList, predictedNssaExternalLSDB, runtimeNssaExternalLSDB)

	t.Run("TestNssaExternalNormalCase", func(t *testing.T) {
		// In normal case, there should be no anomalies
		assert.False(t, ana.AnalysisResult.NssaExternalAnomaly.HasOverAdvertisedPrefixes)
		assert.False(t, ana.AnalysisResult.NssaExternalAnomaly.HasUnderAdvertisedPrefixes)
		assert.False(t, ana.AnalysisResult.NssaExternalAnomaly.HasDuplicatePrefixes)
		assert.Empty(t, ana.AnalysisResult.NssaExternalAnomaly.MissingEntries)
		assert.Empty(t, ana.AnalysisResult.NssaExternalAnomaly.SuperfluousEntries)
		assert.Empty(t, ana.AnalysisResult.NssaExternalAnomaly.DuplicateEntries)
	})
}

func TestNssaExternalAnomaliesUnhappy1(t *testing.T) {
	// Setup test data with intentional anomalies
	ana := initAnalyzer()
	frrMetrics := getNssaRouterFRRdataUnhappy1()

	accessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)
	staticRouteMap := analyzer.GetStaticRouteList(frrMetrics.StaticFrrConfiguration, accessList)

	// Get predicted and runtime NSSA-external LSDBs
	predictedNssaExternalLSDB := analyzer.GetStaticFileNssaExternalData(frrMetrics.StaticFrrConfiguration)
	runtimeNssaExternalLSDB := analyzer.GetNssaExternalData(frrMetrics.OspfNssaExternalData, staticRouteMap, frrMetrics.StaticFrrConfiguration.Hostname)

	fmt.Println("---------------------- Predicted ----------------------")
	fmt.Println(predictedNssaExternalLSDB)
	fmt.Println("---------------------- Predicted ----------------------")

	fmt.Println("---------------------- Runtime ----------------------")
	fmt.Println(runtimeNssaExternalLSDB)
	fmt.Println("---------------------- Runtime ----------------------")

	// Run the analysis
	ana.NssaExternalAnomalyAnalysis(accessList, predictedNssaExternalLSDB, runtimeNssaExternalLSDB)

	t.Run("TestNssaExternalMissingRoutes", func(t *testing.T) {
		// Should detect missing routes that should be advertised
		assert.True(t, ana.AnalysisResult.NssaExternalAnomaly.HasUnderAdvertisedPrefixes)
		assert.NotEmpty(t, ana.AnalysisResult.NssaExternalAnomaly.MissingEntries)
	})

	// t.Run("TestNssaExternalExtraRoutes", func(t *testing.T) {
	// 	// Should detect extra routes that shouldn't be advertised
	// 	assert.True(t, ana.AnalysisResult.NssaExternalAnomaly.HasOverAdvertisedPrefixes)
	// 	assert.NotEmpty(t, ana.AnalysisResult.NssaExternalAnomaly.SuperfluousEntries)
	// })

	// t.Run("TestNssaExternalDuplicates", func(t *testing.T) {
	// 	// Should detect duplicate routes
	// 	assert.True(t, ana.AnalysisResult.NssaExternalAnomaly.HasDuplicatePrefixes)
	// 	assert.NotEmpty(t, ana.AnalysisResult.NssaExternalAnomaly.DuplicateEntries)
	// })
}
