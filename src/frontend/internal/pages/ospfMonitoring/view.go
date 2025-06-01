package ospfMonitoring

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	"github.com/frr-mad/frr-tui/internal/ui/toast"

	"github.com/frr-mad/frr-tui/internal/common"
	backend "github.com/frr-mad/frr-tui/internal/services"
	"github.com/frr-mad/frr-tui/internal/ui/components"
	"github.com/frr-mad/frr-tui/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

var currentSubTabLocal = -1

func (m *Model) OSPFView(currentSubTab int, readOnlyMode bool, textFilter *common.Filter) string {
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
			m.statusMessage,
			m.statusSeverity,
		)
	} else {
		switch currentSubTabLocal {
		case 0:
			body = m.renderLsdbMonitorTab()
		case 1:
			body = m.renderRouterMonitorTab()
		case 2:
			body = m.renderNetworkMonitorTab()
		case 3:
			body = m.renderExternalMonitorTab()
		case 4:
			body = m.renderNeighborMonitorTab()
		case 5:
			body = m.renderRunningConfigTab()
			statusBar = false
		default:
			body = m.renderLsdbMonitorTab()
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
				var cutToSizeMessage string
				maxLength := styles.WidthTwoH1Box - styles.MarginX2 - 3
				if maxLength > 0 && len(m.statusMessage) > maxLength {
					cutToSizeMessage = m.statusMessage[:maxLength] + "..."
				} else if maxLength > 0 {
					cutToSizeMessage = m.statusMessage
				} else {
					cutToSizeMessage = "..."
				}
				statusMessage := styles.StatusTextStyle().Render(cutToSizeMessage)
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

func (m *Model) renderLsdbMonitorTab() string {
	var lsdbBlocks []string

	lsdb, err := backend.GetLSDB(m.logger)
	if err != nil {
		m.statusMessage = "Failed to fetch LSDB data"
		m.statusSeverity = styles.SeverityError
		return common.PrintBackendError(err, "GetLSDB")
	}

	// extract and sort the map keys
	lsdbAreas := make([]string, 0, len(lsdb.Areas))
	for area := range lsdb.Areas {
		lsdbAreas = append(lsdbAreas, area)
	}
	sort.Strings(lsdbAreas)

	// ===== OSPF Internal LSAs (Type 1-4) =====
	for _, areaID := range lsdbAreas {
		lsaTypes := lsdb.Areas[areaID]

		var routerLinkStateTableData [][]string
		var networkLinkStateTableData [][]string
		var summaryLinkStateTableData [][]string
		var asbrSummaryLinkStateTableData [][]string
		var nssaExternalLinkStateTableData [][]string

		var amountOfRouterLS string
		var amountOfNetworkLS string
		var amountOfSummaryLS string
		var amountOfAsSummaryLS string
		var amountOfnssaExternalLS string

		// loop through LSAs (type 1-4 + type 7) and extract data for tables
		for _, routerLinkState := range lsaTypes.RouterLinkStates {
			routerLinkStateTableData = append(routerLinkStateTableData, []string{
				routerLinkState.Base.AdvertisedRouter,
				strconv.Itoa(int(routerLinkState.NumOfRouterLinks)),
				strconv.Itoa(int(routerLinkState.Base.LsaAge)),
			})
		}
		if routerLinkStateTableData != nil {
			amountOfRouterLS = strconv.Itoa(int(lsaTypes.RouterLinkStatesCount))
		} else {
			amountOfRouterLS = "0"
		}
		for _, networkLinkState := range lsaTypes.NetworkLinkStates {
			networkLinkStateTableData = append(networkLinkStateTableData, []string{
				networkLinkState.Base.LsId,
				networkLinkState.Base.AdvertisedRouter,
				strconv.Itoa(int(networkLinkState.Base.LsaAge)),
			})
		}
		if networkLinkStateTableData != nil {
			amountOfNetworkLS = strconv.Itoa(int(lsaTypes.NetworkLinkStatesCount))
		} else {
			amountOfNetworkLS = "0"
		}

		for _, summarLinkState := range lsaTypes.SummaryLinkStates {
			summaryLinkStateTableData = append(summaryLinkStateTableData, []string{
				summarLinkState.SummaryAddress,
				summarLinkState.Base.AdvertisedRouter,
				strconv.Itoa(int(summarLinkState.Base.LsaAge)),
			})
		}
		if summaryLinkStateTableData == nil {
			amountOfSummaryLS = "0"
		} else {
			amountOfSummaryLS = strconv.Itoa(int(lsaTypes.SummaryLinkStatesCount))
		}

		for _, asbrSummaryLinkState := range lsaTypes.AsbrSummaryLinkStates {
			asbrSummaryLinkStateTableData = append(asbrSummaryLinkStateTableData, []string{
				asbrSummaryLinkState.Base.LsId,
				asbrSummaryLinkState.Base.AdvertisedRouter,
				strconv.Itoa(int(asbrSummaryLinkState.Base.LsaAge)),
			})
		}
		if asbrSummaryLinkStateTableData == nil {
			amountOfAsSummaryLS = "0"
		} else {
			amountOfAsSummaryLS = strconv.Itoa(int(lsaTypes.AsbrSummaryLinkStatesCount))
		}

		for _, nssaExternalLinkStates := range lsaTypes.NssaExternalLinkStates {
			nssaExternalLinkStateTableData = append(nssaExternalLinkStateTableData, []string{
				nssaExternalLinkStates.Route,
				nssaExternalLinkStates.MetricType,
				nssaExternalLinkStates.Base.AdvertisedRouter,
				strconv.Itoa(int(nssaExternalLinkStates.Base.LsaAge)),
			})
		}
		if nssaExternalLinkStateTableData == nil {
			amountOfnssaExternalLS = "0"
		} else {
			amountOfnssaExternalLS = strconv.Itoa(int(lsaTypes.NssaExternalLinkStatesCount))
		}

		// Order all Table Data
		common.SortTableByIPColumn(routerLinkStateTableData)
		common.SortTableByIPColumn(networkLinkStateTableData)
		common.SortTableByIPColumn(summaryLinkStateTableData)
		common.SortTableByIPColumn(asbrSummaryLinkStateTableData)
		common.SortTableByIPColumn(nssaExternalLinkStateTableData)

		// apply filters if active
		routerLinkStateTableData = common.FilterRows(routerLinkStateTableData, m.textFilter.Query)
		networkLinkStateTableData = common.FilterRows(networkLinkStateTableData, m.textFilter.Query)
		summaryLinkStateTableData = common.FilterRows(summaryLinkStateTableData, m.textFilter.Query)
		asbrSummaryLinkStateTableData = common.FilterRows(asbrSummaryLinkStateTableData, m.textFilter.Query)
		nssaExternalLinkStateTableData = common.FilterRows(nssaExternalLinkStateTableData, m.textFilter.Query)

		// Create Table for Router Link States and Fill with extracted routerLinkStateTableData
		rowsRouter := len(routerLinkStateTableData)
		routerLinkStateTable := components.NewOspfMonitorTable(
			[]string{
				"Advertised Router ID",
				"Router Links",
				"LSA Age",
			},
			rowsRouter,
		)
		for _, r := range routerLinkStateTableData {
			routerLinkStateTable = routerLinkStateTable.Row(r...)
		}

		// Create Table for Network Link States and Fill with extracted networkLinkStateTableData
		rowsNetwork := len(networkLinkStateTableData)
		networkLinkStateTable := components.NewOspfMonitorTable(
			[]string{
				"Designated Router ID",
				"Advertised Router ID",
				"LSA Age",
			},
			rowsNetwork,
		)
		for _, r := range networkLinkStateTableData {
			networkLinkStateTable = networkLinkStateTable.Row(r...)
		}

		// Create Table for Summary Link States and Fill with extracted summaryLinkStateTableData
		rowsSummary := len(summaryLinkStateTableData)
		summaryLinkStateTable := components.NewOspfMonitorTable(
			[]string{
				"Network ID",
				"Advertised Router ID",
				"LSA Age",
			},
			rowsSummary,
		)
		for _, r := range summaryLinkStateTableData {
			summaryLinkStateTable = summaryLinkStateTable.Row(r...)
		}

		// Create Table for AS Summary Link States and Fill with extracted asbrSummaryLinkStateTableData
		rowsAsSummary := len(asbrSummaryLinkStateTableData)
		asbrSummaryLinkStateTable := components.NewOspfMonitorTable(
			[]string{
				"AS Border Router ID",
				"Advertised Router ID",
				"LSA Age",
			},
			rowsAsSummary,
		)
		for _, r := range asbrSummaryLinkStateTableData {
			asbrSummaryLinkStateTable = asbrSummaryLinkStateTable.Row(r...)
		}

		// Create Table for NSSA External Link States and Fill with extracted nssaExternalLinkStateTableData
		rowsNSSAExternal := len(nssaExternalLinkStateTableData)
		nssaExternalLinkStateTable := components.NewOspfMonitorTable(
			[]string{
				"External Route",
				"Metric Type",
				"Advertising Router ID",
				"LSA Age",
			},
			rowsNSSAExternal,
		)
		for _, r := range nssaExternalLinkStateTableData {
			nssaExternalLinkStateTable = nssaExternalLinkStateTable.Row(r...)
		}

		areaHeader := styles.H1TitleStyleForOne().Render(fmt.Sprintf("Link State Database: Area %s", areaID))

		// create styled boxes for each LSA Type (type 1-4)
		routerTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyleForTwo().Render(amountOfRouterLS+" Router Link States"),
			styles.H2TwoContentBoxesCenterStyle().Render(routerLinkStateTable.String()),
			styles.H2TwoBoxBottomBorderStyle().Render(""),
		)
		networkTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyleForTwo().Render(amountOfNetworkLS+" Network Link States"),
			styles.H2TwoContentBoxesCenterStyle().Render(networkLinkStateTable.String()),
			styles.H2TwoBoxBottomBorderStyle().Render(""),
		)
		summaryTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyleForTwo().Render(amountOfSummaryLS+" Summary Link States"),
			styles.H2TwoContentBoxesCenterStyle().Render(summaryLinkStateTable.String()),
			styles.H2TwoBoxBottomBorderStyle().Render(""),
		)
		asbrSummaryTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyleForTwo().Render(amountOfAsSummaryLS+" ASBR Summary Link States"),
			styles.H2TwoContentBoxesCenterStyle().Render(asbrSummaryLinkStateTable.String()),
			styles.H2TwoBoxBottomBorderStyle().Render(""),
		)
		nssaExternalTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyleForOne().Render(amountOfnssaExternalLS+" NSSA External Link States"),
			styles.H2OneContentBoxCenterStyle().Render(nssaExternalLinkStateTable.String()),
			styles.H2OneBoxBottomBorderStyle().Render(""),
		)

		verticalRouterAndSummaryLinkStates := lipgloss.JoinVertical(lipgloss.Left, routerTableBox, summaryTableBox)
		verticalNetworkAndAsbrSummaryLinkStates := lipgloss.JoinVertical(lipgloss.Left, networkTableBox, asbrSummaryTableBox)

		type1to4Total := lipgloss.JoinHorizontal(lipgloss.Left,
			verticalRouterAndSummaryLinkStates,
			verticalNetworkAndAsbrSummaryLinkStates,
		)

		var optionalLSAType7 []string
		if nssaExternalLinkStateTableData != nil {
			optionalLSAType7 = append(optionalLSAType7, nssaExternalTableBox)
		}

		activeOptionalLSATypes := lipgloss.JoinVertical(lipgloss.Left,
			optionalLSAType7...,
		)

		completeAreaLSDB := lipgloss.JoinVertical(lipgloss.Left,
			areaHeader,
			type1to4Total,
			activeOptionalLSATypes,
		)

		lsdbBlocks = append(lsdbBlocks, completeAreaLSDB)
	}

	// ===== External LSA =====
	var asExternalLinkStateTableData [][]string
	var amountOfExternalLS string
	for _, asExternalLinkState := range lsdb.AsExternalLinkStates {
		asExternalLinkStateTableData = append(asExternalLinkStateTableData, []string{
			asExternalLinkState.Route,
			asExternalLinkState.MetricType,
			asExternalLinkState.Base.AdvertisedRouter,
			strconv.Itoa(int(asExternalLinkState.Base.LsaAge)),
		})
	}
	if asExternalLinkStateTableData == nil {
		amountOfExternalLS = "0"
	} else {
		amountOfExternalLS = strconv.Itoa(int(lsdb.AsExternalCount))
	}

	// Order all Table Data
	common.SortTableByIPColumn(asExternalLinkStateTableData)

	// apply filters if active
	asExternalLinkStateTableData = common.FilterRows(asExternalLinkStateTableData, m.textFilter.Query)

	// Create Table for External Link States and Fill with extracted asExternalLinkStateTableData
	rowsExternal := len(asExternalLinkStateTableData)
	asExternalLinkStateTable := components.NewOspfMonitorTable(
		[]string{
			"External Route",
			"Metric Type",
			"Advertising Router ID",
			"LSA Age",
		},
		rowsExternal,
	)
	for _, r := range asExternalLinkStateTableData {
		asExternalLinkStateTable = asExternalLinkStateTable.Row(r...)
	}

	externalHeader := styles.H1TitleStyleForOne().Render("Link State Database: AS External LSAs")

	// create styled boxes for each external LSA Type (type 5 & 7)
	externalTableBox := lipgloss.JoinVertical(lipgloss.Left,
		styles.H2TitleStyleForOne().Render(amountOfExternalLS+" AS External Link States"),
		styles.H2OneContentBoxCenterStyle().Render(asExternalLinkStateTable.String()),
		styles.H2OneBoxBottomBorderStyle().Render(""),
	)

	completeExternalLSDB := lipgloss.JoinVertical(lipgloss.Left,
		externalHeader,
		externalTableBox,
	)

	lsdbBlocks = append(lsdbBlocks, completeExternalLSDB+"\n\n")

	// Set viewport sizes and assign content to viewport
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = styles.HeightViewPortCompletePage - styles.BodyFooterHeight

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, lsdbBlocks...))

	return m.viewport.View()
}

func (m *Model) renderRouterMonitorTab() string {
	ospfNeighbors, err := backend.GetOspfNeighborInterfaces(m.logger)
	if err != nil {
		m.statusMessage = "Failed to fetch OSPF neighbor interfaces"
		m.statusSeverity = styles.SeverityError
		return common.PrintBackendError(err, "GetOspfNeighborInterfaces")
	}
	routerLSASelf, err := backend.GetOspfRouterDataSelf(m.logger)
	if err != nil {
		m.statusMessage = "Failed to fetch router LSA data"
		m.statusSeverity = styles.SeverityError
		return common.PrintBackendError(err, "GetOspfRouterDataSelf")
	}
	p2pInterfaceMap, err := backend.GetOspfP2PInterfaceMapping(m.logger)
	if err != nil {
		m.statusMessage = "Failed to fetch P2P Interface data"
		m.statusSeverity = styles.SeverityError
		return common.PrintBackendError(err, "GetOspfP2PInterfaceMapping")
	}

	// extract and sort the map keys (areas)
	routerLSAAreas := make([]string, 0, len(routerLSASelf.RouterStates))
	for area := range routerLSASelf.RouterStates {
		routerLSAAreas = append(routerLSAAreas, area)
	}
	sort.Strings(routerLSAAreas)

	var routerLSABlocks []string
	for _, areaID := range routerLSAAreas {
		areaData := routerLSASelf.RouterStates[areaID]

		var transitTableData [][]string
		var stubTableData [][]string
		var point2pointTableData [][]string

		for _, lsa := range areaData.LsaEntries {
			for _, link := range lsa.RouterLinks {
				if strings.Contains(link.LinkType, "Transit Network") {
					name := "No Neighbor"
					if link.DesignatedRouterAddress == link.RouterInterfaceAddress {
						name = "self"
					} else if common.ContainsString(ospfNeighbors, link.DesignatedRouterAddress) {
						name = "Neighbor"
					}
					transitTableData = append(transitTableData, []string{
						link.DesignatedRouterAddress,
						name,
						link.RouterInterfaceAddress,
						strconv.Itoa(int(lsa.LsaAge)),
					})
				} else if strings.Contains(link.LinkType, "Stub Network") {
					stubTableData = append(stubTableData, []string{
						link.NetworkAddress,
						link.NetworkMask,
						strconv.Itoa(int(lsa.LsaAge)),
					})
				} else if strings.Contains(link.LinkType, "point-to-point") {
					var mappedAddress string
					if addr, ok := p2pInterfaceMap.PeerInterfaceToAddress[link.RouterInterfaceAddress]; ok {
						mappedAddress = addr
					} else {
						mappedAddress = "no mapping"
					}

					point2pointTableData = append(point2pointTableData, []string{
						link.RouterInterfaceAddress,
						mappedAddress,
						strconv.Itoa(int(lsa.LsaAge)),
					})
				}
			}
		}

		// Order all Table Data
		common.SortTableByIPColumn(transitTableData)
		common.SortTableByIPColumn(stubTableData)
		common.SortTableByIPColumn(point2pointTableData)

		// apply filters if active
		transitTableData = common.FilterRows(transitTableData, m.textFilter.Query)
		stubTableData = common.FilterRows(stubTableData, m.textFilter.Query)
		point2pointTableData = common.FilterRows(point2pointTableData, m.textFilter.Query)

		rowsTransit := len(transitTableData)
		transitTable := components.NewOspfMonitorTable(
			[]string{
				"DR Address",
				"DR",
				"Interface Address",
				"LSA Age",
			},
			rowsTransit,
		)
		for _, r := range transitTableData {
			transitTable = transitTable.Row(r...)
		}

		rowsStub := len(stubTableData)
		stubTable := components.NewOspfMonitorTable(
			[]string{
				"Network Address",
				"Network Mask",
				"LSA Age",
			},
			rowsStub,
		)
		for _, r := range stubTableData {
			stubTable = stubTable.Row(r...)
		}

		rowsPoint2Point := len(stubTableData)
		point2pointTable := components.NewOspfMonitorTable(
			[]string{
				"Interface Address",
				"Translated Address",
				"LSA Age",
			},
			rowsPoint2Point,
		)
		for _, r := range point2pointTableData {
			point2pointTable = point2pointTable.Row(r...)
		}

		areaHeader := styles.H1TitleStyleForOne().Render(fmt.Sprintf("Area %s", areaID))

		var transitTableBox string
		if len(transitTableData) != 0 {
			transitTableBox = lipgloss.JoinVertical(lipgloss.Left,
				styles.H2TitleStyleForTwo().Render("Transit Networks"),
				styles.H2TwoContentBoxesCenterStyle().Render(transitTable.String()),
				styles.H2TwoBoxBottomBorderStyle().Render(""),
			)
		} else {
			transitTableBox = lipgloss.JoinVertical(lipgloss.Left,
				styles.H2TitleStyleForTwo().Render("No Transit Networks"),
				styles.H2TwoBoxBottomBorderStyle().Render(""),
			)
		}

		var stubTableBox string
		if len(stubTableData) != 0 {
			stubTableBox = lipgloss.JoinVertical(lipgloss.Left,
				styles.H2TitleStyleForTwo().Render("Stub Networks"),
				styles.H2TwoContentBoxesCenterStyle().Render(stubTable.String()),
				styles.H2TwoBoxBottomBorderStyle().Render(""),
			)
		} else {
			stubTableBox = lipgloss.JoinVertical(lipgloss.Left,
				styles.H2TitleStyleForTwo().Render("No Stub Networks"),
				styles.H2TwoBoxBottomBorderStyle().Render(""),
			)
		}

		var point2pointTableBox string
		if len(point2pointTableData) != 0 {
			point2pointTableBox = lipgloss.JoinVertical(lipgloss.Left,
				styles.H2TitleStyleForTwo().Render("Point-to-Point Networks"),
				styles.H2TwoContentBoxesCenterStyle().Render(point2pointTable.String()),
				styles.H2TwoBoxBottomBorderStyle().Render(""),
			)
		} else {
			point2pointTableBox = lipgloss.JoinVertical(lipgloss.Left,
				styles.H2TitleStyleForTwo().Render("No Point-to-Point Networks"),
				styles.H2TwoBoxBottomBorderStyle().Render(""),
			)
		}

		var verticalTables string
		var horizontalTables string

		if len(transitTableBox) < len(stubTableBox) {
			verticalTables = lipgloss.JoinVertical(lipgloss.Left, transitTableBox, point2pointTableBox)
			horizontalTables = lipgloss.JoinHorizontal(lipgloss.Top, verticalTables, stubTableBox)
		} else {
			verticalTables = lipgloss.JoinVertical(lipgloss.Left, stubTableBox, point2pointTableBox)
			horizontalTables = lipgloss.JoinHorizontal(lipgloss.Top, transitTableBox, verticalTables)
		}

		completeAreaRouterLSAs := lipgloss.JoinVertical(lipgloss.Left, areaHeader, horizontalTables)

		routerLSABlocks = append(routerLSABlocks, completeAreaRouterLSAs+"\n\n")
	}

	m.viewport.Width = styles.WidthViewPortCompletePage
	m.viewport.Height = styles.HeightViewPortCompletePage - styles.BodyFooterHeight

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, routerLSABlocks...))

	return m.viewport.View()
}

func (m *Model) renderNetworkMonitorTab() string {
	networkLSASelf, err := backend.GetOspfNetworkDataSelf(m.logger)
	if err != nil {
		m.statusMessage = "Failed to fetch network LSA data"
		m.statusSeverity = styles.SeverityError
		return common.PrintBackendError(err, "GetOspfRouterDataSelf")
	}
	routerName, _, err := backend.GetRouterName(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetRouterName")
	}

	// extract and sort the map keys (areas)
	networkLSAAreas := make([]string, 0, len(networkLSASelf.NetStates))
	for area := range networkLSASelf.NetStates {
		networkLSAAreas = append(networkLSAAreas, area)
	}
	sort.Strings(networkLSAAreas)

	var networkLSABlocks []string
	for _, areaID := range networkLSAAreas {
		areaData := networkLSASelf.NetStates[areaID]
		var networkTableData [][]string

		for lsaID, lsa := range areaData.LsaEntries {
			if lsaID == lsa.LinkStateId {
				var attachedRouterList []string
				for _, attachedRouter := range lsa.AttachedRouters {
					attachedRouterList = append(attachedRouterList, attachedRouter.AttachedRouterId)
				}
				networkTableData = append(networkTableData, []string{
					lsaID,
					strconv.Itoa(int(lsa.NetworkMask)),
					lsa.AdvertisingRouter,
					strings.Join(attachedRouterList, "\n"),
					strconv.Itoa(int(lsa.LsaAge)),
				})
			} else {
				return "Anomaly: LSA Mismatch"
			}
		}

		// Order all Table Data
		common.SortTableByIPColumn(networkTableData)

		// apply filters if active
		networkTableData = common.FilterRows(networkTableData, m.textFilter.Query)

		rowsNetwork := len(networkTableData)
		networkTable := components.NewOspfMonitorMultilineTable(
			[]string{
				"Link State ID",
				"CIDR",
				"Advertising Router",
				"Attached Routers",
				"LSA Age",
			},
			rowsNetwork,
		)
		for _, r := range networkTableData {
			networkTable = networkTable.Row(r...)
		}

		areaHeader := styles.H1TitleStyleForOne().Render(fmt.Sprintf("Area %s", areaID))

		networkTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyleForOne().Render("Network LSAs (Type 2)"),
			styles.H2OneContentBoxCenterStyle().Render(networkTable.String()),
			styles.H2OneBoxBottomBorderStyle().Render(""),
		)

		completeAreaNetworkLSAs := lipgloss.JoinVertical(lipgloss.Left, areaHeader, networkTableBox)

		if areaData.LsaEntries != nil {
			networkLSABlocks = append(networkLSABlocks, completeAreaNetworkLSAs+"\n\n")
		}
	}

	m.viewport.Width = styles.WidthViewPortCompletePage
	m.viewport.Height = styles.HeightViewPortCompletePage - styles.BodyFooterHeight

	if len(networkLSABlocks) == 0 {
		emptyContent := lipgloss.JoinVertical(lipgloss.Left,
			styles.H1TitleStyleForOne().Render(routerName+" does not originate Network LSAs (Type 2)"),
			lipgloss.NewStyle().Height(styles.HeightViewPortCompletePage-styles.BodyFooterHeight-styles.HeightH1).Render(""),
		)
		m.viewport.SetContent(emptyContent)
	} else {
		m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, networkLSABlocks...))
	}

	return m.viewport.View()
}

func (m *Model) renderExternalMonitorTab() string {
	var externalLsaBlock []string
	var nssaExternalLsaBlock []string

	externalLSASelf, err := backend.GetOspfExternalDataSelf(m.logger)
	if err != nil {
		m.statusMessage = "Failed to fetch external LSA data"
		m.statusSeverity = styles.SeverityError
		return common.PrintBackendError(err, "GetOspfExternalDataSelf")
	}
	nssaExternalDataSelf, err := backend.GetOspfNssaExternalDataSelf(m.logger)
	if err != nil {
		m.statusMessage = "Failed to fetch NSSA external LSA data"
		m.statusSeverity = styles.SeverityError
		return common.PrintBackendError(err, "GetOspfNssaExternalDataSelf")
	}
	routerName, _, err := backend.GetRouterName(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetRouterName")
	}

	// ===== OSPF External LSAs (Type 5) =====
	var externalTableData [][]string
	var externalTableDataExpanded [][]string // for future  feature
	for externalLinkState, linkStateData := range externalLSASelf.AsExternalLinkStates {
		externalTableData = append(externalTableData, []string{
			linkStateData.LinkStateId,
			"/" + strconv.Itoa(int(linkStateData.NetworkMask)),
			linkStateData.MetricType,
			linkStateData.ForwardAddress,
			strconv.Itoa(int(linkStateData.LsaAge)),
		})

		externalTableDataExpanded = append(externalTableDataExpanded, []string{
			externalLinkState,
			string(linkStateData.NetworkMask),
			linkStateData.MetricType,
		})
	}

	// Order all Table Data
	common.SortTableByIPColumn(externalTableData)

	// apply filters if active
	externalTableData = common.FilterRows(externalTableData, m.textFilter.Query)

	rowsExternal := len(externalTableData)
	externalTable := components.NewOspfMonitorTable([]string{
		"Link State ID",
		"CIDR",
		"Metric Type",
		"Forwarding Address",
		"LSA Age",
	},
		rowsExternal,
	)

	for _, r := range externalTableData {
		externalTable = externalTable.Row(r...)
	}

	externalHeader := styles.H1TitleStyleForOne().Render("External LSAs (Type 5)")

	externalDataBox := lipgloss.JoinVertical(lipgloss.Left,
		styles.H2TitleStyleForOne().Render("Self Originating"),
		styles.H2OneContentBoxCenterStyle().Render(externalTable.String()),
		styles.H2OneBoxBottomBorderStyle().Render(""),
	)

	var completeExternalBox string
	completeExternalBox = lipgloss.JoinVertical(lipgloss.Left, externalHeader, externalDataBox)

	externalLsaBlock = append(externalLsaBlock, completeExternalBox+"\n\n")

	// extract and sort the map keys
	nssaAreas := make([]string, 0, len(nssaExternalDataSelf.NssaExternalLinkStates))
	for area := range nssaExternalDataSelf.NssaExternalLinkStates {
		nssaAreas = append(nssaAreas, area)
	}
	sort.Strings(nssaAreas)

	hasNssaExternalLSAs := false
	for _, area := range nssaAreas {
		areaData := nssaExternalDataSelf.NssaExternalLinkStates[area]

		if areaData.Data != nil {
			var nssaExternalTableData [][]string
			for _, lsaData := range areaData.Data {
				nssaExternalTableData = append(nssaExternalTableData, []string{
					lsaData.LinkStateId,
					"/" + strconv.Itoa(int(lsaData.NetworkMask)),
					lsaData.MetricType,
					lsaData.NssaForwardAddress,
					strconv.Itoa(int(lsaData.LsaAge)),
				})
			}

			// Order all Table Data
			common.SortTableByIPColumn(nssaExternalTableData)

			// apply filters if active
			nssaExternalTableData = common.FilterRows(nssaExternalTableData, m.textFilter.Query)

			// create table for NSSA Exernal Link States with extracted data (nssaExternalTableData)
			rowsNssaExternal := len(nssaExternalTableData)
			nssaExternalTable := components.NewOspfMonitorTable(
				[]string{
					"Link State ID",
					"CIDR",
					"Metric Type",
					"Forwarding Address",
					"LSA Age",
				},
				rowsNssaExternal,
			)
			for _, r := range nssaExternalTableData {
				nssaExternalTable = nssaExternalTable.Row(r...)
			}

			nssaExternalHeader := styles.H1TitleStyleForOne().Render("NSSA External LSAs (Type 7) in Area " + area)

			nssaExternalDataBox := lipgloss.JoinVertical(lipgloss.Left,
				styles.H2TitleStyleForOne().Render("Self Originating"),
				styles.H2OneContentBoxCenterStyle().Render(nssaExternalTable.String()),
				styles.H2OneBoxBottomBorderStyle().Render(""),
			)
			completeNssaExternalBox := lipgloss.JoinVertical(lipgloss.Left, nssaExternalHeader, nssaExternalDataBox)

			nssaExternalLsaBlock = append(nssaExternalLsaBlock, completeNssaExternalBox+"\n\n")

			hasNssaExternalLSAs = true
		}
	}

	m.viewport.Width = styles.WidthViewPortCompletePage
	m.viewport.Height = styles.HeightViewPortCompletePage - styles.BodyFooterHeight

	var allLsaBlocks []string
	if hasNssaExternalLSAs == false {
		if externalTableData == nil {
			allLsaBlocks = allLsaBlocks[:0]
			allLsaBlocks = append(allLsaBlocks, lipgloss.JoinVertical(lipgloss.Left,
				styles.H1TitleStyleForOne().Render(routerName+" does not originate External LSAs (Type 5 or 7)"),
				lipgloss.NewStyle().Height(styles.HeightH1EmptyContentPadding).Render(""),
			))
		} else {
			allLsaBlocks = externalLsaBlock
		}
	} else {
		allLsaBlocks = append(externalLsaBlock, nssaExternalLsaBlock...)
	}

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, allLsaBlocks...))

	return m.viewport.View()
}

func (m *Model) renderNeighborMonitorTab() string {
	ospfNeighbors, err := backend.GetOspfNeighbors(m.logger)
	if err != nil {
		m.statusMessage = "Failed to fetch OSPF neighbor data"
		m.statusSeverity = styles.SeverityError
		return common.PrintBackendError(err, "GetOspfNeighborInterfaces")
	}

	routerName, _, err := backend.GetRouterName(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetRouterName")
	}

	// extract and sort the map keys
	ospfNeighborIDs := make([]string, 0, len(ospfNeighbors.Neighbors))
	for neighborID := range ospfNeighbors.Neighbors {
		ospfNeighborIDs = append(ospfNeighborIDs, neighborID)
	}
	list := common.SortedIpList(ospfNeighborIDs)
	sort.Sort(&list)

	var ospfNeighborTableData [][]string
	for _, ospfNeighborID := range ospfNeighborIDs {
		ospfNeighborList := ospfNeighbors.Neighbors[ospfNeighborID]

		for _, ospfNeighbor := range ospfNeighborList.Neighbors {
			ospfNeighborTableData = append(ospfNeighborTableData, []string{
				ospfNeighborID,
				ospfNeighbor.IfaceAddress,
				ospfNeighbor.Role,
				ospfNeighbor.Converged,
				ospfNeighbor.IfaceName,
				ospfNeighbor.UpTime,
				ospfNeighbor.DeadTime,
			})
		}

	}

	// Order all Table Data
	common.SortTableByIPColumn(ospfNeighborTableData)

	// apply filters if active
	ospfNeighborTableData = common.FilterRows(ospfNeighborTableData, m.textFilter.Query)

	// Create Table for NSSA External Link States and Fill with extracted nssaExternalLinkStateTableData
	rowsOspfNeighbors := len(ospfNeighborTableData)
	ospfNeighborTable := components.NewOspfMonitorTable(
		[]string{
			"Neighbor ID",
			"Neighbor IP",
			"Role",
			"Converged",
			"Internal Interface",
			"Up Time",
			"Dead Time",
		},
		rowsOspfNeighbors,
	)
	for _, r := range ospfNeighborTableData {
		ospfNeighborTable = ospfNeighborTable.Row(r...)
	}

	ospfNeghborHeader := styles.H1TitleStyleForOne().Render("All OSPF Neighborships")

	// create styled boxes for each external LSA Type (type 5 & 7)
	ospfNeighborTableBox := lipgloss.JoinVertical(lipgloss.Left,
		styles.H2TitleStyleForOne().Render("Router "+routerName+" has "+strconv.Itoa(len(ospfNeighborIDs))+" Neighbors"),
		styles.H2OneContentBoxCenterStyle().Render(ospfNeighborTable.String()),
		styles.H2OneBoxBottomBorderStyle().Render(""),
	)

	m.viewport.Width = styles.WidthViewPortCompletePage
	m.viewport.Height = styles.HeightViewPortCompletePage - styles.BodyFooterHeight

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, ospfNeghborHeader, ospfNeighborTableBox))

	return m.viewport.View()
}

func (m *Model) renderRunningConfigTab() string {
	runningConfigTitle := styles.H1TitleStyleForTwo().Render("Running Config")
	formatedRunningConfigOutput := strings.Join(m.runningConfig, "\n")
	runningConfigBox := styles.H1TwoContentBoxesStyle().Render(formatedRunningConfigOutput)
	completeRunningConfig := lipgloss.JoinVertical(lipgloss.Left,
		runningConfigTitle,
		runningConfigBox,
		styles.H1TwoBoxBottomBorderStyle().Render(""),
	)

	staticFRRConfigTitle := styles.H1TitleStyleForTwo().Render("Parsed Running Config")
	staticFRRConfiguration, err := backend.GetStaticFRRConfigurationPretty(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetStaticFRRConfigurationPretty")
	}
	staticFileBox := styles.H1TwoContentBoxesStyle().Render(staticFRRConfiguration)
	completeStaticConfig := lipgloss.JoinVertical(lipgloss.Left,
		staticFRRConfigTitle,
		staticFileBox,
		styles.H1TwoBoxBottomBorderStyle().Render(""),
	)

	completeContent := lipgloss.JoinHorizontal(lipgloss.Top, completeRunningConfig, completeStaticConfig)

	m.viewport.Width = styles.WidthViewPortCompletePage
	m.viewport.Height = styles.HeightViewPortCompletePage
	m.viewport.SetContent(completeContent)

	return m.viewport.View()
}
