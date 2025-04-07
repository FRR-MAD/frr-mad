package ospfMonitoring

import (
	// "math/rand/v2"

	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming messages and updates the dashboard state.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		// Pressing "r" refreshes metrics (for demonstration).
		case "r":
			// In a real application, you could fetch updated metrics here.

			// these two cases don’t work properly because when key left/right is clicked on tabs then this also triggered so it doesnt match the actual active subtab.
		case "right":
			m.CurrentSubTab = (m.CurrentSubTab + 1) % 3
		case "left":
			m.CurrentSubTab = (m.CurrentSubTab + 3 - 1) % 3
		}
	}
	return m, nil
}
