package rib

import (
	// "math/rand/v2"

	"fmt"
	"sort"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/frr-mad/frr-mad/src/frontend/internal/common"
	"github.com/frr-mad/frr-mad/src/frontend/internal/ui/styles"
	"github.com/frr-mad/frr-mad/src/frontend/internal/ui/toast"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var toastCmd tea.Cmd
	m.toast, toastCmd = m.toast.Update(msg)

	if !m.statusTimer.IsZero() && time.Since(m.statusTimer) > m.statusDuration {
		m.statusMessage = ""
		m.statusTimer = time.Time{}
	}

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

				err := common.CopyOSC52(data)
				if err != nil {
					m.statusSeverity = styles.SeverityWarning
					m.statusMessage = "Could not Copy Clipboard: use a terminal with osc52 enabled"
				} else {
					m.statusSeverity = styles.SeverityInfo
					m.statusMessage = "successfully copied to clipboard"
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
				m.statusSeverity = styles.SeverityInfo
				m.statusMessage = ""
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
			m.statusSeverity = styles.SeverityInfo
			m.statusMessage = ""
			if m.showExportOverlay {
				m.toast = toast.New()
				m.showExportOverlay = false
				return m, nil
			}
		}

	case common.QuitTuiFailedMsg:
		m.statusSeverity = styles.SeverityError
		m.statusMessage = string(msg)
		return m, nil
	}

	return m, nil
}
