package shell

import (
	// "math/rand/v2"

	"fmt"
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

// Update handles incoming messages and updates the dashboard state.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			m.viewport.LineUp(10)
		case "down":
			m.viewport.LineDown(10)
		case "backspace":
			if currentSubTabLocal == 0 && len(m.bashInput) > 0 {
				m.bashInput = m.bashInput[:len(m.bashInput)-1]
			} else if currentSubTabLocal == 1 && len(m.vtyshInput) > 0 {
				m.vtyshInput = m.vtyshInput[:len(m.vtyshInput)-1]
			}
		case "left", "right":
			m.ClearInput()
			m.ClearOutput()
		case "enter":
			if currentSubTabLocal == 0 {
				bashOutput, err := common.RunCustomCommand("bash", m.bashInput, 5*time.Second)
				if err != nil {
					m.bashOutput = fmt.Sprintf("Error: %v", err)
				} else {
					m.bashOutput = bashOutput
				}
			} else if currentSubTabLocal == 1 {
				vtyshOutput, err := common.RunCustomCommand("vtysh", m.vtyshInput, 5*time.Second)
				if err != nil {
					m.vtyshOutput = fmt.Sprintf("Error: %v", err)
				} else {
					m.vtyshOutput = vtyshOutput
				}
			}
		default:
			if currentSubTabLocal == 0 {
				m.bashInput += msg.String()
			} else if currentSubTabLocal == 1 {
				m.vtyshInput += msg.String()
			}
		}
	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			m.viewport.LineUp(1)
		case tea.MouseButtonWheelDown:
			m.viewport.LineDown(1)
		default:
			panic("unhandled default case")
		}
	}

	return m, nil
}
