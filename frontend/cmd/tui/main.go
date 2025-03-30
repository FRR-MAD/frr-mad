package main

import (
	"fmt"
	// "log"
	"os"

	"github.com/ba2025-ysmprc/frr-tui/internal/modules/dashboard"
	"github.com/ba2025-ysmprc/frr-tui/internal/modules/ospfMonitoring"

	tea "github.com/charmbracelet/bubbletea"
)

type AppState int

const (
	ViewDashboard AppState = iota
	ViewOSPF
	// ViewLogs
	// add here new Views
	totalViews
)

type AppModel struct {
	currentView AppState
	dashboard   dashboard.Model
	ospf        ospfMonitoring.Model
	// bgp         bgp.Model
	// logs        logs.Model
}

func (m AppModel) Init() tea.Cmd {
	// Initialize default module, load configs, etc.
	return nil
}

func (m AppModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		// Example: switch view based on key presses
		switch msg.String() {
		case "1":
			m.currentView = ViewDashboard
		case "2":
			m.currentView = ViewOSPF
		// case "3":
		// 	m.currentView = ViewLogs
		case "right":
			m.currentView = (m.currentView + 1) % totalViews
		case "left":
			m.currentView = (m.currentView + totalViews - 1) % totalViews
		case "ctrl+c", "q", "esc":
			return m, tea.Quit
		}
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
		//	m.bgp, cmd = m.bgp.Update(msg)
		//case ViewLogs:
		//	m.logs, cmd = m.logs.Update(msg)
	default:
		panic("unhandled default case")
	}
	return m, cmd
}

func (m AppModel) View() string {
	// Delegate View based on active module
	switch m.currentView {
	case ViewDashboard:
		return m.dashboard.View()
	case ViewOSPF:
		return m.ospf.View()
	//case ViewLogs:
	//	return m.logs.View()
	default:
		return "Unknown view"
	}
}

func main() {
	p := tea.NewProgram(AppModel{
		currentView: ViewDashboard,
		dashboard:   dashboard.New(), // assume each module has its own New() function
		ospf:        ospfMonitoring.New(),
		//logs:        logs.New(),
	})
	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}
