package components

import (
	"fmt"
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	"sort"
)

func RenderExportOptions(
	exportOptions []common.ExportOption,
	exportData map[string]string,
	cursor *int,
	vp *viewport.Model,
) string {
	// adjust viewport dimensions if needed
	vp.Width = styles.ViewPortWidthHalf
	vp.Height = styles.ViewPortHeightCompletePage - styles.HeightH1

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
	s += styles.FooterBoxStyle.Render("\n\n[Tab Shift+Tab] move selection down/up one option")
	s += styles.FooterBoxStyle.Render("\n[↑ ↓ home end] scroll preview\n")
	s += styles.FooterBoxStyle.Render("\n[e] quit export options | [enter] export current selection")

	//i := styles.FooterBoxStyle.Render("[Tab Shift+Tab] move selection down/up one option")
	//i += styles.FooterBoxStyle.Render("\n[↑ ↓ home end] scroll preview\n")
	//i2 := styles.FooterBoxStyle.Render("[enter] export current selection")
	//i2 += styles.FooterBoxStyle.Render("\n[e] quit export options")

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

	// final horizontal layout
	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		styles.VerticallyCenter(menu, styles.HeightBasis),
		exportPreview,
	)
}
