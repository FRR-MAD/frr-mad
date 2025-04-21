package aggregator

import (
	"fmt"
	"log"
	"time"

	frrSocket "github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator/frrsockets"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
)

type Collector struct {
	fetcher     *Fetcher
	configPath  string
	socketPath  string
	logger      *logger.Logger
	FullFrrData *frrProto.FullFRRData
}

func NewFRRCommandExecutor(socketDir string, timeout time.Duration) *frrSocket.FRRCommandExecutor {
	return &frrSocket.FRRCommandExecutor{
		DirPath: socketDir,
		Timeout: timeout,
	}
}

func NewCollector(metricsURL, configPath, socketPath string, logger *logger.Logger) *Collector {
	return &Collector{
		fetcher:     NewFetcher(metricsURL),
		configPath:  configPath,
		socketPath:  socketPath,
		logger:      logger,
		FullFrrData: &frrProto.FullFRRData{},
	}
}

func (c *Collector) Collect() (*frrProto.FullFRRData, error) {
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
		//os.Exit(1)
	}

	c.FullFrrData.StaticFrrConfiguration = staticFRRConfigParsed
	//fmt.Printf("Response of FetchStaticFRRConfig(): \n%+v\n", staticFRRConfigParsed)
	c.logger.Debug("Response of FetchStaticFRRConfig(): " + staticFRRConfigParsed.String())

	ospfRouterData, err := FetchOSPFRouterData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	//fmt.Printf("Response: \n%+v\n", ospfRouterData)
	c.FullFrrData.OspfRouterData = ospfRouterData
	c.logger.Debug("Response of FetchOSPFRouterData(): " + ospfRouterData.String())
	c.logger.Debug(fmt.Sprintf("Response of FetchOSPFRouterData() Address: %p\n", ospfRouterData))

	ospfNetworkData, err := FetchOSPFNetworkData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	//fmt.Printf("Response: \n%+v\n", ospfNetworkData)
	c.FullFrrData.OspfNetworkData = ospfNetworkData
	c.logger.Debug("Response of FetchOSPFNetworkData(): " + ospfNetworkData.String())

	ospfSummaryData, err := FetchOSPFSummaryData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	//fmt.Printf("Response: \n%+v\n", ospfSummaryData)
	c.FullFrrData.OspfSummaryData = ospfSummaryData
	c.logger.Debug("Response of FetchOSPFSummaryData(): " + ospfSummaryData.String())

	ospfAsbrSummaryData, err := FetchOSPFAsbrSummaryData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	//fmt.Printf("Response: \n%+v\n", ospfAsbrSummaryData)
	c.FullFrrData.OspfAsbrSummaryData = ospfAsbrSummaryData
	c.logger.Debug("Response of FetchOSPFAsbrSummaryData(): " + ospfAsbrSummaryData.String())

	ospfExternalData, err := FetchOSPFExternalData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	//fmt.Printf("Response: \n%+v\n", ospfExternalData)
	c.FullFrrData.OspfExternalData = ospfExternalData
	c.logger.Debug("Response of FetchOSPFExternalData(): " + ospfExternalData.String())
	c.logger.Debug(fmt.Sprintf("Response of FetchOSPFExternalData() Address: %p\n", ospfExternalData))

	ospfNssaExternalData, err := FetchOSPFNssaExternalData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	//fmt.Printf("Response: \n%+v\n", ospfNssaExternalData)
	c.FullFrrData.OspfNssaExternalData = ospfNssaExternalData
	c.logger.Debug("Response of FetchOSPFNssaExternalData(): " + ospfNssaExternalData.String())

	out1, err := FetchFullOSPFDatabase(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response FetchFullOSPFDatabase: \n%+v\n", out1)
	c.logger.Debug("Response of FetchFullOSPFDatabase(): " + out1.String())

	out2, err := FetchOSPFDuplicateCandidates(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response FetchOSPFDuplicateCandidates: \n%+v\n", out2)
	c.logger.Debug("Response of FetchOSPFDuplicateCandidates(): " + out2.String())

	out3, err := FetchOSPFNeighbors(executor)
	if err != nil {
		fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response FetchOSPFNeighbors: \n%+v\n", out3)
	c.logger.Debug("Response of FetchOSPFNeighbors(): " + out3.String())
	c.logger.Debug(fmt.Sprintf("Response of FetchOSPFNeighbors() Address: %p", out3))

	out4, err := FetchInterfaceStatus(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response FetchInterfaceStatus: \n%+v\n", out4)

	out5, err := FetchExpectedRoutes(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response FetchExpectedRoutes: \n%+v\n", out5)

	//os.Exit(0)

	//config, err := ParseStaticFRRConfig(c.configPath)
	if err != nil {
		return nil, fmt.Errorf("config parse failed: %w", err)
	}

	systemMetrics, err := c.fetcher.CollectSystemMetrics()
	if err != nil {
		return nil, fmt.Errorf("system metrics failed: %w", err)
	}

	c.FullFrrData.SystemMetrics = systemMetrics

	state := &frrProto.FullFRRData{
		//Ospf:      ospfMetrics,
		//Config: config,
		SystemMetrics: systemMetrics,
	}

	return state, nil
}

// Functions for testing maybe remove later
func (c *Collector) GetFetcherForTesting() *Fetcher {
	return c.fetcher
}

func (c *Collector) GetConfigPathForTesting() string {
	return c.configPath
}
