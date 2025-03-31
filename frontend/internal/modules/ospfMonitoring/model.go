package ospfMonitoring

import tea "github.com/charmbracelet/bubbletea"

// Model defines the state for the dashboard module.
type Model struct {
	Title   string
	Metrics []string // Example: a list of metrics to display
}

// New creates and returns a new dashboard Model.
func New() Model {
	return Model{
		Title:   "OSPF Monitoring",
		Metrics: []string{"Metric 1: 400", "Metric 2: 500", "Metric 3: 600"},
	}
}

func (m Model) GetTitle() string {
	return m.Title
}

// Init returns the initial command (none in this case).
func (m Model) Init() tea.Cmd {
	return nil
}
