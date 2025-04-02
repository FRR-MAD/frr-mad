package dashboard

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"

	tea "github.com/charmbracelet/bubbletea"
)

type Model struct {
	Title         string
	ospfAnomalies []string
	windowSize    *common.WindowSize
}

func New(windowSize *common.WindowSize) Model {
	return Model{
		Title:         "Dashboard",
		ospfAnomalies: []string{"Fetching OSPF data..."},
		windowSize:    windowSize,
	}
}

func (m Model) GetTitle() string {
	return m.Title
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(
		common.FetchOSPFData(),
	)
}
