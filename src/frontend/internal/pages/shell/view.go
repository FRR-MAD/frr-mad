package shell

import (
	"github.com/frr-mad/frr-tui/internal/ui/styles"

	"github.com/charmbracelet/lipgloss"
)

var currentSubTabLocal = -1

func (m *Model) ShellView(currentSubTab int, readOnlyMode bool) string {
	currentSubTabLocal = currentSubTab
	m.readOnlyMode = readOnlyMode
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
	if m.readOnlyMode {
		return "You are in read only mode. Press [ctrl+w] to deactivate it."
	}

	// Update the viewport dimensions.
	m.viewport.Width = styles.WidthViewPortCompletePage
	m.viewport.Height = styles.HeightViewPortCompletePage - styles.HeightH1 - 1

	// Update the viewport content with the latest bashOutput.
	m.viewport.SetContent(m.bashOutput)

	input := "Type bash command: "
	if currentSubTabLocal == -1 {
		input = styles.InactiveBoxStyle.Width(styles.WidthOneH1Box).Render(input)
	} else if currentSubTabLocal == 0 {
		input += m.bashInput
		input = styles.GeneralBoxStyle.Width(styles.WidthOneH1Box).Render(input)
	}

	// return lipgloss.JoinVertical(lipgloss.Left, input, styles.TextOutputStyle.Render(m.bashOutput))

	return lipgloss.JoinVertical(lipgloss.Left,
		input,
		styles.TextOutputStyle.Render(m.viewport.View()))
}

func (m *Model) renderShellTab1() string {
	if m.readOnlyMode {
		return "You are in read only mode. Press [ctrl+w] to deactivate it."
	}

	m.activeShell = "vtysh"

	// Update the viewport dimensions.
	m.viewport.Width = styles.WidthViewPortCompletePage
	m.viewport.Height = styles.HeightViewPortCompletePage - styles.HeightH1 - 1

	// Update the viewport content with the latest vtyshOutput.
	m.viewport.SetContent(m.vtyshOutput)

	input := "Type vtysh command: " + m.vtyshInput
	input = styles.GeneralBoxStyle.Width(styles.WidthOneH1Box).Render(input)

	//return lipgloss.JoinVertical(lipgloss.Left, input, m.vtyshOutput)

	return lipgloss.JoinVertical(lipgloss.Left,
		input,
		styles.TextOutputStyle.Render(m.viewport.View()))
}

func (m *Model) renderBackendTestTab() string {

	testInfo := "To Test the Backend we need a service and a command e.g. 'ospf' / 'database'\n" +
		"Press 'tab' to switch to command input. press 'enter' to send backend call.\n" // +
	// "press '' to copy output to clipboard."

	var serviceBox, commandBox string
	if m.activeBackendInput == "service" {
		serviceBox = lipgloss.JoinVertical(lipgloss.Left,
			styles.TextTitleStyle.Render("Enter Service:"),
			styles.GeneralBoxStyle.Width(20).Render(m.backendServiceInput),
		)

		commandBox = lipgloss.JoinVertical(lipgloss.Left,
			styles.TextTitleStyle.Render("Enter Command:"),
			styles.InactiveBoxStyle.Width(20).Render(m.backendCommandInput),
		)
	} else {
		serviceBox = lipgloss.JoinVertical(lipgloss.Left,
			styles.TextTitleStyle.Render("Enter Service:"),
			styles.InactiveBoxStyle.Width(20).Render(m.backendServiceInput),
		)

		commandBox = lipgloss.JoinVertical(lipgloss.Left,
			styles.TextTitleStyle.Render("Enter Command:"),
			styles.GeneralBoxStyle.Width(20).Render(m.backendCommandInput),
		)
	}

	inputsHorizontal := lipgloss.JoinHorizontal(lipgloss.Top,
		serviceBox,
		lipgloss.NewStyle().Margin(0, 0, 0, 2).Render(commandBox),
		lipgloss.NewStyle().Margin(0, 0, 0, 4).Render(testInfo),
	)

	// Update the viewport dimensions.
	m.viewport.Width = styles.WidthViewPortCompletePage
	m.viewport.Height = styles.HeightViewPortCompletePage - 7

	// Update the viewport content with the latest backendResponse.
	m.viewport.SetContent(m.backendResponse)

	completeTab := lipgloss.JoinVertical(lipgloss.Left,
		inputsHorizontal,
		styles.TextOutputStyle.Render(m.viewport.View()),
	)

	return completeTab
}
