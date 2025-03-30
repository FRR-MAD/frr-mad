package main

import (
	"log"
	"os"

	"github.com/ba2025-ysmprc/frr-tui/frontend/modules/dashboard"

	tea "github.com/charmbracelet/bubbletea"
)

type AppState int

const (
	ViewDashboard AppState = iota
	ViewBGP
	ViewLogs
)

type AppModel struct {
	currentView AppState
	dashboard   dashboard.Model
	bgp         bgp.Model
	logs        logs.Model
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
			m.currentView = ViewBGP
		case "3":
			m.currentView = ViewLogs
		}
	}

	// Delegate Update to active module
	var cmd tea.Cmd
	switch m.currentView {
	case ViewDashboard:
		m.dashboard, cmd = m.dashboard.Update(msg)
	case ViewBGP:
		m.bgp, cmd = m.bgp.Update(msg)
	case ViewLogs:
		m.logs, cmd = m.logs.Update(msg)
	}

	return m, cmd
}

func (m AppModel) View() string {
	// Delegate View based on active module
	switch m.currentView {
	case ViewDashboard:
		return m.dashboard.View()
	case ViewBGP:
		return m.bgp.View()
	case ViewLogs:
		return m.logs.View()
	default:
		return "Unknown view"
	}
}

func main() {
	p := tea.NewProgram(AppModel{
		currentView: ViewDashboard,
		dashboard:   dashboard.New(), // assume each module has its own New() function
		bgp:         bgp.New(),
		logs:        logs.New(),
	})
	if err := p.Start(); err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}
