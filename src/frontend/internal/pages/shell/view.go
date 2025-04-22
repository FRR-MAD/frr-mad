package shell

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

const inputHeight = 3

var currentSubTabLocal = -1

func (m *Model) ShellView(currentSubTab int) string {
	currentSubTabLocal = currentSubTab
	return m.View()
}

func (m *Model) View() string {
	if currentSubTabLocal == 0 {
		return m.renderShellTab0()
	} else if currentSubTabLocal == 1 {
		return m.renderShellTab1()
	}
	return m.renderShellTab0()
}

func (m *Model) renderShellTab0() string {
	// Calculate box width dynamically for one horizontal box based on terminal width
	boxWidthForOne := m.windowSize.Width - 10 // - 6 (padding+margin content) - 2 (for each border)
	if boxWidthForOne < 20 {
		boxWidthForOne = 20 // Minimum width to ensure readability
	}

	outputMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight - inputHeight - 2

	// Update the viewport dimensions.
	m.viewport.Width = boxWidthForOne
	m.viewport.Height = outputMaxHeight

	// Update the viewport content with the latest bashOutput.
	m.viewport.SetContent(m.bashOutput)

	input := "Type bash command: "
	if currentSubTabLocal == -1 {
		input = styles.InactiveBoxStyle.Width(boxWidthForOne).Render(input)
	} else if currentSubTabLocal == 0 {
		input += m.bashInput
		input = styles.GeneralBoxStyle.Width(boxWidthForOne).Render(input)
	}

	// return lipgloss.JoinVertical(lipgloss.Left, input, styles.TextOutputStyle.Render(m.bashOutput))

	return lipgloss.JoinVertical(lipgloss.Left,
		input,
		styles.TextOutputStyle.Render(m.viewport.View()))
}

func (m *Model) renderShellTab1() string {
	// Calculate box width dynamically for one horizontal box based on terminal width
	boxWidthForOne := m.windowSize.Width - 10 // - 6 (padding+margin content) - 2 (for each border)
	if boxWidthForOne < 20 {
		boxWidthForOne = 20 // Minimum width to ensure readability
	}

	outputMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight - inputHeight - 2

	// Update the viewport dimensions.
	m.viewport.Width = boxWidthForOne
	m.viewport.Height = outputMaxHeight

	// Update the viewport content with the latest vtyshOutput.
	m.viewport.SetContent(m.vtyshOutput)

	input := "Type vtysh command: " + m.vtyshInput
	input = styles.GeneralBoxStyle.Width(boxWidthForOne).Render(input)

	//return lipgloss.JoinVertical(lipgloss.Left, input, m.vtyshOutput)

	return lipgloss.JoinVertical(lipgloss.Left,
		input,
		styles.TextOutputStyle.Render(m.viewport.View()))
}
