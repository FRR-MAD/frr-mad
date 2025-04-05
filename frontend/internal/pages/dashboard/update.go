package dashboard

import (
	// "math/rand/v2"

	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.ospfAnomalies = append(m.ospfAnomalies, "Reloading...")

			// FetchOSPFData returns a cmd and eventually triggers case msg.OSPFMsg
			return m, common.FetchOSPFData()
		}

	case common.OSPFMsg:
		m.ospfAnomalies = common.DetectOSPFAnomalies(string(msg))
	}

	return m, nil
}
