package ospfMonitoring

import (
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/frr-mad/frr-tui/internal/common"
	backend "github.com/frr-mad/frr-tui/internal/services"
	"github.com/frr-mad/frr-tui/internal/ui/styles"
	"github.com/frr-mad/frr-tui/internal/ui/toast"
	"google.golang.org/protobuf/proto"
)

// Model defines the state for the dashboard page.
type Model struct {
	title             string
	subTabs           []string
	footer            []string
	readOnlyMode      bool
	toast             toast.Model
	cursor            int
	exportOptions     []common.ExportOption
	exportData        map[string]string
	exportDirectory   string
	runningConfig     []string
	expandedMode      bool // TODO: not used
	showExportOverlay bool
	textFilter        *common.Filter
	windowSize        *common.WindowSize
	viewport          viewport.Model
	viewportRightHalf viewport.Model
	statusMessage     string
	statusSeverity    styles.StatusSeverity
	statusTimer       time.Time
	statusDuration    time.Duration
	logger            *logger.Logger
}

// New creates and returns a new dashboard Model.
func New(windowSize *common.WindowSize, appLogger *logger.Logger, exportPath string) *Model {

	// Create the viewport with the desired dimensions.
	vp := viewport.New(styles.WidthViewPortCompletePage,
		styles.HeightViewPortCompletePage-styles.BodyFooterHeight)
	vprh := viewport.New(styles.WidthViewPortHalf,
		styles.HeightViewPortCompletePage-styles.HeightH1-styles.AdditionalFooterHeight)

	return &Model{
		title: "OSPF Monitoring",
		// 'Running Config' has to remain last in the list
		// because the key '9' is mapped to the last element of the list.
		subTabs:           []string{"LSDB", "Router LSAs", "Network LSAs", "External LSAs", "Neighbors", "Running Config"},
		footer:            []string{"[↑ ↓ home end] scroll", "[ctrl+e] export options", "[ctrl+r] refresh"},
		readOnlyMode:      true,
		cursor:            0,
		exportOptions:     []common.ExportOption{},
		exportData:        make(map[string]string),
		exportDirectory:   exportPath,
		runningConfig:     []string{"Fetching running config..."},
		expandedMode:      false,
		showExportOverlay: false,
		windowSize:        windowSize,
		viewport:          vp,
		viewportRightHalf: vprh,
		statusMessage:     "",
		statusSeverity:    styles.SeverityInfo,
		logger:            appLogger,
	}
}

func (m *Model) GetTitle() common.Tab {
	return common.Tab{
		Title:   m.title,
		SubTabs: m.subTabs,
	}
}

func (m *Model) GetSubTabsLength() int {
	return len(m.subTabs)
}

func (m *Model) GetFooterOptions() common.FooterOption {
	keyBoardOptions := m.footer
	return common.FooterOption{
		PageTitle:   m.title,
		PageOptions: keyBoardOptions,
	}
}

func (m *Model) setTimedStatus(message string, severity styles.StatusSeverity, duration time.Duration) {
	m.statusMessage = message
	m.statusSeverity = severity
	m.statusTimer = time.Now()
	m.statusDuration = duration
}

// fetchLatestData fetches all data from the backend that are possible to export from the ospf monitor exporter
func (m *Model) fetchLatestData() error {
	items := []struct {
		key, label, filename string
		fetch                func() (proto.Message, error)
	}{
		{
			key:      "GetLSDB",
			label:    "complete link-state database",
			filename: "link-state_database.json",
			fetch:    func() (proto.Message, error) { return backend.GetLSDB(m.logger) },
		},
		{
			key:      "GetOspfNeighbors",
			label:    "ospf neighbors",
			filename: "ospf_neighbors.json",
			fetch:    func() (proto.Message, error) { return backend.GetOspfNeighbors(m.logger) },
		},
		{
			key:      "GetOspfRouterDataSelf",
			label:    "lsdb type 1 router self-originating",
			filename: "lsdb_router_self.json",
			fetch:    func() (proto.Message, error) { return backend.GetOspfRouterDataSelf(m.logger) },
		},
		{
			key:      "GetOspfNetworkDataSelf",
			label:    "lsdb type 2 network self-originating",
			filename: "lsdb_network_self.json",
			fetch:    func() (proto.Message, error) { return backend.GetOspfNetworkDataSelf(m.logger) },
		},
		{
			key:      "GetOspfSummaryDataSelf",
			label:    "lsdb type 3 summary self-originating",
			filename: "lsdb_summary_self.json",
			fetch:    func() (proto.Message, error) { return backend.GetOspfSummaryDataSelf(m.logger) },
		},
		{
			key:      "GetOspfAsbrSummaryDataSelf",
			label:    "lsdb type 4 asbr summary self-originating",
			filename: "lsdb_asbr_summary_self.json",
			fetch:    func() (proto.Message, error) { return backend.GetOspfAsbrSummaryDataSelf(m.logger) },
		},
		{
			key:      "GetOspfExternalDataSelf",
			label:    "lsdb type 5 external self-originating",
			filename: "lsdb_external_self.json",
			fetch:    func() (proto.Message, error) { return backend.GetOspfExternalDataSelf(m.logger) },
		},
		{
			key:      "GetOspfNssaExternalDataSelf",
			label:    "lsdb type 7 nssa external self-originating",
			filename: "lsdb_nssa_external_self.json",
			fetch:    func() (proto.Message, error) { return backend.GetOspfNssaExternalDataSelf(m.logger) },
		},
		{
			key:      "GetOspfP2PInterfaceMapping",
			label:    "mapping of P2P interfaces",
			filename: "p2p_mapping.json",
			fetch:    func() (proto.Message, error) { return backend.GetOspfP2PInterfaceMapping(m.logger) },
		},
		{
			key:      "GetStaticFRRConfiguration",
			label:    "parsed frr configuration",
			filename: "frr_configuration.json",
			fetch:    func() (proto.Message, error) { return backend.GetStaticFRRConfiguration(m.logger) },
		},
	}

	var err error
	for _, it := range items {
		m.exportOptions, err = common.ExportProto(
			m.exportData,
			m.exportOptions,
			it.key, it.label, it.filename,
			it.fetch,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		common.FetchRunningConfig(m.logger),
	)
}
