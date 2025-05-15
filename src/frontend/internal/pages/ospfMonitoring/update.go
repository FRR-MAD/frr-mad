package ospfMonitoring

import (
	// "math/rand/v2"

	"fmt"
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/toast"
	tea "github.com/charmbracelet/bubbletea"
	"time"
)

// Update handles incoming messages and updates the dashboard state.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var toastCmd tea.Cmd
	m.toast, toastCmd = m.toast.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.showExportOverlay {
				m.viewportRightHalf.LineUp(10)
			} else {
				m.viewport.LineUp(10)
			}
			return m, nil
		case "down":
			if m.showExportOverlay {
				m.viewportRightHalf.LineDown(10)
			} else {
				m.viewport.LineDown(10)
			}
			return m, nil
		case "tab":
			if m.showExportOverlay && len(m.exportOptions) > 0 {
				m.cursor = (m.cursor + 1) % len(m.exportOptions)
				m.viewportRightHalf.GotoTop()
			}
		case "shift+tab":
			if m.showExportOverlay && len(m.exportOptions) > 0 {
				m.cursor = (m.cursor - 1 + len(m.exportOptions)) % len(m.exportOptions)
				m.viewportRightHalf.GotoTop()
			}
		case "r":
			m.runningConfig = []string{"Reloading..."}
			return m, common.FetchRunningConfig(m.logger)
		case "enter":
			if m.showExportOverlay {
				if len(m.exportOptions) == 0 {
					// m.status = "No export options available"
					break
				}

				opt := m.exportOptions[m.cursor]

				data, ok := m.exportData[opt.MapKey]
				if !ok {
					return m, tea.Batch(
						toastCmd,
						toast.Show("No Data available", 10*time.Second),
					)
				}

				if err := common.WriteExportToFile(data, opt.Filename, m.exportDirectory); err != nil {
					return m, tea.Batch(
						toastCmd,
						toast.Show(fmt.Sprintf("Export failed: %v", err), 10*time.Second),
					)
				}

				return m, tea.Batch(
					toastCmd,
					toast.Show(fmt.Sprintf("Exported %s\nto %s/%s",
						opt.Label, m.exportDirectory, opt.Filename),
						10*time.Second),
				)
			} else {
				return m, common.FetchRunningConfig(m.logger)

			}
		case "e":
			if m.showExportOverlay {
				m.toast = toast.New()
			} else {
				err := m.fetchLatestData()
				if err != nil {
					m.logger.Error("Error while fetching all backend data for OSPF Monitor")
					return m, tea.Batch(
						toastCmd,
						toast.Show(fmt.Sprintf("Fetching data failed:\n%v", err), 10*time.Second),
					)
				}
			}
			m.showExportOverlay = !m.showExportOverlay
		case "esc":
			if m.showExportOverlay {
				m.toast = toast.New()
				m.showExportOverlay = false
				return m, nil
			}
		}
	case common.RunningConfigMsg:
		m.runningConfig = common.ShowRunningConfig(string(msg))
	}
	return m, nil
}
