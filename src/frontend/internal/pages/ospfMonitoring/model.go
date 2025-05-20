package ospfMonitoring

import (
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/frr-mad/frr-tui/internal/common"
	backend "github.com/frr-mad/frr-tui/internal/services"
	"github.com/frr-mad/frr-tui/internal/ui/styles"
	"github.com/frr-mad/frr-tui/internal/ui/toast"
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
	filterQuery       string
	filterActive      bool
	filterInput       textinput.Model
	windowSize        *common.WindowSize
	viewport          viewport.Model
	viewportRightHalf viewport.Model
	logger            *logger.Logger
}

// New creates and returns a new dashboard Model.
func New(windowSize *common.WindowSize, appLogger *logger.Logger) *Model {

	// Create the viewport with the desired dimensions.
	vp := viewport.New(styles.ViewPortWidthCompletePage,
		styles.ViewPortHeightCompletePage-styles.FilterBoxHeight)
	vprh := viewport.New(styles.ViewPortWidthHalf,
		styles.ViewPortHeightCompletePage-styles.HeightH1-styles.AdditionalFooterHeight)

	// Create the text input for the filter
	ti := textinput.New()
	ti.Placeholder = "type to filter..."
	ti.CharLimit = 64
	ti.Width = 20

	return &Model{
		title: "OSPF Monitoring",
		// 'Running Config' has to remain last in the list
		// because the key '9' is mapped to the last element of the list.
		subTabs:           []string{"LSDB", "Router LSAs", "Network LSAs", "External LSAs", "Neighbors", "Running Config"},
		footer:            []string{"[e] export options", "[r] refresh", "[↑ ↓ home end] scroll", "[e] export OSPF data"},
		readOnlyMode:      true,
		cursor:            0,
		exportOptions:     []common.ExportOption{},
		exportData:        make(map[string]string),
		exportDirectory:   "/tmp/frr-mad/exports",
		runningConfig:     []string{"Fetching running config..."},
		expandedMode:      false,
		showExportOverlay: false,
		filterActive:      false,
		filterQuery:       "",
		filterInput:       ti,
		windowSize:        windowSize,
		viewport:          vp,
		viewportRightHalf: vprh,
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

// fetchLatestData fetches all data from the backend that are possible to export from the ospf monitor exporter
func (m *Model) fetchLatestData() error {
	lsdb, err := backend.GetLSDB(m.logger)
	if err != nil {
		return err
	}
	m.exportData["GetLSDB"] = common.PrettyPrintJSON(lsdb)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "complete link-state database",
		MapKey:   "GetLSDB",
		Filename: "link-state_database.json",
	})

	ospfNeighbors, err := backend.GetOspfNeighbors(m.logger)
	if err != nil {
		return nil
	}
	m.exportData["GetOspfNeighbors"] = common.PrettyPrintJSON(ospfNeighbors)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "ospf neighbors",
		MapKey:   "GetOspfNeighbors",
		Filename: "ospf_neighbors.json",
	})

	routerLSASelf, err := backend.GetOspfRouterDataSelf(m.logger)
	if err != nil {
		return nil
	}
	m.exportData["GetOspfRouterDataSelf"] = common.PrettyPrintJSON(routerLSASelf)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "lsdb type 1 router self-originating",
		MapKey:   "GetOspfRouterDataSelf",
		Filename: "lsdb_ruoter_self.json",
	})

	networkLSASelf, err := backend.GetOspfNetworkDataSelf(m.logger)
	if err != nil {
		return nil
	}
	m.exportData["GetOspfNetworkDataSelf"] = common.PrettyPrintJSON(networkLSASelf)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "lsdb type 2 network self-originating",
		MapKey:   "GetOspfNetworkDataSelf",
		Filename: "lsdb_network_self.json",
	})

	summaryLSASelf, err := backend.GetOspfSummaryDataSelf(m.logger)
	if err != nil {
		return nil
	}
	m.exportData["GetOspfSummaryDataSelf"] = common.PrettyPrintJSON(summaryLSASelf)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "lsdb type 3 summary self-originating",
		MapKey:   "GetOspfSummaryDataSelf",
		Filename: "lsdb_summary_self.json",
	})

	asbrSummaryLSASelf, err := backend.GetOspfAsbrSummaryDataSelf(m.logger)
	if err != nil {
		return nil
	}
	m.exportData["GetOspfAsbrSummaryDataSelf"] = common.PrettyPrintJSON(asbrSummaryLSASelf)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "lsdb type 4 asbr summary self-originating",
		MapKey:   "GetOspfAsbrSummaryDataSelf",
		Filename: "lsdb_asbr_summary_self.json",
	})

	externalLSASelf, err := backend.GetOspfExternalDataSelf(m.logger)
	if err != nil {
		return nil
	}
	m.exportData["GetOspfExternalDataSelf"] = common.PrettyPrintJSON(externalLSASelf)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "lsdb type 5 external self-originating",
		MapKey:   "GetOspfExternalDataSelf",
		Filename: "lsdb_external_self.json",
	})

	nssaExternalDataSelf, err := backend.GetOspfNssaExternalDataSelf(m.logger)
	if err != nil {
		return nil
	}
	m.exportData["GetOspfNssaExternalDataSelf"] = common.PrettyPrintJSON(nssaExternalDataSelf)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "lsdb type 7 nssa external self-originating",
		MapKey:   "GetOspfNssaExternalDataSelf",
		Filename: "lsdb_nssa_external_self.json",
	})

	p2pInterfaceMap, err := backend.GetOspfP2PInterfaceMapping(m.logger)
	if err != nil {
		return nil
	}
	m.exportData["GetOspfP2PInterfaceMapping"] = common.PrettyPrintJSON(p2pInterfaceMap)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "mapping of P2P Interfaces",
		MapKey:   "GetOspfP2PInterfaceMapping",
		Filename: "p2p_mapping.json",
	})

	staticFRRConfiguration, err := backend.GetStaticFRRConfiguration(m.logger)
	if err != nil {
		return nil
	}
	m.exportData["GetStaticFRRConfiguration"] = common.PrettyPrintJSON(staticFRRConfiguration)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "parsed frr configuration",
		MapKey:   "GetStaticFRRConfiguration",
		Filename: "frr_configuration.json",
	})

	return nil
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		common.FetchRunningConfig(m.logger),
	)
}
