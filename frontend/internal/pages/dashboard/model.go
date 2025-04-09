package dashboard

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Title         string
	SubTabs       []string
	ospfAnomalies []string
	windowSize    *common.WindowSize
}

func New(windowSize *common.WindowSize) *Model {
	return &Model{
		Title:         "Dashboard",
		SubTabs:       []string{"Dashboard 1", "Dashboard 2"},
		ospfAnomalies: []string{"Fetching OSPF data..."},
		windowSize:    windowSize,
	}
}

func (m *Model) GetTitle() common.Tab {
	return common.Tab{
		Title:   m.Title,
		SubTabs: m.SubTabs,
	}
}

func (m *Model) GetSubTabsLength() int {
	return len(m.SubTabs)
}

func (m *Model) GetFooterOptions() common.FooterOption {
	keyBoardOptions := []string{
		"'r': refresh dashboard",
		"'e': export everything",
	}
	return common.FooterOption{
		PageTitle:   m.Title,
		PageOptions: keyBoardOptions,
	}
}

func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		common.FetchOSPFData(),
	)
}
