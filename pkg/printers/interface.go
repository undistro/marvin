package printers

import (
	"io"

	"github.com/undistro/marvin/pkg/report"
)

type Printer interface {
	PrintObj(report.Report, io.Writer) error
}
