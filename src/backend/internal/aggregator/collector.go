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

func newCollector(configPath, socketPath string, logger *logger.Logger) *Collector {
	fullFrrData := initFullFrrData()

	return &Collector{
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
		c.logger.Debug("Initializing new FullFrrData structure")
		c.FullFrrData = initFullFrrData()
	} else {
		c.ensureFieldsInitialized()
	}

	executor := NewFRRCommandExecutor(c.socketPath, 2*time.Second)

	fetchAndMerge := func(name string, target proto.Message, fetchFunc func() (proto.Message, error)) {
		start := time.Now()
		result, err := fetchFunc()
		if err != nil {
			c.logger.WithAttrs(map[string]any{
				"component": "aggregator",
				"operation": name,
				"error":     err.Error(),
			}).Error("Failed to fetch data")

			if name == "StaticFRRConfig" {
				log.Panic(err)
				c.logger.Error(fmt.Sprintf("%v", err))
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

		c.logger.WithAttrs(map[string]any{
			"component": "aggregator",
			"operation": name,
			"response":  target,
			"duration":  time.Since(start).String(),
		}).Debug("Successfully fetched and merged data")
	}

	fetchAndMerge("StaticFRRConfig", c.FullFrrData.StaticFrrConfiguration, func() (proto.Message, error) {
		return fetchStaticFRRConfig()
	})

	fetchAndMerge("GeneralOSPFInformation", c.FullFrrData.GeneralOspfInformation, func() (proto.Message, error) {
		return fetchGeneralOSPFInformation(executor)
	})

	fetchAndMerge("OSPFRouterData", c.FullFrrData.OspfRouterData, func() (proto.Message, error) {
		return fetchOSPFRouterData(executor)
	})

	fetchAndMerge("OSPFRouterDataAll", c.FullFrrData.OspfRouterDataAll, func() (proto.Message, error) {
		return fetchOSPFRouterDataAll(executor)
	})

	fetchAndMerge("OSPFNetworkData", c.FullFrrData.OspfNetworkData, func() (proto.Message, error) {
		return fetchOSPFNetworkData(executor)
	})

	fetchAndMerge("OSPFNetworkDataAll", c.FullFrrData.OspfNetworkDataAll, func() (proto.Message, error) {
		return fetchOSPFNetworkDataAll(executor)
	})

	fetchAndMerge("OSPFSummaryData", c.FullFrrData.OspfSummaryData, func() (proto.Message, error) {
		return fetchOSPFSummaryData(executor)
	})

	fetchAndMerge("OSPFSummaryDataAll", c.FullFrrData.OspfSummaryDataAll, func() (proto.Message, error) {
		return fetchOSPFSummaryDataAll(executor)
	})

	fetchAndMerge("OSPFAsbrSummaryData", c.FullFrrData.OspfAsbrSummaryData, func() (proto.Message, error) {
		return fetchOSPFAsbrSummaryData(executor)
	})

	fetchAndMerge("OSPFExternalData", c.FullFrrData.OspfExternalData, func() (proto.Message, error) {
		return fetchOSPFExternalData(executor)
	})

	fetchAndMerge("OSPFNssaExternalData", c.FullFrrData.OspfNssaExternalData, func() (proto.Message, error) {
		return fetchOSPFNssaExternalData(executor)
	})

	fetchAndMerge("FullOSPFDatabase", c.FullFrrData.OspfDatabase, func() (proto.Message, error) {
		return fetchFullOSPFDatabase(executor)
	})

	fetchAndMerge("OSPFExternalAll", c.FullFrrData.OspfExternalAll, func() (proto.Message, error) {
		return fetchOSPFExternalAll(executor)
	})

	fetchAndMerge("OSPFNssaExternalAll", c.FullFrrData.OspfNssaExternalAll, func() (proto.Message, error) {
		return fetchOSPFNssaExternalAll(executor)
	})

	fetchAndMerge("OSPFNeighbors", c.FullFrrData.OspfNeighbors, func() (proto.Message, error) {
		return fetchOSPFNeighbors(executor)
	})

	fetchAndMerge("InterfaceStatus", c.FullFrrData.Interfaces, func() (proto.Message, error) {
		return fetchInterfaceStatus(executor)
	})

	fetchAndMerge("ExpectedRoutes", c.FullFrrData.RoutingInformationBase, func() (proto.Message, error) {
		return fetchRib(executor)
	})

	fetchAndMerge("RibFibSummaryRoutes", c.FullFrrData.RibFibSummaryRoutes, func() (proto.Message, error) {
		return fetchRibFibSummary(executor)
	})

	fetchAndMerge("SystemMetrics", c.FullFrrData.SystemMetrics, func() (proto.Message, error) {
		return collectSystemMetrics()
	})

	fetchAndMerge("FrrRouterData", c.FullFrrData.FrrRouterData, func() (proto.Message, error) {
		frrRouterData := &frrProto.FRRRouterData{
			RouterName:   c.FullFrrData.StaticFrrConfiguration.Hostname,
			OspfRouterId: c.FullFrrData.OspfDatabase.RouterId,
		}
		return frrRouterData, nil
	})

	c.logger.WithAttrs(map[string]any{
		"component": "aggregator",
		"action":    "complete_collection",
	}).Info("Completed data collection cycle")

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
