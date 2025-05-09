package rib

import (
	"fmt"
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	backend "github.com/ba2025-ysmprc/frr-tui/internal/services"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/components"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
	"sort"
	"strconv"
	"strings"
)

var currentSubTabLocal = -1

func (m *Model) RibView(currentSubTab int) string {
	currentSubTabLocal = currentSubTab
	return m.View()
}

func (m *Model) View() string {
	if currentSubTabLocal == 0 {
		return m.renderRibTab()
	} else if currentSubTabLocal == 1 {
		return m.renderOSPFRoutesTab()
	} else if currentSubTabLocal == 2 {
		return m.renderConnectedRoutesTab()
	}
	return m.renderRibTab()
}

func (m *Model) renderRibTab() string {
	rib, err := backend.GetRIB()
	if err != nil {
		return common.PrintBackendError(err, "GetRIB")
	}

	// TODO: call backend for correct amount (backend needs to be adjusted)
	amountOfRoutes := 40

	routes := make([]string, 0, len(rib.Routes))
	for route := range rib.Routes {
		routes = append(routes, route)
	}
	sort.Sort(common.IpList(routes))

	var ribTableData [][]string

	for _, route := range routes {
		routeEntry := rib.Routes[route]

		for _, routeEntryData := range routeEntry.Routes {
			var nexthopsList []string
			for _, nexthop := range routeEntryData.Nexthops {
				nexthopsList = append(nexthopsList, nexthop.Ip+" "+nexthop.InterfaceName)
			}
			ribTableData = append(ribTableData, []string{
				routeEntryData.Prefix,
				routeEntryData.Protocol,
				strings.Join(nexthopsList, "\n"),
				strconv.FormatBool(routeEntryData.Installed),
			})
		}
	}

	// Order all Table Data
	common.SortTableByIPColumn(ribTableData)

	rowsRIB := len(ribTableData)
	ribTable := components.NewMultilineTable(
		[]string{
			"Prefix",
			"Protocol",
			"Next Hops",
			"Installed",
		},
		rowsRIB,
	)
	for _, r := range ribTableData {
		ribTable = ribTable.Row(r...)
	}

	ribHeader := styles.H1TitleStyleForOne().
		Render(fmt.Sprintf("Routing Information Base - Received Routes"))
	ribTableHeader := styles.H2TitleStyleForOne().
		Render("The RIB contains " + strconv.Itoa(amountOfRoutes) + " routes")

	// Extract table header and body (top border, header row, bottom border)
	tableStr := ribTable.String()
	lines := strings.Split(tableStr, "\n")
	var headerLines, bodyLines []string
	if len(lines) > 3 {
		headerLines = lines[:3]
		bodyLines = lines[3:]
	} else {
		headerLines = lines
		bodyLines = nil
	}
	// Render header and body
	tableHeaderContent := styles.H2OneContentBoxCenterStyle().Render(strings.Join(headerLines, "\n"))
	bodyContent := strings.Join(bodyLines, "\n")

	headers := lipgloss.JoinVertical(lipgloss.Left, ribHeader, ribTableHeader, tableHeaderContent)

	// Configure viewport
	contentMaxHeight := m.windowSize.Height -
		styles.TabRowHeight -
		styles.FooterHeight -
		styles.HeightH1 -
		styles.HeightH2 -
		3 - 2 // -3 (table Header) -2 (box border bottom style)
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = contentMaxHeight

	// Set only the body into the viewport
	m.viewport.SetContent(
		styles.H2OneContentBoxCenterStyle().Render(bodyContent),
	)

	boxBottomBorder := styles.H2OneBoxBottomBorderStyle().Render("")

	// Render complete view
	completeRIBTab := lipgloss.JoinVertical(lipgloss.Left, headers, m.viewport.View(), boxBottomBorder)
	return completeRIBTab

	//ribTableBox := lipgloss.JoinVertical(lipgloss.Left,
	//	styles.H2OneContentBoxCenterStyle().Render(ribTable.String()),
	//	styles.H2OneBoxBottomBorderStyle().Render(""),
	//)
	//
	//contentMaxHeight := m.windowSize.Height -
	//	styles.TabRowHeight -
	//	styles.FooterHeight -
	//	styles.HeightH1 -
	//	styles.HeightH2
	//m.viewport.Width = styles.WidthBasis
	//m.viewport.Height = contentMaxHeight
	//
	//m.viewport.SetContent(ribTableBox)
	//
	//completeRIBTab := lipgloss.JoinVertical(lipgloss.Left, headers, m.viewport.View())
	//
	//return completeRIBTab
}

func (m *Model) renderOSPFRoutesTab() string {
	return "OSPF learned routes"
}

func (m *Model) renderConnectedRoutesTab() string {

	return "Directly Connected Networks"
}
