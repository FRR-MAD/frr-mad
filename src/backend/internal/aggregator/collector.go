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

type OSPFDatabase struct {
	RouterID             string          `json:"routerId"`
	Areas                map[string]Area `json:"areas"`
	ASExternalLinkStates []ASExternalLSA `json:"asExternalLinkStates"`
	ASExternalCount      int             `json:"asExternalLinkStatesCount"`
}

type Area struct {
	RouterLinkStates           []RouterLSA      `json:"routerLinkStates"`
	RouterLinkStatesCount      int              `json:"routerLinkStatesCount"`
	NetworkLinkStates          []NetworkLSA     `json:"networkLinkStates"`
	NetworkLinkStatesCount     int              `json:"networkLinkStatesCount"`
	SummaryLinkStates          []SummaryLSA     `json:"summaryLinkStates"`
	SummaryLinkStatesCount     int              `json:"summaryLinkStatesCount"`
	ASBRSummaryLinkStates      []ASBRSummaryLSA `json:"asbrSummaryLinkStates"`
	ASBRSummaryLinkStatesCount int              `json:"asbrSummaryLinkStatesCount"`
}

type BaseLSA struct {
	LSID             string `json:"lsId"`
	AdvertisedRouter string `json:"advertisedRouter"`
	LSAAge           int    `json:"lsaAge"`
	SequenceNumber   string `json:"sequenceNumber"`
	Checksum         string `json:"checksum"`
}

type RouterLSA struct {
	BaseLSA
	NumOfRouterLinks int `json:"numOfRouterLinks"`
}

type NetworkLSA struct {
	BaseLSA
}

type SummaryLSA struct {
	BaseLSA
	SummaryAddress string `json:"summaryAddress"`
}

type ASBRSummaryLSA struct {
	BaseLSA
}

type ASExternalLSA struct {
	BaseLSA
	MetricType string `json:"metricType"` // E1 or E2
	Route      string `json:"route"`      // Prefix with mask
	Tag        int    `json:"tag"`
}

// ------------------------------------------------------------------------------------------------------------------------------------

type OSPFDuplicates struct {
	RouterID             string                `json:"routerId"`
	ASExternalLinkStates []ASExternalLinkState `json:"asExternalLinkStates"`
}

type ASExternalLinkState struct {
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
	TOS               int    `json:"tos"`
	Metric            int    `json:"metric"`
	ForwardAddress    string `json:"forwardAddress"`
	ExternalRouteTag  int    `json:"externalRouteTag"`
}

// ------------------------------------------------------------------------------------------------------------------------------------

type OSPFNeighbors struct {
	Neighbors map[string][]Neighbor `json:"neighbors"`
}
type Neighbor struct {
	Priority                           int    `json:"priority"`
	State                              string `json:"state"`
	NbrPriority                        int    `json:"nbrPriority"`
	NbrState                           string `json:"nbrState"`
	Converged                          string `json:"converged"`
	Role                               string `json:"role"`
	UpTimeInMsec                       int64  `json:"upTimeInMsec"`
	DeadTimeMsecs                      int    `json:"deadTimeMsecs"`
	RouterDeadIntervalTimerDueMsec     int    `json:"routerDeadIntervalTimerDueMsec"`
	UpTime                             string `json:"upTime"`
	DeadTime                           string `json:"deadTime"`
	Address                            string `json:"address"`
	IfaceAddress                       string `json:"ifaceAddress"`
	IfaceName                          string `json:"ifaceName"`
	RetransmitCounter                  int    `json:"retransmitCounter"`
	LinkStateRetransmissionListCounter int    `json:"linkStateRetransmissionListCounter"`
	RequestCounter                     int    `json:"requestCounter"`
	LinkStateRequestListCounter        int    `json:"linkStateRequestListCounter"`
	DbSummaryCounter                   int    `json:"dbSummaryCounter"`
	DatabaseSummaryListCounter         int    `json:"databaseSummaryListCounter"`
}

// ------------------------------------------------------------------------------------------------------------------------------------

type InterfaceList map[string]Interface

type Interface struct {
	AdministrativeStatus string      `json:"administrativeStatus"`
	OperationalStatus    string      `json:"operationalStatus"`
	LinkDetection        bool        `json:"linkDetection"`
	LinkUps              int         `json:"linkUps"`
	LinkDowns            int         `json:"linkDowns"`
	LastLinkUp           string      `json:"lastLinkUp,omitempty"`
	LastLinkDown         string      `json:"lastLinkDown,omitempty"`
	VrfName              string      `json:"vrfName"`
	MplsEnabled          bool        `json:"mplsEnabled"`
	LinkDown             bool        `json:"linkDown"`
	LinkDownV6           bool        `json:"linkDownV6"`
	McForwardingV4       bool        `json:"mcForwardingV4"`
	McForwardingV6       bool        `json:"mcForwardingV6"`
	PseudoInterface      bool        `json:"pseudoInterface"`
	Index                int         `json:"index"`
	Metric               int         `json:"metric"`
	Mtu                  int         `json:"mtu"`
	Speed                int         `json:"speed"`
	Flags                string      `json:"flags"`
	Type                 string      `json:"type"`
	HardwareAddress      string      `json:"hardwareAddress,omitempty"`
	IpAddresses          []IpAddress `json:"ipAddresses"`
	InterfaceType        string      `json:"interfaceType"`
	InterfaceSlaveType   string      `json:"interfaceSlaveType"`
	LacpBypass           bool        `json:"lacpBypass"`
	EvpnMh               EvpnMh      `json:"evpnMh"`
	Protodown            string      `json:"protodown"`
	ParentIfindex        int         `json:"parentIfindex,omitempty"`
}

type IpAddress struct {
	Address    string `json:"address"`
	Secondary  bool   `json:"secondary"`
	Unnumbered bool   `json:"unnumbered"`
}

type EvpnMh struct {
	EthernetSegmentID string `json:"ethernetSegmentId,omitempty"`
	ESI               string `json:"esi,omitempty"`
	DFPreference      int    `json:"dfPreference,omitempty"`
	DFAlgorithm       string `json:"dfAlgorithm,omitempty"`
	DFStatus          string `json:"dfStatus,omitempty"`
	MultiHomingMode   string `json:"multihomingMode,omitempty"`
	ActiveMode        bool   `json:"activeMode,omitempty"`
	BypassMode        bool   `json:"bypassMode,omitempty"`
	LocalBias         bool   `json:"localBias,omitempty"`
	FastFailover      bool   `json:"fastFailover,omitempty"`
	UpTime            string `json:"upTime,omitempty"`
	BGPStatus         string `json:"bgpStatus,omitempty"`
	ProtocolStatus    string `json:"protocolStatus,omitempty"`
	ProtocolDown      bool   `json:"protocolDown,omitempty"`
	MacCount          int    `json:"macCount,omitempty"`
	LocalIfindex      int    `json:"localIfindex,omitempty"`
	NetworkCount      int    `json:"networkCount,omitempty"`
	JoinCount         int    `json:"joinCount,omitempty"`
	LeaveCount        int    `json:"leaveCount,omitempty"`
}

// ------------------------------------------------------------------------------------------------------------------------------------

type RouteList map[string][]Route

type Route struct {
	Prefix                   string    `json:"prefix"`
	PrefixLen                int       `json:"prefixLen"`
	Protocol                 string    `json:"protocol"`
	VrfID                    int       `json:"vrfId"`
	VrfName                  string    `json:"vrfName"`
	Selected                 bool      `json:"selected,omitempty"`
	DestSelected             bool      `json:"destSelected,omitempty"`
	Distance                 int       `json:"distance"`
	Metric                   int       `json:"metric"`
	Installed                bool      `json:"installed,omitempty"`
	Table                    int       `json:"table"`
	InternalStatus           int       `json:"internalStatus"`
	InternalFlags            int       `json:"internalFlags"`
	InternalNextHopNum       int       `json:"internalNextHopNum"`
	InternalNextHopActiveNum int       `json:"internalNextHopActiveNum"`
	NexthopGroupID           int       `json:"nexthopGroupId"`
	InstalledNexthopGroupID  int       `json:"installedNexthopGroupId,omitempty"`
	Uptime                   string    `json:"uptime"`
	Nexthops                 []Nexthop `json:"nexthops"`
}

type Nexthop struct {
	Flags             int    `json:"flags"`
	Fib               bool   `json:"fib,omitempty"`
	DirectlyConnected bool   `json:"directlyConnected,omitempty"`
	Duplicate         bool   `json:"duplicate,omitempty"`
	IP                string `json:"ip,omitempty"`
	Afi               string `json:"afi,omitempty"`
	InterfaceIndex    int    `json:"interfaceIndex"`
	InterfaceName     string `json:"interfaceName"`
	Active            bool   `json:"active"`
	Weight            int    `json:"weight,omitempty"`
}

// ------------------------------------------------------------------------------------------------------------------------------------

type Collector struct {
	fetcher    *Fetcher
	configPath string
	socketPath string
	logger     *logger.Logger
	cache      *frrProto.CombinedState
}

func NewFRRCommandExecutor(socketDir string, timeout time.Duration) *frrSocket.FRRCommandExecutor {
	return &frrSocket.FRRCommandExecutor{
		DirPath: socketDir,
		Timeout: timeout,
	}
}

func NewCollector(metricsURL, configPath, socketPath string, logger *logger.Logger) *Collector {
	return &Collector{
		fetcher:    NewFetcher(metricsURL),
		configPath: configPath,
		socketPath: socketPath,
		logger:     logger,
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
	fmt.Printf("Response of FetchStaticFRRConfig(): \n%+v\n", staticFRRConfigParsed)

	ospfRouterData, err := FetchOSPFRouterData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", ospfRouterData)

	ospfNetworkData, err := FetchOSPFNetworkData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", ospfNetworkData)

	ospfSummaryData, err := FetchOSPFSummaryData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", ospfSummaryData)

	ospfAsbrSummaryData, err := FetchOSPFAsbrSummaryData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", ospfAsbrSummaryData)

	ospfExternalData, err := FetchOSPFExternalData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", ospfExternalData)

	ospfNssaExternalData, err := FetchOSPFNssaExternalData(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", ospfNssaExternalData)

	out1, err := FetchFullOSPFDatabase(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", out1)

	out2, err := FetchOSPFDuplicateCandidates(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", out2)

	out3, err := FetchOSPFNeighbors(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", out3)

	out4, err := FetchInterfaceStatus(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", out4)

	out5, err := FetchExpectedRoutes(executor)
	if err != nil {
		//fmt.Print(err)
		c.logger.Error(err.Error())
		//os.Exit(1)
	}

	fmt.Printf("Response: \n%+v\n", out5)

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
