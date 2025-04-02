package dashboard

import (
	// "math/rand/v2"

	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

// Update handles incoming messages and updates the dashboard state.
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			// return m, m.Init()
			// m.ospfAnomalies = append(m.ospfAnomalies, "Reload Placeholder")
		}

	case common.OSPFMsg:
		m.ospfAnomalies = common.DetectOSPFAnomalies(string(msg))
		// Log OSPF anomalies to history
		// history.AddEntry(fmt.Sprintf("OSPF Anomalies Detected:\n%s", strings.Join(m.ospfAnomalies, "\n")))
	}

	return m, nil
}
