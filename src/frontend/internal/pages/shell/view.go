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
	var content string
	var body string
	var bodyFooter string

	statusBar := true

	switch currentSubTabLocal {
	case 0:
		body = m.renderBashShellTab()
	case 1:
		body = m.renderVtyshShellTab()
	case 2:
		body = m.renderBackendTestTab()
	default:
		body = m.renderBashShellTab()
	}

	if statusBar {
		statusBox := lipgloss.NewStyle().Width(styles.WidthTwoH1Box).Margin(0, 2).Render(m.statusMessage)
		if m.statusMessage != "" {
			styles.SetStatusSeverity(m.statusSeverity)
			if len(m.statusMessage) > 50 {
				m.statusMessage = m.statusMessage[:47] + "..."
			}
			renderedStatusMessage := styles.StatusTextStyle().Render(m.statusMessage)
			statusBox = lipgloss.NewStyle().Width(styles.WidthTwoH1Box).Margin(0, 2).Render(renderedStatusMessage)
		}

		bodyFooter = lipgloss.JoinHorizontal(lipgloss.Top, statusBox)

		content = lipgloss.JoinVertical(lipgloss.Left, body, bodyFooter)
	} else {
		content = body
	}

	return content
}

func (m *Model) renderBashShellTab() string {
	if m.readOnlyMode {
		return lipgloss.NewStyle().Height(styles.HeightBasis - styles.BodyFooterHeight).
			Render("You are in read only mode. Press [ctrl+w] to deactivate it.")
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

func (m *Model) renderVtyshShellTab() string {
	if m.readOnlyMode {
		return lipgloss.NewStyle().Height(styles.HeightBasis - styles.BodyFooterHeight).
			Render("You are in read only mode. Press [ctrl+w] to deactivate it.")
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
