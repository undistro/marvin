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

package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/undistro/marvin/pkg/cmd"
)

var (
	scanOptions = cmd.NewScanOptions()

	// scanCmd represents the scan command
	scanCmd = &cobra.Command{
		Use:   "scan [flags]",
		Short: "Scan a Kubernetes cluster",
		Example: fmt.Sprintf(`  # Scan the current cluster
  %[1]s scan

  # Scan the 'foo' namespace of the current cluster
  %[1]s scan -n foo

  # Scan the current cluster providing custom checks
  %[1]s scan --checks ./examples/

  # Scan the current cluster providing custom checks and disabling the built-in checks
  %[1]s scan --disable-builtin --checks ./examples/

  # Scan a specific cluster using a kubeconfig file
  %[1]s scan --kubeconfig /path/to/kubeconfig.yml

  # Scan the current cluster, but do not fail even if there are errors in the report
  %[1]s scan --no-fail

  # Scan the current cluster and generate output in JSON format
  %[1]s scan -o json`, execName()),
		RunE: func(c *cobra.Command, args []string) error {
			if err := scanOptions.Init(c.Context()); err != nil {
				return err
			}
			hasError, err := scanOptions.Run()
			if err != nil {
				return err
			}
			if hasError && !*scanOptions.NoFail {
				os.Exit(2)
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(scanCmd)
	scanOptions.AddFlags(scanCmd.Flags())
}
