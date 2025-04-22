package dashboard

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
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
	return &Model{
		title:         "1 - Dashboard",
		subTabs:       []string{"OSPF", "Dashboard 2"},
		ospfAnomalies: []string{"Fetching OSPF data..."},
		windowSize:    windowSize,
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
		"'r': refresh dashboard",
		"'e': export everything",
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
