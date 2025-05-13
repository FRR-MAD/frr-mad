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
	ospfAnomalies      []string
	hasAnomalyDetected bool
	windowSize         *common.WindowSize
	viewport           viewport.Model
	currentTime        time.Time
	logger             *logger.Logger
}

func New(windowSize *common.WindowSize, appLogger *logger.Logger) *Model {
	boxWidthForOne := windowSize.Width - 6
	// subtract tab row, footer, and border heights.
	outputHeight := windowSize.Height - styles.TabRowHeight - styles.FooterHeight - 2

	// Create the viewport with the desired dimensions.
	vp := viewport.New(boxWidthForOne, outputHeight)

	return &Model{
		title:              "Dashboard",
		subTabs:            []string{"OSPF", "BGP"},
		ospfAnomalies:      []string{"Fetching OSPF data..."},
		hasAnomalyDetected: false,
		windowSize:         windowSize,
		viewport:           vp,
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
	keyBoardOptions := []string{
		"[r] refresh",
		"[↑/↓] scroll",
		// "[e] export everything",
	}
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
