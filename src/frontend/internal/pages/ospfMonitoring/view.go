package ospfMonitoring

import (
	"fmt"
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"google.golang.org/protobuf/encoding/protojson"
	"sort"
	"strconv"

	// frrProto "github.com/ba2025-ysmprc/frr-mad/src/backend/pkg"
	backend "github.com/ba2025-ysmprc/frr-tui/internal/services"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/components"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/ba2025-ysmprc/frr-tui/pkg"
	// "github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
	"strings"
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
		return m.renderExternalMonitorTab()
	} else if currentSubTabLocal == 3 {
		return m.renderRunningConfigTab()
	}
	return m.renderLsdbMonitorTab()
}

func (m *Model) renderLsdbMonitorTab() string {
	var lsdbBlocks []string

	lsdb, _ := getLSDB()

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

		var amountOfRouterLS string
		var amountOfNetworkLS string
		var amountOfSummaryLS string
		var amountOfAsSummaryLS string

		// loop through LSAs (type 1-4) and extract data for tables
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

		// Order all Table Data
		sort.Slice(routerLinkStateTableData, func(i, j int) bool {
			return routerLinkStateTableData[i][0] < routerLinkStateTableData[j][0]
		})
		sort.Slice(networkLinkStateTableData, func(i, j int) bool {
			return networkLinkStateTableData[i][0] < networkLinkStateTableData[j][0]
		})
		sort.Slice(summaryLinkStateTableData, func(i, j int) bool {
			return summaryLinkStateTableData[i][0] < summaryLinkStateTableData[j][0]
		})
		sort.Slice(asbrSummaryLinkStateTableData, func(i, j int) bool {
			return asbrSummaryLinkStateTableData[i][0] < asbrSummaryLinkStateTableData[j][0]
		})

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

		areaHeader := styles.H1TitleStyleForOne().Render(fmt.Sprintf("Link State Database: Area %s", areaID))

		// create styled boxes for each LSA Type (type 1-4)
		routerTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyleForTwo().Render(amountOfRouterLS+" Router Link States"+strconv.Itoa(styles.WidthTwoH2Box)),
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

		horizontalRouterAndNetworkLinkStates := lipgloss.JoinHorizontal(lipgloss.Top, routerTableBox, networkTableBox)
		horizontalSummaryAndAsbrSummaryLinkStates := lipgloss.JoinHorizontal(lipgloss.Top, summaryTableBox, asbrSummaryTableBox)

		completeAreaLSDB := lipgloss.JoinVertical(lipgloss.Left,
			areaHeader,
			horizontalRouterAndNetworkLinkStates,
			horizontalSummaryAndAsbrSummaryLinkStates,
		)

		lsdbBlocks = append(lsdbBlocks, completeAreaLSDB+"\n\n")
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
	ospfNeighbors := getOspfNeighborInterfaces()
	routerLSASelf, _ := getOspfRouterData()

	// extract and sort the map keys (areas)
	routerLSAAreas := make([]string, 0, len(routerLSASelf.RouterStates))
	for area := range routerLSASelf.RouterStates {
		routerLSAAreas = append(routerLSAAreas, area)
	}
	sort.Strings(routerLSAAreas)

	var routerLSABlocks []string
	for _, areaID := range routerLSAAreas {
		areaData := routerLSASelf.RouterStates[areaID]
		var transitData [][]string
		var stubData [][]string

		for _, lsa := range areaData.LsaEntries {
			for _, link := range lsa.RouterLinks {
				if strings.Contains(link.LinkType, "Transit Network") {
					name := "No Neighbor"
					if link.DesignatedRouterAddress == link.RouterInterfaceAddress {
						name = "self"
					} else if common.ContainsString(ospfNeighbors, link.DesignatedRouterAddress) {
						name = "Direct Neighbor"
					}
					transitData = append(transitData, []string{
						link.DesignatedRouterAddress,
						name,
						link.RouterInterfaceAddress,
					})
				} else if strings.Contains(link.LinkType, "Stub Network") {
					stubData = append(stubData, []string{
						link.NetworkAddress,
						link.NetworkMask,
					})
				}
			}
		}

		// Order all Table Data
		sort.Slice(transitData, func(i, j int) bool {
			return transitData[i][0] < transitData[j][0]
		})
		sort.Slice(stubData, func(i, j int) bool {
			return stubData[i][0] < stubData[j][0]
		})

		rowsTransit := len(transitData)
		transitTable := components.NewOspfMonitorTable(
			[]string{
				"Link ID (DR Adr.)",
				"Designated Router",
				"Link Data (own Adr.)",
			},
			rowsTransit,
		)
		for _, r := range transitData {
			transitTable = transitTable.Row(r...)
		}

		rowsStub := len(stubData)
		stubTable := components.NewOspfMonitorTable(
			[]string{
				"Network Address",
				"Network Mask",
			},
			rowsStub,
		)
		for _, r := range stubData {
			stubTable = stubTable.Row(r...)
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

		horizontalTables := lipgloss.JoinHorizontal(lipgloss.Top, transitTableBox, stubTableBox)

		completeAreaRouterLSAs := lipgloss.JoinVertical(lipgloss.Left, areaHeader, horizontalTables)

		routerLSABlocks = append(routerLSABlocks, completeAreaRouterLSAs+"\n\n")
	}

	contentMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = contentMaxHeight

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, routerLSABlocks...))

	return m.viewport.View()
}

func (m *Model) renderExternalMonitorTab() string {
	var externalLsaBlock []string
	var nssaExternalLsaBlock []string

	externalLSASelf, _ := getOspfExternalData()
	nssaExternalDataSelf, _ := getOspfNssaExternalData()

	var externalTableData [][]string
	var externalTableDataExpanded [][]string // for future  feature
	for externalLinkState, linkStateData := range externalLSASelf.AsExternalLinkStates {
		externalTableData = append(externalTableData, []string{
			linkStateData.LinkStateId,
			"/" + strconv.Itoa(int(linkStateData.NetworkMask)),
			linkStateData.MetricType,
			linkStateData.ForwardAddress,
		})

		externalTableDataExpanded = append(externalTableDataExpanded, []string{
			externalLinkState,
			string(linkStateData.NetworkMask),
			linkStateData.MetricType,
		})
	}

	// Order all Table Data
	sort.Slice(externalTableData, func(i, j int) bool {
		return externalTableData[i][0] < externalTableData[j][0]
	})

	rowsExternal := len(externalTableData)
	externalTable := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		Headers("Link State ID", "CIDR", "Metric Type", "Forwarding Address").
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == ltable.HeaderRow:
				return styles.HeaderStyle
			case row == rowsExternal-1:
				return styles.NormalCellStyle.BorderBottom(true)
			default:
				return styles.NormalCellStyle
			}
		})

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

	var nssaExternalTableData [][]string
	for area, areaData := range nssaExternalDataSelf.NssaExternalLinkStates {
		for _, lsaData := range areaData.Data {
			nssaExternalTableData = append(nssaExternalTableData, []string{
				lsaData.LinkStateId,
				"/" + strconv.Itoa(int(lsaData.NetworkMask)),
				lsaData.MetricType,
				lsaData.NssaForwardAddress,
			})
		}

		// Order all Table Data
		sort.Slice(nssaExternalTableData, func(i, j int) bool {
			return nssaExternalTableData[i][0] < nssaExternalTableData[j][0]
		})

		rowsNssaExternal := len(nssaExternalTableData)
		nssaExternalTable := components.NewOspfMonitorTable(
			[]string{
				"Link State ID",
				"CIDR",
				"Metric Type",
				"Forwarding Address",
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
	}

	contentMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight
	m.viewport.Width = styles.WidthBasis
	m.viewport.Height = contentMaxHeight

	var allLsaBlocks []string
	if nssaExternalTableData == nil {
		if externalTableData == nil {
			allLsaBlocks = allLsaBlocks[:0]
			allLsaBlocks = append(allLsaBlocks, lipgloss.JoinVertical(lipgloss.Left,
				externalHeader,
				"no self originating external advertisements",
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

//func (m *Model) renderAdvertisementTab() string {
//	// on this view:
//	// show three advertisement boxes: based on vtysh LSA queries / based on file analysis / based on FIB analysis
//
//	boxWidthForTwo := (m.windowSize.Width - 10) / 2 // - 6 (padding+margin content) - 2 (for gap) - 2 (for border)
//	if boxWidthForTwo < 20 {
//		boxWidthForTwo = 20 // Minimum width to ensure readability
//	}
//
//	shouldAdvertisedTitle := styles.BoxTitleStyle.Render("Should be Advertised")
//
//	shouldAdvertisedContent := "Area 0.0.0.0: \n"
//
//	shouldAdvertisedRouterLSA := styles.H1TitleStyle.
//		Width(boxWidthForTwo - 2).
//		Render("Area 0.0.0.0, Router LSAs (Type 1)")
//	shouldAdvertisedContent += shouldAdvertisedRouterLSA
//
//	shouldAdvertisedVerticalStyle := lipgloss.JoinVertical(lipgloss.Left, shouldAdvertisedTitle, shouldAdvertisedContent)
//	shouldAdvertisedBox := lipgloss.NewStyle().Render(shouldAdvertisedVerticalStyle)
//
//	// ----------------------------------------------------
//	gap := 2
//	// ----------------------------------------------------
//
//	isAdvertisedTitle := styles.BoxTitleStyle.Render("Is Advertised")
//
//	isAdvertisedContent := "Area 0.0.0.0: \n"
//
//	// for each area create area box --> need length of areas
//
//	isAdvertisedRouterLSA := styles.H1TitleStyle.
//		Width(boxWidthForTwo - 2).
//		Render("Area 0.0.0.0, Router LSAs (Type 1)")
//	isAdvertisedContent += isAdvertisedRouterLSA
//
//	isAdvertisedVerticalStyle := lipgloss.JoinVertical(lipgloss.Left, isAdvertisedTitle, isAdvertisedContent)
//	isAdvertisedBox := lipgloss.NewStyle().Render(isAdvertisedVerticalStyle)
//	// returnString := "Advertisement"
//
//	horizontalBoxes := lipgloss.JoinHorizontal(lipgloss.Top,
//		isAdvertisedBox,
//		lipgloss.NewStyle().Width(gap).Render(""),
//		shouldAdvertisedBox,
//	)
//
//	return horizontalBoxes
//}

//func (m *Model) renderOSPFTab0() string {
//	// Calculate box width dynamically for four horizontal boxes based on terminal width
//	boxWidthForFour := (m.windowSize.Width - 16) / 4 // - 6 (padding+margin content) - 10 (for each border)
//	if boxWidthForFour < 20 {
//		boxWidthForFour = 20 // Minimum width to ensure readability
//	}
//
//	ospfAnomalyOne := styles.GeneralBoxStyle.
//		Width(boxWidthForFour).
//		Render(styles.BoxTitleStyle.Render("OSPF Anomaly One") + "\n" + "Call Backend...☎\nEverything Good! amount")
//
//	ospfAnomalyTwo := styles.GeneralBoxStyle.
//		Width(boxWidthForFour).
//		Render(styles.BoxTitleStyle.Render("OSPF Anomaly Two") + "\n" + "Call Backend...☎\nEverything Good!")
//
//	ospfAnomalyThree := styles.BadBoxStyle.
//		Width(boxWidthForFour).
//		Render(styles.BoxTitleStyle.Render("OSPF Anomaly Three") + "\n" + "Call Backend...☎\nVery Bad Anomaly Detected!\n\nReport...\nReport...\nReport...\nReport...\nReport...\n")
//
//	ospfAnomalyFour := styles.GeneralBoxStyle.
//		Width(boxWidthForFour).
//		Render(styles.BoxTitleStyle.Render("OSPF Anomaly Four") + "\n" + "Call Backend...☎\nEverything Good!")
//
//	ospfAnomalies := []struct {
//		Title   string
//		Content string
//		Style   lipgloss.Style
//	}{
//		{
//			Title:   "OSPF Anomaly One",
//			Content: "Call Backend...☎\nVery Bad Anomaly Detected!\n\nReport...\nReport...\nReport...\nReport...\nReport...\n",
//			Style:   styles.BadBoxStyle,
//		},
//		{
//			Title:   "OSPF Anomaly Two",
//			Content: "Call Backend...☎\nEverything Good!",
//			Style:   styles.GeneralBoxStyle,
//		},
//		{
//			Title:   "OSPF Anomaly Three",
//			Content: "Call Backend...☎\nEverything Good!",
//			Style:   styles.GeneralBoxStyle,
//		},
//		{
//			Title:   "OSPF Anomaly Four",
//			Content: "Call Backend...☎\nEverything Good!",
//			Style:   styles.GeneralBoxStyle,
//		},
//	}
//
//	// Build anomaly boxes using the new component
//	var ospfAnomalyBoxes []string
//	for _, a := range ospfAnomalies {
//		box := components.NewAnomalyBox(a.Title, a.Content, a.Style, boxWidthForFour)
//		ospfAnomalyBoxes = append(ospfAnomalyBoxes, box.Render())
//	}
//
//	horizontalBoxes := lipgloss.JoinHorizontal(lipgloss.Top, ospfAnomalyOne, ospfAnomalyTwo, ospfAnomalyThree, ospfAnomalyFour)
//	horizontalBoxes2 := lipgloss.JoinHorizontal(lipgloss.Top, ospfAnomalyBoxes...)
//
//	//infoBox := styles.InfoTextStyle.
//	//	Width(m.windowSize.Width - 12).
//	//	Render("press 'r' to refresh ospf anomalies")
//	//
//	//return lipgloss.JoinVertical(lipgloss.Left, horizontalBoxes, infoBox)
//
//	return lipgloss.JoinVertical(lipgloss.Left, horizontalBoxes, horizontalBoxes2)
//}

//func (m *Model) renderOSPFTab1() string {
//	// Calculate box width dynamically for four horizontal boxes based on terminal width
//	boxWidthForFour := (m.windowSize.Width - 16) / 4 // - 6 (padding+margin content) - 10 (for each border)
//	if boxWidthForFour < 20 {
//		boxWidthForFour = 20 // Minimum width to ensure readability
//	}
//
//	ospfAnomalyOne := styles.GeneralBoxStyle.
//		Width(boxWidthForFour).
//		Render(styles.BoxTitleStyle.Render("OSPF Anomaly One") + "\n" + "Call Backend...☎\nEverything Good! amount")
//
//	ospfAnomalyTwo := styles.GeneralBoxStyle.
//		Width(boxWidthForFour).
//		Render(styles.BoxTitleStyle.Render("OSPF Anomaly Two") + "\n" + "Call Backend...☎\nEverything Good!")
//
//	ospfAnomalyThree := styles.BadBoxStyle.
//		Width(boxWidthForFour).
//		Render(styles.BoxTitleStyle.Render("OSPF Anomaly Three") + "\n" + "Call Backend...☎\nVery Bad Anomaly Detected!\n\nReport...\nReport...\nReport...\nReport...\nReport...\n")
//
//	ospfAnomalyFour := styles.GeneralBoxStyle.
//		Width(boxWidthForFour).
//		Render(styles.BoxTitleStyle.Render("OSPF Anomaly Four") + "\n" + "Call Backend...☎\nEverything Good!")
//
//	return lipgloss.JoinHorizontal(lipgloss.Top, ospfAnomalyThree, ospfAnomalyOne, ospfAnomalyTwo, ospfAnomalyFour)
//}

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
	staticFRRConfiguration := getStaticFRRConfigurationPretty()
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

// ============================== //
// HELPERS: BACKEND CALLS         //
// ============================== //

func getLSDB() (*pkg.OSPFDatabase, error) {
	response, err := backend.SendMessage("ospf", "database", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfDatabase(), nil
}

func getOspfRouterData() (*pkg.OSPFRouterData, error) {
	response, err := backend.SendMessage("ospf", "router", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfRouterData(), nil
}

func getOspfNeighborInterfaces() []string {
	response, err := backend.SendMessage("ospf", "neighbors", nil)
	if err != nil {
		return nil
	}
	ospfNeighbors := response.Data.GetOspfNeighbors()

	var neighborAddresses []string
	for _, neighborGroup := range ospfNeighbors.Neighbors {
		for _, neighbor := range neighborGroup.Neighbors {
			neighborAddresses = append(neighborAddresses, neighbor.IfaceAddress)
		}
	}

	return neighborAddresses
}

func getOspfExternalData() (*pkg.OSPFExternalData, error) {
	response, err := backend.SendMessage("ospf", "externalData", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfExternalData(), nil
}

func getOspfNssaExternalData() (*pkg.OSPFNssaExternalData, error) {
	response, err := backend.SendMessage("ospf", "nssaExternalData", nil)
	if err != nil {
		return nil, err
	}

	return response.Data.GetOspfNssaExternalData(), nil
}

func getStaticFRRConfigurationPretty() string {
	response, err := backend.SendMessage("ospf", "staticConfig", nil)
	if err != nil {
		return ""
	}

	var prettyJson string

	// Pretty‑print the protobuf into nice indented JSON
	marshaler := protojson.MarshalOptions{
		Multiline:     true,
		Indent:        "  ",
		UseProtoNames: true,
	}
	pretty, perr := marshaler.Marshal(response.Data)
	if perr != nil {
		prettyJson = response.Data.String()
	} else {
		prettyJson = string(pretty)
	}

	return prettyJson
}
