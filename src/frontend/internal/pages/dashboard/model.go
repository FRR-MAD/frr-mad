package dashboard

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/viewport"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	title         string
	subTabs       []string
	ospfAnomalies []string
	windowSize    *common.WindowSize
	viewport      viewport.Model
	currentTime   time.Time
}

func New(windowSize *common.WindowSize) *Model {
	boxWidthForOne := windowSize.Width - 6
	// subtract tab row, footer, and border heights.
	outputHeight := windowSize.Height - styles.TabRowHeight - styles.FooterHeight - 2

	// Create the viewport with the desired dimensions.
	vp := viewport.New(boxWidthForOne, outputHeight)

	return &Model{
		title:         "Dashboard",
		subTabs:       []string{"OSPF", "TBD"},
		ospfAnomalies: []string{"Fetching OSPF data..."},
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
		"[r] refresh dashboard",
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

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		common.FetchOSPFData(),
		reloadView(),
	)
}
