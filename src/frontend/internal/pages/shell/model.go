package shell

import (
	"github.com/ba2025-ysmprc/frr-mad/src/logger"
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
)

// Model defines the state for the shell page.
type Model struct {
	title               string
	subTabs             []string
	windowSize          *common.WindowSize
	activeShell         string
	bashInput           string
	bashOutput          string
	vtyshInput          string
	vtyshOutput         string
	backendServiceInput string
	backendCommandInput string
	activeBackendInput  string
	backendResponse     string
	viewport            viewport.Model
	logger              *logger.Logger
}

// New creates and returns a new dashboard Model.
func New(windowSize *common.WindowSize, appLogger *logger.Logger) *Model {
	boxWidthForOne := windowSize.Width - 10
	if boxWidthForOne < 20 {
		boxWidthForOne = 20
	}
	// For example: subtract tab row, footer, and input area heights.
	outputHeight := windowSize.Height - 6 - 1 - 3 - 2

	// Create the viewport with the desired dimensions.
	vp := viewport.New(boxWidthForOne, outputHeight)

	return &Model{
		title:               "Shell",
		subTabs:             []string{"bash", "vtysh", "Backend Test"},
		windowSize:          windowSize,
		activeShell:         "",
		backendServiceInput: "",
		backendCommandInput: "",
		activeBackendInput:  "service",
		viewport:            vp,
		logger:              appLogger,
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

func (m *Model) ClearInput() {
	m.bashInput = ""
	m.vtyshInput = ""
}

func (m *Model) clearBackendInput() {
	m.backendServiceInput = ""
	m.backendCommandInput = ""
}

func (m *Model) ClearOutput() {
	m.bashOutput = ""
	m.vtyshOutput = ""
}

func (m *Model) GetFooterOptions() common.FooterOption {
	keyBoardOptions := []string{
		"[enter]: execute command",
		"[↑/↓] scroll",
		"[backspace]: delete last character",
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
