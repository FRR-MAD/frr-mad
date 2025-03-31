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
	currentView AppState
	dashboard   dashboard.Model
	ospf        ospfMonitoring.Model
	tabs        []string
	width       int
	height      int
}

func initModel() *AppModel {
	return &AppModel{
		currentView: ViewDashboard,
		dashboard:   dashboard.New(),
		ospf:        ospfMonitoring.New(),
		width:       80,
		height:      24,
	}
}

func (m *AppModel) Init() tea.Cmd {
	m.Titles()
	return nil
}

func (m *AppModel) Titles() {
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
		// Example: switch view based on key presses
		switch msg.String() {
		case "1":
			m.currentView = ViewDashboard
		case "2":
			m.currentView = ViewOSPF
		case "right":
			m.currentView = (m.currentView + 1) % totalViews
		case "left":
			m.currentView = (m.currentView + totalViews - 1) % totalViews
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}

	case tea.WindowSizeMsg:
		m.width = msg.Width - 5
		m.height = msg.Height - 5
	}

	// Delegate Update to active module
	var cmd tea.Cmd
	switch m.currentView {
	case ViewDashboard:
		updatedModel, cmd := m.dashboard.Update(msg)
		m.dashboard = updatedModel.(dashboard.Model)
		return m, cmd
	case ViewOSPF:
		updatedModel, cmd := m.ospf.Update(msg)
		m.ospf = updatedModel.(ospfMonitoring.Model)
		return m, cmd
	default:
		panic("unhandled default case")
	}
	return m, cmd
}

func (m *AppModel) View() string {
	// var tabs []string
	//tabs = []string{"aaaa aaa", "bbb bbbb"}
	// tabRowRendered = lipgloss.JoinHorizontal(lipgloss.Top, tabRow...)
	// tabRow := components.GetTabRow(tabs, activeTab, hasSubMenu)
	tabRow := components.CreateTabRow(m.tabs, 0, false)

	var content string
	switch m.currentView {
	case ViewDashboard:
		content = m.dashboard.View()
	case ViewOSPF:
		content = m.ospf.View()
	default:
		return "Unknown view"
	}

	//---------------------------------
	// Table Layout Approach
	//---------------------------------

	//var (
	//	purple = lipgloss.Color("99")
	//	//gray      = lipgloss.Color("245")
	//	//lightGray = lipgloss.Color("241")
	//
	//	headerStyle = lipgloss.NewStyle().Foreground(purple).Bold(true).Align(lipgloss.Center)
	//	//cellStyle    = lipgloss.NewStyle().Padding(0, 1).Width(m.width - 5)
	//	//oddRowStyle  = cellStyle.Foreground(gray)
	//	//evenRowStyle = cellStyle.Foreground(lightGray)
	//)
	//
	////rows := [][]string{
	////	{content},
	////}
	//
	//headers := []string{"Dashboard", "OSPF Monitor", "BGP Monitor", "Custom Command"}
	//
	//// Create a new table with borders, custom style function, headers, and rows.
	//headerTable := table.New().
	//	Border(lipgloss.NormalBorder()).
	//	BorderStyle(lipgloss.NewStyle().Foreground(purple)).
	//	StyleFunc(func(row, col int) lipgloss.Style {
	//		// Only the header row is needed.
	//		if row == table.HeaderRow {
	//			return headerStyle
	//		}
	//		return lipgloss.NewStyle()
	//	}).
	//	Headers(headers...).
	//	Rows() // No content rows.
	//
	//headerStr := headerTable.String()
	//
	//// Create a style for the merged content row.
	//mergedStyle := lipgloss.NewStyle().Width(m.width).Align(lipgloss.Center)
	//mergedContent := mergedStyle.Render(content)
	//
	//return lipgloss.JoinVertical(lipgloss.Left, headerStr, mergedContent)

	return lipgloss.JoinVertical(lipgloss.Left, lipgloss.NewStyle().Width(m.width).Render(tabRow), styles.ContentStyle.Width(m.width).Render(content))
}

func main() {
	p := tea.NewProgram(initModel())
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
