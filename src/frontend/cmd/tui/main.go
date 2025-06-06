package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/viewport"

	"github.com/frr-mad/frr-tui/internal/pages/rib"

	"github.com/frr-mad/frr-mad/src/logger"
	"github.com/frr-mad/frr-tui/internal/common"
	"github.com/frr-mad/frr-tui/internal/configs"
	"github.com/frr-mad/frr-tui/internal/pages/dashboard"
	"github.com/frr-mad/frr-tui/internal/pages/ospfMonitoring"
	"github.com/frr-mad/frr-tui/internal/pages/shell"
	"github.com/frr-mad/frr-tui/internal/ui/components"
	"github.com/frr-mad/frr-tui/internal/ui/styles"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	ViewDashboard common.AppState = iota
	ViewOSPFMonitoring
	ViewRIB
	ViewShell
	// add here new Views
	totalViews
)

var (
	DaemonVersion = "unknown"
	TUIVersion    = "unknown"
	GitCommit     = "unknown"
	BuildDate     = "unknown"
	RepoURL       = "https://github.com/frr-mad/frr-mad"
)

var subTabsLength int

type AppModel struct {
	startupConfig     string
	activeViews       []common.AppState
	currentView       common.AppState
	tabs              []common.Tab
	currentSubTab     int
	readOnlyMode      bool
	windowSize        *common.WindowSize
	viewport          viewport.Model
	dashboard         *dashboard.Model
	ospf              *ospfMonitoring.Model
	rib               *rib.Model
	shell             *shell.Model
	showSystemInfo    bool
	preventSubTabExit bool
	footer            *components.Footer
	footerOptions     []common.FooterOption
	textFilter        *common.Filter
	logger            *logger.Logger
}

func initModel(config *configs.Config) *AppModel {
	windowSize := &common.WindowSize{Width: 157, Height: 38}

	vp := viewport.New(styles.WidthViewPortCompletePage,
		styles.HeightViewPortCompletePage-styles.BodyFooterHeight)

	logLevel := logger.ConvertLogLevelFromConfig(config.Default.DebugLevel)

	appLogger, err := logger.NewApplicationLogger("frr-mad-tui",
		fmt.Sprintf("%v/frr_mad_tui_application.log", config.Default.LogPath))
	if err != nil {
		log.Fatalf("Failed to create application logger: %v", err)
	}
	appLogger.SetDebugLevel(logLevel)
	appLogger.Info("Starting Frontend Application")

	dashboardLogger := appLogger.WithComponent("dashboard")
	ospfLogger := appLogger.WithComponent("ospf")
	ribLogger := appLogger.WithComponent("rib")
	shellLogger := appLogger.WithComponent("shell")

	ti := textinput.New()
	ti.Placeholder = "type to filter..."
	ti.CharLimit = 32
	ti.Width = 20

	enabled := []common.AppState{
		ViewDashboard,
		ViewOSPFMonitoring,
		ViewRIB,
		ViewShell,
	}

	return &AppModel{
		startupConfig:     "",
		activeViews:       enabled,
		currentView:       ViewDashboard,
		tabs:              []common.Tab{},
		currentSubTab:     -1,
		readOnlyMode:      true,
		windowSize:        windowSize,
		viewport:          vp,
		dashboard:         dashboard.New(windowSize, dashboardLogger, config.Default.ExportPath),
		ospf:              ospfMonitoring.New(windowSize, ospfLogger, config.Default.ExportPath),
		rib:               rib.New(windowSize, ribLogger, config.Default.ExportPath),
		shell:             shell.New(windowSize, shellLogger),
		showSystemInfo:    false,
		preventSubTabExit: false,
		footer:            components.NewFooter("[ctrl+c] exit FRR-MAD", "[i] info", "[enter] enter sub tabs"),
		logger:            appLogger,
		textFilter:        &common.Filter{Active: false, Query: "", Input: ti},
	}
}

func (m *AppModel) Init() tea.Cmd {
	m.setTitles()
	startupConfig, err := m.setStartupConfig()
	if err != nil {
		return tea.Batch(
			tea.ClearScreen,
			tea.Quit,
		)
	} else {
		m.startupConfig = startupConfig
	}

	common.SetAppVersionInfo(DaemonVersion, TUIVersion, GitCommit, BuildDate, RepoURL)

	return tea.Batch(
		m.dashboard.Init(),
	)
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			if !m.textFilter.Active && !m.preventSubTabExit {
				if m.currentSubTab == -1 {
					idx := m.indexOfAppState(m.currentView)
					next := (idx + 1) % len(m.activeViews)
					m.currentView = m.activeViews[next]
					m.currentSubTab = -1
					if m.currentView == ViewShell {
						m.footer.CleanInfo()
						m.footer.Append("[enter] enter sub tabs")
					} else {
						m.footer.CleanInfo()
						m.footer.SetMainMenuOptions()
					}
				} else {
					m.currentSubTab = (m.currentSubTab + 1) % subTabsLength
				}
			}
		case "left":
			if !m.textFilter.Active && !m.preventSubTabExit {
				if m.currentSubTab == -1 {
					idx := m.indexOfAppState(m.currentView)
					next := (idx + len(m.activeViews) - 1) % len(m.activeViews)
					m.currentView = m.activeViews[next]
					m.currentSubTab = -1
					if m.currentView == ViewShell {
						m.footer.CleanInfo()
						m.footer.Append("[enter] enter sub tabs")
					} else {
						m.footer.CleanInfo()
						m.footer.SetMainMenuOptions()
					}
				} else {
					m.currentSubTab = (m.currentSubTab + subTabsLength - 1) % subTabsLength
				}
			}
		case "up":
			if m.showSystemInfo {
				m.viewport.LineUp(10)
			} else {
				return m.delegateToActiveView(msg)
			}
		case "down":
			if m.showSystemInfo {
				m.viewport.LineDown(10)
			} else {
				return m.delegateToActiveView(msg)
			}
		case "home":
			if m.showSystemInfo {
				m.viewport.GotoTop()
			} else {
				return m.delegateToActiveView(msg)
			}
		case "end":
			if m.showSystemInfo {
				m.viewport.GotoBottom()
			} else {
				return m.delegateToActiveView(msg)
			}
		case "enter":
			if m.currentSubTab == -1 {
				m.currentSubTab = 0
				m.footer.Clean()
				m.footer.Append("[esc] exit sub tab")
				currentPageOptions := m.getCurrentFooterOptions()
				m.footer.AppendMultiple(currentPageOptions)
			}
		case ":":
			if !m.textFilter.Active {
				m.textFilter.Active = true
				m.textFilter.Input.Focus()
			} else {
				m.textFilter.Active = false
				m.textFilter.Query = ""
				m.textFilter.Input.SetValue("")
				m.textFilter.Input.Blur()
			}
			return m, nil
		case "ctrl+w":
			m.readOnlyMode = !m.readOnlyMode
			styles.ChangeReadWriteMode(m.readOnlyMode)
		case "ctrl+e", "ctrl+a":
			m.preventSubTabExit = !m.preventSubTabExit
		case "i":
			if m.currentView != ViewShell && !m.textFilter.Active {
				m.showSystemInfo = !m.showSystemInfo
			}
		case "esc":
			if m.textFilter.Active {
				m.textFilter.Active = false
				m.textFilter.Query = ""
				m.textFilter.Input.Blur()
				return m, nil
			} else if m.showSystemInfo {
				m.showSystemInfo = false
			} else if m.preventSubTabExit {
				m.preventSubTabExit = false
			} else {
				m.currentSubTab = -1
				m.footer.SetMainMenuOptions()
			}
		case "ctrl+c":
			currentConfig, err := common.GetRunningConfig(m.logger)
			if err != nil {
				m.logger.Error(fmt.Sprintf("Error fetching OSPF Running-Config: %v", err))
			}
			if m.startupConfig == currentConfig {
				return m, tea.Batch(
					tea.ClearScreen,
					tea.Quit,
				)
			} else {
				return m, common.QuitTuiFailedCmd(
					"Config changed: running FRR config must match the config at TUI startup.",
				)
			}
		}
	//case tea.MouseEvent:
	//	return m, nil
	//case tea.MouseMsg:
	//	return m, nil
	//case tea.MouseAction:
	//	return m, nil
	//case tea.MouseButton:
	//	return m, nil
	case tea.WindowSizeMsg:
		m.windowSize.Width = msg.Width
		m.windowSize.Height = msg.Height
		styles.SetWindowSizes(common.WindowSize{
			Width:  msg.Width,
			Height: msg.Height,
		})
	}

	if m.textFilter.Active {
		m.textFilter.Input, _ = m.textFilter.Input.Update(msg)
		m.textFilter.Query = m.textFilter.Input.Value()
		return m, nil
	}

	// Delegate Update to active module
	var cmd tea.Cmd
	switch m.currentView {
	case ViewDashboard:
		updatedModel, cmd := m.dashboard.Update(msg)
		m.dashboard = updatedModel.(*dashboard.Model)
		return m, cmd
	case ViewOSPFMonitoring:
		updatedModel, cmd := m.ospf.Update(msg)
		m.ospf = updatedModel.(*ospfMonitoring.Model)
		return m, cmd
	case ViewRIB:
		updatedModel, cmd := m.rib.Update(msg)
		m.rib = updatedModel.(*rib.Model)
		return m, cmd
	case ViewShell:
		updatedModel, cmd := m.shell.Update(msg)
		m.shell = updatedModel.(*shell.Model)
		return m, cmd
	default:
		panic("unhandled default case")
	}
	return m, cmd
}

func (m *AppModel) View() string {

	// -2 (for content border) -2 (is necessary for error free usage --> leads to style errors without it)
	contentWidth := m.windowSize.Width - 4
	contentHeight := m.windowSize.Height - styles.TabRowHeight - styles.BorderContentBox - styles.FooterHeight

	tabRow := components.CreateTabRow(m.tabs, int(m.currentView), m.currentSubTab, m.windowSize, m.logger)
	footer := m.footer.Get()

	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = styles.HeightViewPortCompletePage

	var content string
	if m.showSystemInfo {
		systemInfo := components.GetSystemInfoOverlay()
		m.viewport.SetContent(systemInfo)
	} else {
		switch m.currentView {
		case ViewDashboard:
			content = m.dashboard.DashboardView(m.currentSubTab, m.readOnlyMode, m.textFilter)
			subTabsLength = m.dashboard.GetSubTabsLength()
		case ViewOSPFMonitoring:
			content = m.ospf.OSPFView(m.currentSubTab, m.readOnlyMode, m.textFilter)
			subTabsLength = m.ospf.GetSubTabsLength()
		case ViewRIB:
			content = m.rib.RibView(m.currentSubTab, m.readOnlyMode, m.textFilter)
			subTabsLength = m.rib.GetSubTabsLength()
		case ViewShell:
			content = m.shell.ShellView(m.currentSubTab, m.readOnlyMode)
			subTabsLength = m.shell.GetSubTabsLength()
		default:
			return "Unknown view"
		}
	}

	var tuiContent string
	if m.showSystemInfo {
		tuiContent = m.viewport.View()
	} else {
		tuiContent = content
	}

	tui := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Width(contentWidth).Margin(0, 1).Render(tabRow),
		styles.ContentBoxStyle().Width(contentWidth).Height(contentHeight).Render(tuiContent),
		styles.FooterBoxStyle.Width(contentWidth).Render(footer),
	)

	return tui
}

func (m *AppModel) delegateToActiveView(msg tea.Msg) (*AppModel, tea.Cmd) {
	var cmd tea.Cmd
	switch m.currentView {
	case ViewDashboard:
		updatedModel, cmd := m.dashboard.Update(msg)
		m.dashboard = updatedModel.(*dashboard.Model)
		return m, cmd
	case ViewOSPFMonitoring:
		updatedModel, cmd := m.ospf.Update(msg)
		m.ospf = updatedModel.(*ospfMonitoring.Model)
		return m, cmd
	case ViewRIB:
		updatedModel, cmd := m.rib.Update(msg)
		m.rib = updatedModel.(*rib.Model)
		return m, cmd
	case ViewShell:
		updatedModel, cmd := m.shell.Update(msg)
		m.shell = updatedModel.(*shell.Model)
		return m, cmd
	default:
		panic("unhandled default case")
	}
	return m, cmd
}

func main() {
	startupConfig, err := common.GetRunningConfigWithoutLog()
	if err != nil || startupConfig == "" {
		fmt.Fprintf(os.Stderr, "Error loading FRR config: %v\n", err)
		os.Exit(1)
	} else {
		startFrrMadTui()
	}
}

func startFrrMadTui() {
	// Load configuration
	config, err := configs.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	maybeUpdateTERM()
	p := tea.NewProgram(initModel(config), tea.WithAltScreen())
	// TODO: find a way to fix the TUI that you cant scroll away (in apple terminal)
	// TODO: the problem with mouseMotion is, you cannot highlight text anymore with the mouse
	// p := tea.NewProgram(initModel(), tea.WithMouseCellMotion()) // start program with msg.MouseMsg options
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
