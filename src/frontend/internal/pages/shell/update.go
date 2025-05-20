package shell

import (
	// "math/rand/v2"

	"fmt"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/frr-mad/frr-tui/internal/common"
	backend "github.com/frr-mad/frr-tui/internal/services"
	"google.golang.org/protobuf/encoding/protojson"
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
		case "home":
			m.viewport.GotoTop()
		case "end":
			m.viewport.GotoBottom()
		case "backspace":
			if currentSubTabLocal == 0 && len(m.bashInput) > 0 {
				m.bashInput = m.bashInput[:len(m.bashInput)-1]
			} else if currentSubTabLocal == 1 && len(m.vtyshInput) > 0 {
				m.vtyshInput = m.vtyshInput[:len(m.vtyshInput)-1]
			} else if currentSubTabLocal == 2 && m.activeBackendInput == "service" && len(m.backendServiceInput) > 0 {
				m.backendServiceInput = m.backendServiceInput[:len(m.backendServiceInput)-1]
			} else if currentSubTabLocal == 2 && m.activeBackendInput == "command" && len(m.backendCommandInput) > 0 {
				m.backendCommandInput = m.backendCommandInput[:len(m.backendCommandInput)-1]
			}
		case "tab":
			if currentSubTabLocal == 2 {
				m.activeBackendInput = "command"
			}
		case "left", "right":
			m.ClearInput()
			m.ClearOutput()
			m.clearBackendInput()
		case "enter":
			if currentSubTabLocal == 0 {
				bashOutput, err := common.RunCustomCommand("bash", m.bashInput, 5*time.Second, m.logger)
				if err != nil {
					m.bashOutput = fmt.Sprintf("Error: %v", err)
				} else {
					m.bashOutput = bashOutput
				}
			} else if currentSubTabLocal == 1 {
				vtyshOutput, err := common.RunCustomCommand("vtysh", m.vtyshInput, 5*time.Second, m.logger)
				if err != nil {
					m.vtyshOutput = fmt.Sprintf("Error: %v", err)
				} else {
					m.vtyshOutput = vtyshOutput
				}
			} else if currentSubTabLocal == 2 {
				res, err := backend.SendMessage(
					m.backendServiceInput,
					m.backendCommandInput,
					nil,
					m.logger,
				)
				if err != nil {
					m.backendResponse = fmt.Sprintf("Error: %v", err)
				} else {
					// Pretty‑print the protobuf into nice indented JSON
					marshaler := protojson.MarshalOptions{
						Multiline:     true,
						Indent:        "  ",
						UseProtoNames: true,
					}
					pretty, perr := marshaler.Marshal(res.Data)
					if perr != nil {
						m.backendResponse = res.Data.String()
					} else {
						m.backendResponse = string(pretty)
					}
				}
				m.clearBackendInput()
				m.activeBackendInput = "service"
			}
		case "ctrl+w":
			m.ClearInput()

		default:
			if currentSubTabLocal == 0 {
				m.bashInput += msg.String()
			} else if currentSubTabLocal == 1 {
				m.vtyshInput += msg.String()
			} else if m.activeBackendInput == "service" {
				m.backendServiceInput += msg.String()
			} else if m.activeBackendInput == "command" {
				m.backendCommandInput += msg.String()
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
