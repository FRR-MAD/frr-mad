package rib

import (
	"time"

	"github.com/charmbracelet/bubbles/viewport"
	"github.com/frr-mad/frr-mad/src/frontend/internal/common"
	backend "github.com/frr-mad/frr-mad/src/frontend/internal/services"
	"github.com/frr-mad/frr-mad/src/frontend/internal/ui/styles"
	"github.com/frr-mad/frr-mad/src/frontend/internal/ui/toast"
	"github.com/frr-mad/frr-mad/src/logger"

	tea "github.com/charmbracelet/bubbletea"
	"google.golang.org/protobuf/proto"
)

type Model struct {
	appState          common.AppState
	title             string
	subTabs           []string
	footer            []string
	readOnlyMode      bool
	toast             toast.Model
	cursor            int
	exportOptions     []common.ExportOption
	exportData        map[string]string
	exportDirectory   string
	showExportOverlay bool
	windowSize        *common.WindowSize
	viewport          viewport.Model
	viewportRightHalf viewport.Model
	textFilter        *common.Filter
	statusMessage     string
	statusSeverity    styles.StatusSeverity
	statusTimer       time.Time
	statusDuration    time.Duration
	logger            *logger.Logger
}

func New(windowSize *common.WindowSize, appLogger *logger.Logger, exportPath string) *Model {

	// Create the viewport with the desired dimensions.
	vp := viewport.New(styles.WidthViewPortCompletePage,
		styles.HeightViewPortCompletePage-styles.BodyFooterHeight)
	vprh := viewport.New(styles.WidthViewPortHalf,
		styles.HeightViewPortCompletePage-styles.HeightH1-styles.AdditionalFooterHeight)

	return &Model{
		appState:          2,
		title:             "RIB",
		subTabs:           []string{"RIB", "FIB", "RIB-OSPF", "RIB-BGP", "RIB-Connected", "RIB-Static"},
		footer:            []string{"[↑ ↓ home end] scroll", "[ctrl+e] export options"},
		readOnlyMode:      true,
		cursor:            0,
		exportOptions:     []common.ExportOption{},
		exportData:        make(map[string]string),
		exportDirectory:   exportPath,
		showExportOverlay: false,
		windowSize:        windowSize,
		viewport:          vp,
		viewportRightHalf: vprh,
		statusMessage:     "",
		statusSeverity:    styles.SeverityInfo,
		logger:            appLogger,
	}
}

func (m *Model) GetAppState() common.AppState {
	return m.appState
}

func (m *Model) GetPageInfo() common.Tab {
	return common.Tab{
		Title:    m.title,
		SubTabs:  m.subTabs,
		AppState: m.appState,
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

// fetchLatestData fetches all data from the backend that are possible to export from the rib exporter
func (m *Model) fetchLatestData() error {
	items := []struct {
		key, label, filename string
		fetch                func() (proto.Message, error)
	}{
		{
			key:      "GetRIB",
			label:    "routing information base",
			filename: "rib.json",
			fetch:    func() (proto.Message, error) { return backend.GetRIB(m.logger) },
		},
		{
			key:      "GetRibFibSummary",
			label:    "summary data for rib and fib",
			filename: "rib_fib_summary.json",
			fetch:    func() (proto.Message, error) { return backend.GetRibFibSummary(m.logger) },
		},
	}

	var err error
	for _, it := range items {
		m.exportOptions, err = common.ExportProto(
			m.exportData,
			m.exportOptions,
			it.key,
			it.label,
			it.filename,
			it.fetch,
		)
		if err != nil {
			return err
		}
	}
	return nil
}

func (m *Model) Init() tea.Cmd {
	return nil
}
