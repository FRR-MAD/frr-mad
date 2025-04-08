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
		case "0":
			m.viewport.LineUp(5) // Scroll up by two lines.
		case "9":
			m.viewport.LineDown(5) // Scroll down by one line.
		case "backspace":
			if currentSubTabLocal == 0 && len(m.BashInput) > 0 {
				m.BashInput = m.BashInput[:len(m.BashInput)-1]
			} else if currentSubTabLocal == 1 && len(m.VtyshInput) > 0 {
				m.VtyshInput = m.VtyshInput[:len(m.VtyshInput)-1]
			}
		case "left", "right", "up":
			m.ClearInput()
			m.ClearOutput()
		case "enter":
			if currentSubTabLocal == 0 {
				bashOutput, err := common.RunCommand("bash", m.BashInput, 5*time.Second)
				if err != nil {
					m.BashOutput = fmt.Sprintf("Error: %v", err)
				} else {
					m.BashOutput = bashOutput
				}
			} else if currentSubTabLocal == 1 {
				vtyshOutput, err := common.RunCommand("vtysh", m.VtyshInput, 5*time.Second)
				if err != nil {
					m.VtyshOutput = fmt.Sprintf("Error: %v", err)
				} else {
					m.VtyshOutput = vtyshOutput
				}
			}
		default:
			if currentSubTabLocal == 0 {
				m.BashInput += msg.String()
			} else if currentSubTabLocal == 1 {
				m.VtyshInput += msg.String()
			}
		}
	case tea.MouseMsg:
		switch msg.Button {
		case tea.MouseButtonWheelUp:
			m.viewport.LineUp(1)
		case tea.MouseButtonWheelDown:
			m.viewport.LineDown(1)

		}
	}

	return m, nil
}
