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

type Model struct {
	title              string
	subTabs            []string
	footer             []string
	ospfAnomalies      []string
	hasAnomalyDetected bool
	showAnomalyOverlay bool
	windowSize         *common.WindowSize
	viewportLeft       viewport.Model
	viewportRight      viewport.Model
	currentTime        time.Time
	logger             *logger.Logger
}

func New(windowSize *common.WindowSize, appLogger *logger.Logger) *Model {

	// Create the viewportLeft with the desired dimensions.
	vpl := viewport.New(styles.ViewPortWidthThreeFourth, styles.ViewPortHeightCompletePage-styles.HeightH1)
	vpr := viewport.New(styles.ViewPortWidthOneFourth, styles.ViewPortHeightCompletePage-styles.HeightH1)

	return &Model{
		title:              "Dashboard",
		subTabs:            []string{"OSPF", "BGP"},
		footer:             []string{"[r] refresh", "[↑/↓] scroll", "[a] show/hide anomaly details"},
		ospfAnomalies:      []string{"Fetching OSPF data..."},
		hasAnomalyDetected: false,
		showAnomalyOverlay: false,
		windowSize:         windowSize,
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

func (m *Model) Init() tea.Cmd {
	m.detectAnomaly()
	return tea.Batch(
		common.FetchOSPFData(m.logger),
		reloadView(),
	)
}
