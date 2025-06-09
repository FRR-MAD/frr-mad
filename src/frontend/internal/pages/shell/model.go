package shell

import (
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/frr-mad/frr-mad/src/frontend/internal/common"
	"github.com/frr-mad/frr-mad/src/frontend/internal/ui/styles"
	"github.com/frr-mad/frr-mad/src/logger"
)

// Model defines the state for the shell page.
type Model struct {
	appState            common.AppState
	title               string
	subTabs             []string
	footer              []string
	readOnlyMode        bool
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
	statusMessage       string
	statusSeverity      styles.StatusSeverity
	logger              *logger.Logger
}

// New creates and returns a new dashboard Model.
func New(windowSize *common.WindowSize, appLogger *logger.Logger) *Model {
	// Create the viewport with the desired dimensions.
	vp := viewport.New(styles.WidthViewPortCompletePage, styles.HeightViewPortCompletePage-styles.HeightH1-2)

	return &Model{
		appState:            3,
		title:               "Shell",
		subTabs:             []string{"bash", "vtysh", "Backend Test"},
		footer:              []string{"[↑ ↓ home end] scroll", "[enter] execute command"},
		readOnlyMode:        true,
		windowSize:          windowSize,
		activeShell:         "",
		backendServiceInput: "",
		backendCommandInput: "",
		activeBackendInput:  "service",
		viewport:            vp,
		statusMessage:       "",
		statusSeverity:      styles.SeverityInfo,
		logger:              appLogger,
	}
}

func (m *Model) GetAppState() common.AppState {
	return m.appState
}

func (m *Model) GetPageInfo() common.Tab {
	return common.Tab{
		Title:    m.title,
		SubTabs:  m.subTabs,
		AppState: m.appState,
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
	keyBoardOptions := m.footer
	return common.FooterOption{
		PageTitle:   m.title,
		PageOptions: keyBoardOptions,
	}
}

// Init returns the initial command (none in this case).
func (m *Model) Init() tea.Cmd {
	return nil
}
