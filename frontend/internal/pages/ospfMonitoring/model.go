package ospfMonitoring

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

// Model defines the state for the dashboard page.
type Model struct {
	Title      string
	SubTabs    []string
	windowSize *common.WindowSize
}

// New creates and returns a new dashboard Model.
func New(windowSize *common.WindowSize) *Model {
	return &Model{
		Title:      "OSPF Monitoring",
		SubTabs:    []string{"OSPF Tab 1", "OSPF Tab 2", "OSPF Tab 3"},
		windowSize: windowSize,
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
		"'r': refresh OSPF monitoring",
		"'e': export OSPF data",
	}
	return common.FooterOption{
		PageTitle:   m.Title,
		PageOptions: keyBoardOptions,
	}
}

// Init returns the initial command (none in this case).
func (m *Model) Init() tea.Cmd {
	return nil
}
