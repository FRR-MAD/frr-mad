package main

import (
	"fmt"
	// "log"
	"os"

	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/components"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"

	"github.com/ba2025-ysmprc/frr-tui/internal/modules/dashboard"
	"github.com/ba2025-ysmprc/frr-tui/internal/modules/ospfMonitoring"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AppState int

const (
	ViewDashboard AppState = iota
	ViewOSPF
	// add here new Views
	totalViews
)

type AppModel struct {
	currentView  AppState
	tabs         []string
	subTabs      []string
	activeSubTab bool
	dashboard    *dashboard.Model
	ospf         *ospfMonitoring.Model
	windowSize   *common.WindowSize
	tabRowHeight int
	footer       *components.Footer
	footerHeight int
}

func initModel() *AppModel {
	windowSize := &common.WindowSize{Width: 80, Height: 24}

	return &AppModel{
		currentView:  ViewDashboard,
		tabs:         []string{},
		subTabs:      []string{},
		activeSubTab: false,
		dashboard:    dashboard.New(windowSize),
		ospf:         ospfMonitoring.New(windowSize),
		windowSize:   windowSize,
		tabRowHeight: 5,
		footer:       components.NewFooter("press 'esc' to quit"),
		footerHeight: 1,
	}
}

func (m *AppModel) Init() tea.Cmd {
	m.setTitles()
	return tea.Batch(
		m.dashboard.Init(),
	)
}

func (m *AppModel) setTitles() {
	modules := []common.TitledModule{
		m.dashboard,
		m.ospf,
	}
	for _, mod := range modules {
		m.tabs = append(m.tabs, mod.GetTitle())
	}
}

func (m *AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "1":
			m.currentView = ViewDashboard
		case "2":
			m.currentView = ViewOSPF
		//case "r":
		//	if m.currentView == ViewDashboard {
		//		m.dashboard = dashboard.New(m.windowSize)
		//		return m, m.dashboard.Init()
		//	}
		case "right":
			m.currentView = (m.currentView + 1) % totalViews
		case "left":
			m.currentView = (m.currentView + totalViews - 1) % totalViews
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		case "f":
			m.footer.Append("press 'f' to append")
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
	default:
		panic("unhandled default case")
	}
	return m, cmd
}

func (m *AppModel) View() string {

	var content string
	switch m.currentView {
	case ViewDashboard:
		content = m.dashboard.View()
		m.footer.Clean()
		m.footer.Append("press 'r' to refresh dashboard")
		m.footer.Append("press 'e' to export everything")
	case ViewOSPF:
		content = m.ospf.View()
		m.footer.Clean()
		m.footer.Append("press 'r' to refresh OSPF monitoring")
		m.footer.Append("press 'e' to export OSPF data")
	default:
		return "Unknown view"
	}

	contentWidth := m.windowSize.Width - 4

	tabRow := components.CreateTabRow(m.tabs, int(m.currentView), m.activeSubTab, m.windowSize)
	footer := m.footer.Get()
	return lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.NewStyle().Width(contentWidth).Margin(0, 1).Render(tabRow),
		styles.ContentBoxStyle.Width(contentWidth).Height(m.windowSize.Height-m.tabRowHeight-m.footerHeight).Render(content),
		styles.FooterBoxStyle.Width(contentWidth).Render(footer),
	)
}

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
