package components

import (
	"github.com/ba2025-ysmprc/frr-tui/internal/common"
	"github.com/ba2025-ysmprc/frr-tui/internal/ui/styles"
	"github.com/charmbracelet/lipgloss"
	"strings"
)

func CreateTabRow(tabs []common.Tab, activeTab int, activeSubTab int, windowSize *common.WindowSize) string {
	var renderedTabs []string
	var renderedSubTabs []string
	for i, tab := range tabs {
		if i == activeTab {
			if activeSubTab != -1 {
				renderedTabs = append(renderedTabs, styles.ActiveTabBoxLockedStyle.Render(tab.Title))
			} else {
				renderedTabs = append(renderedTabs, styles.ActiveTabBoxStyle.Render(tab.Title))
			}
			for j, subTab := range tab.SubTabs {
				if j == activeSubTab {
					renderedSubTabs = append(renderedSubTabs, styles.ActiveSubTabBoxStyle.Render(subTab))
				} else {
					renderedSubTabs = append(renderedSubTabs, styles.InactiveSubTabBoxStyle.Render(subTab))
				}
			}
		} else {
			renderedTabs = append(renderedTabs, styles.InactiveTabBoxStyle.Render(tab.Title))
		}
	}

	// Calculate total tab row width
	tabsWidth := 0
	for _, t := range renderedTabs {
		tabsWidth += lipgloss.Width(t)
	}

	// Calculate total tab row width
	subTabsWidth := 0
	for _, t := range renderedSubTabs {
		subTabsWidth += lipgloss.Width(t)
	}

	// in future call backend to query router name
	routerName := "R101"
	routerNameWidth := lipgloss.Width(routerName)
	routerNameString := "Router Name: "
	routerNameStringWidth := lipgloss.Width(routerNameString)
	routerOSPFID := "65.0.0.1"
	routerOSPFIDWidth := lipgloss.Width(routerOSPFID)
	ospfIdString := "OSPF ID: "
	ospfIdStringWidth := lipgloss.Width(ospfIdString)

	// Build the gap at the right of the last tab based on previous calculation
	remainingWidth := max(0, windowSize.Width-tabsWidth-4)
	leftPadding := max(0, remainingWidth-routerNameWidth-routerNameStringWidth)
	gapContent := strings.Repeat(" ", leftPadding) + routerNameString + routerName
	gap := styles.TabGap.Render(gapContent)
	renderedTabs = append(renderedTabs, gap)

	// Build the gap at the right of the last sub tab based on previous calculation
	remainingWidthSubTab := max(0, windowSize.Width-subTabsWidth-4)
	leftPaddingSubTab := max(0, remainingWidthSubTab-routerOSPFIDWidth-ospfIdStringWidth)
	gapContentSubTab := strings.Repeat(" ", leftPaddingSubTab) + ospfIdString + routerOSPFID
	gapSubTab := gapContentSubTab
	renderedSubTabs = append(renderedSubTabs, gapSubTab)

	horizontalTabs := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedTabs...)
	horizontalSubTabs := lipgloss.JoinHorizontal(lipgloss.Bottom, renderedSubTabs...)

	return lipgloss.JoinVertical(lipgloss.Left, horizontalTabs, horizontalSubTabs)
}
