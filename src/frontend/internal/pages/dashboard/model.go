package dashboard

import (
	"github.com/frr-mad/frr-tui/internal/ui/toast"
	"time"

	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/frr-mad/frr-tui/internal/common"
	backend "github.com/frr-mad/frr-tui/internal/services"
	"github.com/frr-mad/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/viewport"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	title              string
	subTabs            []string
	footer             []string
	toast              toast.Model
	cursor             int
	exportOptions      []common.ExportOption
	exportData         map[string]string
	exportDirectory    string
	ospfAnomalies      []string // to be deleted
	hasAnomalyDetected bool
	showAnomalyOverlay bool
	showExportOverlay  bool
	windowSize         *common.WindowSize
	viewport           viewport.Model
	viewportLeft       viewport.Model
	viewportRight      viewport.Model
	viewportRightHalf  viewport.Model
	currentTime        time.Time
	logger             *logger.Logger
}

func New(windowSize *common.WindowSize, appLogger *logger.Logger) *Model {

	// Create the viewports with the desired dimensions.
	vp := viewport.New(styles.ViewPortWidthCompletePage, styles.ViewPortHeightCompletePage)
	vpl := viewport.New(styles.ViewPortWidthThreeFourth, styles.ViewPortHeightCompletePage-styles.HeightH1)
	vpr := viewport.New(styles.ViewPortWidthOneFourth, styles.ViewPortHeightCompletePage-styles.HeightH1)
	vprh := viewport.New(styles.ViewPortWidthHalf, styles.ViewPortHeightCompletePage-styles.HeightH1)

	return &Model{
		title:              "Dashboard",
		subTabs:            []string{"OSPF", "BGP"},
		footer:             []string{"[e] export options", "[r] refresh", "[↑ ↓ home end] scroll", "[a] anomaly details"},
		cursor:             0,
		exportOptions:      []common.ExportOption{},
		exportData:         make(map[string]string),
		exportDirectory:    "/tmp/frr-mad/exports",
		ospfAnomalies:      []string{"Fetching OSPF data..."},
		hasAnomalyDetected: false,
		showAnomalyOverlay: false,
		showExportOverlay:  false,
		windowSize:         windowSize,
		viewport:           vp,
		viewportLeft:       vpl,
		viewportRight:      vpr,
		viewportRightHalf:  vprh,
		logger:             appLogger,
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

func reloadView() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return common.ReloadMessage(t)
	})
}

func (m *Model) detectAnomaly() {
	ospfRouterAnomalies, _ := backend.GetRouterAnomalies(m.logger)
	ospfExternalAnomalies, _ := backend.GetExternalAnomalies(m.logger)
	ospfNSSAExternalAnomalies, _ := backend.GetNSSAExternalAnomalies(m.logger)

	if common.HasAnyAnomaly(ospfRouterAnomalies) ||
		common.HasAnyAnomaly(ospfExternalAnomalies) ||
		common.HasAnyAnomaly(ospfNSSAExternalAnomalies) {

		m.hasAnomalyDetected = true
	} else {
		m.hasAnomalyDetected = false
	}
}

// fetchLatestData fetches all data from the backend that are possible to export from the dashboard exporter
func (m *Model) fetchLatestData() error {
	ospfRouterAnomalies, err := backend.GetRouterAnomalies(m.logger)
	if err != nil {
		return err
	}
	m.exportData["GetRouterAnomalies"] = common.PrettyPrintJSON(ospfRouterAnomalies)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "anomalies - router (LSA type 1)",
		MapKey:   "GetRouterAnomalies",
		Filename: "type1_router_anomalies.json",
	})

	ospfExternalAnomalies, err := backend.GetExternalAnomalies(m.logger)
	if err != nil {
		return err
	}
	m.exportData["GetExternalAnomalies"] = common.PrettyPrintJSON(ospfExternalAnomalies)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "anomalies - external (LSA type 5)",
		MapKey:   "GetExternalAnomalies",
		Filename: "type5_external_anomalies.json",
	})

	ospfNSSAExternalAnomalies, err := backend.GetNSSAExternalAnomalies(m.logger)
	if err != nil {
		return err
	}
	m.exportData["GetNSSAExternalAnomalies"] = common.PrettyPrintJSON(ospfNSSAExternalAnomalies)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "anomalies - nssa external (LSA type 7)",
		MapKey:   "GetNSSAExternalAnomalies",
		Filename: "type7_nssa_anomalies.json",
	})

	ospfInformation, err := backend.GetOSPF(m.logger)
	if err != nil {
		return err
	}
	m.exportData["GetOSPF"] = common.PrettyPrintJSON(ospfInformation)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "summary of the current OSPF router",
		MapKey:   "GetOSPF",
		Filename: "general_ospf_information.json",
	})

	lsdb, err := backend.GetLSDB(m.logger)
	if err != nil {
		return err
	}
	m.exportData["GetLSDB"] = common.PrettyPrintJSON(lsdb)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "complete Link-State Database",
		MapKey:   "GetLSDB",
		Filename: "link-state_database.json",
	})

	parsedShouldStates, err := backend.GetParsedShouldStates(m.logger)
	if err != nil {
		return err
	}
	//routerShould := parsedShouldStates.ShouldRouterLsdb
	//externalShould := parsedShouldStates.ShouldExternalLsdb
	//nssaExternalShould := parsedShouldStates.ShouldNssaExternalLsdb

	m.exportData["GetParsedShouldStates"] = common.PrettyPrintJSON(parsedShouldStates)
	m.exportOptions = common.AddExportOption(m.exportOptions, common.ExportOption{
		Label:    "predicted should state of lsdb",
		MapKey:   "GetParsedShouldStates",
		Filename: "should_state_lsdb.json",
	})

	return nil
}

func (m *Model) Init() tea.Cmd {
	m.detectAnomaly()
	return tea.Batch(
		common.FetchOSPFData(m.logger),
		reloadView(),
	)
}
