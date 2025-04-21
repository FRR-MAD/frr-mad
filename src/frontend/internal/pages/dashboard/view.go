package dashboard

import (
	"fmt"
	backend "github.com/ba2025-ysmprc/frr-tui/internal/services"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/lipgloss/table"
)

var (
	currentSubTabLocal = -1
)

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
	boxWidthForTwo := (m.windowSize.Width - 10) / 2 // - 6 (padding+border contentbox) - 5 (border + 1 gap)
	if boxWidthForTwo < 20 {
		boxWidthForTwo = 20 // Minimum width to ensure readability
	}

	boxWidthForOne := m.windowSize.Width - 8 // - 6 (padding+margin content) - 2 (for each border)
	if boxWidthForOne < 20 {
		boxWidthForOne = 20 // Minimum width to ensure readability
	}

	boxWidthThreeFourth := boxWidthForTwo / 2 * 3
	boxWidthOneFourth := boxWidthForTwo / 2

	m.viewport.Width = boxWidthThreeFourth
	m.viewport.Height = m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight - 2

	gap := 2

	allGoodRows := backend.GetOSPFMetrics()
	anomalyRows := backend.GetOSPFAnomalies()

	advertisingRouteTitle1 := styles.ActiveContentStyle.
		Width(boxWidthThreeFourth - 2).
		Render("Area 0.0.0.0, Router LSAs (Type 1)")

	advertisingRouteTitle2 := styles.ActiveContentStyle.
		Width(boxWidthThreeFourth - 2).
		Render("Area 0.0.0.0, Autonomous System External LSAs (Type 5)")

	ospfTable := table.New().
		Border(lipgloss.HiddenBorder()).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == table.HeaderRow:
				return styles.HeaderStyle
			//case row == 0:
			//	return styles.FirstNormalRowCellStyle
			default:
				return styles.NormalCellStyle
			}
		}).
		Width(boxWidthThreeFourth).
		//Headers("Advertising Route", "LSA Type", "Status").
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
		Width(boxWidthThreeFourth).
		Headers("Advertised Route", "Anomaly Type", "Details", "Troubleshot").
		Rows(anomalyRows...)

	cpuAmount, cpuUsage, memoryUsage, err := getSystemResources()
	var cpuAmountString, cpuUsageString, memoryString string
	if err != nil {
		cpuAmountString = "N/A"
		cpuUsageString = "N/A"
		memoryString = "N/A"
	} else {
		cpuAmountString = fmt.Sprintf("%v", cpuAmount)
		cpuUsageString = fmt.Sprintf("%.2f%%", cpuUsage*100)
		memoryString = fmt.Sprintf("%.2f%%", memoryUsage*100)
	}

	// Convert 0.0–1.0 → percentage and clamp to [0,100]
	//cpuPct := clamp(cpuUsage*100, 0, 100)
	// memPct := clamp(memoryUsage*100, 0, 100)

	systemResources := lipgloss.JoinVertical(lipgloss.Left,
		styles.BoxTitleStyle.Render("System Resources"),
		styles.GeneralBoxStyle.Width(boxWidthOneFourth-2).
			Render("CPU Usage:\n"+cpuUsageString+
				"\nCores:\n"+cpuAmountString+
				"\n\nMemory Usage:\n"+memoryString),
	)

	// in future either show ospfTable (=no anomaly) or ospfBadTable when anomaly is detected
	verticalTables := lipgloss.JoinVertical(lipgloss.Left,
		styles.BoxTitleStyle.Render("All OSPF Routes are advertised as Expected"),
		advertisingRouteTitle1,
		ospfTable.Render(),
		advertisingRouteTitle2,
		ospfTable.Render(),
		styles.BoxTitleStyle.Render("OSPF Anomaly Detected"),
		ospfBadTable.Render(),
	)

	// Update the viewport content with...
	m.viewport.SetContent(verticalTables)

	horizontalDashboard := lipgloss.JoinHorizontal(lipgloss.Top,
		m.viewport.View(),
		lipgloss.NewStyle().Width(gap).Render(""),
		systemResources,
	)

	return horizontalDashboard
}

func getSystemResources() (int64, float64, float64, error) {

	response, err := backend.SendMessage("system", "allResources", nil)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("rpc error: %w", err)
	}
	if response.Status != "success" {
		return 0, 0, 0, fmt.Errorf("backend returned status %q: %s", response.Status, response.Message)
	}

	systemMetrics := response.Data.GetSystemMetrics()

	cores := systemMetrics.CpuAmount
	cpuUsage := systemMetrics.CpuUsage
	memoryUsage := systemMetrics.MemoryUsage

	return cores, cpuUsage, memoryUsage, nil
}

// clamp returns x clamped to the [min,max] interval.
func clamp(x, min, max float64) float64 {
	switch {
	case x < min:
		return min
	case x > max:
		return max
	default:
		return x
	}
}
