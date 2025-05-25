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
		Render("Info Message")

	WarningMessage := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.WarningStatusColor)).
		Background(lipgloss.Color(styles.WarningStatusBackground)).
		Padding(0, 1).
		Margin(0, 0, 1, 0).
		Render("Warning Message")

	ErrorMessage := lipgloss.NewStyle().
		Foreground(lipgloss.Color(styles.ErrorStatusColor)).
		Background(lipgloss.Color(styles.ErrorStatusBackground)).
		Padding(0, 1).
		Margin(0, 0, 1, 0).
		Render("Error Message")

	messageLegend := lipgloss.JoinVertical(lipgloss.Left,
		messageLegendTitle,
		InfoMessage,
		WarningMessage,
		ErrorMessage,
		"\n",
	)

	keyboardOptionsTitle := styles.TextTitleStyle.Render("Keyboard Options")

	generalOptionsTitle := lipgloss.NewStyle().Bold(true).Render("General")
	generalOptions := []string{
		"Ctrl+C     | Quit FRR-MAD (always active)",
		"Ctrl+W     | Toggle FRR-MAD-TUI mode",
		"i          | Toggle system info (only when no input field is focused)",
		"Enter      | Enter sub-tabs",
		"Esc        | Exit sub-tab",
		"↑ ↓        | Scroll 10 lines",
		"End Home   | Scroll full page",
	}

	filterOptionsTitle := lipgloss.NewStyle().Bold(true).Render("\nFilter Content")
	filterOptions := []string{
		":          | Toggle table filter (works on pages with a filter at the bottom right)",
		"Esc        | Disable filter (retains last query)",
		"Enter      | Apply filter after re-enabling",
	}

	exportOptionsTitle := lipgloss.NewStyle().Bold(true).Render("\nExport Data")
	exportOptions := []string{
		"Ctrl+E     | Toggle export options page",
		"Esc        | Close export options page",
		"Tab        | Move selection down",
		"Shift+Tab  | Move selection up",
		"Enter      | Export selected option to the given path",
		"           | (If OSC52 is supported in local terminal, also copies to clipboard)",
	}

	anomalyInfoTitle := lipgloss.NewStyle().Bold(true).Render("\nAnomaly Details")
	anomalyOptions := []string{
		"Ctrl+A     | Toggle anomaly details page (only on dashboard)",
		"Esc        | Close anomaly details page",
	}

	renderedGeneralOptions := lipgloss.JoinVertical(lipgloss.Left, generalOptions...)
	renderedExportOptions := lipgloss.JoinVertical(lipgloss.Left, exportOptions...)
	renderedFilterOptions := lipgloss.JoinVertical(lipgloss.Left, filterOptions...)
	renderedAnomalyOptions := lipgloss.JoinVertical(lipgloss.Left, anomalyOptions...)

	keyboardOptions := lipgloss.JoinVertical(lipgloss.Left,
		keyboardOptionsTitle,
		generalOptionsTitle,
		renderedGeneralOptions,
		filterOptionsTitle,
		renderedFilterOptions,
		exportOptionsTitle,
		renderedExportOptions,
		anomalyInfoTitle,
		renderedAnomalyOptions,
	)

	firstCol := lipgloss.NewStyle().Width(38).Margin(0, 2, 0, 0).
		Render(lipgloss.JoinVertical(lipgloss.Left, anomalyDetectionLegend))
	secondCol := lipgloss.NewStyle().Width(23).Margin(0, 2, 0, 0).
		Render(lipgloss.JoinVertical(lipgloss.Left, frrMADTUILegend, messageLegend))
	thirdCol := lipgloss.NewStyle().Width(styles.WidthBasis - 65).
		Render(lipgloss.JoinVertical(lipgloss.Left, keyboardOptions))

	return lipgloss.JoinHorizontal(lipgloss.Top, firstCol, secondCol, thirdCol)
}
