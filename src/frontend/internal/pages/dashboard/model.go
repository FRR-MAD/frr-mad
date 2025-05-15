package dashboard

import (
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/logger"
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	backend "github.com/ba2025-ysmprc/frr-tui/internal/services"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/viewport"

	tea "github.com/charmbracelet/bubbletea"
)

type ExportOption struct {
	Label    string
	MapKey   string
	Filename string
}

type Model struct {
	title              string
	subTabs            []string
	footer             []string
	cursor             int
	exportOptions      []ExportOption
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
	currentTime        time.Time
	logger             *logger.Logger
}

func New(windowSize *common.WindowSize, appLogger *logger.Logger) *Model {

	// Create the viewports with the desired dimensions.
	vp := viewport.New(styles.ViewPortWidthCompletePage, styles.ViewPortHeightCompletePage)
	vpl := viewport.New(styles.ViewPortWidthThreeFourth, styles.ViewPortHeightCompletePage-styles.HeightH1)
	vpr := viewport.New(styles.ViewPortWidthOneFourth, styles.ViewPortHeightCompletePage-styles.HeightH1)

	return &Model{
		title:              "Dashboard",
		subTabs:            []string{"OSPF", "BGP"},
		footer:             []string{"[e] export options", "[r] refresh", "[↑/↓] scroll", "[a] anomaly details"},
		cursor:             0,
		exportOptions:      []ExportOption{},
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
	m.exportData["GetRouterAnomalies"] = common.PrettyPrintJSON(ospfRouterAnomalies)
	m.exportOptions = addOption(m.exportOptions, ExportOption{
		Label:    "anomalies - router (LSA type 1)",
		MapKey:   "GetRouterAnomalies",
		Filename: "type1_router_anomalies",
	})

	ospfExternalAnomalies, _ := backend.GetExternalAnomalies(m.logger)
	m.exportData["GetExternalAnomalies"] = common.PrettyPrintJSON(ospfExternalAnomalies)
	m.exportOptions = addOption(m.exportOptions, ExportOption{
		Label:    "anomalies - external (LSA type 5)",
		MapKey:   "GetExternalAnomalies",
		Filename: "type5_external_anomalies",
	})

	ospfNSSAExternalAnomalies, _ := backend.GetNSSAExternalAnomalies(m.logger)
	m.exportData["GetNSSAExternalAnomalies"] = common.PrettyPrintJSON(ospfNSSAExternalAnomalies)
	m.exportOptions = addOption(m.exportOptions, ExportOption{
		Label:    "anomalies - nssa external (LSA type 7)",
		MapKey:   "GetNSSAExternalAnomalies",
		Filename: "type7_nssa_external_anomalies",
	})

	if common.HasAnyAnomaly(ospfRouterAnomalies) ||
		common.HasAnyAnomaly(ospfExternalAnomalies) ||
		common.HasAnyAnomaly(ospfNSSAExternalAnomalies) {

		m.hasAnomalyDetected = true
	} else {
		m.hasAnomalyDetected = false
	}
}

func (m *Model) Init() tea.Cmd {
	m.detectAnomaly()
	return tea.Batch(
		common.FetchOSPFData(m.logger),
		reloadView(),
	)
}

// addOption adds opt to slice only if no existing entry has the same MapKey.
func addOption(opts []ExportOption, opt ExportOption) []ExportOption {
	for _, e := range opts {
		if e.MapKey == opt.MapKey {
			return opts // already present
		}
	}
	return append(opts, opt)
}
