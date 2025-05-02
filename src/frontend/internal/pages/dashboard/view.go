package dashboard

import (
	"fmt"
	// "github.com/ba2025-ysmprc/frr-tui/pkg"
	backend "github.com/ba2025-ysmprc/frr-tui/internal/services"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/components"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
	"strconv"
)

var (
	currentSubTabLocal = -1
	hasAnomalyDetected = false
)

// DashboardView is the updated View function. This allows to call View with an argument.
func (m *Model) DashboardView(currentSubTab int) string {
	currentSubTabLocal = currentSubTab
	return m.View()
}

func (m *Model) View() string {
	if currentSubTabLocal == 0 && !hasAnomalyDetected {
		return m.renderOSPFDashboard()
	} else if currentSubTabLocal == 1 {
		return ""
	}
	return m.renderOSPFDashboard()
}

func (m *Model) renderOSPFDashboard() string {
	ospfDashboardLsdbSelf := getOspfDashboardLsdbSelf()

	// Update the viewport
	m.viewport.Width = styles.WidthTwoH1ThreeFourth + 2
	m.viewport.Height = m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight - 2
	m.viewport.SetContent(ospfDashboardLsdbSelf)

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

func getOspfDashboardLsdbSelf() string {
	var lsdbSelfBlocks []string

	lsdb, _ := backend.GetLSDB()
	_, routerOSPFID, _ := backend.GetRouterName()

	// ===== OSPF Internal LSAs (Type 1-4) =====
	for area, lsaTypes := range lsdb.Areas {
		var routerLinkStateTableData [][]string
		var networkLinkStateTableData [][]string
		var summaryLinkStateTableData [][]string
		var asbrSummaryLinkStateTableData [][]string

		//var amountOfRouterLS string
		//var amountOfNetworkLS string
		//var amountOfSummaryLS string
		//var amountOfAsSummaryLS string

		// loop through LSAs (type 1-4) and extract self-originating data for tables
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

		dashboardHeader := styles.H1TitleStyle().
			Width(styles.WidthTwoH1ThreeFourth).
			BorderBottom(true).
			Render("All OSPF Routes are advertised as Expected")
		areaHeader := styles.H1TitleStyle().Width(styles.WidthTwoH1ThreeFourth).
			Render(fmt.Sprintf("Link State Database (Self): Area %s", area))

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

		var completeType2To4LsdbSelfList []string
		if networkLinkStateTableData != nil {
			completeType2To4LsdbSelfList = append(completeType2To4LsdbSelfList, networkTableBox)
		}
		if summaryLinkStateTableData != nil {
			completeType2To4LsdbSelfList = append(completeType2To4LsdbSelfList, summaryTableBox)
		}
		if asbrSummaryLinkStateTableData != nil {
			completeType2To4LsdbSelfList = append(completeType2To4LsdbSelfList, asbrSummaryTableBox)
		}

		completeType2To4LsdbSelf := lipgloss.JoinVertical(lipgloss.Left,
			completeType2To4LsdbSelfList...,
		)

		completeAreaLSDBSelf := lipgloss.JoinVertical(lipgloss.Left,
			dashboardHeader,
			areaHeader,
			routerTableBox,
			completeType2To4LsdbSelf,
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
		Render("Link State Database: AS External LSAs")

	// create styled boxes for each external LSA Type (type 5 & 7)
	externalTableBox := lipgloss.JoinVertical(lipgloss.Left,
		styles.H2TitleStyle().Width(styles.WidthTwoH2ThreeFourth).
			Render("Self-Originating AS External Link States"),
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

// ============================== //
// HELPERS: BACKEND CALLS         //
// ============================== //

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

//func getLSDB() (*pkg.OSPFDatabase, error) {
//	response, err := backend.SendMessage("ospf", "database", nil)
//	if err != nil {
//		return nil, err
//	}
//
//	return response.Data.GetOspfDatabase(), nil
//}
