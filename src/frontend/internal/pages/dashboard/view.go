package dashboard

import (
	"fmt"
	backend "github.com/ba2025-ysmprc/frr-tui/internal/services"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	frrProto "github.com/ba2025-ysmprc/frr-tui/pkg"
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

	advertisingRouteTitle1 := styles.OSPFMonitoringTableTitleStyle.
		Width(boxWidthThreeFourth - 2).
		Render("Area 0.0.0.0, Router LSAs (Type 1)")

	advertisingRouteTitle2 := styles.OSPFMonitoringTableTitleStyle.
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

	systemResourcesBackend, _ := getSystemResources()

	systemResources := lipgloss.JoinVertical(lipgloss.Left,
		styles.BoxTitleStyle.Render("System Resources"),
		styles.GeneralBoxStyle.Width(boxWidthOneFourth-2).Render("here\nsome\nresources: \n"),
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
		styles.BoxTitleStyle.Render("Original Backend Call"),
		systemResourcesBackend,
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

func getSystemResources() (string, error) {
	params := map[string]*frrProto.ResponseValue{
		"client_id": &frrProto.ResponseValue{
			Kind: &frrProto.ResponseValue_StringValue{
				StringValue: "example_client",
			},
		},
	}

	response, err := backend.SendMessage("system", "allResources", params)
	if err != nil {
		return "err occurred in getSystemResources()", fmt.Errorf("rpc error: %w", err)
	}
	if response.Status != "success" {
		return "response.Status was no success", fmt.Errorf("backend returned status %q: %s", response.Status, response.Message)
	}

	stringValue := response.Data.GetStringValue()

	return stringValue, nil
}
