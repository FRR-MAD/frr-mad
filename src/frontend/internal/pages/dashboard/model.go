package dashboard

import (
	"time"

	"github.com/frr-mad/frr-tui/internal/ui/toast"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/frr-mad/frr-tui/internal/common"
	backend "github.com/frr-mad/frr-tui/internal/services"
	"github.com/frr-mad/frr-tui/internal/ui/styles"
	"google.golang.org/protobuf/proto"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	title              string
	subTabs            []string
	footer             []string
	readOnlyMode       bool
	toast              toast.Model
	cursor             int
	exportOptions      []common.ExportOption
	exportData         map[string]string
	exportDirectory    string
	hasAnomalyDetected bool
	showAnomalyOverlay bool
	showExportOverlay  bool
	textFilter         *common.Filter
	windowSize         *common.WindowSize
	viewport           viewport.Model
	viewportLeft       viewport.Model
	viewportRight      viewport.Model
	viewportRightHalf  viewport.Model
	currentTime        time.Time
	statusMessage      string
	statusSeverity     styles.StatusSeverity
	statusTimer        time.Time
	statusDuration     time.Duration
	logger             *logger.Logger
}

func New(windowSize *common.WindowSize, appLogger *logger.Logger, exportPath string) *Model {

	// Create the viewports with the desired dimensions.
	vp := viewport.New(styles.WidthViewPortCompletePage, styles.HeightViewPortCompletePage)
	vpl := viewport.New(styles.WidthViewPortThreeFourth,
		styles.HeightViewPortCompletePage-styles.HeightH1-styles.BodyFooterHeight)
	vpr := viewport.New(styles.WidthViewPortOneFourth,
		styles.HeightViewPortCompletePage-styles.BodyFooterHeight)
	vprh := viewport.New(styles.WidthViewPortHalf,
		styles.HeightViewPortCompletePage-styles.HeightH1-styles.AdditionalFooterHeight)

	return &Model{
		title:              "Dashboard",
		subTabs:            []string{"Anomalies", "OSPF"},
		footer:             []string{"[↑ ↓ home end] scroll", "[ctrl+e] export options", "[ctrl+a] anomaly details"},
		readOnlyMode:       true,
		cursor:             0,
		exportOptions:      []common.ExportOption{},
		exportData:         make(map[string]string),
		exportDirectory:    exportPath,
		hasAnomalyDetected: false,
		showAnomalyOverlay: false,
		showExportOverlay:  false,

		windowSize:        windowSize,
		viewport:          vp,
		viewportLeft:      vpl,
		viewportRight:     vpr,
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

func reloadView() tea.Cmd {
	return tea.Tick(time.Second, func(t time.Time) tea.Msg {
		return common.ReloadMessage(t)
	})
}

func (m *Model) detectAnomaly() {
	ospfRouterAnomalies, _ := backend.GetRouterAnomalies(m.logger)
	ospfExternalAnomalies, _ := backend.GetExternalAnomalies(m.logger)
	ospfNSSAExternalAnomalies, _ := backend.GetNSSAExternalAnomalies(m.logger)
	ospfLSDBToRibAnomalies, _ := backend.GetLSDBToRibAnomalies(m.logger)
	ribToFibAnomalies, _ := backend.GetRibToFibAnomalies(m.logger)

	if common.HasAnyAnomaly(ospfRouterAnomalies) ||
		common.HasAnyAnomaly(ospfExternalAnomalies) ||
		common.HasAnyAnomaly(ospfNSSAExternalAnomalies) ||
		common.HasAnyAnomaly(ospfLSDBToRibAnomalies) ||
		common.HasAnyAnomaly(ribToFibAnomalies) {

		m.hasAnomalyDetected = true
	} else {
		m.hasAnomalyDetected = false
	}
}

func (m *Model) setTimedStatus(message string, severity styles.StatusSeverity, duration time.Duration) {
	m.statusMessage = message
	m.statusSeverity = severity
	m.statusTimer = time.Now()
	m.statusDuration = duration
}

// fetchLatestData fetches all data from the backend that are possible to export from the dashboard exporter
func (m *Model) fetchLatestData() error {
	var err error

	items := []struct {
		key, label, filename string
		fetch                func() (proto.Message, error)
		withTimestamp        bool
	}{
		{
			key:      "GetRouterAnomalies",
			label:    "anomalies – router (LSA type 1)",
			filename: "type1_router_anomalies.json",
			fetch:    func() (proto.Message, error) { return backend.GetRouterAnomalies(m.logger) },
		},
		{
			key:      "GetExternalAnomalies",
			label:    "anomalies – external (LSA type 5)",
			filename: "type5_external_anomalies.json",
			fetch:    func() (proto.Message, error) { return backend.GetExternalAnomalies(m.logger) },
		},
		{
			key:      "GetNSSAExternalAnomalies",
			label:    "anomalies – nssa external (LSA type 7)",
			filename: "type7_nssa_anomalies.json",
			fetch:    func() (proto.Message, error) { return backend.GetNSSAExternalAnomalies(m.logger) },
		},
		{
			key:      "GetOSPF",
			label:    "summary of the current OSPF router",
			filename: "general_ospf_information.json",
			fetch:    func() (proto.Message, error) { return backend.GetOSPF(m.logger) },
		},
		{
			key:      "GetLSDB",
			label:    "complete Link-State Database",
			filename: "link-state_database.json",
			fetch:    func() (proto.Message, error) { return backend.GetLSDB(m.logger) },
		},
		{
			key:      "GetParsedShouldStates",
			label:    "predicted should state of lsdb",
			filename: "should_state_lsdb.json",
			fetch:    func() (proto.Message, error) { return backend.GetParsedShouldStates(m.logger) },
		},
	}

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
	m.detectAnomaly()
	return tea.Batch(
		common.FetchOSPFData(m.logger),
		reloadView(),
	)
}
