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
	"fmt"
	"io"

	"sigs.k8s.io/yaml"

	"github.com/undistro/marvin/pkg/types"
)

// YAMLPrinter implements a Printer that prints the report in YAML format
type YAMLPrinter struct{}

func (*YAMLPrinter) PrintObj(report types.Report, w io.Writer) error {
	data, err := yaml.Marshal(report)
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(w, string(data))
	return err
}
