package dashboard

import (
	"fmt"
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"sort"
	"strings"
	"time"

	// "github.com/ba2025-ysmprc/frr-tui/pkg"
	backend "github.com/ba2025-ysmprc/frr-tui/internal/services"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/components"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	frrProto "github.com/ba2025-ysmprc/frr-tui/pkg"
	"github.com/charmbracelet/lipgloss"
	"strconv"
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
		m.detectAnomaly()
		return m.renderOSPFDashboard()
	} else if currentSubTabLocal == 1 {
		return "TBD"
	}
	return m.renderOSPFDashboard()
}

func (m *Model) renderOSPFDashboard() string {
	// Update the viewport
	m.viewport.Width = styles.WidthTwoH1ThreeFourth + 2
	m.viewport.Height = m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight - 2

	if m.hasAnomalyDetected {
		ospfDashboardAnomalies := getOspfDashboardAnomalies()
		m.viewport.SetContent(ospfDashboardAnomalies)
	} else {
		ospfDashboardLsdbSelf := getOspfDashboardLsdbSelf()
		m.viewport.SetContent(ospfDashboardLsdbSelf)
	}

	//cpuAmount, cpuUsage, memoryUsage, err := backend.GetSystemResources()
	//var cpuAmountString, cpuUsageString, memoryString string
	//if err != nil {
	//	cpuAmountString = "N/A"
	//	cpuUsageString = "N/A"
	//	memoryString = "N/A"
	//} else {
	//	cpuAmountString = fmt.Sprintf("%v", cpuAmount)
	//	cpuUsageString = fmt.Sprintf("%.2f%%", cpuUsage*100)
	//	memoryString = fmt.Sprintf("%.2f%%", memoryUsage)
	//}
	//
	//cpuStatistics := lipgloss.JoinVertical(lipgloss.Left,
	//	styles.H2TitleStyle().Width(styles.WidthTwoH2OneFourth).Render("CPU Metrics"),
	//	styles.H2TwoContentBoxStyleP1101().Width(styles.WidthTwoH2OneFourthBox).Render(
	//		"CPU Usage: "+cpuUsageString+"\n"+
	//			"Cores: "+cpuAmountString),
	//	styles.H2BoxBottomBorderStyle().Width(styles.WidthTwoH2OneFourth).Render(""),
	//)
	//
	//memoryStatistics := lipgloss.JoinVertical(lipgloss.Left,
	//	styles.H2TitleStyle().Width(styles.WidthTwoH2OneFourth).Render("Memory Metrics"),
	//	styles.H2TwoContentBoxStyleP1101().Width(styles.WidthTwoH2OneFourthBox).Render(
	//		"Memory Usage: "+memoryString),
	//	styles.H2BoxBottomBorderStyle().Width(styles.WidthTwoH2OneFourth).Render(""),
	//)
	//
	//systemResources := lipgloss.JoinVertical(lipgloss.Left,
	//	styles.H1TitleStyle().Width(styles.WidthTwoH1OneFourth).Render("System Resources"),
	//	cpuStatistics,
	//	memoryStatistics,
	//)

	dashboardRight := lipgloss.JoinVertical(lipgloss.Left, getSystemResourcesBox(), getOSPFGeneralInfoBox())

	horizontalDashboard := lipgloss.JoinHorizontal(lipgloss.Top,
		m.viewport.View(),
		dashboardRight,
	)

	return horizontalDashboard
}

func getSystemResourcesBox() string {
	cpuAmount, cpuUsage, memoryUsage, err := backend.GetSystemResources()
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

	systemStatistics := lipgloss.JoinVertical(lipgloss.Left,
		styles.H2TitleStyle().Width(styles.WidthTwoH2OneFourth).Render("CPU Metrics"),
		styles.H2TwoContentBoxStyleP1101().Width(styles.WidthTwoH2OneFourthBox).Render(
			"CPU Usage: "+cpuUsageString+"\n"+
				"Cores: "+cpuAmountString+"\n"+
				"Memory Usage: "+memoryString)+"\n",
	)

	systemResources := lipgloss.JoinVertical(lipgloss.Left,
		styles.H1TitleStyle().Width(styles.WidthTwoH1OneFourth).Render("System Resources"),
		systemStatistics,
	)

	return systemResources
}

func getOSPFGeneralInfoBox() string {
	ospfInformation, err := backend.GetOSPF()
	if err != nil {
		common.PrintBackendError(err, "GetOSPF")
	}

	lastSPFExecution := time.Duration(ospfInformation.SpfLastExecutedMsecs) * time.Millisecond
	lastSPFExecution = lastSPFExecution.Truncate(time.Second) // remove sub-second precision

	ospfRouterInfo := styles.H1TwoContentBoxesStyle().Width(styles.WidthTwoH1OneFourthBox).Render(
		"OSPF Router ID: " + ospfInformation.RouterId + "\n" +
			"Last SPF Execution: " + lastSPFExecution.String() + "\n" +
			"Total External LSAs: " + strconv.Itoa(int(ospfInformation.LsaExternalCounter)) + "\n" +
			"Attached Areas: " + strconv.Itoa(int(ospfInformation.AttachedAreaCounter)) + "\n")

	var ospfAreaInformation []string
	for areaID, areaData := range ospfInformation.Areas {
		ospfAreaInformation = append(ospfAreaInformation,
			styles.H2TitleStyle().Width(styles.WidthTwoH2OneFourth).Render("Area "+areaID))
		ospfAreaInformation = append(ospfAreaInformation,
			styles.H2TwoContentBoxesStyle().Width(styles.WidthTwoH2OneFourthBox).Render(
				"Full Adjencencies: "+strconv.Itoa(int(areaData.NbrFullAdjacentCounter))+"\n"+
					"Total LSAs: "+strconv.Itoa(int(areaData.LsaNumber))+"\n"))
	}

	renderedOSPFAreaInformation := lipgloss.JoinVertical(lipgloss.Left, ospfAreaInformation...)

	ospfInformationBox := lipgloss.JoinVertical(lipgloss.Left,
		styles.H1TitleStyle().Width(styles.WidthTwoH1OneFourth).Render("General OSPF Information"),
		ospfRouterInfo,
		renderedOSPFAreaInformation,
	)

	return ospfInformationBox
}

func getOspfDashboardLsdbSelf() string {
	var lsdbSelfBlocks []string

	dashboardHeader := styles.H1GoodTitleStyle().
		Width(styles.WidthTwoH1ThreeFourth).
		BorderBottom(true).
		Padding(0).
		Render("All OSPF Routes are advertised as Expected")

	lsdbSelfBlocks = append(lsdbSelfBlocks, dashboardHeader)

	lsdb, err := backend.GetLSDB()
	if err != nil {
		return common.PrintBackendError(err, "GetLSDB")
	}

	// extract and sort the map keys
	lsdbAreas := make([]string, 0, len(lsdb.Areas))
	for area := range lsdb.Areas {
		lsdbAreas = append(lsdbAreas, area)
	}
	sort.Strings(lsdbAreas)

	_, routerOSPFID, err := backend.GetRouterName()
	if err != nil {
		return common.PrintBackendError(err, "GetRouterName")
	}

	// ===== OSPF Internal LSAs (Type 1-4) =====
	for _, areaID := range lsdbAreas {
		lsaTypes := lsdb.Areas[areaID]

		var routerLinkStateTableData [][]string
		var networkLinkStateTableData [][]string
		var summaryLinkStateTableData [][]string
		var asbrSummaryLinkStateTableData [][]string
		var nssaExternalLinkStateTableData [][]string

		//var amountOfRouterLS string
		//var amountOfNetworkLS string
		//var amountOfSummaryLS string
		//var amountOfAsSummaryLS string

		// loop through LSAs (type 1-4 + Type 7) and extract self-originating data for tables
		for _, routerLinkState := range lsaTypes.RouterLinkStates {
			if routerLinkState.Base.AdvertisedRouter == routerOSPFID {
				routerLinkStateTableData = append(routerLinkStateTableData, []string{
					routerLinkState.Base.AdvertisedRouter,
					strconv.Itoa(int(routerLinkState.NumOfRouterLinks)),
					strconv.Itoa(int(routerLinkState.Base.LsaAge)),
				})
			}
		}
		//if routerLinkStateTableData != nil {
		//	amountOfRouterLS = strconv.Itoa(int(lsaTypes.RouterLinkStatesCount))
		//} else {
		//	amountOfRouterLS = "0"
		//}
		for _, networkLinkState := range lsaTypes.NetworkLinkStates {
			if networkLinkState.Base.AdvertisedRouter == routerOSPFID {
				networkLinkStateTableData = append(networkLinkStateTableData, []string{
					networkLinkState.Base.LsId,
					networkLinkState.Base.AdvertisedRouter,
					strconv.Itoa(int(networkLinkState.Base.LsaAge)),
				})
			}
		}
		//if networkLinkStateTableData != nil {
		//	amountOfNetworkLS = strconv.Itoa(int(lsaTypes.NetworkLinkStatesCount))
		//} else {
		//	amountOfNetworkLS = "0"
		//}

		for _, summaryLinkState := range lsaTypes.SummaryLinkStates {
			if summaryLinkState.Base.AdvertisedRouter == routerOSPFID {
				summaryLinkStateTableData = append(summaryLinkStateTableData, []string{
					summaryLinkState.SummaryAddress,
					summaryLinkState.Base.AdvertisedRouter,
					strconv.Itoa(int(summaryLinkState.Base.LsaAge)),
				})
			}
		}
		//if summaryLinkStateTableData == nil {
		//	amountOfSummaryLS = "0"
		//} else {
		//	amountOfSummaryLS = strconv.Itoa(int(lsaTypes.SummaryLinkStatesCount))
		//}

		for _, asbrSummaryLinkState := range lsaTypes.AsbrSummaryLinkStates {
			if asbrSummaryLinkState.Base.AdvertisedRouter == routerOSPFID {
				asbrSummaryLinkStateTableData = append(asbrSummaryLinkStateTableData, []string{
					asbrSummaryLinkState.Base.LsId,
					asbrSummaryLinkState.Base.AdvertisedRouter,
					strconv.Itoa(int(asbrSummaryLinkState.Base.LsaAge)),
				})
			}
		}
		//if asbrSummaryLinkStateTableData == nil {
		//	amountOfAsSummaryLS = "0"
		//} else {
		//	amountOfAsSummaryLS = strconv.Itoa(int(lsaTypes.AsbrSummaryLinkStatesCount))
		//}

		for _, nssaExternalLinkStates := range lsaTypes.NssaExternalLinkStates {
			nssaExternalLinkStateTableData = append(nssaExternalLinkStateTableData, []string{
				nssaExternalLinkStates.Route,
				nssaExternalLinkStates.MetricType,
				nssaExternalLinkStates.Base.AdvertisedRouter,
				strconv.Itoa(int(nssaExternalLinkStates.Base.LsaAge)),
			})
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
		sort.Slice(nssaExternalLinkStateTableData, func(i, j int) bool {
			return nssaExternalLinkStateTableData[i][0] < nssaExternalLinkStateTableData[j][0]
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

		areaHeader := styles.H1TitleStyle().Width(styles.WidthTwoH1ThreeFourth).
			Render(fmt.Sprintf("Link State Database (Self): Area %s", areaID))

		// create styled boxes for each LSA Type (type 1-4)
		routerTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyle().Width(styles.WidthTwoH2ThreeFourth).
				Render("Self-Originating Router Link States"),
			styles.H2ContentBoxCenterStyle().Width(styles.WidthTwoH2ThreeFourthBox).
				Render(routerLinkStateTable.String()),
			styles.H2BoxBottomBorderStyle().Width(styles.WidthTwoH2ThreeFourth).Render(""),
		)
		networkTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyle().Width(styles.WidthTwoH2ThreeFourth).
				Render("Self-Originating Network Link States"),
			styles.H2ContentBoxCenterStyle().Width(styles.WidthTwoH2ThreeFourthBox).
				Render(networkLinkStateTable.String()),
			styles.H2BoxBottomBorderStyle().Width(styles.WidthTwoH2ThreeFourth).Render(""),
		)
		summaryTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyle().Width(styles.WidthTwoH2ThreeFourth).
				Render("Self-Originating Summary Link States"),
			styles.H2ContentBoxCenterStyle().Width(styles.WidthTwoH2ThreeFourthBox).
				Render(summaryLinkStateTable.String()),
			styles.H2BoxBottomBorderStyle().Width(styles.WidthTwoH2ThreeFourth).Render(""),
		)
		asbrSummaryTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyle().Width(styles.WidthTwoH2ThreeFourth).
				Render("Self-Originating ASBR Summary Link States"),
			styles.H2ContentBoxCenterStyle().Width(styles.WidthTwoH2ThreeFourthBox).
				Render(asbrSummaryLinkStateTable.String()),
			styles.H2BoxBottomBorderStyle().Width(styles.WidthTwoH2ThreeFourth).Render(""),
		)
		nssaExternalTableBox := lipgloss.JoinVertical(lipgloss.Left,
			styles.H2TitleStyle().Width(styles.WidthTwoH2ThreeFourth).
				Render("Self-Originating NSSA External Link States (Type 7)"),
			styles.H2ContentBoxCenterStyle().Width(styles.WidthTwoH2ThreeFourthBox).
				Render(nssaExternalLinkStateTable.String()),
			styles.H2BoxBottomBorderStyle().Width(styles.WidthTwoH2ThreeFourth).Render(""),
		)

		var optionalLSATypesList []string
		if networkLinkStateTableData != nil {
			optionalLSATypesList = append(optionalLSATypesList, networkTableBox)
		}
		if summaryLinkStateTableData != nil {
			optionalLSATypesList = append(optionalLSATypesList, summaryTableBox)
		}
		if asbrSummaryLinkStateTableData != nil {
			optionalLSATypesList = append(optionalLSATypesList, asbrSummaryTableBox)
		}
		if nssaExternalLinkStateTableData != nil {
			optionalLSATypesList = append(optionalLSATypesList, nssaExternalTableBox)
		}

		activeOptionalLSATypes := lipgloss.JoinVertical(lipgloss.Left,
			optionalLSATypesList...,
		)

		completeAreaLSDBSelf := lipgloss.JoinVertical(lipgloss.Left,
			areaHeader,
			routerTableBox,
			activeOptionalLSATypes,
		)

		lsdbSelfBlocks = append(lsdbSelfBlocks, completeAreaLSDBSelf+"\n\n")
	}

	// ===== External LSA =====
	var asExternalLinkStateTableData [][]string
	//var amountOfExternalLS string
	for _, asExternalLinkState := range lsdb.AsExternalLinkStates {
		if asExternalLinkState.Base.AdvertisedRouter == routerOSPFID {
			asExternalLinkStateTableData = append(asExternalLinkStateTableData, []string{
				asExternalLinkState.Route,
				asExternalLinkState.MetricType,
				asExternalLinkState.Base.AdvertisedRouter,
				strconv.Itoa(int(asExternalLinkState.Base.LsaAge)),
			})
		}
	}
	//if asExternalLinkStateTableData == nil {
	//	amountOfExternalLS = "0"
	//} else {
	//	amountOfExternalLS = strconv.Itoa(int(lsdb.AsExternalCount))
	//}

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

	externalHeader := styles.H1TitleStyle().Width(styles.WidthTwoH1ThreeFourth).
		Render("Link State Database (Self): AS External")

	// create styled boxes for external LSA Type (type 5)
	externalTableBox := lipgloss.JoinVertical(lipgloss.Left,
		styles.H2TitleStyle().Width(styles.WidthTwoH2ThreeFourth).
			Render("Self-Originating AS External Link States (Type 5)"),
		styles.H2ContentBoxCenterStyle().Width(styles.WidthTwoH2ThreeFourthBox).
			Render(asExternalLinkStateTable.String()),
		styles.H2BoxBottomBorderStyle().Width(styles.WidthTwoH2ThreeFourth).Render(""),
	)

	completeExternalLSDB := lipgloss.JoinVertical(lipgloss.Left,
		externalHeader,
		externalTableBox,
	)

	if asExternalLinkStateTableData != nil {
		lsdbSelfBlocks = append(lsdbSelfBlocks, completeExternalLSDB+"\n\n")
	}

	return lipgloss.JoinVertical(lipgloss.Left, lsdbSelfBlocks...)
}

func getOspfDashboardAnomalies() string {
	ospfRouterAnomalies, err := backend.GetRouterAnomalies()
	if err != nil {
		return common.PrintBackendError(err, "GetRouterAnomalies")
	}
	ospfExternalAnomalies, err := backend.GetExternalAnomalies()
	if err != nil {
		return common.PrintBackendError(err, "GetExternalAnomalies")
	}
	ospfNSSAExternalAnomalies, err := backend.GetNSSAExternalAnomalies()
	if err != nil {
		return common.PrintBackendError(err, "GetNSSAExternalAnomalies")
	}

	var routerAnomalyTable string
	if common.HasAnyAnomaly(ospfRouterAnomalies) {
		routerAnomalyTable = createAnomalyTable(
			ospfRouterAnomalies,
			"Router Anomalies (Type 1 LSAs)",
		)
	}

	var externalAnomalyTable string
	if common.HasAnyAnomaly(ospfExternalAnomalies) {
		externalAnomalyTable = createAnomalyTable(
			ospfExternalAnomalies,
			"External Link State Anomalies (Type 5 LSAs)",
		)
	}

	var nssaExternalAnomalyTable string
	if common.HasAnyAnomaly(ospfNSSAExternalAnomalies) {
		nssaExternalAnomalyTable = createAnomalyTable(
			ospfNSSAExternalAnomalies,
			"NSSA External Link State Anomalies (Type 7 LSAs)",
		)
	}

	// prevents printing empty strings
	var allAnomaliesList []string
	if routerAnomalyTable != "" {
		allAnomaliesList = append(allAnomaliesList, routerAnomalyTable)
	}
	if externalAnomalyTable != "" {
		allAnomaliesList = append(allAnomaliesList, externalAnomalyTable)
	}
	if nssaExternalAnomalyTable != "" {
		allAnomaliesList = append(allAnomaliesList, nssaExternalAnomalyTable)
	}

	allAnomalies := lipgloss.JoinVertical(lipgloss.Left, allAnomaliesList...)

	return allAnomalies
}

func createAnomalyTable(a *frrProto.AnomalyDetection, lsaTypeHeader string) string {
	// extract data for tables
	var tableData [][]string

	// TODO: add all anomily types
	if a.HasOverAdvertisedPrefixes {
		for _, superfluousEntry := range a.SuperfluousEntries {
			var firstCol string
			if strings.Contains(lsaTypeHeader, "Router") {
				firstCol = superfluousEntry.InterfaceAddress
			} else {
				firstCol = superfluousEntry.LinkStateId
			}

			tableData = append(tableData, []string{
				firstCol,
				"/" + superfluousEntry.PrefixLength,
				superfluousEntry.LinkType,
				"Overadvertised Route",
			})
		}
	}

	if a.HasUnderAdvertisedPrefixes {
		for _, missingEntry := range a.MissingEntries {
			var firstCol string
			if strings.Contains(lsaTypeHeader, "Router") {
				firstCol = missingEntry.InterfaceAddress
			} else {
				firstCol = missingEntry.LinkStateId
			}

			tableData = append(tableData, []string{
				firstCol,
				"/" + missingEntry.PrefixLength,
				missingEntry.LinkType,
				"Underadvertised Route",
			})
		}
	}

	// Order all Table Data
	sort.Slice(tableData, func(i, j int) bool {
		return tableData[i][0] < tableData[j][0]
	})

	// create the tables and fill it with collected data
	rows := len(tableData)
	table := components.NewAnomalyTable(
		[]string{
			"Network Address",
			"CIDR",
			"Link Type",
			"Anomaly Type",
		},
		rows,
	)
	for _, r := range tableData {
		table = table.Row(r...)
	}

	// style the output
	// anomalyHeader := styles.H1BadTitleStyle().Width(styles.WidthTwoH1ThreeFourth).Render("Router (Type 1) Anomalies")

	tableBox := lipgloss.JoinVertical(lipgloss.Left,
		styles.H1BadTitleStyle().Width(styles.WidthTwoH1ThreeFourth).Render(lsaTypeHeader),
		styles.H1ContentBoxCenterStyle().Width(styles.WidthTwoH1ThreeFourthBox).Render(table.String()),
		styles.H1BadBoxBottomBorderStyle().Width(styles.WidthTwoH1ThreeFourth).Render(""),
	)

	return tableBox
}

// ============================== //
// HELPERS: BACKEND CALLS         //
// ============================== //
