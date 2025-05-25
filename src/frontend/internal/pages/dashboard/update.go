package dashboard

import (
	// "math/rand/v2"

	"fmt"
	"github.com/frr-mad/frr-tui/internal/ui/styles"
	"github.com/frr-mad/frr-tui/internal/ui/toast"
	"sort"
	"strconv"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/frr-mad/frr-tui/internal/common"
)

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var toastCmd tea.Cmd
	m.toast, toastCmd = m.toast.Update(msg)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up":
			if m.showAnomalyOverlay {
				m.viewport.LineUp(10)
			} else if m.showExportOverlay {
				m.viewportRightHalf.LineUp(10)
			} else {
				m.viewportLeft.LineUp(10)
				m.viewportRight.LineUp(10)
			}
		case "down":
			if m.showAnomalyOverlay {
				m.viewport.LineDown(10)
			} else if m.showExportOverlay {
				m.viewportRightHalf.LineDown(10)
			} else {
				m.viewportLeft.LineDown(10)
				m.viewportRight.LineDown(10)
			}

			// FetchOSPFData returns a cmd and eventually triggers case msg.OSPFMsg
			return m, common.FetchOSPFData(m.logger)
		case "home":
			if m.showExportOverlay {
				m.viewportRightHalf.GotoTop()
			} else if m.showAnomalyOverlay {
				m.viewport.GotoTop()
			} else {
				m.viewportRight.GotoTop()
				m.viewportLeft.GotoTop()
			}
		case "end":
			if m.showExportOverlay {
				m.viewportRightHalf.GotoBottom()
			} else if m.showAnomalyOverlay {
				m.viewport.GotoBottom()
			} else {
				m.viewportRight.GotoBottom()
				m.viewportLeft.GotoBottom()
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
						toast.Show(fmt.Sprintf(strconv.Itoa(int(m.cursor))+" Export failed: %v", err)),
					)
				}

				return m, tea.Batch(
					toastCmd,
					toast.Show(fmt.Sprintf("Exported to: %s/%s",
						m.exportDirectory, opt.Filename)),
				)
			}
			return m, nil
		case "ctrl+a":
			if !m.showExportOverlay {
				m.showAnomalyOverlay = !m.showAnomalyOverlay
			}
		case "ctrl+e":
			if !m.showAnomalyOverlay {
				if m.showExportOverlay {
					m.toast = toast.New()
					m.statusSeverity = styles.SeverityInfo
					m.statusMessage = ""
				} else {
					err := m.fetchLatestData()
					if err != nil {
						m.logger.Error("Error while fetching all backend data for dashboard")
						m.statusSeverity = styles.SeverityError
						m.statusMessage = "Please re-open export options: Error while fetching data"
						return m, tea.Batch(
							toastCmd,
							toast.Show(fmt.Sprintf("Fetching data failed:\n%v", err)),
						)
					}
				}
				m.showExportOverlay = !m.showExportOverlay
			}
		case "esc":
			m.statusSeverity = styles.SeverityInfo
			m.statusMessage = ""
			if m.showAnomalyOverlay {
				m.showAnomalyOverlay = false
				return m, nil
			} else if m.showExportOverlay {
				m.toast = toast.New()
				m.showExportOverlay = false
				return m, nil
			}
		}

	case common.ReloadMessage:
		m.currentTime = time.Time(msg)
		return m, reloadView()

	case common.QuitTuiFailedMsg:
		m.statusSeverity = styles.SeverityError
		m.statusMessage = string(msg)
		return m, nil
	}

	return m, nil
}
