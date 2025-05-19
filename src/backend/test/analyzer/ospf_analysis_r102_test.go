package analyzer_test

import (
	"strings"
	"testing"

	"github.com/frr-mad/frr-mad/src/backend/internal/analyzer"
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/stretchr/testify/assert"
)

func TestRouterLsaHappy2(t *testing.T) {

	ana := initAnalyzer()

	frrMetrics := getR102FRRdata()
	expectedAccessList := getExpectedAccessListr102Happy()
	expectedIsRouterLSDB := getExpectedIsRouterLSDBr102Happy()
	expectedShouldRouterLSDB := getExpectedShouldRouterLSDBr102Happy()
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
		"192.168.11.0": {
			IpAddress:    "192.168.11.0",
			PrefixLength: 24,
			NextHop:      "192.168.101.93",
		},
	}

	actualStaticList := analyzer.GetStaticRouteList(frrMetrics.StaticFrrConfiguration, actualAccessList)
	peerInterfaceMap := analyzer.GetPeerNetworkAddress(frrMetrics.StaticFrrConfiguration)

	var actualStaticListKeys []string
	for k, _ := range actualStaticList {
		actualStaticListKeys = append(actualStaticListKeys, k)
	}
	var expectedStaticListKeys []string
	for k, _ := range actualStaticList {
		expectedStaticListKeys = append(expectedStaticListKeys, k)
	}

	// Runtime parsing of router
	actualRuntimeRouterLSDB, _ := analyzer.GetRuntimeRouterDataSelf(frrMetrics.OspfRouterData, frrMetrics.StaticFrrConfiguration.Hostname, peerInterfaceMap, ana.Logger)

	expectedRuntimeRouterLSDBAreaLength := len(expectedIsRouterLSDB.Areas)
	actualRuntimeRouterLSDBAreaLength := len(actualRuntimeRouterLSDB.Areas)

	expectedRuntimeRouterLSDBAreaMap := make(map[string][]*frrProto.Advertisement)
	expectedRuntimeRouterLSDBAreas := []string{}
	expectedRuntimeRouterLSDBLsaType := map[string]string{}
	for _, area := range expectedIsRouterLSDB.Areas {
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
		assert.Equal(t, expectedIsRouterLSDB.RouterId, actualRuntimeRouterLSDB.RouterId)
		assert.Equal(t, expectedIsRouterLSDB.Hostname, actualRuntimeRouterLSDB.Hostname)
		assert.Equal(t, expectedRuntimeRouterLSDBAreaLength, actualRuntimeRouterLSDBAreaLength)
		assert.True(t, cmp.Diff(expectedRuntimeRouterLSDBAreas, actualRuntimeRouterLSDBAreas, cmpopts.SortSlices(less)) == "")
		for _, key := range expectedRuntimeRouterLSDBAreas {
			assert.Equal(t, expectedRuntimeRouterLSDBLsaType[key], actualRuntimeRouterLSDBLsaType[key])
		}
		for _, key := range expectedRuntimeRouterLSDBAreas {
			t.Logf("%v\n", actualRuntimeRouterLSDBAreaMap[key])

		}

	})

	isNssa, actualPredictedRouterLSDB := ana.GetStaticFileRouterData(frrMetrics.StaticFrrConfiguration)

	ana.RouterAnomalyAnalysisLSDB(actualAccessList, actualPredictedRouterLSDB, actualRuntimeRouterLSDB)

	//analyzer.RouterAnomalyAnalysis(accessList, shouldState, isState)
	t.Run("TestGetAccessList", func(t *testing.T) {
		assert.Equal(t, len(expectedAccessList), len(actualAccessList))
		assert.True(t, cmp.Diff(expectedAccessListKeys, actualAccessListKeys, cmpopts.SortSlices(less)) == "")
		for _, v := range actualAccessListKeys {
			assert.Equal(t, expectedAccessList[v], actualAccessList[v])
		}
	})

	t.Run("TestGetStaticRouteList", func(t *testing.T) {
		assert.Equal(t, len(expectedStaticList), len(actualStaticList))
		assert.Equal(t, expectedStaticListKeys, actualStaticListKeys)
		for _, v := range actualStaticListKeys {
			assert.Equal(t, actualStaticList[v], expectedStaticList[v])
		}
	})

	expectedPredictedRouterLSDBAreas := []string{}
	actualPredictedRouterLSDBAreas := []string{}

	expectedPredictedRouterLSDBLsaTypePerArea := make(map[string]string)
	actualPredictedRouterLSDBLsaTypePerArea := make(map[string]string)

	expectedPredictedRouterLSDBIntPerArea := make(map[string][]*frrProto.Advertisement)
	actualPredictedRouterLSDBIntPerArea := make(map[string][]*frrProto.Advertisement)

	for _, area := range expectedShouldRouterLSDB.Areas {
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
		assert.True(t, isNssa)
		assert.Equal(t, expectedShouldRouterLSDB.Hostname, actualPredictedRouterLSDB.Hostname)
		assert.Equal(t, expectedShouldRouterLSDB.RouterId, actualPredictedRouterLSDB.RouterId)
		assert.Equal(t, len(expectedPredictedRouterLSDBAreas), len(actualPredictedRouterLSDBAreas))
		assert.True(t, cmp.Diff(expectedPredictedRouterLSDBAreas, actualPredictedRouterLSDBAreas, cmpopts.SortSlices(less)) == "")

		for _, key := range expectedPredictedRouterLSDBAreas {
			assert.Equal(t, expectedPredictedRouterLSDBLsaTypePerArea[key], actualPredictedRouterLSDBLsaTypePerArea[key])
		}

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

}

func TestRouterLsaUnhappy2(t *testing.T) {

	ana := initAnalyzer()
	frrMetrics := getR102FRRdata()
	accessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)
	isRouterLSDB := getExpectedIsRouterLSDBr102MissingEntries()
	shouldRouterLSDB := getExpectedIsRouterLSDBr102Happy()
	ana.RouterAnomalyAnalysisLSDB(accessList, shouldRouterLSDB, isRouterLSDB)
	expectedMissingEntrires := []*frrProto.Advertisement{
		{
			InterfaceAddress: "10.0.12.2",
			LinkType:         "transit network",
		},
	}

	t.Run("TestUnadvertised", func(t *testing.T) {

		assert.True(t, ana.AnalysisResult.RouterAnomaly.HasUnAdvertisedPrefixes)
		assert.Equal(t, 1, len(ana.AnalysisResult.RouterAnomaly.MissingEntries))

		assert.Equal(t, len(expectedMissingEntrires), len(ana.AnalysisResult.RouterAnomaly.MissingEntries))
		missingOne := false
		if expectedMissingEntrires[0].InterfaceAddress == ana.AnalysisResult.RouterAnomaly.MissingEntries[0].InterfaceAddress {
			missingOne = strings.ToLower(expectedMissingEntrires[0].LinkType) == strings.ToLower(ana.AnalysisResult.RouterAnomaly.MissingEntries[0].LinkType)
		}
		assert.True(t, missingOne)
	})

}

func TestExternalLsaHappy2(t *testing.T) {

	ana := initAnalyzer()
	less := func(a, b string) bool { return a < b }
	frrMetrics := getR102FRRdata()
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
	actualPredictedExternalLSDB := ana.GetStaticFileExternalData(frrMetrics.StaticFrrConfiguration, accessList, staticList)

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
		Hostname: "r102",
		RouterId: "65.0.1.2",
		Areas: []*frrProto.AreaAnalyzer{
			{
				LsaType: "AS-external-LSA",
				Links: []*frrProto.Advertisement{
					{
						LinkStateId:  "192.168.11.0",
						PrefixLength: "24",
						LinkType:     "external",
					},
				},
			},
		},
	}
	actualRuntimeExternalLSDB := analyzer.GetRuntimeExternalDataSelf(frrMetrics.OspfExternalData, staticList, frrMetrics.StaticFrrConfiguration.Hostname, ana.Logger)

	//TODO: maybe add AreaName testing? For that area assignment needs to be done. It doesn't seem too easy and it's not really necessary. Considering that static and connected redistributions happen via LSA Type 5 anyway and if it's connected to an NSSA it will still show a type 5 lsa but in type 7 lsa testing it will correctly show the correct static and connected redistributions.
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

		expectedTmp := map[string][]*frrProto.Advertisement{}
		actualTmp := map[string][]*frrProto.Advertisement{}

		for _, area := range expectedRuntimeExternalLSDB.Areas {
			for _, link := range area.Links {
				expectedTmp[link.LinkStateId] = append(expectedTmp[link.LinkStateId], link)
			}
		}

		for _, area := range actualRuntimeExternalLSDB.Areas {
			for _, link := range area.Links {
				actualTmp[link.LinkStateId] = append(actualTmp[link.LinkStateId], link)
			}
		}

		assert.Equal(t, len(expectedTmp), len(actualTmp), "Expected and actual maps should have the same number of LinkStateIds")

		for linkStateId, expectedAdvs := range expectedTmp {
			actualAdvs, exists := actualTmp[linkStateId]
			assert.True(t, exists, "LinkStateId %s should exist in actual data", linkStateId)
			assert.Equal(t, len(expectedAdvs), len(actualAdvs), "Expected and actual should have same number of advertisements for LinkStateId %s", linkStateId)

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

func TestExternalLsaUnhappy2(t *testing.T) {
	ana := initAnalyzer()
	frrMetrics := getR102FRRdata()

	accessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)
	staticList := analyzer.GetStaticRouteList(frrMetrics.StaticFrrConfiguration, accessList)

	shouldExternalLSDB := ana.GetStaticFileExternalData(frrMetrics.StaticFrrConfiguration, accessList, staticList)
	isExternalLSDB := getIsExternalLSDBr102OverUnhappy()

	// Unadvertised: isExternalLSDB is empty
	ana.ExternalAnomalyAnalysisLSDB(shouldExternalLSDB, isExternalLSDB)
	t.Run("TestUnadvertisedPrefix", func(t *testing.T) {
		assert.True(t, ana.AnalysisResult.ExternalAnomaly.HasUnAdvertisedPrefixes)
		assert.Equal(t, 1, len(ana.AnalysisResult.ExternalAnomaly.MissingEntries))
		expectedMissingEntrires := []*frrProto.Advertisement{
			{
				InterfaceAddress: "192.168.11.0",
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
	frrMetrics = getR102FRRdata()
	isExternalLSDB = getIsExternalLSDBr102UnUnhappy()
	shouldExternalLSDB = ana.GetStaticFileExternalData(frrMetrics.StaticFrrConfiguration, accessList, staticList)
	ana.ExternalAnomalyAnalysisLSDB(shouldExternalLSDB, isExternalLSDB)
	t.Run("TestUnadvertisedPrefix", func(t *testing.T) {
		assert.False(t, ana.AnalysisResult.ExternalAnomaly.HasUnAdvertisedPrefixes)
		assert.True(t, ana.AnalysisResult.ExternalAnomaly.HasOverAdvertisedPrefixes)
		assert.Equal(t, 1, len(ana.AnalysisResult.ExternalAnomaly.SuperfluousEntries))
		expectedMissingEntrires := []*frrProto.Advertisement{
			{
				InterfaceAddress: "192.168.11.0",
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
