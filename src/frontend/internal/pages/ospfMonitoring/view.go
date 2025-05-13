package ospfMonitoring

import (
	"fmt"
	"sort"
	"strconv"

	"github.com/ba2025-ysmprc/frr-tui/internal/common"

	// frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	backend "github.com/ba2025-ysmprc/frr-tui/internal/services"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/components"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"

	// "github.com/charmbracelet/bubbles/table"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var currentSubTabLocal = -1

func (m *Model) OSPFView(currentSubTab int) string {
	currentSubTabLocal = currentSubTab
	return m.View()
}

func (m *Model) View() string {
	if currentSubTabLocal == 0 {
		return m.renderLsdbMonitorTab()
	} else if currentSubTabLocal == 1 {
		return m.renderRouterMonitorTab()
	} else if currentSubTabLocal == 2 {
		return m.renderNetworkMonitorTab()
	} else if currentSubTabLocal == 3 {
		return m.renderExternalMonitorTab()
	} else if currentSubTabLocal == 4 {
		return m.renderNeighborMonitorTab()
	} else if currentSubTabLocal == 5 {
		return m.renderRunningConfigTab()
	}
	return m.renderLsdbMonitorTab()
}

func (m *Model) renderLsdbMonitorTab() string {
	var lsdbBlocks []string

	lsdb, err := backend.GetLSDB(m.logger)
	if err != nil {
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

		horizontalRouterAndNetworkLinkStates := lipgloss.JoinHorizontal(lipgloss.Top, routerTableBox, networkTableBox)
		horizontalSummaryAndAsbrSummaryLinkStates := lipgloss.JoinHorizontal(lipgloss.Top, summaryTableBox, asbrSummaryTableBox)

		var optionalLSAType7 []string
		if nssaExternalLinkStateTableData != nil {
			optionalLSAType7 = append(optionalLSAType7, nssaExternalTableBox)
		}

		activeOptionalLSATypes := lipgloss.JoinVertical(lipgloss.Left,
			optionalLSAType7...,
		)

		completeAreaLSDB := lipgloss.JoinVertical(lipgloss.Left,
			areaHeader,
			horizontalRouterAndNetworkLinkStates,
			horizontalSummaryAndAsbrSummaryLinkStates,
			activeOptionalLSATypes,
		)

		lsdbBlocks = append(lsdbBlocks, completeAreaLSDB+"\n")
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
	contentMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = contentMaxHeight

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, lsdbBlocks...))

	return m.viewport.View()
}

func (m *Model) renderRouterMonitorTab() string {
	ospfNeighbors, err := backend.GetOspfNeighborInterfaces(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetOspfNeighborInterfaces")
	}
	routerLSASelf, err := backend.GetOspfRouterDataSelf(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetOspfRouterDataSelf")
	}
	p2pInterfaceMap, err := backend.GetOspfP2PInterfaceMapping(m.logger)
	if err != nil {
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

		transitTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyleForTwo().Render("Transit Networks"),
			styles.H2TwoContentBoxesCenterStyle().Render(transitTable.String()),
			styles.H2TwoBoxBottomBorderStyle().Render(""),
		)
		stubTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyleForTwo().Render("Stub Networks"),
			styles.H2TwoContentBoxesCenterStyle().Render(stubTable.String()),
			styles.H2TwoBoxBottomBorderStyle().Render(""),
		)
		point2pointTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyleForTwo().Render("Point-to-Point Networks"),
			styles.H2TwoContentBoxesCenterStyle().Render(point2pointTable.String()),
			styles.H2TwoBoxBottomBorderStyle().Render(""),
		)

		var verticalTables string
		if len(transitTableBox) < len(stubTableBox) {
			verticalTables = lipgloss.JoinVertical(lipgloss.Left, transitTableBox, point2pointTableBox)
		} else {
			verticalTables = lipgloss.JoinVertical(lipgloss.Left, stubTableBox, point2pointTableBox)
		}

		horizontalTables := lipgloss.JoinHorizontal(lipgloss.Top, transitTableBox, verticalTables)

		completeAreaRouterLSAs := lipgloss.JoinVertical(lipgloss.Left, areaHeader, horizontalTables)

		routerLSABlocks = append(routerLSABlocks, completeAreaRouterLSAs+"\n\n")
	}

	contentMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = contentMaxHeight

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, routerLSABlocks...))

	return m.viewport.View()
}

func (m *Model) renderNetworkMonitorTab() string {
	networkLSASelf, err := backend.GetOspfNetworkDataSelf(m.logger)
	if err != nil {
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
				// TODO: print a pretty anomaly error when lsaID != linkstateID
				return "Anomaly: LSA Mismatch"
			}
		}

		// Order all Table Data
		common.SortTableByIPColumn(networkTableData)

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

	if networkLSABlocks == nil {
		return lipgloss.JoinHorizontal(lipgloss.Left,
			styles.H1TitleStyleForOne().Render(routerName+" does not originate Network LSAs (Type 2)"))
	}

	contentMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = contentMaxHeight

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, networkLSABlocks...))

	return m.viewport.View()
}

func (m *Model) renderExternalMonitorTab() string {
	var externalLsaBlock []string
	var nssaExternalLsaBlock []string

	externalLSASelf, err := backend.GetOspfExternalDataSelf(m.logger)
	if err != nil {
		return common.PrintBackendError(err, "GetOspfExternalDataSelf")
	}
	nssaExternalDataSelf, err := backend.GetOspfNssaExternalDataSelf(m.logger)
	if err != nil {
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
			// var completeNssaExternalBox string
			completeNssaExternalBox := lipgloss.JoinVertical(lipgloss.Left, nssaExternalHeader, nssaExternalDataBox)

			nssaExternalLsaBlock = append(nssaExternalLsaBlock, completeNssaExternalBox+"\n\n")

			hasNssaExternalLSAs = true
		}
	}

	contentMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = contentMaxHeight

	var allLsaBlocks []string
	if hasNssaExternalLSAs == false {
		if externalTableData == nil {
			allLsaBlocks = allLsaBlocks[:0]
			allLsaBlocks = append(allLsaBlocks, lipgloss.JoinVertical(lipgloss.Left,
				styles.H1TitleStyleForOne().Render(routerName+" does not originate External LSAs (Type 5 or 7)"),
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

	contentMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = contentMaxHeight

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

	// completeColoredContent := lipgloss.NewStyle().Foreground(lipgloss.Color("#ffffff")).Render(completeContent)
	outputMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = outputMaxHeight
	m.viewport.SetContent(completeContent)

	// runningConfigBox := lipgloss.NewStyle().Padding(0, 5).Render(m.viewport.View())

	return m.viewport.View()
}
