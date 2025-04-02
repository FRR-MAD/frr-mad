package dashboard

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// View renders the dashboard UI.
func (m Model) View() string {
	// Calculate box width dynamically based on terminal width
	boxWidth := (m.windowSize.Width - 14) / 2 // - 6 (padding+margin content) - 6 (for each border)
	if boxWidth < 20 {
		boxWidth = 20 // Minimum width to ensure readability
	}

	// Create boxes for OSPF and BGP
	ospfBox := styles.GeneralBoxStyle.Width(boxWidth).Render(
		"OSPF Anomalies:\n" + strings.Join(m.ospfAnomalies, "\n"),
	)

	bgpBox := styles.GeneralBoxStyle.Width(boxWidth).Render(
		"Helloooo...",
		// "BGP Info:\n" + bgpInfo,
	)

	horizontalBoxes := lipgloss.JoinHorizontal(lipgloss.Top, ospfBox, bgpBox)

	infoBox := styles.InfoTextStyle.Width(m.windowSize.Width - 12).Foreground(lipgloss.Color("#C0C0C0")).Render("press 'r' to refresh dashboard")

	return lipgloss.JoinVertical(lipgloss.Left, horizontalBoxes, infoBox)

	//// Use the styles defined in styles.go
	//title := styles.TitleStyle.Render(m.Title)
	//metrics := strings.Join(m.Metrics, "\n")
	//body := styles.BodyStyle.Render(metrics)
	//
	//return lipgloss.JoinVertical(lipgloss.Left, title, body, "\nPress 'r' to refresh metrics.")
}
