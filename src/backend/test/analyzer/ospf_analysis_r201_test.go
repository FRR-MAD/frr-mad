package analyzer_test

import (
	"testing"

	"github.com/frr-mad/frr-mad/src/backend/internal/analyzer"
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/stretchr/testify/assert"
)

func TestRouterLsaHappy3(t *testing.T) {

	ana := initAnalyzer()
	frrMetrics := getR201FRRdata()
	accessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)

	actualPeerInterfaceMap := analyzer.GetPeerNetworkAddress(frrMetrics.StaticFrrConfiguration)
	expectedPeerInterfaceMap := map[string]string{
		"eth3": "10.20.13.1",
		"eth4": "10.20.14.1",
	}
	actualPeerNeighborMap := analyzer.GetPeerNeighbor(frrMetrics.OspfNeighbors, actualPeerInterfaceMap)
	expectedPeerNeighborMap := map[string]string{
		"65.0.2.3": "10.20.13.1",
		"65.0.2.4": "10.20.14.1",
	}

	actualIsRouterLSDB, _ := analyzer.GetRuntimeRouterDataSelf(frrMetrics.OspfRouterData, frrMetrics.StaticFrrConfiguration.Hostname, actualPeerNeighborMap)
	expectedIsRouterLSDB := &frrProto.IntraAreaLsa{
		RouterId: "65.0.2.1",
		Hostname: "r201",
		Areas: []*frrProto.AreaAnalyzer{
			{
				AreaName: "0.0.0.0",
				LsaType:  "router-LSA",
				Links: []*frrProto.Advertisement{
					{
						InterfaceAddress: "10.20.12.1",
						LinkType:         "transit network",
					},
					{
						InterfaceAddress: "10.20.14.1",
						LinkType:         "point-to-point",
					},
					{
						InterfaceAddress: "10.20.13.1",
						LinkType:         "point-to-point",
					},
					{
						InterfaceAddress: "10.20.13.3",
						PrefixLength:     "32",
						LinkType:         "stub network",
					},
					{
						InterfaceAddress: "10.20.14.4",
						PrefixLength:     "32",
						LinkType:         "stub network",
					},
				},
			},
		},
	}

	_, shouldRouterLSDB := ana.GetStaticFileRouterData(frrMetrics.StaticFrrConfiguration)
	isRouterLSDB, _ := analyzer.GetRuntimeRouterDataSelf(frrMetrics.OspfRouterData, frrMetrics.StaticFrrConfiguration.Hostname, actualPeerNeighborMap)

	t.Run("TestHelperFunctionParsingR201", func(t *testing.T) {
		assert.Equal(t, len(expectedPeerInterfaceMap), len(actualPeerInterfaceMap))
		assert.Equal(t, expectedPeerInterfaceMap["eth3"], actualPeerInterfaceMap["eth3"])
		assert.Equal(t, expectedPeerInterfaceMap["eth4"], actualPeerInterfaceMap["eth4"])

		assert.Equal(t, len(expectedPeerNeighborMap), len(actualPeerNeighborMap))
		assert.Equal(t, expectedPeerNeighborMap["65.0.2.3"], actualPeerNeighborMap["65.0.2.3"])
		assert.Equal(t, expectedPeerNeighborMap["65.0.2.4"], actualPeerNeighborMap["65.0.2.4"])

		assert.Equal(t, expectedIsRouterLSDB.Hostname, actualIsRouterLSDB.Hostname)
		assert.Equal(t, expectedIsRouterLSDB.RouterId, actualIsRouterLSDB.RouterId)
		assert.Equal(t, len(expectedIsRouterLSDB.Areas), len(expectedIsRouterLSDB.Areas))

		expectedAreaList := []string{}
		expectedTmpList := []string{}
		expectedTmpMap := map[string]string{}
		for _, area := range expectedIsRouterLSDB.Areas {
			expectedAreaList = append(expectedAreaList, area.AreaName)
			for _, link := range area.Links {
				expectedTmpList = append(expectedTmpList, link.InterfaceAddress)
				expectedTmpMap[link.InterfaceAddress] = link.LinkType
			}
		}

		actualAreaList := []string{}
		actualTmpList := []string{}
		actualTmpMap := map[string]string{}
		for _, area := range actualIsRouterLSDB.Areas {
			actualAreaList = append(actualAreaList, area.AreaName)
			for _, link := range area.Links {
				actualTmpList = append(actualTmpList, link.InterfaceAddress)
				actualTmpMap[link.InterfaceAddress] = link.LinkType
			}
		}

		assert.Equal(t, len(expectedAreaList), len(actualAreaList))
		assert.Equal(t, len(expectedTmpList), len(actualTmpList))
		assert.Equal(t, len(expectedTmpMap), len(actualTmpMap))
		for _, entry := range expectedTmpList {
			_, exists := actualTmpMap[entry]
			assert.True(t, exists)
		}
		for _, entry := range actualTmpList {
			_, exists := expectedTmpMap[entry]
			assert.True(t, exists)
		}
		for _, entry := range actualTmpList {
			assert.Equal(t, expectedTmpMap[entry], actualTmpMap[entry])
		}

		// pretty1, _ := json.MarshalIndent(actualIsRouterLSDB, "", "  ")
		// pretty2, _ := json.MarshalIndent(expectedIsRouterLSDB, "", "  ")
		// t.Log(string(pretty1))
		// t.Log(string(pretty2))
	})

	ana.RouterAnomalyAnalysisLSDB(accessList, shouldRouterLSDB, isRouterLSDB)

	t.Run("TestAnomalyAnalysisR201", func(t *testing.T) {
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasUnAdvertisedPrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasDuplicatePrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasMisconfiguredPrefixes)
		assert.Equal(t, len(ana.AnalysisResult.RouterAnomaly.SuperfluousEntries), 0)
		assert.Equal(t, len(ana.AnalysisResult.RouterAnomaly.MissingEntries), 0)
		assert.Equal(t, len(ana.AnalysisResult.RouterAnomaly.DuplicateEntries), 0)
	})

	t.Run("TestCheckUnknown", func(t *testing.T) {
		_, shouldRouterState := ana.GetStaticFileRouterData(frrMetrics.StaticFrrConfiguration)
		isUnknown := getUnknown(shouldRouterState)

		assert.True(t, isUnknown)
	})

}

func TestRouterLsaUnhappy3(t *testing.T) {

	ana := initAnalyzer()
	frrMetrics := getR201FRRdata()
	accessList := analyzer.GetAccessList(frrMetrics.StaticFrrConfiguration)

	_, shouldRouterLSDB := ana.GetStaticFileRouterData(frrMetrics.StaticFrrConfiguration)
	isRouterLSDB := &frrProto.IntraAreaLsa{
		RouterId: "65.0.2.1",
		Hostname: "r201",
		Areas: []*frrProto.AreaAnalyzer{
			{
				AreaName: "0.0.0.0",
				LsaType:  "router-LSA",
				Links: []*frrProto.Advertisement{
					{
						InterfaceAddress: "10.20.13.1",
						LinkType:         "point-to-point",
					},
					{
						InterfaceAddress: "10.20.12.1",
						LinkType:         "transit network",
					},
					{
						InterfaceAddress: "10.20.14.1",
						LinkType:         "point-to-point",
					},
					{
						InterfaceAddress: "10.20.13.3",
						PrefixLength:     "32",
						LinkType:         "stub network",
					},
					{
						InterfaceAddress: "10.20.14.4",
						PrefixLength:     "32",
						LinkType:         "stub network",
					},
					{
						InterfaceAddress: "10.20.15.1",
						LinkType:         "transit network",
					},
				},
			},
		},
	}

	ana.RouterAnomalyAnalysisLSDB(accessList, shouldRouterLSDB, isRouterLSDB)

	t.Run("TestAnomalyAnalysisR201Overadvertised", func(t *testing.T) {
		assert.True(t, ana.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasUnAdvertisedPrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasDuplicatePrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasMisconfiguredPrefixes)
		assert.Equal(t, len(ana.AnalysisResult.RouterAnomaly.SuperfluousEntries), 1)
		assert.Equal(t, len(ana.AnalysisResult.RouterAnomaly.MissingEntries), 0)
		assert.Equal(t, len(ana.AnalysisResult.RouterAnomaly.DuplicateEntries), 0)

		assert.Equal(t, ana.AnalysisResult.RouterAnomaly.SuperfluousEntries[0].InterfaceAddress, "10.20.15.1")
	})

	isRouterLSDB = &frrProto.IntraAreaLsa{
		RouterId: "65.0.2.1",
		Hostname: "r201",
		Areas: []*frrProto.AreaAnalyzer{
			{
				AreaName: "0.0.0.0",
				LsaType:  "router-LSA",
				Links: []*frrProto.Advertisement{
					{
						InterfaceAddress: "10.20.13.1",
						LinkType:         "point-to-point",
					},
					{
						InterfaceAddress: "10.20.13.3",
						PrefixLength:     "32",
						LinkType:         "stub network",
					},
					{
						InterfaceAddress: "10.20.12.1",
						LinkType:         "transit network",
					},
					// {
					// 	InterfaceAddress: "10.20.14.1",
					// 	LinkType:         "point-to-point",
					// },
					{
						InterfaceAddress: "10.20.14.4",
						PrefixLength:     "32",
						LinkType:         "stub network",
					},
					// {
					// 	InterfaceAddress: "10.20.15.1",
					// 	LinkType:         "transit network",
					// },
				},
			},
		},
	}

	ana.RouterAnomalyAnalysisLSDB(accessList, shouldRouterLSDB, isRouterLSDB)

	t.Run("TestAnomalyAnalysisR201Unadvertised", func(t *testing.T) {
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes)
		assert.True(t, ana.AnalysisResult.RouterAnomaly.HasUnAdvertisedPrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasDuplicatePrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasMisconfiguredPrefixes)
		assert.Equal(t, len(ana.AnalysisResult.RouterAnomaly.SuperfluousEntries), 0)
		assert.Equal(t, len(ana.AnalysisResult.RouterAnomaly.MissingEntries), 1)
		assert.Equal(t, len(ana.AnalysisResult.RouterAnomaly.DuplicateEntries), 0)

		assert.Equal(t, ana.AnalysisResult.RouterAnomaly.MissingEntries[0].InterfaceAddress, "10.20.14.1")
	})

	isRouterLSDB = &frrProto.IntraAreaLsa{
		RouterId: "65.0.2.1",
		Hostname: "r201",
		Areas: []*frrProto.AreaAnalyzer{
			{
				AreaName: "0.0.0.0",
				LsaType:  "router-LSA",
				Links: []*frrProto.Advertisement{
					{
						InterfaceAddress: "10.20.13.1",
						LinkType:         "point-to-point",
					},
					{
						InterfaceAddress: "10.20.13.3",
						PrefixLength:     "32",
						LinkType:         "stub network",
					},
					{
						InterfaceAddress: "10.20.12.1",
						LinkType:         "transit network",
					},
					// missing entry
					// {
					// 	InterfaceAddress: "10.20.14.1",
					// 	LinkType:         "point-to-point",
					// },
					{
						InterfaceAddress: "10.20.14.4",
						PrefixLength:     "32",
						LinkType:         "stub network",
					},
					// overadvertised entry
					{
						InterfaceAddress: "10.20.15.1",
						LinkType:         "transit network",
					},
				},
			},
		},
	}

	ana.RouterAnomalyAnalysisLSDB(accessList, shouldRouterLSDB, isRouterLSDB)
	t.Run("TestAnomalyAnalysisR201WrongPeerAddress", func(t *testing.T) {

		assert.True(t, ana.AnalysisResult.RouterAnomaly.HasOverAdvertisedPrefixes)
		assert.True(t, ana.AnalysisResult.RouterAnomaly.HasUnAdvertisedPrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasDuplicatePrefixes)
		assert.False(t, ana.AnalysisResult.RouterAnomaly.HasMisconfiguredPrefixes)
		assert.Equal(t, len(ana.AnalysisResult.RouterAnomaly.SuperfluousEntries), 1)
		assert.Equal(t, len(ana.AnalysisResult.RouterAnomaly.MissingEntries), 1)
		assert.Equal(t, len(ana.AnalysisResult.RouterAnomaly.DuplicateEntries), 0)

		assert.Equal(t, ana.AnalysisResult.RouterAnomaly.SuperfluousEntries[0].InterfaceAddress, "10.20.15.1")
		assert.Equal(t, ana.AnalysisResult.RouterAnomaly.MissingEntries[0].InterfaceAddress, "10.20.14.1")
	})
}
