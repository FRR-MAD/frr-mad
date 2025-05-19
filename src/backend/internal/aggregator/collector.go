package aggregator

import (
	"fmt"
	"log"
	"time"

	frrSocket "github.com/frr-mad/frr-mad/src/backend/internal/aggregator/frrsockets"
	frrProto "github.com/frr-mad/frr-mad/src/backend/pkg"
	"github.com/frr-mad/frr-mad/src/logger"
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
	}
}

func initFullFrrData() *frrProto.FullFRRData {
	fullFrrData := &frrProto.FullFRRData{
		OspfDatabase:           &frrProto.OSPFDatabase{},
		GeneralOspfInformation: &frrProto.GeneralOspfInformation{},
		OspfRouterData:         &frrProto.OSPFRouterData{},
		OspfRouterDataAll:      &frrProto.OSPFRouterData{},
		OspfNetworkData:        &frrProto.OSPFNetworkData{},
		OspfNetworkDataAll:     &frrProto.OSPFNetworkData{},
		OspfSummaryData:        &frrProto.OSPFSummaryData{},
		OspfSummaryDataAll:     &frrProto.OSPFSummaryData{},
		OspfAsbrSummaryData:    &frrProto.OSPFAsbrSummaryData{},
		OspfExternalData:       &frrProto.OSPFExternalData{},
		OspfNssaExternalData:   &frrProto.OSPFNssaExternalData{},
		OspfExternalAll:        &frrProto.OSPFExternalAll{},
		OspfNssaExternalAll:    &frrProto.OSPFNssaExternalAll{},
		OspfNeighbors:          &frrProto.OSPFNeighbors{},
		Interfaces:             &frrProto.InterfaceList{},
		RoutingInformationBase: &frrProto.RoutingInformationBase{},
		RibFibSummaryRoutes:    &frrProto.RibFibSummaryRoutes{},
		StaticFrrConfiguration: &frrProto.StaticFRRConfiguration{},
		SystemMetrics:          &frrProto.SystemMetrics{},
		FrrRouterData:          &frrProto.FRRRouterData{},
	}

	return fullFrrData
}

func (c *Collector) Collect() error {
	c.logger.Debug(fmt.Sprintf("Address of collector: %p\n", c))

	if c.FullFrrData == nil {
		c.FullFrrData = initFullFrrData()
	} else {
		c.ensureFieldsInitialized()
	}

	executor := NewFRRCommandExecutor(c.socketPath, 2*time.Second)

	fetchAndMerge := func(name string, target proto.Message, fetchFunc func() (proto.Message, error)) {
		result, err := fetchFunc()
		if err != nil {
			c.logger.Error(err.Error())
			if name == "StaticFRRConfig" {
				log.Panic(err)
			}
			return
		}

		// TODO: Yes, do something here
		// Merge the fetched data into the target.
		// Reset target by creating a new instance of the same type

		if p, ok := target.(interface{ Reset() }); ok {
			p.Reset()
		}
		proto.Merge(target, result)

		c.logger.Debug(fmt.Sprintf("Response of Fetch%s(): %v\n", name, target))
		c.logger.Debug(fmt.Sprintf("Response of Fetch%s() Address: %p\n", name, target))
	}

	fetchAndMerge("StaticFRRConfig", c.FullFrrData.StaticFrrConfiguration, func() (proto.Message, error) {
		return fetchStaticFRRConfig()
	})

	fetchAndMerge("GeneralOSPFInformation", c.FullFrrData.GeneralOspfInformation, func() (proto.Message, error) {
		return FetchGeneralOSPFInformation(executor)
	})

	fetchAndMerge("OSPFRouterData", c.FullFrrData.OspfRouterData, func() (proto.Message, error) {
		return FetchOSPFRouterData(executor)
	})

	fetchAndMerge("OSPFRouterDataAll", c.FullFrrData.OspfRouterDataAll, func() (proto.Message, error) {
		return FetchOSPFRouterDataAll(executor)
	})

	fetchAndMerge("OSPFNetworkData", c.FullFrrData.OspfNetworkData, func() (proto.Message, error) {
		return FetchOSPFNetworkData(executor)
	})

	fetchAndMerge("OSPFNetworkDataAll", c.FullFrrData.OspfNetworkDataAll, func() (proto.Message, error) {
		return FetchOSPFNetworkDataAll(executor)
	})

	fetchAndMerge("OSPFSummaryData", c.FullFrrData.OspfSummaryData, func() (proto.Message, error) {
		return FetchOSPFSummaryData(executor)
	})

	fetchAndMerge("OSPFSummaryDataAll", c.FullFrrData.OspfSummaryDataAll, func() (proto.Message, error) {
		return FetchOSPFSummaryDataAll(executor)
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

	fetchAndMerge("OSPFExternalAll", c.FullFrrData.OspfExternalAll, func() (proto.Message, error) {
		return FetchOSPFExternalAll(executor)
	})

	fetchAndMerge("OSPFNssaExternalAll", c.FullFrrData.OspfNssaExternalAll, func() (proto.Message, error) {
		return FetchOSPFNssaExternalAll(executor)
	})

	fetchAndMerge("OSPFNeighbors", c.FullFrrData.OspfNeighbors, func() (proto.Message, error) {
		return FetchOSPFNeighbors(executor)
	})

	fetchAndMerge("InterfaceStatus", c.FullFrrData.Interfaces, func() (proto.Message, error) {
		return FetchInterfaceStatus(executor)
	})

	fetchAndMerge("ExpectedRoutes", c.FullFrrData.RoutingInformationBase, func() (proto.Message, error) {
		return FetchRib(executor)
	})

	fetchAndMerge("RibFibSummaryRoutes", c.FullFrrData.RibFibSummaryRoutes, func() (proto.Message, error) {
		return FetchRibFibSummary(executor)
	})

	fetchAndMerge("SystemMetrics", c.FullFrrData.SystemMetrics, func() (proto.Message, error) {
		return c.fetcher.CollectSystemMetrics()
	})

	fetchAndMerge("FrrRouterData", c.FullFrrData.FrrRouterData, func() (proto.Message, error) {
		frrRouterData := &frrProto.FRRRouterData{
			RouterName:   c.FullFrrData.StaticFrrConfiguration.Hostname,
			OspfRouterId: c.FullFrrData.OspfDatabase.RouterId,
		}
		return frrRouterData, nil
	})

	return nil
}

func (c *Collector) ensureFieldsInitialized() {
	if c.FullFrrData.StaticFrrConfiguration == nil {
		c.FullFrrData.StaticFrrConfiguration = &frrProto.StaticFRRConfiguration{}
	}

	if c.FullFrrData.GeneralOspfInformation == nil {
		c.FullFrrData.GeneralOspfInformation = &frrProto.GeneralOspfInformation{}
	}

	if c.FullFrrData.OspfRouterData == nil {
		c.FullFrrData.OspfRouterData = &frrProto.OSPFRouterData{}
	}

	if c.FullFrrData.OspfRouterDataAll == nil {
		c.FullFrrData.OspfRouterDataAll = &frrProto.OSPFRouterData{}

	}
	if c.FullFrrData.OspfNetworkData == nil {
		c.FullFrrData.OspfNetworkData = &frrProto.OSPFNetworkData{}
	}

	if c.FullFrrData.OspfNetworkDataAll == nil {
		c.FullFrrData.OspfNetworkDataAll = &frrProto.OSPFNetworkData{}
	}

	if c.FullFrrData.OspfSummaryData == nil {
		c.FullFrrData.OspfSummaryData = &frrProto.OSPFSummaryData{}
	}

	if c.FullFrrData.OspfSummaryDataAll == nil {
		c.FullFrrData.OspfSummaryDataAll = &frrProto.OSPFSummaryData{}
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

	if c.FullFrrData.OspfExternalAll == nil {
		c.FullFrrData.OspfExternalAll = &frrProto.OSPFExternalAll{}
	}

	if c.FullFrrData.OspfNssaExternalAll == nil {
		c.FullFrrData.OspfNssaExternalAll = &frrProto.OSPFNssaExternalAll{}
	}

	if c.FullFrrData.OspfNeighbors == nil {
		c.FullFrrData.OspfNeighbors = &frrProto.OSPFNeighbors{}
	}

	if c.FullFrrData.Interfaces == nil {
		c.FullFrrData.Interfaces = &frrProto.InterfaceList{}
	}

	if c.FullFrrData.RoutingInformationBase == nil {
		c.FullFrrData.RoutingInformationBase = &frrProto.RoutingInformationBase{}
	}

	if c.FullFrrData.RibFibSummaryRoutes == nil {
		c.FullFrrData.RibFibSummaryRoutes = &frrProto.RibFibSummaryRoutes{}
	}

	if c.FullFrrData.SystemMetrics == nil {
		c.FullFrrData.SystemMetrics = &frrProto.SystemMetrics{}
	}
}

func (c *Collector) GetFetcherForTesting() *Fetcher {
	return c.fetcher
}

func (c *Collector) GetConfigPathForTesting() string {
	return c.configPath
}
