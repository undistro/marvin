/*
Copyright © 2023 Undistro Authors

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package cmd

import (
	"github.com/spf13/cobra"

	"github.com/undistro/marvin/pkg/cmd"
)

var (
	scanOptions = cmd.NewScanOptions()

	// scanCmd represents the scan command
	scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Scan a Kubernetes cluster",
		RunE: func(c *cobra.Command, args []string) error {
			if err := scanOptions.Run(); err != nil {
				return err
			}
			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(scanCmd)
	scanOptions.AddFlags(scanCmd.Flags())
}
