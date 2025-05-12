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
		return m.renderFibTab()
	} else if currentSubTabLocal == 2 {
		return m.renderRibWithProtocolFilterTab("ospf")
	} else if currentSubTabLocal == 3 {
		return m.renderRibWithProtocolFilterTab("bgp")
	} else if currentSubTabLocal == 4 {
		return m.renderRibWithProtocolFilterTab("connected")
	} else if currentSubTabLocal == 5 {
		return m.renderRibWithProtocolFilterTab("static")
	}
	return m.renderRibTab()
}

func (m *Model) renderRibTab() string {
	rib, err := backend.GetRIB()
	if err != nil {
		return common.PrintBackendError(err, "GetRIB")
	}
	ribFibSummary, err := backend.GetRibFibSummary()
	if err != nil {
		return common.PrintBackendError(err, "GetRibFibSummary")
	}

	amountOfRIBRoutes := strconv.Itoa(int(ribFibSummary.RoutesTotal))

	routes := make([]string, 0, len(rib.Routes))
	for route := range rib.Routes {
		routes = append(routes, route)
	}
	sort.Sort(common.SortedPrefixList(routes))

	// return strings.Join(routes, "\n")

	var ribTableData [][]string

	for _, route := range routes {
		routeEntry := rib.Routes[route]

		for _, routeEntryData := range routeEntry.Routes {
			var nexthopsList []string
			var fibList []string
			for _, nexthop := range routeEntryData.Nexthops {
				if nexthop == nil {
					continue
				}

				var entry string
				if nexthop.Ip == "" {
					entry = nexthop.InterfaceName
				} else {
					entry = nexthop.Ip + " " + nexthop.InterfaceName
				}
				nexthopsList = append(nexthopsList, entry)
				fibList = append(fibList, strconv.FormatBool(nexthop.Fib))
			}
			ribTableData = append(ribTableData, []string{
				routeEntryData.Prefix,
				routeEntryData.Protocol,
				strings.Join(nexthopsList, "\n"),
				strings.Join(fibList, "\n"),
				strconv.FormatBool(routeEntryData.Installed),
				strconv.Itoa(int(routeEntryData.Distance)),
				strconv.Itoa(int(routeEntryData.Metric)),
				routeEntryData.Uptime,
			})
		}
	}

	rowsRIB := len(ribTableData)
	ribTable := components.NewRibMonitorTable(rowsRIB)
	for _, r := range ribTableData {
		ribTable = ribTable.Row(r...)
	}

	ribHeader := styles.H1TitleStyleForOne().
		Render(fmt.Sprintf("Routing Information Base - " + amountOfRIBRoutes + " Received Routes"))

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

	headers := lipgloss.JoinVertical(lipgloss.Left, ribHeader, tableHeaderContent)

	// Configure viewport
	contentMaxHeight := m.windowSize.Height -
		styles.TabRowHeight -
		styles.FooterHeight -
		styles.HeightH1 -
		3 - 2 // -3 (table Header) -2 (box border bottom style)
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = contentMaxHeight

	// Set only the body into the viewport
	m.viewport.SetContent(
		styles.H2OneContentBoxCenterStyle().Render(bodyContent),
	)

	boxBottomBorder := styles.H1OneSmallBoxBottomBorderStyle().Render("")

	// Render complete view
	completeRIBTab := lipgloss.JoinVertical(lipgloss.Left, headers, m.viewport.View(), boxBottomBorder)
	return completeRIBTab
}

func (m *Model) renderFibTab() string {
	rib, err := backend.GetRIB()
	if err != nil {
		return common.PrintBackendError(err, "GetRIB")
	}
	ribFibSummary, err := backend.GetRibFibSummary()
	if err != nil {
		return common.PrintBackendError(err, "GetRibFibSummary")
	}

	amountOfFIBRoutes := strconv.Itoa(int(ribFibSummary.RoutesTotalFib))

	routes := make([]string, 0, len(rib.Routes))
	for route := range rib.Routes {
		routes = append(routes, route)
	}
	sort.Sort(common.SortedPrefixList(routes))

	// return strings.Join(routes, "\n")

	var fibTableData [][]string

	for _, route := range routes {
		routeEntry := rib.Routes[route]

		for _, routeEntryData := range routeEntry.Routes {
			var nexthopsList []string
			var fibList []string
			for _, nexthop := range routeEntryData.Nexthops {
				if nexthop == nil {
					continue
				}

				// To confirm that all listed nexthops are indeed in the kernel FIB, check each "fib": true status
				if nexthop.Fib {
					var entry string
					if nexthop.Ip == "" {
						entry = nexthop.InterfaceName
					} else {
						entry = nexthop.Ip + " " + nexthop.InterfaceName
					}
					nexthopsList = append(nexthopsList, entry)
					fibList = append(fibList, strconv.FormatBool(nexthop.Fib))
				}
			}
			// "installed": true = FRR has pushed a forwarding entry for that prefix to the kernel (at least one)
			if routeEntryData.Installed {
				fibTableData = append(fibTableData, []string{
					routeEntryData.Prefix,
					routeEntryData.Protocol,
					strings.Join(nexthopsList, "\n"),
					strings.Join(fibList, "\n"),
					strconv.FormatBool(routeEntryData.Installed),
					strconv.Itoa(int(routeEntryData.Distance)),
					strconv.Itoa(int(routeEntryData.Metric)),
					routeEntryData.Uptime,
				})
			}
		}
	}

	rowsFIB := len(fibTableData)
	fibTable := components.NewRibMonitorTable(rowsFIB)
	for _, r := range fibTableData {
		fibTable = fibTable.Row(r...)
	}

	fibHeader := styles.H1TitleStyleForOne().
		Render(fmt.Sprintf("Forwarding Information Base - " + amountOfFIBRoutes + " Installed Routes"))

	// Extract table header and body (top border, header row, bottom border)
	tableStr := fibTable.String()
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

	headers := lipgloss.JoinVertical(lipgloss.Left, fibHeader, tableHeaderContent)

	// Configure viewport
	contentMaxHeight := m.windowSize.Height -
		styles.TabRowHeight -
		styles.FooterHeight -
		styles.HeightH1 -
		3 - 2 // -3 (table Header) -2 (box border bottom style)
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = contentMaxHeight

	// Set only the body into the viewport
	m.viewport.SetContent(
		styles.H2OneContentBoxCenterStyle().Render(bodyContent),
	)

	boxBottomBorder := styles.H1OneSmallBoxBottomBorderStyle().Render("")

	// Render complete view
	completeFIBTab := lipgloss.JoinVertical(lipgloss.Left, headers, m.viewport.View(), boxBottomBorder)
	return completeFIBTab
}

func (m *Model) renderRibWithProtocolFilterTab(protocolName string) string {
	rib, err := backend.GetRIB()
	if err != nil {
		return common.PrintBackendError(err, "GetRIB")
	}
	ribFibSummary, err := backend.GetRibFibSummary()
	if err != nil {
		return common.PrintBackendError(err, "GetRibFibSummary")
	}

	protocolName = strings.ToLower(protocolName)

	amountOfRibRoutes := "0"

	for _, routeSummary := range ribFibSummary.RouteSummaries {
		if routeSummary == nil {
			continue
		}
		if strings.ToLower(routeSummary.Type) == protocolName {
			amountOfRibRoutes = strconv.Itoa(int(routeSummary.Rib))
			break
		}
	}

	routes := make([]string, 0, len(rib.Routes))
	for route := range rib.Routes {
		routes = append(routes, route)
	}
	sort.Sort(common.SortedPrefixList(routes))

	// return strings.Join(routes, "\n")

	var partialRIBRoutesTableData [][]string

	for _, route := range routes {
		routeEntry := rib.Routes[route]

		for _, routeEntryData := range routeEntry.Routes {
			var nexthopsList []string
			var fibList []string
			for _, nexthop := range routeEntryData.Nexthops {
				if nexthop == nil {
					continue
				}

				var entry string
				if nexthop.Ip == "" {
					entry = nexthop.InterfaceName
				} else {
					entry = nexthop.Ip + " " + nexthop.InterfaceName
				}
				nexthopsList = append(nexthopsList, entry)
				fibList = append(fibList, strconv.FormatBool(nexthop.Fib))
			}
			if strings.ToLower(routeEntryData.Protocol) == protocolName {
				partialRIBRoutesTableData = append(partialRIBRoutesTableData, []string{
					routeEntryData.Prefix,
					routeEntryData.Protocol,
					strings.Join(nexthopsList, "\n"),
					strings.Join(fibList, "\n"),
					strconv.FormatBool(routeEntryData.Installed),
					strconv.Itoa(int(routeEntryData.Distance)),
					strconv.Itoa(int(routeEntryData.Metric)),
					routeEntryData.Uptime,
				})
			}
		}
	}

	rowsPartialRIBRoutesRIB := len(partialRIBRoutesTableData)
	partialRIBRoutesTable := components.NewRibMonitorTable(rowsPartialRIBRoutesRIB)
	for _, r := range partialRIBRoutesTableData {
		partialRIBRoutesTable = partialRIBRoutesTable.Row(r...)
	}

	partialRoutesHeader := styles.H1TitleStyleForOne().
		Render(fmt.Sprintf("Routing Information Base received " + amountOfRibRoutes + " routes via " + strings.ToUpper(protocolName)))

	// Extract table header and body (top border, header row, bottom border)
	tableStr := partialRIBRoutesTable.String()
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

	headers := lipgloss.JoinVertical(lipgloss.Left, partialRoutesHeader, tableHeaderContent)

	// Configure viewport
	contentMaxHeight := m.windowSize.Height -
		styles.TabRowHeight -
		styles.FooterHeight -
		styles.HeightH1 -
		3 - 2 // -3 (table Header) -2 (box border bottom style)
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = contentMaxHeight

	// Set only the body into the viewport
	m.viewport.SetContent(
		styles.H2OneContentBoxCenterStyle().Render(bodyContent),
	)

	boxBottomBorder := styles.H1OneSmallBoxBottomBorderStyle().Render("")

	// Render complete view
	completePartialRIBRoutesTab := lipgloss.JoinVertical(lipgloss.Left, headers, m.viewport.View(), boxBottomBorder)
	return completePartialRIBRoutesTab
}
