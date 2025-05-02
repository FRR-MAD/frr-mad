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
	m.viewport.Width = styles.WidthTwoH1ThreeFourth + 2
	m.viewport.Height = m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight - 2

	allGoodRows := backend.GetOSPFMetrics()
	anomalyRows := backend.GetOSPFAnomalies()

	advertisingRouteTitle1 := styles.H2TitleStyle().
		Width(styles.WidthTwoH2ThreeFourth).
		Render("Area 0.0.0.0, Router LSAs (Type 1)")

	advertisingRouteTitle2 := styles.H2TitleStyle().
		Width(styles.WidthTwoH2ThreeFourth).
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
		Width(styles.WidthTwoH2ThreeFourthBox).
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
		Width(styles.WidthTwoH2ThreeFourthBox).
		Headers("Advertised Route", "Anomaly Type", "Details", "Troubleshot").
		Rows(anomalyRows...)

	// in future either show ospfTable (=no anomaly) or ospfBadTable when anomaly is detected
	verticalTables := lipgloss.JoinVertical(lipgloss.Left,
		styles.H1TitleStyle().Width(styles.WidthTwoH1ThreeFourth).Render("All OSPF Routes are advertised as Expected"),
		advertisingRouteTitle1,
		ospfTable.Render(),
		advertisingRouteTitle2,
		ospfTable.Render(),
		styles.H1TitleStyle().Width(styles.WidthTwoH1ThreeFourth).Render("OSPF Anomaly Detected"),
		ospfBadTable.Render(),
	)

	// Update the viewport content with...
	m.viewport.SetContent(verticalTables)

	cpuAmount, cpuUsage, memoryUsage, err := getSystemResources()
	var cpuAmountString, cpuUsageString, memoryString string
	if err != nil {
		cpuAmountString = "N/A"
		cpuUsageString = "N/A"
		memoryString = "N/A"
	} else {
		cpuAmountString = fmt.Sprintf("%v", cpuAmount)
		cpuUsageString = fmt.Sprintf("%.2f%%", cpuUsage*100)
		memoryString = fmt.Sprintf("%.2f%%", memoryUsage)
	}

	cpuStatistics := lipgloss.JoinVertical(lipgloss.Left,
		styles.H2TitleStyle().Width(styles.WidthTwoH2OneFourth).Render("CPU Metrics"),
		styles.H2TwoContentBoxStyleP1101().Width(styles.WidthTwoH2OneFourthBox).Render(
			"CPU Usage: "+cpuUsageString+"\n"+
				"Cores: "+cpuAmountString),
		styles.H2BoxBottomBorderStyle().Width(styles.WidthTwoH2OneFourth).Render(""),
	)

	memoryStatistics := lipgloss.JoinVertical(lipgloss.Left,
		styles.H2TitleStyle().Width(styles.WidthTwoH2OneFourth).Render("Memory Metrics"),
		styles.H2TwoContentBoxStyleP1101().Width(styles.WidthTwoH2OneFourthBox).Render(
			"Memory Usage: "+memoryString),
		styles.H2BoxBottomBorderStyle().Width(styles.WidthTwoH2OneFourth).Render(""),
	)

	systemResources := lipgloss.JoinVertical(lipgloss.Left,
		styles.H1TitleStyle().Width(styles.WidthTwoH1OneFourth).Render("System Resources"),
		cpuStatistics,
		memoryStatistics,
	)

	horizontalDashboard := lipgloss.JoinHorizontal(lipgloss.Top,
		m.viewport.View(),
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
