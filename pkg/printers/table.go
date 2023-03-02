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

	"github.com/fatih/color"
	"github.com/olekukonko/tablewriter"

	"github.com/undistro/marvin/pkg/checks"
	"github.com/undistro/marvin/pkg/report"
)

var (
	red     = color.New(color.FgRed).SprintfFunc()
	redBold = color.New(color.FgRed, color.Bold).SprintfFunc()
	yellow  = color.New(color.FgYellow).SprintfFunc()
	blue    = color.New(color.FgBlue).SprintfFunc()
	green   = color.New(color.FgGreen).SprintfFunc()
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
	t.SetHeader([]string{"SEVERITY", "ID", "CHECK", "STATUS", "FAILED", "PASSED", "SKIPPED"})
	sort.Slice(r.Checks, func(i, j int) bool {
		if r.Checks[i].Severity != r.Checks[j].Severity {
			return r.Checks[i].Severity > r.Checks[j].Severity
		}
		if r.Checks[i].TotalFailed != r.Checks[j].TotalFailed {
			return r.Checks[i].TotalFailed > r.Checks[j].TotalFailed
		}
		return r.Checks[i].TotalPassed > r.Checks[j].TotalPassed
	})
	for _, c := range r.Checks {
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

func colorSeverity(s checks.Severity) string {
	switch s {
	case checks.SeverityLow:
		return blue("%s", s)
	case checks.SeverityMedium:
		return yellow("%s", s)
	case checks.SeverityHigh:
		return red("%s", s)
	case checks.SeverityCritical:
		return redBold("%s", s)
	default:
		return s.String()
	}
}
func colorStatus(s report.CheckStatus) string {
	switch s {
	case report.StatusPassed:
		return green("%s", s)
	case report.StatusSkipped:
		return blue("%s", s)
	case report.StatusFailed:
		return red("%s", s)
	case report.StatusError:
		return redBold("%s", s)
	default:
		return s.String()
	}
}
