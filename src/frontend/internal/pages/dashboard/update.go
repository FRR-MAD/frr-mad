package dashboard

import (
	// "math/rand/v2"

	"fmt"
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
			} else if m.showExportOverlay {
				if m.cursor > 0 {
					m.cursor--
				}
			} else {
				m.viewportLeft.LineUp(10)
				m.viewportRight.LineUp(10)
			}
		case "down":
			if m.showAnomalyOverlay {
				m.viewport.LineDown(10)
			} else if m.showExportOverlay {
				if m.cursor < len(m.exportData)-1 {
					m.cursor++
				}
			} else {
				m.viewportLeft.LineDown(10)
				m.viewportRight.LineDown(10)
			}

			// FetchOSPFData returns a cmd and eventually triggers case msg.OSPFMsg
			return m, common.FetchOSPFData(m.logger)
		case "enter":
			if m.showExportOverlay {
				if len(m.exportOptions) == 0 {
					// m.status = "No export options available"
					break
				}

				opt := m.exportOptions[m.cursor]

				data, ok := m.exportData[opt.MapKey]
				if !ok {
					// TODO: add a status to the model to dispay status messages to the user
					// m.status = fmt.Sprintf("No data for %q", opt.MapKey)
					break
				}

				err := common.WriteExportToFile(data, opt.Filename, m.exportDirectory)
				if err != nil {
					m.logger.Error(fmt.Sprintf("Failed to Export %s", opt.Filename))
					// m.status = fmt.Sprintf("Export error: %v", err)
				} else {
					m.logger.Info(fmt.Sprintf("Exported %s", opt.Filename))
					//m.status = fmt.Sprintf("Exported %q to /tmp/frr-mad/exports/%s", opt.Label, opt.Filename)
				}
			} else {
				return m, nil
			}
		case "a":
			if !m.showExportOverlay {
				m.showAnomalyOverlay = !m.showAnomalyOverlay
			}
		case "e":
			if !m.showAnomalyOverlay {
				m.showExportOverlay = !m.showExportOverlay
			}
		case "esc":
			if m.showAnomalyOverlay {
				m.showAnomalyOverlay = false
				return m, nil
			} else if m.showExportOverlay {
				m.showExportOverlay = false
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
