package ospfMonitoring

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/components"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
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
		return m.renderOSPFTab0()
	} else if currentSubTabLocal == 2 {
		return m.renderOSPFTab1()
	} else if currentSubTabLocal == 3 {
		return m.renderRunningConfigTab()
	}
	return m.renderAdvertisementTab()
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

	shouldAdvertisedRouterLSA := styles.ActiveContentStyle.
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

	isAdvertisedRouterLSA := styles.ActiveContentStyle.
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
