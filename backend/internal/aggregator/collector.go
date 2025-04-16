package aggregator

import (
	"fmt"

	frrProto "github.com/ba2025-ysmprc/frr-tui/backend/pkg"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Collector struct {
	fetcher    *Fetcher
	configPath string
	cache      *frrProto.CombinedState
}

func NewCollector(metricsURL, configPath string) *Collector {
	return &Collector{
		fetcher:    NewFetcher(metricsURL),
		configPath: configPath,
	}
}

func (c *Collector) Collect() (*frrProto.CombinedState, error) {
	ospfMetrics, err := c.fetcher.FetchOSPF()
	if err != nil {
		return nil, fmt.Errorf("OSPF fetch failed: %w", err)
	}

	config, err := ParseConfig(c.configPath)
	if err != nil {
		return nil, fmt.Errorf("config parse failed: %w", err)
	}

	systemMetrics, err := c.fetcher.CollectSystemMetrics()
	if err != nil {
		return nil, fmt.Errorf("system metrics failed: %w", err)
	}

	state := &frrProto.CombinedState{
		Timestamp: timestamppb.Now(),
		Ospf:      ospfMetrics,
		Config:    config,
		System:    systemMetrics,
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
