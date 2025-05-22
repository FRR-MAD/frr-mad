package main

import (
	"fmt"
	"log"
	"os"

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

type AppState int

const (
	ViewDashboard AppState = iota
	ViewOSPF
	ViewRIB
	ViewShell
	// add here new Views
	totalViews
)

var subTabsLength int

type AppModel struct {
	currentView   AppState
	tabs          []common.Tab
	currentSubTab int
	readOnlyMode  bool
	windowSize    *common.WindowSize
	dashboard     *dashboard.Model
	ospf          *ospfMonitoring.Model
	rib           *rib.Model
	shell         *shell.Model
	footer        *components.Footer
	footerOptions []common.FooterOption
	textFilter    *common.Filter
	logger        *logger.Logger
}

func initModel(config *configs.Config) *AppModel {
	windowSize := &common.WindowSize{Width: 80, Height: 24}

	debugLevel := getDebugLevel(config.Default.DebugLevel)
	appLogger := createLogger("frr_mad_frontend", fmt.Sprintf("%v/frr_mad_frontend.log", config.Default.LogPath))
	appLogger.SetDebugLevel(debugLevel)
	appLogger.Info("Starting Frontend Application")

	dashboardLogger := createLogger("dashboard_frontend", fmt.Sprintf("%v/dashboard_frontend.log", config.Default.LogPath))
	dashboardLogger.SetDebugLevel(debugLevel)

	ospfLogger := createLogger("ospf_frontend", fmt.Sprintf("%v/ospf_frontend.log", config.Default.LogPath))
	ospfLogger.SetDebugLevel(debugLevel)

	ribLogger := createLogger("rib_frontend", fmt.Sprintf("%v/rib_frontend.log", config.Default.LogPath))
	ribLogger.SetDebugLevel(debugLevel)

	shellLogger := createLogger("shell_frontend", fmt.Sprintf("%v/shell_frontend.log", config.Default.LogPath))
	shellLogger.SetDebugLevel(debugLevel)

	ti := textinput.New()
	ti.Placeholder = "type to filter..."
	ti.CharLimit = 32
	ti.Width = 20

	return &AppModel{
		currentView:   ViewDashboard,
		tabs:          []common.Tab{},
		currentSubTab: -1,
		readOnlyMode:  true,
		windowSize:    windowSize,
		dashboard:     dashboard.New(windowSize, dashboardLogger),
		ospf:          ospfMonitoring.New(windowSize, ospfLogger),
		rib:           rib.New(windowSize, ribLogger),
		shell:         shell.New(windowSize, shellLogger),
		footer:        components.NewFooter("[ctrl+c] exit FRR-MAD", "[enter] enter sub tabs"),
		logger:        appLogger,
		textFilter:    &common.Filter{Active: false, Query: "", Input: ti},
	}
}

func (m *AppModel) Init() tea.Cmd {
	m.setTitles()
	return tea.Batch(
		m.dashboard.Init(),
	)
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "right":
			if !m.textFilter.Active {
				if m.currentSubTab == -1 {
					m.currentView = (m.currentView + 1) % totalViews
					m.currentSubTab = -1
				} else {
					m.currentSubTab = (m.currentSubTab + 1) % subTabsLength
				}
			}
		case "left":
			if !m.textFilter.Active {
				if m.currentSubTab == -1 {
					m.currentView = (m.currentView + totalViews - 1) % totalViews
					m.currentSubTab = -1
				} else {
					m.currentSubTab = (m.currentSubTab + subTabsLength - 1) % subTabsLength
				}
			}
		case "up", "down", "home", "end":
			return m.delegateToActiveView(msg)
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
		case "esc":
			if m.textFilter.Active {
				m.textFilter.Active = false
				m.textFilter.Query = ""
				m.textFilter.Input.Blur()
				return m, nil
			} else {
				m.currentSubTab = -1
				m.footer.SetMainMenuOptions()
			}
		case "ctrl+c":
			return m, tea.Batch(
				tea.ClearScreen,
				tea.Quit,
			)
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
	case ViewOSPF:
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

	var content string
	switch m.currentView {
	case ViewDashboard:
		content = m.dashboard.DashboardView(m.currentSubTab, m.readOnlyMode, m.textFilter)
		subTabsLength = m.dashboard.GetSubTabsLength()
	case ViewOSPF:
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

	// -2 (for content border) -2 (is necessary for error free usage --> leads to style errors without it)
	contentWidth := m.windowSize.Width - 4
	contentHeight := m.windowSize.Height - styles.TabRowHeight - styles.BorderContentBox - styles.FooterHeight

	tabRow := components.CreateTabRow(m.tabs, int(m.currentView), m.currentSubTab, m.windowSize, m.logger)
	footer := m.footer.Get()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Width(contentWidth).Margin(0, 1).Render(tabRow),
		styles.ContentBoxStyle().Width(contentWidth).Height(contentHeight).Render(content),
		styles.FooterBoxStyle.Width(contentWidth).Render(footer),
	)
}

func (m *AppModel) delegateToActiveView(msg tea.Msg) (*AppModel, tea.Cmd) {
	var cmd tea.Cmd
	switch m.currentView {
	case ViewDashboard:
		updatedModel, cmd := m.dashboard.Update(msg)
		m.dashboard = updatedModel.(*dashboard.Model)
		return m, cmd
	case ViewOSPF:
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

// Create a new logger instance
func createLogger(name, filePath string) *logger.Logger {
	tuiLogger, err := logger.NewLogger(name, filePath)
	if err != nil {
		log.Fatalf("Failed to create logger %s: %v", name, err)
	}
	return tuiLogger
}

// Convert debug level string to int
func getDebugLevel(level string) int {
	switch level {
	case "debug":
		return 2
	case "error":
		return 1
	default:
		return 0
	}
}
