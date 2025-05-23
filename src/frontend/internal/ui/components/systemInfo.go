package components

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/frr-mad/frr-tui/internal/ui/styles"
)

func GetSystemInfoOverlay() string {
	frrMADLegendTitle := styles.TextTitleStyle.Render("FRR-MAD-TUI Mode")
	frrMADRead := styles.ContentBoxStyle().
		BorderForeground(lipgloss.Color(styles.ReadModeBlue)).
		Render("read mode")
	frrMADWrite := styles.ContentBoxStyle().
		BorderForeground(lipgloss.Color(styles.WriteModeCoral)).
		Render("write mode")

	frrMADTUILegend := lipgloss.JoinVertical(lipgloss.Left,
		frrMADLegendTitle,
		frrMADRead,
		frrMADWrite,
		"\n",
	)

	anomalyDetectionLegendTitle := styles.TextTitleStyle.Render("Anomaly Detection")
	monitoringLegendTitle := styles.TextTitleStyle.Render("Monitoring")

	titleAnomaly := styles.H1BadTitleStyle().
		Padding(0, 2).Margin(0, 0, 1, 0).
		BorderBottom(true).
		Render("Title for detected anomalies")
	titleNoAnomaly := styles.H1GoodTitleStyle().
		Padding(0, 2).Margin(0, 0, 1, 0).
		BorderBottom(true).
		Render("Title for no issues detected")

	titleMonitoringH1 := styles.H1TitleStyle().
		Padding(1, 2, 0, 2).Margin(0, 0, 1, 0).
		Render("Title for monitoring (H1)")
	titleMonitoringH2 := styles.H2TitleStyle().
		Padding(0, 1).Margin(0, 0, 1, 0).
		Render("Title for monitoring (H2)")

	anomalyDetectionLegend := lipgloss.JoinVertical(lipgloss.Left,
		anomalyDetectionLegendTitle,
		titleAnomaly,
		titleNoAnomaly,
		monitoringLegendTitle,
		titleMonitoringH1,
		titleMonitoringH2,
		"\n",
	)

	messageLegendTitle := styles.TextTitleStyle.Render("Status Messages")

	InfoMessage := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.InfoStatusColor)).
		Background(lipgloss.Color(styles.InfoStatusBackground)).
		Padding(0, 1).
		Margin(0, 0, 1, 0).
		Render("Info Status Message")

	WarningMessage := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.WarningStatusColor)).
		Background(lipgloss.Color(styles.WarningStatusBackground)).
		Padding(0, 1).
		Margin(0, 0, 1, 0).
		Render("Warning Status Message")

	ErrorMessage := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ErrorStatusColor)).
		Background(lipgloss.Color(styles.ErrorStatusBackground)).
		Padding(0, 1).
		Margin(0, 0, 1, 0).
		Render("Error Status Message")

	messageLegend := lipgloss.JoinVertical(lipgloss.Left,
		messageLegendTitle,
		InfoMessage,
		WarningMessage,
		ErrorMessage,
		"\n",
	)

	firstCol := lipgloss.NewStyle().Width(styles.WidthBasis / 3).
		Render(lipgloss.JoinVertical(lipgloss.Left, anomalyDetectionLegend))
	secondCol := lipgloss.NewStyle().Width(styles.WidthBasis / 3).
		Render(lipgloss.JoinVertical(lipgloss.Left, frrMADTUILegend, messageLegend))

	return lipgloss.JoinHorizontal(lipgloss.Top, firstCol, secondCol)
}
