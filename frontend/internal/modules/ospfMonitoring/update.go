package ospfMonitoring

import (
	// "math/rand/v2"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming messages and updates the dashboard state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Pressing "r" refreshes metrics (for demonstration).
		case "r":
			// In a real application, you could fetch updated metrics here.
			m.Metrics = append(m.Metrics, "New Metric: 400")
		}
	}
	return m, nil
}
