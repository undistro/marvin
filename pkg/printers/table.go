// Copyright 2023 Undistro Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package printers

import (
	"io"
	"sort"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"

	"github.com/undistro/marvin/pkg/types"
)

var (
	red        = color.New(color.FgRed).SprintfFunc()
	redBold    = color.New(color.FgRed, color.Bold).SprintfFunc()
	yellow     = color.New(color.FgYellow).SprintfFunc()
	yellowBold = color.New(color.FgYellow, color.Bold)
	blue       = color.New(color.FgBlue).SprintfFunc()
	green      = color.New(color.FgGreen).SprintfFunc()

	zoraBanner = `
           Now you can use Marvin as a Zora plugin and see the results in a dashboard.
           Access the documentation for more details:    https://zora-docs.undistro.io
`
)

// TablePrinter implements a Printer that prints the report in table format
type TablePrinter struct {
	DisableZoraBanner bool
}

func (r *TablePrinter) PrintObj(report types.Report, w io.Writer) error {
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

	renderTable(report, t)

	if !r.DisableZoraBanner {
		yellowBold.Fprintln(w, zoraBanner)
	}
	return nil
}

func renderTable(report types.Report, t *tablewriter.Table) {
	t.SetHeader([]string{"SEVERITY", "ID", "CHECK", "STATUS", "FAILED", "PASSED", "SKIPPED"})
	sort.Slice(report.Checks, func(i, j int) bool {
		if report.Checks[i].Severity != report.Checks[j].Severity {
			return report.Checks[i].Severity > report.Checks[j].Severity
		}
		if report.Checks[i].TotalFailed != report.Checks[j].TotalFailed {
			return report.Checks[i].TotalFailed > report.Checks[j].TotalFailed
		}
		if report.Checks[i].TotalPassed != report.Checks[j].TotalPassed {
			return report.Checks[i].TotalPassed > report.Checks[j].TotalPassed
		}
		return strings.Compare(report.Checks[i].ID, report.Checks[j].ID) < 0
	})
	for _, c := range report.Checks {
		t.Append([]string{
			colorSeverity(c.Severity),
			c.ID,
			c.Message,
			colorStatus(c.Status),
			strconv.Itoa(c.TotalFailed),
			strconv.Itoa(c.TotalPassed),
			strconv.Itoa(c.TotalSkipped),
		})
	}
	t.Render()
}

func colorSeverity(s types.Severity) string {
	switch s {
	case types.SeverityLow:
		return blue("%s", s)
	case types.SeverityMedium:
		return yellow("%s", s)
	case types.SeverityHigh:
		return red("%s", s)
	case types.SeverityCritical:
		return redBold("%s", s)
	default:
		return s.String()
	}
}
func colorStatus(s types.CheckStatus) string {
	switch s {
	case types.StatusPassed:
		return green("%s", s)
	case types.StatusSkipped:
		return blue("%s", s)
	case types.StatusFailed:
		return red("%s", s)
	case types.StatusError:
		return redBold("%s", strings.ToUpper(s.String()))
	default:
		return s.String()
	}
}
