package components

import (
	"github.com/charmbracelet/lipgloss"
	ltable "github.com/charmbracelet/lipgloss/table"
	"github.com/frr-mad/frr-mad/src/frontend/internal/ui/styles"
)

func NewOspfMonitorTable(headers []string, rows int) *ltable.Table {
	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == ltable.HeaderRow:
				return styles.HeaderStyle
			default:
				return styles.NormalCellStyle
			}
		})
	return t
}

func NewOspfMonitorMultilineTable(headers []string, rows int) *ltable.Table {
	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == ltable.HeaderRow:
				return styles.HeaderStyle
			case row == rows-1:
				return styles.LastCellOfMultiline
			default:
				return styles.MultilineCellStyle
			}
		})
	return t
}

func NewAnomalyTable(headers []string, rows int) *ltable.Table {
	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == ltable.HeaderRow:
				return styles.HeaderStyle
			default:
				return styles.NormalCellStyle
			}
		})
	return t
}

func NewRibMonitorTable(rows int) *ltable.Table {
	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		Headers(
			"Prefix",
			"Protocol",
			"Next Hops",
			"FIB",
			"Installed",
			"Distance",
			"Metric",
			"Uptime",
		).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == ltable.HeaderRow:
				return styles.HeaderStyle
			case row%2 == 0:
				return styles.EvenRowCell
			default:
				return styles.NormalCellStyle
			}
		})
	return t
}

func NewAnomalyTypesTable(headers []string, rows int) *ltable.Table {
	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == ltable.HeaderRow:
				return styles.HeaderStyle
			default:
				return styles.NormalCellStyle
			}
		})
	return t
}

func NewMultilineTable(headers []string, rows int) *ltable.Table {
	t := ltable.New().
		Border(lipgloss.NormalBorder()).
		BorderTop(true).
		BorderBottom(true).
		BorderLeft(true).
		BorderRight(true).
		BorderHeader(true).
		BorderColumn(true).
		Headers(headers...).
		StyleFunc(func(row, col int) lipgloss.Style {
			switch {
			case row == ltable.HeaderRow:
				return styles.HeaderStyle
			case row == rows-1:
				return styles.LastCellOfMultiline
			default:
				return styles.MultilineCellStyle
			}
		})
	return t
}
