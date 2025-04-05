package ospfMonitoring

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/components"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
)

func (m *Model) View() string {
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
