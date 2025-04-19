package aggregator

import (
	"fmt"
	"os"
	"time"

	frrSocket "github.com/ba2025-ysmprc/frr-mad/src/backend/internal/aggregator/frrsockets"
	frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

type OSPFRouterData struct {
	RouterID     string                    `json:"routerId"`
	RouterStates map[string]OSPFRouterArea `json:"Router Link States"`
}

type OSPFRouterArea map[string]OSPFRouterLSA
type OSPFRouterLSA struct {
	LSAAge            int                          `json:"lsaAge"`
	Options           string                       `json:"options"`
	LSAFlags          int                          `json:"lsaFlags"`
	Flags             int                          `json:"flags"`
	ASBR              bool                         `json:"asbr"`
	LSAType           string                       `json:"lsaType"`
	LinkStateID       string                       `json:"linkStateId"`
	AdvertisingRouter string                       `json:"advertisingRouter"`
	LSASeqNumber      string                       `json:"lsaSeqNumber"`
	Checksum          string                       `json:"checksum"`
	Length            int                          `json:"length"`
	NumOfLinks        int                          `json:"numOfLinks"`
	RouterLinks       map[string]OSPFRouterLSALink `json:"routerLinks"`
}
type OSPFRouterLSALink struct {
	LinkType                string `json:"linkType"`
	DesignatedRouterAddress string `json:"designatedRouterAddress,omitempty"`
	RouterInterfaceAddress  string `json:"routerInterfaceAddress,omitempty"`
	NetworkAddress          string `json:"networkAddress,omitempty"`
	NetworkMask             string `json:"networkMask,omitempty"`
	NumOfTosMetrics         int    `json:"numOfTosMetrics"`
	Tos0Metric              int    `json:"tos0Metric"`
}

type OSPFNetworkData struct {
	RouterID  string                  `json:"routerId"`
	NetStates map[string]NetAreaState `json:"Net Link States"`
}

type NetAreaState map[string]NetworkLSA

type NetworkLSA struct {
	LSAAge            int                       `json:"lsaAge"`
	Options           string                    `json:"options"`
	LSAFlags          int                       `json:"lsaFlags"`
	LSAType           string                    `json:"lsaType"`
	LinkStateID       string                    `json:"linkStateId"`
	AdvertisingRouter string                    `json:"advertisingRouter"`
	LSASeqNumber      string                    `json:"lsaSeqNumber"`
	Checksum          string                    `json:"checksum"`
	Length            int                       `json:"length"`
	NetworkMask       int                       `json:"networkMask"`
	AttachedRouters   map[string]AttachedRouter `json:"attchedRouters"`
}

type AttachedRouter struct {
	AttachedRouterID string `json:"attachedRouterId"`
}

type OSPFSummaryData struct {
	RouterID      string                      `json:"routerId"`
	NetStates     map[string]NetAreaState     `json:"Net Link States"`
	SummaryStates map[string]SummaryAreaState `json:"Summary Link States"`
}

type SummaryAreaState map[string]SummaryLSA

type SummaryLSA struct {
	LSAAge            int    `json:"lsaAge"`
	Options           string `json:"options"`
	LSAFlags          int    `json:"lsaFlags"`
	LSAType           string `json:"lsaType"`
	LinkStateID       string `json:"linkStateId"`
	AdvertisingRouter string `json:"advertisingRouter"`
	LSASeqNumber      string `json:"lsaSeqNumber"`
	Checksum          string `json:"checksum"`
	Length            int    `json:"length"`
	NetworkMask       int    `json:"networkMask"`
	Tos0Metric        int    `json:"tos0Metric"`
}

type OSPFAsbrSummaryData struct {
	RouterID          string                          `json:"routerId"`
	ASBRSummaryStates map[string]AsbrSummaryAreaState `json:"ASBR-Summary Link States"`
}

type AsbrSummaryAreaState map[string]AsbrSummaryLSA

type AsbrSummaryLSA struct {
	LSAAge            int    `json:"lsaAge"`
	Options           string `json:"options"`
	LSAFlags          int    `json:"lsaFlags"`
	LSAType           string `json:"lsaType"`
	LinkStateID       string `json:"linkStateId"`
	AdvertisingRouter string `json:"advertisingRouter"`
	LSASeqNumber      string `json:"lsaSeqNumber"`
	Checksum          string `json:"checksum"`
	Length            int    `json:"length"`
	NetworkMask       int    `json:"networkMask"`
	Tos0Metric        int    `json:"tos0Metric"`
}

type OSPFExternalData struct {
	RouterID         string                 `json:"routerId"`
	ASExternalStates map[string]ExternalLSA `json:"AS External Link States"`
}

type ExternalLSA struct {
	LSAAge            int    `json:"lsaAge"`
	Options           string `json:"options"`
	LSAFlags          int    `json:"lsaFlags"`
	LSAType           string `json:"lsaType"`
	LinkStateID       string `json:"linkStateId"`
	AdvertisingRouter string `json:"advertisingRouter"`
	LSASeqNumber      string `json:"lsaSeqNumber"`
	Checksum          string `json:"checksum"`
	Length            int    `json:"length"`
	NetworkMask       int    `json:"networkMask"`
	MetricType        string `json:"metricType"`
	Tos               int    `json:"tos"`
	Metric            int    `json:"metric"`
	ForwardAddress    string `json:"forwardAddress"`
	ExternalRouteTag  int    `json:"externalRouteTag"`
}

type OSPFNssaExternalData struct {
	RouterID           string                      `json:"routerId"`
	NSSAExternalStates map[string]NssaExternalArea `json:"NSSA External Link States"`
}

type NssaExternalArea map[string]NssaExternalLSA

type NssaExternalLSA struct {
	LSAAge             int    `json:"lsaAge"`
	Options            string `json:"options"`
	LSAFlags           int    `json:"lsaFlags"`
	LSAType            string `json:"lsaType"`
	LinkStateID        string `json:"linkStateId"`
	AdvertisingRouter  string `json:"advertisingRouter"`
	LSASeqNumber       string `json:"lsaSeqNumber"`
	Checksum           string `json:"checksum"`
	Length             int    `json:"length"`
	NetworkMask        int    `json:"networkMask"`
	MetricType         string `json:"metricType"`
	Tos                int    `json:"tos"`
	Metric             int    `json:"metric"`
	NSSAForwardAddress string `json:"nssaForwardAddress"`
	ExternalRouteTag   int    `json:"externalRouteTag"`
}

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

	fmt.Printf("Response: \n%+v\n", ospfRouterData)

	ospfNetworkData, err := FetchOSPFNetworkData(executor)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", ospfNetworkData)

	ospfSummaryData, err := FetchOSPFSummaryData(executor)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", ospfSummaryData)

	ospfAsbrSummaryData, err := FetchOSPFAsbrSummaryData(executor)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", ospfAsbrSummaryData)

	ospfExternalData, err := FetchOSPFExternalData(executor)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", ospfExternalData)

	ospfNssaExternalData, err := FetchOSPFNssaExternalData(executor)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", ospfNssaExternalData)

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
