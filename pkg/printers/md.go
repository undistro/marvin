package printers

import (
	"io"

	"github.com/olekukonko/tablewriter"

	"github.com/undistro/marvin/pkg/report"
)

type MarkdownPrinter struct{}

func (*MarkdownPrinter) PrintObj(r report.Report, w io.Writer) error {
	t := tablewriter.NewWriter(w)
	t.SetAutoWrapText(false)
	t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetBorders(tablewriter.Border{Left: true, Top: false, Right: true, Bottom: false})
	t.SetCenterSeparator("|")

	renderTable(r, t)
	return nil
}
