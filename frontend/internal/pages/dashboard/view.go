package dashboard

import (
	backend "github.com/ba2025-ysmprc/frr-tui/internal/services"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var currentSubTabLocal = -1

// DashboardView is the updated View function. This allows to call View with an argument.
func (m *Model) DashboardView(currentSubTab int) string {
	currentSubTabLocal = currentSubTab
	return m.View()
}

func (m *Model) View() string {
	if currentSubTabLocal == 0 {
		return m.renderOSPFDashboard()
	} else if currentSubTabLocal == 1 {
		return ""
	}
	return m.renderOSPFDashboard()
}

func (m *Model) renderOSPFDashboard() string {
	// Calculate box width dynamically for two horizontal boxes based on terminal width
	boxWidthForTwo := (m.windowSize.Width - 12) / 2 // - 6 (padding+border contentbox) - 5 (border + 1 gap)
	if boxWidthForTwo < 20 {
		boxWidthForTwo = 20 // Minimum width to ensure readability
	}

	boxWidthForOne := m.windowSize.Width - 8 // - 6 (padding+margin content) - 2 (for each border)
	if boxWidthForOne < 20 {
		boxWidthForOne = 20 // Minimum width to ensure readability
	}

	allGoodRows := backend.GetOSPFMetrics()
	anomalyRows := backend.GetOSPFAnomalies()

	ospfTable := table.New().
		Border(lipgloss.RoundedBorder()).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(styles.NormalBeige))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return styles.HeaderStyle
			case row == 0:
				return styles.FirstNormalRowCellStyle
			default:
				return styles.NormalCellStyle
			}
		}).
		Width(boxWidthForOne).
		Headers("Advertising Route", "LSA Type", "Status").
		Rows(allGoodRows...)

	ospfBadTable := table.New().
		BorderRow(true).
		BorderStyle(lipgloss.NewStyle().Foreground(lipgloss.Color(styles.BadRed))).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return styles.HeaderStyle
			default:
				return styles.BadCellStyle
			}
		}).
		Width(boxWidthForOne).
		Headers("Advertised Route", "Anomaly Type", "Details", "Troubleshot").
		Rows(anomalyRows...)

	// in future either show ospfTable (=no anomaly) or ospfBadTable when anomaly is detected
	verticalTables := lipgloss.JoinVertical(lipgloss.Left,
		styles.BoxTitleStyle.Render("All OSPF Routes are advertised as Expected"),
		ospfTable.Render(),
		styles.BoxTitleStyle.Render("OSPF Anomaly Detected"),
		ospfBadTable.Render(),
	)

	return verticalTables
}
