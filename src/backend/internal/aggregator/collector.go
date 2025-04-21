package aggregator

import (
	"fmt"
	"log"
	"os"
	"time"

	frrSocket "github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator/frrsockets"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type Collector struct {
	fetcher        *Fetcher
	configPath     string
	socketPath     string
	logger         *logger.Logger
	cache          *frrProto.CombinedState
	TempFrrMetrics *TempFRRMetrics
}

type TempFRRMetrics struct {
	StaticFRRConfiguration *frrProto.StaticFRRConfiguration
	OspfRouterData         *frrProto.OSPFRouterData
	OspfNetworkData        *frrProto.OSPFNetworkData
	OspfSummaryData        *frrProto.OSPFSummaryData
	OspfAsbrSummaryData    *frrProto.OSPFAsbrSummaryData
	OspfExternalData       *frrProto.OSPFExternalData
	OspfNssaExternalData   *frrProto.OSPFNssaExternalData
	SystemMetrics          *frrProto.SystemMetrics
}

func NewFRRCommandExecutor(socketDir string, timeout time.Duration) *frrSocket.FRRCommandExecutor {
	return &frrSocket.FRRCommandExecutor{
		DirPath: socketDir,
		Timeout: timeout,
	}
}

func NewCollector(metricsURL, configPath, socketPath string, logger *logger.Logger) *Collector {
	return &Collector{
		fetcher:        NewFetcher(metricsURL),
		configPath:     configPath,
		socketPath:     socketPath,
		logger:         logger,
		TempFrrMetrics: &TempFRRMetrics{},
	}
}

func (c *Collector) Collect() (*frrProto.CombinedState, error) {
	// ospfMetrics, err := c.fetcher.FetchOSPF()
	// if err != nil {
	// 	return nil, fmt.Errorf("OSPF fetch failed: %w", err)
	// }

	// Previously hard coded socket path to /var/run/frr
	executor := NewFRRCommandExecutor(c.socketPath, 2*time.Second)

	staticFRRConfigParsed, err := fetchStaticFRRConfig()
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		log.Panic(err)
		os.Exit(1)
	}

	c.TempFrrMetrics.StaticFRRConfiguration = staticFRRConfigParsed
	//fmt.Printf("Response of FetchStaticFRRConfig(): \n%+v\n", staticFRRConfigParsed)
	c.logger.Debug("Response of FetchStaticFRRConfig(): " + staticFRRConfigParsed.String())

	ospfRouterData, err := FetchOSPFRouterData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	//fmt.Printf("Response: \n%+v\n", ospfRouterData)
	c.TempFrrMetrics.OspfRouterData = ospfRouterData
	c.logger.Debug("Response of FetchOSPFRouterData(): " + ospfRouterData.String())

	ospfNetworkData, err := FetchOSPFNetworkData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	//fmt.Printf("Response: \n%+v\n", ospfNetworkData)
	c.TempFrrMetrics.OspfNetworkData = ospfNetworkData
	c.logger.Debug("Response of FetchOSPFNetworkData(): " + ospfNetworkData.String())

	ospfSummaryData, err := FetchOSPFSummaryData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	//fmt.Printf("Response: \n%+v\n", ospfSummaryData)
	c.TempFrrMetrics.OspfSummaryData = ospfSummaryData
	c.logger.Debug("Response of FetchOSPFSummaryData(): " + ospfSummaryData.String())

	ospfAsbrSummaryData, err := FetchOSPFAsbrSummaryData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	//fmt.Printf("Response: \n%+v\n", ospfAsbrSummaryData)
	c.TempFrrMetrics.OspfAsbrSummaryData = ospfAsbrSummaryData
	c.logger.Debug("Response of FetchOSPFAsbrSummaryData(): " + ospfAsbrSummaryData.String())

	ospfExternalData, err := FetchOSPFExternalData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	//fmt.Printf("Response: \n%+v\n", ospfExternalData)
	c.TempFrrMetrics.OspfExternalData = ospfExternalData
	c.logger.Debug("Response of FetchOSPFExternalData(): " + ospfExternalData.String())

	ospfNssaExternalData, err := FetchOSPFNssaExternalData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	//fmt.Printf("Response: \n%+v\n", ospfNssaExternalData)
	c.TempFrrMetrics.OspfNssaExternalData = ospfNssaExternalData
	c.logger.Debug("Response of FetchOSPFNssaExternalData(): " + ospfNssaExternalData.String())

	//os.Exit(0)

	//config, err := ParseStaticFRRConfig(c.configPath)
	if err != nil {
		return nil, fmt.Errorf("config parse failed: %w", err)
	}

	systemMetrics, err := c.fetcher.CollectSystemMetrics()
	if err != nil {
		return nil, fmt.Errorf("system metrics failed: %w", err)
	}

	c.TempFrrMetrics.SystemMetrics = systemMetrics

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
