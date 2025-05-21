package analyzer_test

import (
	"encoding/json"
	"testing"

	"github.com/frr-mad/frr-mad/src/backend/internal/aggregator"
	"github.com/frr-mad/frr-mad/src/backend/internal/analyzer"
)

func TestTempCust(t *testing.T) {
	configPath := "/path/to/config/file"

	configCust, _ := aggregator.ParseStaticFRRConfig(configPath)

	prettyConfig, _ := json.MarshalIndent(configCust, "", "  ")
	t.Log("Parsed Config")
	t.Log(string(prettyConfig))

	accessList := analyzer.GetAccessList(configCust)
	prettyAccessList, _ := json.MarshalIndent(accessList, "", "  ")
	t.Log("Parsed Config")
	t.Log(string(prettyAccessList))

	staticRouteMap := analyzer.GetStaticRouteList(configCust, accessList)
	prettyStaticRouteMap, _ := json.MarshalIndent(staticRouteMap, "", "  ")
	t.Log("Parsed Config")
	t.Log(string(prettyStaticRouteMap))

	peerInterfaceMap := analyzer.GetPeerNetworkAddress(configCust)
	prettyPeerInterface, _ := json.MarshalIndent(peerInterfaceMap, "", "  ")
	t.Log("Parsed Config")
	t.Log(string(prettyPeerInterface))

	//peerNeighborMap := analyzer.GetPeerNeighbor(a.metrics.OspfNeighbors, peerInterfaceMap)
	//prettyPeerNeighborMap, _ := json.MarshalIndent(staticRouteMap, "", "  ")
	//t.Log("Parsed Config")
	//t.Log(string(prettyStaticRouteMap))

	hostname := configCust.Hostname
	t.Log(hostname)

	ana := initAnalyzer()

	_, shouldRouterLSDB := ana.GetStaticFileRouterData(configCust)
	prettyShouldRouterLSDB, _ := json.MarshalIndent(shouldRouterLSDB, "", "  ")
	t.Log("Parsed Config")
	t.Log(string(prettyShouldRouterLSDB))

	t.Run("TestFoobar", func(t *testing.T) {
	})

}
