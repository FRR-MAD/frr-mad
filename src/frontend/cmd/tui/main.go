package main

import (
	"fmt"

	// "log"
	"os"

	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/ba2025-ysmprc/frr-tui/internal/pages/dashboard"
	"github.com/ba2025-ysmprc/frr-tui/internal/pages/ospfMonitoring"
	"github.com/ba2025-ysmprc/frr-tui/internal/pages/shell"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/components"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AppState int

const (
	ViewDashboard AppState = iota
	ViewOSPF
	// code for Presentation slides
	//ViewOSPF2
	//ViewOSPF3
	ViewShell
	// add here new Views
	totalViews
)

var subTabsLength int

type AppModel struct {
	currentView   AppState
	tabs          []common.Tab
	currentSubTab int
	windowSize    *common.WindowSize
	dashboard     *dashboard.Model
	ospf          *ospfMonitoring.Model
	// code for Presentation slides
	//ospf2         *ospfMonitoring.Model
	//ospf3         *ospfMonitoring.Model
	shell         *shell.Model
	footer        *components.Footer
	footerOptions []common.FooterOption
}

func initModel() *AppModel {
	windowSize := &common.WindowSize{Width: 80, Height: 24}

	return &AppModel{
		currentView:   ViewDashboard,
		tabs:          []common.Tab{},
		currentSubTab: -1,
		windowSize:    windowSize,
		dashboard:     dashboard.New(windowSize),
		ospf:          ospfMonitoring.New(windowSize),
		// code for Presentation slides
		//ospf2:         ospfMonitoring.New(windowSize),
		//ospf3:         ospfMonitoring.New(windowSize),
		shell:  shell.New(windowSize),
		footer: components.NewFooter("[ctrl+c] exit FRR-MAD", "[enter] enter sub tabs"),
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
		case "1":
			if m.currentSubTab == -1 {
				m.currentView = ViewDashboard
			} else {
				m.currentSubTab = 0
			}
		case "2":
			if m.currentSubTab == -1 {
				m.currentView = ViewOSPF
			} else if subTabsLength >= 2 {
				m.currentSubTab = 1
			}
		case "3":
			if m.currentSubTab == -1 {
				m.currentView = ViewShell
			} else if subTabsLength >= 3 {
				m.currentSubTab = 2
			}
		case "9":
			if m.currentSubTab == -1 {
				break
			} else {
				m.currentSubTab = subTabsLength - 1
			}
		// code for Presentation slides
		//case "4":
		//	if m.currentSubTab == -1 {
		//		m.currentView = ViewOSPF2
		//	}
		//case "5":
		//	if m.currentSubTab == -1 {
		//		m.currentView = ViewOSPF3
		//	}
		case "right":
			if m.currentSubTab == -1 {
				m.currentView = (m.currentView + 1) % totalViews
				m.currentSubTab = -1
			} else {
				m.currentSubTab = (m.currentSubTab + 1) % subTabsLength
			}
		case "left":
			if m.currentSubTab == -1 {
				m.currentView = (m.currentView + totalViews - 1) % totalViews
				m.currentSubTab = -1
			} else {
				m.currentSubTab = (m.currentSubTab + subTabsLength - 1) % subTabsLength
			}
		case "enter":
			if m.currentSubTab == -1 {
				m.currentSubTab = 0
				m.footer.Clean()
				m.footer.Append("[esc] exit sub tab")
				currentPageOptions := m.getCurrentFooterOptions()
				m.footer.AppendMultiple(currentPageOptions)
			}
		case "esc":
			m.currentSubTab = -1
			m.footer.SetMainMenuOptions()
		case "ctrl+c":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.windowSize.Width = msg.Width
		m.windowSize.Height = msg.Height
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
	// code for Presentation slides
	//case ViewOSPF2:
	//	updatedModel, cmd := m.ospf.Update(msg)
	//	m.ospf = updatedModel.(*ospfMonitoring.Model)
	//	return m, cmd
	//case ViewOSPF3:
	//	updatedModel, cmd := m.ospf.Update(msg)
	//	m.ospf = updatedModel.(*ospfMonitoring.Model)
	//	return m, cmd
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
		content = m.dashboard.DashboardView(m.currentSubTab)
		subTabsLength = m.dashboard.GetSubTabsLength()
	case ViewOSPF:
		content = m.ospf.OSPFView(m.currentSubTab)
		subTabsLength = m.ospf.GetSubTabsLength()
	// code for Presentation slides
	//case ViewOSPF2:
	//	content = m.ospf.OSPFView(m.currentSubTab)
	//	subTabsLength = m.ospf.GetSubTabsLength()
	//case ViewOSPF3:
	//	content = m.ospf.OSPFView(m.currentSubTab)
	//	subTabsLength = m.ospf.GetSubTabsLength()
	case ViewShell:
		content = m.shell.ShellView(m.currentSubTab)
		subTabsLength = m.shell.GetSubTabsLength()
	default:
		return "Unknown view"
	}

	contentWidth := m.windowSize.Width - 4
	contentHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight

	tabRow := components.CreateTabRow(m.tabs, int(m.currentView), m.currentSubTab, m.windowSize)
	footer := m.footer.Get()

	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Width(contentWidth).Margin(0, 1).Render(tabRow),
		styles.ContentBoxStyle.Width(contentWidth).Height(contentHeight).Render(content),
		styles.FooterBoxStyle.Width(contentWidth).Render(footer),
	)
}

func main() {
	maybeUpdateTERM()
	p := tea.NewProgram(initModel())
	// p := tea.NewProgram(initModel(), tea.WithMouseCellMotion()) // start program with msg.MouseMsg options
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
