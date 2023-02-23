package printers

import (
	"io"
	"strconv"

	"github.com/olekukonko/tablewriter"

	"github.com/undistro/marvin/pkg/report"
)

type TablePrinter struct {
}

func (*TablePrinter) PrintObj(r report.Report, w io.Writer) error {
	t := tablewriter.NewWriter(w)
	t.SetAutoWrapText(false)
	t.SetAutoFormatHeaders(true)
	t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetCenterSeparator("")
	t.SetColumnSeparator("")
	t.SetRowSeparator("")
	t.SetHeaderLine(false)
	t.SetBorder(false)
	t.SetTablePadding("\t")
	t.SetNoWhiteSpace(true)
	t.SetHeader([]string{"SEVERITY", "CHECK", "FAILED"})
	for _, c := range r.Checks {
		t.Append([]string{c.Severity.String(), c.Message, strconv.Itoa(len(c.Failed))})
	}
	t.Render()
	return nil
}
