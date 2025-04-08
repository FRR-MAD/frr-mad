package shell

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

const tabRowHeight = 6
const footerHeight = 1

var currentSubTabLocal = -1

func (m *Model) ShellView(currentSubTab int) string {
	currentSubTabLocal = currentSubTab
	return m.View()
}

func (m *Model) View() string {
	if currentSubTabLocal == 0 {
		return m.RenderSubTab0()
	} else if currentSubTabLocal == 1 {
		return m.RenderSubTab1()
	}
	return m.RenderSubTab0()
}

func (m *Model) RenderSubTab0() string {
	// Calculate box width dynamically for one horizontal box based on terminal width
	boxWidthForOne := m.windowSize.Width - 10 // - 6 (padding+margin content) - 2 (for each border)
	if boxWidthForOne < 20 {
		boxWidthForOne = 20 // Minimum width to ensure readability
	}

	inputHeight := 3

	outputMaxHeight := m.windowSize.Height - tabRowHeight - footerHeight - inputHeight - 2

	// Update the viewport dimensions.
	m.viewport.Width = boxWidthForOne
	m.viewport.Height = outputMaxHeight

	// Update the viewport content with the latest BashOutput.
	m.viewport.SetContent(m.BashOutput)

	input := "Type bash command: "
	if currentSubTabLocal == -1 {
		input = styles.InactiveBoxStyle.Width(boxWidthForOne).Render(input)
	} else if currentSubTabLocal == 0 {
		input += m.BashInput
		input = styles.GeneralBoxStyle.Width(boxWidthForOne).Render(input)
	}

	// return lipgloss.JoinVertical(lipgloss.Left, input, styles.TextOutputStyle.Render(m.BashOutput))

	return lipgloss.JoinVertical(lipgloss.Left,
		input,
		styles.TextOutputStyle.Render(m.viewport.View()))
}

func (m *Model) RenderSubTab1() string {
	// Calculate box width dynamically for one horizontal box based on terminal width
	boxWidthForOne := m.windowSize.Width - 10 // - 6 (padding+margin content) - 2 (for each border)
	if boxWidthForOne < 20 {
		boxWidthForOne = 20 // Minimum width to ensure readability
	}

	input := "Type vtysh command: " + m.VtyshInput
	input = styles.GeneralBoxStyle.Width(boxWidthForOne).Render(input)

	return lipgloss.JoinVertical(lipgloss.Left, input, m.VtyshOutput)
}
