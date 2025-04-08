package shell

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

// Model defines the state for the shell page.
type Model struct {
	Title       string
	SubTabs     []string
	windowSize  *common.WindowSize
	ActiveShell string
	BashInput   string
	BashOutput  string
	VtyshInput  string
	VtyshOutput string
}

// New creates and returns a new dashboard Model.
func New(windowSize *common.WindowSize) *Model {
	return &Model{
		Title:       "Shell",
		SubTabs:     []string{"bash", "vtysh"},
		windowSize:  windowSize,
		ActiveShell: "",
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

func (m *Model) ClearInput() {
	m.BashInput = ""
	m.VtyshInput = ""
}

func (m *Model) ClearOutput() {
	m.BashOutput = ""
	m.VtyshOutput = ""
}

// Init returns the initial command (none in this case).
func (m *Model) Init() tea.Cmd {
	return nil
}
