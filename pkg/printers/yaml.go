package printers

import (
	"fmt"
	"io"

	"sigs.k8s.io/yaml"

	"github.com/undistro/marvin/pkg/report"
)

type YAMLPrinter struct{}

func (*YAMLPrinter) PrintObj(r report.Report, w io.Writer) error {
	data, err := yaml.Marshal(r)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(w, string(data))
	return err
}
