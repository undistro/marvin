package printers

import (
	"encoding/json"
	"io"

	"github.com/undistro/marvin/pkg/report"
)

type JSONPrinter struct{}

func (*JSONPrinter) PrintObj(r report.Report, w io.Writer) error {
	data, err := json.MarshalIndent(r, "", "    ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	_, err = w.Write(data)
	return err
}
