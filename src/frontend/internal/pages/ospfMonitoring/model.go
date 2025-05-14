package ospfMonitoring

import (
	"github.com/ba2025-ysmprc/frr-mad/src/logger"
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Model defines the state for the dashboard page.
type Model struct {
	title         string
	subTabs       []string
	footer        []string
	runningConfig []string
	expandedMode  bool
	windowSize    *common.WindowSize
	viewport      viewport.Model
	logger        *logger.Logger
}

// New creates and returns a new dashboard Model.
func New(windowSize *common.WindowSize, appLogger *logger.Logger) *Model {

	// Create the viewport with the desired dimensions.
	vp := viewport.New(styles.ViewPortWidthCompletePage, styles.ViewPortHeightCompletePage)

	return &Model{
		title: "OSPF Monitoring",
		// 'Running Config' has to remain last in the list
		// because the key '9' is mapped to the last element of the list.
		subTabs:       []string{"LSDB", "Router LSAs", "Network LSAs", "External LSAs", "Neighbors", "Running Config"},
		footer:        []string{"[r] refresh", "[↑/↓] scroll", "[e] export OSPF data"},
		runningConfig: []string{"Fetching running config..."},
		expandedMode:  false,
		windowSize:    windowSize,
		viewport:      vp,
		logger:        appLogger,
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

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		common.FetchRunningConfig(m.logger),
	)
}
