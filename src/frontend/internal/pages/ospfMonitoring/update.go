package ospfMonitoring

import (
	// "math/rand/v2"

	"fmt"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/frr-mad/frr-tui/internal/common"
	"github.com/frr-mad/frr-tui/internal/ui/toast"
	"sort"
)

// Update handles incoming messages and updates the OSPF Monitor state.
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
		case "home":
			if m.showExportOverlay {
				m.viewportRightHalf.GotoTop()
			} else {
				m.viewport.GotoTop()
			}
		case "end":
			if m.showExportOverlay {
				m.viewportRightHalf.GotoBottom()
			} else {
				m.viewport.GotoBottom()
			}
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
		case "ctrl+r":
			m.runningConfig = []string{"Reloading..."}
			return m, common.FetchRunningConfig(m.logger)
		case "enter":
			if m.showExportOverlay {
				if len(m.exportOptions) == 0 {
					return m, nil
				}

				opts := make([]common.ExportOption, len(m.exportOptions))
				copy(opts, m.exportOptions)
				sort.Slice(opts, func(i, j int) bool {
					return opts[i].Label < opts[j].Label
				})

				opt := opts[m.cursor]

				data, ok := m.exportData[opt.MapKey]
				if !ok {
					return m, tea.Batch(
						toastCmd,
						toast.Show("No Data available"),
					)
				}

				if err := common.WriteExportToFile(data, opt.Filename, m.exportDirectory); err != nil {
					return m, tea.Batch(
						toastCmd,
						toast.Show(fmt.Sprintf("Export failed: %v", err)),
					)
				}

				return m, tea.Batch(
					toastCmd,
					toast.Show(fmt.Sprintf("Exported to: %s/%s",
						m.exportDirectory, opt.Filename)),
				)
			} else {
				return m, common.FetchRunningConfig(m.logger)

			}
		case "ctrl+e":
			if m.showExportOverlay {
				m.toast = toast.New()
			} else {
				err := m.fetchLatestData()
				if err != nil {
					m.logger.Error("Error while fetching all backend data for OSPF Monitor")
					return m, tea.Batch(
						toastCmd,
						toast.Show(fmt.Sprintf("Fetching data failed:\n%v", err)),
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
