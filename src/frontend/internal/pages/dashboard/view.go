package dashboard

import (
	"fmt"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/toast"
	"sort"
	"strings"
	"time"

	"github.com/ba2025-ysmprc/frr-mad/src/logger"
	"github.com/ba2025-ysmprc/frr-tui/internal/common"

	// "github.com/ba2025-ysmprc/frr-tui/pkg"
	"strconv"

	backend "github.com/ba2025-ysmprc/frr-tui/internal/services"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/components"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	frrProto "github.com/ba2025-ysmprc/frr-tui/pkg"
	"github.com/charmbracelet/lipgloss"
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
	var body string
	switch currentSubTabLocal {
	case 0:
		if m.showAnomalyOverlay {
			body = m.renderAnomalyDetails()
		} else if m.showExportOverlay {
			body = components.RenderExportOptions(
				m.exportOptions,
				m.exportData,
				&m.cursor,
				&m.viewportRightHalf,
			)
		} else {
			m.detectAnomaly()
			body = m.renderOSPFDashboard()
		}
	case 1:
		body = "TBD"
	default:
		body = m.renderOSPFDashboard()
	}

	toastView := m.toast.View()
	if toastView == "" {
		return body
	}

	totalW := styles.WidthBasis
	totalH := styles.HeightBasis
	x := 2
	y := 0

	return toast.Overlay(body, toastView, x, y, totalW, totalH)
}

func (m *Model) renderOSPFDashboard() string {
	// Update the viewportLeft
	m.viewportLeft.Width = styles.ViewPortWidthThreeFourth
	m.viewportLeft.Height = styles.ViewPortHeightCompletePage - styles.HeightH1

	m.viewportRight.Width = styles.ViewPortWidthOneFourth
	m.viewportRight.Height = styles.ViewPortHeightCompletePage - styles.HeightH1

	var statusHeader string
	if m.hasAnomalyDetected {
		anomalyHeader := styles.H1BadTitleStyle().
			Width(styles.WidthTwoH1ThreeFourth).
			BorderBottom(true).
			Padding(0).
			Render("Anomalies Detected!")
		statusHeader = anomalyHeader
		ospfDashboardAnomalies := m.getOspfDashboardAnomalies(m.logger)
		m.viewportLeft.SetContent(ospfDashboardAnomalies)
	} else {
		dashboardHeader := styles.H1GoodTitleStyle().
			Width(styles.WidthTwoH1ThreeFourth).
			BorderBottom(true).
			Padding(0).
			Render("All OSPF Routes are advertised as Expected")
		statusHeader = dashboardHeader
		ospfDashboardLsdbSelf := m.getOspfDashboardLsdbSelf(m.logger)
		m.viewportLeft.SetContent(ospfDashboardLsdbSelf)
	}

	systemResourceHeader := styles.H1TitleStyle().Width(styles.WidthTwoH1OneFourth).Render("System Resourcess")
	rightSideDashboardContent := lipgloss.JoinVertical(lipgloss.Left, getSystemResourcesBox(m.logger), m.getOSPFGeneralInfoBox(m.logger))
	m.viewportRight.SetContent(rightSideDashboardContent)

	rightSideDashboard := lipgloss.JoinVertical(lipgloss.Left, systemResourceHeader, m.viewportRight.View())

	leftSideDashboard := lipgloss.JoinVertical(lipgloss.Left, statusHeader, m.viewportLeft.View())

	horizontalDashboard := lipgloss.JoinHorizontal(lipgloss.Top,
		leftSideDashboard,
		rightSideDashboard,
	)

	return horizontalDashboard
}

func getSystemResourcesBox(logger *logger.Logger) string {
	cpuAmount, cpuUsage, memoryUsage, err := backend.GetSystemResources(logger)
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

	return systemStatistics
}

func (m *Model) getOSPFGeneralInfoBox(logger *logger.Logger) string {
	ospfInformation, err := backend.GetOSPF(logger)
	if err != nil {
		return common.PrintBackendError(err, "GetOSPF")
	}

	lastSPFExecution := time.Duration(ospfInformation.SpfLastExecutedMsecs) * time.Millisecond
	lastSPFExecution = lastSPFExecution.Truncate(time.Second) // remove sub-second precision

	var routerType []string
	switch {
	case ospfInformation.AsbrRouter != "" && ospfInformation.AbrType != "":
		routerType = append(routerType, "Router Type: ASBR / ABR")
	case ospfInformation.AsbrRouter != "":
		routerType = append(routerType, "Router Type: ASBR")
	case ospfInformation.AbrType != "":
		routerType = append(routerType, "Router Type: ABR")
	default:
		routerType = append(routerType, "Router Type: Internal")
	}

	if ospfInformation.AbrType != "" {
		routerType = append(routerType, "ABR Type: "+ospfInformation.AbrType)
	}

	ospfRouterInfo := styles.H1TwoContentBoxesStyle().Width(styles.WidthTwoH1OneFourthBox).Render(
		"OSPF Router ID: " + ospfInformation.RouterId + "\n" +
			strings.Join(routerType, "\n") + "\n" +
			"Last SPF Execution: " + lastSPFExecution.String() + "\n" +
			"Total External LSAs: " + strconv.Itoa(int(ospfInformation.LsaExternalCounter)) + "\n" +
			"Attached Areas: " + strconv.Itoa(int(ospfInformation.AttachedAreaCounter)) + "\n")

	ospfAreas := make([]string, 0, len(ospfInformation.Areas))
	for area := range ospfInformation.Areas {
		ospfAreas = append(ospfAreas, area)
	}
	sort.Strings(ospfAreas)

	var ospfAreaInformation []string
	for _, areaID := range ospfAreas {
		areaData := ospfInformation.Areas[areaID]

		ospfAreaInformation = append(ospfAreaInformation,
			styles.H2TitleStyle().Width(styles.WidthTwoH2OneFourth).Render("Area "+areaID))
		ospfAreaInformation = append(ospfAreaInformation,
			styles.H2TwoContentBoxesStyle().Width(styles.WidthTwoH2OneFourthBox).Render(
				"Full Adjencencies: "+strconv.Itoa(int(areaData.NbrFullAdjacentCounter))+"\n"+
					"Total LSAs: "+strconv.Itoa(int(areaData.LsaNumber))+"\n"+
					"Router LSAs: "+strconv.Itoa(int(areaData.LsaRouterNumber))+"\n"+
					"Network LSAs: "+strconv.Itoa(int(areaData.LsaNetworkNumber))+"\n"+
					"Summary LSAs: "+strconv.Itoa(int(areaData.LsaSummaryNumber))+"\n"+
					"ASBR Summary LSAs: "+strconv.Itoa(int(areaData.LsaAsbrNumber))+"\n"+
					"NSSA External LSAs: "+strconv.Itoa(int(areaData.LsaNssaNumber))))
	}

	renderedOSPFAreaInformation := lipgloss.JoinVertical(lipgloss.Left, ospfAreaInformation...)

	ospfInformationBox := lipgloss.JoinVertical(lipgloss.Left,
		styles.H1TitleStyle().Width(styles.WidthTwoH1OneFourth).Render("General OSPF Information"),
		ospfRouterInfo,
		renderedOSPFAreaInformation,
	)

	return ospfInformationBox
}

func (m *Model) getOspfDashboardLsdbSelf(logger *logger.Logger) string {
	var lsdbSelfBlocks []string

	lsdb, err := backend.GetLSDB(logger)
	if err != nil {
		return common.PrintBackendError(err, "GetLSDB")
	}

	// extract and sort the map keys
	lsdbAreas := make([]string, 0, len(lsdb.Areas))
	for area := range lsdb.Areas {
		lsdbAreas = append(lsdbAreas, area)
	}
	sort.Strings(lsdbAreas)

	_, routerOSPFID, err := backend.GetRouterName(logger)
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

		for _, networkLinkState := range lsaTypes.NetworkLinkStates {
			if networkLinkState.Base.AdvertisedRouter == routerOSPFID {
				networkLinkStateTableData = append(networkLinkStateTableData, []string{
					networkLinkState.Base.LsId,
					networkLinkState.Base.AdvertisedRouter,
					strconv.Itoa(int(networkLinkState.Base.LsaAge)),
				})
			}
		}

		for _, summaryLinkState := range lsaTypes.SummaryLinkStates {
			if summaryLinkState.Base.AdvertisedRouter == routerOSPFID {
				summaryLinkStateTableData = append(summaryLinkStateTableData, []string{
					summaryLinkState.SummaryAddress,
					summaryLinkState.Base.AdvertisedRouter,
					strconv.Itoa(int(summaryLinkState.Base.LsaAge)),
				})
			}
		}

		for _, asbrSummaryLinkState := range lsaTypes.AsbrSummaryLinkStates {
			if asbrSummaryLinkState.Base.AdvertisedRouter == routerOSPFID {
				asbrSummaryLinkStateTableData = append(asbrSummaryLinkStateTableData, []string{
					asbrSummaryLinkState.Base.LsId,
					asbrSummaryLinkState.Base.AdvertisedRouter,
					strconv.Itoa(int(asbrSummaryLinkState.Base.LsaAge)),
				})
			}
		}

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

func (m *Model) getOspfDashboardAnomalies(logger *logger.Logger) string {
	ospfRouterAnomalies, err := backend.GetRouterAnomalies(logger)
	if err != nil {
		return common.PrintBackendError(err, "GetRouterAnomalies")
	}

	ospfExternalAnomalies, err := backend.GetExternalAnomalies(logger)
	if err != nil {
		return common.PrintBackendError(err, "GetExternalAnomalies")
	}

	ospfNSSAExternalAnomalies, err := backend.GetNSSAExternalAnomalies(logger)
	if err != nil {
		return common.PrintBackendError(err, "GetNSSAExternalAnomalies")
	}

	var routerAnomalyTable string
	var routerAnomalyCount int
	if common.HasAnyAnomaly(ospfRouterAnomalies) {
		routerAnomalyCount = countAnomalies(ospfRouterAnomalies)
		routerAnomalyTable = createAnomalyTable(
			ospfRouterAnomalies,
			"Router Anomalies (Type 1 LSAs)",
		)

		logger.WithAttrs(map[string]interface{}{
			"anomaly_type":         "Router (Type 1)",
			"count":                routerAnomalyCount,
			"has_under_advertised": ospfRouterAnomalies.HasUnAdvertisedPrefixes,
			"has_over_advertised":  ospfRouterAnomalies.HasOverAdvertisedPrefixes,
			"has_duplicates":       ospfRouterAnomalies.HasDuplicatePrefixes,
		}).Info("Router anomalies detected")
	}

	var externalAnomalyTable string
	var externalAnomalyCount int
	if common.HasAnyAnomaly(ospfExternalAnomalies) {
		externalAnomalyCount = countAnomalies(ospfExternalAnomalies)
		externalAnomalyTable = createAnomalyTable(
			ospfExternalAnomalies,
			"External Link State Anomalies (Type 5 LSAs)",
		)

		logger.WithAttrs(map[string]interface{}{
			"anomaly_type":         "External (Type 5)",
			"count":                externalAnomalyCount,
			"has_under_advertised": ospfExternalAnomalies.HasUnAdvertisedPrefixes,
			"has_over_advertised":  ospfExternalAnomalies.HasOverAdvertisedPrefixes,
			"has_duplicates":       ospfExternalAnomalies.HasDuplicatePrefixes,
		}).Info("External anomalies detected")
	}

	var nssaExternalAnomalyTable string
	var nssaAnomalyCount int
	if common.HasAnyAnomaly(ospfNSSAExternalAnomalies) {
		nssaAnomalyCount = countAnomalies(ospfNSSAExternalAnomalies)
		nssaExternalAnomalyTable = createAnomalyTable(
			ospfNSSAExternalAnomalies,
			"NSSA External Link State Anomalies (Type 7 LSAs)",
		)

		logger.WithAttrs(map[string]interface{}{
			"anomaly_type":         "NSSA External (Type 7)",
			"count":                nssaAnomalyCount,
			"has_under_advertised": ospfNSSAExternalAnomalies.HasUnAdvertisedPrefixes,
			"has_over_advertised":  ospfNSSAExternalAnomalies.HasOverAdvertisedPrefixes,
			"has_duplicates":       ospfNSSAExternalAnomalies.HasDuplicatePrefixes,
		}).Info("NSSA External anomalies detected")
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

	// Log summary if any anomalies were found
	if len(allAnomaliesList) > 0 {
		logger.WithAttrs(map[string]interface{}{
			"total_anomalies":    routerAnomalyCount + externalAnomalyCount + nssaAnomalyCount,
			"router_anomalies":   routerAnomalyCount,
			"external_anomalies": externalAnomalyCount,
			"nssa_anomalies":     nssaAnomalyCount,
		}).Info("OSPF anomalies summary")
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

	if a.HasUnAdvertisedPrefixes {
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
				"Unadvertised Route",
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

func (m *Model) renderAnomalyDetails() string {
	m.viewport.Width = styles.ViewPortWidthCompletePage
	m.viewport.Height = styles.ViewPortHeightCompletePage

	// ===== IMPORTANT: If a line break happens automatically in the TUI,                ===== //
	// =====            lipgloss renders an extra line which breaks the viewport Height. ===== //
	// ===== Solution:  Use newline '\n' after maximum 149 characters                    ===== //
	// =====            to ensure minimum supported width of FRR-MAD-TUI (157)           ===== //

	anomalyProcessTitle := styles.TextTitleStyle.Padding(0, 0, 0, 0).Render("Anomaly Detection Process")
	anomalyProcessText1 := "The frr-mad-analyzer predicts a 'should-state' for the router based on its static FRR configuration. This includes:\n"
	anomalyPossibilities := []string{
		"Interface addresses that should be announced in Type 1 Router LSAs",
		"Type 5 External LSAs and Type 7 NSSA External LSAs expected from static routes",
	}
	for i, item := range anomalyPossibilities {
		anomalyPossibilities[i] = " > " + item // →
	}
	anomalyProcessText2 := "\nIt then retrieves the 'is-state' using vtysh queries and compares it against the predicted state.\n" +
		"If a mismatch is detected, the anomaly is identified and classified into one of the defined types listed below."

	anomalyTypesTitle := styles.TextTitleStyle.Padding(1, 0, 0, 0).Render("OSPF Anomaly Types")
	anomalyTypes := [][]string{
		{"Unadvertised", "A prefix that is expected to be announced (advertised) to other devices in the network but is missing."},
		{"Overadvertised", "A prefix that is being announced (advertised) to other devices in the network but should not be."},
		{"Duplicated", "A prefix that is present multiple times in the Link-State Database."},
	}
	anomalyTypesTable := components.NewAnomalyTypesTable(
		[]string{
			"Anomaly Type",
			"Description",
		},
		3,
	)
	for _, r := range anomalyTypes {
		anomalyTypesTable = anomalyTypesTable.Row(r...)
	}

	anomalyDetailsOverlay := lipgloss.JoinVertical(lipgloss.Left,
		anomalyProcessTitle,
		anomalyProcessText1,
		strings.Join(anomalyPossibilities, "\n"),
		anomalyProcessText2,
		anomalyTypesTitle,
		anomalyTypesTable.String(),
	)

	m.viewport.SetContent(anomalyDetailsOverlay)

	return m.viewport.View()
}

//func (m *Model) renderExportOptions() string {
//	m.viewportRightHalf.Width = styles.ViewPortWidthHalf
//	m.viewportRightHalf.Height = styles.ViewPortHeightCompletePage - styles.HeightH1
//
//	options := make([]common.ExportOption, len(m.exportOptions))
//	copy(options, m.exportOptions)
//	sort.Slice(options, func(i, j int) bool {
//		return options[i].Label < options[j].Label
//	})
//
//	// 2) Clamp the cursor into [0, len(options)-1]
//	if m.cursor < 0 {
//		m.cursor = 0
//	} else if m.cursor >= len(options) {
//		m.cursor = len(options) - 1
//	}
//
//	s := styles.TextTitleStyle.Render("Choose an option to export:") + "\n\n"
//	for i, opt := range options {
//		cursor := "   "
//		label := opt.Label
//		if i == m.cursor {
//			cursor = styles.SelectedOptionCursorStyle.Render(" ➔ ")
//			label = styles.SelectedOptionStyle.Render(label + " ")
//		}
//		s += fmt.Sprintf("%s%s\n", cursor, label)
//	}
//	s += styles.FooterBoxStyle.Render("\n\n[Tab / Shift+Tab] move selection down/up one option")
//	s += styles.FooterBoxStyle.Render("\n[↑/↓] scroll preview\n")
//	s += styles.FooterBoxStyle.Render("\n[e] quit export options.")
//
//	styledMenu := styles.H1TwoContentBoxCenterStyle().Render(s)
//
//	activeOption := options[m.cursor]
//	preview := m.exportData[activeOption.MapKey]
//	if preview == "" {
//		preview = "<no data for " + activeOption.MapKey + ">"
//	}
//
//	m.viewportRightHalf.SetContent(preview)
//
//	previewHeader := styles.H1TitleStyleForTwo().
//		Render("Preview for: " + activeOption.Filename)
//	exportPreview := lipgloss.JoinVertical(lipgloss.Left,
//		previewHeader,
//		styles.H1TwoContentBoxesStyle().Render(m.viewportRightHalf.View()),
//	)
//
//	return lipgloss.JoinHorizontal(
//		lipgloss.Top,
//		styles.VerticallyCenter(styledMenu, styles.HeightBasis),
//		exportPreview,
//	)
//}

// ============================== //
// HELPERS:                       //
// ============================== //

// countAnomalies return the total amount of detected anomalies
func countAnomalies(a *frrProto.AnomalyDetection) int {
	count := 0
	if a.HasUnAdvertisedPrefixes {
		count += len(a.MissingEntries)
	}
	if a.HasOverAdvertisedPrefixes {
		count += len(a.SuperfluousEntries)
	}
	if a.HasDuplicatePrefixes {
		count += len(a.DuplicateEntries)
	}
	return count
}
