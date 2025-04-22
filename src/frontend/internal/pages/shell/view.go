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
	} else if currentSubTabLocal == 2 {
		return m.renderBackendTestTab()
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

func (m *Model) renderBackendTestTab() string {

	testInfo := "To Test the Backend we need a service and a command e.g. 'ospf' / 'database'\n" +
		"press 'tab' to switch to command input. press 'enter' to send backend call.\n"

	var serviceBox, commandBox string
	if m.activeBackendInput == "service" {
		serviceBox = lipgloss.JoinVertical(lipgloss.Left,
			styles.BoxTitleStyle.Render("Enter Service:"),
			styles.GeneralBoxStyle.Width(20).Render(m.backendServiceInput),
		)

		commandBox = lipgloss.JoinVertical(lipgloss.Left,
			styles.BoxTitleStyle.Render("Enter Command:"),
			styles.InactiveBoxStyle.Width(20).Render(m.backendCommandInput),
		)
	} else {
		serviceBox = lipgloss.JoinVertical(lipgloss.Left,
			styles.BoxTitleStyle.Render("Enter Service:"),
			styles.InactiveBoxStyle.Width(20).Render(m.backendServiceInput),
		)

		commandBox = lipgloss.JoinVertical(lipgloss.Left,
			styles.BoxTitleStyle.Render("Enter Command:"),
			styles.GeneralBoxStyle.Width(20).Render(m.backendCommandInput),
		)
	}

	inputsHorizontal := lipgloss.JoinHorizontal(lipgloss.Top,
		serviceBox,
		lipgloss.NewStyle().Margin(0, 0, 0, 2).Render(commandBox),
		lipgloss.NewStyle().Margin(0, 0, 0, 4).Render(testInfo),
	)

	outputMaxHeight := m.windowSize.Height - styles.TabRowHeight - styles.FooterHeight - 8

	// Update the viewport dimensions.
	m.viewport.Width = m.windowSize.Width - 6
	m.viewport.Height = outputMaxHeight

	// Update the viewport content with the latest backendResponse.
	m.viewport.SetContent(m.backendResponse)

	completeTab := lipgloss.JoinVertical(lipgloss.Left,
		inputsHorizontal,
		styles.TextOutputStyle.Render(m.viewport.View()),
	)

	return completeTab
}
