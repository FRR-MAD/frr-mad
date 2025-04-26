package ospfMonitoring

import (
	"fmt"
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
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
		return m.renderAdvertisementTab()
	} else if currentSubTabLocal == 1 {
		return m.renderRouterMonitorTab()
	} else if currentSubTabLocal == 2 {
		return m.renderExternalMonitorTab()
	} else if currentSubTabLocal == 3 {
		return m.renderOSPFTab1()
	} else if currentSubTabLocal == 4 {
		return m.renderRunningConfigTab()
	}
	return m.renderAdvertisementTab()
}

//func (m *Model) renderLsdbMonitorTab() string {
//	boxWidthForOneH1 := m.windowSize.Width - 10       // - 6 (padding+margin content) - 2 (for each border)
//	boxWidthForTwoH2 := (m.windowSize.Width - 16) / 2 // - 6 (padding+margin content) - 2 (for border) - 8 (for margin)
//
//	var lsdbBlocks []string
//
//	lsdb, _ := getLSDB()
//
//	var databaseBlocks []string
//	for area, lsaTypes := range lsdb.Areas {
//		var routerLinkStateData [][]string
//		var networkLinkStateData [][]string
//		var summaryLinkStateData [][]string
//		var asbrSummaryLinkStateData [][]string
//		for _, routerLinkState := range lsaTypes.RouterLinkStates {
//			routerLinkStateData = append(routerLinkStateData, []string{
//				routerLinkState.Base.AdvertisedRouter,
//				strconv.Itoa(int(routerLinkState.Base.LsaAge)),
//				strconv.Itoa(int(routerLinkState.NumOfRouterLinks)),
//			})
//		}
//		for _, networkLinkState := range lsaTypes.NetworkLinkStates {
//			networkLinkStateData = append(networkLinkStateData, []string{
//				networkLinkState.Base.LsId,
//				networkLinkState.Base.AdvertisedRouter,
//				strconv.Itoa(int(networkLinkState.Base.LsaAge)),
//			})
//		}
//		for _, summarLinkState := range lsaTypes.SummaryLinkStates {
//			summaryLinkStateData = append(summaryLinkStateData, []string{
//				summarLinkState.SummaryAddress,
//				summarLinkState.Base.AdvertisedRouter,
//				strconv.Itoa(int(summarLinkState.Base.LsaAge)),
//			})
//		}
//		for _, asbrSummaryLinkState := range lsaTypes.AsbrSummaryLinkStates {
//			asbrSummaryLinkStateData = append(asbrSummaryLinkStateData, []string{
//				asbrSummaryLinkState.Base.LsId,
//				asbrSummaryLinkState.Base.AdvertisedRouter,
//				asbrSummaryLinkState.Base.LsaAge
//			})
//		}
//
//		areaHeader := styles.ContentTitleH1Style.
//			Width(boxWidthForOneH1).
//			Margin(0, 0, 1, 0).
//			Padding(1, 0, 0, 0).
//			Render(fmt.Sprintf("Area %s", area))
//	}
//	return "nothing"
//}

func (m *Model) renderRouterMonitorTab() string {
	boxWidthForOneH1 := m.windowSize.Width - 10 // - 6 (padding+margin content) - 2 (for each border)
	if boxWidthForOneH1 < 20 {
		boxWidthForOneH1 = 20 // Minimum width to ensure readability
	}
	boxWidthForTwoH2 := (m.windowSize.Width - 16) / 2 // - 6 (padding+margin content) - 2 (for border) - 8 (for margin)
	if boxWidthForTwoH2 < 20 {
		boxWidthForTwoH2 = 20 // Minimum width to ensure readability
	}

	ospfNeighbors := getOspfNeighborInterfaces()
	routerLSASelf, _ := getOspfRouterData()

	var routerLSABlocks []string
	for area, areaData := range routerLSASelf.RouterStates {
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

		rowsTransit := len(transitData)
		transitTable := ltable.New().
			Border(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true).
			BorderHeader(true).
			BorderColumn(true).
			Headers("Link ID (DR Adr.)", "Designated Router", "Link Data (own Adr.)").
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == ltable.HeaderRow:
					return styles.HeaderStyle
				case row == rowsTransit-1:
					return styles.NormalCellStyle.BorderBottom(true)
				default:
					return styles.NormalCellStyle
				}
			})

		for _, r := range transitData {
			transitTable = transitTable.Row(r...)
		}

		rowsStub := len(stubData)
		stubTable := ltable.New().
			Border(lipgloss.NormalBorder()).
			BorderTop(true).
			BorderBottom(true).
			BorderLeft(true).
			BorderRight(true).
			BorderHeader(true).
			BorderColumn(true).
			Headers("Network Address", "Network Mask").
			StyleFunc(func(row, col int) lipgloss.Style {
				switch {
				case row == ltable.HeaderRow:
					return styles.HeaderStyle
				case row == rowsStub-1:
					return styles.NormalCellStyle
				default:
					return styles.NormalCellStyle
				}
			})

		for _, r := range stubData {
			stubTable = stubTable.Row(r...)
		}

		areaHeader := styles.ContentTitleH1Style.
			Width(boxWidthForOneH1).
			Margin(0, 0, 1, 0).
			Padding(1, 0, 0, 0).
			Render(fmt.Sprintf("Area %s", area))

		correctBoxWidthTransit := lipgloss.JoinVertical(lipgloss.Left,
			styles.ContentTitleH2Style.Width(boxWidthForTwoH2-2).Render("Transit Networks"),
			lipgloss.NewStyle().Align(lipgloss.Center).Margin(0, 2).Width(boxWidthForTwoH2).Render(transitTable.String()),
			styles.ContentBottomBorderStyle.Width(boxWidthForTwoH2-2).Render(""),
		)
		correctBoxWidthStub := lipgloss.JoinVertical(lipgloss.Left,
			styles.ContentTitleH2Style.Width(boxWidthForTwoH2-2).Render("Stub Networks"),
			lipgloss.NewStyle().Align(lipgloss.Center).Margin(0, 2).Width(boxWidthForTwoH2).Render(stubTable.String()),
			styles.ContentBottomBorderStyle.Width(boxWidthForTwoH2-2).Render(""),
		)

		horizontalTables := lipgloss.JoinHorizontal(lipgloss.Top, correctBoxWidthTransit, correctBoxWidthStub)

		completeAreaRouterLSAs := lipgloss.JoinVertical(lipgloss.Left, areaHeader, horizontalTables)

		routerLSABlocks = append(routerLSABlocks, completeAreaRouterLSAs+"\n\n")
	}

	contentMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight
	m.viewport.Width = boxWidthForOneH1 + 2
	m.viewport.Height = contentMaxHeight

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, routerLSABlocks...))

	return m.viewport.View()
}

func (m *Model) renderExternalMonitorTab() string {
	boxWidthForTwoH2 := (m.windowSize.Width - 16) / 2 // - 6 (padding+margin content) - 2 (for border) - 8 (for margin)
	if boxWidthForTwoH2 < 20 {
		boxWidthForTwoH2 = 20 // Minimum width to ensure readability
	}
	boxWidthForOneH1 := m.windowSize.Width - 10 // - 6 (padding+margin content) - 2 (for each border)
	if boxWidthForOneH1 < 20 {
		boxWidthForOneH1 = 20 // Minimum width to ensure readability
	}
	boxWidthForOneH2 := boxWidthForOneH1 - 4 // -4 (for margin) -2 (for border)

	var externalLsaBlock []string
	var nssaExternalLsaBlock []string

	externalLSASelf, _ := getOspfExternalData()
	nssaExternalDataSelf, _ := getOspfNssaExternalData()

	var externalTableData [][]string
	var externalTableDataExpanded [][]string
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

	externalHeader := styles.ContentTitleH1Style.
		Width(boxWidthForOneH1).
		Margin(0, 0, 1, 0).
		Padding(1, 0, 0, 0).
		Render("External LSAs (Type 5)")

	externalDataBox := lipgloss.JoinVertical(lipgloss.Left,
		styles.ContentTitleH2Style.Width(boxWidthForOneH2).Render("Self Originating"),
		lipgloss.NewStyle().Align(lipgloss.Center).Margin(0, 2).Width(boxWidthForOneH2).Render(externalTable.String()),
		styles.ContentBottomBorderStyle.Width(boxWidthForOneH2).Render(""),
	)

	var completeExternalBox string
	completeExternalBox = lipgloss.JoinVertical(lipgloss.Left, externalHeader, externalDataBox)

	externalLsaBlock = append(externalLsaBlock, completeExternalBox+"\n\n")

	for area, areaData := range nssaExternalDataSelf.NssaExternalLinkStates {
		var nssaExternalTableData [][]string
		// var nssaExternalTableDataExpanded [][]string

		for _, lsaData := range areaData.Data {
			nssaExternalTableData = append(nssaExternalTableData, []string{
				lsaData.LinkStateId,
				"/" + strconv.Itoa(int(lsaData.NetworkMask)),
				lsaData.MetricType,
				lsaData.NssaForwardAddress,
			})
		}
		rowsNssaExternal := len(nssaExternalTableData)
		nssaExternalTable := ltable.New().
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
				case row == rowsNssaExternal-1:
					return styles.NormalCellStyle.BorderBottom(true)
				default:
					return styles.NormalCellStyle
				}
			})

		for _, r := range nssaExternalTableData {
			nssaExternalTable = nssaExternalTable.Row(r...)
		}

		nssaExternalHeader := styles.ContentTitleH1Style.
			Width(boxWidthForOneH1).
			Margin(0, 0, 1, 0).
			Padding(1, 0, 0, 0).
			Render("NSSA External LSAs (Type 7) in Area " + area)

		nssaExternalDataBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.ContentTitleH2Style.Width(boxWidthForOneH2).Render("Self Originating"),
			lipgloss.NewStyle().Align(lipgloss.Center).Margin(0, 2).Width(boxWidthForOneH2).Render(nssaExternalTable.String()),
			styles.ContentBottomBorderStyle.Width(boxWidthForOneH2).Render(""),
		)
		// var completeNssaExternalBox string
		completeNssaExternalBox := lipgloss.JoinVertical(lipgloss.Left, nssaExternalHeader, nssaExternalDataBox)

		nssaExternalLsaBlock = append(nssaExternalLsaBlock, completeNssaExternalBox+"\n\n")
	}

	contentMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight
	m.viewport.Width = boxWidthForOneH1 + 2
	m.viewport.Height = contentMaxHeight

	var allLsaBlocks []string
	if externalLSASelf.AsExternalLinkStates != nil {
		if len(nssaExternalDataSelf.NssaExternalLinkStates) == 0 {
			allLsaBlocks = append(externalLsaBlock, nssaExternalLsaBlock...)
		} else {
			allLsaBlocks = externalLsaBlock
		}
	} else {
		allLsaBlocks = allLsaBlocks[:0]
		allLsaBlocks = append(allLsaBlocks, lipgloss.JoinVertical(lipgloss.Left,
			externalHeader,
			"no self originating external advertisements",
		))
	}

	m.viewport.SetContent(lipgloss.JoinVertical(lipgloss.Left, allLsaBlocks...))

	return m.viewport.View()
}

func (m *Model) renderAdvertisementTab() string {
	// on this view:
	// show three advertisement boxes: based on vtysh LSA queries / based on file analysis / based on FIB analysis

	boxWidthForTwo := (m.windowSize.Width - 10) / 2 // - 6 (padding+margin content) - 2 (for gap) - 2 (for border)
	if boxWidthForTwo < 20 {
		boxWidthForTwo = 20 // Minimum width to ensure readability
	}

	shouldAdvertisedTitle := styles.BoxTitleStyle.Render("Should be Advertised")

	shouldAdvertisedContent := "Area 0.0.0.0: \n"

	shouldAdvertisedRouterLSA := styles.ContentTitleH1Style.
		Width(boxWidthForTwo - 2).
		Render("Area 0.0.0.0, Router LSAs (Type 1)")
	shouldAdvertisedContent += shouldAdvertisedRouterLSA

	shouldAdvertisedVerticalStyle := lipgloss.JoinVertical(lipgloss.Left, shouldAdvertisedTitle, shouldAdvertisedContent)
	shouldAdvertisedBox := lipgloss.NewStyle().Render(shouldAdvertisedVerticalStyle)

	// ----------------------------------------------------
	gap := 2
	// ----------------------------------------------------

	isAdvertisedTitle := styles.BoxTitleStyle.Render("Is Advertised")

	isAdvertisedContent := "Area 0.0.0.0: \n"

	// for each area create area box --> need length of areas

	isAdvertisedRouterLSA := styles.ContentTitleH1Style.
		Width(boxWidthForTwo - 2).
		Render("Area 0.0.0.0, Router LSAs (Type 1)")
	isAdvertisedContent += isAdvertisedRouterLSA

	isAdvertisedVerticalStyle := lipgloss.JoinVertical(lipgloss.Left, isAdvertisedTitle, isAdvertisedContent)
	isAdvertisedBox := lipgloss.NewStyle().Render(isAdvertisedVerticalStyle)
	// returnString := "Advertisement"

	horizontalBoxes := lipgloss.JoinHorizontal(lipgloss.Top,
		isAdvertisedBox,
		lipgloss.NewStyle().Width(gap).Render(""),
		shouldAdvertisedBox,
	)

	return horizontalBoxes
}

func (m *Model) renderOSPFTab0() string {
	// Calculate box width dynamically for four horizontal boxes based on terminal width
	boxWidthForFour := (m.windowSize.Width - 16) / 4 // - 6 (padding+margin content) - 10 (for each border)
	if boxWidthForFour < 20 {
		boxWidthForFour = 20 // Minimum width to ensure readability
	}

	ospfAnomalyOne := styles.GeneralBoxStyle.
		Width(boxWidthForFour).
		Render(styles.BoxTitleStyle.Render("OSPF Anomaly One") + "\n" + "Call Backend...☎\nEverything Good! amount")

	ospfAnomalyTwo := styles.GeneralBoxStyle.
		Width(boxWidthForFour).
		Render(styles.BoxTitleStyle.Render("OSPF Anomaly Two") + "\n" + "Call Backend...☎\nEverything Good!")

	ospfAnomalyThree := styles.BadBoxStyle.
		Width(boxWidthForFour).
		Render(styles.BoxTitleStyle.Render("OSPF Anomaly Three") + "\n" + "Call Backend...☎\nVery Bad Anomaly Detected!\n\nReport...\nReport...\nReport...\nReport...\nReport...\n")

	ospfAnomalyFour := styles.GeneralBoxStyle.
		Width(boxWidthForFour).
		Render(styles.BoxTitleStyle.Render("OSPF Anomaly Four") + "\n" + "Call Backend...☎\nEverything Good!")

	ospfAnomalies := []struct {
		Title   string
		Content string
		Style   lipgloss.Style
	}{
		{
			Title:   "OSPF Anomaly One",
			Content: "Call Backend...☎\nVery Bad Anomaly Detected!\n\nReport...\nReport...\nReport...\nReport...\nReport...\n",
			Style:   styles.BadBoxStyle,
		},
		{
			Title:   "OSPF Anomaly Two",
			Content: "Call Backend...☎\nEverything Good!",
			Style:   styles.GeneralBoxStyle,
		},
		{
			Title:   "OSPF Anomaly Three",
			Content: "Call Backend...☎\nEverything Good!",
			Style:   styles.GeneralBoxStyle,
		},
		{
			Title:   "OSPF Anomaly Four",
			Content: "Call Backend...☎\nEverything Good!",
			Style:   styles.GeneralBoxStyle,
		},
	}

	// Build anomaly boxes using the new component
	var ospfAnomalyBoxes []string
	for _, a := range ospfAnomalies {
		box := components.NewAnomalyBox(a.Title, a.Content, a.Style, boxWidthForFour)
		ospfAnomalyBoxes = append(ospfAnomalyBoxes, box.Render())
	}

	horizontalBoxes := lipgloss.JoinHorizontal(lipgloss.Top, ospfAnomalyOne, ospfAnomalyTwo, ospfAnomalyThree, ospfAnomalyFour)
	horizontalBoxes2 := lipgloss.JoinHorizontal(lipgloss.Top, ospfAnomalyBoxes...)

	//infoBox := styles.InfoTextStyle.
	//	Width(m.windowSize.Width - 12).
	//	Render("press 'r' to refresh ospf anomalies")
	//
	//return lipgloss.JoinVertical(lipgloss.Left, horizontalBoxes, infoBox)

	return lipgloss.JoinVertical(lipgloss.Left, horizontalBoxes, horizontalBoxes2)
}

func (m *Model) renderOSPFTab1() string {
	// Calculate box width dynamically for four horizontal boxes based on terminal width
	boxWidthForFour := (m.windowSize.Width - 16) / 4 // - 6 (padding+margin content) - 10 (for each border)
	if boxWidthForFour < 20 {
		boxWidthForFour = 20 // Minimum width to ensure readability
	}

	ospfAnomalyOne := styles.GeneralBoxStyle.
		Width(boxWidthForFour).
		Render(styles.BoxTitleStyle.Render("OSPF Anomaly One") + "\n" + "Call Backend...☎\nEverything Good! amount")

	ospfAnomalyTwo := styles.GeneralBoxStyle.
		Width(boxWidthForFour).
		Render(styles.BoxTitleStyle.Render("OSPF Anomaly Two") + "\n" + "Call Backend...☎\nEverything Good!")

	ospfAnomalyThree := styles.BadBoxStyle.
		Width(boxWidthForFour).
		Render(styles.BoxTitleStyle.Render("OSPF Anomaly Three") + "\n" + "Call Backend...☎\nVery Bad Anomaly Detected!\n\nReport...\nReport...\nReport...\nReport...\nReport...\n")

	ospfAnomalyFour := styles.GeneralBoxStyle.
		Width(boxWidthForFour).
		Render(styles.BoxTitleStyle.Render("OSPF Anomaly Four") + "\n" + "Call Backend...☎\nEverything Good!")

	return lipgloss.JoinHorizontal(lipgloss.Top, ospfAnomalyThree, ospfAnomalyOne, ospfAnomalyTwo, ospfAnomalyFour)
}

func (m *Model) renderRunningConfigTab() string {
	// Calculate box width dynamically for one horizontal box based on terminal width
	boxWidthForOne := m.windowSize.Width - 10 // - 6 (padding+margin content) - 2 (for each border)
	if boxWidthForOne < 20 {
		boxWidthForOne = 20 // Minimum width to ensure readability
	}

	outputMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight
	m.viewport.Width = boxWidthForOne
	m.viewport.Height = outputMaxHeight

	m.viewport.SetContent(strings.Join(m.runningConfig, "\n"))

	runningConfigBox := lipgloss.NewStyle().Padding(0, 5).Render(m.viewport.View())

	return runningConfigBox
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
