package ospfMonitoring

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

// Model defines the state for the dashboard page.
type Model struct {
	title      string
	subTabs    []string
	windowSize *common.WindowSize
}

// New creates and returns a new dashboard Model.
func New(windowSize *common.WindowSize) *Model {
	return &Model{
		title: "2 - OSPF Monitoring",
		// '9 - Running Config' has to remain last in the list
		// because the key '9' is mapped to the last element of the list.
		subTabs:    []string{"1 - Advertisement", "2 - TBD", "3 - TBD", "9 - Running Config"},
		windowSize: windowSize,
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
		"'r': refresh OSPF monitoring",
		"'e': export OSPF data",
	}
	return common.FooterOption{
		PageTitle:   m.title,
		PageOptions: keyBoardOptions,
	}
}

// Init returns the initial command (none in this case).
func (m *Model) Init() tea.Cmd {
	return nil
}
