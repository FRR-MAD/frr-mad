package analyzer_test

import (
	"strings"
	"testing"

	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/analyzer"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestRouterLsaHappy1(t *testing.T) {

	ana := initAnalyzer()
	frrMetrics := getR101FRRdata()
	expectedAccessList := getExpectedAccessListr101Happy()
	expectedIsRouterLSDB := getExpectedIsRouterLSDBr101Happy()
	expectedShouldRouterLSDB := getExpectedShouldRouterLSDBr101Happy()
	actualAccessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)
	less := func(a, b string) bool { return a < b }

	var actualAccessListKeys []string
	for k := range actualAccessList {
		actualAccessListKeys = append(actualAccessListKeys, k)
	}
	var expectedAccessListKeys []string
	for k := range actualAccessList {
		expectedAccessListKeys = append(expectedAccessListKeys, k)
	}

	expectedStaticList := getExpectedStaticListr101Happy()
	actualStaticList := analyzer.GetStaticRouteList(frrMetrics.StaticFrrConfiguration, actualAccessList)

	var actualStaticListKeys []string
	for k := range actualStaticList {
		actualStaticListKeys = append(actualStaticListKeys, k)
	}
	var expectedStaticListKeys []string
	for k := range actualStaticList {
		expectedStaticListKeys = append(expectedStaticListKeys, k)
	}

	actualRuntimeRouterLSDB := analyzer.GetRuntimeRouterData(frrMetrics.OspfRouterData, frrMetrics.StaticFrrConfiguration.Hostname)
	expectedRuntimeRouterLSDBAreaLength := len(expectedIsRouterLSDB.Areas)
	actualRuntimeRouterLSDBAreaLength := len(actualRuntimeRouterLSDB.Areas)
	expectedRuntimeRouterLSDBAreaMap := make(map[string][]*frrProto.Advertisement)
	expectedIsRouterLSDBAreas := []string{}
	expectedRuntimeRouterLSDBLsaType := map[string]string{}
	for _, area := range expectedIsRouterLSDB.Areas {
		expectedIsRouterLSDBAreas = append(expectedIsRouterLSDBAreas, area.AreaName)
		expectedRuntimeRouterLSDBLsaType[area.AreaName] = area.LsaType
		expectedRuntimeRouterLSDBAreaMap[area.AreaName] = append(expectedRuntimeRouterLSDBAreaMap[area.AreaName], area.Links...)
	}

	actualIsRouterLSDBAreaMap := make(map[string][]*frrProto.Advertisement)
	actualRuntimeRouterLSDBAreas := []string{}
	actualRuntimeRouterLSDBLsaType := map[string]string{}
	for _, area := range actualRuntimeRouterLSDB.Areas {
		actualRuntimeRouterLSDBAreas = append(actualRuntimeRouterLSDBAreas, area.AreaName)
		actualRuntimeRouterLSDBLsaType[area.AreaName] = area.LsaType
		actualIsRouterLSDBAreaMap[area.AreaName] = append(actualIsRouterLSDBAreaMap[area.AreaName], area.Links...)
	}

	t.Run("TestRuntimeRouterParsing", func(t *testing.T) {
		assert.Equal(t, expectedIsRouterLSDB.RouterId, actualRuntimeRouterLSDB.RouterId)
		assert.Equal(t, expectedIsRouterLSDB.Hostname, actualRuntimeRouterLSDB.Hostname)
		assert.Equal(t, expectedRuntimeRouterLSDBAreaLength, actualRuntimeRouterLSDBAreaLength)
		assert.True(t, cmp.Diff(expectedIsRouterLSDBAreas, actualRuntimeRouterLSDBAreas, cmpopts.SortSlices(less)) == "")
		for _, key := range expectedIsRouterLSDBAreas {
			assert.Equal(t, expectedRuntimeRouterLSDBLsaType[key], actualRuntimeRouterLSDBLsaType[key])
		}
	})

	// test should state Static File Router Data
	//expectedPredictedRouterLSDB := &frrProto.IntraAreaLsa{}

	isNssa, actualPredictedRouterLSDB := analyzer.GetStaticFileRouterData(frrMetrics.StaticFrrConfiguration)

	// Write Router Testing now, because parsing of static config, router config, static list and access list is successful

	ana.RouterAnomalyAnalysisLSDB(actualAccessList, actualPredictedRouterLSDB, actualRuntimeRouterLSDB)

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

	for _, area := range expectedShouldRouterLSDB.Areas {
		tmp := []*frrProto.Advertisement{}
		expectedPredictedRouterLSDBAreas = append(expectedPredictedRouterLSDBAreas, area.AreaName)
		expectedPredictedRouterLSDBLsaTypePerArea[area.GetAreaName()] = area.LsaType
		tmp = append(tmp, area.Links...)
		expectedPredictedRouterLSDBIntPerArea[area.AreaName] = tmp
	}

	for _, area := range actualPredictedRouterLSDB.Areas {
		tmp := []*frrProto.Advertisement{}
		actualPredictedRouterLSDBAreas = append(actualPredictedRouterLSDBAreas, area.AreaName)
		actualPredictedRouterLSDBLsaTypePerArea[area.GetAreaName()] = area.LsaType
		tmp = append(tmp, area.Links...)
		actualPredictedRouterLSDBIntPerArea[area.AreaName] = tmp
	}

	t.Run("TestStaticFileRouterDataFunction", func(t *testing.T) {
		assert.False(t, isNssa)
		assert.Equal(t, expectedShouldRouterLSDB.Hostname, actualPredictedRouterLSDB.Hostname)
		assert.Equal(t, expectedShouldRouterLSDB.RouterId, actualPredictedRouterLSDB.RouterId)
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

func TestRouterLsaUnhappy1(t *testing.T) {

	ana := initAnalyzer()
	frrMetrics := getR101FRRdata()
	accessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)
	isRouterLSDB := getExpectedShouldRouterLSDBr101MissingEntries()
	shouldRouterLSDB := getExpectedIsRouterLSDBr101Happy()
	ana.RouterAnomalyAnalysisLSDB(accessList, &shouldRouterLSDB, isRouterLSDB)
	// unhappy c.ExternalAnomalyAnalysisLSDB(shouldExternalLSDB, isExternalLSDB)
	expectedMissingEntrires := []*frrProto.Advertisement{
		{
			InterfaceAddress: "10.0.2.0",
			LinkType:         "Stub Network",
		},
		{
			InterfaceAddress: "10.0.12.1",
			LinkType:         "Stub Network",
		},
	}

	t.Run("TestUnderadvertised", func(t *testing.T) {
		assert.True(t, ana.AnalysisResult.RouterAnomaly.HasUnderAdvertisedPrefixes)
		assert.Equal(t, 2, len(ana.AnalysisResult.RouterAnomaly.MissingEntries))
		assert.Equal(t, len(expectedMissingEntrires), len(ana.AnalysisResult.RouterAnomaly.MissingEntries))
		missingOne := false
		if expectedMissingEntrires[0].InterfaceAddress == ana.AnalysisResult.RouterAnomaly.MissingEntries[0].InterfaceAddress {
			missingOne = strings.ToLower(expectedMissingEntrires[0].LinkType) == strings.ToLower(ana.AnalysisResult.RouterAnomaly.MissingEntries[0].LinkType)
		} else if expectedMissingEntrires[0].InterfaceAddress == ana.AnalysisResult.RouterAnomaly.MissingEntries[1].InterfaceAddress {
			missingOne = strings.ToLower(expectedMissingEntrires[0].LinkType) == strings.ToLower(ana.AnalysisResult.RouterAnomaly.MissingEntries[1].LinkType)
		}
		assert.True(t, missingOne)
	})

	isRouterLSDB2 := analyzer.GetRuntimeRouterData(frrMetrics.OspfRouterData, frrMetrics.StaticFrrConfiguration.Hostname)
	shouldRouterLSDB2 := getExpectedShouldRouterLSDBr101SuperfluousEntriesUnhappy()
	expectedSuperfluousEntrires := []*frrProto.Advertisement{
		{
			InterfaceAddress: "10.0.2.0",
			LinkType:         "Stub Network",
		},
		{
			InterfaceAddress: "10.0.12.1",
			LinkType:         "Stub Network",
		},
	}

	ana.RouterAnomalyAnalysisLSDB(accessList, &shouldRouterLSDB2, isRouterLSDB2)
	t.Run("TestOveradvertised", func(t *testing.T) {
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasUnderAdvertisedPrefixes)
		assert.True(t, ana.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes)
		assert.Equal(t, 2, len(ana.AnalysisResult.RouterAnomaly.SuperfluousEntries))
		assert.Equal(t, len(expectedSuperfluousEntrires), len(ana.AnalysisResult.RouterAnomaly.SuperfluousEntries))
		missingOne := false
		if expectedSuperfluousEntrires[0].InterfaceAddress == ana.AnalysisResult.RouterAnomaly.SuperfluousEntries[0].InterfaceAddress {
			missingOne = strings.ToLower(expectedSuperfluousEntrires[0].LinkType) == strings.ToLower(ana.AnalysisResult.RouterAnomaly.SuperfluousEntries[0].LinkType)
		} else if expectedSuperfluousEntrires[0].InterfaceAddress == ana.AnalysisResult.RouterAnomaly.SuperfluousEntries[1].InterfaceAddress {
			missingOne = strings.ToLower(expectedSuperfluousEntrires[0].LinkType) == strings.ToLower(ana.AnalysisResult.RouterAnomaly.SuperfluousEntries[1].LinkType)
		}
		assert.True(t, missingOne)
	})
}

func TestExternalLsaHappy1(t *testing.T) {

	less := func(a, b string) bool { return a < b }
	frrMetrics := getR101FRRdata()
	accessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)
	staticList := analyzer.GetStaticRouteList(frrMetrics.StaticFrrConfiguration, accessList)

	expectedShouldExternalLSDB := getExpectedShouldExternalLSDBr101(frrMetrics.StaticFrrConfiguration.Hostname, frrMetrics.StaticFrrConfiguration.OspfConfig.RouterId)

	actualshouldExternalLSDB := analyzer.GetStaticFileExternalData(frrMetrics.StaticFrrConfiguration, accessList, staticList)

	// - actualExternalLSDB -> GetRuntimeExternalRouterData
	// - predictedExternalLSDB
	t.Run("TestExternalDataStaticShouldAndIs", func(t *testing.T) {

		assert.Equal(t, expectedShouldExternalLSDB.Hostname, actualshouldExternalLSDB.Hostname)
		assert.Equal(t, expectedShouldExternalLSDB.RouterId, actualshouldExternalLSDB.RouterId)
		expectedAreaMapTmp := make(map[string][]string)
		expectedAreaListTmp := []string{}

		assert.Equal(t, len(expectedShouldExternalLSDB.Areas), len(actualshouldExternalLSDB.Areas))

		assert.Equal(t, expectedShouldExternalLSDB.Areas[0].LsaType, actualshouldExternalLSDB.Areas[0].LsaType)
		assert.Equal(t, len(expectedShouldExternalLSDB.Areas[0].Links), len(actualshouldExternalLSDB.Areas[0].Links))
		assert.Equal(t, expectedShouldExternalLSDB.Areas[0].Links[0].LinkStateId, actualshouldExternalLSDB.Areas[0].Links[0].LinkStateId)
		assert.Equal(t, expectedShouldExternalLSDB.Areas[0].Links[0].PrefixLength, actualshouldExternalLSDB.Areas[0].Links[0].PrefixLength)
		assert.Equal(t, expectedShouldExternalLSDB.Areas[0].Links[0].LinkType, actualshouldExternalLSDB.Areas[0].Links[0].LinkType)
		for _, area := range expectedShouldExternalLSDB.Areas {
			expectedAreaListTmp = append(expectedAreaListTmp, area.AreaName)
			for _, adv := range area.Links {
				expectedAreaMapTmp[area.AreaName] = append(expectedAreaMapTmp[area.AreaName], adv.LinkStateId)
			}
		}
		actualAreaMapTmp := make(map[string][]string)
		actualAreaListTmp := []string{}
		for _, area := range expectedShouldExternalLSDB.Areas {
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
	expectedIsExternalLSDB := getExternalIsExternalLSDBr101()
	actualIsExternalLSDB := analyzer.GetRuntimeExternalData(frrMetrics.OspfExternalData, staticList, frrMetrics.StaticFrrConfiguration.Hostname)

	// TODO: maybe add AreaName testing? For that area assignment needs to be done. It doesn't seem too easy and it's not really necessary. Considering that static and connected redistributions happen via LSA Type 5 anyway and if it's connected to an NSSA it will still show a type 5 lsa but in type 7 lsa testing it will correctly show the correct static and connected redistributions.
	t.Run("TestExternalDataRuntimeShouldAndIs", func(t *testing.T) {
		assert.Equal(t, expectedIsExternalLSDB.Hostname, actualIsExternalLSDB.Hostname)
		assert.Equal(t, expectedIsExternalLSDB.RouterId, actualIsExternalLSDB.RouterId)

		assert.Equal(t, len(expectedIsExternalLSDB.Areas), len(actualIsExternalLSDB.Areas))
		expectedTotalLinks := 0
		actualTotalLinks := 0

		for _, area := range expectedIsExternalLSDB.Areas {
			expectedTotalLinks += len(area.Links)
		}
		for _, area := range actualIsExternalLSDB.Areas {
			actualTotalLinks += len(area.Links)
		}

		assert.Equal(t, expectedTotalLinks, actualTotalLinks)

		// Create maps with LinkStateId as keys for comparison
		expectedTmp := map[string][]*frrProto.Advertisement{}
		actualTmp := map[string][]*frrProto.Advertisement{}

		// Populate the map for expected data
		for _, area := range expectedIsExternalLSDB.Areas {
			for _, link := range area.Links {
				expectedTmp[link.LinkStateId] = append(expectedTmp[link.LinkStateId], link)
			}
		}

		// Populate the map for actual data
		for _, area := range actualIsExternalLSDB.Areas {
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

func TestExternalLsaUnhappy1(t *testing.T) {
	ana := initAnalyzer()
	frrMetrics := getR101FRRdata()

	accessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)
	staticList := analyzer.GetStaticRouteList(frrMetrics.StaticFrrConfiguration, accessList)

	shouldExternalLSDB := analyzer.GetStaticFileExternalData(frrMetrics.StaticFrrConfiguration, accessList, staticList)
	isExternalLSDB := getIsExternalLSDBr101OverUnhappy()

	// Unadvertised: isExternalLSDB is empty
	ana.ExternalAnomalyAnalysisLSDB(shouldExternalLSDB, isExternalLSDB)
	t.Run("TestUnadvertisedPrefix", func(t *testing.T) {
		assert.True(t, ana.AnalysisResult.ExternalAnomaly.HasUnderAdvertisedPrefixes)
		assert.Equal(t, 1, len(ana.AnalysisResult.ExternalAnomaly.MissingEntries))
		expectedMissingEntrires := []*frrProto.Advertisement{
			{
				InterfaceAddress: "192.168.1.0",
				LinkType:         "external",
			},
		}

		assert.Equal(t, len(expectedMissingEntrires), len(ana.AnalysisResult.ExternalAnomaly.MissingEntries))
		expectedME := expectedMissingEntrires[0]
		actualME := ana.AnalysisResult.ExternalAnomaly.MissingEntries[0]

		missingOne := strings.EqualFold(actualME.LinkType, expectedME.LinkType)
		if !missingOne {
			missingOne = strings.EqualFold(actualME.InterfaceAddress, expectedME.InterfaceAddress)
		}
		assert.True(t, missingOne)

	})

	// Overadvertised: isExternalLSDB is empty
	ana = initAnalyzer()
	frrMetrics = getR101FRRdata()
	isExternalLSDB = getIsExternalLSDBr101UnUnhappy()
	shouldExternalLSDB = analyzer.GetStaticFileExternalData(frrMetrics.StaticFrrConfiguration, accessList, staticList)
	ana.ExternalAnomalyAnalysisLSDB(shouldExternalLSDB, isExternalLSDB)
	t.Run("TestUnadvertisedPrefix", func(t *testing.T) {
		assert.False(t, ana.AnalysisResult.ExternalAnomaly.HasUnderAdvertisedPrefixes)
		assert.True(t, ana.AnalysisResult.ExternalAnomaly.HasOverAdvertisedPrefixes)
		assert.Equal(t, 1, len(ana.AnalysisResult.ExternalAnomaly.SuperfluousEntries))
		expectedMissingEntrires := []*frrProto.Advertisement{
			{
				InterfaceAddress: "192.168.2.0",
				LinkType:         "external",
			},
		}

		assert.Equal(t, len(expectedMissingEntrires), len(ana.AnalysisResult.ExternalAnomaly.SuperfluousEntries))
		expectedME := expectedMissingEntrires[0]
		actualME := ana.AnalysisResult.ExternalAnomaly.SuperfluousEntries[0]

		missingOne := strings.EqualFold(actualME.LinkType, expectedME.LinkType)
		if !missingOne {
			missingOne = strings.EqualFold(actualME.InterfaceAddress, expectedME.InterfaceAddress)
		}
		assert.True(t, missingOne)

	})

}

func TestAnomalyAnalysisHappy1(t *testing.T) {

	ana := initAnalyzer()
	frrMetrics := getR101FRRdata()
	accessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)
	staticRouteMap := analyzer.GetStaticRouteList(frrMetrics.StaticFrrConfiguration, accessList)

	isRouterLSDB := analyzer.GetRuntimeRouterData(frrMetrics.OspfRouterData, frrMetrics.StaticFrrConfiguration.Hostname)

	_, shouldRouterLSDB := analyzer.GetStaticFileRouterData(frrMetrics.StaticFrrConfiguration)
	ana.RouterAnomalyAnalysisLSDB(accessList, shouldRouterLSDB, isRouterLSDB)

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
	ana.ExternalAnomalyAnalysisLSDB(predictedExternalLSDB, runtimeExternalLSDB)

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

func TestAnomalyAnalysisLsaFiveHappy1(t *testing.T) {

}

// TestAnomalyAnalysis1
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

	// fmt.Println("---------------------- Predicted ----------------------")
	// fmt.Println(predictedNssaExternalLSDB)
	// fmt.Println("---------------------- Predicted ----------------------")

	// fmt.Println("---------------------- Runtime ----------------------")
	// fmt.Println(runtimeNssaExternalLSDB)
	// fmt.Println("---------------------- Runtime ----------------------")

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
