package rib

import (
	"fmt"
	"github.com/frr-mad/frr-tui/internal/ui/toast"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/frr-mad/frr-tui/internal/common"
	backend "github.com/frr-mad/frr-tui/internal/services"
	"github.com/frr-mad/frr-tui/internal/ui/components"
	"github.com/frr-mad/frr-tui/internal/ui/styles"
)

var currentSubTabLocal = -1

func (m *Model) RibView(currentSubTab int, readOnlyMode bool, textFilter *common.Filter) string {
	currentSubTabLocal = currentSubTab
	m.readOnlyMode = readOnlyMode
	m.textFilter = textFilter
	return m.View()
}

func (m *Model) View() string {
	var body string
	var bodyFooter string
	var content string

	statusBar := true

	if m.showExportOverlay {
		content = components.RenderExportOptions(
			m.exportOptions,
			m.exportData,
			&m.cursor,
			&m.viewportRightHalf,
		)
	} else {
		switch currentSubTabLocal {
		case 0:
			body = m.renderRibTab()
		case 1:
			body = m.renderFibTab()
		case 2:
			body = m.renderRibWithProtocolFilterTab("ospf")
		case 3:
			body = m.renderRibWithProtocolFilterTab("bgp")
		case 4:
			body = m.renderRibWithProtocolFilterTab("connected")
		case 5:
			body = m.renderRibWithProtocolFilterTab("static")
		default:
			body = m.renderRibTab()
		}

		if statusBar {
			var filterBox string
			if m.textFilter.Active {
				filterBox = "Filter: " + m.textFilter.Input.View()
			} else {
				filterBox = "Filter: " + styles.FooterBoxStyle.Render("press [:] to activate filter")
			}
			filterBox = styles.FilterTextStyle().Render(filterBox)

			statusBox := lipgloss.NewStyle().Width(styles.WidthTwoH1Box).Margin(0, 2).Render(m.statusMessage)
			if m.statusMessage != "" {
				styles.SetStatusSeverity(m.statusSeverity)
				if len(m.statusMessage) > (styles.WidthTwoH1Box - styles.MarginX2) {
					m.statusMessage = m.statusMessage[:styles.WidthTwoH1Box-styles.MarginX2-3] + "..."
				}
				statusMessage := styles.StatusTextStyle().Render(m.statusMessage)
				statusBox = lipgloss.NewStyle().Width(styles.WidthTwoH1Box).Margin(0, 2).Render(statusMessage)
			}

			bodyFooter = lipgloss.JoinHorizontal(lipgloss.Top, statusBox, filterBox)

			content = lipgloss.JoinVertical(lipgloss.Left, body, bodyFooter)
		} else {
			content = body
		}
	}

	toastView := m.toast.View()
	if toastView == "" {
		return content
	}

	totalW := styles.WidthBasis
	totalH := styles.HeightBasis
	x := 0
	y := 0

	return toast.Overlay(content, toastView, x, y, totalW, totalH)
}

func (m *Model) renderRibTab() string {
	m.statusSeverity = styles.SeverityError
	m.statusMessage = "You opened RIB Tab And this is 100% a message that is way to long and should be cut that the system dont burn down."

	rib, err := backend.GetRIB(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetRIB")
	}
	ribFibSummary, err := backend.GetRibFibSummary(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetRibFibSummary")
	}

	amountOfRIBRoutes := strconv.Itoa(int(ribFibSummary.RoutesTotal))

	routes := make([]string, 0, len(rib.Routes))
	for route := range rib.Routes {
		routes = append(routes, route)
	}
	list := common.SortedPrefixList(routes)
	sort.Sort(&list)

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

	// apply filters if active
	ribTableData = common.FilterRows(ribTableData, m.textFilter.Query)

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
	m.viewport.Width = styles.WidthBasis
	// -3 (table Header) -2 (box border bottom style)
	m.viewport.Height = styles.HeightViewPortCompletePage - styles.HeightH1 - 3 - 2 - styles.FilterBoxHeight

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
	m.statusSeverity = styles.SeverityInfo
	m.statusMessage = "You opened FIB Tab"

	rib, err := backend.GetRIB(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetRIB")
	}
	ribFibSummary, err := backend.GetRibFibSummary(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetRibFibSummary")
	}

	amountOfFIBRoutes := strconv.Itoa(int(ribFibSummary.RoutesTotalFib))

	routes := make([]string, 0, len(rib.Routes))
	for route := range rib.Routes {
		routes = append(routes, route)
	}
	list := common.SortedPrefixList(routes)
	sort.Sort(&list)

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

	// apply filters if active
	fibTableData = common.FilterRows(fibTableData, m.textFilter.Query)

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
	m.viewport.Width = styles.WidthBasis
	// -3 (table Header) -2 (box border bottom style)
	m.viewport.Height = styles.HeightViewPortCompletePage - styles.HeightH1 - 3 - 2 - styles.FilterBoxHeight

	// Set only the body into the viewport
	m.viewport.SetContent(
		styles.H2OneContentBoxCenterStyle().Render(bodyContent),
	)

	boxBottomBorder := styles.H1OneSmallBoxBottomBorderStyle().Render("")

	completeFIBTab := lipgloss.JoinVertical(lipgloss.Left, headers, m.viewport.View(), boxBottomBorder)

	return completeFIBTab
}

func (m *Model) renderRibWithProtocolFilterTab(protocolName string) string {
	m.statusSeverity = styles.SeverityWarning
	m.statusMessage = "You opened a Partial RIB Tab"

	rib, err := backend.GetRIB(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetRIB")
	}
	ribFibSummary, err := backend.GetRibFibSummary(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetRibFibSummary")
	}

	protocolName = strings.ToLower(protocolName)

	amountOfRibRoutes := "0"

	for _, routeSummary := range ribFibSummary.RouteSummaries {
		if routeSummary == nil {
			continue
		}
		if strings.Contains(strings.ToLower(routeSummary.Type), protocolName) {
			amountOfRibRoutes = strconv.Itoa(int(routeSummary.Rib))
			break
		}
	}

	routes := make([]string, 0, len(rib.Routes))
	for route := range rib.Routes {
		routes = append(routes, route)
	}
	list := common.SortedPrefixList(routes)
	sort.Sort(&list)

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

	// apply filters if active
	partialRIBRoutesTableData = common.FilterRows(partialRIBRoutesTableData, m.textFilter.Query)

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
	m.viewport.Width = styles.WidthBasis
	// -3 (table Header) -2 (box border bottom style)
	m.viewport.Height = styles.HeightViewPortCompletePage - styles.HeightH1 - 3 - 2 - styles.FilterBoxHeight

	// Set only the body into the viewport
	m.viewport.SetContent(
		styles.H2OneContentBoxCenterStyle().Render(bodyContent),
	)

	boxBottomBorder := styles.H1OneSmallBoxBottomBorderStyle().Render("")

	// Render complete view
	completePartialRIBRoutesTab := lipgloss.JoinVertical(lipgloss.Left, headers, m.viewport.View(), boxBottomBorder)
	return completePartialRIBRoutesTab
}
