package coverage

import (
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

// Table represents a table definition with headers and rows
type Table struct {
	Table table.Writer
}

const (
	colTitleIndex      = "#"
	colTitleCoverage   = "Coverage"
	colTitlePercentage = "Percentage"
	colTitleRatio      = "Ratio"
)

var (
	rowHeader = table.Row{colTitleIndex, colTitleCoverage, colTitlePercentage, colTitleRatio}
)

// InitializeTable initializes and returns a new table writer with the specified header and style settings.
// Parameters:
//
//	header - the header row to be set for the table.
//	source - a pointer to a slice of objects containing the information to be displayed.
//
// Returns:
//
//	A configured table instance.
func InitializeTable() Table {
	t := table.NewWriter()
	t.SetTitle("Coverage Summary")
	t.Style().Title.Align = text.AlignCenter
	t.Style().Title.Colors = text.Colors{text.BgBlue, text.FgHiMagenta, text.Bold}
	t.AppendHeader(rowHeader)
	t.SetPageSize(0) // disables paging
	t.Style().Options.DrawBorder = false
	colorBOnW := text.Colors{text.BgWhite, text.FgBlack}
	t.SetColumnConfigs([]table.ColumnConfig{
		{Name: colTitleIndex, Colors: text.Colors{text.FgYellow}, ColorsHeader: colorBOnW},
		{Name: colTitleCoverage, Colors: text.Colors{text.FgHiRed}, ColorsHeader: colorBOnW},
		{Name: colTitlePercentage, Colors: text.Colors{text.FgHiRed}, ColorsHeader: colorBOnW},
		{Name: colTitleRatio, Colors: text.Colors{text.FgHiRed}, ColorsHeader: colorBOnW, ColorsFooter: colorBOnW},
		{Number: 5, Colors: text.Colors{text.FgCyan}, ColorsHeader: colorBOnW},
	})

	tt := Table{Table: t}
	tt.Table.AppendRow(table.Row{1, "Manifests", "0%", "0/2"})
    tt.Table.AppendRow(table.Row{2, "Documents", "0%", "0/10"})
	tt.Table.AppendRow(table.Row{3, "Lines", "0%", "0/50"})
	return tt
}

// TODO: one more table e.g. per manifest

//           Coverage Summary
//  # | COVERAGE  | PERCENTAGE | RATIO
// ---+-----------+------------+-------
//  1 | Manifests | 0%         | 0/2
//  2 | Documents | 0%         | 0/10
//  3 | Lines     | 0%         | 0/50
