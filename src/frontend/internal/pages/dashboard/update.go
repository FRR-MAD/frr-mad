package dashboard

import (
	// "math/rand/v2"

	"time"

	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "r":
			m.ospfAnomalies = append(m.ospfAnomalies, "Reloading...")
		case "up":
			if m.showAnomalyOverlay {
				m.viewport.LineUp(10)
			} else {
				m.viewportLeft.LineUp(10)
				m.viewportRight.LineUp(10)
			}
		case "down":
			if m.showAnomalyOverlay {
				m.viewport.LineDown(10)
			} else {
				m.viewportLeft.LineDown(10)
				m.viewportRight.LineDown(10)
			}

			// FetchOSPFData returns a cmd and eventually triggers case msg.OSPFMsg
			return m, common.FetchOSPFData(m.logger)
		case "a":
			m.showAnomalyOverlay = !m.showAnomalyOverlay
		case "esc":
			if m.showAnomalyOverlay {
				m.showAnomalyOverlay = false
				return m, nil
			}
		}

	case common.OSPFMsg:
		m.ospfAnomalies = common.DetectOSPFAnomalies(string(msg))
	case common.ReloadMessage:
		m.currentTime = time.Time(msg)
		return m, reloadView()
	}

	return m, nil
}
