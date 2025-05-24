package components

import (
	"fmt"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"github.com/frr-mad/frr-tui/internal/common"
	"github.com/frr-mad/frr-tui/internal/ui/styles"
	"sort"
)

func RenderExportOptions(
	exportOptions []common.ExportOption,
	exportData map[string]string,
	cursor *int,
	vp *viewport.Model,
	statusMessage string,
	statusSeverity styles.StatusSeverity,
) string {
	// adjust viewport dimensions if needed
	vp.Width = styles.WidthViewPortHalf
	vp.Height = styles.HeightViewPortCompletePage - styles.HeightH1 - styles.AdditionalFooterHeight - styles.BodyFooterHeight

	// copy & sort options by label
	opts := make([]common.ExportOption, len(exportOptions))
	copy(opts, exportOptions)
	sort.Slice(opts, func(i, j int) bool {
		return opts[i].Label < opts[j].Label
	})

	// clamp cursor
	if *cursor < 0 {
		*cursor = 0
	} else if *cursor >= len(opts) {
		*cursor = len(opts) - 1
	}

	// build menu
	s := styles.TextTitleStyle.Render("Choose an option to export:") + "\n\n"
	for i, opt := range opts {
		prefix := "   "
		label := opt.Label
		if i == *cursor {
			prefix = styles.SelectedOptionCursorStyle.Render(" ➔ ")
			label = styles.SelectedOptionStyle.Render(label + " ")
		}
		s += fmt.Sprintf("%s%s\n", prefix, label)
	}

	menu := styles.H1TwoContentBoxCenterStyle().Render(s)

	// select active and preview content
	active := opts[*cursor]
	preview := exportData[active.MapKey]
	if preview == "" {
		preview = "<no data for " + active.MapKey + ">"
	}
	vp.SetContent(preview)

	// assemble preview pane
	header := styles.H1TitleStyleForTwo().Render("Preview for: " + active.Filename)
	exportPreview := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		styles.H1TwoContentBoxesStyle().Render(vp.View()),
	)

	horizontalContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		styles.VerticallyCenter(menu, styles.HeightBasis-styles.AdditionalFooterHeight),
		exportPreview,
	)

	statusBox := lipgloss.NewStyle().Width(styles.WidthTwoH1Box).Margin(0, 2).Render(statusMessage)
	if statusMessage != "" {
		styles.SetStatusSeverity(statusSeverity)
		var cutToSizeMessage string
		if len(statusMessage) > (styles.WidthTwoH1Box - styles.MarginX2) {
			cutToSizeMessage = statusMessage[:styles.WidthTwoH1Box-styles.MarginX2-3] + "..."
		} else {
			cutToSizeMessage = statusMessage
		}
		renderedStatusMessage := styles.StatusTextStyle().Render(cutToSizeMessage)
		statusBox = lipgloss.NewStyle().Width(styles.WidthTwoH1Box).Margin(0, 2).Render(renderedStatusMessage)
	}

	keyboardOptions := styles.FooterBoxStyle.Render("\n" +
		"[Tab Shift+Tab] move selection down/up one option | " +
		"[enter] export current selection to file and clipboard | " +
		"[ctrl+e] quit export options")

	return lipgloss.JoinVertical(lipgloss.Left,
		horizontalContent,
		statusBox,
		keyboardOptions,
	)
}
