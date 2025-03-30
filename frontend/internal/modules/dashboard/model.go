package dashboard

import tea "github.com/charmbracelet/bubbletea"

// Model defines the state for the dashboard module.
type Model struct {
	Title   string
	Metrics []string // Example: a list of metrics to display
}

// New creates and returns a new dashboard Model.
func New() Model {
	return Model{
		Title:   "dashboard Overview",
		Metrics: []string{"Metric 1: 100", "Metric 2: 200", "Metric 3: 300"},
	}
}

// Init returns the initial command (none in this case).
func (m Model) Init() tea.Cmd {
	return nil
}
