package ospfMonitoring

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Model defines the state for the dashboard page.
type Model struct {
	title         string
	subTabs       []string
	runningConfig []string
	windowSize    *common.WindowSize
	viewport      viewport.Model
}

// New creates and returns a new dashboard Model.
func New(windowSize *common.WindowSize) *Model {
	boxWidthForOne := windowSize.Width - 8
	if boxWidthForOne < 20 {
		boxWidthForOne = 20
	}
	// subtract tab row, footer, and border heights.
	outputHeight := windowSize.Height - styles.TabRowHeight - styles.FooterHeight - 2

	// Create the viewport with the desired dimensions.
	vp := viewport.New(boxWidthForOne, outputHeight)

	return &Model{
		title: "OSPF Monitoring",
		// '9 - Running Config' has to remain last in the list
		// because the key '9' is mapped to the last element of the list.
		subTabs:       []string{"Advertisement", "Router LSAs", "TBD", "TBD", "Running Config"},
		runningConfig: []string{"Fetching running config..."},
		windowSize:    windowSize,
		viewport:      vp,
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
		"[e] export OSPF data",
	}
	return common.FooterOption{
		PageTitle:   m.title,
		PageOptions: keyBoardOptions,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		common.FetchRunningConfig(),
	)
}
