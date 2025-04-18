package aggregator

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	frrSocket "github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator/frrsockets"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Collector struct {
	fetcher    *Fetcher
	configPath string
	cache      *frrProto.CombinedState
}

func NewFRRCommandExecutor(socketDir string, timeout time.Duration) *frrSocket.FRRCommandExecutor {
	return &frrSocket.FRRCommandExecutor{
		DirPath: socketDir,
		Timeout: timeout,
	}
}

func NewCollector(metricsURL, configPath string) *Collector {
	return &Collector{
		fetcher:    NewFetcher(metricsURL),
		configPath: configPath,
	}
}

func (c *Collector) Collect() (*frrProto.CombinedState, error) {
	// ospfMetrics, err := c.fetcher.FetchOSPF()
	// if err != nil {
	// 	return nil, fmt.Errorf("OSPF fetch failed: %w", err)
	// }
	executor := NewFRRCommandExecutor("/var/run/frr", 2*time.Second)

	ospfRouterData, err := FetchOSPFRouterData(executor)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	data, err := json.MarshalIndent(ospfRouterData, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling struct:", err)
	} else {
		fmt.Println(string(data))
	}

	ospfNetworkData, err := FetchOSPFNetworkData(executor)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	println(ospfNetworkData.String())

	ospfSummaryData, err := FetchOSPFSummaryData(executor)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	println(ospfSummaryData.String())

	ospfAsbrSummaryData, err := FetchOSPFAsbrSummaryData(executor)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	println(ospfAsbrSummaryData.String())

	ospfExternalData, err := FetchOSPFExternalData(executor)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	println(ospfExternalData.String())

	ospfNssaExternalData, err := FetchOSPFNssaExternalData(executor)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	println(ospfNssaExternalData.String())

	os.Exit(0)

	//config, err := ParseStaticFRRConfig(c.configPath)
	if err != nil {
		return nil, fmt.Errorf("config parse failed: %w", err)
	}

	systemMetrics, err := c.fetcher.CollectSystemMetrics()
	if err != nil {
		return nil, fmt.Errorf("system metrics failed: %w", err)
	}

	state := &frrProto.CombinedState{
		Timestamp: timestamppb.Now(),
		//Ospf:      ospfMetrics,
		//Config: config,
		System: systemMetrics,
	}

	c.cache = state
	return state, nil
}

func (c *Collector) GetCache() *frrProto.CombinedState {
	return c.cache
}

// Functions for testing maybe remove later
func (c *Collector) GetFetcherForTesting() *Fetcher {
	return c.fetcher
}

func (c *Collector) GetConfigPathForTesting() string {
	return c.configPath
}

func (c *Collector) GetCacheForTesting() *frrProto.CombinedState {
	return c.cache
}
