package rib

import (
	// "math/rand/v2"

	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	tea "github.com/charmbracelet/bubbletea"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.viewport.LineUp(10)
		case "down":
			m.viewport.LineDown(10)

			// FetchOSPFData returns a cmd and eventually triggers case msg.OSPFMsg
			return m, common.FetchOSPFData()
		}
	}

	return m, nil
}
