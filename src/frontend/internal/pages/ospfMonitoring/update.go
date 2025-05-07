package ospfMonitoring

import (
	// "math/rand/v2"

	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming messages and updates the dashboard state.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.viewport.LineUp(10)
			return m, nil
		case "down":
			m.viewport.LineDown(10)
			return m, nil
		case "r":
			m.runningConfig = []string{"Reloading..."}
			return m, common.FetchRunningConfig()
		case "enter":
			return m, common.FetchRunningConfig()
		}
	case common.RunningConfigMsg:
		m.runningConfig = common.ShowRunningConfig(string(msg))
	}
	return m, nil
}
