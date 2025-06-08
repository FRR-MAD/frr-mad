package shell

import (
	"fmt"
	"strings"

	"github.com/frr-mad/frr-tui/internal/ui/styles"

	"slices"

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
			var cutToSizeMessage string
			maxLength := styles.WidthTwoH1Box - styles.MarginX2 - 3
			if maxLength > 0 && len(m.statusMessage) > maxLength {
				cutToSizeMessage = m.statusMessage[:maxLength] + "..."
			} else if maxLength > 0 {
				cutToSizeMessage = m.statusMessage
			} else {
				cutToSizeMessage = "..."
			}
			statusMessage := styles.StatusTextStyle().Render(cutToSizeMessage)
			statusBox = lipgloss.NewStyle().Width(styles.WidthTwoH1Box).Margin(0, 2).Render(statusMessage)
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

	if strings.TrimSpace(m.bashOutput) != "" {
		lines := strings.Split(strings.TrimSpace(m.bashOutput), "\n")
		if len(lines) > 50 {
			m.statusMessage = fmt.Sprintf("Output contains %d lines - scroll to view all", len(lines))
			m.statusSeverity = styles.SeverityInfo
		}
	}

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
		m.statusMessage = "Read-only mode active"
		m.statusSeverity = styles.SeverityWarning
		return lipgloss.NewStyle().Height(styles.HeightBasis - styles.BodyFooterHeight).
			Render("You are in read only mode. Press [ctrl+w] to deactivate it.")
	}

	m.activeShell = "vtysh"

	if currentSubTabLocal == 1 && m.statusMessage == "Read-only mode active" {
		m.statusMessage = "VTY shell ready"
		m.statusSeverity = styles.SeverityInfo
	}

	// Update the viewport dimensions.
	m.viewport.Width = styles.WidthViewPortCompletePage
	m.viewport.Height = styles.HeightViewPortCompletePage - styles.HeightH1 - 1

	// Update the viewport content with the latest vtyshOutput.
	m.viewport.SetContent(m.vtyshOutput)

	if strings.TrimSpace(m.vtyshOutput) != "" {
		lines := strings.Split(strings.TrimSpace(m.vtyshOutput), "\n")
		if len(lines) > 50 {
			m.statusMessage = fmt.Sprintf("VTY output contains %d lines - scroll to view all", len(lines))
			m.statusSeverity = styles.SeverityInfo
		}
		if strings.Contains(strings.ToLower(m.vtyshOutput), "unknown command") {
			m.statusMessage = "Unknown command - Try 'show running-config' or '?'"
			m.statusSeverity = styles.SeverityWarning
		} else if strings.Contains(strings.ToLower(m.vtyshOutput), "error") {
			m.statusMessage = "Command returned an error - check syntax"
			m.statusSeverity = styles.SeverityError
		}
	}

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

		if strings.TrimSpace(m.backendServiceInput) == "" {
			m.statusMessage = "Enter service name (e.g., frr, ospf, analysis)"
			m.statusSeverity = styles.SeverityInfo
		} else {
			validServices := []string{"frr", "ospf", "analysis", "system", "zebra"}
			isValid := slices.Contains(validServices, strings.ToLower(m.backendServiceInput))
			if !isValid {
				m.statusMessage = fmt.Sprintf("Service '%s' may not be valid", m.backendServiceInput)
				m.statusSeverity = styles.SeverityWarning
			} else {
				m.statusMessage = fmt.Sprintf("Service '%s' selected - press TAB for command input", m.backendServiceInput)
				m.statusSeverity = styles.SeverityInfo
			}
		}
	} else {
		serviceBox = lipgloss.JoinVertical(lipgloss.Left,
			styles.TextTitleStyle.Render("Enter Service:"),
			styles.InactiveBoxStyle.Width(20).Render(m.backendServiceInput),
		)

		commandBox = lipgloss.JoinVertical(lipgloss.Left,
			styles.TextTitleStyle.Render("Enter Command:"),
			styles.GeneralBoxStyle.Width(20).Render(m.backendCommandInput),
		)

		if strings.TrimSpace(m.backendCommandInput) == "" {
			m.statusMessage = "Enter command (e.g., database, generalInfo, summary)"
			m.statusSeverity = styles.SeverityInfo
		} else if strings.TrimSpace(m.backendServiceInput) != "" {
			m.statusMessage = fmt.Sprintf("Ready to test: %s/%s - Press ENTER to execute", m.backendServiceInput, m.backendCommandInput)
			m.statusSeverity = styles.SeverityInfo
		}
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
