package printers

import (
	"io"
	"sort"
	"strconv"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"

	"github.com/undistro/marvin/pkg/checks"
	"github.com/undistro/marvin/pkg/report"
)

type TablePrinter struct{}

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
	t.SetTablePadding("   ")
	t.SetNoWhiteSpace(true)

	renderTable(r, t)
	return nil
}

func renderTable(r report.Report, t *tablewriter.Table) {
	t.SetHeader([]string{"SEVERITY", "ID", "CHECK", "FAILED", "PASSED", "SKIPPED"})
	sort.Slice(r.Checks, func(i, j int) bool {
		if r.Checks[i].Severity != r.Checks[j].Severity {
			return r.Checks[i].Severity > r.Checks[j].Severity
		}
		return len(r.Checks[i].Failed) > len(r.Checks[j].Failed)
	})
	for _, c := range r.Checks {
		t.Append([]string{
			colorSeverity(c.Severity),
			c.ID,
			c.Message,
			strconv.Itoa(len(c.Failed)),
			strconv.Itoa(len(c.Passed)),
			strconv.Itoa(len(c.Skipped)),
		})
	}
	t.Render()
}

func colorSeverity(s checks.Severity) string {
	switch s {
	case checks.SeverityLow:
		return color.BlueString("%s", s)
	case checks.SeverityMedium:
		return color.YellowString("%s", s)
	case checks.SeverityHigh:
		return color.RedString("%s", s)
	case checks.SeverityCritical:
		return color.RedString("%s", s)
	default:
		return s.String()
	}
}
