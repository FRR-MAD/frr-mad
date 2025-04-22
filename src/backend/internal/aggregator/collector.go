package aggregator

import (
	"fmt"
	"log"
	"time"

	frrSocket "github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator/frrsockets"
	"github.com/ba2025-ysmprc/frr-mad/src/backend/internal/logger"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	"google.golang.org/protobuf/proto"
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

func newCollector(metricsURL, configPath, socketPath string, logger *logger.Logger) *Collector {
	fullFrrData := initFullFrrData()

	return &Collector{
		fetcher:     NewFetcher(metricsURL),
		configPath:  configPath,
		socketPath:  socketPath,
		logger:      logger,
		FullFrrData: fullFrrData,
		//FullFrrData: &frrProto.FullFRRData{},
	}
}

func initFullFrrData() *frrProto.FullFRRData {
	//var fullFrrData frrProto.FullFRRData

	fullFrrData := frrProto.FullFRRData{
		OspfDatabase:           &frrProto.OSPFDatabase{},
		OspfRouterData:         &frrProto.OSPFRouterData{},
		OspfNetworkData:        &frrProto.OSPFNetworkData{},
		OspfSummaryData:        &frrProto.OSPFSummaryData{},
		OspfAsbrSummaryData:    &frrProto.OSPFAsbrSummaryData{},
		OspfExternalData:       &frrProto.OSPFExternalData{},
		OspfNssaExternalData:   &frrProto.OSPFNssaExternalData{},
		OspfDuplicates:         &frrProto.OSPFDuplicates{},
		OspfNeighbors:          &frrProto.OSPFNeighbors{},
		Interfaces:             &frrProto.InterfaceList{},
		Routes:                 &frrProto.RouteList{},
		StaticFrrConfiguration: &frrProto.StaticFRRConfiguration{},
		SystemMetrics:          &frrProto.SystemMetrics{},
	}

	return &fullFrrData
}

func (c *Collector) Collect() error {
	c.logger.Debug(fmt.Sprintf("Address of collector: %p\n", c))

	executor := NewFRRCommandExecutor(c.socketPath, 2*time.Second)

	// Generic fetch function
	fetchAndMerge := func(name string, target proto.Message, fetchFunc func() (proto.Message, error)) {
		result, err := fetchFunc()
		if err != nil {
			c.logger.Error(err.Error())
			if name == "StaticFRRConfig" {
				log.Panic(err)
			}
			return
		}

		// Merge the fetched data into the target
		proto.Merge(target, result)

		// Log results consistently
		c.logger.Debug(fmt.Sprintf("Response of Fetch%s(): %v\n", name, target))
		c.logger.Debug(fmt.Sprintf("Response of Fetch%s() Address: %p\n", name, target))
	}

	c.initDataContainers()

	// Fetch each type of data using the generic function
	fetchAndMerge("StaticFRRConfig", c.FullFrrData.StaticFrrConfiguration, func() (proto.Message, error) {
		return fetchStaticFRRConfig()
	})

	fetchAndMerge("OSPFRouterData", c.FullFrrData.OspfRouterData, func() (proto.Message, error) {
		return FetchOSPFRouterData(executor)
	})

	fetchAndMerge("OSPFNetworkData", c.FullFrrData.OspfNetworkData, func() (proto.Message, error) {
		return FetchOSPFNetworkData(executor)
	})

	fetchAndMerge("OSPFSummaryData", c.FullFrrData.OspfSummaryData, func() (proto.Message, error) {
		return FetchOSPFSummaryData(executor)
	})

	fetchAndMerge("OSPFAsbrSummaryData", c.FullFrrData.OspfAsbrSummaryData, func() (proto.Message, error) {
		return FetchOSPFAsbrSummaryData(executor)
	})

	fetchAndMerge("OSPFExternalData", c.FullFrrData.OspfExternalData, func() (proto.Message, error) {
		return FetchOSPFExternalData(executor)
	})

	fetchAndMerge("OSPFNssaExternalData", c.FullFrrData.OspfNssaExternalData, func() (proto.Message, error) {
		return FetchOSPFNssaExternalData(executor)
	})

	fetchAndMerge("FullOSPFDatabase", c.FullFrrData.OspfDatabase, func() (proto.Message, error) {
		return FetchFullOSPFDatabase(executor)
	})

	fetchAndMerge("OSPFDuplicateCandidates", c.FullFrrData.OspfDuplicates, func() (proto.Message, error) {
		return FetchOSPFDuplicateCandidates(executor)
	})

	fetchAndMerge("OSPFNeighbors", c.FullFrrData.OspfNeighbors, func() (proto.Message, error) {
		return FetchOSPFNeighbors(executor)
	})

	fetchAndMerge("InterfaceStatus", c.FullFrrData.Interfaces, func() (proto.Message, error) {
		return FetchInterfaceStatus(executor)
	})

	fetchAndMerge("ExpectedRoutes", c.FullFrrData.Routes, func() (proto.Message, error) {
		return FetchExpectedRoutes(executor)
	})

	fetchAndMerge("SystemMetrics", c.FullFrrData.SystemMetrics, func() (proto.Message, error) {
		return c.fetcher.CollectSystemMetrics()
	})

	return nil
}

func (c *Collector) initDataContainers() {
	if c.FullFrrData.StaticFrrConfiguration == nil {
		c.FullFrrData.StaticFrrConfiguration = &frrProto.StaticFRRConfiguration{}
	}

	if c.FullFrrData.OspfRouterData == nil {
		c.FullFrrData.OspfRouterData = &frrProto.OSPFRouterData{}
	}

	if c.FullFrrData.OspfNetworkData == nil {
		c.FullFrrData.OspfNetworkData = &frrProto.OSPFNetworkData{}
	}

	if c.FullFrrData.OspfSummaryData == nil {
		c.FullFrrData.OspfSummaryData = &frrProto.OSPFSummaryData{}
	}

	if c.FullFrrData.OspfAsbrSummaryData == nil {
		c.FullFrrData.OspfAsbrSummaryData = &frrProto.OSPFAsbrSummaryData{}
	}

	if c.FullFrrData.OspfExternalData == nil {
		c.FullFrrData.OspfExternalData = &frrProto.OSPFExternalData{}
	}

	if c.FullFrrData.OspfNssaExternalData == nil {
		c.FullFrrData.OspfNssaExternalData = &frrProto.OSPFNssaExternalData{}
	}

	if c.FullFrrData.OspfDatabase == nil {
		c.FullFrrData.OspfDatabase = &frrProto.OSPFDatabase{}
	}

	if c.FullFrrData.OspfDuplicates == nil {
		c.FullFrrData.OspfDuplicates = &frrProto.OSPFDuplicates{}
	}

	if c.FullFrrData.OspfNeighbors == nil {
		c.FullFrrData.OspfNeighbors = &frrProto.OSPFNeighbors{}
	}

	if c.FullFrrData.Interfaces == nil {
		c.FullFrrData.Interfaces = &frrProto.InterfaceList{}
	}

	if c.FullFrrData.Routes == nil {
		c.FullFrrData.Routes = &frrProto.RouteList{}
	}

	if c.FullFrrData.SystemMetrics == nil {
		c.FullFrrData.SystemMetrics = &frrProto.SystemMetrics{}
	}
}

// Functions for testing maybe remove later
func (c *Collector) GetFetcherForTesting() *Fetcher {
	return c.fetcher
}

func (c *Collector) GetConfigPathForTesting() string {
	return c.configPath
}
